package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

type ProductData struct {
	Name        string   `json:"name"`
	Price       int      `json:"price"`
	OldPrice    int      `json:"old_price"`
	Size        []Size   `json:"size"`
	Count       []Stock  `json:"stock"`
	Description string   `json:"description"`
	Color       []Color  `json:"color"`
	Article     string   `json:"article"`
	ProductURL  []string `json:"product_url"`
	Category    string   `json:"category"`
	ProductID   int      `json:"product_id"`
	Composition string   `json:"composition"`
	Care        string   `json:"care"`
	Medias      []string `json:"medias"`
}

type Size struct {
	Value string `json:"value"`
}

type Stock struct {
	Online  int `json:"online"`
	Offline int `json:"offline"`
}

type Color struct {
	Name string `json:"name"`
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

type Media struct {
	ID     int    `json:"id"`
	Slot   int    `json:"slot"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Type   int    `json:"type"`
	URL    string `json:"url"`
}

type Sku struct {
	Price    int   `json:"price"`
	OldPrice int   `json:"old_price"`
	Size     Size  `json:"size"`
	Stock    Stock `json:"stock"`
}

type Category struct {
	URL    string     `json:"url"`
	Items  []Category `json:"items"`
}



func main() {
	menuURLs := []string{
		"https://lime-shop.com/api/menu/left_women",
		"https://lime-shop.com/api/menu/left_men",
		"https://lime-shop.com/api/menu/left_kids",
	}

	var wg sync.WaitGroup
	productChan := make(chan ProductData, 100)

	for _, menuURL := range menuURLs {
		urls, err := getCategoryURLs(menuURL)
		if err != nil {
			fmt.Println("Error getting category URLs:", err)
			continue
		}

		for _, urlTemplate := range urls {
			for page := 1; page <= 7; page++ {
				wg.Add(1)
				go func(urlTemplate string, page int) {
					defer wg.Done()
					url := fmt.Sprintf(urlTemplate, page)
					parseURL(url, productChan)
				}(urlTemplate, page)
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
}

func getCategoryURLs(url string) ([]string, error) {
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
		// Заменяем "catalog" на "section"
		category.URL = strings.Replace(category.URL, "catalog", "section", 1)
		// Пропускаем URL, содержащие "#gift"
		if !strings.Contains(category.URL, "#gift") && !strings.Contains(category.URL, "/section/kids_all"){
			urls = append(urls, generateCatalogURL(category.URL))
		}
	}
	for _, item := range category.Items {
		urls = append(urls, extractURLs(item)...)
	}
	return urls
}

func generateCatalogURL(catalogCode string) string {
	return fmt.Sprintf("https://lime-shop.com/api%s?page_size=30", catalogCode)
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

	var data struct {
		Items []struct {
			Cells []struct {
				Entity Item `json:"entity"`
			} `json:"cells"`
		} `json:"items"`
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Printf("Error unmarshalling JSON for URL %s: %v\n", url, err)
		return
	}

	for _, item := range data.Items {
		for _, cell := range item.Cells {
			for _, model := range cell.Entity.Models {
				productURL := generateProductURL(cell.Entity.Code, model.Code)

				productData := ProductData{
					Name:        cell.Entity.Name,
					Description: cell.Entity.Description,
					Article:     cell.Entity.Article,
					ProductURL:  []string{productURL},
					Category:    model.Category,
					ProductID:   cell.Entity.ProductID,
					Composition: cell.Entity.Composition,
					Care:        cell.Entity.Care,
					Medias:      getMediaURLs(model.Medias),
					Color:       []Color{model.Color},
				}

				// Добавляем информацию о размерах и ценах
				for _, sku := range model.Skus {
					productData.Size = append(productData.Size, sku.Size)
					productData.Count = append(productData.Count, sku.Stock)
					productData.Price = sku.Price
					productData.OldPrice = sku.OldPrice
				}

				productChan <- productData
			}
		}
	}
}

func getMediaURLs(medias []Media) []string {
	var urls []string
	for _, media := range medias {
		urls = append(urls, media.URL)
	}
	return urls
}

func generateProductURL(productCode string, modelCode string) string {
	return fmt.Sprintf("https://lime-shop.com/api/v2/product/%s?id=%s&force=false&model=%s", productCode, productCode, modelCode)
}