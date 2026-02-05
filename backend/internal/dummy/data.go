package dummy

import (
	"fmt"
	"sync"
	"time"
)

// Order source constants
const (
	SourceCashier    = "cashier"
	SourceTableQR    = "table_qr"
	SourceGrabFood   = "grabfood"
	SourceGoFood     = "gofood"
	SourceShopeeFood = "shopee_food"
)

// In-memory dummy data store
var (
	mu sync.RWMutex

	// Categories
	Categories = []Category{
		{ID: "cat-1", Name: "Coffee", Description: "Hot and cold coffee drinks", SortOrder: 1},
		{ID: "cat-2", Name: "Non-Coffee", Description: "Tea, chocolate, and more", SortOrder: 2},
		{ID: "cat-3", Name: "Food", Description: "Snacks and meals", SortOrder: 3},
		{ID: "cat-4", Name: "Dessert", Description: "Sweet treats", SortOrder: 4},
	}

	// Products with variants and modifiers
	Products = []Product{
		{
			ID: "prod-1", CategoryID: "cat-1", Name: "Espresso", Description: "Strong Italian coffee",
			BasePrice: 18000, ImageURL: "", IsAvailable: true,
			Variants: []Variant{
				{ID: "var-1", Name: "Single Shot", PriceAdjustment: 0},
				{ID: "var-2", Name: "Double Shot", PriceAdjustment: 8000},
			},
			Modifiers: []Modifier{
				{ID: "mod-1", Name: "Extra Shot", Price: 8000, MaxQty: 3},
			},
		},
		{
			ID: "prod-2", CategoryID: "cat-1", Name: "Americano", Description: "Espresso with hot water",
			BasePrice: 22000, ImageURL: "", IsAvailable: true,
			Variants: []Variant{
				{ID: "var-3", Name: "Hot", PriceAdjustment: 0},
				{ID: "var-4", Name: "Iced", PriceAdjustment: 3000},
			},
			Modifiers: []Modifier{
				{ID: "mod-1", Name: "Extra Shot", Price: 8000, MaxQty: 3},
				{ID: "mod-2", Name: "Vanilla Syrup", Price: 5000, MaxQty: 2},
			},
		},
		{
			ID: "prod-3", CategoryID: "cat-1", Name: "Cappuccino", Description: "Espresso with steamed milk foam",
			BasePrice: 28000, ImageURL: "", IsAvailable: true,
			Variants: []Variant{
				{ID: "var-5", Name: "Regular", PriceAdjustment: 0},
				{ID: "var-6", Name: "Large", PriceAdjustment: 8000},
			},
			Modifiers: []Modifier{
				{ID: "mod-1", Name: "Extra Shot", Price: 8000, MaxQty: 3},
				{ID: "mod-3", Name: "Oat Milk", Price: 8000, MaxQty: 1},
			},
		},
		{
			ID: "prod-4", CategoryID: "cat-1", Name: "Latte", Description: "Smooth espresso with milk",
			BasePrice: 28000, ImageURL: "", IsAvailable: true,
			Variants: []Variant{
				{ID: "var-7", Name: "Hot", PriceAdjustment: 0},
				{ID: "var-8", Name: "Iced", PriceAdjustment: 3000},
			},
			Modifiers: []Modifier{
				{ID: "mod-2", Name: "Vanilla Syrup", Price: 5000, MaxQty: 2},
				{ID: "mod-4", Name: "Caramel Syrup", Price: 5000, MaxQty: 2},
			},
		},
		{
			ID: "prod-5", CategoryID: "cat-2", Name: "Matcha Latte", Description: "Japanese green tea latte",
			BasePrice: 32000, ImageURL: "", IsAvailable: true,
			Variants: []Variant{
				{ID: "var-9", Name: "Hot", PriceAdjustment: 0},
				{ID: "var-10", Name: "Iced", PriceAdjustment: 3000},
			},
			Modifiers: []Modifier{},
		},
		{
			ID: "prod-6", CategoryID: "cat-2", Name: "Chocolate", Description: "Rich hot chocolate",
			BasePrice: 25000, ImageURL: "", IsAvailable: true,
			Variants: []Variant{
				{ID: "var-11", Name: "Hot", PriceAdjustment: 0},
				{ID: "var-12", Name: "Iced", PriceAdjustment: 3000},
			},
			Modifiers: []Modifier{
				{ID: "mod-5", Name: "Whipped Cream", Price: 5000, MaxQty: 1},
			},
		},
		{
			ID: "prod-7", CategoryID: "cat-3", Name: "Croissant", Description: "Buttery French pastry",
			BasePrice: 25000, ImageURL: "", IsAvailable: true,
			Variants:  []Variant{},
			Modifiers: []Modifier{},
		},
		{
			ID: "prod-8", CategoryID: "cat-3", Name: "Sandwich", Description: "Grilled cheese sandwich",
			BasePrice: 35000, ImageURL: "", IsAvailable: true,
			Variants: []Variant{
				{ID: "var-13", Name: "Cheese", PriceAdjustment: 0},
				{ID: "var-14", Name: "Ham & Cheese", PriceAdjustment: 10000},
			},
			Modifiers: []Modifier{},
		},
		{
			ID: "prod-9", CategoryID: "cat-4", Name: "Cheesecake", Description: "New York style cheesecake",
			BasePrice: 38000, ImageURL: "", IsAvailable: true,
			Variants:  []Variant{},
			Modifiers: []Modifier{},
		},
		{
			ID: "prod-10", CategoryID: "cat-4", Name: "Brownies", Description: "Chocolate fudge brownies",
			BasePrice: 28000, ImageURL: "", IsAvailable: true,
			Variants:  []Variant{},
			Modifiers: []Modifier{},
		},
	}

	// Tables
	Tables = []Table{
		{ID: "table-1", Number: 1, Capacity: 2, Status: "available", QRCode: "QR001"},
		{ID: "table-2", Number: 2, Capacity: 4, Status: "available", QRCode: "QR002"},
		{ID: "table-3", Number: 3, Capacity: 4, Status: "occupied", QRCode: "QR003"},
		{ID: "table-4", Number: 4, Capacity: 6, Status: "available", QRCode: "QR004"},
		{ID: "table-5", Number: 5, Capacity: 2, Status: "reserved", QRCode: "QR005"},
	}

	// Users
	Users = []User{
		{ID: "11111111-1111-1111-1111-111111111111", Email: "admin@kaori.pos", Name: "Admin", Role: "super_admin", PIN: "1234"},
		{ID: "22222222-2222-2222-2222-222222222222", Email: "store@kaori.pos", Name: "Store Manager", Role: "store_admin", PIN: "5678"},
		{ID: "33333333-3333-3333-3333-333333333333", Email: "cashier@kaori.pos", Name: "John Cashier", Role: "cashier", PIN: "1111"},
		{ID: "44444444-4444-4444-4444-444444444444", Email: "kitchen@kaori.pos", Name: "Chef Mike", Role: "kitchen", PIN: "2222"},
	}

	// Orders (in-memory)
	Orders   = []Order{}
	orderSeq = 1000
)

// Types
type Category struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

type Variant struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	PriceAdjustment int    `json:"price_adjustment"`
}

type Modifier struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Price  int    `json:"price"`
	MaxQty int    `json:"max_qty"`
}

type Product struct {
	ID          string     `json:"id"`
	CategoryID  string     `json:"category_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	BasePrice   int        `json:"base_price"`
	ImageURL    string     `json:"image_url"`
	IsAvailable bool       `json:"is_available"`
	Variants    []Variant  `json:"variants"`
	Modifiers   []Modifier `json:"modifiers"`
}

type Table struct {
	ID       string `json:"id"`
	Number   int    `json:"number"`
	Capacity int    `json:"capacity"`
	Status   string `json:"status"`
	QRCode   string `json:"qr_code"`
}

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
	PIN   string `json:"-"`
}

type OrderItem struct {
	ID          string   `json:"id"`
	ProductID   string   `json:"product_id"`
	ProductName string   `json:"product_name"`
	VariantID   string   `json:"variant_id,omitempty"`
	VariantName string   `json:"variant_name,omitempty"`
	Modifiers   []string `json:"modifiers,omitempty"`
	Quantity    int      `json:"quantity"`
	UnitPrice   int      `json:"unit_price"`
	Subtotal    int      `json:"subtotal"`
	Notes       string   `json:"notes,omitempty"`
}

type Order struct {
	ID              string      `json:"id"`
	OrderNumber     string      `json:"order_number"`
	ExternalOrderID string      `json:"external_order_id,omitempty"` // GrabFood/GoFood/Shopee order ID
	TableID         string      `json:"table_id,omitempty"`
	TableNumber     int         `json:"table_number,omitempty"`
	OrderType       string      `json:"order_type"`   // dine_in, takeaway, delivery
	OrderSource     string      `json:"order_source"` // cashier, table_qr, grabfood, gofood, shopee_food
	Status          string      `json:"status"`       // pending, confirmed, cooking, ready, completed, cancelled
	PaymentStatus   string      `json:"payment_status"`
	Items           []OrderItem `json:"items"`
	Subtotal        int         `json:"subtotal"`
	Tax             int         `json:"tax"`
	Total           int         `json:"total"`
	Notes           string      `json:"notes,omitempty"`
	CustomerName    string      `json:"customer_name,omitempty"`
	CustomerPhone   string      `json:"customer_phone,omitempty"`
	DeliveryAddress string      `json:"delivery_address,omitempty"`
	DriverName      string      `json:"driver_name,omitempty"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
	CashierID       string      `json:"cashier_id,omitempty"`
	CashierName     string      `json:"cashier_name,omitempty"`
}

// Helper functions
func GetNextOrderNumber() string {
	mu.Lock()
	defer mu.Unlock()
	orderSeq++
	return fmt.Sprintf("ORD-%04d", orderSeq)
}

func GetNextDeliveryOrderNumber(source string) string {
	mu.Lock()
	defer mu.Unlock()
	orderSeq++
	prefix := "DEL"
	switch source {
	case SourceGrabFood:
		prefix = "GRAB"
	case SourceGoFood:
		prefix = "GOFOOD"
	case SourceShopeeFood:
		prefix = "SHOPEE"
	}
	return fmt.Sprintf("%s-%04d", prefix, orderSeq)
}

func AddOrder(order Order) {
	mu.Lock()
	defer mu.Unlock()
	Orders = append(Orders, order)
}

func UpdateOrderStatus(orderID, status string) bool {
	mu.Lock()
	defer mu.Unlock()
	for i := range Orders {
		if Orders[i].ID == orderID {
			Orders[i].Status = status
			Orders[i].UpdatedAt = time.Now()
			return true
		}
	}
	return false
}

func GetOrdersByStatus(statuses ...string) []Order {
	mu.RLock()
	defer mu.RUnlock()
	var result []Order
	for _, order := range Orders {
		for _, s := range statuses {
			if order.Status == s {
				result = append(result, order)
				break
			}
		}
	}
	return result
}

func GetOrdersBySource(source string) []Order {
	mu.RLock()
	defer mu.RUnlock()
	var result []Order
	for _, order := range Orders {
		if order.OrderSource == source {
			result = append(result, order)
		}
	}
	return result
}

func GetAllOrders() []Order {
	mu.RLock()
	defer mu.RUnlock()
	return append([]Order{}, Orders...)
}

func GetOrderByID(orderID string) *Order {
	mu.RLock()
	defer mu.RUnlock()
	for i := range Orders {
		if Orders[i].ID == orderID {
			return &Orders[i]
		}
	}
	return nil
}

// IsDeliverySource checks if the source is from a delivery platform
func IsDeliverySource(source string) bool {
	return source == SourceGrabFood || source == SourceGoFood || source == SourceShopeeFood
}
