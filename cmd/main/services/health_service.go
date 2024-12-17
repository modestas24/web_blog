package services

import (
	"net/http"
	"web_blog/cmd/main/utils"
)

type HealthService struct {
	HealthEnvelope HealthEnvelope
}

type HealthEnvelope struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Url         string `json:"url"`
}

// CheckHealth godoc
//
//	@Summary		Show server health information.
//	@Description	get server health object
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	EnvelopeJson{data=HealthEnvelope}
//	@Router			/health [get]
func (service *HealthService) CheckHealth(w http.ResponseWriter, _ *http.Request) {
	utils.WriteJsonData(w, http.StatusOK, service.HealthEnvelope)
}
