package db

import (
	"context"
	"errors"
	"sync"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

type M = bson.M

type DBMongo struct {
	Client *mongo.Client
	Config *MongoConfig

	sessions   map[string]*mongo.Session
	sessionMux *sync.RWMutex
}

type MongoConfig struct {
	Connection string
	SSLCert    string
	Database   string
	Debug      bool
}

func NewDBMongo(config *MongoConfig) (*DBMongo, error) {
	cli, err := mongo.Connect(context.TODO(), config.Connection)
	if err != nil {
		return nil, err
	}
	db := &DBMongo{
		Client:     cli,
		Config:     config,
		sessions:   make(map[string]*mongo.Session),
		sessionMux: new(sync.RWMutex),
	}
	return db, nil
}

func (d *DBMongo) GetDatabase() *mongo.Database {
	return d.Client.Database(d.Config.Database)
}

func (d *DBMongo) GetCollection(name string) *mongo.Collection {
	return d.Client.Database(d.Config.Database).Collection(name)
}

func (d *DBMongo) GetSession() (mongo.Session, error) {
	return d.Client.StartSession()
}

// ErrNotFound is returned from the CheckMongoError method when no results are found
var ErrNotFound = errors.New("object not found in mongo repository")

func CheckMongoError(err error) (bool, error) {
	if err == nil {
		return false, nil
	}
	if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
		return false, ErrNotFound
	}
	return true, err
}

func MakeIndex(unique bool, keys interface{}) mongo.IndexModel {
	idx := mongo.IndexModel{
		Keys:    keys,
		Options: options.Index(),
	}
	if unique {
		idx.Options = idx.Options.SetUnique(true)
	}
	return idx
}
