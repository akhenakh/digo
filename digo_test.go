package digo

import (
	"testing"
)

func Test_EmailField(t *testing.T) {
	type User struct {
		Email string `digo:"email,emailfield()"`
	}

	b := []byte(`{"Email":"toto"}`)
	var u User
	err := UnmarshalJSON(b, &u)
	if err.Error() != "field: email is not a valid email address" {
		t.Error("email test failed")
	}
	b = []byte(`{"Email":"toto@gmail.com"}`)
	if err = UnmarshalJSON(b, &u); err != nil {
		t.Error("should be a valid email")
	}

}

func Test_NoDigoType(t *testing.T) {
	type User struct {
		Name string `digo:"name, required()"`
	}

	b := []byte(`{"name":"toto"}`)
	var u User
	err := UnmarshalJSON(b, &u)
	if err.Error() != "field: name invalid digo type" {
		t.Error("required field not detected")
	}
}

func Test_Required(t *testing.T) {
	type User struct {
		Name string `digo:"name,stringfield(), required()"`
	}

	b := []byte(`{"error":"toto"}`)
	var u User
	err := UnmarshalJSON(b, &u)
	if err.Error() != "field: name is required" {
		t.Error("required field not detected")
	}
}

func Test_MinMax(t *testing.T) {
	type User struct {
		Name string `digo:"name,stringfield(), minmax(2|5)"`
	}

	b := []byte(`{"name":"toto"}`)
	var u User
	err := UnmarshalJSON(b, &u)
	if err != nil {
		t.Error("should not detect an error")
	}
}
