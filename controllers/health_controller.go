package controllers

import (
	"net/http"
)

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (c *HealthController) Routes() []Route {
	return []Route{
		{
			Path:    "/health",
			Method:  http.MethodGet,
			Handler: c.Check,
		},
	}
}

func (c *HealthController) Check(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}