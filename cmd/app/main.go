package main
import (
	"fmt"
	"github.com/riperaspberry/subscription-service/internal/config"
)
func main() {
fmt.Println("Starting subscription service...")
cfg := config.Load()
fmt.Println("Port:", cfg.AppPort)
}
