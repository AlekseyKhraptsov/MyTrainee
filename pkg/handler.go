package main

import (
	"NatsStreaming/storage"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

var db sql.DB

func getOne(w http.ResponseWriter, r *http.Request) {
	const (
		HOST     = "localhost"
		DATABASE = "postgres"
		PORT     = "5432"
		USER     = "postgres"
		PASSWORD = "fs"
	)
	params := mux.Vars(r)
	var id string
	for _, v := range params {
		id = v
	}

	connection := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", HOST, PORT, USER, PASSWORD, DATABASE)
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}
	//_ = json.NewDecoder(r.Body).Decode(&id)
	var result storage.Order

	query := fmt.Sprintf("SELECT * FROM orders WHERE order_uid = '%v'", id)

	err = db.QueryRow(query).Scan(&result.OrderUid,
		&result.TrackNumber,
		&result.Entry,
		&result.Locale,
		&result.InternalSignature,
		&result.CustomerId,
		&result.DeliveryService,
		&result.Shardkey,
		&result.SmId,
		&result.OofShard,
		&result.DateCreated)
	if err != nil {
		fmt.Println(err.Error())
	}

	//get delivery by ID
	query = fmt.Sprintf("SELECT * FROM delivery WHERE order_uid='%v'", id)

	err = db.QueryRow(query).Scan(&result.OrderUid,
		&result.Delivery.Name,
		&result.Delivery.Phone,
		&result.Delivery.Zip,
		&result.Delivery.City,
		&result.Delivery.Address,
		&result.Delivery.Region,
		&result.Delivery.Email)
	if err != nil {
		fmt.Println(err.Error())
	}

	//get payment by transaction
	query = fmt.Sprintf("SELECT * FROM payment WHERE order_uid='%v'", id)
	err = db.QueryRow(query).Scan(&result.OrderUid,
		&result.Payment.Transaction,
		&result.Payment.RequestId,
		&result.Payment.Currency,
		&result.Payment.Provider,
		&result.Payment.Amount,
		&result.Payment.PaymentDt,
		&result.Payment.Bank,
		&result.Payment.DeliveryCost,
		&result.Payment.GoodsTotal,
		&result.Payment.CustomFee)
	if err != nil {
		fmt.Println(err.Error())
	}

	//get each item by []ids

	var result2 storage.ItemsDB
	query = fmt.Sprintf("SELECT * FROM items WHERE order_uid='%v'", id)
	err = db.QueryRow(query).Scan(&result.OrderUid,
		&result2.ChrtId,
		&result2.TrackNumber,
		&result2.Price,
		&result2.Rid,
		&result2.Name,
		&result2.Sale,
		&result2.Size,
		&result2.TotalPrice,
		&result2.NmId,
		&result2.Brand,
		&result2.Status)
	if err == nil {
		result.Items = append(result.Items, result2)

	}

	fmt.Println("this is sparta", result2.ChrtId)
	json.NewEncoder(w).Encode(result)
}
func getOrder(w http.ResponseWriter, r *http.Request) { //get one...
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var dat storage.Order
	for _, order := range dat.OrderUid {
		if dat.OrderUid == params["order_uid"] {
			json.NewEncoder(w).Encode(order)
			return
		}
		json.NewEncoder(w).Encode(&storage.Order{})
	}

}

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/{id}", getOne).Methods("GET")

	log.Fatal(http.ListenAndServe("localhost:8080", r))

}
