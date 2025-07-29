package main

import (
	"context"
	"log"
	"net/http"

	"google.golang.org/genai"
)

func main() {
	ctx := context.Background()
	config, err := GetConfig()
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	clientConfig := genai.ClientConfig{
		APIKey: config.GeminiAPIKey,
	}
	client, err := genai.NewClient(ctx, &clientConfig)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/search", dishSearch(ctx, *client))
	http.ListenAndServe(":8080", nil)
}
