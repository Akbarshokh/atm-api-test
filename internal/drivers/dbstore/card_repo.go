package repo

import (
	"atm-test/internal/errs"
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"atm-test/helpers"
	"atm-test/internal/domain"
	"atm-test/internal/pkg/logger"
)

type Repo struct {
	db  *sql.DB
	log logger.Logger
}

func New(db *sql.DB, log logger.Logger) *Repo {
	return &Repo{
		db:  db,
		log: log,
	}
}

// Получить провайдера по коду
func (r *Repo) GetProviderByCode(ctx context.Context, code string) (domain.Provider, error) {
	var (
		logMsg = "repo.GetProviderByCode "
		query  = `SELECT id, code, name, created_at FROM card_providers WHERE code=$1`
		data   providerModel
	)

	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&data.ID, &data.Code, &data.Name, &data.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.log.Error(logMsg+"provider not found", logger.Error(err))
			return domain.Provider{}, errs.ErrProviderNotAllowed
		}
		r.log.Error(logMsg+"r.db.QueryRowContext failed", logger.Error(err))
		return domain.Provider{}, err
	}

	return data.toDomain(), nil
}

// Получить карту без блокировки (для аутентификации/чтения)
func (r *Repo) GetCardByPAN(ctx context.Context, pan string) (domain.Card, error) {
	var (
		logMsg = "repo.GetCardByPAN "
		query  = `
		SELECT id, provider_id, pan, hashed_pin, pin_attempts, is_blocked,
		       balance::text, currency_code, status, created_at, updated_at
		FROM cards WHERE pan=$1`
		data cardModel
	)

	err := r.db.QueryRowContext(ctx, query, pan).Scan(
		&data.ID, &data.ProviderID, &data.PAN, &data.HashedPIN,
		&data.PINAttempts, &data.IsBlocked, &data.BalanceStr, &data.Currency,
		&data.Status, &data.CreatedAt, &data.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.log.Error(logMsg+"card not found", logger.Error(err))
			return domain.Card{}, errs.ErrCardNotFound
		}
		r.log.Error(logMsg+"r.db.QueryRowContext failed", logger.Error(err))
		return domain.Card{}, err
	}

	return data.toDomain()
}

// Заблокировать карту
func (r *Repo) BlockCard(ctx context.Context, cardID string) error {
	var (
		logMsg = "repo.BlockCard "
		query  = `UPDATE cards SET is_blocked=TRUE, updated_at=now() WHERE id=$1`
	)

	if _, err := r.db.ExecContext(ctx, query, cardID); err != nil {
		r.log.Error(logMsg+"r.db.ExecContext failed", logger.Error(err))
		return err
	}
	return nil
}

// Обновить баланс карты на дельту (delta может быть отрицательной при снятии)
func (r *Repo) UpdateCardBalance(ctx context.Context, cardID string, delta domain.Money) error {
	var (
		logMsg = "repo.UpdateCardBalance "
		query  = `UPDATE cards SET balance = balance + $2::numeric(14,2), updated_at=now() WHERE id=$1`
	)

	if _, err := r.db.ExecContext(ctx, query, cardID, helpers.TiyinToNumericString(delta.Tiyin())); err != nil {
		r.log.Error(logMsg+"r.db.ExecContext failed", logger.Error(err))
		return err
	}
	return nil
}

// Добавить транзакцию
func (r *Repo) AddTransaction(ctx context.Context, tx domain.CardTx) error {
	var (
		logMsg = "repo.AddTransaction "
		query  = `
		INSERT INTO card_transactions (card_id, type, amount, fee, success, message, metadata)
		VALUES ($1, $2, $3::numeric(14,2), $4::numeric(14,2), $5, $6, $7)`
	)

	metaBytes, _ := json.Marshal(tx.Metadata)
	var msg sql.NullString
	if tx.Message != "" {
		msg.Valid = true
		msg.String = tx.Message
	}

	if _, err := r.db.ExecContext(ctx, query,
		tx.CardID, string(tx.Type),
		helpers.TiyinToNumericString(tx.Amount.Tiyin()),
		helpers.TiyinToNumericString(tx.Fee.Tiyin()),
		tx.Success, msg, metaBytes,
	); err != nil {
		r.log.Error(logMsg+"r.db.ExecContext failed", logger.Error(err))
		return err
	}
	return nil
}
