package repository

import (
	"database/sql"
)

// Repositories holds all repository instances
type Repositories struct {
	User     *UserRepository
	Store    *StoreRepository
	Table    *TableRepository
	Category *CategoryRepository
	Product  *ProductRepository
	Order    *OrderRepository
	Payment  *PaymentRepository
	Member   *MemberRepository
	Voucher  *VoucherRepository
}

// NewRepositories creates all repository instances
func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		User:     NewUserRepository(db),
		Store:    NewStoreRepository(db),
		Table:    NewTableRepository(db),
		Category: NewCategoryRepository(db),
		Product:  NewProductRepository(db),
		Order:    NewOrderRepository(db),
		Payment:  NewPaymentRepository(db),
		Member:   NewMemberRepository(db),
		Voucher:  NewVoucherRepository(db),
	}
}

// UserRepository handles user database operations
type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// StoreRepository handles store database operations
type StoreRepository struct {
	db *sql.DB
}

func NewStoreRepository(db *sql.DB) *StoreRepository {
	return &StoreRepository{db: db}
}

// TableRepository handles table database operations
type TableRepository struct {
	db *sql.DB
}

func NewTableRepository(db *sql.DB) *TableRepository {
	return &TableRepository{db: db}
}

// CategoryRepository handles category database operations
type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// ProductRepository handles product database operations
type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// OrderRepository handles order database operations
type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// PaymentRepository handles payment database operations
type PaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

// MemberRepository handles member database operations
type MemberRepository struct {
	db *sql.DB
}

func NewMemberRepository(db *sql.DB) *MemberRepository {
	return &MemberRepository{db: db}
}

// VoucherRepository handles voucher database operations
type VoucherRepository struct {
	db *sql.DB
}

func NewVoucherRepository(db *sql.DB) *VoucherRepository {
	return &VoucherRepository{db: db}
}
