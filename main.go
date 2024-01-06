package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const sampleDataPath = "sample"

type Indexes struct {
	Prices     map[string][]string
	Categories map[string][]string
}

type Product struct {
	ID          string   `json:"id"`
	Status      string   `json:"status"` // 'available' or 'unavailable'
	Name        string   `json:"name"`
	Price       float64  `json:"price"`
	Categories  []string `json:"categories"`
	Description string   `json:"description"`
}

func sourceData() ([]Product, error) {
	var products []Product

	file, err := os.Open(sampleDataPath)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	fileScanner := bufio.NewScanner(file)

	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		p := Product{}
		byt := fileScanner.Bytes()

		err := json.Unmarshal(byt, &p)
		if err != nil {
			log.Fatal(err)
		}

		products = append(products, p)
	}

	return products, nil
}

func generatePriceIndexes(products []Product, indexObj *Indexes) {
	minRange := 0
	maxRange := 99

	for _, p := range products {
		for p.Price > float64(maxRange) {
			minRange += 100
			maxRange += 100
		}
		priceRange := fmt.Sprintf("%d-%d", minRange, maxRange)
		indexObj.Prices[priceRange] = append(indexObj.Prices[priceRange], p.ID)

		minRange = 0
		maxRange = 99
	}
}

func generateCategoriesIndexes(products []Product, indexObj *Indexes) {
	for _, p := range products {
		for _, c := range p.Categories {
			indexObj.Categories[c] = append(indexObj.Categories[c], p.ID)
		}
	}
}

func writeResult(indexObj *Indexes) {
	f, err := os.Create("indexes")
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	jsonData, err := json.MarshalIndent(indexObj, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	// Write the JSON data to the file
	_, err = f.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}

func main() {
	products, err := sourceData()
	if err != nil {
		log.Fatal(err)
	}

	var indexes = Indexes{
		Prices:     make(map[string][]string),
		Categories: make(map[string][]string),
	}

	generatePriceIndexes(products, &indexes)
	generateCategoriesIndexes(products, &indexes)

	writeResult(&indexes)
}
