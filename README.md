# Tech WB L0 - Система управления заказами

## Описание проекта

Микросервисная система для управления заказами. Система включает в себя backend API, Kafka для обработки сообщений(заказов), PostgreSQL для хранения данных, кэш и веб-интерфейс для просмотра заказов.

## Архитектура системы

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   Backend API   │    │  Kafka Producer │
│   (Port 8000)   │◄──►│   (Port 8080)   │◄──►│                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │                       │
                                ▼                       ▼
                       ┌─────────────────┐    ┌─────────────────┐
                       │   PostgreSQL    │    │     Kafka       │
                       │   + Cache       │    │   (Port 9094)   │
                       └─────────────────┘    └─────────────────┘
```

## Основные компоненты

### 1. Backend API (`/backend`)
- **Язык**: Go 1.23.4
- **Фреймворк**: Chi Router
- **Порт**: 8080
- **Функции**:
  - REST API для получения заказов
  - Интеграция с Kafka Consumer
  - Управление базой данных PostgreSQL
  - Кэширование в Redis
  - Автоматические миграции БД

### 2. Frontend (`/frontend`)
- **Язык**: HTML, JavaScript
- **Порт**: 8000
- **Функции**:
  - Веб-интерфейс для поиска заказов
  - Отображение детальной информации о заказах
  - Интеграция с Backend API

### 3. Kafka Producer (`/kafka-producer`)
- **Язык**: Go
- **Функции**:
  - Отправка заказов в Kafka топик
  - Чтение данных из JSON файла
  - Имитация потока заказов

### 4. Инфраструктура
- **PostgreSQL**: Основная база данных
- **Cache**: Кэш для быстрого доступа к данным
- **Kafka**: Очередь сообщений для асинхронной обработки
- **Zookeeper**: Координация Kafka кластера

## API Endpoints

### GET `/order/{orderUID}`
Получение информации о заказе по его уникальному идентификатору.

**Параметры:**
- `orderUID` (path) - уникальный идентификатор заказа

**Ответ:**
- `200 OK` - заказ найден и возвращен
- `400 Bad Request` - отсутствует orderUID
- `404 Not Found` - заказ не найден
- `500 Internal Server Error` - внутренняя ошибка сервера

**Пример ответа:**
```json
{
  "order_uid": "b563feb7b2b84b6test",
  "track_number": "WBILMTESTTRACK",
  "entry": "WBIL",
  "delivery": {
    "name": "Test Testov",
    "phone": "+9720000000",
    "zip": "2639809",
    "city": "Kiryat Mozkin",
    "address": "Ploshad Mira 15",
    "region": "Kraiot",
    "email": "test@gmail.com"
  },
  "payment": {
    "transaction": "b563feb7b2b84b6test",
    "request_id": "",
    "currency": "USD",
    "provider": "wbpay",
    "amount": 1817,
    "payment_dt": 1637901234,
    "bank": "alpha",
    "delivery_cost": 150,
    "goods_total": 1667,
    "custom_fee": 0
  },
  "items": [
    {
      "chrt_id": 9934930,
      "track_number": "WBILMTESTTRACK",
      "price": 453,
      "rid": "ab4219087a764ae0btest",
      "name": "Mascaras",
      "sale": 30,
      "size": "0",
      "total_price": 317,
      "nm_id": 2389212,
      "brand": "Vivienne Sabo",
      "status": 202
    }
  ],
  "locale": "en",
  "internal_signature": "",
  "customer_id": "test",
  "delivery_service": "meest",
  "shardkey": "9",
  "sm_id": 99,
  "date_created": "2021-11-26T06:22:19Z",
  "oof_shard": "1"
}
```

## Функциональные возможности

### 1. Создание заказов
- Автоматическое создание заказов через Kafka Consumer
- Валидация данных заказов
- Сохранение в PostgreSQL и кэширование в Redis

### 2. Поиск заказов
- Поиск заказов по уникальному идентификатору (orderUID)
- Быстрый доступ через Redis кэш
- Fallback на PostgreSQL при отсутствии в кэше

### 3. Восстановление кэша
- Автоматическое восстановление кэша из базы данных при запуске
- Синхронизация данных между Redis и PostgreSQL

### 4. Асинхронная обработка
- Получение заказов из Kafka топика
- Обработка сообщений в фоновом режиме
- Graceful shutdown при получении сигналов

### 5. Веб-интерфейс
- Поиск заказов по ID
- Детальное отображение информации о заказе
- Структурированный вывод данных о доставке, оплате и товарах
- Адаптивный дизайн

## Установка и запуск

### Предварительные требования
- Go 1.23.4+
- Docker и Docker Compose
- PostgreSQL

### 1. Клонирование репозитория
```bash
git clone <repository-url>
cd tech-wb-L0
```

### 2. Настройка переменных окружения
Создайте файл `.env` в папке `backend`:
```env
DATABASE_URL=postgres://username:password@localhost:5432/dbname?sslmode=disable
```

### 3. Запуск инфраструктуры
```bash
cd backend
docker-compose up -d
```

### 4. Запуск Backend API
```bash
cd backend
go mod tidy
go run main.go
```

### 5. Запуск Frontend
```bash
cd frontend
go run main.go
```

### 6. Запуск Kafka Producer (опционально)
```bash
cd kafka-producer
go mod tidy
go run main.go
```

## Конфигурация

### Порты
- **Frontend**: 8000
- **Backend API**: 8080
- **Kafka**: 9094
- **Zookeeper**: 2181

### Kafka настройки
- **Топик**: `orders`
- **Группа потребителей**: `orders-consumer-group`
- **Брокер**: `localhost:9094`

## Использование

### 1. Поиск заказа через веб-интерфейс
1. Откройте `http://localhost:8000` в браузере
2. Введите `order_uid` в поле поиска
3. Нажмите кнопку "Загрузить заказ"
4. Просмотрите детальную информацию о заказе

### 2. API вызовы
```bash
# Получение заказа по ID
curl http://localhost:8080/order/b563feb7b2b84b6test
```

### 3. Отправка заказов в Kafka
1. Поместите файл `orders.json` в папку `kafka-producer`
2. Запустите producer: `go run main.go`
3. Заказы автоматически поступят в систему через Kafka Consumer

## Структура проекта

```
tech-wb-L0/
├── backend/                 # Backend API сервер
│   ├── domain/             # Модели данных
│   ├── internal/           # Внутренняя логика
│   │   ├── api/           # HTTP handlers
│   │   ├── kafka/         # Kafka consumer
│   │   ├── repository/    # Слой доступа к данным
│   │   └── service/       # Бизнес-логика
│   ├── migrations/        # SQL миграции
│   ├── pkg/              # Общие пакеты
│   ├── docker-compose.yml # Docker конфигурация
│   └── main.go           # Точка входа
├── frontend/              # Веб-интерфейс
│   ├── index.html        # HTML страница
│   ├── script.js         # JavaScript логика
│   └── main.go           # HTTP сервер
├── kafka-producer/        # Kafka producer
│   ├── domain/           # Модели данных
│   ├── orders.json       # Тестовые данные
│   └── main.go           # Точка входа
├── go.mod                # Go модули
└── README.md             # Документация
```

Проект разработан для образовательных целей.
