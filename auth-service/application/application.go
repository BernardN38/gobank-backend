package application

import (
	"database/sql"
	"github.com/bernardn38/gobank/auth-service/handler"
	"github.com/bernardn38/gobank/auth-service/sql/users"
	"github.com/bernardn38/gobank/auth-service/token"
	"github.com/cristalhq/jwt/v4"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

type Config struct {
	jwtSecretKey     string
	jwtSigningMethod jwt.Algorithm
}
type App struct {
	srv          server
	pgDb         *sql.DB
	tokenManager *token.Manager
}

type server struct {
	router  *chi.Mux
	handler *handler.Handler
}

func New() *App {
	app := App{}
	config := Config{jwtSecretKey: "superSecretKey", jwtSigningMethod: jwt.HS256}
	app.runAppSetup(config)
	return &app
}
func (app *App) Run() {
	log.Fatal(http.ListenAndServe(":80", app.srv.router))
}

func (app *App) runAppSetup(config Config) {
	db, err := sql.Open("postgres", "user=bernardn host=host.docker.internal dbname=identity_service sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	queries := users.New(db)
	tokenManger := token.NewManager([]byte(config.jwtSecretKey), config.jwtSigningMethod)
	h := &handler.Handler{UsersDb: queries, TokenManager: tokenManger}

	app.srv.router = SetupRouter(h, tokenManger)
	app.pgDb = db
	app.tokenManager = tokenManger
	app.srv.handler = h
}

func SetupRouter(handler *handler.Handler, tm *token.Manager) *chi.Mux {
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	router.Get("/health", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Server is up and running"))
	})
	router.Post("/register", handler.RegisterUser)
	router.Post("/login", handler.LoginUser)
	return router
}
