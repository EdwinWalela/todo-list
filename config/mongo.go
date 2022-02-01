package mongoDriver

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() *mongo.Client {
	godotenv.Load()
	MONGO_URL := os.Getenv("MONGO_DB_URL")

	client, err := mongo.NewClient(options.Client().ApplyURI(MONGO_URL))

	if err != nil {
		log.Fatal("Unable to connect to Mongo", err)
		return nil
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)

	if err := client.Connect(ctx); err != nil {
		log.Fatal("Unable to connect to Mongo")
		return nil
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil
	}

	log.Println("Connected to MongoDB")

	return client
}

func GetCollection(client *mongo.Client, collection string) *mongo.Collection {
	return client.Database("demo").Collection(collection)
}
