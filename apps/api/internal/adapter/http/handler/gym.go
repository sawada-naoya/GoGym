package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"gogym-api/internal/adapter/http/dto"
	"gogym-api/internal/domain/common"
	"gogym-api/internal/domain/gym"
	gymUsecase "gogym-api/internal/usecase/gym"
)

// GymHandler handles gym-related HTTP requests
type GymHandler struct {
	gymUsecase *gymUsecase.UseCase
}

// NewGymHandler creates a new GymHandler instance
func NewGymHandler(gymUsecase *gymUsecase.UseCase) *GymHandler {
	return &GymHandler{
		gymUsecase: gymUsecase,
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
			"error": "invalid_request",
			"message": err.Error(),
		})
	}
	
	// ユースケースリクエストに変換
	var location *common.Location
	if req.Lat != nil && req.Lon != nil {
		location = &common.Location{
			Latitude:  *req.Lat,
			Longitude: *req.Lon,
		}
	}
	
	usecaseReq := gymUsecase.SearchGymRequest{
		Query:    req.Query,
		Location: location,
		RadiusM:  req.RadiusM,
		Cursor:   req.Cursor,
		Limit:    req.Limit,
	}
	
	// ジム検索実行
	result, err := h.gymUsecase.SearchGyms(c.Request().Context(), usecaseReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "search_failed",
			"message": "Failed to search gyms",
		})
	}
	
	// レスポンスに変換
	response := dto.SearchGymResponse{
		Gyms:       make([]dto.GymResponse, len(result.Gyms)),
		NextCursor: result.NextCursor,
		HasMore:    result.HasMore,
	}
	
	for i, gym := range result.Gyms {
		response.Gyms[i] = h.convertToGymResponse(gym)
	}
	
	return c.JSON(http.StatusOK, response)
}

// GetGym handles GET /gyms/:id request
func (h *GymHandler) GetGym(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid_gym_id",
			"message": "Invalid gym ID format",
		})
	}
	
	gym, err := h.gymUsecase.GetGym(c.Request().Context(), common.ID(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "gym_not_found",
			"message": "Gym not found",
		})
	}
	
	response := h.convertToGymResponse(*gym)
	return c.JSON(http.StatusOK, response)
}

// GetRecommendedGyms handles GET /gyms/recommended request
func (h *GymHandler) GetRecommendedGyms(c echo.Context) error {
	limitParam := c.QueryParam("limit")
	limit := 10 // デフォルト値
	if limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 && parsedLimit <= 50 {
			limit = parsedLimit
		}
	}

	// おすすめジムのユースケースリクエストを作成
	recommendReq := gymUsecase.RecommendGymRequest{
		UserLocation: nil, // 今回は位置情報なし
		Limit:        limit,
		Cursor:       "",
	}

	// おすすめジム取得実行
	result, err := h.gymUsecase.RecommendGyms(c.Request().Context(), recommendReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "recommend_failed",
			"message": "Failed to fetch recommended gyms",
		})
	}

	// レスポンスに変換
	response := dto.SearchGymResponse{
		Gyms:       make([]dto.GymResponse, len(result.Gyms)),
		NextCursor: result.NextCursor,
		HasMore:    result.HasMore,
	}

	for i, gym := range result.Gyms {
		response.Gyms[i] = h.convertToGymResponse(gym)
	}

	return c.JSON(http.StatusOK, response)
}

// convertToGymResponse converts gym domain model to response DTO
func (h *GymHandler) convertToGymResponse(gym gym.Gym) dto.GymResponse {
	tags := make([]dto.TagResponse, len(gym.Tags))
	for i, tag := range gym.Tags {
		tags[i] = dto.TagResponse{
			ID:   int64(tag.ID),
			Name: tag.Name,
		}
	}
	
	return dto.GymResponse{
		ID:            int64(gym.ID),
		Name:          gym.Name,
		Description:   gym.Description,
		Location:      gym.Location,
		Address:       gym.Address,
		City:          gym.City,
		Prefecture:    gym.Prefecture,
		PostalCode:    gym.PostalCode,
		Tags:          tags,
		AverageRating: gym.AverageRating,
		ReviewCount:   gym.ReviewCount,
		CreatedAt:     gym.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     gym.UpdatedAt.Format(time.RFC3339),
	}
}

