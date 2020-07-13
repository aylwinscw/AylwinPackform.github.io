package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var clientG mongo.Client

var db *sql.DB

type Customer struct {
	User_Id      string `json:"user_id"`
	Login        string `json:"login"`
	Password     string `json:"password"`
	Name         string `json:"name"`
	Company_Id   string `json:"company_id"`
	Credit_Cards string `json:"credit_cards"`
}

type Company struct {
	Company_Id   string `json:"company_id"`
	Company_Name string `json:"company_name"`
}

type Order struct {
	Created_At  string `json:"created_at"`
	Order_Name  string `json:"order_name"`
	Customer_Id string `json:"customer_id"`
}

func getCustomersEndPoint(response http.ResponseWriter, request *http.Request) {

	// Connect to MongoDB
	var customers []Customer
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("content-type", "application/json")

	collection := clientG.Database("customers").Collection("customers")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var customer Customer
		err := cursor.Decode(&customer)
		if err != nil {
			log.Fatal(err)
		}
		customers = append(customers, customer)
	}

	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	json.NewEncoder(response).Encode(customers)
}

func getCompaniesEndPoint(response http.ResponseWriter, request *http.Request) {
	// Connect to MongoDB
	var companies []Company
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("content-type", "application/json")

	collection := clientG.Database("customers").Collection("customer_companies")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var company Company
		err := cursor.Decode(&company)
		if err != nil {
			log.Fatal(err)
		}
		companies = append(companies, company)
	}

	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	json.NewEncoder(response).Encode(companies)
}

func getOrdersEndPoint(response http.ResponseWriter, request *http.Request) {
	var orders []Order
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("content-type", "application/json")

	collection := clientG.Database("customers").Collection("orders")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var order Order
		err := cursor.Decode(&order)
		if err != nil {
			log.Fatal(err)
		}
		orders = append(orders, order)
	}

	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	json.NewEncoder(response).Encode(orders)
}

func main() {

	port := os.Getenv("PORT")

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://dev1:dev1234@customerorders.lkrhm.mongodb.net/customers?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}

	clientG = *client

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/customers", getCustomersEndPoint).Methods("GET")
	router.HandleFunc("/companies", getCompaniesEndPoint).Methods("GET")
	router.HandleFunc("/orders", getOrdersEndPoint).Methods("GET")
	http.ListenAndServe(":"+port, router)

}
