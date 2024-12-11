package handlers

import (
	"net/http"
)

type HealthHandler struct {
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
func (handler *HealthHandler) CheckHealthHandler(w http.ResponseWriter, _ *http.Request) {
	writeJsonData(w, http.StatusOK, HealthEnvelope{
		Title:       handler.HealthEnvelope.Title,
		Description: handler.HealthEnvelope.Description,
		Version:     handler.HealthEnvelope.Version,
		Url:         handler.HealthEnvelope.Url,
	})
}
