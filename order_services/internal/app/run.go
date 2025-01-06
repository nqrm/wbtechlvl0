package app

import (
	"context"
	"log"
	"net/http"
	controller "nqrm/wbtechlvl0/order_services/internal/controller/http"
	"nqrm/wbtechlvl0/order_services/internal/repository/cache"
	"nqrm/wbtechlvl0/order_services/internal/repository/pgrepo"
	"nqrm/wbtechlvl0/order_services/internal/services"
	"os"
	"sync"

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

	orderService := services.NewOrderService(postgre, cache)
	err = orderService.Recovery(ctx)
	if err != nil {
		log.Printf("Unable to recovery cache", err)
	}
	orderRouter := controller.NewOrderRouter(orderService)
	srv := &http.Server{
		Handler: orderRouter,
		Addr:    "localhost:8000",
	}

	// kafka consumer
	opts := []kgo.Opt{
		kgo.SeedBrokers("localhost:19092"),
		kgo.ConsumeTopics("orders"),
		kgo.ClientID("consumer-client-id"),
	}
	kafkaServ := services.NewKafkaService(opts, cache, postgre)
	defer kafkaServ.CloseClient()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		log.Printf("kafka consumer starting\n")
		kafkaServ.Consuming(ctx)
		defer wg.Done()
	}()

	go func() {
		log.Printf("http server listening on %s\n", srv.Addr)
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatalf("error listening and serving: %s\n", err)
		}
		defer wg.Done()
	}()

	wg.Wait()
}
