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

var GENAI_MODEL = "gemini-2.5-flash"

var SEARCH_PROMPT_TEMPLATE = "Give me the names (only the name, not the description or ingredients) of %d dishes that " +
	"use the following ingredients as the base ingredients: %s. Separate the names using the | symbol"

var RESTAURANT_PROMPT_TEMPLATE = `Find a restaurant that serves %s in or near %s.
	format your response like this, replacing the placeholders with the relevant information.
	If you can't find any particular field, replace it with "None found":
	[restaurant name]|[website]|[telephone no.]|[address]|[bool indicating whether takes reservations]`

type SearchRequest struct {
	Ingredients []string `json: "ingredients"`
	DishCount   int      `json: "dishCount"`
	Location    string   `json: "location`
}

type SearchResponse struct {
	Dishes []DishWithInfo `json: "dishes"`
}

type RestaurantInfo struct {
	Name         string `json: "name"`
	Website      string `json: "website"`
	Telephone    string `json: "telephone"`
	Address      string `json: "address"`
	Reservations string `json: "reservations"`
}

type DishWithInfo struct {
	Name        string `json: "name`
	WhereToFind RestaurantInfo
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
			GENAI_MODEL,
			genai.Text(fmt.Sprintf(SEARCH_PROMPT_TEMPLATE, search.DishCount, strings.Join(search.Ingredients, ", "))),
			nil,
		)
		if err != nil {
			log.Fatal(err)
		}

		var responseDishes []DishWithInfo

		dishes := strings.Split(result.Text(), "|")
		for dish := range dishes {
			restaurant := findRestaurantWithDish(ctx, aiClient, dishes[dish], search.Location)
			responseDishes = append(responseDishes, DishWithInfo{
				Name:        dishes[dish],
				WhereToFind: restaurant,
			})
		}

		response := SearchResponse{
			Dishes: responseDishes,
		}

		json.NewEncoder(w).Encode(&response)
	}
}

func findRestaurantWithDish(ctx context.Context, aiClient genai.Client, dishName string, location string) RestaurantInfo {
	result, err := aiClient.Models.GenerateContent(
		ctx,
		GENAI_MODEL,
		genai.Text(fmt.Sprintf(RESTAURANT_PROMPT_TEMPLATE, dishName, location)),
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	resultFields := strings.Split(result.Text(), "|")
	return RestaurantInfo{
		Name:         resultFields[0],
		Website:      resultFields[1],
		Telephone:    resultFields[2],
		Address:      resultFields[3],
		Reservations: resultFields[4],
	}
}
