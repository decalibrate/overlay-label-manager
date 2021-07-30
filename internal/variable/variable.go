package variable

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/decalibrate/overlay-label-manager/internal/configuration"
	"github.com/decalibrate/overlay-label-manager/internal/helper"
)

var Variables []Variable = make([]Variable, 0)
var Cfg *configuration.ConfigStruct

var reVariableName = regexp.MustCompile("^[a-zA-Z_][0-9a-zA-Z_]{0,19}$")

type Variable interface {
	Id() string
	Type() string
	GetBindings(*[]string) (string, interface{})
	Unmarshal(map[string]interface{}) error
	Set(string, map[string][]string) error
}

func Create(t string) (Variable, error) {
	switch t {
	case "counter":
		v := &Counter{VarType: t}
		return v, nil
	case "text":
		v := &Text{VarType: t}
		return v, nil
	}

	return nil, fmt.Errorf("unknown variable type \"%s\"", t)
}

func validateName(s string) error {
	if reVariableName.Match([]byte(s)) {
		return nil
	}

	return fmt.Errorf("variable name should be 1-20 characters long, and only contain ascii letters and numbers, dash or underscore")
}

func GetVariableByName(n string) (int, Variable) {
	for k, v := range Variables {
		if v.Id() == n {
			return k, v
		}
	}

	return -1, &Text{}
}

func SaveVariables() error {
	if Cfg.VariablesFile != nil {
		if ba, err := MarshalJSON(Variables); err == nil {
			return os.WriteFile(*Cfg.VariablesFile, ba, 0755)
		} else {
			return err
		}
	}

	return errors.New("file path not specified")
}

func RemoveVariableByName(n string) error {
	vk, _ := GetVariableByName(n)

	if vk == -1 {
		return fmt.Errorf("variable not found")
	} else if vk == 0 {
		if len(Variables) > 1 {
			Variables = Variables[vk+1:]
		} else {
			Variables = []Variable{}
		}
	} else if vk == len(Variables)-1 {
		Variables = Variables[:vk]
	} else {
		Variables = append(Variables[:vk], Variables[vk+1:]...)
	}

	SaveVariables()

	return nil
}

func ReadFromFile() error {
	vs := make([]interface{}, 0)

	if e := helper.ReadJSONFile(*Cfg.VariablesFile, &vs); e != nil {
		return e
	} else {
		for _, o := range vs {
			o := o.(map[string]interface{})

			if v, e2 := Create(o["type"].(string)); e2 == nil {
				e := v.Unmarshal(o)

				if e != nil {
					return fmt.Errorf("unable to decode variables file: %s", e)
				} else {
					Variables = append(Variables, v)
				}
			} else {
				return fmt.Errorf("unable to decode variables file: %s", e2)
			}

		}
	}

	return nil
}

func MarshalJSON(vs []Variable) ([]byte, error) {
	var bf bytes.Buffer

	bf.WriteByte('[')

	for k, v := range vs {

		b, e := json.Marshal(v)

		if e != nil {
			return []byte{}, e
		} else {
			if bf.Cap() < len(b)+bf.Len() {
				bf.Grow(len(b) + 10)
			}
			if k > 0 && len(b) > 1 {
				bf.Write([]byte{','})
			}

			if len(b) < 2 {
				return []byte{}, e
			} else {
				bf.Write(b)
				//bf.Write(b[:len(b)-1])
				//bf.Write([]byte(",\"t\":\"" + v.Type() + "\"}"))
			}
		}
	}

	bf.WriteByte(']')

	return bf.Bytes(), nil
}

func wrapBindings(s string, usedBindings *[]string, val interface{}) func() interface{} {
	return func() interface{} {
		*usedBindings = append(*usedBindings, s)
		return val
	}
}
