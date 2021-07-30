package template

import (
	"fmt"
	"html"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/decalibrate/overlay-label-manager/internal/configuration"
	"github.com/decalibrate/overlay-label-manager/internal/helper"
	"github.com/decalibrate/overlay-label-manager/internal/variable"
	"github.com/osteele/liquid"
)

var Cfg *configuration.ConfigStruct
var Templates []*Template

var Template2VariableMap = make(map[string]map[string]bool)
var Variable2TemplateMap = make(map[string]map[string]bool)

type Template struct {
	Name     string `json:"name,omitempty"`
	Template string `json:"template,omitempty"`
	Hidden   bool   `json:"hidden,omitempty"`
	Label    string `json:"-"`
}

var reTemplateName = regexp.MustCompile("^[a-zA-Z_][0-9a-zA-Z_]{0,19}$")

var engine *liquid.Engine

func prepareEngine() {
	engine = liquid.NewEngine()
	engine.RegisterFilter("title", strings.Title)

	engine.RegisterFilter("sarcastic", func(s string) string {
		rs, upper := []rune(s), false
		for i, r := range rs {
			if unicode.IsLetter(r) {
				if upper = !upper; upper {
					rs[i] = unicode.ToUpper(r)
				} else {
					rs[i] = unicode.ToLower(r)
				}
			}
		}
		return string(rs)
	})

	engine.RegisterFilter("unescape", html.UnescapeString)

	engine.RegisterFilter("unescape", html.EscapeString)

	engine.RegisterFilter("at_least", func(a float64, b float64) float64 {
		if a < b {
			return b
		}
		return a
	})

	engine.RegisterFilter("at_most", func(a float64, b float64) float64 {
		if a > b {
			return b
		}
		return a
	})

	engine.RegisterFilter("wrap_span", func(a string) string {
		return "<span class=\"variable\">" + a + "</span>"
	})
}

func (t *Template) RenderLiquidTemplate() error {

	if engine == nil {
		prepareEngine()
	}

	liquid.NewEngine()

	template := t.Template
	bindings := map[string]interface{}{}
	var err error

	var usedBindings = make([]string, 0)

	for _, v := range variable.Variables {
		k, b := v.GetBindings(&usedBindings)
		bindings[k] = b
	}

	clearVariable2TemplateMapping(t.Name)

	t.Label, err = engine.ParseAndRenderString(template, bindings)
	if err != nil {
		return err
	}

	for _, b := range usedBindings {
		Template2VariableMap[t.Name][b] = true
		if Variable2TemplateMap[b] == nil {
			Variable2TemplateMap[b] = make(map[string]bool)
		}
		Variable2TemplateMap[b][t.Name] = true
	}

	return err
}

func GetTemplateByName(n string) (int, Template) {
	for k, t := range Templates {
		if t.Name == n {
			t.Refresh()
			return k, *t
		}
	}

	return -1, Template{}
}

func SaveTemplates() {
	helper.SaveJSONFile(*Cfg.TemplatesFile, Templates)
}

func RebuildTemplatesReferencingVariable(n string) {
	if Variable2TemplateMap[n] != nil {
		for tn, u := range Variable2TemplateMap[n] {
			if u {
				tk, t := GetTemplateByName(tn)
				if tk != -1 {
					t.Refresh()
					t.SaveLabel()
				}
			}
		}
	}
}
func clearVariable2TemplateMapping(n string) {
	if Template2VariableMap[n] != nil {
		for vn, u := range Template2VariableMap[n] {
			if u {
				if Variable2TemplateMap[vn] != nil {
					delete(Variable2TemplateMap[vn], n)
				}
			}
			delete(Template2VariableMap[n], vn)
		}
	} else {
		Template2VariableMap[n] = make(map[string]bool)
	}
}

func RemoveTemplateByName(n string) error {

	tk, _ := GetTemplateByName(n)
	if tk == -1 {
		return fmt.Errorf("template not found")
	} else if tk == 0 {
		if len(Templates) > 1 {
			Templates = Templates[tk+1:]
		} else {
			Templates = []*Template{}
		}
	} else if tk == len(Templates)-1 {
		Templates = Templates[:tk]
	} else {
		Templates = append(Templates[:tk], Templates[tk+1:]...)
	}

	clearVariable2TemplateMapping(n)
	delete(Template2VariableMap, n)

	SaveTemplates()

	return nil
}

func BuildAllLabels() {
	for _, t := range Templates {
		t.RenderLiquidTemplate()
		t.SaveLabel()
	}

	SaveTemplates()
}

func (t *Template) Set(templateString string, name string, hide string) error {
	var e, e2 error

	if name == "" && !reTemplateName.Match([]byte(name)) {
		return fmt.Errorf("template names should be 1-20 characters long, and only contain ascii letters and numbers, dash or underscore")
	}

	t.Name = name
	t.Template = templateString

	if hide == "true" {
		t.Hide()
	} else if hide == "false" {
		t.Show()
	}

	e2 = t.Refresh()

	if e != nil {
		return e
	}

	if e2 != nil {
		return e2
	}

	return nil
}

func (t *Template) Refresh() error {

	e2 := t.RenderLiquidTemplate()

	if e2 != nil {
		return e2
	}

	return nil
}

func (t *Template) Hide() {
	t.Hidden = true
}

func (t *Template) Show() {
	t.Hidden = false
}

func (t *Template) SaveLabel() {
	if Cfg.LabelDirectory != nil && *Cfg.LabelDirectory != "" {
		tn := strings.ReplaceAll(strings.ReplaceAll(t.Name, "/", ""), ",", "")

		var b []byte
		if !t.Hidden {
			b = []byte(t.Label)
		}

		if err := os.WriteFile(filepath.Join(*Cfg.LabelDirectory, tn+".txt"), b, 0755); err != nil {
			panic(err)
		}
	}
}
