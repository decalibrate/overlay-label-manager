package variable

import "strconv"

type Toggle struct {
	VarType string `json:"type"`
	Name    string `json:"name"`
	Value   bool   `json:"value,omitempty"`
}

func (v *Toggle) Unmarshal(o map[string]interface{}) error {
	v.VarType = "toggle"
	if o["name"] != nil {
		v.Name = o["name"].(string)

		if o["value"] != nil {
			val := o["value"].(bool)
			v.Value = val
		}
	}

	return nil
}

func (v Toggle) Id() string {
	return v.Name
}

func (v Toggle) Type() string {
	return v.VarType
}

func (v *Toggle) Set(s ...string) error {
	v.VarType = "toggle"
	if len(s) == 0 {
		return nil
	}

	if len(s) >= 1 {
		if s[0] != "" {
			if e := validateName(s[0]); e == nil {
				v.Name = s[0]
			} else {
				return e
			}
		}
	}

	if len(s) >= 2 {
		if b, e := strconv.ParseBool(s[1]); e == nil {
			v.Value = b
		} else {
			v.Value = false
		}
	}

	return nil
}

func (v Toggle) GetBindings(bs *[]string) (string, interface{}) {
	var out = struct {
		Value    func() interface{} `liquid:"value"`
		Enabled  func() interface{} `liquid:"enabled"`
		Disabled func() interface{} `liquid:"disabled"`
		On       func() interface{} `liquid:"on"`
		Off      func() interface{} `liquid:"off"`
	}{
		Value:    wrapBindings(v.Name, bs, v.Value),
		Enabled:  wrapBindings(v.Name, bs, v.Value),
		Disabled: wrapBindings(v.Name, bs, !v.Value),
		On:       wrapBindings(v.Name, bs, v.Value),
		Off:      wrapBindings(v.Name, bs, !v.Value),
	}

	return v.Name, out
}
