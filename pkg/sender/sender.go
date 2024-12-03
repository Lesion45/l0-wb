package sender

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"math/rand"
	"time"
)

func Run() {
	var orders []Order
	for i := 0; i < 10; i++ {
		order := generateOrder(i)
		fmt.Println(order)
		orders = append(orders, order)
	}
	topic := "orders"
	brokerAddress := "kafka:9092"

	writer := kafka.Writer{
		Addr:  kafka.TCP(brokerAddress),
		Topic: topic,
	}

	for _, order := range orders {
		message, err := json.Marshal(order)
		if err != nil {
			log.Printf("Failed to marshal order: %v", err)
			continue
		}
		err = writer.WriteMessages(context.Background(), kafka.Message{
			Value: message,
		})

		if err != nil {
			log.Fatalf("Failed to write messages: %v", err)
		}

		log.Println("Message written successfully!")

	}
}

func generateOrder(index int) Order {
	rand.Seed(time.Now().UnixNano())
	order := Order{
		OrderUID:        randomString(10) + "test",
		TrackNumber:     "WBILMTESTTRACK",
		Entry:           "WBIL",
		Locale:          "en",
		CustomerID:      "customer " + randomString(5),
		DeliveryService: "meest",
		Shardkey:        "123",
		SmID:            91,
		DateCreated:     time.Now().Format(time.RFC3339),
		OofShard:        "A1234",
		Delivery: struct {
			Name    string `json:"name"`
			Phone   string `json:"phone"`
			Zip     string `json:"zip"`
			City    string `json:"city"`
			Address string `json:"address"`
			Region  string `json:"region"`
			Email   string `json:"email"`
		}{
			Name:    "Test User " + randomString(4),
			Phone:   "+1234567890",
			Zip:     "00000",
			City:    "Sample City",
			Address: "Sample Address " + randomString(4),
			Region:  "Sample Region",
			Email:   "test" + randomString(5) + "@example.com",
		},
		Payment: struct {
			Transaction  string `json:"transaction"`
			RequestID    string `json:"request_id"`
			Currency     string `json:"currency"`
			Provider     string `json:"provider"`
			Amount       int    `json:"amount"`
			PaymentDt    int64  `json:"payment_dt"`
			Bank         string `json:"bank"`
			DeliveryCost int    `json:"delivery_cost"`
			GoodsTotal   int    `json:"goods_total"`
			CustomFee    int    `json:"custom_fee"`
		}{
			Transaction:  randomString(15),
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       rand.Intn(1000) + 100,
			PaymentDt:    time.Now().Unix(),
			Bank:         "TestBank",
			DeliveryCost: rand.Intn(200) + 50,
			GoodsTotal:   rand.Intn(500) + 50,
			CustomFee:    0,
		},
		Items: []struct {
			ChrtID      int    `json:"chrt_id"`
			TrackNumber string `json:"track_number"`
			Price       int    `json:"price"`
			Rid         string `json:"rid"`
			Name        string `json:"name"`
			Sale        int    `json:"sale"`
			Size        string `json:"size"`
			TotalPrice  int    `json:"total_price"`
			NmID        int    `json:"nm_id"`
			Brand       string `json:"brand"`
			Status      int    `json:"status"`
		}{
			{
				ChrtID:      rand.Intn(100000),
				TrackNumber: "TRACK" + randomString(3),
				Price:       rand.Intn(100),
				Rid:         randomString(10),
				Name:        "Test Item",
				Sale:        10,
				Size:        "L",
				TotalPrice:  rand.Intn(100),
				NmID:        rand.Intn(10000),
				Brand:       "TestBrand",
				Status:      202,
			},
		},
	}
	return order
}

func randomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

type Order struct {
	OrderUID    string `json:"order_uid"`
	TrackNumber string `json:"track_number"`
	Entry       string `json:"entry"`
	Delivery    struct {
		Name    string `json:"name"`
		Phone   string `json:"phone"`
		Zip     string `json:"zip"`
		City    string `json:"city"`
		Address string `json:"address"`
		Region  string `json:"region"`
		Email   string `json:"email"`
	} `json:"delivery"`
	Payment struct {
		Transaction  string `json:"transaction"`
		RequestID    string `json:"request_id"`
		Currency     string `json:"currency"`
		Provider     string `json:"provider"`
		Amount       int    `json:"amount"`
		PaymentDt    int64  `json:"payment_dt"`
		Bank         string `json:"bank"`
		DeliveryCost int    `json:"delivery_cost"`
		GoodsTotal   int    `json:"goods_total"`
		CustomFee    int    `json:"custom_fee"`
	} `json:"payment"`
	Items []struct {
		ChrtID      int    `json:"chrt_id"`
		TrackNumber string `json:"track_number"`
		Price       int    `json:"price"`
		Rid         string `json:"rid"`
		Name        string `json:"name"`
		Sale        int    `json:"sale"`
		Size        string `json:"size"`
		TotalPrice  int    `json:"total_price"`
		NmID        int    `json:"nm_id"`
		Brand       string `json:"brand"`
		Status      int    `json:"status"`
	} `json:"items"`
	Locale            string `json:"locale"`
	InternalSignature string `json:"internal_signature"`
	CustomerID        string `json:"customer_id"`
	DeliveryService   string `json:"delivery_service"`
	Shardkey          string `json:"shardkey"`
	SmID              int    `json:"sm_id"`
	DateCreated       string `json:"date_created"`
	OofShard          string `json:"oof_shard"`
}
