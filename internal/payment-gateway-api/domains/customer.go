package domains

type Customer struct {
	Id int64
}

func (customer *Customer) IsInvalid() bool {
	return customer == nil || customer.Id <= 0
}
