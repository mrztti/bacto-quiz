package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

const SERVER_PORT = ":80"

type Quizlet struct {
	Id            uint16
	ChosenProject string
	Question      string
	Option1       string
	Option2       string
	Option3       string
	Option4       string
	Answer        int
}

type Quizzy struct {
	Link     string
	Name     string
	Property string
}

func main() {
	fmt.Printf("Starting bacto-quiz backend on port %s\n", SERVER_PORT)

	//create a new router
	router := mux.NewRouter()

	//specify endpoints, handler functions and HTTP method
	router.HandleFunc("/quiz/", quizHandler).Methods("GET")
	router.HandleFunc("/quiz", quizHandler).Methods("GET")
	router.HandleFunc("/", homeHandler).Methods("GET")
	http.Handle("/", router)
	fmt.Println("Set handlers...")
	//start and listen to requests
	fmt.Println("Listening now!")
	http.ListenAndServe(SERVER_PORT, router)

}

func timeStamp() string {
	return time.Now().Format(time.ANSIC)
}

func readQuizzy(prefix string, href string) Quizzy {
	file, err := os.Open(prefix + href)
	if err != nil {
		panic(err)
	}
	fileScanner := bufio.NewScanner(file)
	lineCount := 0
	var lines []string
	for fileScanner.Scan() {
		lineCount++
		lines = append(lines, fileScanner.Text())
	}

	lineTry := rand.Intn(lineCount)
	for len(lines[lineTry]) <= 1 {
		lineTry = rand.Intn(lineCount)
	}

	qz := Quizzy{href, href[:(len(href) - 3)], lines[lineTry]}
	return qz

}

func createQuizlet() Quizlet {
	arcs, err := ioutil.ReadDir("archives/")
	if err != nil {
		panic(err)
	}
	qid := len(arcs)

	projs := getProjects()
	fmt.Println(len(projs))
	r := rand.Intn(len(projs))
	rd_proj := projs[r].Namespace
	files, err := ioutil.ReadDir("data/" + rd_proj + "/")
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UnixNano())
	candidates := rand.Perm(len(files))
	target := rand.Intn(4)

	qz := Quizlet{uint16(qid), rd_proj, "This is a test question.", "option", "option", "option", "option", target}
	refMap := make(map[int]int)

	q := readQuizzy("data/"+rd_proj+"/", files[candidates[0]].Name())
	refMap[0] = candidates[0]
	qz.Option1 = q.Name

	q = readQuizzy("data/"+rd_proj+"/", files[candidates[1]].Name())
	refMap[1] = candidates[1]
	qz.Option2 = q.Name

	q = readQuizzy("data/"+rd_proj+"/", files[candidates[2]].Name())
	refMap[2] = candidates[2]
	qz.Option3 = q.Name

	q = readQuizzy("data/"+rd_proj+"/", files[candidates[3]].Name())
	refMap[3] = candidates[3]
	qz.Option4 = q.Name

	q = readQuizzy("data/"+rd_proj+"/", files[target].Name())
	qz.Question = q.Property
	qz.Answer = refMap[target+1]
	fmt.Println(target)
	fmt.Println(qz.Answer)

	file, err := json.MarshalIndent(qz, "", " ")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("archives/"+fmt.Sprint(qid)+".json", file, 0666)
	if err != nil {
		panic(err)
	}

	return qz
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

type HomeData struct {
	WelcomeMessage string
	Projects       []Project
}

type Project struct {
	Namespace string
	Objects   []DataObject
}

type DataObject struct {
	Name       string
	Properties []Property
}

type Property struct {
	Value string
}

func getProjects() []Project {
	files, err := ioutil.ReadDir("./data/")
	if err != nil {
		panic(err)
	}
	var projects []Project

	for _, f := range files {
		projects = append(projects, Project{f.Name(), getObjects(f.Name())})

	}

	return projects
}

func getObjects(namespace string) []DataObject {
	objects, err := ioutil.ReadDir("./data/" + namespace)
	if err != nil {
		panic(err)
	}
	var objs []DataObject

	for _, f := range objects {
		objs = append(objs, DataObject{f.Name()[:(len(f.Name()) - 3)], getProperties("./data/" + namespace + "/" + f.Name())})

	}

	return objs
}

func getProperties(href string) []Property {
	file, err := os.Open(href)
	if err != nil {
		panic(err)
	}
	fileScanner := bufio.NewScanner(file)
	lineCount := 0
	var lines []Property
	for fileScanner.Scan() {
		lineCount++
		l := fileScanner.Text()
		if len(l) <= 1 {
			continue
		}
		lines = append(lines, Property{l})
	}
	return lines
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("welcome home")
	home := HomeData{}
	home.WelcomeMessage = "Available topics"
	home.Projects = getProjects()

	tmpl := template.Must(template.ParseFiles("templates/home.tmpl"))
	tmpl.Execute(w, home)
}
