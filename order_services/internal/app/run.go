package app

import (
	"context"
	"log"
	"nqrm/wbtechlvl0/order_services/internal/repository/pgrepo"
	"os"

	"github.com/joho/godotenv"
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
	defer orderDB.Close() // закрытие соединениня

}
