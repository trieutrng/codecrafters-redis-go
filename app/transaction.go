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

func (tx *Transaction) IsActive(txId string) bool {
	_, ok := tx.Active[txId]
	return ok
}

func (tx *Transaction) Inactive(txId string) {
	_, ok := tx.Active[txId]
	if ok {
		delete(tx.Active, txId)
	}
}

func (tx *Transaction) GetTx(txId string) []*RESP {
	queued, ok := tx.Active[txId]
	if !ok {
		return []*RESP{}
	}
	return queued
}
