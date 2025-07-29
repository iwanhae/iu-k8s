package handlers

import (
	"iu-k8s.linecorp.com/server/internal/api"
)

var _ api.StrictServerInterface = (*aggregated)(nil)

type aggregated struct {
	*ManagementHandler
}

func New() *aggregated {
	return &aggregated{
		ManagementHandler: &ManagementHandler{},
	}
}
