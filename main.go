package main

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	password = "password"
	dbname   = "clayte"
	username = "clay"
	//transfer_collections = "transfer"
	collections = "clayte"
)

type User struct {
	Name     string `json:"name"`
	Sex      string `json:"sex"`
	Interest string `json:"interest"`
	Age      int    `json:"age"`
	//Token     string `json:"token"`
}

type ResponseResult struct {
	ResponseCode int    `json:"rc"`
	Error        string `json:"error"`
	Result       string `json:"result"`
}

func GetDBCollection() (*mongo.Collection, error) {
	log.SetFormatter(&log.JSONFormatter{})
	// client, err := mongo.NewClient((options.Client().ApplyURI("mongodb://localhost:27017")))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb://"+username+":"+password+"@localhost:27017",
		//"mongodb+srv://standard:"+password+"@cluster0.pdpui.mongodb.net/"+dbname+"?retryWrites=true&w=majority",
	))

	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	//Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	collection := client.Database(dbname).Collection(collections)
	return collection, nil
}

func insert_after_query(collection *mongo.Collection) bool {
	var result User = User{Name: "clayte", Sex: "man", Interest: "Baseball", Age: 28}
	var res User
	err := collection.FindOne(context.TODO(), bson.D{{Key: "name", Value: result.Name}}).Decode(&res)
	if err != nil {
		_, err = collection.InsertOne(context.TODO(), bson.D{{Key: "name", Value: result.Name}, {Key: "sex", Value: result.Sex}, {Key: "age", Value: result.Age}})
		if err != nil {
			return false
		}
	}
	return true
}

func query_single(collection *mongo.Collection) bool {
	var result User = User{Name: "clayer", Sex: "man"}
	filter := bson.D{{Key: "name", Value: result.Name}, {Key: "sex", Value: result.Sex}}
	opts := options.Find().SetSort(bson.D{{Key: "age", Value: -1}}).SetLimit(1)
	cur, err := collection.Find(context.TODO(), filter, opts)
	if err != nil {
		log.Fatal(err)
	}
	var results []User

	if err = cur.All(context.Background(), &results); err != nil {
		log.Fatal(err)
	}

	for _, val := range results {
		fmt.Println(val)
	}
	return true
}

func update(collection *mongo.Collection) bool {
	var result User = User{Name: "clayer", Sex: "man", Interest: "Table"}
	filter := bson.D{{Key: "name", Value: result.Name}, {Key: "sex", Value: result.Sex}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "interest", Value: result.Interest}}}}

	updateres, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(updateres.ModifiedCount)
	return true
}

func main() {
	collection, err := GetDBCollection()
	if err != nil {
		return
	}
	//	insert_after_query(collection)
	//query_single(collection)
	update(collection)

	fmt.Println("123")
}
