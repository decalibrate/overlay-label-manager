package variable

import (
	"github.com/decalibrate/overlay-label-manager/internal/helper"
)

type Counter struct {
	VarType        string   `json:"type"`
	Name           string   `json:"name"`
	Value          *float64 `json:"value,omitempty"`
	Goal           *float64 `json:"goal,omitempty"`
	CompletionText *string  `json:"completion_text,omitempty"`
}

func (v *Counter) Unmarshal(o map[string]interface{}) error {

	if o["name"] != nil {
		v.Name = o["name"].(string)

		if o["value"] != nil {
			val := o["value"].(float64)
			v.Value = &val
		}
		if o["goal"] != nil {
			g := o["goal"].(float64)
			v.Goal = &g
		}
		if o["completion_text"] != nil {
			ct := o["completion_text"].(string)
			v.CompletionText = &ct
		}
	}

	return nil
}

func (v Counter) Id() string {
	return v.Name
}

func (v Counter) Type() string {
	return v.VarType
}

func (v *Counter) Set(n string, s map[string][]string) error {
	v.VarType = "counter"

	if e := validateName(n); e == nil {
		v.Name = n
	} else {
		return e
	}

	if s["value"] != nil && len(s["value"]) > 0 {
		if val, e := helper.ParseFloat(s["value"][0]); e != nil {
			return e
		} else {
			v.Value = &val
		}
	} else {
		v.Value = nil
	}

	if s["goal"] != nil && len(s["goal"]) > 0 {
		if g, e := helper.ParseFloat(s["goal"][0]); e != nil {
			return e
		} else {
			v.Goal = &g
		}
	} else {
		v.Goal = nil
	}

	if s["completion_text"] != nil && len(s["completion_text"]) > 0 {
		v.CompletionText = &s["completion_text"][0]
	} else {
		v.CompletionText = nil
	}

	return nil
}

func (v *Counter) GetValue() float64 {
	if v.Value != nil {
		return *v.Value
	}
	return float64(0)

}
func (v *Counter) GetGoal() float64 {
	if v.Goal != nil {
		return *v.Goal
	}
	return float64(0)

}
func (v *Counter) GetCompletionText() string {
	if v.CompletionText != nil {
		return *v.CompletionText
	}
	return ""

}

func (v Counter) GetBindings(bs *[]string) (string, interface{}) {

	var out = struct {
		Value          func() interface{} `liquid:"value"`
		Goal           func() interface{} `liquid:"goal"`
		CompletionText func() interface{} `liquid:"completion_text"`
	}{
		Value:          wrapBindings(v.Name, bs, v.GetValue()),
		Goal:           wrapBindings(v.Name, bs, v.GetGoal()),
		CompletionText: wrapBindings(v.Name, bs, v.GetCompletionText())}

	return v.Name, out
}
