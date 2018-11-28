package data

import (
	"errors"

	"github.com/astranet/galaxy/db"
)

type mongoRepoFactory struct {
	db *db.DBMongo
}

func NewMongoRepoFactory(cfg *db.MongoConfig) ExtendedRepoFactory {
	return &mongoRepoFactory{
		db: db.NewDBMongo(cfg),
	}
}

type collectionFunction func(collection *db.Collection) error
type withCollectionFunction func(fn collectionFunction) error

func (f *mongoRepoFactory) GetRepo(collection string) BaseRepo {
	return &mongoRepo{
		withCollection: func(fn collectionFunction) error {
			c := f.db.GetCollectionWrapper(collection)
			defer c.Database.Session.Close()
			return fn(c)
		},
	}
}
func (f *mongoRepoFactory) GetExtendedRepo(collection string, indexes ...interface{}) ExtendedBaseRepo {
	return &mongoRepo{
		withCollection: func(fn collectionFunction) error {
			c := f.db.GetCollectionWrapper(collection)

			for _, idx := range indexes {
				if mgoIdx, ok := idx.(db.MongoIndex); ok {
					c.EnsureIndex(mgoIdx)
				}
			}

			defer c.Database.Session.Close()
			return fn(c)
		},
	}
}

type MongoRepo interface {
	GetById(id interface{}, out interface{}) error
	GetOne(query interface{}, out interface{}) error
	PipeOne(pipeline interface{}, out interface{}) error
	GetMany(query interface{}, out interface{}) error
	PipeMany(pipeline interface{}, out interface{}) error
	Upsert(key interface{}, doc interface{}) error
	UpdateSingle(key interface{}, doc interface{}) error
	UpdateMany(query interface{}, doc interface{}) error
	Insert(doc interface{}) error
	DeleteById(id interface{}) error
	UpdateAndReturn(query interface{}, doc interface{}, out interface{}, upsert bool) error
}

type mongoRepo struct {
	withCollection withCollectionFunction
}

func (r *mongoRepo) isValidQuery(query interface{}) (error, bool) {
	if _, ok := query.(db.M); !ok {
		return errors.New("MongoRepo requires queries to be of type db.M or bson.M"), false
	}
	return nil, true
}

func (r *mongoRepo) GetById(id interface{}, out interface{}) error {
	return r.withCollection(func(c *db.Collection) error {
		return c.FindId(id).One(out)
	})
}

func (r *mongoRepo) GetOne(query interface{}, out interface{}) error {
	return r.withCollection(func(c *db.Collection) error {
		return c.Find(query).One(out)
	})
}
func (r *mongoRepo) PipeOne(pipeline interface{}, out interface{}) error {
	return r.withCollection(func(c *db.Collection) error {
		return c.Pipe(pipeline).One(out)
	})
}
func (r *mongoRepo) GetMany(query interface{}, out interface{}) error {
	return r.withCollection(func(c *db.Collection) error {
		return c.Find(query).All(out)
	})
}
func (r *mongoRepo) PipeMany(pipeline interface{}, out interface{}) error {
	return r.withCollection(func(c *db.Collection) error {
		return c.Pipe(pipeline).All(out)
	})
}

func (r *mongoRepo) Upsert(key interface{}, doc interface{}) error {
	return r.withCollection(func(c *db.Collection) error {
		keySet, err := db.GeneratePartialSet(key)
		if err != nil {
			return err
		}
		_, err = c.Upsert(keySet, doc)
		return err
	})
}

func (r *mongoRepo) UpdateSingle(key interface{}, doc interface{}) error {
	return r.withCollection(func(c *db.Collection) error {
		err := c.Update(key, doc)
		return err
	})
}

func (r *mongoRepo) UpdateMany(query interface{}, doc interface{}) error {
	return r.withCollection(func(c *db.Collection) error {
		_, err := c.UpdateAll(query, doc)
		return err
	})
}

func (r *mongoRepo) Insert(doc interface{}) error {
	return r.withCollection(func(c *db.Collection) error {
		return c.Insert(doc)
	})
}

func (r *mongoRepo) DeleteById(id interface{}) error {
	return r.withCollection(func(c *db.Collection) error {
		return c.RemoveId(id)
	})
}

func (r *mongoRepo) UpdateAndReturn(query interface{}, doc interface{}, out interface{}, upsert bool) error {
	return r.withCollection(func(c *db.Collection) error {
		change := db.MongoChange{
			Update:    doc,
			ReturnNew: true,
			Upsert:    upsert,
		}
		_, err := c.Find(query).Apply(change, out)
		return err
	})
}
