package app

import (
	"context"
	"log"
	"nqrm/wbtechlvl0/order_services/internal/repository/pgrepo"
	"nqrm/wbtechlvl0/order_services/internal/services"
	"os"

	"github.com/joho/godotenv"
	"github.com/twmb/franz-go/pkg/kgo"
)

func Run() {
	godotenv.Load("../../.env")
	orderDB, err := pgrepo.NewPG(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	err = orderDB.Ping(context.Background())
	if err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}
	defer orderDB.Close()

	// kafka consumer
	opts := []kgo.Opt{
		kgo.SeedBrokers("localhost:19092"),
		kgo.ConsumeTopics("orders"),
		kgo.ClientID("consumer-client-id"),
	}
	kafkaServ, err := services.NewKafkaService(opts)
	if err != nil {
		log.Fatalf("Problems when creating KafkaService: %v", err)
	}
	kafkaServ.StartConsume(context.Background())
	defer kafkaServ.CloseClient()
}
