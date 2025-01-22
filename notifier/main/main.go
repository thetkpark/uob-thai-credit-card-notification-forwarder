package main

import (
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	_ "github.com/thetkpark/uob-thai-credit-card-notification-forwarder/notifier"
	"log"
	"os"
)

func main() {
	port := "8080"
	os.Setenv("FUNCTION_TARGET", "Notifier")
	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}
