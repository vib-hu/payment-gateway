package enums

type ErrorType int

const (
	ValidationError     ErrorType = 1
	IntegrationError    ErrorType = 2
	InternalServerError ErrorType = 3
	PaymentError        ErrorType = 4
)
