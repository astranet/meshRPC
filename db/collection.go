package db

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Collection is a wrapper for the mgo.Collection. Use GetCollection to retrieve one.
type Collection struct {
	*mgo.Collection
}

// GetCollectionWrapper returns a collection object (which is just a wrapped mgo.Collection object).
func (d *DBMongo) GetCollectionWrapper(name string) *Collection {
	return &Collection{
		d.GetCollection(name),
	}
}

// PartialUpdate helper which does a query and apply the update and returns the new object
func (c *Collection) PartialUpdate(query bson.M, update bson.M, result interface{}) (err error) {
	change := mgo.Change{
		Update:    update,
		ReturnNew: true,
	}
	_, err = c.Find(query).Apply(change, result)
	_, err = CheckMongoError(err)
	return
}

// PartialUpdateNoReturn helper which does a query and apply the update without returning the new object
func (c *Collection) PartialUpdateNoReturn(query bson.M, update bson.M) (err error) {
	change := mgo.Change{
		Update:    update,
		ReturnNew: false,
	}
	_, err = c.Find(query).Apply(change, nil)
	_, err = CheckMongoError(err)
	return
}
