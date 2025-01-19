package utils

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Db *mongo.Client
var UserDb *mongo.Collection

func ConnectDB() {
	var err error
	mongoURI := os.Getenv("MONGO_URI")

	Db, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	err = Db.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB!")
	InitUsersCollection()
}

func InitUsersCollection() {
	UserDb = Db.Database("streamit").Collection("users")
}

