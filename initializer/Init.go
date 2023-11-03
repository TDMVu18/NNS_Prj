package initializer

import (
	"context"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

func ConnectEnv() {
	err := godotenv.Load(".env")
	//loi ket noi moi truong
	if err != nil {
		log.Fatalf("không kết nối được đến file môi trường; lỗi: %s", err)
	}
}
func ConnectMongo() *mongo.Client {
	ConnectEnv()
	uri := os.Getenv("MG_CONNECT")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	return client
}

var mongoClient *mongo.Client
var collection *mongo.Collection

func ConnectDB(collname string) *mongo.Collection {
	mongoClient = ConnectMongo()
	collection = mongoClient.Database("personlist").Collection(collname)
	return collection
}

func UserDB() *mongo.Collection {
	mongoClient = ConnectMongo()
	collection = mongoClient.Database("personlist").Collection("user_info")
	return collection
}

func DisconnectDB() {
	if err := mongoClient.Disconnect(context.TODO()); err != nil {
		err.Error()
	}
}
