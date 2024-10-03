package http_utils

import (
	"fmt"
	"log"
	"net/http"
)

const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Gray    = "\033[37m"
	White   = "\033[97m"
)

func BadRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	writeHeaderResponse(w, "400 Bad Request", r)
}

func ConflictRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusConflict)
	writeHeaderResponse(w, "409 Conflict", r)
}

func OkRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	writeHeaderResponse(w, "200 OK", r)
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	writeHeaderResponse(w, "405 Method Not Allowed", r)
}

func writeHeaderResponse(w http.ResponseWriter, code string, r *http.Request) {
	msg := fmt.Sprint(r.Method, " ", r.URL, " ", r.Proto, " ", code)
	status := code[0]

	switch status {
	case '1':
		log.Print(Cyan + msg + Reset)
	case '2':
		log.Print(Green + msg + Reset)
	case '3':
		log.Print(Blue + msg + Reset)
	case '4':
		log.Print(Red + msg + Reset)
	case '5':
		log.Print(Yellow + msg + Reset)
	default:
		log.Print(Magenta + msg + Reset)
	}
}
