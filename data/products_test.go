package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name:  "BK",
		Price: 1.00,
		SKU:   "abc-abs-abs",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
