package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Category struct {
	URL   string     `json:"url"`
	Items []Category `json:"items"`
}

type Item struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Code        string  `json:"code"`
	Article     string  `json:"article"`
	Composition string  `json:"composition"`
	Care        string  `json:"care"`
	ProductID   int     `json:"id"`
	Models      []Model `json:"models"`
}

type Model struct {
	Code     string  `json:"code"`
	Category string  `json:"category"`
	Medias   []Media `json:"medias"`
	Skus     []Sku   `json:"skus"`
	Color    Color   `json:"color"`
}

type Color struct {
	Name string `json:"name"`
}

type Media struct {
	URL string `json:"url"`
}

type Sku struct {
	Id int `json:"id"`
	Price    int   `json:"price"`
	OldPrice int   `json:"old_price"`
	Size     Size  `json:"size"`
	Stock    Stock `json:"stock"`
	Sku      string `json:"sku"`
}

type Size struct {
	Value string `json:"value"`
}

type Stock struct {
	Online  int `json:"online"`
	Offline int `json:"offline"`
}

type ProductData struct {
	Name        string            `json:"name"`
	Price       int               `json:"price"`
	Composition string            `json:"composition"`
	OldPrice    *int              `json:"old_price"`
	Description string            `json:"description"`
	Colors      []ColorData       `json:"colors"`
	Article     string            `json:"article"`
	ProductURL  string            `json:"product_url"`
	Category    string            `json:"category"`
	ProductID   int               `json:"product_id"`
	Care        string            `json:"care"`
	Medias      []string          `json:"medias"`
}

type ColorData struct {
	Name      string      `json:"name"`
	SizeStock []SizeStock `json:"size_stock"`
}

type SizeStock struct {
	Id    int    `json:"id"`
	Size  string `json:"size"`
	Unit  string `json:"unit"`
	Stock int    `json:"stock"`
}

func main() {
	start := time.Now()

	menuURLs := []string{
		"https://lime-shop.com/api/menu/left_kids",
		"https://lime-shop.com/api/menu/left_women",
		"https://lime-shop.com/api/menu/left_men",

	}

	var wg sync.WaitGroup
	productChan := make(chan ProductData, 100)

	for _, menuURL := range menuURLs {
		urls, err := getCategoryURLs(menuURL)
		if err != nil {
			fmt.Println("Error getting category URLs:", err)
			continue
		}

		for _, url := range urls {
			url = strings.Replace(url, "catalog", "section", 1)
			for page := 1; page <= 7; page++ {
				newURL := generateCatalogURL(url, page)
				codes, err := getProductCodes(newURL)
				if err != nil {
					fmt.Println("Error getting product codes:", err)
					continue
				}
				for _, code := range codes {
					for _, model := range code.Models {
						productURL := generateProductURL(code.Code, model.Code)
						if strings.Contains(productURL, "#gift") {
							continue
						}

						wg.Add(1)
						go func(url string) {
							defer wg.Done()
							parseURL(url, productChan)
						}(productURL)
					}
				}
			}
		}
	}

	go func() {
		wg.Wait()
		close(productChan)
	}()

	var products []ProductData
	productMap := make(map[int]bool)

	for product := range productChan {
		if !productMap[product.ProductID] {
			productMap[product.ProductID] = true
			products = append(products, product)
		}
	}

	jsonData, err := json.MarshalIndent(products, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	err = ioutil.WriteFile("products.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing JSON to file:", err)
		return
	}

	fmt.Println("Data saved to products.json")

	elapsed := time.Since(start)
	fmt.Printf("Program execution time: %s\n", elapsed)
}

func getCategoryURLs(url string) ([]string, error) {
	fmt.Println("Fetching category URLs from", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data struct {
		Items []Category `json:"items"`
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	var urls []string
	for _, item := range data.Items {
		urls = append(urls, extractURLs(item)...)
	}

	return urls, nil
}

func extractURLs(category Category) []string {
	var urls []string
	if category.URL != "" {
		urls = append(urls, category.URL)
	}
	for _, item := range category.Items {
		urls = append(urls, extractURLs(item)...)
	}
	return urls
}

func generateCatalogURL(catalogCode string, page int) string {
	return fmt.Sprintf("https://lime-shop.com/api%s?page=%d&page_size=30", catalogCode, page)
}

func generateProductURL(productCode string, modelCode string) string {
	return fmt.Sprintf("https://lime-shop.com/api/v2/product/%s?id=%s&force=false&model=%s", productCode, productCode, modelCode)
}

func getProductCodes(url string) ([]Item, error) {
	fmt.Println("Fetching product codes from", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data struct {
		Items []struct {
			Cells []struct {
				Entity Item `json:"entity"`
			} `json:"cells"`
		} `json:"items"`
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	var items []Item
	for _, item := range data.Items {
		for _, cell := range item.Cells {
			items = append(items, cell.Entity)
		}
	}

	return items, nil
}

func parseURL(url string, productChan chan<- ProductData) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error making HTTP request for URL %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body for URL %s: %v\n", url, err)
		return
	}

	var item Item
	err = json.Unmarshal(body, &item)
	if err != nil {
		fmt.Printf("Error unmarshalling JSON for URL %s: %v\n", url, err)
		return
	}

	productData := ProductData{
		Name:        item.Name,
		Description: item.Description,
		Article:     item.Article,
		Composition: item.Composition,
		Care:        item.Care,
		ProductID:   item.ProductID,
	}
	for _, model := range item.Models {
		productData.Category = model.Category
		productData.Medias = getMediaURLs(model.Medias)
		colorData := ColorData{
			Name: model.Color.Name,
		}
		for _, sku := range model.Skus {
			sizeStock := SizeStock{
				Id:    sku.Id,
				Size:  sku.Size.Value,
				Unit:  "шт",
				Stock: sku.Stock.Online + sku.Stock.Offline,
			}
			colorData.SizeStock = append(colorData.SizeStock, sizeStock)
			productData.Price = sku.Price
			if sku.OldPrice > 0 {
				productData.OldPrice = &sku.OldPrice
			}
		}
		productData.Colors = append(productData.Colors, colorData)
		productData.ProductURL = generateLandingURL(item.Code, model.Code)
	}

	productChan <- productData
}

func getMediaURLs(medias []Media) []string {
	var urls []string
	for _, media := range medias {
		urls = append(urls, media.URL)
	}
	return urls
}

func generateLandingURL(productCode string, colorCode string) string {
	return fmt.Sprintf("https://lime-shop.com/ru_ru/product/%s-%s", productCode, colorCode)
}