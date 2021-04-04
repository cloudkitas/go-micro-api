package handlers

import (
	"context"
	"fmt"
	"go-api/practice/data"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Courses struct {
	l *log.Logger
}

func NewCourses(l *log.Logger) *Courses {
	return &Courses{l}
}

func (c *Courses) GetCourses(rw http.ResponseWriter, r *http.Request) {
	c.l.Println("Handle GET Courses")
	cl := data.GetCourses()
	err := cl.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to Marshal json", http.StatusBadRequest)
	}
}

func (c *Courses) AddCourse(rw http.ResponseWriter, r *http.Request) {
	c.l.Println("Handle POST Course")
	cors := r.Context().Value(keyCourse{}).(data.Course)
	data.AddCourse(&cors)
}

func (c Courses) UpdateCourse(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to Convert ID", http.StatusBadRequest)
		return
	}

	c.l.Println("Handle PUT Course")
	cors := r.Context().Value(keyCourse{}).(data.Course)
	err = data.UpdateCourse(id, &cors)

	if err == data.ErrCourseNotFound {
		http.Error(rw, "cant find course", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Cant Find course dam ", http.StatusInternalServerError)
		return
	}

}

type keyCourse struct{}

func (c Courses) MiddleWareValidateCourse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		cors := data.Course{}
		err := cors.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}

		//validate couse

		err = cors.Validate()
		if err != nil {
			http.Error(
				rw,
				fmt.Sprintf("Error validating course: %s", err),
				http.StatusBadRequest,
			)
		}

		ctx := context.WithValue(r.Context(), keyCourse{}, cors)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}
