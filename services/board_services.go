package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/justynwen0/project-management-API/models"
	"github.com/justynwen0/project-management-API/repositories"
)

type BoardService interface {
	Create(board *models.Board) error
	Update(board *models.Board) error
	Delete(publicID string) error
	GetByPublicID(publicID string) (*models.Board, error)
	AddMembers(boardPublicID string, userPublicIDs []string) error
	RemoveMembers(boardPublicID string, userPublicIDs []string) error
	GetAllByUserPaginate(userPublicID, filter, sort string, limit,
		offset int) ([]models.Board, int64, error)
	GetMembersByBoardID(boardPublicID string) ([]models.User, error)
}

type boardService struct {
	boardRepo       repositories.BoardRepository
	userRepo        repositories.UserRepository
	boardMemberRepo repositories.BoardMemberRepository
}

func NewBoardService(
	boardRepo repositories.BoardRepository,
	userRepo repositories.UserRepository,
	boardMemberRepo repositories.BoardMemberRepository,
) BoardService {
	return &boardService{boardRepo, userRepo, boardMemberRepo}
}

func (s *boardService) Create(board *models.Board) error {
	user, err := s.userRepo.FindByPublicID(board.OwnerPublicID.String())
	if err != nil {
		return errors.New("owner not found")
	}
	board.PublicID = uuid.New()
	board.OwnerID = user.InternalID
	return s.boardRepo.Create(board)
}

func (s *boardService) Update(board *models.Board) error {
	return s.boardRepo.Update(board)
}

func (s *boardService) Delete(publicID string) error {

	_, err := s.boardRepo.FindByPublicID(publicID)

	if err != nil {
		return errors.New("board not found")
	}

	return s.boardRepo.Delete(publicID)
}

func (s *boardService) GetByPublicID(publicID string) (*models.Board, error) {
	return s.boardRepo.FindByPublicID(publicID)
}

func (s *boardService) AddMembers(boardPublicID string, userPublicIDs []string) error {
	board, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		return errors.New("board not found")
	}

	var userInternalIDs []uint
	for _, userPublicID := range userPublicIDs {
		user, err := s.userRepo.FindByPublicID(userPublicID)
		if err != nil {
			return errors.New("user not found: " + userPublicID)
		}
		userInternalIDs = append(userInternalIDs, uint(user.InternalID))
	}

	// cek keanggotaan
	existingMembers, err := s.boardMemberRepo.GetMembers(string(board.PublicID.String()))
	if err != nil {
		return err
	}

	// cek cepat menggunakan map
	memberMap := make(map[uint]bool)
	for _, member := range existingMembers {
		memberMap[uint(member.InternalID)] = true //membermap[1] = true
	}

	var newMembersIDs []uint
	for _, userID := range userInternalIDs {
		if !memberMap[userID] {
			newMembersIDs = append(newMembersIDs, userID)
		}
	}
	if len(newMembersIDs) == 0 {
		return nil
	}
	return s.boardRepo.AddMember(uint(board.InternalID), newMembersIDs)
}

func (s *boardService) RemoveMembers(boardPublicID string, userPublicIDs []string) error {
	board, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		return errors.New("board not found")
	}

	// Validasi user
	var userInternalIDs []uint
	for _, userPublicID := range userPublicIDs {
		user, err := s.userRepo.FindByPublicID(userPublicID)
		if err != nil {
			return errors.New("user not found: " + userPublicID)
		}
		userInternalIDs = append(userInternalIDs, uint(user.InternalID))
	}

	// cek keanggotaan
	existingMembers, err := s.boardMemberRepo.GetMembers(string(board.PublicID.String()))
	if err != nil {
		return err
	}

	// cek cepat menggunakan map
	memberMap := make(map[uint]bool)
	for _, member := range existingMembers {
		memberMap[uint(member.InternalID)] = true //membermap[1] = true
	}

	var membersToRemove []uint
	for _, userID := range userInternalIDs {
		if memberMap[userID] {
			membersToRemove = append(membersToRemove, userID)
		}
	}

	return s.boardRepo.RemoveMember(uint(board.InternalID), membersToRemove)
}

func (s *boardService) GetAllByUserPaginate(userPublicID, filter, sort string, limit,
	offset int) ([]models.Board, int64, error) {
	return s.boardRepo.FindAllByUserPaginate(userPublicID, filter, sort, limit, offset)
}

func (s *boardService) GetMembersByBoardID(boardPublicID string) ([]models.User, error) {
	return s.boardMemberRepo.GetMembers(boardPublicID)
}
