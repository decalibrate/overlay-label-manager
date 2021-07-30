package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/decalibrate/overlay-label-manager/internal/template"
	"github.com/decalibrate/overlay-label-manager/internal/variable"
	"github.com/gorilla/mux"
)

type templateResponseStruct struct {
	Template  template.Template   `json:"t"`
	Variables []variable.Variable `json:"v"`
}

func TemplateHandler(w http.ResponseWriter, r *http.Request) {

	qs := r.URL.Query()

	vars := mux.Vars(r)
	tnam := vars["name"]

	switch r.Method {
	case "GET":
		if tnam == "" {
			rsp, err := json.Marshal(template.Templates)
			if err != nil {
				apiResponse(w, r, http.StatusInternalServerError, formatErrorForResponse(err), "")
				return
			}
			apiResponse(w, r, http.StatusOK, rsp, "")
			return
		}

		tk, t := template.GetTemplateByName(tnam)
		if tk == -1 {
			apiResponse(w, r, http.StatusNotFound, formatErrorForResponse(nil), "")
			return
		}

		var rsp []byte
		var err error

		rsp, err = json.Marshal(t)
		if err == nil {
			if qs.Get("bv") == "true" {
				vs := make([]variable.Variable, 0)
				vsn := template.Template2VariableMap[t.Name]

				for vn, u := range vsn {
					if u {
						vk, v := variable.GetVariableByName(vn)
						if vk != -1 {
							vs = append(vs, v)
						}
					}
				}
				rsp, err = json.Marshal(templateResponseStruct{Template: t, Variables: vs})
			}
		}

		if err != nil {
			apiResponse(w, r, http.StatusInternalServerError, formatErrorForResponse(err), "")
			return
		}

		apiResponse(w, r, http.StatusOK, rsp, "")
	case "DELETE":

		if err := template.RemoveTemplateByName(tnam); err != nil {
			apiResponse(w, r, http.StatusNotFound, formatErrorForResponse(nil), "")
			return
		}

		apiResponse(w, r, http.StatusOK, nil, "")

	case "POST", "PUT":

		hid := qs.Get("hidden")

		tk, t := template.GetTemplateByName(tnam)

		var tem []byte

		if r.Method == "POST" {
			tem, _ = ioutil.ReadAll(r.Body)
			r.Body.Close()

			if tk == -1 {
				t = template.Template{}
			}

			t.Set(string(tem), tnam, hid)

		} else if r.Method == "PUT" {
			if tk == -1 {
				apiResponse(w, r, http.StatusNotFound, formatErrorForResponse(nil), "")
				return
			}

			if hid == "true" {
				t.Hide()
			} else {
				t.Show()
			}
		}

		t.Refresh()

		rsp, err := json.Marshal(t)
		if err != nil {
			apiResponse(w, r, http.StatusInternalServerError, formatErrorForResponse(err), "")
			return
		}

		if tk == -1 {
			template.Templates = append(template.Templates, &t)
		} else {
			template.Templates[tk] = &t
		}

		template.SaveTemplates()
		t.SaveLabel()

		apiResponse(w, r, http.StatusOK, rsp, "")
	}
}
