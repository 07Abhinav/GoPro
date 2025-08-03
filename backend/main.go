// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"
// )

// type Message struct {
// 	ID int `json:"id"`
// 	Text string `json: "test"`
// }

// func handler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Welcome to the API:")
// }
// func messageHandler(w http.ResponseWriter, r *http.Request) {
// 	messages := []Message{
// 		{ID: 1, Text: "Heelo jangjang"},
// 		{ID: 2, Text: "This is a json file"},
// 	}

// 	w.Header().Set("Content-Type", "application/json")

// 	if err := json.NewEncoder(w).Encode(messages); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// }
// func main(){
// 	http.HandleFunc("/", handler)
// 	http.HandleFunc("/messages", messageHandler)
// 	fmt.Println("Starting server on :8080")

// 	if err := http.ListenAndServe(":8000", nil); err != nil {
// 		log.Fatal(err)
// 	}
// }