package model

// LoginRequest for email/password login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginPINRequest for PIN-based login
type LoginPINRequest struct {
	Email string `json:"email" binding:"required,email"`
	PIN   string `json:"pin" binding:"required,min=4,max=6"`
}

// LoginResponse returns tokens after successful login
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	User         User   `json:"user"`
}

// RefreshTokenRequest for refreshing access token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// CreateStoreRequest for creating a new store
type CreateStoreRequest struct {
	Name          string  `json:"name" binding:"required"`
	Code          string  `json:"code" binding:"required"`
	Address       *string `json:"address"`
	Phone         *string `json:"phone"`
	LogoURL       *string `json:"logo_url"`
	ReceiptHeader *string `json:"receipt_header"`
	ReceiptFooter *string `json:"receipt_footer"`
}

// CreateTableRequest for creating a new table
type CreateTableRequest struct {
	StoreID     string  `json:"store_id" binding:"required,uuid"`
	TableNumber string  `json:"table_number" binding:"required"`
	Capacity    int     `json:"capacity"`
	Location    *string `json:"location"`
}

// CreateUserRequest for creating a new staff user
type CreateUserRequest struct {
	Email    string  `json:"email" binding:"required,email"`
	Password string  `json:"password" binding:"required,min=6"`
	Name     string  `json:"name" binding:"required"`
	Role     string  `json:"role" binding:"required,oneof=cashier kitchen store_admin super_admin"`
	StoreID  *string `json:"store_id"`
	PIN      *string `json:"pin"`
}

// UpdateUserRequest for updating a user
type UpdateUserRequest struct {
	Name     *string `json:"name"`
	Role     *string `json:"role"`
	StoreID  *string `json:"store_id"`
	PIN      *string `json:"pin"`
	IsActive *bool   `json:"is_active"`
}

// CreateCategoryRequest for creating a category
type CreateCategoryRequest struct {
	StoreID   *string `json:"store_id"`
	Name      string  `json:"name" binding:"required"`
	Icon      *string `json:"icon"`
	SortOrder int     `json:"sort_order"`
}

// CreateProductRequest for creating a product
type CreateProductRequest struct {
	StoreID          *string                    `json:"store_id"`
	CategoryID       string                     `json:"category_id" binding:"required,uuid"`
	Name             string                     `json:"name" binding:"required"`
	Description      *string                    `json:"description"`
	BasePrice        float64                    `json:"base_price" binding:"required,min=0"`
	ImageURL         *string                    `json:"image_url"`
	HasVariants      bool                       `json:"has_variants"`
	HasModifiers     bool                       `json:"has_modifiers"`
	PointsMultiplier float64                    `json:"points_multiplier"`
	Variants         []CreateVariantRequest     `json:"variants"`
	Modifiers        []CreateModifierRequest    `json:"modifiers"`
}

// CreateVariantRequest for product variants
type CreateVariantRequest struct {
	Name            string  `json:"name" binding:"required"`
	PriceAdjustment float64 `json:"price_adjustment"`
	IsDefault       bool    `json:"is_default"`
	SortOrder       int     `json:"sort_order"`
}

// CreateModifierRequest for product modifiers
type CreateModifierRequest struct {
	Name      string  `json:"name" binding:"required"`
	Price     float64 `json:"price"`
	SortOrder int     `json:"sort_order"`
}

// CreateOrderRequest for creating a new order
type CreateOrderRequest struct {
	StoreID     string                   `json:"store_id" binding:"required,uuid"`
	OrderSource string                   `json:"order_source" binding:"required,oneof=table_qr client_app cashier"`
	OrderType   string                   `json:"order_type" binding:"required,oneof=dine_in takeaway"`
	TableID     *string                  `json:"table_id"`
	MemberID    *string                  `json:"member_id"`
	Items       []CreateOrderItemRequest `json:"items" binding:"required,min=1"`
	Notes       *string                  `json:"notes"`
	VoucherCode *string                  `json:"voucher_code"`
}

// CreateOrderItemRequest for order items
type CreateOrderItemRequest struct {
	ProductID   string   `json:"product_id" binding:"required,uuid"`
	VariantID   *string  `json:"variant_id"`
	ModifierIDs []string `json:"modifier_ids"`
	Quantity    int      `json:"quantity" binding:"required,min=1"`
	Notes       *string  `json:"notes"`
}

// UpdateOrderStatusRequest for kitchen status updates
type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending confirmed cooking ready completed cancelled"`
}

// SyncOrderRequest for offline order sync
type SyncOrderRequest struct {
	Orders []CreateOrderRequest `json:"orders" binding:"required"`
}

// ProcessCashPaymentRequest for cash payments
type ProcessCashPaymentRequest struct {
	OrderID    string  `json:"order_id" binding:"required,uuid"`
	Amount     float64 `json:"amount" binding:"required,min=0"`
	CashierID  string  `json:"cashier_id" binding:"required,uuid"`
}

// CreateMidtransPaymentRequest for digital payments
type CreateMidtransPaymentRequest struct {
	OrderID string `json:"order_id" binding:"required,uuid"`
	Method  string `json:"method" binding:"required,oneof=qris card ewallet"`
}

// CreateMemberRequest for registering a new member
type CreateMemberRequest struct {
	Phone string  `json:"phone" binding:"required"`
	Name  string  `json:"name" binding:"required"`
	Email *string `json:"email"`
}

// RedeemPointsRequest for redeeming points
type RedeemPointsRequest struct {
	Points      int    `json:"points" binding:"required,min=1"`
	Description string `json:"description"`
}

// ApplyVoucherRequest for applying a voucher
type ApplyVoucherRequest struct {
	OrderID     string `json:"order_id" binding:"required,uuid"`
	VoucherCode string `json:"voucher_code" binding:"required"`
}

// CreateVoucherRequest for creating a voucher
type CreateVoucherRequest struct {
	StoreID     *string  `json:"store_id"`
	Code        string   `json:"code" binding:"required"`
	Type        string   `json:"type" binding:"required,oneof=percentage fixed freeitem"`
	Value       float64  `json:"value" binding:"required,min=0"`
	MinPurchase float64  `json:"min_purchase"`
	MaxUses     *int     `json:"max_uses"`
	ValidFrom   *string  `json:"valid_from"`
	ValidUntil  *string  `json:"valid_until"`
}

// PaginationParams for list endpoints
type PaginationParams struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
	SortBy   string `form:"sort_by"`
	SortDir  string `form:"sort_dir" binding:"oneof=asc desc"`
}

// OrderFilterParams for filtering orders
type OrderFilterParams struct {
	PaginationParams
	StoreID   string `form:"store_id"`
	Status    string `form:"status"`
	Source    string `form:"source"`
	DateFrom  string `form:"date_from"`
	DateTo    string `form:"date_to"`
	CashierID string `form:"cashier_id"`
}
