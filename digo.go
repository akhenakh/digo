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

func emailField(email string, fieldName string) error {
	if emailRegex == nil {
		emailRegex = regexp.MustCompile(EMAIL_REG)
	}
	if !emailRegex.MatchString(email) {
		return errors.New("field: " + fieldName + " is not a valid email address")
	}
	return nil
}

func minMax(field interface{}, fieldName string, min, max int) error {
	v := reflect.ValueOf(field)

	//if v.Type().Name() != "string" ...
	// this is useless pass until we are scanning the payload ourself

	if v.Len() < min {
		return errors.New("field: " + fieldName + " is too short")
	}
	if v.Len() > max {
		return errors.New("field: " + fieldName + " is too long")
	}
	return nil
}

func UnmarshalJSON(data []byte, dst interface{}) error {
	err := json.Unmarshal(data, dst)

	if err != nil {
		return err
	}

	val := reflect.ValueOf(dst)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return errors.New("digo: interface must be a pointer to struct")
	}
	val = val.Elem()

	// TODO we should parse the dst interface first then ask for json ...

	for i := 0; i < val.NumField(); i++ {
		vf := val.Field(i)
		tf := val.Type().Field(i)
		tag := tf.Tag

		log.Printf("Field Name: %s,\t Field Type: %s Field Value: %v,\t Tag Value: %s\n",
			tf.Name,
			tf.Type,
			vf.Interface(),
			tag.Get("digo"))

		stags := strings.Split(tag.Get("digo"), ",")
		// nothing to validate
		if len(stags) < 2 {
			return errors.New("field: " + tf.Name + " is incorrect")
		}
		fieldName := strings.TrimSpace(stags[0])
		fieldType := strings.TrimSpace(stags[1])

		if strings.HasPrefix(fieldType, "emailfield") {
			if err := emailField(vf.Interface().(string), fieldName); err != nil {
				return err
			}
		} else if strings.HasPrefix(fieldType, "stringfield") {
			//TODO: ensure field type
		} else {
			return errors.New("field: " + fieldName + " invalid digo type")
		}

		if len(stags) < 3 {
			continue
		}
		log.Println(stags)

		for _, stag := range stags[2:] {
			stag = strings.TrimSpace(stag)
			if strings.HasPrefix(stag, "required") {
				log.Println("LLAAA")
				if vf.Interface().(string) == "" {
					return errors.New("field: " + fieldName + " is required")
				}
			} else if strings.HasPrefix(stag, "minmax") {
				if minMaxRegex == nil {
					minMaxRegex = regexp.MustCompile(MINMAX_REG)
				}
				res := minMaxRegex.FindStringSubmatch(stag)
				log.Println(res)
				if len(res) != 3 {
					return errors.New("field: " + fieldName + " minmax invalid call")
				}
				min, _ := strconv.Atoi(res[1])
				max, _ := strconv.Atoi(res[2])
				if err = minMax(vf.Interface(), fieldName, min, max); err != nil {
					return err
				}

			}

		}

		log.Println(stags)
	}

	return nil
}
