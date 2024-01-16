package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type Index struct {
	Name        map[string][]string `json:"name"`
	Description map[string][]string `json:"description"`
	Price       map[string][]string `json:"price"`
	Categories  map[string][]string `json:"categories"`
}

func NewIndex() *Index {
	return &Index{
		Name:        make(map[string][]string),
		Description: make(map[string][]string),
		Price:       make(map[string][]string),
		Categories:  make(map[string][]string),
	}
}

func (index *Index) WriteResult() {
	f, err := os.Create("indexes")
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	jsonData, err := json.MarshalIndent(index, "", "  ")
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

func (index *Index) MountPriceIndex(products []Product) {
	minRange := 0
	maxRange := 99

	for _, p := range products {
		for p.Price > float64(maxRange) {
			minRange += 100
			maxRange += 100
		}
		priceRange := fmt.Sprintf("%d-%d", minRange, maxRange)
		index.Price[priceRange] = append(index.Price[priceRange], p.ID)

		minRange = 0
		maxRange = 99
	}
}

func (index *Index) MountCategoriesIndex(products []Product) {
	for _, p := range products {
		for _, c := range p.Categories {
			index.Categories[c] = append(index.Categories[c], p.ID)
		}
	}
}

func parseName(name string) []string {
	return strings.Split(name, " ")
}

func (index *Index) MountNameIndex(products []Product) {
	for _, p := range products {
		for _, s := range parseName(p.Name) {
			index.Name[s] = append(index.Name[s], p.ID)
		}
	}
}

func parseDescription(description string) []string {
	return strings.Split(description, " ")
}

func (index *Index) MountDescriptionIndex(products []Product) {
	for _, p := range products {
		for _, s := range parseDescription(p.Description) {
			index.Description[s] = append(index.Description[s], p.ID)
		}
	}
}

func (index *Index) FindProductsByPrice(price float64) []string {
	minPrice := 0
	maxPrice := 99
	priceRange := "0-99"

	for {
		if price >= float64(minPrice) && price <= float64(maxPrice) {
			break
		}
		maxPrice += 100
		minPrice += 100
		priceRange = fmt.Sprintf("%d-%d", minPrice, maxPrice)
	}

	ids, ok := index.Price[priceRange]
	if !ok {
		return []string{}
	}

	return ids
}

func (index *Index) FindProductsByName(name string) []string {
	ids, ok := index.Name[name]
	if !ok {
		return []string{}
	}

	return ids
}

func (index *Index) FindProductsByCategory(category string) []string {
	ids, ok := index.Categories[category]
	if !ok {
		return []string{}
	}

	return ids
}

func (index *Index) FindProductsByDescription(term string) []string {
	ids, ok := index.Description[term]
	if !ok {
		return []string{}
	}

	return ids
}
