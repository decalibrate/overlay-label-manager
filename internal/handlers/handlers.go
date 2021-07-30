package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/decalibrate/overlay-label-manager/internal/configuration"
	"github.com/gorilla/mux"
)

var Cfg *configuration.ConfigStruct

var Srv *http.Server
var Router *mux.Router

var HttpServerExitDone *sync.WaitGroup

// var clients = make(map[*websocket.Conn]bool)
// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

type errorResponseStruct struct {
	Message string `json:"error"`
}

func apiResponse(w http.ResponseWriter, r *http.Request, statusCode int, message []byte, messageType string) {

	if statusCode == http.StatusOK && r.URL.Query().Get("jsonp") == "1" {
		qs := r.URL.Query()
		cbp := qs.Get("cbp")
		if cbp == "" {
			cbp = "window._jsonp"
		}
		cbn := qs.Get("cbn")

		if strings.ContainsAny(cbp, "{}();&|=+-^*/") {
			apiResponse(w, r, http.StatusBadRequest, formatErrorForResponse(nil), "")
			return
		}
		if cbn == "" || strings.Contains(cbn, "\"") {
			apiResponse(w, r, http.StatusBadRequest, formatErrorForResponse(nil), "")
			return
		}

		message = []byte(cbp + "[\"" + cbn + "\"](" + string(message) + ");")
		w.Header().Set("Content-Type", "text/javascript")
	} else if r.Header.Get("Accepts") == "text/plain" {
		w.Header().Set("Content-Type", "text/plain")
	} else if messageType != "" {
		w.Header().Set("Content-Type", messageType)
	} else {
		w.Header().Set("Content-Type", "application/json")
	}
	w.WriteHeader(statusCode)
	w.Write(message)
	log.Printf("[api] %s %d %s %s", r.Method, statusCode, r.URL, message)
}

func formatErrorForResponse(e error) []byte {
	msg := "Bad Request"
	if e != nil {
		msg = e.Error()
	}
	m, _ := json.Marshal(errorResponseStruct{Message: msg})
	return []byte(m)
}

// func WSHandler(w http.ResponseWriter, r *http.Request) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	for {
// 		// Read message from browser
// 		msgType, msg, err := conn.ReadMessage()
// 		if err != nil {
// 			return
// 		}

// 		// Print the message to the console
// 		log.Printf("[ws] %s sent: %s\n", conn.RemoteAddr(), string(msg))

// 		// Write message back to browser
// 		if err = conn.WriteMessage(msgType, msg); err != nil {
// 			return
// 		}
// 	}
// }

func RestartHttpServer() {
	HttpServerExitDone.Add(1)
	if Srv != nil {
		log.Println("[load] Restarting Server...")
		ctx, ctxCFn := context.WithTimeout(context.Background(), time.Millisecond*5000)
		defer ctxCFn()
		if err := Srv.Shutdown(ctx); err != nil {
			panic(err) // failure/timeout shutting down the server gracefully
		}
	}
	Srv = startHttpServer()
}

func startHttpServer() *http.Server {
	srv := &http.Server{Addr: ":" + strconv.Itoa(*Cfg.Port), Handler: Router}

	go func() {
		defer HttpServerExitDone.Done() // let main know we are done cleaning up

		log.Printf("\n\n[ready] UI is available at http://localhost:%d\n\n", *Cfg.Port)

		// always returns error. ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// unexpected error. port in use?
			log.Fatalf("ListenAndServe(): %v", err)
		}

	}()

	// returning reference so caller can call Shutdown()
	return srv
}
