package repo

import (
	"atm-test/helpers"
	"atm-test/internal/domain"
	"database/sql"
	"encoding/json"
	"time"
)

type providerModel struct {
	ID        string    `db:"id"`
	Code      string    `db:"code"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

func (m providerModel) toDomain() domain.Provider {
	return domain.Provider{
		ID:        m.ID,
		Code:      m.Code,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
	}
}

type cardModel struct {
	ID          string    `db:"id"`
	ProviderID  string    `db:"provider_id"`
	PAN         string    `db:"pan"`
	HashedPIN   string    `db:"hashed_pin"`
	PINAttempts int       `db:"pin_attempts"`
	IsBlocked   bool      `db:"is_blocked"`
	BalanceStr  string    `db:"balance"`
	Currency    string    `db:"currency_code"`
	Status      string    `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (m cardModel) toDomain() (domain.Card, error) {
	tiyin, err := helpers.NumericStringToTiyin(m.BalanceStr)
	if err != nil {
		return domain.Card{}, err
	}
	return domain.Card{
		ID:          m.ID,
		ProviderID:  m.ProviderID,
		PAN:         m.PAN,
		HashedPIN:   m.HashedPIN,
		PINAttempts: m.PINAttempts,
		IsBlocked:   m.IsBlocked,
		Balance:     domain.Money(tiyin),
		Currency:    m.Currency,
		Status:      domain.CardStatus(m.Status),
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}, nil
}

// transactionModel — структура, соответствующая строке из card_transactions
type transactionModel struct {
	ID        string          `db:"id"`
	CardID    string          `db:"card_id"`
	Type      string          `db:"type"`
	AmountStr string          `db:"amount"`
	FeeStr    string          `db:"fee"`
	Success   bool            `db:"success"`
	Message   sql.NullString  `db:"message"`
	Metadata  json.RawMessage `db:"metadata"`
	CreatedAt time.Time       `db:"created_at"`
}

func (m transactionModel) toDomain() (domain.CardTx, error) {
	amt, err := helpers.NumericStringToTiyin(m.AmountStr)
	if err != nil {
		return domain.CardTx{}, err
	}
	fee, err := helpers.NumericStringToTiyin(m.FeeStr)
	if err != nil {
		return domain.CardTx{}, err
	}
	msg := ""
	if m.Message.Valid {
		msg = m.Message.String
	}
	var metadata map[string]any
	_ = json.Unmarshal(m.Metadata, &metadata)
	return domain.CardTx{
		ID:        m.ID,
		CardID:    m.CardID,
		Type:      domain.TxType(m.Type),
		Amount:    domain.Money(amt),
		Fee:       domain.Money(fee),
		Success:   m.Success,
		Message:   msg,
		Metadata:  metadata,
		CreatedAt: m.CreatedAt,
	}, nil
}
