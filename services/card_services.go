package services

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/justynwen0/project-management-API/config"
	"github.com/justynwen0/project-management-API/models"
	"github.com/justynwen0/project-management-API/models/types"
	"github.com/justynwen0/project-management-API/repositories"
	"gorm.io/gorm"
)

type CardService interface {
	Create(card *models.Card, listPublicID string) error
	Update(card *models.Card, listPublicID string) error
	Delete(id uint) error

	GetByListID(listPublicID string) ([]models.Card, error)
	GetByID(id uint) (*models.Card, error)
	GetByPublicID(publicID string) (*models.Card, error)

	AddAssignees(cardPublicID string, userPublicIDs []string) error
	RemoveAssignees(cardPublicID string, userPublicIDs []string) error
	AddAssigneesByPublicID(cardPublicID string, userPublicIDs []string) error
}

type cardService struct {
	cardRepo repositories.CardRepository
	listRepo repositories.ListRepositories
	userRepo repositories.UserRepository
}

func NewCardService(
	cardRepo repositories.CardRepository,
	listRepo repositories.ListRepositories,
	userRepo repositories.UserRepository,
) CardService {
	return &cardService{cardRepo, listRepo, userRepo}
}

func (s *cardService) Create(card *models.Card, listPublicID string) error {
	// 1. Ambil list dari listPublicID
	list, err := s.listRepo.FindByPublicID(listPublicID)
	if err != nil {
		return fmt.Errorf("list not found: %w", err)
	}

	// 2. Set list_internal_id
	card.ListID = list.InternalID

	// 3. Generate public_id jika belum ada\
	if card.PublicID == uuid.Nil {
		card.PublicID = uuid.New()
	}
	card.CreatedAt = time.Now()

	// 4. Mulai transaksi
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// 5. Simpan card
	if err := tx.Create(card).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create card: %w", err)
	}

	// 6. Update atau buat card position
	var position models.CardPosition
	if err := tx.Model(&models.CardPosition{}).
		Where("list_internal_id = ?", list.InternalID).
		First(&position).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Buat baru jika belum ada
			position = models.CardPosition{
				PublicID:  uuid.New(),
				ListID:    list.InternalID,
				CardOrder: types.UUIDArray{card.PublicID},
			}

			if err := tx.Create(&position).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to create card position: %w", err)
			}
		} else {
			tx.Rollback()
			return fmt.Errorf("failed to get card position: %w", err)
		}
	} else {
		// Tambahkan card baru ke urutan
		position.CardOrder = append(position.CardOrder, card.PublicID)
		if err := tx.Model(&models.CardPosition{}).
			Where("internal_id = ?", position.InternalID).
			Update("card_order", position.CardOrder).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update card position: %w", err)
		}
	}

	// 7. Commit transaksi
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("transaction commit failed: %w", err)
	}
	return nil
}

func (s *cardService) Update(card *models.Card, listPublicID string) error {
	exitingCard, err := s.cardRepo.FindByPublicID(card.PublicID.String())
	if err != nil {
		return fmt.Errorf("card not found: %w", err)
	}

	newList, err := s.listRepo.FindByPublicID(listPublicID)
	if err != nil {
		return fmt.Errorf("list not found: %w", err)
	}

	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if exitingCard.ListID != newList.InternalID {
		var oldPos models.CardPosition
		err := tx.Where("list_internal_id = ?", exitingCard.ListID).First(&oldPos).Error
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				tx.Rollback()
				return fmt.Errorf("failed to get old card position: %w", err)
			}
		} else {
			filtered := make(types.UUIDArray, 0, len(oldPos.CardOrder))
			for _, id := range oldPos.CardOrder {
				if id != exitingCard.PublicID {
					filtered = append(filtered, id)
				}
			}

			if err := tx.Model(&models.CardPosition{}).
				Where("internal_id = ?", oldPos.InternalID).
				Update("card_order", types.UUIDArray(filtered)).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to update old card position: %w", err)
			}
		}

		var newPos models.CardPosition
		res := tx.Where("list_internal_id = ?", newList.InternalID).First(&newPos)
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			newPos = models.CardPosition{
				PublicID:  uuid.New(),
				ListID:    newList.InternalID,
				CardOrder: types.UUIDArray{exitingCard.PublicID},
			}
			if err := tx.Create(&newPos).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to create card position for new list: %w", err)
			}
		} else if res.Error == nil {
			updateOrder := append(newPos.CardOrder, exitingCard.PublicID)
			if err := tx.Model(&models.CardPosition{}).
				Where("internal_id = ?", newPos.InternalID).
				Update("card_order", types.UUIDArray(updateOrder)).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to update new card position: %w", err)
			}
		} else {
			tx.Rollback()
			return fmt.Errorf("failed to get new card position: %w", res.Error)
		}

		card.ListID = newList.InternalID
	} else {
		card.ListID = exitingCard.ListID
	}

	card.InternalID = exitingCard.InternalID
	card.PublicID = exitingCard.PublicID

	if err := tx.Save(card).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update card: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("transaction commit failed: %w", err)
	}
	return nil
}

func (s *cardService) Delete(id uint) error {
	return s.cardRepo.Delete(id)
}

func (s *cardService) GetByListID(listPublicID string) ([]models.Card, error) {
	//Verifikasi listnya ada
	list, err := s.listRepo.FindByPublicID(listPublicID)
	if err != nil {
		return nil, fmt.Errorf("List tidak ditemukan: %w", err)
	}
	// ambil card positionnya
	position, err := s.cardRepo.FindCardPositionByListID(list.InternalID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("Gagal mengambil posisi card: %w", err)
	}

	if position == nil {
		return []models.Card{}, nil
	}
	// ambil semua card di list tersebut
	cards, err := s.cardRepo.FindByListID(listPublicID)
	if err != nil {
		return nil, fmt.Errorf("Gagal mengambil cards: %w", err)
	}
	// sorting
	if len(position.CardOrder) > 0 {
		cards = sortCardByPosition(cards, position.CardOrder)
	}

	return cards, nil
}

func sortCardByPosition(cards []models.Card, order []uuid.UUID) []models.Card {
	// buat map untuk pencarian cepat
	orderMap := make(map[uuid.UUID]int)
	for i, id := range order {
		orderMap[id] = i
	}
	defaultIndex := len(order)

	// sorting table
	sort.SliceStable(cards, func(i, j int) bool {
		idxI, okI := orderMap[cards[i].PublicID]
		if !okI {
			idxI = defaultIndex
		}

		idxJ, okJ := orderMap[cards[j].PublicID]
		if !okJ {
			idxJ = defaultIndex
		}

		if idxI == idxJ {
			return cards[i].CreatedAt.Before(cards[j].CreatedAt)
		}
		return idxI < idxJ
	})
	return cards
}

func (s *cardService) GetByID(id uint) (*models.Card, error) {
	return s.cardRepo.FindByID(id)
}

func (s *cardService) GetByPublicID(publicID string) (*models.Card, error) {
	return s.cardRepo.FindByPublicID(publicID)
}

func (s *cardService) AddAssignees(cardPublicID string, userPublicIDs []string) error {
	card, err := s.cardRepo.FindByPublicID(cardPublicID)
	if err != nil {
		return fmt.Errorf("card not found: %w", err)
	}

	var userInternalIDs []uint
	for _, userPublicID := range userPublicIDs {
		user, err := s.userRepo.FindByPublicID(userPublicID)
		if err != nil {
			return fmt.Errorf("user not found: %w", err)
		}
		userInternalIDs = append(userInternalIDs, uint(user.InternalID))
	}

	return s.cardRepo.AddAssignees(uint(card.InternalID), userInternalIDs)
}

func (s *cardService) RemoveAssignees(cardPublicID string, userPublicIDs []string) error {
	card, err := s.cardRepo.FindByPublicID(cardPublicID)
	if err != nil {
		return fmt.Errorf("card not found: %w", err)
	}

	var userInternalIDs []uint
	for _, userPublicID := range userPublicIDs {
		user, err := s.userRepo.FindByPublicID(userPublicID)
		if err != nil {
			return fmt.Errorf("user not found: %w", err)
		}
		userInternalIDs = append(userInternalIDs, uint(user.InternalID))
	}

	return s.cardRepo.RemoveAssignees(uint(card.InternalID), userInternalIDs)
}

func (s *cardService) AddAssigneesByPublicID(cardPublicID string, userPublicIDs []string) error {
	// PROTEKSI: Cek input agar tidak kosong
	if cardPublicID == "" || cardPublicID == "undefined" {
		return fmt.Errorf("card ID tidak valid")
	}

	if len(userPublicIDs) == 0 {
		return nil
	}

	// DISINI SOLUSINYA:
	// Langsung panggil fungsi Repository yang baru.
	// Kita TIDAK PERLU lagi loop cari card dan cari user di sini (No Double Work).
	// Biarkan Repository yang mengurus pencarian ID di database.
	return s.cardRepo.AddAssigneesByPublicID(cardPublicID, userPublicIDs)
}
