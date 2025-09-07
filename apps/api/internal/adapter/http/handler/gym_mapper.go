// internal/adapter/http/handler/gym_mapper.go
// 役割: Gym Domain ↔ HTTP DTO 変換（Adapter Layer）
// ドメインエンティティとHTTP DTO間の双方向変換を担当
package handler

import (
	"time"

	"gogym-api/internal/adapter/http/dto"
	"gogym-api/internal/domain/gym"
	gymUsecase "gogym-api/internal/usecase/gym"
)

// ToUseCaseSearchRequest はHTTP DTOをユースケースリクエストに変換する
func ToUseCaseSearchRequest(req dto.SearchGymRequest) gymUsecase.SearchGymRequest {
	var location *gym.Location
	if req.Lat != nil && req.Lon != nil {
		location = &gym.Location{
			Latitude:  *req.Lat,
			Longitude: *req.Lon,
		}
	}

	return gymUsecase.SearchGymRequest{
		Query:    req.Query,
		Location: location,
		RadiusM:  req.RadiusM,
		Cursor:   req.Cursor,
		Limit:    req.Limit,
	}
}

// ToUseCaseRecommendRequest はHTTP DTOをユースケースリクエストに変換する
func ToUseCaseRecommendRequest(limit int) gymUsecase.RecommendGymRequest {
	return gymUsecase.RecommendGymRequest{
		UserLocation: nil,
		Limit:        limit,
		Cursor:       "",
	}
}

// ToSearchResponse はユースケースレスポンスをHTTP DTOに変換する
func ToSearchResponse(result *gymUsecase.SearchGymsResponse) dto.SearchGymResponse {
	response := dto.SearchGymResponse{
		Gyms:       make([]dto.GymResponse, len(result.Gyms)),
		NextCursor: result.NextCursor,
		HasMore:    result.HasMore,
	}

	for i, gym := range result.Gyms {
		response.Gyms[i] = ToGymResponse(gym)
	}

	return response
}

// ToRecommendResponse はユースケースレスポンスをHTTP DTOに変換する
func ToRecommendResponse(result *gymUsecase.RecommendGymsResponse) dto.SearchGymResponse {
	response := dto.SearchGymResponse{
		Gyms:       make([]dto.GymResponse, len(result.Gyms)),
		NextCursor: result.NextCursor,
		HasMore:    result.HasMore,
	}

	for i, gym := range result.Gyms {
		response.Gyms[i] = ToGymResponse(gym)
	}

	return response
}

// ToGymResponse はジムドメインエンティティをHTTP DTOに変換する
func ToGymResponse(gym gym.Gym) dto.GymResponse {
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