package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	const totalRows = 1_000_000
	categories := []string{"books", "electronics", "toys", "clothing", "sports"}
	startDate := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)

	file, err := os.Create("../../internal/data/large_dataset.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"id", "user_id", "product_id", "category", "price", "quantity", "date"})

	for i := 1; i <= totalRows; i++ {
		id := strconv.Itoa(i)
		userID := strconv.Itoa(rand.Intn(10_000) + 1)
		productID := strconv.Itoa(rand.Intn(1_000) + 1)
		category := categories[rand.Intn(len(categories))]
		price := fmt.Sprintf("%.2f", rand.Float64()*200+1)
		quantity := strconv.Itoa(rand.Intn(5) + 1)
		date := startDate.AddDate(0, 0, rand.Intn(365*2)).Format("2006-01-02")

		record := []string{id, userID, productID, category, price, quantity, date}
		writer.Write(record)
	}
}
