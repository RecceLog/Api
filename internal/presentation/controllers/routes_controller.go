package controllers

import (
	services2 "Api/internal/application/services"
	"Api/internal/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type routesController struct {
	routesService services2.RoutesService
	notesService  services2.NotesService
}

func NewRoutesController(_routesService services2.RoutesService, _notesService services2.NotesService) *routesController {
	return &routesController{
		routesService: _routesService,
		notesService:  _notesService,
	}
}

type CreateRouteRequest struct {
	Route   domain.Route   `json:"route" binding:"required"`
	NoteSet domain.NoteSet `json:"note_set" binding:"required"`
}

type RouteNotesResponse struct {
	Route     domain.Route     `json:"route"`
	NotesSets []domain.NoteSet `json:"note_sets"`
}

func (c *routesController) CreateRoute(ctx *gin.Context) {

	// Get route and notes from request body
	var routeWithNotes CreateRouteRequest
	if err := ctx.ShouldBindJSON(&routeWithNotes); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, "Invalid body format", err)
		return
	}

	// Insert route in db
	err := c.routesService.Insert(ctx.Request.Context(), &routeWithNotes.Route)
	if err != nil {
		ErrorResponse(ctx, http.StatusInternalServerError, "Error inserting route in database", err)
		return
	}

	// Insert notes in db
	err = c.notesService.InsertSet(ctx.Request.Context(), &routeWithNotes.NoteSet, routeWithNotes.Route.Id)
	if err != nil {
		ErrorResponse(ctx, http.StatusInternalServerError, "Error inserting notes in database", err)
		return
	}

	// Return response
	ctx.JSON(http.StatusCreated, gin.H{
		"id": routeWithNotes.Route.Id,
	})
}

func (c *routesController) GetRoutesInRange(ctx *gin.Context) {

	// Get range from url
	_range, err := strconv.Atoi(ctx.Param("range"))
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, "Invalid range", err)
		return
	}

	// Get position coordinates from headers
	latHeader := ctx.GetHeader("Latitude")
	lngHeader := ctx.GetHeader("Longitude")
	lat, err := strconv.ParseFloat(latHeader, 64)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid latitude value",
		})
		return
	}

	lng, err := strconv.ParseFloat(lngHeader, 64)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid longitude value",
		})
		return
	}

	// Get routes in range
	routes, err := c.routesService.FindInRange(ctx.Request.Context(), domain.Coordinate{Lat: lat, Lng: lng}, _range)
	if err != nil {
		ErrorResponse(ctx, http.StatusInternalServerError, "Error retrieving routes from database", err)
		return
	}

	// Return response
	ctx.JSON(http.StatusOK, gin.H{
		"routes": routes,
	})
}

func (c *routesController) GetRouteById(ctx *gin.Context) {

	// Get id url parameter
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, "Invalid UUID format", err)
		return
	}

	// Retrieve route data
	route, err := c.routesService.FindById(ctx.Request.Context(), id)
	if err != nil {
		ErrorResponse(ctx, http.StatusInternalServerError, "Error retrieving route from database", err)
		return
	}

	// Retrieve note sets about route
	noteSets, err := c.notesService.FindSetsByRouteId(ctx.Request.Context(), route.Id)
	if err != nil {
		ErrorResponse(ctx, http.StatusInternalServerError, "Error retrieving notes from database", err)
		return
	}

	// Write response body
	ctx.JSON(http.StatusOK, RouteNotesResponse{
		Route:     route,
		NotesSets: noteSets,
	})
}

func (c *routesController) GetNoteSetOfRoute(ctx *gin.Context) {

	// Check route id and get note set id from url parameter
	_, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, "Invalid UUID format for id", err)
		return
	}

	setId, err := uuid.Parse(ctx.Param("set"))
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, "Invalid UUID format for set", err)
		return
	}

	// Retrieve data from db
	notes, err := c.notesService.FindSetById(ctx.Request.Context(), setId)
	if err != nil {
		ErrorResponse(ctx, http.StatusInternalServerError, "Error retrieving notes from db", err)
		return
	}

	// Return response
	ctx.JSON(http.StatusOK, gin.H{
		"notes": notes,
	})
}

func (c *routesController) GetRoutes(ctx *gin.Context) {

	routes, err := c.routesService.FindAll(ctx.Request.Context())
	if err != nil {
		ErrorResponse(ctx, http.StatusInternalServerError, "Error retrieving routes from database", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"routes": routes,
	})
}

func (c *routesController) AddNoteSetToRoute(ctx *gin.Context) {

	// Get route id form url parameter
	routeId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, "Invalid UUID format", err)
		return
	}

	// Get note set from request body
	var noteSet domain.NoteSet
	if err = ctx.ShouldBindJSON(&noteSet); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Insert note set in db
	err = c.notesService.InsertSet(ctx.Request.Context(), &noteSet, routeId)
	if err != nil {
		ErrorResponse(ctx, http.StatusInternalServerError, "Error inserting route in database", err)
		return
	}

	// Return response
	ctx.JSON(http.StatusCreated, gin.H{
		"id": noteSet.Id,
	})
}

func (c *routesController) DeleteRoute(ctx *gin.Context) {

	// Get route id from url parameter
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, "Invalid UUID format", err)
		return
	}

	// Delete route (and notes) from db
	err = c.routesService.Delete(ctx, id)
	if err != nil {
		ErrorResponse(ctx, http.StatusInternalServerError, "Error deleting route", err)
		return
	}

	// Return response
	ctx.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func (c *routesController) DeleteNoteSetFromRoute(ctx *gin.Context) {

	// Check route id and get note set id from url parameter
	_, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, "Invalid UUID format for id", err)
		return
	}

	setId, err := uuid.Parse(ctx.Param("set"))
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, "Invalid UUID format for set", err)
		return
	}

	// Remove note set from db
	err = c.notesService.DeleteSet(ctx.Request.Context(), setId)
	if err != nil {
		ErrorResponse(ctx, http.StatusInternalServerError, "Error removing note set", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": setId,
	})
}
