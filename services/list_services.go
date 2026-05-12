package services

import (
	"github.com/Filbertfelix888/project-management-API/repositories"
	"github.com/Filbertfelix888/project-management-API/models"
	"github.com/Filbertfelix888/project-management-API/utils"
	"github.com/Filbertfelix888/project-management-API/config"
	"github.com/google/uuid"	
	"errors"
	"fmt"
	"gorm.io/gorm"
	"github.com/Filbertfelix888/project-management-API/models/types"	
) 

type listService struct {
	listRepo repositories.ListRepositories
	boardRepo repositories.BoardRepository
	listPosRepo repositories.ListPositionRepository	
}

type ListWithOrder struct {
	positions []uuid.UUID
	Lists []models.List
}

type ListService interface {
	GetByBoardID(boardPublicID string) (*ListWithOrder, error)
	GetByID (id uint) (*models.List, error)
	GetByPublicID(publicID string) (*models.List, error)
	Create(list *models.List) error
	Update(list *models.List) error
	Delete(id uint) error
	UpdatePositions(boardPublicID string, positions []uuid.UUID) error
}

func NewListService(listRepo repositories.ListRepositories, boardRepo repositories.BoardRepository, 
	listPosRepo repositories.ListPositionRepository) ListService {
		return &listService{listRepo, boardRepo, listPosRepo}
	}

func (s *listService) GetByBoardID(boardPublicID string) (*ListWithOrder, error) {
	//Verifikasi board 

	_, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		return nil,errors.New("board not found")
	}

	position, err := s.listPosRepo.GetListOrder(boardPublicID)
	if err != nil {
		return nil, errors.New("failed to get list order : " + err.Error())
	}
	if len (position) == 0 {
		return nil, errors.New("list position not found")
	}

	lists, err := s.listRepo.FindByBoardID(boardPublicID)
	if err != nil {
		return nil, errors.New("failed to get list : " + err.Error())
	}
	
	fmt.Println(position)
	fmt.Println(lists)

	//sorting by position
	orderedLists := utils.SortListsByPosition(lists, position)

	return &ListWithOrder{
		positions: position,
		Lists: orderedLists,
	}, nil 
}

func (s *listService) GetByID (id uint) (*models.List, error) {
	return s.listRepo.FindByID(id)
}

func (s *listService) GetByPublicID(publicID string) (*models.List, error) {
	return s.listRepo.FindByPublicID(publicID)	
}

func (s *listService) Create(list *models.List) error {
	// validasi board, transaction, update position, commit trx
	board, err := s.boardRepo.FindByPublicID(list.BoardPublicID.String())
	if err != nil {
		if  errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("board not found")
		}
		return fmt.Errorf("failed to find board: %w", err)
	}
	list.BoardInternalID = board.InternalID
	
	if list.PublicID == uuid.Nil {
		list.PublicID = uuid.New()
	}

	//mulai transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	//update posisi list
	var position models.ListPosition
	res := tx.Where("board_internal_id = ?", board.InternalID).First(&position)

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
        list.Position = 0
    } else {
        // Posisinya adalah jumlah list yang sudah ada di array ListOrder
		list.Position = len(position.ListOrder)
    }

	//Simpan list baru
	if err := tx.Create(list).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create list: %w", err)
	}
	
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		//buat baru jika belum ada
		position = models.ListPosition{
			PublicID: uuid.New(),
			BoardID: board.InternalID,
			ListOrder: types.UUIDArray{list.PublicID},
		}
		if err := tx.Create(&position).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create list position: %w", err)
		}
	} else if res.Error != nil {
		tx.Rollback()
			return fmt.Errorf("failed to create list position: %w", err)
	} else {
		//Tambahkan ID baru 
		position.ListOrder = append(position.ListOrder, list.PublicID)	
		//Update ke DB	
		if err := tx.Model(&position).Update("list_order", position.ListOrder).Error; err != nil {
			tx.Rollback()
				return fmt.Errorf("failed to update list position: %w", err)
		}
	}

	//commit trx
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("transaction commit failed: %w", err)
	}
	return nil
}

func (s *listService) Update(list *models.List) error {
	return s.listRepo.Update(list)
}

func (s *listService) Delete(id uint) error {
	return s.listRepo.Delete(id)
}

func (s *listService) UpdatePositions(boardPublicID string, positions []uuid.UUID) error {
	// verifikasi board
	board, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		return errors.New("board not found")
	}

	// get list position
	position, err := s.listPosRepo.GetByBoard(board.PublicID.String())
	if err != nil {
		return errors.New("list position not found")
	}

	//update list ordernya	
	position.ListOrder = positions
	return s.listPosRepo.UpdateListOrder(position)
}