package handlers

import (
	"encoding/json"
	"net/http"

	"task-crud-kategori/services"
)

type ReportHandler struct {
	service *services.ReportService
}

func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

// GET /api/report/hari-ini
// GET /api/report?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD
func (h *ReportHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	result, err := h.service.GetSummary(startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
