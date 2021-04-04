package data

import (
	"log"
	"testing"
)

func TestPeopleValidation(t *testing.T) {
	p := &People{
		ID:       1,
		Name:     "test",
		Age:      27,
		YearBorn: 1992,
	}

	err := p.Validate()
	if err != nil {
		log.Fatalln(err)
	}
}
