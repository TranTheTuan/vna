package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/TranTheTuan/vna/internal/dto"
	"github.com/TranTheTuan/vna/internal/service"
)

// MessageHandler handles chat message endpoints.
type MessageHandler struct {
	svc service.MessageService
}

// NewMessageHandler creates a MessageHandler with the given MessageService.
func NewMessageHandler(svc service.MessageService) *MessageHandler {
	return &MessageHandler{svc: svc}
}

// Send handles POST /api/v1/messages.
// Requires JWT auth — user_id is extracted from context set by JWTMiddleware.
//
// @Summary      Send a chat message
// @Description  Sends a message to the AI and returns the response. Requires a valid Bearer access token.
// @Tags         messages
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.SendMessageRequest  true  "Message request"
// @Success      201   {object}  dto.MessageResponse
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Failure      504   {object}  map[string]string
// @Router       /api/v1/messages [post]
func (h *MessageHandler) Send(c echo.Context) error {
	var req dto.SendMessageRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if req.Message == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "message is required")
	}

	userID := c.Get("user_id").(string)

	msg, err := h.svc.Send(c.Request().Context(), userID, req.Message)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmptyMessage):
			return echo.NewHTTPError(http.StatusBadRequest, "message is required")
		case errors.Is(err, service.ErrUpstreamTimeout):
			return echo.NewHTTPError(http.StatusGatewayTimeout, "upstream API timed out")
		case errors.Is(err, service.ErrUpstreamFailed):
			return echo.NewHTTPError(http.StatusBadGateway, "upstream API error")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to process message")
		}
	}

	return c.JSON(http.StatusCreated, dto.MessageResponse{
		ID:        msg.ID,
		Question:  msg.Question,
		Answer:    msg.Answer,
		CreatedAt: msg.CreatedAt,
	})
}

// List handles GET /api/v1/messages.
// Supports cursor-based pagination via ?limit=N&cursor=<uuid>.
//
// @Summary      List chat messages
// @Description  Returns a paginated list of the authenticated user's messages. Requires a valid Bearer access token.
// @Tags         messages
// @Produce      json
// @Security     BearerAuth
// @Param        limit   query     int     false  "Number of results (1-100, default 20)"
// @Param        cursor  query     string  false  "Pagination cursor from previous response"
// @Success      200     {object}  dto.ListResponse
// @Failure      400     {object}  map[string]string
// @Failure      401     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /api/v1/messages [get]
func (h *MessageHandler) List(c echo.Context) error {
	userID := c.Get("user_id").(string)

	limit := 20
	if raw := c.QueryParam("limit"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 || n > 100 {
			return echo.NewHTTPError(http.StatusBadRequest, "limit must be an integer between 1 and 100")
		}
		limit = n
	}

	cursor := c.QueryParam("cursor")

	msgs, nextCursor, err := h.svc.List(c.Request().Context(), userID, limit, cursor)
	if err != nil {
		if errors.Is(err, service.ErrInvalidLimit) {
			return echo.NewHTTPError(http.StatusBadRequest, "limit must be between 1 and 100")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list messages")
	}

	// Convert domain slice to response DTOs
	data := make([]*dto.MessageResponse, len(msgs))
	for i, m := range msgs {
		data[i] = &dto.MessageResponse{
			ID:        m.ID,
			Question:  m.Question,
			Answer:    m.Answer,
			CreatedAt: m.CreatedAt,
		}
	}

	return c.JSON(http.StatusOK, dto.ListResponse{
		Data:       data,
		NextCursor: nextCursor,
	})
}
