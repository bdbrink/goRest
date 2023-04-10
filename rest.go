package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	googlesearch "github.com/rocketlaunchr/google-search"
)

// credentials to interact with server
type User struct {
	Username string `json:"user"`
	Password string `json:"password"`
}

var user User

// may need to open port otherwise server will fail to server traffic
var PORT = ":1234"
var DATA = make(map[string]string)

// default response from server, should be only top level domain, anything else give a 404
func defaultHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("Serving:", r.URL.Path, "from", r.Host)
	w.WriteHeader(http.StatusNotFound)
	Body := "Thanks for visiting! \n"
	fmt.Fprintf(w, "%s", Body)
}

// handle /time to the server
func timeHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("Serving:", r.URL.Path, "from", r.Host)
	t := time.Now().Format(time.RFC1123)
	Body := "The current time is: " + t + "\n"
	fmt.Fprintf(w, "%s", Body)
}

// add user to the server
func addHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("Serving:", r.URL.Path, "from", r.Host, r.Method)
	if r.Method != http.MethodPost {
		http.Error(w, "Error:", http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "%s\n", "Method not allowed!")
		return
	}

	d, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error:", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(d, &user)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error:", http.StatusBadRequest)
		return
	}

	// need a valid username
	if user.Username != "" {
		DATA[user.Username] = user.Password
		log.Println(DATA)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Error:", http.StatusBadRequest)
		return
	}
}

// return the info from the server, with valid credentials only
func getHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("Serving:", r.URL.Path, "from", r.Host, r.Method)
	if r.Method != http.MethodGet {
		http.Error(w, "Error:", http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "%s \n", "Method not Allowed!")
		return
	}

	d, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "ReadAll - Error", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(d, &user)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unmarshal - Error", http.StatusBadRequest)
		return
	}
	fmt.Println(user)

	_, ok := DATA[user.Username]
	if ok && user.Username != "" {
		log.Println("Found!")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s \n", d)
	} else {
		log.Println("Not Found!")
		w.WriteHeader(http.StatusNotFound)
		http.Error(w, "Map - Resource not found!", http.StatusNotFound)
	}
	return
}

// delete a mapped user/password
func deleteHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("Serving:", r.URL.Path, "from", r.Host, r.Method)
	if r.Method != http.MethodDelete {
		http.Error(w, "Error:", http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "%s \n", "Method not allowed!")
		return
	}

	d, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "ReadAll - Error", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(d, &user)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unmarshal - Error", http.StatusBadRequest)
		return
	}
	log.Println(user)

	_, ok := DATA[user.Username]
	if ok && user.Username != "" {
		if user.Password == DATA[user.Username] {
			delete(DATA, user.Username)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "%s \n", d)
			log.Println(DATA)
		}
	} else {
		log.Println("User", user.Username, "Not found!")
		w.WriteHeader(http.StatusNotFound)
		http.Error(w, "Delete - Resource not found!", http.StatusNotFound)
	}
	log.Println("After:", DATA)
	return
}

// if we receieve a bad request kill the server
func errorHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("Serving:", r.URL.Path, "from", r.Host, r.Method)
	Body := "Error with request terminating server \n"
	fmt.Fprintf(w, "%s", Body)
	panic("bad request")
}

// display top result from dank memes
func memeHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("Serving:", r.URL.Path, "from", r.Host, r.Method)
	Body := "Searching dank memes \n"
	fmt.Fprintf(w, "%s", Body)

	searchResults, _ := googlesearch.Search(nil, "dank memes")
	resultTop := searchResults[0].Description

	fmt.Fprintln(w, resultTop)
}

func main() {

	args := os.Args
	if len(args) != 1 {
		PORT = ":" + args[1]
	}

	mux := http.NewServeMux()
	s := &http.Server{
		Addr:         PORT,
		Handler:      mux,
		IdleTimeout:  10 * time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	}

	mux.Handle("/time", http.HandlerFunc(timeHandler))
	mux.Handle("/add", http.HandlerFunc(addHandler))
	mux.Handle("/get", http.HandlerFunc(getHandler))
	mux.Handle("/delete", http.HandlerFunc(deleteHandler))
	mux.Handle("/error", http.HandlerFunc(errorHandler))
	mux.Handle("/", http.HandlerFunc(defaultHandler))
	mux.Handle("/meme", http.HandlerFunc(memeHandler))

	fmt.Println("Ready to serve at", PORT)
	err := s.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		return
	}
}
