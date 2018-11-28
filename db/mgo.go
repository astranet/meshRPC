package db

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	bugsnag "github.com/bugsnag/bugsnag-go"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	log "github.com/sirupsen/logrus"
)

type M = bson.M
type MongoIndex = mgo.Index
type MongoChange = mgo.Change
type ObjectId = bson.ObjectId

type DBMongo struct {
	Config *MongoConfig

	sessions   map[string]*mgo.Session
	sessionMux *sync.RWMutex
}

type MongoConfig struct {
	Connection string
	SSLCert    string
	Database   string
	Debug      bool
}

func NewDBMongo(mgoConfig *MongoConfig) *DBMongo {
	return &DBMongo{
		Config: mgoConfig,
	}
}

func (d *DBMongo) Remove(collection *mgo.Collection, selector interface{}) error {
	err := collection.Remove(selector)
	if err != nil {
		collection.Database.Session.Refresh()
		err = collection.RemoveId(selector)
	}
	return err
}

func (d *DBMongo) RemoveId(collection *mgo.Collection, id string) error {
	objectId := bson.ObjectIdHex(id)
	err := collection.RemoveId(objectId)
	if err != nil {
		collection.Database.Session.Refresh()
		err = collection.RemoveId(objectId)
	}
	return err
}

func FieldSelector(q []string) (r bson.M) {
	r = make(bson.M, len(q))
	for _, s := range q {
		r[s] = 1
	}
	return
}

func (d *DBMongo) Upsert(collection *mgo.Collection,
	selector interface{}, data interface{}) (*mgo.ChangeInfo, error) {

	info, err := collection.Upsert(selector, bson.M{"$set": data})
	if err != nil {
		collection.Database.Session.Refresh()
		info, err = collection.Upsert(selector, bson.M{"$set": data})
	}
	return info, err
}

func (d *DBMongo) UpsertId(collection *mgo.Collection,
	id interface{}, data interface{}) (*mgo.ChangeInfo, error) {

	info, err := collection.UpsertId(id, bson.M{"$set": data})
	if err != nil {
		collection.Database.Session.Refresh()
		info, err = collection.UpsertId(id, bson.M{"$set": data})
	}
	return info, err
}

func (d *DBMongo) Query(q *mgo.Query, result interface{}) error {
	timer1 := time.Now()

	err := q.All(result)

	if d.Config.Debug {
		endTime := time.Since(timer1)
		log.Infof("mongo query: %v. execution time: %v (err: %+v)", q, endTime, err)
	}
	return err
}

func (d *DBMongo) QuerySingle(q *mgo.Query, result interface{}) error {
	timer1 := time.Now()

	err := q.One(result)

	if d.Config.Debug {
		endTime := time.Since(timer1)
		log.Infof("mongo query: %v. execution time: %v (err: %+v)", q, endTime, err)
	}
	return err
}

func (d *DBMongo) UpdateId(collection *mgo.Collection, id interface{}, data interface{}) (err error) {
	err = collection.UpdateId(id, data)
	if err != nil {
		collection.Database.Session.Refresh()
		err = collection.UpdateId(id, data)
	}
	return
}

func (d *DBMongo) Insert(collection *mgo.Collection, data interface{}) (err error) {
	err = collection.Insert(data)
	if err != nil {
		collection.Database.Session.Refresh()
		err = collection.Insert(data)
	}
	return
}

func (d *DBMongo) GetCollection(name string) *mgo.Collection {
	mdb := d.GetSession(true).DB(d.Config.Database)
	collection := mdb.C(name)
	runtime.SetFinalizer(collection, func(obj *mgo.Collection) {
		if obj != nil && obj.Database != nil && obj.Database.Session != nil {
			obj.Database.Session.Close()
		}
	})
	return collection
}

func (d *DBMongo) GetDatabase() *mgo.Database {
	mdb := d.GetSession(true).DB(d.Config.Database)
	runtime.SetFinalizer(mdb, func(obj *mgo.Database) {
		if obj != nil && obj.Session != nil {
			obj.Session.Close()
		}
	})
	return mdb
}

func (d *DBMongo) GetSession(monotonic bool) *mgo.Session {
	var s *mgo.Session
	if len(d.Config.SSLCert) > 0 {
		s = d.getSession(d.Config.Connection, d.Config.SSLCert)
	} else {
		s = d.getSession(d.Config.Connection)
	}
	s.SetMode(mgo.Monotonic, monotonic)
	return s
}

func (d *DBMongo) getSession(connectionVariable string, cert ...string) *mgo.Session {
	d.sessionMux.Lock()
	defer d.sessionMux.Unlock()
	session, ok := d.sessions[connectionVariable]
	if !ok || session == nil {
		var cs string
		var ssl bool

		if strings.HasPrefix(connectionVariable, "mongodb://") {
			cs = connectionVariable
		} else {
			cs = os.Getenv(connectionVariable)
		}
		if strings.Contains(cs, "ssl=true") {
			cs = strings.Replace(cs, "ssl=true", "", -1)
			cs = strings.Replace(cs, "?&", "?", -1)
			cs = strings.Replace(cs, "&&", "&", -1)
			ssl = true
		}

		var err error
		if cert != nil || ssl {
			session, err = dialWithSSL(cs, cert)
		} else {
			session, err = mgo.Dial(cs)
		}
		if err != nil {
			bugsnag.Notify(err)
		}

		session.SetSocketTimeout(30 * time.Second)
		session.SetSyncTimeout(30 * time.Second)
		d.sessions[connectionVariable] = session
		return session
	}
	session.Refresh()
	return session.Copy()
}

func dialWithSSL(cs string, certs []string) (session *mgo.Session, err error) {
	tlsConfig := &tls.Config{}

	if certs != nil {
		roots := x509.NewCertPool()
		roots.AppendCertsFromPEM([]byte(certs[0]))
		tlsConfig.RootCAs = roots
	} else {
		tlsConfig.InsecureSkipVerify = true
	}

	dialInfo, err := mgo.ParseURL(cs)
	if err != nil {
		return
	}
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}
	session, err = mgo.DialWithInfo(dialInfo)
	return
}

// ErrNotFound is returned from the CheckMongoError method when no results are found
var ErrNotFound = errors.New("object not found in mongo repository")

// CheckMongoError will check the error to see if it is the 'not found' error thrown when
// a document is not found in mongo.
func CheckMongoError(mgoError error) (bool, error) {
	if mgoError == nil {
		return false, nil
	}
	if mgoError == mgo.ErrNotFound {
		return false, ErrNotFound
	}
	return true, mgoError
}

func ToObjectID(id string) (objectID bson.ObjectId, err error) {
	if !bson.IsObjectIdHex(id) {
		err = errors.New("Not a valid ObjectId")
		return
	}
	objectID = bson.ObjectIdHex(id)
	return
}
