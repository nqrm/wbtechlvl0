package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"nqrm/wbtechlvl0/order_services/internal/model"
	"os"

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

	var order model.Order
	byteValue, _ := io.ReadAll(file)
	err = json.Unmarshal(byteValue, &order)
	if err != nil {
		log.Fatalf(err.Error())
	}

	id := uuid.New()
	order.OrderUID = id.String()

	/*for {
		id = uuid.New()
		order.OrderUID = id.String()

		time.Sleep(5 * time.Second)
	}*/

	jsonOrder, err := json.Marshal(order)
	if err != nil {
		log.Printf("JSON Marshall error : %v", err)
	}
	record := &kgo.Record{
		Value: jsonOrder,
	}

	if err := client.ProduceSync(context.Background(), record).FirstErr(); err != nil { // mb TryProduce
		log.Fatalf("Produce failed: %v\n", err)
		return
	}
}
