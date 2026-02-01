package service

import (
	"github.com/kaori/backend/internal/config"
	"github.com/kaori/backend/internal/repository"
	"github.com/kaori/backend/internal/websocket"
)

// Services holds all service instances
type Services struct {
	Auth     *AuthService
	Store    *StoreService
	Table    *TableService
	Category *CategoryService
	Product  *ProductService
	Order    *OrderService
	Payment  *PaymentService
	Member   *MemberService
	Voucher  *VoucherService
	Report   *ReportService
	User     *UserService
}

// NewServices creates all service instances
func NewServices(repos *repository.Repositories, cfg *config.Config, hub *websocket.Hub) *Services {
	return &Services{
		Auth:     NewAuthService(cfg),
		Store:    NewStoreService(repos.Store),
		Table:    NewTableService(repos.Table),
		Category: NewCategoryService(repos.Category),
		Product:  NewProductService(repos.Product),
		Order:    NewOrderService(repos.Order, repos.Product, repos.Voucher, hub),
		Payment:  NewPaymentService(repos.Payment, repos.Order, cfg),
		Member:   NewMemberService(repos.Member),
		Voucher:  NewVoucherService(repos.Voucher),
		Report:   NewReportService(repos.Order, repos.Payment),
		User:     NewUserService(repos.User),
	}
}

// StoreService handles store business logic
type StoreService struct {
	repo *repository.StoreRepository
}

func NewStoreService(repo *repository.StoreRepository) *StoreService {
	return &StoreService{repo: repo}
}

// TableService handles table business logic
type TableService struct {
	repo *repository.TableRepository
}

func NewTableService(repo *repository.TableRepository) *TableService {
	return &TableService{repo: repo}
}

// CategoryService handles category business logic
type CategoryService struct {
	repo *repository.CategoryRepository
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

// ProductService handles product business logic
type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

// OrderService handles order business logic
type OrderService struct {
	repo        *repository.OrderRepository
	productRepo *repository.ProductRepository
	voucherRepo *repository.VoucherRepository
	hub         *websocket.Hub
}

func NewOrderService(repo *repository.OrderRepository, productRepo *repository.ProductRepository, voucherRepo *repository.VoucherRepository, hub *websocket.Hub) *OrderService {
	return &OrderService{
		repo:        repo,
		productRepo: productRepo,
		voucherRepo: voucherRepo,
		hub:         hub,
	}
}

// PaymentService handles payment business logic
type PaymentService struct {
	repo      *repository.PaymentRepository
	orderRepo *repository.OrderRepository
	cfg       *config.Config
}

func NewPaymentService(repo *repository.PaymentRepository, orderRepo *repository.OrderRepository, cfg *config.Config) *PaymentService {
	return &PaymentService{
		repo:      repo,
		orderRepo: orderRepo,
		cfg:       cfg,
	}
}

// MemberService handles membership business logic
type MemberService struct {
	repo *repository.MemberRepository
}

func NewMemberService(repo *repository.MemberRepository) *MemberService {
	return &MemberService{repo: repo}
}

// VoucherService handles voucher business logic
type VoucherService struct {
	repo *repository.VoucherRepository
}

func NewVoucherService(repo *repository.VoucherRepository) *VoucherService {
	return &VoucherService{repo: repo}
}

// ReportService handles report generation
type ReportService struct {
	orderRepo   *repository.OrderRepository
	paymentRepo *repository.PaymentRepository
}

func NewReportService(orderRepo *repository.OrderRepository, paymentRepo *repository.PaymentRepository) *ReportService {
	return &ReportService{
		orderRepo:   orderRepo,
		paymentRepo: paymentRepo,
	}
}

// UserService handles user management
type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}
