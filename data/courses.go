package data

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"regexp"

	"github.com/go-playground/validator/v10"
)

type Course struct {
	ID     int    `json:"id"`
	Number string `json:"number" validate:"required,nameVal"`
	Name   string `json:"name" validate:"required"`
	Units  string `json:"units" validate:"required"`
}

func (c *Course) Validate() error {

	validate := validator.New()
	validate.RegisterValidation("nameVal", nameVal)

	return validate.Struct(c)

}

func nameVal(fl validator.FieldLevel) bool {
	reg := regexp.MustCompile(`[a-z]+-[a-z]+`)
	matches := reg.FindAllString(fl.Field().String(), -1)
	if len(matches) != 1 {
		return false
	}
	return true
}

func (c *Course) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(c)
}

type courses []*Course

func GetCourses() courses {
	return courseList
}

func AddCourse(c *Course) {
	c.ID = getNextCourseID()
	courseList = append(courseList, c)
}

func UpdateCourse(id int, c *Course) error {
	_, pos, err := findCourse(id)
	if err != nil {
		log.Println(err)
	}

	c.ID = id
	courseList[pos] = c
	return nil
}

var ErrCourseNotFound = fmt.Errorf("Course Not Found")

func findCourse(id int) (*Course, int, error) {
	for i, c := range courseList {
		if c.ID == id {
			return c, i, nil
		}
	}

	return nil, -1, ErrCourseNotFound
}

func getNextCourseID() int {
	lc := courseList[len(courseList)-1]
	return lc.ID + 1
}

func (c *courses) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(c)
}

var courseList = []*Course{
	&Course{
		ID:     1,
		Number: "abc-0123",
		Name:   "Test Course",
		Units:  "4",
	},
	&Course{
		ID:     2,
		Number: "abcd-0233",
		Name:   "2nd Test Course",
		Units:  "5",
	},
}
