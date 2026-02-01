package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/kaori/backend/internal/config"
	"github.com/kaori/backend/internal/handler"
	"github.com/kaori/backend/internal/middleware"
	"github.com/kaori/backend/internal/repository"
	"github.com/kaori/backend/internal/service"
	"github.com/kaori/backend/internal/websocket"
	"github.com/kaori/backend/pkg/database"
)

func main() {
	// Load .env file in development
	if os.Getenv("GIN_MODE") != "release" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, using environment variables")
		}
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Initialize repositories
	repos := repository.NewRepositories(db)

	// Initialize services
	services := service.NewServices(repos, cfg, hub)

	// Initialize handlers
	handlers := handler.NewHandlers(services, hub)

	// Setup Gin router
	if cfg.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSAllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	// Root endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "Kaori POS API",
			"version": "1.0.0",
			"status":  "online",
			"docs":    "/api/health for health check",
		})
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "kaori-api"})
	})

	// API routes
	api := r.Group("/api")
	{
		// Health check under /api as well
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok", "service": "kaori-api"})
		})
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/login", handlers.Auth.Login)
			auth.POST("/login/pin", handlers.Auth.LoginWithPIN)
			auth.POST("/refresh", handlers.Auth.RefreshToken)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			// Auth
			protected.GET("/auth/me", handlers.Auth.Me)
			protected.POST("/auth/logout", handlers.Auth.Logout)

			// Stores
			stores := protected.Group("/stores")
			{
				stores.GET("", handlers.Store.List)
				stores.GET("/:id", handlers.Store.GetByID)
				stores.POST("", middleware.RequireRole("super_admin"), handlers.Store.Create)
				stores.PUT("/:id", middleware.RequireRole("super_admin"), handlers.Store.Update)
				stores.GET("/:id/stats", handlers.Store.GetStats)
			}

			// Tables
			tables := protected.Group("/tables")
			{
				tables.GET("", handlers.Table.List)
				tables.GET("/:id", handlers.Table.GetByID)
				tables.POST("", middleware.RequireRole("store_admin", "super_admin"), handlers.Table.Create)
				tables.PUT("/:id", middleware.RequireRole("store_admin", "super_admin"), handlers.Table.Update)
				tables.DELETE("/:id", middleware.RequireRole("store_admin", "super_admin"), handlers.Table.Delete)
				tables.GET("/:id/qr", handlers.Table.GetQRCode)
			}

			// Categories
			categories := protected.Group("/categories")
			{
				categories.GET("", handlers.Category.List)
				categories.POST("", middleware.RequireRole("store_admin", "super_admin"), handlers.Category.Create)
				categories.PUT("/:id", middleware.RequireRole("store_admin", "super_admin"), handlers.Category.Update)
				categories.DELETE("/:id", middleware.RequireRole("store_admin", "super_admin"), handlers.Category.Delete)
			}

			// Products
			products := protected.Group("/products")
			{
				products.GET("", handlers.Product.List)
				products.GET("/:id", handlers.Product.GetByID)
				products.POST("", middleware.RequireRole("store_admin", "super_admin"), handlers.Product.Create)
				products.PUT("/:id", middleware.RequireRole("store_admin", "super_admin"), handlers.Product.Update)
				products.PATCH("/:id/availability", middleware.RequireRole("store_admin", "super_admin", "cashier"), handlers.Product.ToggleAvailability)
			}

			// Orders
			orders := protected.Group("/orders")
			{
				orders.GET("", handlers.Order.List)
				orders.GET("/active", handlers.Order.GetActive)
				orders.GET("/incoming", handlers.Order.GetIncoming)
				orders.GET("/:id", handlers.Order.GetByID)
				orders.POST("", handlers.Order.Create)
				orders.PATCH("/:id/confirm", middleware.RequireRole("cashier", "store_admin", "super_admin"), handlers.Order.Confirm)
				orders.PATCH("/:id/status", handlers.Order.UpdateStatus)
				orders.POST("/:id/cancel", handlers.Order.Cancel)
				orders.POST("/sync", handlers.Order.SyncOffline)
			}

			// Payments
			payments := protected.Group("/payments")
			{
				payments.POST("/cash", handlers.Payment.ProcessCash)
				payments.POST("/midtrans", handlers.Payment.CreateMidtrans)
				payments.GET("/:id/status", handlers.Payment.GetStatus)
			}

			// Members (schema ready, feature on hold)
			members := protected.Group("/members")
			{
				members.GET("/lookup", handlers.Member.Lookup)
				members.POST("", handlers.Member.Create)
				members.GET("/:id/points", handlers.Member.GetPoints)
				members.POST("/:id/redeem", handlers.Member.Redeem)
			}

			// Vouchers
			vouchers := protected.Group("/vouchers")
			{
				vouchers.GET("/validate/:code", handlers.Voucher.Validate)
				vouchers.POST("/apply", handlers.Voucher.Apply)
			}

			// Reports
			reports := protected.Group("/reports")
			reports.Use(middleware.RequireRole("store_admin", "super_admin"))
			{
				reports.GET("/daily", handlers.Report.GetDaily)
				reports.GET("/daily/:date", handlers.Report.GetDailyByDate)
				reports.GET("/products", handlers.Report.GetProductSales)
				reports.GET("/cashiers", handlers.Report.GetCashierSales)
				reports.GET("/hourly", handlers.Report.GetHourly)
			}

			// Users (staff management)
			users := protected.Group("/users")
			users.Use(middleware.RequireRole("store_admin", "super_admin"))
			{
				users.GET("", handlers.User.List)
				users.POST("", handlers.User.Create)
				users.PUT("/:id", handlers.User.Update)
				users.DELETE("/:id", handlers.User.Delete)
			}
		}

		// WebSocket for real-time updates
		api.GET("/ws", func(c *gin.Context) {
			websocket.ServeWS(hub, c.Writer, c.Request)
		})
	}

	// Midtrans webhook (public)
	r.POST("/api/payments/midtrans/callback", handlers.Payment.MidtransCallback)

	// Public table info for QR ordering
	r.GET("/api/public/tables/:id", handlers.Table.GetPublicInfo)

	// Start server
	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Kaori API starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
