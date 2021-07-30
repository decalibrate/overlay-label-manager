package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/decalibrate/overlay-label-manager/internal/configuration"
	"github.com/decalibrate/overlay-label-manager/internal/helper"
	"github.com/decalibrate/overlay-label-manager/internal/template"
)

type RestartServer func()

func testForLocalIP(r *http.Request) bool {
	return strings.HasPrefix(r.Host, "127.0.0.1:") || strings.HasPrefix(r.Host, "[::1]:") || strings.HasPrefix(r.Host, "localhost:")
}

func ConfigurationHandler(w http.ResponseWriter, r *http.Request) {

	// configuration can only be changed from local context
	if !testForLocalIP(r) {
		apiResponse(w, r, http.StatusForbidden, []byte(http.StatusText(http.StatusForbidden)), "")
		return
	}

	switch r.Method {
	case "GET":

		rsp, err := json.Marshal(Cfg)
		if err != nil {
			apiResponse(w, r, http.StatusInternalServerError, formatErrorForResponse(err), "")
			return
		}

		apiResponse(w, r, http.StatusOK, rsp, "")
	case "POST":
		tem, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()

		newcfg := &configuration.ConfigStruct{}
		if err := json.Unmarshal([]byte(tem), newcfg); err != nil {
			apiResponse(w, r, http.StatusBadRequest, formatErrorForResponse(err), "")
			return
		}

		newcfg.ConfigFile = Cfg.ConfigFile
		newcfg.VariablesFile = Cfg.VariablesFile
		newcfg.TemplatesFile = Cfg.TemplatesFile

		if newcfg.LabelDirectory != nil && *newcfg.LabelDirectory != *Cfg.LabelDirectory {
			if _, err := os.Stat(*newcfg.LabelDirectory); err != nil {
				log.Printf("[conf-change] Labels directory does not exist - %s", *newcfg.LabelDirectory)

				if err := os.Mkdir(*newcfg.LabelDirectory, 0755); err != nil {
					apiResponse(w, r, http.StatusBadRequest, []byte("Could not create labels file"), "")
					return
				} else {
					log.Printf("[conf-change] Successfully created labels directory")
				}
			} else {
				log.Printf("[conf-change] Found existing folder %s", *newcfg.LabelDirectory)
				log.Printf("[conf-change] Set labels directory %s", *newcfg.LabelDirectory)
			}
		}

		restardNeeded := false

		if newcfg.Port != nil && *newcfg.Port != *Cfg.Port && *newcfg.Port > 0 && *newcfg.Port < 65536 {
			Cfg.Port = newcfg.Port
			restardNeeded = true
		}

		if err := helper.SaveJSONFile(*Cfg.ConfigFile, newcfg); err != nil {
			apiResponse(w, r, http.StatusInternalServerError, formatErrorForResponse(err), "")
			return
		}

		if newcfg.LabelDirectory != nil && *newcfg.LabelDirectory != *Cfg.LabelDirectory {
			Cfg.LabelDirectory = newcfg.LabelDirectory
			template.BuildAllLabels()
		}

		resp, _ := json.Marshal(Cfg)

		apiResponse(w, r, http.StatusOK, resp, "")

		if restardNeeded {
			go RestartHttpServer()
		}
	}
}
