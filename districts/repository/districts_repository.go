package repository

import (
	"context"
	"fmt"
	"travel_advisor/domain"
	"travel_advisor/pkg/conn"

	"gorm.io/gorm"
)

type DistrictPostgreSQL struct {
	db *conn.DB
}

func NewDistrictPostgreSQL(db *conn.DB) domain.DistrictRepository {
	return &DistrictPostgreSQL{
		db: db,
	}
}

func (r *DistrictPostgreSQL) List(ctx context.Context, ctr *domain.DistrictCriteria) ([]*domain.District, error) {
	qry := r.db.DB.WithContext(ctx)

	if ctr.DistrictName != nil && *ctr.DistrictName != "" {
		qry = qry.Where("name = ?", *ctr.DistrictName)
	}
	var districtList = make([]*domain.District, 0)
	if err := qry.Find(&districtList).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrDistrictNotFound
		}
		return nil, fmt.Errorf("repository:postgreSQL: failed to Fetch city: %v", err)
	}

	return districtList, nil
}
