package template

import (
	"testing"

	"github.com/decalibrate/overlay-label-manager/internal/variable"
)

func Test(t *testing.T) {

	fv := 10.0
	fg := 20.0
	fct := "Completed!"

	variable.Variables = []variable.Variable{
		&variable.Counter{Name: "variable", Value: &fv, Goal: &fg, CompletionText: &fct},
		&variable.Counter{Name: "complete", Value: &fg, Goal: &fg, CompletionText: &fct},
		&variable.Counter{Name: "nogoal", Value: &fv},
		&variable.Counter{Name: "novalue"},
		&variable.Text{Name: "text", Value: &fct},
		&variable.Text{Name: "textnovalue"},
	}

	t.Run("Processor", testProcessor)

	t.Run("TextVariable", testTextVariable)

	t.Run("CounterVariable", testCounterVariable)

	t.Run("SetTemplate", testSetTemplate)
}

func testSetTemplate(t *testing.T) {

	template := Template{}

	tr := "true"
	fa := "false"
	tmpl := "string {{ variable }} string"
	tmplerr := "string {{ variable } string"
	tplerrmsg := "unexpected '}' on line 0, position 19; expected '|', ':'. Second '}' may be missing"
	namerrmsg := "template names should be 1-20 characters long, and only contain ascii letters and numbers, dash or underscore"
	nm := "variable_test"

	lbl := "string 10/20 string"

	tests := []struct {
		tm, nm   string
		hd       string
		lbl, err string
	}{
		{lbl: lbl, err: namerrmsg},
		{tm: tmpl, lbl: lbl, err: namerrmsg},
		{nm: nm, lbl: lbl},
		{hd: tr, lbl: lbl, err: namerrmsg},
		{hd: fa, lbl: lbl, err: namerrmsg},
		{tm: tmpl, nm: nm, hd: fa, lbl: lbl},
		{nm: nm, tm: tmplerr, err: tplerrmsg},
	}

	for _, tt := range tests {

		t.Run("", func(t *testing.T) {

			e := template.Set(tt.tm, tt.nm, tt.hd)
			if e != nil && tt.err == "" {
				t.Errorf("\nInput:  %s\nExepcted: no error\nGot:    %s", template.Template, e)
			} else if e != nil && tt.err != e.Error() {
				t.Errorf("\nInput: %s\nExpected: %s\nGot:      %s", template.Template, tt.err, e)
			}

		})
	}
}

func testTextVariable(t *testing.T) {

	tests := []struct {
		k            int
		in, err, out string
	}{
		{in: `{{ text }}`, out: "Completed!"},
		{in: `{{ textnovalue }}`, out: ""},
		{in: `{{ text.unknownkey }}`, out: ""},
	}

	for _, tt := range tests {

		//testname := fmt.Sprintf("%d", tt.k)
		t.Run(tt.in, func(t *testing.T) {

			tmplt := Template{Template: tt.in}

			e := tmplt.RenderLiquidTemplate()
			str := tmplt.Label

			if tt.err != "" && (e == nil || e.Error() != tt.err) {
				t.Errorf("\nInput: %s\nExpected: err %s\nGot:      %s", tt.in, tt.err, e)
			} else if tt.err == "" && e != nil {
				t.Errorf("\nInput: %s\nExpected: no error\nGot:      %s", tt.in, e)
			} else if tt.err == "" {
				if string(str) != tt.out {
					t.Errorf("\nInput: %s\nExpected: %s\nGot:      %s", tt.in, tt.out, str)
				}
			}
		})
	}
}

func testCounterVariable(t *testing.T) {

	tests := []struct {
		k            int
		in, err, out string
	}{
		{in: `string {{ variable }} string`, out: "string 10/20 string"},
		{in: `string {{ complete }} string`, out: "string Completed! string"},
		{in: `string {{ nogoal }} string`, out: "string 10 string"},
		{in: `string {{ novalue }} string`, out: "string  string"},
		{in: `string {{ unknown }} string`, out: "string  string"},
	}

	for _, tt := range tests {

		//testname := fmt.Sprintf("%d", tt.k)
		t.Run(tt.in, func(t *testing.T) {

			tmplt := Template{Template: tt.in}

			e := tmplt.RenderLiquidTemplate()
			str := tmplt.Label

			if tt.err != "" && (e == nil || e.Error() != tt.err) {
				t.Errorf("\nInput: %s\nExpected: err %s\nGot:      %s", tt.in, tt.err, e)
			} else if tt.err == "" && e != nil {
				t.Errorf("\nInput: %s\nExpected: no error\nGot:      %s", tt.in, e)
			} else if tt.err == "" {
				if string(str) != tt.out {
					t.Errorf("\nInput: %s\nExpected: %s\nGot:      %s", tt.in, tt.out, str)
				}
			}
		})
	}
}

func testProcessor(t *testing.T) {

	tests := []struct {
		k            int
		in, err, out string
	}{
		{in: "string", out: "string"},
		{in: `string {{ }}`, out: "string "},
		{in: `string {{ "string" }}`, out: "string string"},
		{in: `{{ "string" }} string`, out: "string string"},
		{in: `string {{ "string" }} string`, out: "string string string"},
		{in: `string {{ "string" | split: "" | reverse | join: "" }} string`, out: "string gnirts string"},
		{in: `string {{ 123 | reverse: 2, 2, 2 }} string`, out: "string 321 string"},
		{in: `string {{ "string" | plus: 3 }} string`, out: "string string string"},
		{in: `string {{ 123 | plus: 3 }} string`, out: "string 126 string"},
		{in: `{{ variable }} {{ variable }}`, out: "10/20 10/20"},
		{in: `{{ variable }} {{ unknown }}`, out: "10/20 "},
		{in: `string {{ variable }} string {{ complete }} string`, out: "string 10/20 string Completed! string"},
		{in: `string {{ variable.value }} string`, out: "string 10 string"},
		{in: `string {{ variable.value | plus: 3 }} string`, out: "string 13 string"},
		{in: `string {{ variable.goal | plus: 3 }} string`, out: "string 23 string"},
		{in: `string {{ variable.completion_text | plus: 3 }} string`, out: "string Completed! string"},
		{in: `string {{ variable.completion_text | append: variable.value }} string`, out: "string Completed!10 string"},
		{in: `string {{ variable.goal | minus: variable.value | append: " remaining" | split: "" | reverse | join: "" }} string`, out: "string gniniamer 01 string"},
		{in: `{{ "string" | upcase }}`, out: "STRING"},
		{in: `{{ "STRING" | downcase }}`, out: "string"},
		{in: `{{ "string string" | title }}`, out: "String String"},
		{in: `{{ "string string" | capitalize }}`, out: "String string"},
		{in: `{{ "string string" | sarcastic }}`, out: "StRiNg StRiNg"},
		{in: `{{ "      string     " | strip }}`, out: "string"},
		{in: "{{ \"  string\n  \" | strip_newlines }}", out: "  string  "},
		{in: `{{ "      string     " | lstrip }}`, out: "string     "},
		{in: `{{ "      string     " | rstrip }}`, out: "      string"},
		{in: `{{ "&lt;string&gt;" | unescape }}`, out: "<string>"},
		{in: `{{ "<string>" | escape }}`, out: "&lt;string&gt;"},
		{in: `{{ "<string>" | escape }}`, out: "&lt;string&gt;"},
		{in: "{{ \"  string\n  \" | newline_to_br }}", out: "  string<br />\n  "},
		{in: `{{ "stringstring" | remove: "str" }}`, out: "inging"},
		{in: `{{ "stringstring" | remove_first: "ing" }}`, out: "strstring"},
		{in: `{{ "stringstring" | replace: "ing", "ung" }}`, out: "strungstrung"},
		{in: `{{ "stringstring" | replace: "ing" }}`, out: "strstr"},
		{in: `{{ "stringstring" | replace_first: "ing" }}`, out: "strstring"},
		{in: `{{ "stringstring" | replace_first: "ing", "ung" }}`, out: "strungstring"},
		{in: `{{ "string" | prepend: "fisrt" }}`, out: "fisrtstring"},
		{in: `{{ "string" | append: "last" }}`, out: "stringlast"},
		{in: `{{ "" | default: "d" }}`, out: "d"},
		{in: `{{ "hi" | default: "d" }}`, out: "hi"},
		{in: `{{ "0" | default: "d" }}`, out: "d"},
		{in: `{{ "false" | default: "d" }}`, out: "d"},
		{in: `{{ .5 | default: "d" }}`, out: ".5"},
		{in: `{{ -1 | default: "d" }}`, out: "-1"},
		{in: `{{ 0 | default: "d" }}`, out: "d"},
		{in: `{{ 1 | default: "d" }}`, out: "1"},
		{in: `{{ false | default: "d" }}`, out: "d"},
		{in: `{{ true | default: "d" }}`, out: "true"},
		{in: `{{ nil | default: "d" }}`, out: "d"},
		{in: `{{ "nil" | default: "d" }}`, out: "nil"},
		{in: `{{ "88" | size }}`, out: "2"},
		{in: `{{ 88 | size }}`, out: "2"},
		{in: `{{ nil | size }}`, out: "0"},
		{in: `{{ true | size }}`, out: "4"},
		{in: `{{ false | size }}`, out: "5"},
		{in: `{{ "string" | truncate }}`, out: "string"},
		{in: `{{ "string" | truncate: 3 }}`, out: "..."},
		{in: `{{ "string" | truncate: 5 }}`, out: "st..."},
		{in: `{{ "string" | truncate: 7 }}`, out: "string"},
		{in: `{{ "string" | truncate: 7, "!" }}`, out: "string"},
		{in: `{{ "string" | truncate: 5, "!" }}`, out: "stri!"},
		{in: `{{ "string" | round }}`, out: "string"},
		{in: `{{ "string" | round: 2 }}`, out: "string"},
		{in: `{{ 9.500001 | round }}`, out: "9.50"},
		{in: `{{ 9.500001 | round: 3 }}`, out: "9.500"},
		{in: `{{ 9.500001 | ceil }}`, out: "10"},
		{in: `{{ 9.500001 | floor }}`, out: "9"},
		{in: `{{ 9.500001 | abs }}`, out: "9.500001"},
		{in: `{{ -9.500001 | abs }}`, out: "9.500001"},
		{in: `{{ -9.5 | minus: 3.2 }}`, out: "-12.7"},
		{in: `{{ -9.5 | plus: 3.2 }}`, out: "-6.3"},
		{in: `{{ -9.5 | times: 100 }}`, out: "-950"},
		{in: `{{ -9.5 | times: 0.5 }}`, out: "-4.75"},
		{in: `{{ -9.5 | divided_by: 2 }}`, out: "-4.75"},
		{in: `{{ -9.5 | divided_by: 0.5 }}`, out: "-19"},
		{in: `{{ -9.5 | at_most: 1 }}`, out: "-9.5"},
		{in: `{{ -9.5 | at_most: -10 }}`, out: "-10"},
		{in: `{{ -9.5 | at_least: 1 }}`, out: "1"},
		{in: `{{ -9.5 | at_least: -10 }}`, out: "-9.5"},
		{in: `{{ 5 | modulo: 2 }}`, out: "1"},
		{in: `{{ 5 | times: 0 | default: true }}`, out: "true"},
		{in: `{{ 5 | times: 0 | default: false }}`, out: "false"},
		{in: `{{ "string" | abs }}`, out: "string"},
		{in: `{{ "-3" | abs }}`, out: "3"},
		{in: `{{ "-3.01" | abs }}`, out: "3.01"},
	}

	for _, tt := range tests {

		//testname := fmt.Sprintf("%d", tt.k)
		t.Run(tt.in, func(t *testing.T) {

			tmplt := Template{Template: tt.in}

			e := tmplt.RenderLiquidTemplate()
			str := tmplt.Label

			if tt.err != "" && (e == nil || e.Error() != tt.err) {
				t.Errorf("\nInput: %s\nExpected: err %s\nGot:      %s", tt.in, tt.err, e)
			} else if tt.err == "" && e != nil {
				t.Errorf("\nInput: %s\nExpected: no error\nGot:      %s", tt.in, e)
			} else if tt.err == "" {
				if string(str) != tt.out {
					t.Errorf("\nInput: %s\nExpected: %s\nGot:      %s", tt.in, tt.out, str)
				}
			}
		})
	}
}
