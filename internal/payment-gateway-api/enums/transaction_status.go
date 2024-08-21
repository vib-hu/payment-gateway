package enums

import "fmt"

type TransactionStatus int

const (
	TransactionInitiated  TransactionStatus = 1
	TransactionSuccessful TransactionStatus = 2
	TransactionFailed     TransactionStatus = 3
	TransactionRefunded   TransactionStatus = 4
)

func ParseTransactionStatus(value int) (TransactionStatus, error) {
	switch value {
	case int(TransactionInitiated):
		return TransactionInitiated, nil
	case int(TransactionSuccessful):
		return TransactionSuccessful, nil
	case int(TransactionFailed):
		return TransactionFailed, nil
	case int(TransactionRefunded):
		return TransactionRefunded, nil
	default:
		return 0, fmt.Errorf("invalid transaction status: %d", value)
	}
}
