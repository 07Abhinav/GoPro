package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"strings"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var messagesCollection *mongo.Collection

type Messages struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Text string `bson:"text" json:"text"`
}

func messageRouter(w http.ResponseWriter, r *http.Request) {
	id := ""
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) > 2 {
		id = parts[2]
	}

	switch r.Method {
	case "GET":
		if id != "" {
			getMessage(w,r,id)
		} else{
			getAllMessage(w,r)
		}
	case "POST":
		createMessage(w,r)
	case "PUT":
		if id == "" {
			http.Error(w, "ID is required for update", http.StatusBadRequest)
			return
		}
		updateMessage(w,r,id)
	case "DELETE":
		if id == "" {
			http.Error(w, "ID is required for delete", http.StatusBadRequest)
			return
		}
		deleteMessage(w,r,id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

//CREATE
func createMessage(w http.ResponseWriter, r *http.Request) {
	var message Messages

	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := messagesCollection.InsertOne(context.TODO(), message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

//GET ALL
func getAllMessage(w http.ResponseWriter, r *http.Request) {
	var messages []Messages
	cursor, err := messagesCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = cursor.All(context.TODO(), &messages); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

//GET ONE
func getMessage(w http.ResponseWriter, r *http.Request, id string) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	var message Messages
	err = messagesCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&message)
	if err != nil {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

//UPDATE
func updateMessage(w http.ResponseWriter, r *http.Request, id string) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	var message Messages
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := messagesCollection.ReplaceOne(context.TODO(), bson.M{"_id": objId}, message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if result.MatchedCount == 0 {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Message updated successfully")
}

//DELETE
func deleteMessage(w http.ResponseWriter, r *http.Request, id string) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	result, err := messagesCollection.DeleteOne(context.TODO(), bson.M{"_id": objId})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if result.DeletedCount == 0 {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Message deleted successfully")
}

// func messageHandlers(w http.ResponseWriter, r *http.Request) {
// 	var messages []Messages

// 	cursor, err := messagesCollection.Find(context.TODO(), bson.M{})
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	defer cursor.Close(context.TODO())

// 	if err = cursor.All(context.TODO(), &messages); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(messages)
// }

func main(){
	err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
	uri := os.Getenv("MONGODB_URI")

	if uri == "" {
		log.Fatal("MONGODB_URI environment variable is not set")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	fmt.Println("Connected to MongoDB!")

	messagesCollection = client.Database("test").Collection("users")

	http.HandleFunc("/messages/", messageRouter)

	fmt.Println("Starting server on :8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}