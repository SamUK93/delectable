package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"google.golang.org/genai"
)

var PROMPT_TEMPLATE = "Give me the names (only the name, not the description or ingredients) of %d dishes that " +
	"use the following ingredients as the base ingredients: %s. Separate the names using the | symbol"

type SearchRequest struct {
	Ingredients []string `json: "ingredients"`
	DishCount   int      `json: "dishCount"`
}

type SearchResponse struct {
	Dishes []string `json: "dishes"`
}

func dishSearch(ctx context.Context, aiClient genai.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var search SearchRequest
		err := json.NewDecoder(r.Body).Decode(&search)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := aiClient.Models.GenerateContent(
			ctx,
			"gemini-2.5-flash",
			genai.Text(fmt.Sprintf(PROMPT_TEMPLATE, search.DishCount, strings.Join(search.Ingredients, ", "))),
			nil,
		)
		if err != nil {
			log.Fatal(err)
		}

		response := SearchResponse{
			Dishes: strings.Split(result.Text(), "|"),
		}

		json.NewEncoder(w).Encode(&response)
	}
}
