package main

type Transaction struct {
	Active map[string][]*RESP
}

func NewTransaction() *Transaction {
	return &Transaction{
		Active: make(map[string][]*RESP),
	}
}

func (tx *Transaction) Start(txId string) {
	tx.Active[txId] = make([]*RESP, 0)
}

func (tx *Transaction) isActive(txId string) bool {
	_, ok := tx.Active[txId]
	return ok
}
