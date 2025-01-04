package services

import (
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

/*
Забираем все сообщения из топика
Записываем в бд и кэш только те сообщения, которые соответсвуют типу Order
*/
func (k *KafkaService) Consuming(ctx context.Context) {
	fetches := k.client.PollFetches(ctx)
	if errs := fetches.Errors(); len(errs) > 0 {
		log.Printf("Fetching erros: %v\n", errs)
	}

	for {
		fetches.EachTopic(func(p kgo.FetchTopic) {
			p.EachRecord(func(record *kgo.Record) {
				var order model.Order
				err := json.Unmarshal(record.Value, &order) // проверка, что в топик записали Order
				if err != nil {
					log.Printf("Failed to unmarshal JSON: %v\n", err)
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
