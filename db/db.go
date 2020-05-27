package db

import (
	"context"
	"fmt"
	"github.com/semyon-dev/znai-krai/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

func Connect() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb+srv://"+config.MongoDBLogin+":"+config.MongoDBPass+"@main-h6nko.mongodb.net/test?retryWrites=true&w=majority",
	))
	if err != nil {
		log.Fatal(err)
	}

	if client.Ping(ctx, readpref.Primary()) == nil {
		fmt.Println("✔ Подключение MongoDB успешно (ping) ")
	} else {
		fmt.Println("× Подключение к MongoDB не удалось:", err)
	}
}
