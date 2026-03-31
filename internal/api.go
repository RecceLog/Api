package internal

import (
	services2 "Api/internal/application/services"
	"log/slog"
	"net/http"
	"time"

	"Api/internal/infrastructure/database"
	"Api/internal/presentation/controllers"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	Config ApiConfig
	Db     *pgxpool.Pool
}

func (app Application) Mount() http.Handler {

	r := gin.New()

	r.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	routesService := services2.NewRoutesService(app.Db)
	notesService := services2.NewNotesService(app.Db)

	routesController := controllers.NewRoutesController(routesService, notesService)

	v1 := r.Group("/v1")
	{
		routes := v1.Group("/routes")

		// CREATE
		routes.POST("/", routesController.CreateRoute)
		routes.POST("/:id/notes", routesController.AddNoteSetToRoute)

		// READ
		routes.GET("/", routesController.GetRoutes)
		routes.GET("/:id", routesController.GetRouteById)
		routes.GET("/range/:range", routesController.GetRoutesInRange)
		routes.GET("/:id/note-set/:set", routesController.GetNoteSetOfRoute)

		// UPDATE

		// DELETE
		routes.DELETE("/:id", routesController.DeleteRoute)
		routes.DELETE("/:id/note-set/:set", routesController.DeleteNoteSetFromRoute)
	}

	return r
}

func (app Application) Run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.Config.Address,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	slog.Info("Starting server", "listening at", app.Config.Address)

	return srv.ListenAndServe()
}

type ApiConfig struct {
	Address string
	Db      database.DbConfig
}
