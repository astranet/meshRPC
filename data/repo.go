package data

import (
	db "github.com/astranet/galaxy/db_mongo"
)

func NewRepo(db *db.DBMongo) Repo {
	return &repo{
		db: db,
	}
}

type Repo interface {
	DB() *db.DBMongo
}

type repo struct {
	db *db.DBMongo
}

func (r *repo) DB() *db.DBMongo {
	return r.db
}
