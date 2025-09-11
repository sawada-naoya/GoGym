package handler

import (
	"net/http"
	"strconv"

	"gogym-api/internal/adapter/http/dto"
	"gogym-api/internal/domain/gym"
	gymUsecase "gogym-api/internal/usecase/gym"

	"github.com/labstack/echo/v4"
)

// GymHandler handles gym-related HTTP requests
type GymHandler struct {
	gu *gymUsecase.UseCase
}

// NewGymHandler creates a new GymHandler instance
func NewGymHandler(gu *gymUsecase.UseCase) *GymHandler {
	return &GymHandler{
		gu: gu,
	}
}

// SearchGyms handles GET /gyms request
func (h *GymHandler) SearchGyms(c echo.Context) error {
	var req dto.SearchGymRequest

	// クエリパラメータをバインド
	req.Query = c.QueryParam("q")
	req.Cursor = c.QueryParam("cursor")

	if lat := c.QueryParam("lat"); lat != "" {
		if latVal, err := strconv.ParseFloat(lat, 64); err == nil {
			req.Lat = &latVal
		}
	}

	if lon := c.QueryParam("lon"); lon != "" {
		if lonVal, err := strconv.ParseFloat(lon, 64); err == nil {
			req.Lon = &lonVal
		}
	}

	if radius := c.QueryParam("radius_m"); radius != "" {
		if radiusVal, err := strconv.Atoi(radius); err == nil {
			req.RadiusM = &radiusVal
		}
	}

	if limit := c.QueryParam("limit"); limit != "" {
		if limitVal, err := strconv.Atoi(limit); err == nil {
			req.Limit = limitVal
		}
	}

	// バリデーション
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":   "invalid_request",
			"message": err.Error(),
		})
	}

	// ユースケースリクエストに変換
	usecaseReq := ToUseCaseSearchRequest(req)

	// ジム検索実行
	result, err := h.gu.SearchGyms(c.Request().Context(), usecaseReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "search_failed",
			"message": "Failed to search gyms",
		})
	}

	// レスポンシに変換
	response := ToSearchResponse(result)

	return c.JSON(http.StatusOK, response)
}

// GetGym handles GET /gyms/:id request
func (h *GymHandler) GetGym(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":   "invalid_gym_id",
			"message": "Invalid gym ID format",
		})
	}

	gym, err := h.gu.GetGym(c.Request().Context(), gym.ID(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error":   "gym_not_found",
			"message": "Gym not found",
		})
	}

	response := ToGymResponse(*gym)
	return c.JSON(http.StatusOK, response)
}

// GetRecommendedGyms handles GET /gyms/recommended request
func (h *GymHandler) GetRecommendedGyms(c echo.Context) error {
	ctx := c.Request().Context()
	limitParam := c.QueryParam("limit")
	limit := 10
	if limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 && parsedLimit <= 50 {
			limit = parsedLimit
		}
	}

	// おすすめジムのユースケースリクエストを作成
	recommendReq := ToUseCaseRecommendRequest(limit)

	// おすすめジム取得実行
	result, err := h.gu.RecommendGyms(ctx, recommendReq)
	if err != nil {
		println("Handler error:", err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// レスポンシに変換
	response := ToRecommendResponse(result)

	// UTF-8の文字エンコーディングを明示的に設定
	c.Response().Header().Set("Content-Type", "application/json; charset=utf-8")
	return c.JSON(http.StatusOK, response)
}

// CreateGym handles POST /gyms request  
func (h *GymHandler) CreateGym(c echo.Context) error {
	// TODO: Create gym
	return c.JSON(http.StatusCreated, map[string]string{
		"message": "Gym created",
	})
}

// UpdateGym handles PUT /gyms/:id request
func (h *GymHandler) UpdateGym(c echo.Context) error {
	// TODO: Update gym
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Gym updated", 
	})
}

// DeleteGym handles DELETE /gyms/:id request
func (h *GymHandler) DeleteGym(c echo.Context) error {
	// TODO: Delete gym
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Gym deleted",
	})
}

// GetGymImages handles GET /gyms/:id/images request
func (h *GymHandler) GetGymImages(c echo.Context) error {
	// TODO: Get gym images from reviews
	return c.JSON(http.StatusOK, map[string]interface{}{
		"images": []interface{}{},
	})
}

// AutocompleteGyms handles GET /gyms/autocomplete request
func (h *GymHandler) AutocompleteGyms(c echo.Context) error {
	// TODO: Gym name autocomplete
	return c.JSON(http.StatusOK, map[string]interface{}{
		"suggestions": []interface{}{},
	})
}
