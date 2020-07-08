package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"

	"github.com/go-playground/validator"
)

type Company struct {
	ID      int     `json:"id"`
	Name    string  `json:"name" validate:"required"`
	Sector  string  `json:"sector" validate:"required,sectorVal"`
	Revenue float64 `json:"revenue"`
}

func (c *Company) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(c)
}

func (c *Company) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("sectorVal", sectorVal)
	return validate.Struct(c)
}

func sectorVal(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`[a-z]+-[0-9]+`)
	matches := re.FindAllString(fl.Field().String(), -1)

	if len(matches) != 1 {
		return false
	}

	return true
}

type companies []*Company

func GetCompany() companies {
	return companyList
}

func (c *companies) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(c)
}

func AddCompany(c *Company) {
	c.ID = getNextCompanyID()
	companyList = append(companyList, c)
}

func getNextCompanyID() int {
	lc := companyList[len(companyList)-1]
	return lc.ID + 1
}

func UpdateCompany(id int, c *Company) error {
	_, pos, err := findCompany(id)
	if err != nil {
		return err
	}

	c.ID = id
	companyList[pos] = c
	return nil

}

var ErrCompanyNotFound = fmt.Errorf("Company Not Found")

func findCompany(id int) (*Company, int, error) {
	for i, c := range companyList {
		if c.ID == id {
			return c, i, nil
		}
	}

	return nil, -1, ErrCompanyNotFound
}

var companyList = []*Company{
	&Company{
		ID:      1,
		Name:    "EkuyTrade",
		Sector:  "Financial Services",
		Revenue: 2000000.00,
	},
	&Company{
		ID:      2,
		Name:    "FKitas Mgmt",
		Sector:  "Consulting",
		Revenue: 1000000.00,
	},
}
