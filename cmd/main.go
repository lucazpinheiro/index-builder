package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/lucazpinheiro/index-seeker/internal"
)

const (
	redisAddr      = "localhost:6379"
	sampleDataPath = "data/sample"
)

func main() {
	db := internal.NewDB(redisAddr)
	defer db.Close()

	products, err := extractData(sampleDataPath, db)
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
		case "id":
			p, err := db.GetProductByID(values[1])
			if err != nil {
				log.Fatal(err)
			}
			prettyPrint(&p)
			continue
		}

		if len(ids) == 0 {
			fmt.Println("No products found")
		}

		for _, id := range ids {
			p, err := db.GetProductByID(id)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(p.ID)
		}
	}

	if scanner.Err() != nil {
		fmt.Println("Error: ", scanner.Err())
	}
}

func prettyPrint(p *internal.Product) {
	prettyProduct, err := json.MarshalIndent(p, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(prettyProduct))
}

func extractData(sourceData string, destiny *internal.DB) ([]internal.Product, error) {
	log.Printf("extracting products data, source: %s", sourceData)

	file, err := os.Open(sourceData)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var products []internal.Product

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		p := internal.Product{}
		byt := fileScanner.Bytes()

		err := json.Unmarshal(byt, &p)
		if err != nil {
			return nil, err
		}

		ok, err := destiny.SaveProduct(p)
		if !ok {
			return nil, err
		}

		products = append(products, p)
	}

	return products, nil
}
