package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/naxumi/bnsp-jwd/internal/config"
	appHTTP "github.com/naxumi/bnsp-jwd/internal/handler/http"
	"github.com/naxumi/bnsp-jwd/internal/pkg/database"
	"github.com/naxumi/bnsp-jwd/internal/pkg/storage"
	"github.com/naxumi/bnsp-jwd/internal/repository/postgresql"
	"github.com/naxumi/bnsp-jwd/internal/service/file"
	"github.com/naxumi/bnsp-jwd/internal/service/product"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	dsn := cfg.DatabaseURL()
	db, err := database.NewPostgreSQLDB(dsn)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}

	productRepo := postgresql.NewProductRepository(db)

	var fileStorage storage.FileStorage
	switch cfg.Storage.Type {
	case "local":
		fileStorage, err = storage.NewLocalStorage(
			cfg.Storage.BasePath,
			cfg.Storage.BaseURL,
		)
		if err != nil {
			log.Fatal("Failed to initialize local storage:", err)
		}
	case "minio":
		// Future: minIO implementation
		log.Fatal("Minio storage not yet implemented")
	default:
		log.Fatal("Unsupported storage types: ", cfg.Storage.Type)
	}

	fileService := file.NewFileService(fileStorage)
	productService := product.NewProductService(db, productRepo, fileService)

	productHandler := appHTTP.NewProductHandler(productService)

	router := appHTTP.NewRouter(
		productHandler,
		cfg.Storage.BasePath,
	)

	port := fmt.Sprintf(":%d", cfg.App.Port)
	fmt.Printf("Server running at http://localhost%s\n", port)
	if err := http.ListenAndServe(port, router); err != nil {
		fmt.Println("Server error:", err)
	}
}
