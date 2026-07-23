package main
import (
	"fmt"
	"log"
	"github.com/riperaspberry/subscription-service/internal/database"
	"github.com/riperaspberry/subscription-service/internal/config"
)

func main() {
	fmt.Println("Starting subscription service...")

	cfg := config.Load()

	fmt.Println("Port:", cfg.AppPort)

	db, err := database.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	fmt.Println("Connected to database")
}
