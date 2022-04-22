// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

type APIStatus struct {
	Name   string `json:"name"`
	Ref    string `json:"ref"`
	Commit string `json:"commit"`
}

type NewRecipe struct {
	Name string `json:"name"`
}

type Recipe struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Planned      bool       `json:"planned"`
	LastCookedAt *time.Time `json:"lastCookedAt"`
}

type Result struct {
	Status Status `json:"status"`
}

type UpdateRecipe struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Status string

const (
	StatusSuccess  Status = "Success"
	StatusError    Status = "Error"
	StatusNotFound Status = "NotFound"
)

var AllStatus = []Status{
	StatusSuccess,
	StatusError,
	StatusNotFound,
}

func (e Status) IsValid() bool {
	switch e {
	case StatusSuccess, StatusError, StatusNotFound:
		return true
	}
	return false
}

func (e Status) String() string {
	return string(e)
}

func (e *Status) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Status(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Status", str)
	}
	return nil
}

func (e Status) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
