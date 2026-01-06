package models

type AccrualOrderStatus string

const (
	AccrualStatusRegistered AccrualOrderStatus = "REGISTERED"
	AccrualStatusProcessing AccrualOrderStatus = "PROCESSING"
	AccrualStatusInvalid    AccrualOrderStatus = "INVALID"
	AccrualStatusProcessed  AccrualOrderStatus = "PROCESSED"
)

type GetOrderAccrualResponse struct {
	Order   string
	Status  AccrualOrderStatus
	Accrual float64
}
