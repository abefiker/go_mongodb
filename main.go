package main

import (
	"context"
	"net/http"

	"github.com/abefiker/go_mongodb/controllers"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background() // Create a context

	r := httprouter.New()
	client, err := getSession(ctx)
	if err != nil {
		panic(err)
	}
	uc := controllers.NewUserController(client)
	r.GET("/user/:id", uc.GetUser)
	r.POST("/user", uc.CreateUser)
	r.DELETE("/user/:id", uc.DeleteUser)
	http.ListenAndServe(":8080", r)
}

func getSession(ctx context.Context) (*mongo.Client, error) {
	// Replace with your actual connection URI (remove "mongodb://")
	uri := "mongodb+srv://go-mongodb:MongoDb@be1994@cluster0.sxvbxay.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	// Test the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}
