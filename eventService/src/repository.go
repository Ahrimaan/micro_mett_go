package main

import (
	"context"
	"log"
	"os"

	"github.com/mongodb/mongo-go-driver/mongo"
)

type Event struct {
	ID        int    `bson:"_id,omitempty" json:"_id,omitempty"`
	EventName string `bson:"eventname" json:"eventname"`
	Date      string `bson:"date" json:"date"`
}

func GetAllEvents() ([]Event, error) {
	connectionString, present := os.LookupEnv("MONGODB_URL")
	if !present {
		log.Fatal("No URL proiveded, please set MONGODB_URL to your environment")
	}
	var events []Event
	client, err := mongo.NewClient(connectionString)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(nil)

	if err != nil {
		log.Fatal(err)
	}
	c := client.Database("mett").Collection("events")
	cur, err := c.Find(nil, nil, nil)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	for cur.Next(nil) {
		item := Event{}
		err := cur.Decode(&item)
		if err != nil {
			log.Fatal("Decode error ", err)
		}
		events = append(events, item)
	}

	if err := cur.Err(); err != nil {
		log.Fatal("Cursor error ", err)
	}

	return events, nil
}
