package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

type ProductData struct {
	Name        string        `json:"name"`
	Price       int           `json:"price"`
	OldPrice    int           `json:"old_price"`
	Size        []Size        `json:"size"`
	Count 		[]Stock `json:"stock"`
	Description string        `json:"description"`
	Color       []Color       `json:"color"`
	Article     string        `json:"article"`
	ProductURL  []string      `json:"product_url"`
	Category    string        `json:"category"`
	ProductID   int           `json:"product_id"`
	Composition string        `json:"composition"`
	Care        string        `json:"care"`
	Medias      []string      `json:"medias"`
}

type Size struct {
	Value string `json:"value"`

}
type Stock struct{
		Online int `json:"online"`
	Offline int `json:"offline"`
}
type Color struct {

	Name     string `json:"name"`
}

type Product struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Code string `json:"code"`
	Article     string   `json:"article"`
	Composition string   `json:"composition"`
	Care        string   `json:"care"`
	ProductID   int      `json:"id"` 
	Models      []Model  `json:"models"`
}

type Model struct {
	Code string			`json:"code"`
	Category string   `json:"category"`
	Medias   []Media  `json:"medias"`
	Skus     []Sku    `json:"skus"`
	Color    Color    `json:"color"`
}

type Media struct {
	URLs string `json:"url"`
}

type Sku struct {
	Price    int    `json:"price"`
	OldPrice int    `json:"old_price"`
	Size     Size   `json:"size"`
	Stock Stock `json:"stock"`
}

func main() {
	urls := []string{
"https://lime-shop.com/api/section/outerwear?page=%d&page_size=30",
		"https://lime-shop.com/api/section/knitwear?page=%d&page_size=30",
    "https://lime-shop.com/api/section/blazers?page=%d&page_size=30",
    "https://lime-shop.com/api/section/waistcoats?page=%d&page_size=30",
    "https://lime-shop.com/api/section/trousers?page=%d&page_size=30",
    "https://lime-shop.com/api/section/suits?page=%d&page_size=30",
    "https://lime-shop.com/api/section/dresses?page=%d&page_size=30",
    "https://lime-shop.com/api/section/women_shirts_all?page=%d&page_size=30",
    "https://lime-shop.com/api/section/skirts?page=%d&page_size=30",
    "https://lime-shop.com/api/section/co_ord_sets?page=%d&page_size=30",
    "https://lime-shop.com/api/section/down_jackets?page=%d&page_size=30",
    "https://lime-shop.com/api/section/women_sweaters_cardigans?page=%d&page_size=30",
    "https://lime-shop.com/api/section/t_shirts?page=%d&page_size=30",
    "https://lime-shop.com/api/section/women_jeans?page=%d&page_size=30",
    "https://lime-shop.com/api/section/sweatshirts?page=%d&page_size=30",
    "https://lime-shop.com/api/section/tops?page=%d&page_size=30",
    "https://lime-shop.com/api/section/shorty?page=%d&page_size=30",
    "https://lime-shop.com/api/section/sportswear?page=%d&page_size=30",
    "https://lime-shop.com/api/section/all_shoes?page=%d&page_size=30",
    "https://lime-shop.com/api/section/bags?page=%d&page_size=30",
    "https://lime-shop.com/api/section/accessories?page=%d&page_size=30",
    "https://lime-shop.com/api/section/jewellery?page=%d&page_size=30",
    "https://lime-shop.com/api/section/underwear?page=%d&page_size=30",
    "https://lime-shop.com/api/section/loungewear?page=%d&page_size=30",
    "https://lime-shop.com/api/section/women_winter_sets?page=%d&page_size=30",
    "https://lime-shop.com/api/section/women_wool?page=%d&page_size=30",
    "https://lime-shop.com/api/section/women_basic_wardrobe?page=%d&page_size=30",
    "https://lime-shop.com/api/section/last_sizes?page=%d&page_size=30",
    "https://lime-shop.com/api/section/limited_edition?page=%d&page_size=30",
    "https://lime-shop.com/api/section/kids_girls_view_all?page=%d&page_size=30",
    "https://lime-shop.com/api/section/kids_boys_view_all?page=%d&page_size=30",
    "https://lime-shop.com/api/section/kids_baby_girls_view_all?page=%d&page_size=30",
    "https://lime-shop.com/api/section/kids_baby_boys_view_all?page=%d&page_size=30",
	}

	var wg sync.WaitGroup
	productChan := make(chan ProductData, 100)


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

	go func() {
		wg.Wait()
		close(productChan)
	}()


	var products []ProductData
	for product := range productChan {
		products = append(products, product)
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
				Entity Product `json:"entity"`
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
				for _, sku := range model.Skus {
					productURL := generateProductURL(cell.Entity.Code, model.Code)


					productData := ProductData{
						Name:        cell.Entity.Name,
						Price:       sku.Price,
						OldPrice:    sku.OldPrice,
						Size:        []Size{sku.Size},
						Count: 		 []Stock{sku.Stock},
						Description: cell.Entity.Description,
						Color:       []Color{model.Color},
						Article:     cell.Entity.Article,
						ProductURL:  []string{productURL},
						Category:    model.Category,
						ProductID:   cell.Entity.ProductID, 
						Composition: cell.Entity.Composition,
						Care:        cell.Entity.Care,
						Medias:      getMediaURLs(model.Medias),
					}

					productChan <- productData
					
				}
			}
		}
	}
}

func getMediaURLs(medias []Media) []string {
	var urls []string
	for _, media := range medias {
		urls = append(urls, media.URLs)
	}
	return urls
}

func generateProductURL(productCode string, modelCode string) string {
	return fmt.Sprintf("https://lime-shop.com/api/v2/product/%s?id=%s&force=false&model=%s", productCode, productCode, modelCode)
}