package main

import (
	"embed"
	"flag"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/decalibrate/overlay-label-manager/internal/configuration"
	"github.com/decalibrate/overlay-label-manager/internal/handlers"
	"github.com/decalibrate/overlay-label-manager/internal/helper"
	"github.com/decalibrate/overlay-label-manager/internal/template"
	"github.com/decalibrate/overlay-label-manager/internal/variable"
	"github.com/gorilla/mux"
)

//go:embed static
var embededFiles embed.FS

func main() {

	var cfg = &configuration.ConfigStruct{}

	cfgP := flag.String("cfg", "", "path to the settings storage file")
	cfgPort := flag.Int("p", -1, "port to run the webapp on")
	cfgDev := flag.Bool("dev", false, "run in development mode")
	flag.Parse()

	if err := helper.ReadJSONFile(*cfgP, cfg); err != nil {
		log.Printf("[load]: -cfg not specified or specified file doesn't exist '%s'. Using Default config", *cfgP)

		defaultPath := ".overlaylabelmanager"

		v := defaultPath + "/variables.json"
		t := defaultPath + "/templates.json"
		c := defaultPath + "/config.json"
		l := defaultPath + "/../labels"

		if home, err := os.UserHomeDir(); err != nil {
			v = filepath.Join(home, v)
			t = filepath.Join(home, t)
			c = filepath.Join(home, c)
			l = filepath.Join(home, l)

			log.Printf("[load] Using %s as app storage location", filepath.Join(home, defaultPath))

			_ = os.MkdirAll(filepath.Join(home, defaultPath), 0600)

		} else {
			ex, err := os.Executable()
			if err != nil {
				panic(err)
			}
			exPath := filepath.Dir(ex)

			v = filepath.Join(exPath, v)
			t = filepath.Join(exPath, t)
			c = filepath.Join(exPath, c)
			l = filepath.Join(exPath, l)

			log.Printf("[load] Using %s as app storage location", filepath.Join(exPath, defaultPath))

			_ = os.Mkdir(defaultPath, 0600)

		}

		cfg.VariablesFile = &v
		cfg.TemplatesFile = &t
		cfg.ConfigFile = &c
		cfg.LabelDirectory = &l

		if err := helper.ReadJSONFile(*cfg.ConfigFile, cfg); err == nil {
			log.Printf("[load] Previous configs retrieved from the defaulted app storage directory")
		}

	}

	if cfgPort != nil && *cfgPort != -1 && *cfgPort > 0 {
		cfg.Port = cfgPort
	} else if cfg.Port == nil {
		p := 9144
		cfg.Port = &p
	}

	updateAvailable := cfg.IsUpdateAvailable()
	if !updateAvailable {
		log.Printf("[load] You're running the most recent version - last checked: %s", (*cfg.LastVersionCheckDate).UTC().Format(time.RFC3339))
	}

	handlers.Cfg = cfg
	variable.Cfg = cfg
	template.Cfg = cfg

	if _, err := os.Stat(*cfg.LabelDirectory); err != nil {
		log.Printf("[load] Labels directory does not exist - %s", *cfg.LabelDirectory)

		if err := os.Mkdir(*cfg.LabelDirectory, 0755); err != nil {
			log.Fatalf("[load] could not create labels directory - %s", err)
		} else {
			log.Printf("[load] Successfully created labels directory")
		}
	} else {
		log.Printf("[load] Labels directory already exists %s", *cfg.LabelDirectory)
	}

	if err := helper.SaveJSONFile(*cfg.ConfigFile, cfg); err != nil {
		log.Fatalf("[load] Could not save config file - %s", err)
	}

	if err := variable.ReadFromFile(); err != nil {
		log.Printf("[load] Variables file not found, or corrupted, starting from fresh. %s load-error: %s", *cfg.VariablesFile, err)
	} else {
		log.Printf("[load] Successfully loaded previous variables")
	}

	if err := helper.ReadJSONFile(*cfg.TemplatesFile, &template.Templates); err != nil {
		log.Printf("[load] Templates file not found, or corrupted, starting from fresh. %s load-error: %s", *cfg.TemplatesFile, err)
	} else {
		log.Printf("[load] Successfully loaded previous templates")
		template.BuildAllLabels()
	}

	router := mux.NewRouter()

	handlers.GetFileSystem(embededFiles, cfgDev)

	//router.HandleFunc("/ws", handlers.WSHandler)
	router.HandleFunc("/conf", handlers.ConfigurationHandler).Methods("GET", "POST")

	router.HandleFunc("/v/{name}", handlers.VariableHandler).Methods("GET", "DELETE", "PUT")
	router.HandleFunc("/v", handlers.VariableHandler).Methods("GET")

	router.HandleFunc("/t/{name}", handlers.TemplateHandler).Methods("GET", "DELETE", "POST", "PUT")
	router.HandleFunc("/t", handlers.TemplateHandler).Methods("GET")

	router.HandleFunc("/l/{name}", handlers.LabelHandler).Methods("GET")

	router.PathPrefix("/bv").HandlerFunc(handlers.StaticHandlerBrowserView).Methods("GET")

	router.PathPrefix("/").HandlerFunc(handlers.StaticHandler).Methods("GET")

	router.Use(mux.CORSMethodMiddleware(router))

	handlers.Router = router

	//log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*cfg.Port), router))

	if updateAvailable {
		log.Print("\n\n*** An Update for Overlay Label Manager is available *** \nhttps://github.com/decalibrate/overlay-label-manager/releases\n\n")
	}

	handlers.HttpServerExitDone = &sync.WaitGroup{}

	handlers.RestartHttpServer()
	handlers.HttpServerExitDone.Wait()
}
