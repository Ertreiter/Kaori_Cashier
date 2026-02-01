package model

import (
	"time"

	"github.com/google/uuid"
)

// Store represents a store/branch
type Store struct {
	ID            uuid.UUID `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"`
	Code          string    `json:"code" db:"code"`
	Address       *string   `json:"address" db:"address"`
	Phone         *string   `json:"phone" db:"phone"`
	LogoURL       *string   `json:"logo_url" db:"logo_url"`
	ReceiptHeader *string   `json:"receipt_header" db:"receipt_header"`
	ReceiptFooter *string   `json:"receipt_footer" db:"receipt_footer"`
	IsActive      bool      `json:"is_active" db:"is_active"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// Table represents a dining table with QR code
type Table struct {
	ID          uuid.UUID `json:"id" db:"id"`
	StoreID     uuid.UUID `json:"store_id" db:"store_id"`
	TableNumber string    `json:"table_number" db:"table_number"`
	QRCodeURL   *string   `json:"qr_code_url" db:"qr_code_url"`
	Capacity    int       `json:"capacity" db:"capacity"`
	Location    *string   `json:"location" db:"location"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// User represents a staff member
type User struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	Email        string     `json:"email" db:"email"`
	PasswordHash string     `json:"-" db:"password_hash"`
	Name         string     `json:"name" db:"name"`
	Role         string     `json:"role" db:"role"`
	StoreID      *uuid.UUID `json:"store_id" db:"store_id"`
	PIN          *string    `json:"-" db:"pin"`
	IsActive     bool       `json:"is_active" db:"is_active"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// Member represents a customer membership
type Member struct {
	ID             uuid.UUID `json:"id" db:"id"`
	Phone          string    `json:"phone" db:"phone"`
	Name           string    `json:"name" db:"name"`
	Email          *string   `json:"email" db:"email"`
	TotalPoints    int       `json:"total_points" db:"total_points"`
	LifetimePoints int       `json:"lifetime_points" db:"lifetime_points"`
	Tier           string    `json:"tier" db:"tier"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// MemberPoints represents points history
type MemberPoints struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	MemberID    uuid.UUID  `json:"member_id" db:"member_id"`
	OrderID     *uuid.UUID `json:"order_id" db:"order_id"`
	Points      int        `json:"points" db:"points"`
	Type        string     `json:"type" db:"type"`
	Description *string    `json:"description" db:"description"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
}

// Category represents a product category
type Category struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	StoreID   *uuid.UUID `json:"store_id" db:"store_id"`
	Name      string     `json:"name" db:"name"`
	Icon      *string    `json:"icon" db:"icon"`
	SortOrder int        `json:"sort_order" db:"sort_order"`
	IsActive  bool       `json:"is_active" db:"is_active"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

// Product represents a menu item
type Product struct {
	ID               uuid.UUID         `json:"id" db:"id"`
	StoreID          *uuid.UUID        `json:"store_id" db:"store_id"`
	CategoryID       uuid.UUID         `json:"category_id" db:"category_id"`
	Name             string            `json:"name" db:"name"`
	Description      *string           `json:"description" db:"description"`
	BasePrice        float64           `json:"base_price" db:"base_price"`
	ImageURL         *string           `json:"image_url" db:"image_url"`
	IsAvailable      bool              `json:"is_available" db:"is_available"`
	HasVariants      bool              `json:"has_variants" db:"has_variants"`
	HasModifiers     bool              `json:"has_modifiers" db:"has_modifiers"`
	PointsMultiplier float64           `json:"points_multiplier" db:"points_multiplier"`
	CreatedAt        time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at" db:"updated_at"`
	Variants         []ProductVariant  `json:"variants,omitempty"`
	Modifiers        []ProductModifier `json:"modifiers,omitempty"`
	Category         *Category         `json:"category,omitempty"`
}

// ProductVariant represents size/variant options
type ProductVariant struct {
	ID              uuid.UUID `json:"id" db:"id"`
	ProductID       uuid.UUID `json:"product_id" db:"product_id"`
	Name            string    `json:"name" db:"name"`
	PriceAdjustment float64   `json:"price_adjustment" db:"price_adjustment"`
	IsDefault       bool      `json:"is_default" db:"is_default"`
	IsAvailable     bool      `json:"is_available" db:"is_available"`
	SortOrder       int       `json:"sort_order" db:"sort_order"`
}

// ProductModifier represents add-ons/toppings
type ProductModifier struct {
	ID          uuid.UUID `json:"id" db:"id"`
	ProductID   uuid.UUID `json:"product_id" db:"product_id"`
	Name        string    `json:"name" db:"name"`
	Price       float64   `json:"price" db:"price"`
	IsAvailable bool      `json:"is_available" db:"is_available"`
	SortOrder   int       `json:"sort_order" db:"sort_order"`
}

// Order represents a customer order
type Order struct {
	ID            uuid.UUID   `json:"id" db:"id"`
	StoreID       uuid.UUID   `json:"store_id" db:"store_id"`
	OrderNumber   string      `json:"order_number" db:"order_number"`
	OrderSource   string      `json:"order_source" db:"order_source"`
	OrderType     string      `json:"order_type" db:"order_type"`
	TableID       *uuid.UUID  `json:"table_id" db:"table_id"`
	MemberID      *uuid.UUID  `json:"member_id" db:"member_id"`
	CashierID     *uuid.UUID  `json:"cashier_id" db:"cashier_id"`
	Status        string      `json:"status" db:"status"`
	PaymentStatus string      `json:"payment_status" db:"payment_status"`
	Subtotal      float64     `json:"subtotal" db:"subtotal"`
	Discount      float64     `json:"discount" db:"discount"`
	Total         float64     `json:"total" db:"total"`
	PointsEarned  int         `json:"points_earned" db:"points_earned"`
	Notes         *string     `json:"notes" db:"notes"`
	CreatedAt     time.Time   `json:"created_at" db:"created_at"`
	ConfirmedAt   *time.Time  `json:"confirmed_at" db:"confirmed_at"`
	CompletedAt   *time.Time  `json:"completed_at" db:"completed_at"`
	Items         []OrderItem `json:"items,omitempty"`
	Table         *Table      `json:"table,omitempty"`
	Cashier       *User       `json:"cashier,omitempty"`
	Member        *Member     `json:"member,omitempty"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID             uuid.UUID           `json:"id" db:"id"`
	OrderID        uuid.UUID           `json:"order_id" db:"order_id"`
	ProductID      uuid.UUID           `json:"product_id" db:"product_id"`
	VariantID      *uuid.UUID          `json:"variant_id" db:"variant_id"`
	ProductName    string              `json:"product_name" db:"product_name"`
	VariantName    *string             `json:"variant_name" db:"variant_name"`
	BasePrice      float64             `json:"base_price" db:"base_price"`
	VariantPrice   float64             `json:"variant_price" db:"variant_price"`
	ModifiersPrice float64             `json:"modifiers_price" db:"modifiers_price"`
	Quantity       int                 `json:"quantity" db:"quantity"`
	Notes          *string             `json:"notes" db:"notes"`
	Modifiers      []OrderItemModifier `json:"modifiers,omitempty"`
}

// OrderItemModifier represents a modifier applied to an order item
type OrderItemModifier struct {
	ID           uuid.UUID `json:"id" db:"id"`
	OrderItemID  uuid.UUID `json:"order_item_id" db:"order_item_id"`
	ModifierID   uuid.UUID `json:"modifier_id" db:"modifier_id"`
	ModifierName string    `json:"modifier_name" db:"modifier_name"`
	Price        float64   `json:"price" db:"price"`
}

// Payment represents a payment transaction
type Payment struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	OrderID    uuid.UUID  `json:"order_id" db:"order_id"`
	Method     string     `json:"method" db:"method"`
	Amount     float64    `json:"amount" db:"amount"`
	MidtransID *string    `json:"midtrans_id" db:"midtrans_id"`
	Status     string     `json:"status" db:"status"`
	PaidAt     *time.Time `json:"paid_at" db:"paid_at"`
}

// Voucher represents a discount voucher
type Voucher struct {
	ID               uuid.UUID  `json:"id" db:"id"`
	StoreID          *uuid.UUID `json:"store_id" db:"store_id"`
	Code             string     `json:"code" db:"code"`
	Type             string     `json:"type" db:"type"`
	Value            float64    `json:"value" db:"value"`
	MinPurchase      float64    `json:"min_purchase" db:"min_purchase"`
	MaxUses          *int       `json:"max_uses" db:"max_uses"`
	CurrentUses      int        `json:"current_uses" db:"current_uses"`
	ValidFrom        *time.Time `json:"valid_from" db:"valid_from"`
	ValidUntil       *time.Time `json:"valid_until" db:"valid_until"`
	AutoGenerateRule *string    `json:"auto_generate_rule" db:"auto_generate_rule"`
	IsActive         bool       `json:"is_active" db:"is_active"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
}

// VoucherUsage represents voucher usage history
type VoucherUsage struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	VoucherID       uuid.UUID  `json:"voucher_id" db:"voucher_id"`
	OrderID         uuid.UUID  `json:"order_id" db:"order_id"`
	MemberID        *uuid.UUID `json:"member_id" db:"member_id"`
	DiscountApplied float64    `json:"discount_applied" db:"discount_applied"`
	UsedAt          time.Time  `json:"used_at" db:"used_at"`
}

// RefreshToken represents a JWT refresh token
type RefreshToken struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
