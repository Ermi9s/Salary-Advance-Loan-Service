package entity

import "time"

type BankCustomer struct {
	ID              int64   `json:"id"`
	CustomerName    string  `json:"customerName"`
	Mobile          string  `json:"mobile"`
	AccountNo       string  `json:"accountNo"`
	BranchName      string  `json:"branchName"`
	ProductName     string  `json:"productName"`
	CustomerID      string  `json:"customerId"`
	BranchCode      string  `json:"branchCode"`
	CustomerBalance float64 `json:"customerBalance"`
}

type InputCustomer struct {
	CustomerName string `json:"customerName"`
	AccountNo    string `json:"accountNo"`
}

type Transaction struct {
	ID                  string  `json:"id"`
	FromAccount         string  `json:"fromAccount"`
	ToAccount           *string `json:"toAccount"`
	Amount              string  `json:"amount"`
	Remark              string  `json:"remark"`
	TransactionType     string  `json:"transactionType"`
	RequestID           string  `json:"requestId"`
	Reference           string  `json:"reference"`
	ThirdPartyReference *string `json:"thirdPartyReference"`
	InstitutionID       *string `json:"institutionId"`
	ClearedBalance      *string `json:"clearedBalance"`
	TransactionDate     string  `json:"transactionDate"`
	BillerID            *string `json:"billerId"`
}

type EnrichedTransaction struct {
	ID              string    `json:"id"`
	AccountNo       string    `json:"account_no"`
	Amount          float64   `json:"amount"`
	Direction       string    `json:"direction"`
	TransactionDate time.Time `json:"transaction_date"`
	Source          string    `json:"source"`
}
