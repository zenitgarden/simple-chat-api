package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/zenitgarden/simple-chat-api/internal/dto"
	"github.com/zenitgarden/simple-chat-api/internal/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id uuid.UUID, trx *gorm.DB) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindAll(ctx context.Context, filter dto.UserFilter) ([]*entity.User, int64, error)
	Update(ctx context.Context, user *entity.User, trx *gorm.DB) error
	BeginTrx(ctx context.Context) *gorm.DB
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	// Start the transaction
	tx := r.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check if the email already exists
	var existingUser entity.User
	if err := tx.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		// Return custom error when email already exists
		tx.Rollback()
		return errors.New("email")
	}

	// Insert user
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	return tx.Commit().Error
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID, trx *gorm.DB) (*entity.User, error) {
	var user entity.User

	db := r.db
	if trx != nil {
		db = trx
	}

	if err := db.WithContext(ctx).
		Select("id, name, email, created_at, updated_at").
		First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).Select("id", "email", "name", "password").First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindAll(ctx context.Context, filter dto.UserFilter) ([]*entity.User, int64, error) {
	var users []*entity.User
	var total int64

	// query := r.db.WithContext(ctx).
	// 	Model(&entity.User{}).
	// 	Where("id <> ?", filter.UserID)

	// if filter.Name != "" {
	// 	query = query.Or("name LIKE ?", "%"+filter.Name+"%")
	// }

	// if err := query.Count(&total).Error; err != nil {
	// 	return nil, 0, err
	// }

	// // Apply ordering, pagination, and retrieve users
	// if err := query.
	// 	Select("id", "name").
	// 	Order("name asc").
	// 	Limit(filter.Limit).
	// 	Offset(filter.Offset).
	// 	Find(&users).Error; err != nil {
	// 	return nil, 0, err
	// }
	// return users, total, nil

	subquery := r.db.
		Table("conversations c").
		Select("1").
		Joins("JOIN participants p1 ON p1.conversation_id = c.id AND p1.user_id = ?", filter.UserID).
		Joins("JOIN participants p2 ON p2.conversation_id = c.id AND p2.user_id = users.id").
		Where("c.is_group = FALSE")

	query := r.db.WithContext(ctx).
		Model(&entity.User{}).
		Where("id <> ?", filter.UserID).
		Where("NOT EXISTS (?)", subquery)

	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Select("id", "name").
		Order("name asc").
		Limit(filter.Limit).
		Offset(filter.Offset).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User, trx *gorm.DB) error {

	db := r.db
	if trx != nil {
		db = trx
	}

	if err := db.WithContext(ctx).Save(user).Error; err != nil {
		return err
	}

	return nil
}

func (r *userRepository) BeginTrx(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).Begin()
}
