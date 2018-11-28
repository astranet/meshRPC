package data

import (
	bugsnag "github.com/bugsnag/bugsnag-go"

	"github.com/astranet/galaxy/db"
)

type BaseRepoFactory interface {
	GetRepo(string) BaseRepo
}

type ExtendedRepoFactory interface {
	BaseRepoFactory
	GetExtendedRepo(repo string, indexes ...interface{}) ExtendedBaseRepo
}

type BaseRepo interface {
	GetById(id interface{}, out interface{}) error
	GetMany(query interface{}, out interface{}) error
	GetOne(query interface{}, out interface{}) error
	Upsert(key interface{}, doc interface{}) error
	UpdateSingle(key interface{}, doc interface{}) error
	UpdateMany(query interface{}, update interface{}) error
	Insert(doc interface{}) error
	DeleteById(id interface{}) error
}

type ExtendedBaseRepo interface {
	BaseRepo
	PipeMany(pipeline interface{}, out interface{}) error
	PipeOne(pipeline interface{}, out interface{}) error
	UpdateAndReturn(query interface{}, update interface{}, out interface{}, upsert bool) error
}

func NewRepo(mgo *db.DBMongo) Repo {
	return &repo{
		db: mgo,
	}
}

type Repo interface {
	DB() *db.DBMongo
	SavePartialModel(collectionName string, query db.M, updateModel interface{}, returnModel interface{}) error
	SavePartialModelNoReturn(collectionName string, query db.M, updateModel interface{}) error
	SavePartial(collectionName string, query db.M, updates db.M, returnModel interface{}) error
	SavePartialNoReturn(collectionName string, query db.M, updates db.M) error
	DocumentExists(collectionName string, query db.M) (bool, error)
}

type repo struct {
	db *db.DBMongo
}

func (r *repo) DB() *db.DBMongo {
	return r.db
}

func (r *repo) SavePartialModel(collectionName string, query db.M, updateModel interface{}, returnModel interface{}) error {
	updates, err := db.GeneratePartialSet(updateModel)
	if err != nil {
		bugsnag.Notify(err)
		return err
	}

	err = r.SavePartial(collectionName, query, updates, returnModel)
	_, err = db.CheckMongoError(err)
	return err
}

func (r *repo) SavePartialModelNoReturn(collectionName string, query db.M, updateModel interface{}) error {
	updates, err := db.GeneratePartialSet(updateModel)
	if err != nil {
		bugsnag.Notify(err)
		return err
	}

	err = r.SavePartialNoReturn(collectionName, query, updates)
	_, err = db.CheckMongoError(err)
	return err
}

func (r *repo) SavePartial(collectionName string, query db.M, updates db.M, returnModel interface{}) error {
	db.SetUpdatedOn(&updates)

	collection := r.db.GetCollectionWrapper(collectionName)
	err := collection.PartialUpdate(query, updates, returnModel)
	_, err = db.CheckMongoError(err)
	return err
}

func (r *repo) SavePartialNoReturn(collectionName string, query db.M, updates db.M) (err error) {
	db.SetUpdatedOn(&updates)

	collection := r.db.GetCollectionWrapper(collectionName)
	err = collection.PartialUpdateNoReturn(query, updates)
	_, err = db.CheckMongoError(err)
	return
}

func (r *repo) DocumentExists(collectionName string, query db.M) (bool, error) {
	collection := r.db.GetCollection(collectionName)
	count, err := collection.Find(query).Limit(1).Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
