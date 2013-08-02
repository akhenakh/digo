package digo

import (
	"testing"
)

func Test_EmailField(t *testing.T) {
	assertEqual := func(val interface{}, exp interface{}) {
		if val != exp {
			t.Errorf("Expected %v, got %v.", exp, val)
		}
	}
	type User struct {
		Email string `digo:"email,emailfield()"`
	}

	b := []byte(`{"email":"toto"}`)
	var u User
	err := UnmarshalJSON(b, &u)
	assertEqual(err.Error(), "field: email is not a valid email address")
	b = []byte(`{"email":"toto@gmail.com"}`)
	if err = UnmarshalJSON(b, &u); err != nil {
		t.Error("should be a valid email")
	}

	b = []byte(`{"email":2}`)
	err = UnmarshalJSON(b, &u)
	assertEqual(err.Error(), "field: email is not a string")

	type BadUser struct {
		Email int `digo:"email,emailfield()"`
	}
	var bu BadUser
	b = []byte(`{"email":"akh@titi.fr"}`)
	err = UnmarshalJSON(b, &bu)
	assertEqual(err.Error(), "struct: Email is not a string")
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
	assertEqual := func(val interface{}, exp interface{}) {
		if val != exp {
			t.Errorf("Expected %v, got %v.", exp, val)
		}
	}
	type User struct {
		Name string `digo:"name,stringfield(),required()"`
	}

	b := []byte(`{"error":"toto"}`)
	var u User
	err := UnmarshalJSON(b, &u)

	assertEqual(err.Error(), "field: name is required")
	b = []byte(`{"name":""}`)
	err = UnmarshalJSON(b, &u)

	assertEqual(err.Error(), "field: name is required")
}

func Test_Intfield(t *testing.T) {
	assertEqual := func(val interface{}, exp interface{}) {
		if val != exp {
			t.Errorf("Expected %v, got %v.", exp, val)
		}
	}
	type UserInt struct {
		Uid int `digo:"uid,intfield()"`
	}

	var ui UserInt
	b := []byte(`{"uid":"toto"}`)
	err := UnmarshalJSON(b, &ui)
	assertEqual(err.Error(), "field: uid is not an integer")

	type User struct {
		Uid string `digo:"uid,intfield()"`
	}
	var u User
	b = []byte(`{"uid":3}`)
	err = UnmarshalJSON(b, &u)
	assertEqual(err.Error(), "struct: Uid is not an integer")
}

func Test_MinMax(t *testing.T) {
	assertEqual := func(val interface{}, exp interface{}) {
		if val != exp {
			t.Errorf("Expected %v, got %v.", exp, val)
		}
	}

	type User struct {
		Name string `digo:"name,stringfield(),minmax(2|5)"`
	}

	b := []byte(`{"name":"toto"}`)
	var u User
	err := UnmarshalJSON(b, &u)
	if err != nil {
		t.Error("should not detect an error")
	}

	type UserInt struct {
		Uid int `digo:"uid,intfield(),minmax(2|5)"`
	}
	var ui UserInt
	b = []byte(`{"uid":4}`)
	err = UnmarshalJSON(b, &ui)
	if err != nil {
		t.Error("should not detect an error")
	}

	b = []byte(`{"uid":7}`)
	err = UnmarshalJSON(b, &ui)
	assertEqual(err.Error(), "field: uid is too big")
}
