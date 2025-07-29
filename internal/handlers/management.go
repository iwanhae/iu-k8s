package handlers

import (
	"context"
	"time"

	"iu-k8s.linecorp.com/server/internal/api"
)

type ManagementHandler struct{}

func (h *ManagementHandler) GetReadiness(ctx context.Context, request api.GetReadinessRequestObject) (api.GetReadinessResponseObject, error) {
	return api.GetReadiness200JSONResponse{
		Status:    api.Ready,
		Timestamp: time.Now(),
		Version:   "1.0.0",
	}, nil
}
