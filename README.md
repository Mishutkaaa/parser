#Lime-shop
Сбор информации о товарах с lime-shop.com
## Расположение
## Запуск
- Локально 
    1. Клонировать репозиторий: `git clone https://github.com/Mishutkaaa/parser.git`
## Конфигурация
  Парсер не требует дополнительной конфигурации, так как использует предопределенные URL-адреса для сбора данных. Однако, если необходимо изменить источники данных, это можно сделать в коде:
  1. Измените список URL-адресов категорий в переменной menuURLs в функции main.
  2. Настройте параметры парсинга в функциях getCategoryURLs, getProductCodes и parseURL.
## Модели/Структуры данных

### Категория (Category)

```json
{
    "url": "string",
    "items": [
        {
           "url": "string",
            "items": []
        }
    ]
}
```
  ### Товар (Item)
  
  ```json
{
    "name": "string",
    "description": "string",
    "code": "string",
    "article": "string",
    "composition": "string",
    "care": "string",
    "product_id": 0,
    "models": [
        {
            "code": "string",
            "category": "string",
            "medias": [
                {
                    "url": "string"
                }
            ],
            "skus": [
                {
                    "id": 0,
                    "price": 0,
                    "old_price": 0,
                    "size": {
                        "value": "string"
                    },
                    "stock": {
                        "online": 0,
                        "offline": 0
                    },
                    "sku": "string"
                }
            ],
            "color": {
                "name": "string"
            }
        }
    ]
}
```
  ### Модель (Model)

  ```json
{
    "code": "string",
    "category": "string",
    "medias": [
        {
            "url": "string"
        }
    ],
    "skus": [
        {
            "id": 0,
            "price": 0,
            "old_price": 0,
            "size": {
                "value": "string"
            },
            "stock": {
                "online": 0,
                "offline": 0
            },
            "sku": "string"
        }
    ],
    "color": {
        "name": "string"
    }
}
```

  ### SKU (Sku)
  
  ```json
{
    "id": 0,
    "price": 0,
    "old_price": 0,
    "size": {
        "value": "string"
    },
    "stock": {
        "online": 0,
        "offline": 0
    },
    "sku": "string"
}
```
  ### ProductData
  
  ```json
{
    "name": "string",
    "price": 0,
    "composition": "string",
    "old_price": 0,
    "description": "string",
    "colors": [
        {
            "name": "string",
            "size_stock": [
                {
                    "id": 0,
                    "size": "string",
                    "unit": "string",
                    "stock": 0
                }
            ]
        }
    ],
    "article": "string",
    "product_url": "string",
    "category": "string",
    "product_id": 0,
    "care": "string",
    "medias": [
        "string"
    ]
}
```
## Основные функции

  1. getCategoryURLs(url string) ([]string, error)

      Получает URL-адреса категорий с указанного URL.
      Параметры:
        url (string): URL для получения категорий.
      Возвращает:
        Список URL-адресов категорий.
        Ошибку, если запрос не удался.
     
  2. extractURLs(category Category) []string

      Рекурсивно извлекает все URL-адреса из структуры категории.
      Параметры:
        category (Category): Структура категории.
      Возвращает:
        Список URL-адресов.

  3. generateCatalogURL(catalogCode string, page int) string

      Генерирует URL для страницы каталога с указанным кодом и номером страницы.
      Параметры:
        catalogCode (string): Код каталога.
        page (int): Номер страницы.
      Возвращает:
        Сгенерированный URL.

  4. generateProductURL(productCode string, modelCode string) string

     Генерирует URL для страницы товара с указанным кодом товара и кодом модели.
     Параметры:
       productCode (string): Код товара.
       modelCode (string): Код модели.
     Возвращает:
       Сгенерированный URL.
     
  5. getProductCodes(url string) ([]Item, error)

      Получает коды товаров с указанного URL.
      Параметры:
        url (string): URL для получения товаров.
      Возвращает:
        Список товаров.
        Ошибку, если запрос не удался.

  6. parseURL(url string, productChan chan<- ProductData)

      Парсит данные о товаре с указанного URL и отправляет результат в канал.
      Параметры:
        url (string): URL для парсинга.
        productChan (chan<- ProductData): Канал для отправки данных о товаре.
      Возвращает:
        Данные в канал productChan
     
  7.  getMediaURLs(medias []Media) []string
  
      Извлекает URL-адреса медиа (изображений) из списка медиа.
      Параметры:
        medias ([]Media): Список медиа.
      Возвращает:
      Список URL-адресов.

  8. generateLandingURL(productCode string, colorCode string) string

      Генерирует URL для целевой страницы товара с указанным кодом товара и кодом цвета.
      Параметры:
        productCode (string): Код товара.
        colorCode (string): Код цвета.
      Возвращает:
        Сгенерированный URL.

## Основные зависимости
Сторонние зависимости отсутствуют
