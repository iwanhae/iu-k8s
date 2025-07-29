package handlers

import (
	"context"
	"time"

	"iu-k8s.linecorp.com/server/internal/api"
	"iu-k8s.linecorp.com/server/internal/log"
)

type ManagementHandler struct{}

// GetReadiness checks if the service is ready
// (GET /readyz)
func (h *ManagementHandler) GetReadiness(ctx context.Context, request api.GetReadinessRequestObject) (api.GetReadinessResponseObject, error) {
	return api.GetReadiness200JSONResponse{
		Status:    api.Ready,
		Timestamp: time.Now(),
		Version:   "v0.0.1", // TODO: build time injection
	}, nil
}

// SetLogLevel sets the log level dynamically
// (GET /debug/log)
func (h *ManagementHandler) SetLogLevel(ctx context.Context, request api.SetLogLevelRequestObject) (api.SetLogLevelResponseObject, error) {
	if request.Params.Level != nil {
		if err := log.SetLevel(string(*request.Params.Level)); err != nil {
			return api.SetLogLevel400JSONResponse{
				Error:   "invalid_log_level",
				Message: err.Error(),
			}, nil
		}
	}

	if request.Params.Format != nil {
		if err := log.SetFormat(string(*request.Params.Format)); err != nil {
			return api.SetLogLevel400JSONResponse{
				Error:   "invalid_log_format",
				Message: err.Error(),
			}, nil
		}
	}

	level := log.GetLevel()
	format := log.GetFormat()
	return api.SetLogLevel200JSONResponse{
		Level:  &level,
		Format: &format,
	}, nil
}
