package services

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"nqrm/wbtechlvl0/order_services/internal/model"
	"nqrm/wbtechlvl0/order_services/internal/repository"

	"github.com/twmb/franz-go/pkg/kgo"
)

type KafkaService struct {
	client *kgo.Client
	cache  repository.CacheOrder
	db     repository.OrderDB
}

func NewKafkaService(opts []kgo.Opt, cache repository.CacheOrder, db repository.OrderDB) *KafkaService {
	client, err := kgo.NewClient(opts...)
	if err != nil {
		log.Fatalf("Kafka client creation error %v\n", err)
	}
	return &KafkaService{client, cache, db}
}

func (k *KafkaService) Consuming(ctx context.Context) {
	for {
		fetches := k.client.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			log.Printf("Fetching erros: %v\n", errs)
		}

		fetches.EachTopic(func(p kgo.FetchTopic) {
			p.EachRecord(func(record *kgo.Record) {
				orderData := record.Value
				ok := json.Valid(orderData)
				if !ok {
					log.Printf("Producer message is not valid in JSON format: %v\n", string(orderData))
					return
				}
				dec := json.NewDecoder(bytes.NewReader(orderData))
				dec.DisallowUnknownFields() // проверка что сообщение из топика содержит все поля Order

				var order model.Order
				if err := dec.Decode(&order); err != nil {
					log.Printf("Failed to decode Order:", err)
					return
				}
				k.db.AddOrder(ctx, order)
				k.cache.Set(&order)
			})
		})
	}
}

func (k *KafkaService) CloseClient() {
	k.client.Close()
}
