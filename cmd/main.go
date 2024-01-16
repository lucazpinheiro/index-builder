package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/lucazpinheiro/indexing-system/internal"
)

const (
	redisAddr      = "localhost:6379"
	sampleDataPath = "sample"
)

func sourceData(saveData func(p internal.Product) (bool, error)) ([]internal.Product, error) {
	var products []internal.Product

	file, err := os.Open(sampleDataPath)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	fileScanner := bufio.NewScanner(file)

	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		p := internal.Product{}
		byt := fileScanner.Bytes()

		err := json.Unmarshal(byt, &p)
		if err != nil {
			log.Fatal(err)
		}

		ok, err := saveData(p)
		if !ok {
			log.Fatal(err)
		}

		products = append(products, p)
	}

	return products, nil
}

func main() {
	db := internal.NewDB(redisAddr)
	defer db.Close()

	products, err := sourceData(db.SaveProduct)
	if err != nil {
		log.Fatal(err)
	}

	index := internal.NewIndex()

	index.MountNameIndex(products)
	index.MountDescriptionIndex(products)
	index.MountPriceIndex(products)
	index.MountCategoriesIndex(products)

	index.WriteResult()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter query: ")
		scanner.Scan()

		query := scanner.Text()

		if len(query) == 0 {
			fmt.Println("Missing query...")
			continue
		}

		var ids []string

		values := strings.Split(query, ":")
		if len(values) != 2 {
			break
		}

		switch strings.ToLower(values[0]) {
		case "price":
			price, err := strconv.ParseFloat(values[1], 64)
			if err != nil {
				log.Fatal(err)
				break
			}
			ids = index.FindProductsByPrice(price)
		case "name":
			ids = index.FindProductsByName(values[1])
		case "categories":
			ids = index.FindProductsByCategory(values[1])
		case "description":
			ids = index.FindProductsByDescription(values[1])
		}

		if len(ids) == 0 {
			fmt.Println("No products found")
		}

		for _, id := range ids {
			p, err := db.GetProduct(id)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(p.ID)
		}
	}

	// handle error
	if scanner.Err() != nil {
		fmt.Println("Error: ", scanner.Err())
	}
}
