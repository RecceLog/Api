package controllers

import (
	"Api/internal/application/services"
	"Api/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type notesController struct {
	routesService services.RoutesService
}

func NewNotesController(routesService services.RoutesService) *notesController {
	return &notesController{
		routesService: routesService,
	}
}

func (c *notesController) AddNoteSetToRoute(ctx *gin.Context) {

	routeId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, "Invalid UUID format for route id", err)
		return
	}

	var noteSet domain.NoteSet
	if err = ctx.ShouldBindJSON(&noteSet); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	err = c.routesService.AddNoteSet(ctx.Request.Context(), &noteSet, routeId)
	if err != nil {
		ErrorResponse(ctx, http.StatusInternalServerError, "Error adding note set", err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": noteSet.Id})
}

func (c *notesController) GetNoteSetOfRoute(ctx *gin.Context) {

	setId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, "Invalid UUID format for set id", err)
		return
	}

	noteSet, err := c.routesService.FindNoteSet(ctx.Request.Context(), setId)
	if err != nil {
		ErrorResponse(ctx, http.StatusInternalServerError, "Error retrieving note set", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"notes": noteSet})
}

func (c *notesController) DeleteNoteSetFromRoute(ctx *gin.Context) {

	if _, err := uuid.Parse(ctx.Param("id")); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, "Invalid UUID format for route id", err)
		return
	}

	setId, err := uuid.Parse(ctx.Param("set"))
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, "Invalid UUID format for set id", err)
		return
	}

	if err := c.routesService.DeleteNoteSet(ctx.Request.Context(), setId); err != nil {
		ErrorResponse(ctx, http.StatusInternalServerError, "Error deleting note set", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": setId})
}
