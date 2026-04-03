// Package http contains HTTP handlers for the API server.
package http

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/TranTheTuan/vna/internal/dto"
	"github.com/TranTheTuan/vna/internal/service"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	svc service.UserService
}

// NewAuthHandler creates an AuthHandler with the given UserService.
func NewAuthHandler(svc service.UserService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Register handles POST /api/v1/auth/register.
//
// @Summary      Register a new user
// @Description  Creates a new user account with email and password.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.RegisterRequest   true  "Registration request"
// @Success      201   {object}  dto.RegisterResponse
// @Failure      400   {object}  map[string]string
// @Failure      409   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /api/v1/auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	user, err := h.svc.Register(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidEmail):
			return echo.NewHTTPError(http.StatusBadRequest, "invalid email format")
		case errors.Is(err, service.ErrPasswordTooShort):
			return echo.NewHTTPError(http.StatusBadRequest, "password must be at least 8 characters")
		case errors.Is(err, service.ErrDuplicateEmail):
			return echo.NewHTTPError(http.StatusConflict, "email already registered")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "registration failed")
		}
	}

	return c.JSON(http.StatusCreated, dto.RegisterResponse{
		UserID: user.ID,
		Email:  user.Email,
	})
}

// Login handles POST /api/v1/auth/login.
//
// @Summary      Login
// @Description  Authenticates a user and returns access and refresh tokens.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.LoginRequest  true  "Login request"
// @Success      200   {object}  dto.AuthResponse
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /api/v1/auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	accessToken, refreshToken, err := h.svc.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid email or password")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "login failed")
	}

	return c.JSON(http.StatusOK, dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900, // 15 minutes in seconds
	})
}

// Refresh handles POST /api/v1/auth/refresh.
//
// @Summary      Refresh access token
// @Description  Issues a new access token using a valid refresh token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.RefreshRequest  true  "Refresh token request"
// @Success      200   {object}  dto.RefreshResponse
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /api/v1/auth/refresh [post]
func (h *AuthHandler) Refresh(c echo.Context) error {
	var req dto.RefreshRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if req.RefreshToken == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "refresh_token is required")
	}

	accessToken, err := h.svc.RefreshToken(c.Request().Context(), req.RefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrTokenInvalid),
			errors.Is(err, service.ErrTokenRevoked),
			errors.Is(err, service.ErrTokenExpired):
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired refresh token")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "token refresh failed")
		}
	}

	return c.JSON(http.StatusOK, dto.RefreshResponse{
		AccessToken: accessToken,
		ExpiresIn:   900,
	})
}

// Logout handles POST /api/v1/auth/logout.
// Requires a valid access token (enforced by JWT middleware on the route).
//
// @Summary      Logout
// @Description  Revokes the provided refresh token. Requires a valid Bearer access token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.LogoutRequest  true  "Logout request"
// @Success      200   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c echo.Context) error {
	var req dto.LogoutRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if req.RefreshToken == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "refresh_token is required")
	}

	if err := h.svc.Logout(c.Request().Context(), req.RefreshToken); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "logout failed")
	}

	return c.JSON(http.StatusOK, map[string]string{})
}
