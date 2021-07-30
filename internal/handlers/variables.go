package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/decalibrate/overlay-label-manager/internal/template"
	"github.com/decalibrate/overlay-label-manager/internal/variable"
	"github.com/gorilla/mux"
)

func VariableHandler(w http.ResponseWriter, r *http.Request) {

	qs := r.URL.Query()

	vars := mux.Vars(r)
	vnam := vars["name"]

	switch r.Method {
	case "GET":
		if vnam == "" {
			rsp, err := variable.MarshalJSON(variable.Variables)
			if err != nil {
				apiResponse(w, r, http.StatusInternalServerError, formatErrorForResponse(err), "")
				return
			}
			apiResponse(w, r, http.StatusOK, rsp, "")
			return
		}

		vk, v := variable.GetVariableByName(vnam)
		if vk == -1 {
			apiResponse(w, r, http.StatusNotFound, formatErrorForResponse(nil), "")
			return
		}

		rsp, err := json.Marshal(v)
		if err != nil {
			apiResponse(w, r, http.StatusInternalServerError, formatErrorForResponse(err), "")
			return
		}

		apiResponse(w, r, http.StatusOK, rsp, "")
	case "DELETE":

		if err := variable.RemoveVariableByName(vnam); err != nil {
			apiResponse(w, r, http.StatusNotFound, formatErrorForResponse(err), "")
			return
		}

		template.RebuildTemplatesReferencingVariable(vnam)

		apiResponse(w, r, http.StatusOK, nil, "")
	case "PUT":
		typ := qs.Get("type")

		vk, v := variable.GetVariableByName(vnam)

		if vk != -1 && typ != v.Type() {
			apiResponse(w, r, http.StatusInternalServerError, formatErrorForResponse(errors.New("mismatched type")), "")
			return
		} else if vk == -1 {
			if v2, err := variable.Create(typ); err != nil {
				apiResponse(w, r, http.StatusInternalServerError, formatErrorForResponse(err), "")
				return
			} else {
				v = v2
			}
		}

		if e := v.Set(vnam, qs); e != nil {
			apiResponse(w, r, http.StatusInternalServerError, formatErrorForResponse(e), "")
			return
		}

		rsp, err := json.Marshal(v)
		if err != nil {
			apiResponse(w, r, http.StatusInternalServerError, formatErrorForResponse(err), "")
			return
		}

		if vk == -1 {
			variable.Variables = append(variable.Variables, v)
		} else {
			variable.Variables[vk] = v
		}

		variable.SaveVariables()
		template.RebuildTemplatesReferencingVariable(vnam)

		apiResponse(w, r, http.StatusOK, rsp, "")
	}
}
