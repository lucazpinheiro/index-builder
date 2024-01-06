package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"

	"github.com/bxcodec/faker/v3"
)

type Product struct {
	ID          string   `json:"id"`
	Status      string   `json:"status"`
	Name        string   `json:"name"`
	Price       float64  `json:"price"`
	Categories  []string `json:"categories"`
	Description string   `json:"description"`
}

func generateFakeProduct() Product {
	product := Product{
		ID:          faker.UUIDDigit(),
		Status:      getRandomStatus(),
		Name:        faker.Word(),
		Price:       math.Round(rand.Float64()*(900.0-10.0)*100) / 100,
		Categories:  getCategories(),
		Description: faker.Sentence(),
	}
	return product
}

func getCategories() []string {
	categories := []string{}

	numProducts := rand.Intn(5) + 1

	for i := 0; i < numProducts; i++ {
		categories = append(categories, faker.Word())
	}

	return categories
}

func getRandomStatus() string {
	statuses := []string{"available", "unavailable"}
	return statuses[rand.Intn(len(statuses))]
}

func main() {
	rand.Seed(42) // For reproducibility

	numProducts := 500 // Adjust the number of products you want to generate

	for i := 0; i < numProducts; i++ {
		product := generateFakeProduct()
		productJSON, _ := json.Marshal(product)
		fmt.Println(string(productJSON))
	}
}
