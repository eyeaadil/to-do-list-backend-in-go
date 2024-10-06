package config

import (
    "context"
    "log"
    "time"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

// ConnectDB establishes a connection to the MongoDB database.
func ConnectDB() {
    // Set client options
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
    // Connect to MongoDB
    client, err := mongo.NewClient(clientOptions)
    if err != nil {
        log.Fatal("Error creating MongoDB client:", err)
    }

    // Set a timeout context for the connection
    ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
    defer cancel()

    err = client.Connect(ctx)
    if err != nil {
        log.Fatal("Error connecting to MongoDB:", err)
    }

    // Check the connection
    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }

    // Get a reference to the database
    DB = client.Database("to_do_db")
    log.Println("Connected to mongoDb")
}
