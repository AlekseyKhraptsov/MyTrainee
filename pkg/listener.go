package main

import (
	"NatsStreaming/storage"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"os/signal"
	"sync"
)

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)
	var wg sync.WaitGroup

	for {
		wg.Add(1)
		var dat storage.Order

		const (
			HOST     = "localhost"
			DATABASE = "postgres"
			PORT     = "5432"
			USER     = "postgres"
			PASSWORD = "fs"
		)
		var connectionString string = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", HOST, PORT, USER, PASSWORD, DATABASE)
		db, err := sql.Open("postgres", connectionString)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		err = db.Ping()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Successfully created connection to database")

		nc, _ := nats.Connect(nats.DefaultURL)

		nc.Subscribe("orders", func(m *nats.Msg) {

			if err := json.Unmarshal(m.Data, &dat); err != nil {
				log.Fatal(err)

			}

			defer wg.Done()

		})
		wg.Wait()
		sql_statement := "INSERT INTO orders VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11);"
		_, err = db.Exec(sql_statement,
			dat.OrderUid,
			dat.TrackNumber,
			dat.Entry,
			dat.Locale,
			dat.InternalSignature,
			dat.CustomerId,
			dat.DeliveryService,
			dat.Shardkey,
			dat.SmId,
			dat.OofShard,
			dat.DateCreated)
		if err != nil {
			log.Fatal(err)
		}

		sql_statement = "INSERT INTO delivery VALUES ($1,$2,$3,$4,$5,$6,$7,$8);"
		_, err = db.Exec(sql_statement,
			dat.OrderUid,
			dat.Delivery.Name,
			dat.Delivery.Phone,
			dat.Delivery.Zip,
			dat.Delivery.City,
			dat.Delivery.Address,
			dat.Delivery.Region,
			dat.Delivery.Email)
		if err != nil {
			log.Fatal(err)
		}

		sql_statement = "INSERT INTO payment VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11);"
		_, err = db.Exec(sql_statement,
			dat.OrderUid,
			dat.Payment.Transaction,
			dat.Payment.RequestId,
			dat.Payment.Currency,
			dat.Payment.Provider,
			dat.Payment.Amount,
			dat.Payment.PaymentDt,
			dat.Payment.Bank,
			dat.Payment.DeliveryCost,
			dat.Payment.GoodsTotal,
			dat.Payment.CustomFee)
		if err != nil {
			log.Fatal(err)
		}
		for _, item := range dat.Items {
			sql_statement = "INSERT INTO items VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12);"
			_, err = db.Exec(sql_statement,
				dat.OrderUid,
				item.ChrtId,
				item.TrackNumber,
				item.Price,
				item.Rid,
				item.Name,
				item.Sale,
				item.Size,
				item.TotalPrice,
				item.NmId,
				item.Brand,
				item.Status)
			if err != nil {
				log.Fatal(err)
			}
		}

		fmt.Println("Записано сообщение в таблицу с ключем:", dat.OrderUid)
	}
}
