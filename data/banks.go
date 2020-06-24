package data

import (
	"encoding/json"
	"fmt"
	"io"
)

type Bank struct {
	ID          int
	Name        string
	Description string
	USrate      float64
}

func (b *Bank) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(b)
}

type banks []*Bank

func GetBanks() banks {
	return banksList
}

func AddBank(b *Bank) {
	b.ID = getNextBankID()
	banksList = append(banksList, b)
}

func UpdateBank(id int, b *Bank) error {
	_, pos, err := findBank(id)
	if err != nil {
		return err
	}

	b.ID = id
	banksList[pos] = b
	return nil
}

var ErrBankNotFound = fmt.Errorf("Bank Not Found")

func findBank(id int) (*Bank, int, error) {
	for i, b := range banksList {
		if b.ID == id {
			return b, i, nil
		}
	}
	return nil, -1, ErrBankNotFound
}

func getNextBankID() int {
	lb := banksList[len(banksList)-1]
	return lb.ID + 1
}

func (b *banks) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(b)
}

var banksList = []*Bank{
	&Bank{
		ID:          1,
		Name:        "BIC",
		Description: "Test 1",
		USrate:      200,
	},
	&Bank{
		ID:          2,
		Name:        "BAI",
		Description: "Test 2",
		USrate:      300,
	},
}
