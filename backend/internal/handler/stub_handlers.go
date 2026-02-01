package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kaori/backend/internal/dummy"
	"github.com/kaori/backend/pkg/response"
)

// --- Store Handler ---

func (h *StoreHandler) List(c *gin.Context) {
	stores := []gin.H{{
		"id": "store-1", "name": "Kaori Coffee", "address": "Jl. Sample No. 123",
		"phone": "0812345678", "is_active": true,
	}}
	response.Success(c, http.StatusOK, stores)
}

func (h *StoreHandler) GetByID(c *gin.Context) {
	response.Success(c, http.StatusOK, gin.H{
		"id": "store-1", "name": "Kaori Coffee", "address": "Jl. Sample No. 123",
		"phone": "0812345678", "is_active": true,
	})
}

func (h *StoreHandler) Create(c *gin.Context) {
	response.Success(c, http.StatusCreated, gin.H{"message": "Store created"})
}

func (h *StoreHandler) Update(c *gin.Context) {
	response.Success(c, http.StatusOK, gin.H{"message": "Store updated"})
}

func (h *StoreHandler) GetStats(c *gin.Context) {
	response.Success(c, http.StatusOK, gin.H{
		"total_orders_today":  len(dummy.Orders),
		"total_revenue_today": calculateTotalRevenue(),
		"active_orders":       len(dummy.GetOrdersByStatus("confirmed", "cooking", "ready")),
	})
}

func calculateTotalRevenue() int {
	total := 0
	for _, o := range dummy.Orders {
		if o.PaymentStatus == "paid" {
			total += o.Total
		}
	}
	return total
}

// --- Table Handler ---

func (h *TableHandler) List(c *gin.Context) {
	response.Success(c, http.StatusOK, dummy.Tables)
}

func (h *TableHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	for _, t := range dummy.Tables {
		if t.ID == id {
			response.Success(c, http.StatusOK, t)
			return
		}
	}
	response.NotFound(c, "Table not found")
}

func (h *TableHandler) Create(c *gin.Context) {
	response.Success(c, http.StatusCreated, gin.H{"message": "Table created"})
}

func (h *TableHandler) Update(c *gin.Context) {
	response.Success(c, http.StatusOK, gin.H{"message": "Table updated"})
}

func (h *TableHandler) Delete(c *gin.Context) {
	response.Success(c, http.StatusOK, gin.H{"message": "Table deleted"})
}

func (h *TableHandler) GetQRCode(c *gin.Context) {
	id := c.Param("id")
	response.Success(c, http.StatusOK, gin.H{
		"table_id": id,
		"qr_url":   "https://kaori.pos/order?table=" + id,
	})
}

func (h *TableHandler) GetPublicInfo(c *gin.Context) {
	id := c.Param("id")
	for _, t := range dummy.Tables {
		if t.ID == id {
			response.Success(c, http.StatusOK, gin.H{
				"table_number": t.Number,
				"store_name":   "Kaori Coffee",
				"store_id":     "store-1",
			})
			return
		}
	}
	response.NotFound(c, "Table not found")
}

// --- Category Handler ---

func (h *CategoryHandler) List(c *gin.Context) {
	response.Success(c, http.StatusOK, dummy.Categories)
}

func (h *CategoryHandler) Create(c *gin.Context) {
	response.Success(c, http.StatusCreated, gin.H{"message": "Category created"})
}

func (h *CategoryHandler) Update(c *gin.Context) {
	response.Success(c, http.StatusOK, gin.H{"message": "Category updated"})
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	response.Success(c, http.StatusOK, gin.H{"message": "Category deleted"})
}

// --- Product Handler ---

func (h *ProductHandler) List(c *gin.Context) {
	categoryID := c.Query("category_id")
	if categoryID != "" {
		var filtered []dummy.Product
		for _, p := range dummy.Products {
			if p.CategoryID == categoryID {
				filtered = append(filtered, p)
			}
		}
		response.Success(c, http.StatusOK, filtered)
		return
	}
	response.Success(c, http.StatusOK, dummy.Products)
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	for _, p := range dummy.Products {
		if p.ID == id {
			response.Success(c, http.StatusOK, p)
			return
		}
	}
	response.NotFound(c, "Product not found")
}

func (h *ProductHandler) Create(c *gin.Context) {
	response.Success(c, http.StatusCreated, gin.H{"message": "Product created"})
}

func (h *ProductHandler) Update(c *gin.Context) {
	response.Success(c, http.StatusOK, gin.H{"message": "Product updated"})
}

func (h *ProductHandler) ToggleAvailability(c *gin.Context) {
	response.Success(c, http.StatusOK, gin.H{"message": "Availability updated"})
}

// --- Order Handler ---

func (h *OrderHandler) List(c *gin.Context) {
	response.Success(c, http.StatusOK, dummy.Orders)
}

func (h *OrderHandler) GetActive(c *gin.Context) {
	active := dummy.GetOrdersByStatus("confirmed", "cooking", "ready")
	response.Success(c, http.StatusOK, active)
}

func (h *OrderHandler) GetIncoming(c *gin.Context) {
	pending := dummy.GetOrdersByStatus("pending")
	response.Success(c, http.StatusOK, pending)
}

func (h *OrderHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	for _, o := range dummy.Orders {
		if o.ID == id {
			response.Success(c, http.StatusOK, o)
			return
		}
	}
	response.NotFound(c, "Order not found")
}

type CreateOrderRequest struct {
	OrderType string `json:"order_type"` // dine_in, takeaway
	TableID   string `json:"table_id,omitempty"`
	Items     []struct {
		ProductID string   `json:"product_id"`
		VariantID string   `json:"variant_id,omitempty"`
		Modifiers []string `json:"modifiers,omitempty"`
		Quantity  int      `json:"quantity"`
		Notes     string   `json:"notes,omitempty"`
	} `json:"items"`
	Notes string `json:"notes,omitempty"`
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	// Build order items
	var items []dummy.OrderItem
	subtotal := 0
	for _, item := range req.Items {
		// Find product
		var product *dummy.Product
		for i := range dummy.Products {
			if dummy.Products[i].ID == item.ProductID {
				product = &dummy.Products[i]
				break
			}
		}
		if product == nil {
			continue
		}

		unitPrice := product.BasePrice
		variantName := ""

		// Apply variant price
		if item.VariantID != "" {
			for _, v := range product.Variants {
				if v.ID == item.VariantID {
					unitPrice += v.PriceAdjustment
					variantName = v.Name
					break
				}
			}
		}

		// Apply modifier prices
		for _, modID := range item.Modifiers {
			for _, m := range product.Modifiers {
				if m.ID == modID {
					unitPrice += m.Price
				}
			}
		}

		itemSubtotal := unitPrice * item.Quantity
		subtotal += itemSubtotal

		items = append(items, dummy.OrderItem{
			ID:          uuid.New().String(),
			ProductID:   product.ID,
			ProductName: product.Name,
			VariantID:   item.VariantID,
			VariantName: variantName,
			Modifiers:   item.Modifiers,
			Quantity:    item.Quantity,
			UnitPrice:   unitPrice,
			Subtotal:    itemSubtotal,
			Notes:       item.Notes,
		})
	}

	tax := subtotal * 11 / 100 // 11% tax
	total := subtotal + tax

	// Get table number if dine-in
	tableNumber := 0
	if req.OrderType == "dine_in" && req.TableID != "" {
		for _, t := range dummy.Tables {
			if t.ID == req.TableID {
				tableNumber = t.Number
				break
			}
		}
	}

	order := dummy.Order{
		ID:            uuid.New().String(),
		OrderNumber:   dummy.GetNextOrderNumber(),
		TableID:       req.TableID,
		TableNumber:   tableNumber,
		OrderType:     req.OrderType,
		OrderSource:   "cashier",
		Status:        "confirmed",
		PaymentStatus: "unpaid",
		Items:         items,
		Subtotal:      subtotal,
		Tax:           tax,
		Total:         total,
		Notes:         req.Notes,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		CashierID:     "user-3",
		CashierName:   "John Cashier",
	}

	dummy.AddOrder(order)

	// Broadcast to kitchen via WebSocket
	h.hub.BroadcastOrder(order.ID, "new_order", order)

	response.Success(c, http.StatusCreated, order)
}

func (h *OrderHandler) Confirm(c *gin.Context) {
	id := c.Param("id")
	if dummy.UpdateOrderStatus(id, "confirmed") {
		response.Success(c, http.StatusOK, gin.H{"message": "Order confirmed"})
		return
	}
	response.NotFound(c, "Order not found")
}

func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request")
		return
	}

	if dummy.UpdateOrderStatus(id, req.Status) {
		h.hub.BroadcastOrder(id, "order_status", gin.H{"id": id, "status": req.Status})
		response.Success(c, http.StatusOK, gin.H{"message": "Status updated", "status": req.Status})
		return
	}
	response.NotFound(c, "Order not found")
}

func (h *OrderHandler) Cancel(c *gin.Context) {
	id := c.Param("id")
	if dummy.UpdateOrderStatus(id, "cancelled") {
		response.Success(c, http.StatusOK, gin.H{"message": "Order cancelled"})
		return
	}
	response.NotFound(c, "Order not found")
}

func (h *OrderHandler) SyncOffline(c *gin.Context) {
	response.Success(c, http.StatusOK, gin.H{"synced": 0})
}

// --- Payment Handler ---

func (h *PaymentHandler) ProcessCash(c *gin.Context) {
	var req struct {
		OrderID    string `json:"order_id"`
		AmountPaid int    `json:"amount_paid"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request")
		return
	}

	// Find and update order
	for i := range dummy.Orders {
		if dummy.Orders[i].ID == req.OrderID {
			dummy.Orders[i].PaymentStatus = "paid"
			change := req.AmountPaid - dummy.Orders[i].Total
			response.Success(c, http.StatusOK, gin.H{
				"message":     "Payment successful",
				"order_id":    req.OrderID,
				"total":       dummy.Orders[i].Total,
				"amount_paid": req.AmountPaid,
				"change":      change,
			})
			return
		}
	}
	response.NotFound(c, "Order not found")
}

func (h *PaymentHandler) CreateMidtrans(c *gin.Context) {
	response.Success(c, http.StatusOK, gin.H{
		"payment_url":    "https://midtrans.com/pay/sample",
		"transaction_id": uuid.New().String(),
	})
}

func (h *PaymentHandler) MidtransCallback(c *gin.Context) {
	response.Success(c, http.StatusOK, gin.H{"message": "Callback received"})
}

func (h *PaymentHandler) GetStatus(c *gin.Context) {
	response.Success(c, http.StatusOK, gin.H{"status": "pending"})
}

// --- Member Handler ---

func (h *MemberHandler) Lookup(c *gin.Context) {
	response.NotFound(c, "Member not found")
}

func (h *MemberHandler) Create(c *gin.Context) {
	response.Success(c, http.StatusCreated, gin.H{"message": "Member created"})
}

func (h *MemberHandler) GetPoints(c *gin.Context) {
	response.Success(c, http.StatusOK, []interface{}{})
}

func (h *MemberHandler) Redeem(c *gin.Context) {
	response.Success(c, http.StatusOK, gin.H{"message": "Points redeemed"})
}

// --- Voucher Handler ---

func (h *VoucherHandler) Validate(c *gin.Context) {
	response.Success(c, http.StatusOK, gin.H{"valid": false, "message": "Voucher not found"})
}

func (h *VoucherHandler) Apply(c *gin.Context) {
	response.Success(c, http.StatusOK, gin.H{"discount": 0})
}

// --- Report Handler ---

func (h *ReportHandler) GetDaily(c *gin.Context) {
	totalOrders := len(dummy.Orders)
	totalRevenue := 0
	for _, o := range dummy.Orders {
		if o.PaymentStatus == "paid" {
			totalRevenue += o.Total
		}
	}
	response.Success(c, http.StatusOK, gin.H{
		"date":          time.Now().Format("2006-01-02"),
		"total_orders":  totalOrders,
		"total_revenue": totalRevenue,
	})
}

func (h *ReportHandler) GetDailyByDate(c *gin.Context) {
	response.Success(c, http.StatusOK, gin.H{})
}

func (h *ReportHandler) GetProductSales(c *gin.Context) {
	response.Success(c, http.StatusOK, []interface{}{})
}

func (h *ReportHandler) GetCashierSales(c *gin.Context) {
	response.Success(c, http.StatusOK, []interface{}{})
}

func (h *ReportHandler) GetHourly(c *gin.Context) {
	response.Success(c, http.StatusOK, []interface{}{})
}

// --- User Handler ---

func (h *UserHandler) List(c *gin.Context) {
	response.Success(c, http.StatusOK, dummy.Users)
}

func (h *UserHandler) Create(c *gin.Context) {
	response.Success(c, http.StatusCreated, gin.H{"message": "User created"})
}

func (h *UserHandler) Update(c *gin.Context) {
	response.Success(c, http.StatusOK, gin.H{"message": "User updated"})
}

func (h *UserHandler) Delete(c *gin.Context) {
	response.Success(c, http.StatusOK, gin.H{"message": "User deleted"})
}
