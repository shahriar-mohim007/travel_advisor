package repository

import (
	"context"
	"fmt"
	"travel_advisor/domain"
	"travel_advisor/pkg/conn"

	"gorm.io/gorm"
)

type UserPostgreSQL struct {
	db *conn.DB
}

func NewUserPostgreSQL(db *conn.DB) domain.UserRepository {
	return &UserPostgreSQL{
		db: db,
	}
}

func (r *UserPostgreSQL) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	if err := r.db.DB.WithContext(ctx).Create(user).Error; err != nil {
		return nil, fmt.Errorf("repository:postgreSQL: failed to create user: %v", err)
	}
	return user, nil
}

func (r *UserPostgreSQL) Get(ctx context.Context, ctr *domain.UserCriteria) (*domain.User, error) {
	qry := r.db.DB.WithContext(ctx)

	if ctr.ID != nil && *ctr.ID != 0 {
		qry = qry.Where("id = ?", *ctr.ID)
	}

	if ctr.Email != nil && *ctr.Email != "" {
		qry = qry.Where("email = ?", *ctr.Email)
	}

	var user domain.User
	if err := qry.First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("repository:postgreSQL: failed to fetch user: %v", err)
	}

	return &user, nil
}
