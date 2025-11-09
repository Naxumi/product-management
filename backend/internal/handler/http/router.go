package http

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v3"
)

func NewRouter(productHandler ProductHandler, storageBasePath string) *chi.Mux {
	r := chi.NewRouter()
	logFormat := httplog.SchemaECS.Concise(false)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: logFormat.ReplaceAttr,
	})).With(
		slog.String("app", "bnsp-jwd"),
		slog.String("version", "v1.0.0"),
		slog.String("env", "development"),
	)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		MaxAge:           300,
	}))

	// r.Use(chiMiddleware.RealIP)

	r.Use(httplog.RequestLogger(logger, &httplog.Options{
		Level:  slog.LevelDebug,
		Schema: httplog.SchemaECS,
	}))

	r.Use(chiMiddleware.AllowContentType("application/json", "multipart/form-data"))
	r.Use(chiMiddleware.CleanPath)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Heartbeat("/"))

	fileServer := http.FileServer(http.Dir(storageBasePath))
	r.Handle("/uploads/*", http.StripPrefix("/uploads/", fileServer))

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/product", func(r chi.Router) {
			r.Post("/", productHandler.CreateProduct)
			r.Get("/{id}", productHandler.GetProduct)
			r.Get("/sku/{sku}", productHandler.GetProductBySKU)
			r.Put("/", productHandler.UpdateProduct)
			r.Delete("/{id}", productHandler.DeleteProduct)
			r.Get("/", productHandler.ListProducts)
			r.Post("/{id}/image", productHandler.UploadImage)
			r.Delete("/{id}/image", productHandler.DeleteImage)
		})
	})
	return r
}
