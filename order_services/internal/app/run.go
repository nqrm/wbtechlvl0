package app

import (
	"context"
	"log"
	"nqrm/wbtechlvl0/order_services/internal/repository/cache"
	"nqrm/wbtechlvl0/order_services/internal/repository/pgrepo"
	"nqrm/wbtechlvl0/order_services/internal/services"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/twmb/franz-go/pkg/kgo"
)

func Run() {
	ctx := context.Background()

	godotenv.Load("../../.env")
	orderDB, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer orderDB.Close()

	postgre := pgrepo.NewPG(orderDB)
	cache := cache.NewCacheStorage()

	// kafka consumer
	opts := []kgo.Opt{
		kgo.SeedBrokers("localhost:19092"),
		kgo.ConsumeTopics("orders"),
		kgo.ClientID("consumer-client-id"),
	}
	kafkaServ := services.NewKafkaService(opts, cache, postgre)
	defer kafkaServ.CloseClient()
	kafkaServ.Consuming(ctx)

	orderService := services.NewOrderService(postgre, cache)
}
