package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/abefiker/go_mongodb/models"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type UserController struct {
	session *mongo.Client
}

func NewUserController(s *mongo.Client) *UserController {
	return &UserController{s}
}

func (uc *UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context() // Use context from request

	id := p.ByName("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var u models.User
	err = uc.session.Database("mongo-golang").Collection("users").FindOne(ctx, bson.M{"_id": objectID}).Decode(&u)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		fmt.Println(err) // Handle other errors
		return
	}

	uj, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", uj)
}

func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context() // Use context from request

	var u models.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest) // Handle bad request
		return
	}

	_, err = uc.session.Database("mongo-golang").Collection("users").InsertOne(ctx, u)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError) // Handle internal server error
		return
	}

	uj, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s\n", uj)
}

func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context() // Use context from request

	id := p.ByName("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	result, err := uc.session.Database("mongo-golang").Collection("users").DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError) // Handle internal server error
		return
	}

	if result.DeletedCount == 0 {
		w.WriteHeader(http.StatusNotFound) // Handle document not found
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Deleted user: %s\n", id)
}
