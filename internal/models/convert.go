package models

func NewWithdrawalFromDB(db WithdrawalDB) Withdrawal {
	return Withdrawal{
		Order:       db.Order,
		Sum:         NewMoney(db.Amount),
		ProcessedAt: db.ProcessedAt,
	}
}

func NewWithdrawalsFromDB(db []WithdrawalDB) []Withdrawal {
	res := make([]Withdrawal, 0, len(db))
	for _, v := range db {
		res = append(res, NewWithdrawalFromDB(v))
	}
	return res
}

func NewOrderFromDB(db OrderDB) Order {
	return Order{
		Number:     db.Number,
		Status:     db.Status,
		Accrual:    NewMoney(db.Accrual),
		UploadedAt: db.UploadedAt,
	}
}

func NewOrdersFromDB(db []OrderDB) []Order {
	res := make([]Order, 0, len(db))
	for _, v := range db {
		res = append(res, NewOrderFromDB(v))
	}
	return res
}
