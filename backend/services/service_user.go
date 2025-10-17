package services

import (
	"backend/models"
	"backend/repositories"
	"errors"
	"unicode/utf8"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetAll() ([]models.User, error) {
	return s.repo.GetAll()
}

func (s *UserService) GetByID(id int64) (*models.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) Create(req *models.CreateUserRequest) (*models.User, error) {
	// Validation: names must not exceed 3 characters
	if utf8.RuneCountInString(req.FirstName) > 3 {
		return nil, errors.New("first_name must not exceed 3 characters")
	}
	if utf8.RuneCountInString(req.LastName) > 3 {
		return nil, errors.New("last_name must not exceed 3 characters")
	}
	if req.FirstName == "" || req.LastName == "" {
		return nil, errors.New("first_name and last_name are required")
	}

	now := models.Now()
	user := &models.User{
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Email:         req.Email,
		Phone:         req.Phone,
		AvatarURL:     req.AvatarURL,
		Bio:           req.Bio,
		PointsBalance: 0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	err := s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Update(id int64, req *models.UpdateUserRequest) (*models.User, error) {
	// Validation: names must not exceed 3 characters
	if req.FirstName != "" && utf8.RuneCountInString(req.FirstName) > 3 {
		return nil, errors.New("first_name must not exceed 3 characters")
	}
	if req.LastName != "" && utf8.RuneCountInString(req.LastName) > 3 {
		return nil, errors.New("last_name must not exceed 3 characters")
	}

	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.FirstName != "" {
		existing.FirstName = req.FirstName
	}
	if req.LastName != "" {
		existing.LastName = req.LastName
	}
	if req.Email != "" {
		existing.Email = req.Email
	}
	if req.Phone != "" {
		existing.Phone = req.Phone
	}
	if req.AvatarURL != "" {
		existing.AvatarURL = req.AvatarURL
	}
	if req.Bio != "" {
		existing.Bio = req.Bio
	}

	existing.UpdatedAt = models.Now()

	err = s.repo.Update(id, existing)
	if err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *UserService) Delete(id int64) error {
	return s.repo.Delete(id)
}
