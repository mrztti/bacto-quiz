package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const SERVER_PORT = ":4444"

type Quizlet struct {
	Id       uint16
	Question string
	Ref1     uint16
	Ref2     uint16
	Ref3     uint16
	Ref4     uint16
	Option1  string
	Option2  string
	Option3  string
	Option4  string
	Answer   uint16
}

func main() {
	fmt.Printf("Starting bacto-quiz backend on port %s\n", SERVER_PORT)

	//create a new router
	router := mux.NewRouter()

	//specify endpoints, handler functions and HTTP method
	router.HandleFunc("/data/", dataHandler).Methods("GET")
	router.HandleFunc("/data", dataHandler).Methods("GET")
	router.HandleFunc("/quiz/", quizHandler).Methods("GET")
	router.HandleFunc("/quiz", quizHandler).Methods("GET")
	router.HandleFunc("/", homeHandler).Methods("GET")
	http.Handle("/", router)

	//start and listen to requests
	http.ListenAndServe(SERVER_PORT, router)
}

func timeStamp() string {
	return time.Now().Format(time.ANSIC)
}

func createQuizlet() Quizlet {
	qz := Quizlet{0, "This is a test question.", 0, 0, 0, 0, "option", "option", "option", "option", 2}
	return qz
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("entering data end point")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "API is up and running")
}

func quizHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("generate quiz")
	paths := []string{
		"templates/quizlet.tmpl",
	}
	t := template.Must(template.New("quizlet.tmpl").ParseFiles(paths...))
	qz := createQuizlet()
	err := t.Execute(w, qz)
	if err != nil {
		panic(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("welcome home")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, timeStamp())
}
