package database

type Bank struct {
	ID   string
	Name string
}

type Statement struct {
	ID          string
	Bank        string
	Amount      int64
	Description string
	Asset       string
	Qtd         int16
}
