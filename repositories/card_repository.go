package repositories

import (
	"fmt"
	"path/filepath"

	"github.com/Filbertfelix888/project-management-API/config"
	"github.com/Filbertfelix888/project-management-API/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CardRepository interface {
	Create(card *models.Card) error
	Update(card *models.Card) error
	Delete(id uint) error
	FindByID(id uint) (*models.Card, error)
	FindByPublicID(publicID string) (*models.Card, error)
	FindByListID(listID string) ([]models.Card, error)

	FindCardPositionByListID (id int64) (*models.CardPosition, error)
	UpdatePosition(listID string, position []string) error

	AddAssignees(cardID uint, userIDs []uint) error
	RemoveAssignees(cardID uint, userIDs []uint) error
	AddAssigneesByPublicID(publicID string, userUUIDs []string) error
}

type cardRepository struct {
}

func NewCardRepository() CardRepository {
	return &cardRepository{}
}

func (r *cardRepository) Create(card *models.Card) error {
	return config.DB.Create(card).Error
}

func (r *cardRepository) Update(card *models.Card) error {
	return config.DB.Save(card).Error
}

func (r *cardRepository) Delete(id uint) error {
	return config.DB.Delete(&models.Card{}, id).Error
}

func (r *cardRepository) FindByID(id uint) (*models.Card, error) {
	var card models.Card
	err := config.DB.Preload("Labels").Preload("Assignees").First(&card, id).Error

	return &card, err
}

func (r *cardRepository) FindByPublicID(publicID string) (*models.Card, error) {
	var card models.Card
	if err := config.DB.Preload("Assignees.User", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("Internal_id", "public_id", "name", "email")
	}).Preload("Attachments").Where("public_id = ?", publicID).First(&card).Error; err != nil {
		return nil, err
	}

	baseUrl := config.AppConfig.APPURL

	for i := range card.Attachments {
		card.Attachments[i].FileUrl = fmt.Sprintf("%s/files/%s",
			baseUrl,
			filepath.Base(card.Attachments[i].File),
		)
	}

	return &card, nil
}

func (r *cardRepository) FindByListID(listID string) ([]models.Card, error) {
	var cards []models.Card
	err := config.DB.Joins("JOIN lists ON lists.internal_id = cards.list_internal_id").
		Where("lists.public_id = ?", listID).
		Order("position ASC").
		Find(&cards).Error
	return cards, err
}

func (r * cardRepository) FindCardPositionByListID(id int64) (*models.CardPosition, error) {
	var position models.CardPosition
	err := config.DB.Where("list_internal_id = ?",id).First(&position).Error
	if err != nil {
		return nil, err
	}
	return &position, nil
}

func (r *cardRepository) UpdatePosition(listID string, position []string) error {
	return config.DB.Model(&models.CardPosition{}).
	Where("list_internal_id = (SELECT internal_id FROM lists Where public_id = ? )",listID). 
	Update("card_order", position).Error
}

func (r *cardRepository) AddAssignees(cardID uint, userIDs []uint) error {
    var assignees []models.CardAssignee
    for _, uid := range userIDs {
        assignees = append(assignees, models.CardAssignee{
            CardID: int64(cardID),
            UserID: int64(uid),
        })
    }
   return config.DB.Clauses(clause.OnConflict{DoNothing: true}).
        Create(&assignees).Error
}

func (r *cardRepository) RemoveAssignees(cardID uint, userIDs []uint) error {
    return config.DB.Where("card_internal_id = ? AND user_internal_id IN ?", cardID, userIDs).
        Delete(&models.CardAssignee{}).Error
}

// Tambahkan implementasinya di bawah
func (r *cardRepository) AddAssigneesByPublicID(cardPublicID string, userUUIDs []string) error {
    var card models.Card
    // 1. Cari card-nya
    if err := config.DB.Where("public_id = ?", cardPublicID).First(&card).Error; err != nil {
        return err
    }

    // 2. Cari semua Internal ID User berdasarkan UUID (string) yang dikirim frontend
    var userInternalIDs []uint
    if err := config.DB.Model(&models.User{}).
        Where("public_id IN ?", userUUIDs).
        Pluck("internal_id", &userInternalIDs).Error; err != nil {
        return err
    }

    // 3. Panggil fungsi insert internal yang sudah kamu buat sebelumnya
    // Gunakan uint() untuk card.InternalID agar tidak error int64 vs uint
    return r.AddAssignees(uint(card.InternalID), userInternalIDs)
}