package handler

import (
	"net/http"

	"github.com/joshey40/database_aethervault/internal/services"
)

type SyncHandler struct {
	sv *services.SyncService
}

func NewSyncHandler(sv *services.SyncService) *SyncHandler {
	return &SyncHandler{
		sv: sv,
	}
}

func (h *SyncHandler) UserPush(w http.ResponseWriter, r *http.Request) {

}

func (h *SyncHandler) UserPull(w http.ResponseWriter, r *http.Request) {

}
