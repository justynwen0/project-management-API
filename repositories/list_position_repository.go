package repositories

import (
	"errors"

	"github.com/google/uuid"
	"github.com/justynwen0/project-management-API/config"
	"github.com/justynwen0/project-management-API/models"
	"gorm.io/gorm"
)

type listPositionRepository struct {
}

type ListPositionRepository interface {
	GetByBoard(boardPublicID string) (*models.ListPosition, error)
	CreateOrUpdate(boardPublicID string, listOrder []uuid.UUID) error
	GetListOrder(boardPublicID string) ([]uuid.UUID, error)
	UpdateListOrder(position *models.ListPosition) error
}

func NewListPositionRepository() ListPositionRepository {
	return &listPositionRepository{}
}

func (r *listPositionRepository) GetByBoard(boardPublicID string) (*models.ListPosition, error) {
	var position models.ListPosition

	err := config.DB.
		Joins("JOIN boards ON boards.internal_id = list_positions.board_internal_id").
		Where("boards.public_id = ?", boardPublicID).
		First(&position).Error

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &position, nil
}

func (r *listPositionRepository) CreateOrUpdate(boardPublicID string, listOrder []uuid.UUID) error {
	return config.DB.Exec(`
	INSERT INTO list_positions (board_internal_id, list_order)
	SELECT internal_id, ? FROM boards Where public_id = ?
	ON CONFLICT (board_internal_id)
	DO UPDATE SET list_order = EXCLUDED.list_order
	`, listOrder, boardPublicID).Error
}

func (r *listPositionRepository) GetListOrder(boardPublicID string) ([]uuid.UUID, error) {
	position, err := r.GetByBoard(boardPublicID)

	if err != nil {
		return nil, err
	}

	if position == nil {
		return []uuid.UUID{}, nil
	}

	return position.ListOrder, nil
}

func (r *listPositionRepository) UpdateListOrder(position *models.ListPosition) error {
	return config.DB.Model(&models.ListPosition{}).
		Where("internal_id = ?", position.InternalID).
		Update("list_order", position.ListOrder).Error
}
