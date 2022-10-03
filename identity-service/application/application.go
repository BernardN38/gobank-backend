package application

import (
	"database/sql"
	"github.com/bernardn38/gobank/identity-service/handler"
	"github.com/bernardn38/gobank/identity-service/sql/users"
	"github.com/bernardn38/gobank/identity-service/token"
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
	defer app.pgDb.Close()
	log.Fatal(http.ListenAndServe(":80", app.srv.router))
}

func (app *App) runAppSetup(config Config) {
	db, err := sql.Open("postgres", "user=bernardn dbname=identity_service sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	queries := users.New(db)
	handler := &handler.Handler{UsersDb: queries}
	tokenManger := token.NewManager([]byte(config.jwtSecretKey), config.jwtSigningMethod)

	handler.TokenManager = tokenManger
	app.srv.router = SetupRouter(handler, tokenManger)
	app.pgDb = db
	app.srv.handler = handler
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
	router.Get("/users/all", handler.GetAllUsers)
	router.Mount("/", setupProtectedRoutes(handler, tm))
	return router
}

func setupProtectedRoutes(handler *handler.Handler, tokenManager *token.Manager) http.Handler {
	r := chi.NewRouter()
	r.Use(handler.VerifyJwtToken)
	r.Get("/users/{userId}", handler.GetUserView)
	return r
}
