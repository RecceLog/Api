package controllers

import (
	"Api/internal/application/services"
	"Api/internal/domain"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RoutesController interface {
	CreateRoute(ctx *gin.Context)
	GetRoutes(ctx *gin.Context)
	GetRoutesInRange(ctx *gin.Context)
	GetRouteById(ctx *gin.Context)
	DeleteRoute(ctx *gin.Context)
}

type routesController struct {
	routesService services.RoutesService
}

func NewRoutesController(routesService services.RoutesService) RoutesController {
	return &routesController{
		routesService: routesService,
	}
}

// --- Request / Response types ---

type CreateRouteRequest struct {
	Route   domain.Route   `json:"route"    binding:"required"`
	NoteSet domain.NoteSet `json:"note_set" binding:"required"`
}

type RouteNotesResponse struct {
	Route     domain.Route     `json:"route"`
	NotesSets []domain.NoteSet `json:"note_sets"`
}

// --- Handlers ---

func (c *routesController) CreateRoute(ctx *gin.Context) {

	var body CreateRouteRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		slog.Info("il body ricevuto è", "body", body)
		ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	err := c.routesService.CreateWithNotes(ctx.Request.Context(), &body.Route, &body.NoteSet)
	if err != nil {
		ErrorResponse(ctx, http.StatusInternalServerError, "Error creating route", err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": body.Route.Id})
}

func (c *routesController) GetRoutes(ctx *gin.Context) {

	routes, err := c.routesService.FindAll(ctx.Request.Context())
	if err != nil {
		ErrorResponse(ctx, http.StatusInternalServerError, "Error retrieving routes", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"routes": routes})
}

func (c *routesController) GetRoutesInRange(ctx *gin.Context) {

	rangeM, err := strconv.Atoi(ctx.Param("range"))
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, "Invalid range parameter", err)
		return
	}

	lat, err := strconv.ParseFloat(ctx.GetHeader("Latitude"), 64)
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, "Invalid Latitude header", err)
		return
	}

	lng, err := strconv.ParseFloat(ctx.GetHeader("Longitude"), 64)
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, "Invalid Longitude header", err)
		return
	}

	routes, err := c.routesService.FindInRange(ctx.Request.Context(), domain.Coordinate{Lat: lat, Lng: lng}, rangeM*1000 /* to km */)
	if err != nil {
		ErrorResponse(ctx, http.StatusInternalServerError, "Error retrieving routes in range", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"routes": routes})
}

func (c *routesController) GetRouteById(ctx *gin.Context) {

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, "Invalid UUID format", err)
		return
	}

	route, err := c.routesService.FindById(ctx.Request.Context(), id)
	if err != nil {
		ErrorResponse(ctx, http.StatusInternalServerError, "Error retrieving route", err)
		return
	}

	noteSets, err := c.routesService.FindNoteSets(ctx.Request.Context(), route.Id)
	if err != nil {
		ErrorResponse(ctx, http.StatusInternalServerError, "Error retrieving note sets", err)
		return
	}

	ctx.JSON(http.StatusOK, RouteNotesResponse{Route: route, NotesSets: noteSets})
}

func (c *routesController) DeleteRoute(ctx *gin.Context) {

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, "Invalid UUID format", err)
		return
	}

	if err = c.routesService.Delete(ctx.Request.Context(), id); err != nil {
		ErrorResponse(ctx, http.StatusInternalServerError, "Error deleting route", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": id})
}
