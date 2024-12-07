package storage

import (
	"Questify/internal/user"
	"Questify/pkg/adapters/storage/entities"
	"Questify/pkg/adapters/storage/mappers"
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) user.Repo {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) Create(ctx context.Context, user *user.User) error {
	newUser := mappers.UserDomainToEntity(user)
	err := r.db.WithContext(ctx).Create(&newUser).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil
		}
		return err
	}
	user.ID = newUser.ID
	return nil
}

func (r *userRepo) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	var u entities.User

	err := r.db.WithContext(ctx).Model(&entities.User{}).Where("id = ?", id).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	uu := mappers.UserEntityToDomain(u)
	return &uu, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).Model(&entities.User{}).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	uu := mappers.UserEntityToDomain(user)
	return &uu, nil
}

func (r *userRepo) UpdateUser(ctx context.Context, user *user.User) error {
	updatedUser := mappers.UserDomainToEntity(user)
	err := r.db.WithContext(ctx).Model(&entities.User{}).Where("id = ?", user.ID).Updates(updatedUser).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepo) GetUsers(ctx context.Context, page, pageSize int) ([]user.User, int64, error) {
	offset := (page - 1) * pageSize

	var totalCount int64
	var users []entities.User

	err := r.db.WithContext(ctx).Model(&entities.User{}).Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).
		Model(&entities.User{}).
		Limit(pageSize).
		Offset(offset).
		Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	result := []user.User{}
	for _, u := range users {
		domainUser := mappers.UserEntityToDomain(u)
		result = append(result, domainUser)
	}

	return result, totalCount, nil
}

func (r *userRepo) GetByNationalCode(ctx context.Context, nationalCode string) (*user.User, error) {
    var userEntity entities.User
    result := r.db.WithContext(ctx).Where("national_code = ?", nationalCode).First(&userEntity)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return nil, nil
        }
        return nil, result.Error
    }

    // Map the entity to the domain type
    userDomain := mappers.UserEntityToDomain(userEntity)
    return &userDomain, nil
}
