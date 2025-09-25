# API Документация

## Обзор

API для работы с ботом Cos AI. Все запросы должны содержать заголовок `X-Telegram-ID` с Telegram ID пользователя.

## Базовый URL

```
http://localhost:8080
```

## Аутентификация

Все запросы должны содержать заголовок `X-Telegram-ID` с Telegram ID пользователя:

```
X-Telegram-ID: 374892056
```

## Endpoints

### 1. Получение продуктов пользователя

**GET** `/api/user/products`

Получает все продукты, связанные с пользователем.

#### Заголовки
```
X-Telegram-ID: 374892056
Content-Type: application/json
```

#### Пример запроса
```bash
curl -X GET "http://localhost:8080/api/user/products" \
  -H "X-Telegram-ID: 374892056" \
  -H "Content-Type: application/json"
```

#### Пример ответа
```json
{
  "success": true,
  "data": [
    {
      "brand": "La Roche-Posay",
      "product_title": "Effaclar Duo+",
      "image": "https://example.com/image.jpg",
      "ingredients": "Aqua, Glycerin, Niacinamide...",
      "description": "Корректирующий гель для проблемной кожи"
    }
  ]
}
```

### 2. Получение состояния пользователя

**GET** `/api/user/state`

Получает текущее состояние пользователя (данные формы).

#### Заголовки
```
X-Telegram-ID: 374892056
Content-Type: application/json
```

#### Пример запроса
```bash
curl -X GET "http://localhost:8080/api/user/state" \
  -H "X-Telegram-ID: 374892056" \
  -H "Content-Type: application/json"
```

#### Пример ответа
```json
{
  "success": true,
  "data": {
    "user_id": 374892056,
    "step": 5,
    "skin_type": "skin_oily",
    "age": "age_25_34",
    "gender": "gender_female",
    "pregnancy": "none_of_above",
    "concerns": "Акне и черные точки",
    "goal": "goal_clear",
    "budget": "budget_medium",
    "current_routine": "Умываюсь гелем",
    "allergies": "allergies_none",
    "preferences": "Без отдушек"
  }
}
```

### 3. Получение всех продуктов

**GET** `/api/products`

Получает все продукты из базы данных.

#### Заголовки
```
X-Telegram-ID: 374892056
Content-Type: application/json
```

#### Пример запроса
```bash
curl -X GET "http://localhost:8080/api/products" \
  -H "X-Telegram-ID: 374892056" \
  -H "Content-Type: application/json"
```

### 4. Health Check

**GET** `/health`

Проверка работоспособности сервиса.

#### Пример запроса
```bash
curl -X GET "http://localhost:8080/health"
```

#### Пример ответа
```json
{
  "success": true,
  "message": "Service is healthy"
}
```

## Коды ошибок

### 400 Bad Request
- Отсутствует заголовок `X-Telegram-ID`
- Неверный формат Telegram ID

### 404 Not Found
- Пользователь не найден
- Ресурс не существует

### 405 Method Not Allowed
- Неподдерживаемый HTTP метод

### 500 Internal Server Error
- Внутренняя ошибка сервера

## Примеры использования в коде

### JavaScript (Fetch API)
```javascript
const telegramID = '374892056';

fetch('http://localhost:8080/api/user/products', {
  method: 'GET',
  headers: {
    'X-Telegram-ID': telegramID,
    'Content-Type': 'application/json'
  }
})
.then(response => response.json())
.then(data => {
  if (data.success) {
    console.log('Продукты пользователя:', data.data);
  } else {
    console.error('Ошибка:', data.error);
  }
});
```

### Python (requests)
```python
import requests

telegram_id = '374892056'
headers = {
    'X-Telegram-ID': telegram_id,
    'Content-Type': 'application/json'
}

response = requests.get('http://localhost:8080/api/user/products', headers=headers)
data = response.json()

if data['success']:
    print('Продукты пользователя:', data['data'])
else:
    print('Ошибка:', data['error'])
```

### Go
```go
package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
)

func main() {
    client := &http.Client{}
    
    req, err := http.NewRequest("GET", "http://localhost:8080/api/user/products", nil)
    if err != nil {
        panic(err)
    }
    
    req.Header.Set("X-Telegram-ID", "374892056")
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }
    
    fmt.Println(string(body))
}
```

## Примечания

1. **Telegram ID** должен быть числом (например: 374892056)
2. Все ответы возвращаются в формате JSON
3. При отсутствии данных возвращается пустой массив `[]`
4. Сервер автоматически логирует все запросы с Telegram ID
5. Все запросы должны содержать валидный Telegram ID в заголовке
