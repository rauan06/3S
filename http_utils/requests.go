package http_utils

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"

	. "triples/bucket_struct"
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
	writeXML(w, "400 Bad Request")
	writeHeaderResponse("400 Bad Request", r)
}

func ConflictRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusConflict)
	writeXML(w, "409 Conflict")
	writeHeaderResponse("409 Conflict", r)
}

func NoContentRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
	writeXML(w, "204 No Content")
	writeHeaderResponse("204 No Content", r)
}

func NotFoundRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	writeXML(w, "404 Not Found")
	writeHeaderResponse("404 Not Found", r)
}

func OkRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	writeXML(w, "200 OK")
	writeHeaderResponse("200 OK", r)
}

func OkRequestWithHeaders(w http.ResponseWriter, r *http.Request) {
	if len(AllBuckets[SessionUser.UserID]) != 0 {
		w.Header().Add("Location", r.URL.Path)
		w.Header().Add("Connection", "close")
		w.Header().Add("Server", "triple-s")
	}
	w.WriteHeader(http.StatusOK)
	writeHeaderResponse("200 OK", r)
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	writeXML(w, "405 Method Not Allowed")
	writeHeaderResponse("405 Method Not Allowed", r)
}

func writeHeaderResponse(code string, r *http.Request) {
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

func writeXML(w http.ResponseWriter, msg string) {
	var nf Response = Response(msg)
	out, _ := xml.MarshalIndent(nf, "  ", "  ")
	fmt.Fprint(w, xml.Header)
	fmt.Fprintln(w, string(out))
}
