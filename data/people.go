package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"

	"github.com/go-playground/validator/v10"
)

type People struct {
	ID       int    `json:"id"`
	Name     string `json:"name" validate:"required,valName"`
	Age      int    `json:"age" validate:"required"`
	YearBorn int    `json:"yearborn"`
}

func (p *People) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

func (p *People) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("valName", valName)
	return validate.Struct(p)
}

func valName(fl validator.FieldLevel) bool {
	reg := regexp.MustCompile(`[a-z]+-[a-z]+`)
	matches := reg.FindAllString(fl.Field().String(), -1)

	if len(matches) != 1 {
		return false
	}
	return true
}

type peoples []*People

func GetPeople() peoples {
	return peopleList
}

func AddPeople(p *People) {
	p.ID = getNextPeopleID()
	peopleList = append(peopleList, p)
}

func getNextPeopleID() int {
	lp := peopleList[len(peopleList)-1]
	return lp.ID + 1
}

func UpdatePeople(id int, p *People) error {
	_, pos, err := findPeople(id)
	if err != nil {
		return err
	}

	p.ID = id
	peopleList[pos] = p
	return nil
}

var ErrPeopleNotFound = fmt.Errorf("People NoT Found")

func findPeople(id int) (*People, int, error) {
	for i, p := range peopleList {
		if p.ID == id {
			return p, i, nil
		}
	}
	return nil, -1, ErrPeopleNotFound
}

func (p *peoples) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

var peopleList = []*People{
	&People{
		ID:       1,
		Name:     "Bruno",
		Age:      28,
		YearBorn: 1992,
	},
	&People{
		ID:       2,
		Name:     "Shiara",
		Age:      28,
		YearBorn: 1992,
	},
	&People{
		ID:       3,
		Name:     "Ayanni",
		Age:      02,
		YearBorn: 2018,
	},
}
