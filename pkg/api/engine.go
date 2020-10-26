package api

import (
	"coding-challenge-go/pkg/api/product"
	"coding-challenge-go/pkg/api/seller"
	"database/sql"
	"github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

// CreateAPIEngine creates engine instance that serves API endpoints,
// consider it as a router for incoming requests.
func CreateAPIEngine(db *sql.DB) (*gin.Engine, error) {
	r := gin.New()

	logger, _ := zap.NewProduction()

	// Add a ginzap middleware, which:
	//   - RFC3339 with UTC time format.
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	// Logs all panic to error log
	//   - stack means whether output the stack info.
	r.Use(ginzap.RecoveryWithZap(logger, true))
	v1 := r.Group("api/v1")
	productRepository := product.NewRepository(db)
	sellerRepository := seller.NewRepository(db)
	emailProvider := seller.NewEmailProvider()
	productController := product.NewController(productRepository, sellerRepository, emailProvider)
	{
		v1.GET("products", productController.List)
		v1.GET("product", productController.Get)
		v1.POST("product", productController.Post)
		v1.PUT("product", productController.Put)
		v1.DELETE("product", productController.Delete)
		sellerController := seller.NewController(sellerRepository)
		v1.GET("sellers", sellerController.List)
	}
	v2 := r.Group("api/v2")
	{
		v2.GET("products", productController.ListV2)
		v2.GET("product", productController.GetV2)
		v2.POST("product", productController.Post)
		v2.PUT("product", productController.Put)
		v2.GET("sellers/top10", productController.ListTopTen)
	}
	return r, nil
}
