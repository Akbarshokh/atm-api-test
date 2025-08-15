package domain

import "time"

type Provider struct {
	ID        string
	Code      string //  UZCARD, HUMO
	Name      string
	CreatedAt time.Time
}

type CardStatus string

const (
	CardStatusActive CardStatus = "active"
	CardStatusClosed CardStatus = "closed"
	CardStatusStolen CardStatus = "stolen"
)

type Card struct {
	ID          string
	ProviderID  string
	Provider    *Provider
	PAN         string
	HashedPIN   string
	PINAttempts int
	IsBlocked   bool
	Balance     Money
	Currency    string
	Status      CardStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TxType string

const (
	TxWithdrawal   TxType = "withdrawal"
	TxDeposit      TxType = "deposit"
	TxBalanceCheck TxType = "balance_check"
	TxPINChange    TxType = "pin_change"
	TxBlock        TxType = "block"
	TxUnblock      TxType = "unblock"
)

type CardTx struct {
	ID        string
	CardID    string
	Type      TxType
	Amount    Money
	Fee       Money
	Success   bool
	Message   string
	Metadata  map[string]any
	CreatedAt time.Time
}

type DailyCounters struct {
	CardID       string
	Day          time.Time // date-only
	WithdrawnSum Money
	DepositSum   Money
	UpdatedAt    time.Time
}
