package bootstrap

import (
	_ "codegen/docs"
	"codegen/internal/database"
	"codegen/internal/handler"
	"codegen/internal/repository"
	"codegen/internal/service"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"

	echoSwagger "github.com/swaggo/echo-swagger"
)

// Register the routes
func RegisterRoutes(e *echo.Echo, db *database.DB, conf *Config) {

	authorize := echojwt.JWT([]byte(conf.JWT.SigningKey))
	api := e.Group("/api")

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	healthHandler := handler.NewHealthHandler(service.NewHealthService(db))
	e.GET("/health", healthHandler.HealthCheck)

	// handlers
	productHandler := handler.NewProductHandler(service.NewProductService(repository.NewProductRepository(db)))

	// routes
	products := api.Group("/products") // authorize
	{
		products.GET("/search", productHandler.Search)
		products.GET("", productHandler.GetList)
		products.GET("/:id", productHandler.GetSingle)
		products.POST("", productHandler.Create, authorize)
		products.PUT("/:id", productHandler.Update, authorize)
		products.DELETE("/:id", productHandler.Delete, authorize)
	}

}
