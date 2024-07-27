package main

import (
	"context"
	"log"
	"time"

	"github.com/shubhamkumar96/go-microservices/logger-service/data"
)

// Declare a type on which we will be creating RPC methods on.
type RPCServer struct{}

// Declare a type for the kind of data that we will receive for
// any methods tied to RPCServer
type RPCPayload struct {
	Name string
	Data string
}

// Declare RPC Methods
// 'resp' is the response that we will be sending back
func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Println("error writing to mongo", err)
		return err
	}

	*resp = "Processed payload via RPC:" + payload.Name
	return nil
}
