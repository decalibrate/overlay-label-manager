package variable

type Text struct {
	VarType string  `json:"type"`
	Name    string  `json:"name"`
	Value   *string `json:"value,omitempty"`
}

func (v *Text) Unmarshal(o map[string]interface{}) error {
	v.VarType = "text"
	if o["name"] != nil {
		v.Name = o["name"].(string)

		if o["value"] != nil {
			val := o["value"].(string)
			v.Value = &val
		}
	}

	return nil
}

func (v Text) Id() string {
	return v.Name
}

func (v Text) Type() string {
	return v.VarType
}

func (v *Text) Set(n string, s map[string][]string) error {
	v.VarType = "text"

	if e := validateName(n); e == nil {
		v.Name = n
	} else {
		return e
	}

	if s["value"] != nil && len(s["value"]) > 0 {
		v.Value = &s["value"][0]
	} else {
		v.Value = nil
	}

	return nil
}

func (v Text) GetBindings(bs *[]string) (string, interface{}) {
	var val string

	if v.Value != nil {
		val = *v.Value
	}
	var out = struct {
		Value func() interface{} `liquid:"value"`
	}{
		Value: wrapBindings(v.Name, bs, val)}

	return v.Name, out
}
