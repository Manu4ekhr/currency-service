package entity

type OperationStatus string

const (
	StatusPending    OperationStatus = "pending"
	StatusProcessing OperationStatus = "processing"
	StatusSended     OperationStatus = "sended"
)

type OperationType string

const (
	TypeTransfers OperationType = "Transfers"
	TypeCashDesks OperationType = "Cash Desks"
	TypeNonCash   OperationType = "Non-cash"
	TypeNBT       OperationType = "NBT"
)

type CurrencyOperation struct {
	Currency string          `json:"currency" gorm:"column:currency"`
	Rate     float64         `json:"rate" gorm:"column:rate"`
	ID       int64           `json:"id" gorm:"primaryKey;autoIncrement"`
	From     string          `json:"from" gorm:"column:from"`
	To       string          `json:"to" gorm:"column:to"`
	ExtID    string          `json:"ext_id" gorm:"column:ext_id;uniqueIndex"`
	Type     OperationType   `json:"type" gorm:"column:type"`
	Status   OperationStatus `json:"status" gorm:"column:status"`
}

func (CurrencyOperation) TableName() string {
	return "currency_operations"
}
