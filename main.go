package main

import (
	"encoding/json"
	"fmt"
	"github/kerrai1990/phonebook-rest-api/data"
	"github/kerrai1990/phonebook-rest-api/middleware"
	"github/kerrai1990/phonebook-rest-api/models"
	"github/kerrai1990/phonebook-rest-api/responses"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var people []models.Person
var db *data.Client

// Main function
func main() {
	db = data.New()
	db.CheckStatus()

	withMiddleware := middleware.ChainMiddleware(middleware.WithContentType, middleware.WithLogging)

	router := mux.NewRouter()
	router.HandleFunc("/people", withMiddleware(index)).Methods("GET")
	router.HandleFunc("/people/{id}", withMiddleware(show)).Methods("GET")
	router.HandleFunc("/people", withMiddleware(create)).Methods("POST")
	router.HandleFunc("/people/{id}", withMiddleware(update)).Methods("PATCH")
	router.HandleFunc("/people/{id}", withMiddleware(delete)).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}

// Helper Functions
func nextID() int {
	currentID := 0
	if len(people) != 0 {
		currentID = people[len(people)-1].ID
	}
	return currentID + 1
}

func errorResponse(err error, code int, w http.ResponseWriter) {
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(responses.JSONErrors{
		Code:    code,
		Message: err.Error(),
	})
}

// Rest Functions
func index(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(people)

}

func show(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	person, err := db.Get(id)
	if err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(person)
}

func create(w http.ResponseWriter, r *http.Request) {
	var person models.Person

	//try and save new contact
	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		errorResponse(err, http.StatusInternalServerError, w)
		return
	}

	if err := person.ValidatePerson(); err != nil {
		errorResponse(err, http.StatusUnprocessableEntity, w)
		return
	}

	person.ID = nextID()
	newPerson, err := db.Create(person)
	if err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(newPerson)
}

func update(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person models.Person

	_ = json.NewDecoder(r.Body).Decode(&person)
	id, _ := strconv.Atoi(params["id"])

	personToUpdate, err := db.Get(id)
	if err != nil {
		errorResponse(err, http.StatusNotFound, w)
		return
	}

	fmt.Println(personToUpdate)

	updatedPerson, err := db.Update(personToUpdate)
	if err != nil {
		panic(err)
	}

	// json.NewEncoder(w).Encode(personToUpdate)
}

func delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	for index, item := range people {
		if item.ID == id {
			people = append(people[:index], people[index+1:]...)
			break
		}
		json.NewEncoder(w).Encode(people)
		return
	}
}
