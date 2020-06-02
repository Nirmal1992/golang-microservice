package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

//User is a type that is basically an employee object
type User struct {
	Name        string    `json:"name"`
	ID          string    `json:"id"`
	Birthday    time.Time `json:"birthday"`
	OnboardedOn time.Time `json:"onBoardedOn"`
	IsActive    bool      `json:"isActive"`
}

//AllUsers is collection of users
type AllUsers []User

var birthDay, _ = time.Parse("2006-01-02", "1992-12-12")
var onBoardedOn, _ = time.Parse("2006-01-02", "2020-12-05")

var users = AllUsers{
	{
		Name:        "Nirmal",
		ID:          "342",
		Birthday:    birthDay,
		OnboardedOn: onBoardedOn,
		IsActive:    true,
	},
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(users)
}

func addUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	json.Unmarshal(reqBody, &newUser)
	users = append(users, newUser)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
	fmt.Println(users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	var existingUser User
	userID := mux.Vars(r)["id"]
	for _, u := range users {
		if u.ID == userID {
			existingUser = u
			break
		}
	}
	if (User{}) != existingUser {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(existingUser)
	} else {
		w.WriteHeader(http.StatusNotFound)

	}

	fmt.Println(existingUser)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	var updatedUser User
	userID := mux.Vars(r)["id"]
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	json.Unmarshal(reqBody, &updatedUser)
	for i, u := range users {
		if u.ID == userID {
			users[i] = updatedUser
			break
		}
	}
	if (User{}) != updatedUser {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(updatedUser)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Println(users)
}

func removeUser(w http.ResponseWriter, r *http.Request) {
	var id string
	userID := mux.Vars(r)["id"]
	for i, u := range users {
		if u.ID == userID {
			id = userID
			w.WriteHeader(http.StatusNoContent)
			users = append(users[:i], users[i+1:]...)
			fmt.Fprintf(w, "The user with ID %v has been deleted successfully", userID)
			break
		}
	}
	if id == "" {
		w.WriteHeader(http.StatusNotFound)
	}
	fmt.Println(users)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/users", getUsers).Methods("GET")
	router.HandleFunc("/addUser", addUser).Methods("POST")
	router.HandleFunc("/users/{id}", getUser).Methods("GET")
	router.HandleFunc("/users/{id}", updateUser).Methods("PATCH")
	router.HandleFunc("/users/{id}", removeUser).Methods("DELETE")
	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe("localhost:8080", router))
}
