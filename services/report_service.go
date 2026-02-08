package services

import "task-crud-kategori/repositories"

type ReportService struct {
	repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) GetSummary(startDate, endDate string) (interface{}, error) {
	return s.repo.GetSummary(startDate, endDate)
}
