package services

import (
	"context"
	"fmt"
	"log"

	"github.com/twmb/franz-go/pkg/kgo"
)

type KafkaService struct {
	client *kgo.Client
}

func NewKafkaService(opts []kgo.Opt) (*KafkaService, error) {
	client, err := kgo.NewClient(opts...)
	if err != nil {
		log.Fatal("Error")
		return nil, err
	}
	return &KafkaService{client: client}, nil
}

func (k *KafkaService) StartConsume(ctx context.Context) {
	for {
		fetches := k.client.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			log.Fatalf("Fetching erros: %v", errs)
		}

		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()
			fmt.Println(string(record.Value))
		}
	}
}

func (k *KafkaService) CloseClient() {
	k.client.Close()
}
