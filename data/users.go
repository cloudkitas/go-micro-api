package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"

	"github.com/go-playground/validator/v10"
)

type User struct {
	ID        int    `json:"id" validate:"required"`
	Username  string `json:"username" validate:"required,valusername"`
	AdminUser bool
}

func (u *User) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(u)
}

func (u *User) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("valusername", valUser)

	return validate.Struct(u)
}

func valUser(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`[a-z]+-[a-z]+`)
	matches := re.FindAllString(fl.Field().String(), -1)

	if len(matches) != 1 {
		return false
	}
	return true
}

type users []*User

func GetUsers() users {
	return usersList
}

func AddUsers(u *User) {
	u.ID = getUserNextID()
	usersList = append(usersList, u)
}

func UpdateUser(id int, u *User) error {
	_, pos, err := findUser(id)
	if err != nil {
		return err
	}

	u.ID = id
	usersList[pos] = u
	return nil
}

var ErrUserNotFound = fmt.Errorf("User Not Found")

func findUser(id int) (*User, int, error) {
	for i, u := range usersList {
		if u.ID == id {
			return u, i, nil
		}
	}

	return nil, -1, ErrUserNotFound
}

func getUserNextID() int {
	lu := usersList[len(usersList)-1]
	return lu.ID + 1
}

func (u *users) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(u)
}

var usersList = []*User{
	&User{
		ID:        1,
		Username:  "brunokitas",
		AdminUser: true,
	},
	&User{
		ID:        2,
		Username:  "krabbe",
		AdminUser: true,
	},
	&User{
		ID:        3,
		Username:  "Anni",
		AdminUser: false,
	},
}
