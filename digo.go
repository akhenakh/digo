package digo

import (
	"encoding/json"
	"errors"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const EMAIL_REG = `(\w[-._\w]*\w@\w[-._\w]*\w\.\w{2,3})`
const MINMAX_REG = `^minmax\(\s*(\d+)\s*\|\s*(\d+)\s*\)`

var (
	emailRegex  *regexp.Regexp
	minMaxRegex *regexp.Regexp
)

func init() {
	log.SetFlags(log.Lshortfile | log.Ltime)
}

func emailField(email string, fieldName string) error {
	if emailRegex == nil {
		emailRegex = regexp.MustCompile(EMAIL_REG)
	}
	if !emailRegex.MatchString(email) {
		return errors.New("field: " + fieldName + " is not a valid email address")
	}
	return nil
}

func minMax(field interface{}, fieldName string, dstType string, min, max int) error {

	switch dstType {
	case "int":
		if int(field.(float64)) < min {
			return errors.New("field: " + fieldName + " is too small")
		}
		if int(field.(float64)) > max {
			return errors.New("field: " + fieldName + " is too big")
		}
	case "string":
		if reflect.ValueOf(field).Len() < min {
			return errors.New("field: " + fieldName + " is too short")
		}
		if reflect.ValueOf(field).Len() > max {
			return errors.New("field: " + fieldName + " is too long")
		}

	}

	return nil
}

func UnmarshalJSON(data []byte, dst interface{}) error {
	var in interface{}

	err := json.Unmarshal(data, &in)

	// map for in
	im := in.(map[string]interface{})

	if err != nil {
		return err
	}

	// analyze the target struct
	val := reflect.ValueOf(dst)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return errors.New("digo: interface must be a pointer to struct")
	}
	val = val.Elem()

	for i := 0; i < val.NumField(); i++ {
		vf := val.Field(i)
		tf := val.Type().Field(i)
		tag := tf.Tag

		// no tag for us no decoding
		if tag.Get("digo") == "" {
			continue
		}
		stags := strings.Split(tag.Get("digo"), ",")

		// nothing to validate
		if len(stags) < 2 {
			return errors.New("field: " + tf.Name + " is incorrect")
		}
		fieldName := strings.TrimSpace(stags[0])
		fieldType := strings.TrimSpace(stags[1])

		// log.Printf("SRCfield value: %s, SRCfield type :%T\nDSTField Name: %s, DSTField Type: %s Tag Value: %s",
		// 	im[fieldName],
		// 	im[fieldName],
		// 	tf.Name,
		// 	tf.Type,
		// 	tag.Get("digo"))

		// we need to manually set the dst type at least for numbers (defaulting to float64)
		var dstType string
		// realvalue in the decoded json
		rv, _ := im[fieldName]
		if rv != nil {
			if strings.HasPrefix(fieldType, "emailfield") {
				dstType = "string"
				// check the src type
				if reflect.TypeOf(rv).String() != "string" {
					return errors.New("field: " + fieldName + " is not a string")
				}
				// check the dst type
				if tf.Type.String() != "string" {
					return errors.New("struct: " + tf.Name + " is not a string")
				}
				if err := emailField(rv.(string), fieldName); err != nil {
					return err
				}

			} else if strings.HasPrefix(fieldType, "stringfield") {
				dstType = "string"
				// check the src type
				if reflect.TypeOf(rv).String() != "string" {
					return errors.New("field: " + fieldName + " is not a string")
				}
				// check the dst type
				if tf.Type.String() != "string" {
					return errors.New("struct: " + tf.Name + " is not a string")
				}
			} else if strings.HasPrefix(fieldType, "intfield") {
				dstType = "int"
				// check the src type
				// all numbers default to float 64
				if reflect.TypeOf(rv).String() != "float64" {
					return errors.New("field: " + fieldName + " is not an integer")
				}
				// check the dst type
				if tf.Type.String() != "int" {
					return errors.New("struct: " + tf.Name + " is not an integer")
				}
			} else {
				return errors.New("field: " + fieldName + " invalid digo type")
			}

			// set the target value
			switch dstType {
			case "string":
				vf.SetString(rv.(string))
			case "int":
				vf.SetInt(int64(rv.(float64)))
			}
		}

		// check the target value is a same type
		for _, stag := range stags[2:] {
			stag = strings.TrimSpace(stag)

			// tag required
			if strings.HasPrefix(stag, "required") {
				if rv == nil {
					return errors.New("field: " + fieldName + " is required")
				}
				// for each type apply revelant test
				if tf.Type.String() == "string" {
					if rv.(string) == "" {
						return errors.New("field: " + fieldName + " is required")
					}
				}
			} else if strings.HasPrefix(stag, "minmax") {
				if minMaxRegex == nil {
					minMaxRegex = regexp.MustCompile(MINMAX_REG)
				}
				res := minMaxRegex.FindStringSubmatch(stag)
				if len(res) != 3 {
					return errors.New("field: " + fieldName + " minmax invalid call")
				}
				// no need to read err, regexp already check that for us
				min, _ := strconv.Atoi(res[1])
				max, _ := strconv.Atoi(res[2])

				if err = minMax(rv, fieldName, dstType, min, max); err != nil {
					return err
				}

			}

		}

	}

	return nil
}
