package services

import (
	"github.com/justynwen0/project-management-API/models"
	"github.com/justynwen0/project-management-API/repositories"
)

type DashboardService interface {
	GetTaskCountByAssignee() ([]models.AssigneeTaskCount, error)
	GetTaskPercentage() ([]models.TaskPercentage, error)
}

type dashboardService struct {
	cardRepo repositories.CardRepository
}

func NewDashboardService(cardRepo repositories.CardRepository) DashboardService {
	return &dashboardService{
		cardRepo: cardRepo,
	}
}

func (s *dashboardService) GetTaskCountByAssignee() ([]models.AssigneeTaskCount, error) {
	return s.cardRepo.GetTaskCountByAssignee()
}

func (s *dashboardService) GetTaskPercentage() ([]models.TaskPercentage, error) {
	return s.cardRepo.GetTaskPercentage()
}
