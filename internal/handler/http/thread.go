package http

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/TranTheTuan/vna/internal/dto"
	"github.com/TranTheTuan/vna/internal/service"
)

// ThreadHandler handles chat thread endpoints.
type ThreadHandler struct {
	svc service.ThreadService
}

// NewThreadHandler creates a ThreadHandler with the given ThreadService.
func NewThreadHandler(svc service.ThreadService) *ThreadHandler {
	return &ThreadHandler{svc: svc}
}

// List handles GET /api/v1/threads.
// Returns all threads for the authenticated user, ordered newest-first.
//
// @Summary      List chat threads
// @Description  Returns all chat threads for the authenticated user. Requires a valid Bearer access token.
// @Tags         threads
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  dto.ListThreadsResponse
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/threads [get]
func (h *ThreadHandler) List(c echo.Context) error {
	userID := c.Get("user_id").(string)

	threads, err := h.svc.ListByUser(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list threads")
	}

	data := make([]*dto.ThreadResponse, len(threads))
	for i, t := range threads {
		data[i] = &dto.ThreadResponse{
			ID:        t.ID,
			Name:      t.Name,
			CreatedAt: t.CreatedAt,
		}
	}

	return c.JSON(http.StatusOK, dto.ListThreadsResponse{Data: data})
}

// Rename handles PATCH /api/v1/threads/:id.
// Updates the thread name. Validates that the thread belongs to the authenticated user.
//
// @Summary      Rename a chat thread
// @Description  Updates the name of a thread. Requires a valid Bearer access token and thread ownership.
// @Tags         threads
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                  true  "Thread ID"
// @Param        body  body      dto.RenameThreadRequest true  "New thread name"
// @Success      200   {object}  dto.ThreadResponse
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /api/v1/threads/{id} [put]
func (h *ThreadHandler) Rename(c echo.Context) error {
	userID := c.Get("user_id").(string)
	threadID := c.Param("id")

	var req dto.RenameThreadRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if req.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "name is required")
	}

	t, err := h.svc.Rename(c.Request().Context(), userID, threadID, req.Name)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrThreadNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "thread not found")
		case errors.Is(err, service.ErrInvalidThreadName):
			return echo.NewHTTPError(http.StatusBadRequest, "name is required")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to rename thread")
		}
	}

	return c.JSON(http.StatusOK, dto.ThreadResponse{
		ID:        t.ID,
		Name:      t.Name,
		CreatedAt: t.CreatedAt,
	})
}
