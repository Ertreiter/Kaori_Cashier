package handler

import (
	"github.com/kaori/backend/internal/config"
	"github.com/kaori/backend/internal/service"
	"github.com/kaori/backend/internal/websocket"
)

// Handlers holds all handler instances
type Handlers struct {
	Auth     *AuthHandler
	Store    *StoreHandler
	Table    *TableHandler
	Category *CategoryHandler
	Product  *ProductHandler
	Order    *OrderHandler
	Payment  *PaymentHandler
	Member   *MemberHandler
	Voucher  *VoucherHandler
	Report   *ReportHandler
	User     *UserHandler
}

// NewHandlers creates all handler instances
// If services is nil, handlers will use stub implementations with dummy data
func NewHandlers(services *service.Services, hub *websocket.Hub) *Handlers {
	if services == nil {
		// Dummy data mode - create handlers with dummy auth service
		cfg := config.Load()
		authService := service.NewAuthService(cfg)
		return &Handlers{
			Auth:     NewAuthHandler(authService),
			Store:    &StoreHandler{},
			Table:    &TableHandler{},
			Category: &CategoryHandler{},
			Product:  &ProductHandler{},
			Order:    &OrderHandler{hub: hub},
			Payment:  &PaymentHandler{},
			Member:   &MemberHandler{},
			Voucher:  &VoucherHandler{},
			Report:   &ReportHandler{},
			User:     &UserHandler{},
		}
	}

	return &Handlers{
		Auth:     NewAuthHandler(services.Auth),
		Store:    NewStoreHandler(services.Store),
		Table:    NewTableHandler(services.Table),
		Category: NewCategoryHandler(services.Category),
		Product:  NewProductHandler(services.Product),
		Order:    NewOrderHandler(services.Order, hub),
		Payment:  NewPaymentHandler(services.Payment),
		Member:   NewMemberHandler(services.Member),
		Voucher:  NewVoucherHandler(services.Voucher),
		Report:   NewReportHandler(services.Report),
		User:     NewUserHandler(services.User),
	}
}

// StoreHandler handles store endpoints
type StoreHandler struct {
	service *service.StoreService
}

func NewStoreHandler(s *service.StoreService) *StoreHandler {
	return &StoreHandler{service: s}
}

// TableHandler handles table endpoints
type TableHandler struct {
	service *service.TableService
}

func NewTableHandler(s *service.TableService) *TableHandler {
	return &TableHandler{service: s}
}

// CategoryHandler handles category endpoints
type CategoryHandler struct {
	service *service.CategoryService
}

func NewCategoryHandler(s *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: s}
}

// ProductHandler handles product endpoints
type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(s *service.ProductService) *ProductHandler {
	return &ProductHandler{service: s}
}

// OrderHandler handles order endpoints
type OrderHandler struct {
	service *service.OrderService
	hub     *websocket.Hub
}

func NewOrderHandler(s *service.OrderService, hub *websocket.Hub) *OrderHandler {
	return &OrderHandler{service: s, hub: hub}
}

// PaymentHandler handles payment endpoints
type PaymentHandler struct {
	service *service.PaymentService
}

func NewPaymentHandler(s *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: s}
}

// MemberHandler handles member endpoints
type MemberHandler struct {
	service *service.MemberService
}

func NewMemberHandler(s *service.MemberService) *MemberHandler {
	return &MemberHandler{service: s}
}

// VoucherHandler handles voucher endpoints
type VoucherHandler struct {
	service *service.VoucherService
}

func NewVoucherHandler(s *service.VoucherService) *VoucherHandler {
	return &VoucherHandler{service: s}
}

// ReportHandler handles report endpoints
type ReportHandler struct {
	service *service.ReportService
}

func NewReportHandler(s *service.ReportService) *ReportHandler {
	return &ReportHandler{service: s}
}

// UserHandler handles user management endpoints
type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}
