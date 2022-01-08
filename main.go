package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

// Get Random Number - Helper function
func getRandomIndex(max int) int {
	min := 0
	random := rand.Intn(max-min) + min
	return random
}

// Dino Data
var dinos []Dino

// Dino struct (Model)
type Dino struct {
	Name        string `json:"Name"`
	Description string `json:"Description"`
}

// Name Model
type Name struct {
	Name string `json:"Name"`
}

// Description Model
type Description struct {
	Description string `json:"Description"`
}

// Get All Dinos
func getAllDinosaurs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dinos)
}

// Get Random Dino
func getRandomDinosaur(w http.ResponseWriter, r *http.Request) {
	index := getRandomIndex(len(dinos))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dinos[index])
}

// Get Random Dino Name
func getRandomDinosaurName(w http.ResponseWriter, r *http.Request) {
	index := getRandomIndex(len(dinos))
	var dinoName Name
	dinoName.Name = dinos[index].Name

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dinoName)
}

// Get Random Dino Description
func getRandomDinosaurDescription(w http.ResponseWriter, r *http.Request) {
	index := getRandomIndex(len(dinos))
	var dinoDesc Description
	dinoDesc.Description = dinos[index].Description

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dinoDesc)
}

// spaHandler implements the http.Handler interface, so we can use it
// to respond to HTTP requests. The path to the static directory and
// path to the index file within that static directory are used to
// serve the SPA in the given static directory.
type spaHandler struct {
	staticPath string
	indexPath  string
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func main() {
	// Init Router
	r := mux.NewRouter()

	// Open our jsonFile
	jsonFile, err := os.Open("dinosaurs.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened dinosaurs.json")
	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)
	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &dinos)

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// Route Handlers / Endpoints
	r.HandleFunc("/dinosaurs", getAllDinosaurs).Methods("GET")
	r.HandleFunc("/dinosaurs/random", getRandomDinosaur).Methods("GET")
	r.HandleFunc("/dinosaurs/random/name", getRandomDinosaurName).Methods("GET")
	r.HandleFunc("/dinosaurs/random/description", getRandomDinosaurDescription).Methods("GET")

	// Home Page
	spa := spaHandler{staticPath: "build", indexPath: "index.html"}
	r.PathPrefix("/").Handler(spa)

	log.Fatal(http.ListenAndServe(":8080", r))
}
