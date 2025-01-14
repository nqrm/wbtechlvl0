package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"nqrm/wbtechlvl0/order_services/internal/model"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/twmb/franz-go/pkg/kgo"
)

func main() {
	opts := []kgo.Opt{
		kgo.SeedBrokers("localhost:19092"),
		kgo.DefaultProduceTopic("orders"),
		kgo.ClientID("producer-client-id"),
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		log.Fatal("Error")
		return
	}
	defer client.Close()

	file, err := os.Open("../../model.json")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	ctx := context.Background()

	var order model.Order
	byteValue, _ := io.ReadAll(file)
	err = json.Unmarshal(byteValue, &order)
	if err != nil {
		log.Fatalf(err.Error())
	}

	for {
		id := uuid.New()
		fmt.Println(id)
		order.OrderUID = id.String()
		jsonOrder, err := json.Marshal(order)
		if err != nil {
			log.Printf("JSON Marshall error : %v", err)
		}
		record := &kgo.Record{
			Value: jsonOrder,
		}

		if err := client.ProduceSync(ctx, record).FirstErr(); err != nil {
			log.Fatalf("Produce failed: %v\n", err)
			return
		}

		time.Sleep(5 * time.Second)
	}

}
