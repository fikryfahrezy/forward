package bankstatement

type TransactionType string

const (
	TransactionTypeDebit   TransactionType = "DEBIT"
	TransactionTypeCredit  TransactionType = "CREDIT"
	TransactionTypeUnknown TransactionType = "UNKNOWN"
)

func (t TransactionType) Valid() bool {
	return t == TransactionTypeCredit || t == TransactionTypeDebit
}

type TransactionStatus string

const (
	TransactionStatusSuccess TransactionStatus = "SUCCESS"
	TransactionStatusPending TransactionStatus = "PENDING"
	TransactionStatusFailed  TransactionStatus = "FAILED"
	TransactionStatusUnknown TransactionStatus = "UNKNOWN"
)

func (t TransactionStatus) Valid() bool {
	return t == TransactionStatusSuccess || t == TransactionStatusPending || t == TransactionStatusFailed
}

type Transaction struct {
	Timestamp   int64
	Name        string
	Type        TransactionType
	Amount      float64
	Status      TransactionStatus
	Description string
}
