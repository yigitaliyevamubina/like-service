package storage

import (
	"database/sql"
	"like-service/storage/postgres"
	"like-service/storage/repo"
)

type IStorage interface {
	Like() repo.LikeServiceI
}

type storagePg struct {
	db       *sql.DB
	likeRepo repo.LikeServiceI
}

func NewStoragePg(db *sql.DB) *storagePg {
	return &storagePg{
		db:       db,
		likeRepo: postgres.NewLikeRepo(db),
	}
}

func (s storagePg) Like() repo.LikeServiceI {
	return s.likeRepo
}
