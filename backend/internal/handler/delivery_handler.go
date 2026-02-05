package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kaori/backend/internal/dummy"
	"github.com/kaori/backend/internal/websocket"
	"github.com/kaori/backend/pkg/response"
)

// ItemInput is a common type for order items
type ItemInput struct {
	Name     string
	Quantity int
	Price    int
	Notes    string
}

// DeliveryHandler handles delivery platform webhooks
type DeliveryHandler struct {
	hub *websocket.Hub
}

// NewDeliveryHandler creates a new delivery handler
func NewDeliveryHandler(hub *websocket.Hub) *DeliveryHandler {
	return &DeliveryHandler{hub: hub}
}

// GrabFood webhook - POST /api/webhooks/grabfood
type GrabFoodOrder struct {
	OrderID       string `json:"orderId"`
	CustomerName  string `json:"customerName"`
	CustomerPhone string `json:"customerPhone"`
	Address       string `json:"address"`
	Items         []struct {
		Name     string `json:"name"`
		Quantity int    `json:"quantity"`
		Price    int    `json:"price"`
		Notes    string `json:"notes,omitempty"`
	} `json:"items"`
	DriverName string `json:"driverName,omitempty"`
	Total      int    `json:"total"`
}

func (h *DeliveryHandler) HandleGrabFood(c *gin.Context) {
	var req GrabFoodOrder
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid GrabFood order format")
		return
	}

	items := make([]ItemInput, len(req.Items))
	for i, item := range req.Items {
		items[i] = ItemInput{Name: item.Name, Quantity: item.Quantity, Price: item.Price, Notes: item.Notes}
	}

	order := h.createDeliveryOrder(req.OrderID, dummy.SourceGrabFood, req.CustomerName, req.CustomerPhone, req.Address, req.DriverName, items, req.Total)
	dummy.AddOrder(order)
	h.hub.BroadcastOrder(order.ID, "new_order", order)

	response.Success(c, http.StatusCreated, gin.H{
		"status":       "accepted",
		"order_id":     order.ID,
		"order_number": order.OrderNumber,
	})
}

// GoFood webhook - POST /api/webhooks/gofood
type GoFoodOrder struct {
	TransactionID string `json:"transaction_id"`
	Customer      struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
	} `json:"customer"`
	DeliveryAddress string `json:"delivery_address"`
	Items           []struct {
		ProductName string `json:"product_name"`
		Qty         int    `json:"qty"`
		Price       int    `json:"price"`
		Note        string `json:"note,omitempty"`
	} `json:"items"`
	Driver struct {
		Name string `json:"name"`
	} `json:"driver"`
	TotalAmount int `json:"total_amount"`
}

func (h *DeliveryHandler) HandleGoFood(c *gin.Context) {
	var req GoFoodOrder
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid GoFood order format")
		return
	}

	items := make([]ItemInput, len(req.Items))
	for i, item := range req.Items {
		items[i] = ItemInput{Name: item.ProductName, Quantity: item.Qty, Price: item.Price, Notes: item.Note}
	}

	order := h.createDeliveryOrder(req.TransactionID, dummy.SourceGoFood, req.Customer.Name, req.Customer.Phone, req.DeliveryAddress, req.Driver.Name, items, req.TotalAmount)
	dummy.AddOrder(order)
	h.hub.BroadcastOrder(order.ID, "new_order", order)

	response.Success(c, http.StatusCreated, gin.H{
		"status":       "accepted",
		"order_id":     order.ID,
		"order_number": order.OrderNumber,
	})
}

// Shopee Food webhook - POST /api/webhooks/shopee
type ShopeeFoodOrder struct {
	OrderNo    string `json:"order_no"`
	BuyerName  string `json:"buyer_name"`
	BuyerPhone string `json:"buyer_phone"`
	Address    struct {
		Full string `json:"full"`
	} `json:"address"`
	OrderItems []struct {
		ItemName string `json:"item_name"`
		Quantity int    `json:"quantity"`
		Price    int    `json:"price"`
		Remark   string `json:"remark,omitempty"`
	} `json:"order_items"`
	ShipperName string `json:"shipper_name"`
	TotalPrice  int    `json:"total_price"`
}

func (h *DeliveryHandler) HandleShopeeFood(c *gin.Context) {
	var req ShopeeFoodOrder
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid Shopee Food order format")
		return
	}

	items := make([]ItemInput, len(req.OrderItems))
	for i, item := range req.OrderItems {
		items[i] = ItemInput{Name: item.ItemName, Quantity: item.Quantity, Price: item.Price, Notes: item.Remark}
	}

	order := h.createDeliveryOrder(req.OrderNo, dummy.SourceShopeeFood, req.BuyerName, req.BuyerPhone, req.Address.Full, req.ShipperName, items, req.TotalPrice)
	dummy.AddOrder(order)
	h.hub.BroadcastOrder(order.ID, "new_order", order)

	response.Success(c, http.StatusCreated, gin.H{
		"status":       "accepted",
		"order_id":     order.ID,
		"order_number": order.OrderNumber,
	})
}

// Simulate incoming order (for testing) - POST /api/simulate/order
type SimulateOrderRequest struct {
	Source       string `json:"source" binding:"required"`
	CustomerName string `json:"customer_name"`
	Items        []struct {
		Name     string `json:"name"`
		Quantity int    `json:"quantity"`
		Price    int    `json:"price"`
	} `json:"items"`
}

func (h *DeliveryHandler) SimulateOrder(c *gin.Context) {
	var req SimulateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request")
		return
	}

	if !dummy.IsDeliverySource(req.Source) && req.Source != dummy.SourceCashier && req.Source != dummy.SourceTableQR {
		response.BadRequest(c, "Invalid source. Use: cashier, table_qr, grabfood, gofood, shopee_food")
		return
	}

	items := make([]ItemInput, len(req.Items))
	total := 0
	for i, item := range req.Items {
		items[i] = ItemInput{Name: item.Name, Quantity: item.Quantity, Price: item.Price, Notes: ""}
		total += item.Price * item.Quantity
	}

	extID := uuid.New().String()[:8]
	order := h.createDeliveryOrder(extID, req.Source, req.CustomerName, "08123456789", "Jl. Delivery No. 123", "Driver", items, total)
	dummy.AddOrder(order)
	h.hub.BroadcastOrder(order.ID, "new_order", order)

	response.Success(c, http.StatusCreated, order)
}

// GetOrdersBySource - GET /api/orders/source/:source
func (h *DeliveryHandler) GetOrdersBySource(c *gin.Context) {
	source := c.Param("source")
	orders := dummy.GetOrdersBySource(source)
	response.Success(c, http.StatusOK, orders)
}

// Helper to create delivery order
func (h *DeliveryHandler) createDeliveryOrder(externalID, source, customerName, phone, address, driver string, items []ItemInput, total int) dummy.Order {
	orderItems := make([]dummy.OrderItem, len(items))
	subtotal := 0
	for i, item := range items {
		itemTotal := item.Price * item.Quantity
		subtotal += itemTotal
		orderItems[i] = dummy.OrderItem{
			ID:          uuid.New().String(),
			ProductID:   "",
			ProductName: item.Name,
			Quantity:    item.Quantity,
			UnitPrice:   item.Price,
			Subtotal:    itemTotal,
			Notes:       item.Notes,
		}
	}

	if total == 0 {
		total = subtotal
	}
	tax := total * 11 / 100

	return dummy.Order{
		ID:              uuid.New().String(),
		OrderNumber:     dummy.GetNextDeliveryOrderNumber(source),
		ExternalOrderID: externalID,
		OrderType:       "delivery",
		OrderSource:     source,
		Status:          "pending",
		PaymentStatus:   "paid",
		Items:           orderItems,
		Subtotal:        subtotal,
		Tax:             tax,
		Total:           total,
		CustomerName:    customerName,
		CustomerPhone:   phone,
		DeliveryAddress: address,
		DriverName:      driver,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}
