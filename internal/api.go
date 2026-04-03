package internal

import (
	"Api/internal/application/services"
	"Api/internal/infrastructure/repositories"
	"cmp"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"Api/internal/presentation/controllers"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Application struct {
	Addr  string
	Db    *pgxpool.Pool
	Cache *redis.Client
}

func (app Application) Mount() http.Handler {

	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	routesRepo := repositories.NewRoutesRepository()
	notesRepo := repositories.NewNotesRepository()

	routesService := services.NewRoutesService(app.Db, routesRepo, notesRepo)
	//notesService := services.NewNotesService(app.Db)

	routesController := controllers.NewRoutesController(routesService)
	notesController := controllers.NewNotesController(routesService)

	v1 := r.Group("/v1")
	{
		routes := v1.Group("/routes")
		notes := v1.Group("/notes")

		// CREATE
		routes.POST("/", routesController.CreateRoute)
		routes.POST("/:id/notes", notesController.AddNoteSetToRoute)

		// READ
		routes.GET("/", routesController.GetRoutes)
		routes.GET("/:id", routesController.GetRouteById)
		routes.GET("/range/:range", routesController.GetRoutesInRange)
		notes.GET("/:id", notesController.GetNoteSetOfRoute)

		// UPDATE

		// DELETE
		routes.DELETE("/:id", routesController.DeleteRoute)
		notes.DELETE("/:id", notesController.DeleteNoteSetFromRoute)
	}

	return r
}

func (app Application) Run(h http.Handler) error {
	srv := &http.Server{
		Addr:         cmp.Or(os.Getenv("ADDR"), ":8080"),
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	slog.Info(fmt.Sprintf("Starting server, listening at %s", srv.Addr))

	return srv.ListenAndServe()
}
