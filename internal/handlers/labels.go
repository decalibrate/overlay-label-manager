package handlers

import (
	"net/http"

	"github.com/decalibrate/overlay-label-manager/internal/template"
	"github.com/gorilla/mux"
)

func LabelHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	lnam := vars["name"]

	switch r.Method {
	case "GET":
		if lnam == "" {
			apiResponse(w, r, http.StatusBadRequest, formatErrorForResponse(nil), "")
			return
		}

		tk, t := template.GetTemplateByName(lnam)

		if tk == -1 {
			apiResponse(w, r, http.StatusBadRequest, formatErrorForResponse(nil), "")
			return
		}

		err := t.Refresh()
		if err != nil {
			apiResponse(w, r, http.StatusBadRequest, formatErrorForResponse(err), "")
			return
		}

		if !t.Hidden {
			apiResponse(w, r, http.StatusOK, []byte(t.Label), "text/plain")
		} else {
			apiResponse(w, r, http.StatusOK, []byte{}, "text/plain")
		}

	}
}
