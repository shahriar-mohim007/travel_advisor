package dependencies

import (
	"travel_advisor/districts/repository"
	"travel_advisor/domain"
	"travel_advisor/pkg/cache"
	"travel_advisor/pkg/conn"
)

type RepositoryInterfaces struct {
	Districts domain.DistrictRepository
	Cacher    cache.Cache
}

func InjectRepositories() RepositoryInterfaces {
	db := conn.DefaultDB()
	districRepository := repository.NewDistrictPostgreSQL(db)
	cacher := conn.DefaultCache()
	return RepositoryInterfaces{
		Districts: districRepository,
		Cacher:    cacher,
	}
}
