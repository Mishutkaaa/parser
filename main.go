package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

type ProductData struct {
	Name        string                 `json:"name"`
	Price       int                    `json:"price"`
	OldPrice    int                    `json:"old_price"`
	Size        []map[string]interface{} `json:"size"`
	Description string                 `json:"description"`
	Color       []map[string]interface{} `json:"color"`
	Article     string                 `json:"article"`
	ProductURL  []string               `json:"product_url"`
	Category    string                 `json:"category"`
	ProductID   int                    `json:"product_id"`
	Composition string                 `json:"composition"`
	Care        string                 `json:"care"`
}

func main() {
	// Список ссылок для парсинга
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

	// Запуск горутин для каждой ссылки
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

	// Закрытие канала после завершения всех горутин
	go func() {
		wg.Wait()
		close(productChan)
	}()

	// Сбор результатов из канала
	var products []ProductData
	for product := range productChan {
		products = append(products, product)
	}

	// Сериализация результатов в JSON
	jsonData, err := json.MarshalIndent(products, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	// Сохранение результатов в файл
	err = ioutil.WriteFile("products.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing JSON to file:", err)
		return
	}

	// Вывод результатов
	fmt.Println("Data saved to products.json")
}

func parseURL(url string, productChan chan<- ProductData) {
	// Выполнение HTTP-запроса
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error making HTTP request for URL %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close()

	// Чтение тела ответа
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body for URL %s: %v\n", url, err)
		return
	}

	// Десериализация JSON
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Printf("Error unmarshalling JSON for URL %s: %v\n", url, err)
		return
	}

	// Проверка наличия ключа "items"
	items, ok := data["items"].([]interface{})
	if !ok {
		fmt.Printf("Key 'items' not found or not a slice for URL %s\n", url)
		return
	}

	// Извлечение данных
	for _, item := range items {
		cells, ok := item.(map[string]interface{})["cells"].([]interface{})
		if !ok {
			fmt.Printf("Key 'cells' not found or not a slice for item in URL %s\n", url)
			continue
		}

		for _, cell := range cells {
			product, ok := cell.(map[string]interface{})["entity"].(map[string]interface{})
			if !ok {
				fmt.Printf("Key 'entity' not found or not a map for cell in URL %s\n", url)
				continue
			}

			models, ok := product["models"].([]interface{})
			if !ok || len(models) == 0 {
				fmt.Printf("Key 'models' not found or empty for product in URL %s\n", url)
				continue
			}

			for _, model := range models {
				modelMap, ok := model.(map[string]interface{})
				if !ok {
					fmt.Printf("Model is not a map for product in URL %s\n", url)
					continue
				}

				skus, ok := modelMap["skus"].([]interface{})
				if !ok || len(skus) == 0 {
					fmt.Printf("Key 'skus' not found or empty for model in URL %s\n", url)
					continue
				}

				sku := skus[0].(map[string]interface{})

				// Проверка наличия ключа "medias"
				medias, ok := modelMap["medias"].([]interface{})
				if !ok {
					fmt.Printf("Key 'medias' not found or not a slice for model in URL %s\n", url)
					continue
				}

				var photoURLs []string
				for _, media := range medias {
					mediaMap, ok := media.(map[string]interface{})
					if !ok {
						fmt.Printf("Media is not a map for model in URL %s\n", url)
						continue
					}

					mediaURL, ok := mediaMap["url"].(string)
					if !ok {
						fmt.Printf("Key 'url' not found or not a string in 'media' for model in URL %s\n", url)
						continue
					}

					photoURLs = append(photoURLs, mediaURL)
				}

				// Проверка наличия ключей "name", "price", "old_price", "description", "article", "category", "id", "composition", "care"
				name, ok := product["name"].(string)
				if !ok {
					fmt.Printf("Key 'name' not found or not a string for product in URL %s\n", url)
					continue
				}

				price, ok := sku["price"].(float64)
				if !ok {
					fmt.Printf("Key 'price' not found or not a float64 for SKU in URL %s\n", url)
					continue
				}

				// Проверка наличия ключа "old_price"
				var oldPrice float64
				if oldPriceVal, ok := sku["old_price"].(float64); ok {
					oldPrice = oldPriceVal
				} else {
					oldPrice = 0 // Устанавливаем значение по умолчанию, если ключ отсутствует
				}

				description, ok := product["description"].(string)
				if !ok {
					fmt.Printf("Key 'description' not found or not a string for product in URL %s\n", url)
					continue
				}

				article, ok := product["article"].(string)
				if !ok {
					fmt.Printf("Key 'article' not found or not a string for product in URL %s\n", url)
					continue
				}

				category, ok := modelMap["category"].(string)
				if !ok {
					fmt.Printf("Key 'category' not found or not a string for model in URL %s\n", url)
					continue
				}

				productID, ok := product["id"].(float64)
				if !ok {
					fmt.Printf("Key 'id' not found or not a float64 for product in URL %s\n", url)
					continue
				}

				composition, ok := product["composition"].(string)
				if !ok {
					fmt.Printf("Key 'composition' not found or not a string for product in URL %s\n", url)
					continue
				}

				care, ok := product["care"].(string)
				if !ok {
					fmt.Printf("Key 'care' not found or not a string for product in URL %s\n", url)
					continue
				}

				// Извлечение полей "product.code" и "models.code"
				productCode, ok := product["code"].(string)
				
				if !ok {
					fmt.Printf("Key 'code' not found or not a string for product in URL %s\n", url)
					continue
				}

				modelCode, ok := modelMap["code"].(string)
				
				if !ok {
					fmt.Printf("Key 'code' not found or not a string for model in URL %s\n", url)
					continue
				}



				// Генерация ссылки на продукт
				productURL := generateProductURL(productCode, modelCode)

				// Извлечение данных
				productData := ProductData{
					Name:        name,
					Price:       int(price),
					OldPrice:    int(oldPrice),
					Size:        getSizes(skus),
					Description: description,
					Color:       getColors(models),
					Article:     article,
					ProductURL:  []string{productURL}, // Добавляем сгенерированную ссылку
					Category:    category,
					ProductID:   int(productID),
					Composition: composition,
					Care:        care,
				}

				// Отправка данных в канал
				productChan <- productData
			}
		}
	}
}

func getSizes(skus []interface{}) []map[string]interface{} {
	var sizes []map[string]interface{}
	for _, sku := range skus {
		size := sku.(map[string]interface{})["size"].(map[string]interface{})
		sizes = append(sizes, map[string]interface{}{
			"value": size["value"],
			"stock": sku.(map[string]interface{})["stock"],
		})
	}
	return sizes
}

func getColors(models []interface{}) []map[string]interface{} {
	var colors []map[string]interface{}
	for _, model := range models {
		color := model.(map[string]interface{})["color"].(map[string]interface{})
		colors = append(colors, map[string]interface{}{
			"unique_id": model.(map[string]interface{})["id"],
			"name":      color["name"],
		})
	}
	return colors
}

func generateProductURL(productCode string, modelCode string) string {
	return fmt.Sprintf("https://lime-shop.com/api/v2/product/%s?id=%s&force=false&model=false&color=%s", productCode, productCode, modelCode)
}