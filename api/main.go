package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Product model
type Product struct {
	ID       uint    `gorm:"primaryKey"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

// Database instance
var DB *gorm.DB

// Initialize database
func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	var errDB error
	DB, errDB = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if errDB != nil {
		log.Fatalf("Failed to connect to database: %v", errDB)
	}

	DB.AutoMigrate(&Product{})
}

// Product Handlers
func GetProducts(c *gin.Context) {
	var products []Product
	DB.Find(&products)
	c.JSON(200, products)
}

func CreateProduct(c *gin.Context) {
	var product Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	DB.Create(&product)
	c.JSON(200, product)
}

func main() {
	InitDB()
	router := gin.Default()

	router.GET("/products", GetProducts)
	router.POST("/products", CreateProduct)

	router.Run(":8080")
}
