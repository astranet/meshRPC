package db

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"
	"sync"

	"github.com/bmatsuo/lmdb-go/lmdb"
	"github.com/tinylib/msgp/msgp"
	"github.com/xlab/catcher"
)

var (
	ErrScanStop  = errors.New("scan stop")
	ErrPeekStop  = ErrScanStop
	ErrRangeStop = ErrScanStop
)

type PeekFunc func(k, v []byte) error
type ModifyFunc func(k, v []byte) ([]byte, error)

type LMDBStore interface {
	Env() *lmdb.Env
	MakeBucket(name string) error
	GetBucket(name string, create bool) (lmdb.DBI, error)
	Put(bucket string, key []byte, value interface{}) error
	Delete(bucket string, key []byte) error
	View(bucket string, key []byte, fn PeekFunc) error
	Update(bucket string, key []byte, fn ModifyFunc) error
	Move(bucket1, bucket2 string, key []byte, fn ModifyFunc) error
	ForEach(bucket string, fn PeekFunc) error
	Range(bucket string, offset []byte, limit uint64, fn PeekFunc) (pos []byte, err error)
	Close() error
}

type store struct {
	mux     *sync.RWMutex
	env     *lmdb.Env
	buckets map[string]lmdb.DBI
}

func NewLMDBStore(envPath string, mapSize int64) (s LMDBStore, err error) {
	defer catcher.Catch(
		catcher.RecvError(&err, false),
		catcher.RecvLog(true),
	)
	env, err := lmdb.NewEnv()
	orPanic(err)
	err = env.SetMapSize(mapSize)
	orPanic(err)
	err = env.SetMaxDBs(256)
	orPanic(err)
	err = env.SetMaxReaders(1024)
	orPanic(err)
	err = env.Open(envPath, 0, 0644)
	orPanic(err)
	s = &store{
		env:     env,
		mux:     new(sync.RWMutex),
		buckets: make(map[string]lmdb.DBI),
	}
	return
}

func (s *store) Env() *lmdb.Env {
	return s.env
}

func (s *store) ForEach(bucket string, fn PeekFunc) error {
	_, err := s.Range(bucket, nil, 0, fn)
	return err
}

func (s *store) MakeBucket(name string) error {
	s.mux.RLock()
	if _, ok := s.buckets[name]; ok {
		s.mux.RUnlock()
		return nil
	}
	s.mux.RUnlock()

	s.mux.Lock()
	defer s.mux.Unlock()
	if err := s.env.Update(func(txn *lmdb.Txn) (err error) {
		dbi, err := txn.OpenDBI(name, lmdb.Create)
		if err != nil {
			return err
		}
		s.buckets[name] = dbi
		return nil
	}); err != nil {
		err := fmt.Errorf("failed to open '%s' bucket: %v", name, err)
		return err
	}
	return nil
}

func (s *store) GetBucket(name string, create bool) (lmdb.DBI, error) {
	s.mux.RLock()
	dbi, ok := s.buckets[name]
	s.mux.RUnlock()
	if ok {
		return dbi, nil
	}
	var flags uint
	if create {
		flags = lmdb.Create
	}
	s.mux.Lock()
	defer s.mux.Unlock()
	if err := s.env.Update(func(txn *lmdb.Txn) (err error) {
		dbi, err := txn.OpenDBI(name, flags)
		if err != nil {
			return err
		}
		s.buckets[name] = dbi
		return nil
	}); err != nil {
		err := fmt.Errorf("failed to open '%s' bucket: %v", name, err)
		return 0, err
	}
	return s.buckets[name], nil
}

func (s *store) Move(bucket1, bucket2 string, key []byte, fn ModifyFunc) error {
	return s.env.Update(func(txn *lmdb.Txn) (err error) {
		checkErr(&err)

		dbi1, err := txn.OpenDBI(bucket1, 0)
		orPanic(err)
		data, err := txn.Get(dbi1, key)
		orPanic(err)
		if fn == nil {
			dbi2, err := txn.CreateDBI(bucket2)
			orPanic(err)
			orPanic(txn.Del(dbi1, key, nil))
			orPanic(txn.Put(dbi2, key, data, 0))
			return nil
		}
		dbi2, err := txn.CreateDBI(bucket2)
		orPanic(err)
		data, err = fn(key, data)
		orPanic(err)
		orPanic(txn.Del(dbi1, key, nil))
		orPanic(txn.Put(dbi2, key, data, 0))
		return nil
	})
}

func (s *store) View(bucket string, key []byte, fn PeekFunc) error {
	return s.env.View(func(txn *lmdb.Txn) error {
		if dbi, err := txn.OpenDBI(bucket, 0); err == nil {
			data, _ := txn.Get(dbi, key)
			return fn(key, data)
		} else {
			return err
		}
		return nil
	})
}

func (s *store) Update(bucket string, key []byte, fn ModifyFunc) error {
	return s.env.Update(func(txn *lmdb.Txn) (err error) {
		if fn == nil {
			return nil
		}
		checkErr(&err)
		dbi, err := txn.OpenDBI(bucket, 0)
		orPanic(err)
		data, err := txn.Get(dbi, key)
		if lmdb.IsNotFound(err) {
			data, err = fn(key, nil)
			orPanic(err)
			orPanic(txn.Put(dbi, key, data, 0))
			return nil
		} else {
			orPanic(err)
		}
		data, err = fn(key, data)
		orPanic(err)
		orPanic(txn.Put(dbi, key, data, 0))
		return nil
	})
}

func (s *store) Range(bucket string, offset []byte, limit uint64, fn PeekFunc) (pos []byte, err error) {
	copyPos := func(key, value []byte, err error) {
		pos = make([]byte, len(key))
		copy(pos, key)
	}
	var read uint64
	err = s.env.RunTxn(lmdb.Readonly, func(txn *lmdb.Txn) (err error) {
		defer catcher.Catch(
			catcher.RecvError(&err, false),
		)
		dbi, err := s.GetBucket(bucket, false)
		orPanic(err)
		txn.RawRead = true

		cur, err := txn.OpenCursor(dbi)
		orPanic(err)
		defer cur.Close()

		var beginWith uint = lmdb.Next
		if len(offset) > 0 {
			beginWith = lmdb.SetRange
		}
		k, v, err := cur.Get(offset, nil, beginWith)
		for err == nil {
			read++
			if err := fn(k, v); err != ErrPeekStop {
				pos = nil
				orPanic(err)
				if limit > 0 && read >= limit {
					copyPos(cur.Get(nil, nil, lmdb.Next))
					return nil
				}
			} else {
				copyPos(cur.Get(nil, nil, lmdb.Next))
				return nil
			}
			k, v, err = cur.Get(nil, nil, lmdb.Next)
		}
		if lmdb.IsNotFound(err) {
			pos = nil
			return nil
		}
		pos = nil
		return err
	})
	return
}

func PutValue(txn *lmdb.Txn, dbi lmdb.DBI, key []byte, value interface{}) error {
	buf := new(bytes.Buffer)
	if v, ok := value.(msgp.Encodable); ok {
		if err := msgp.Encode(buf, v); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("value is not msgp.Encodable: %T", value)
	}
	return txn.Put(dbi, key, buf.Bytes(), 0)
}

func (s *store) Delete(bucketName string, key []byte) error {
	return s.env.Update(func(txn *lmdb.Txn) error {
		if dbi, err := txn.OpenDBI(bucketName, 0); err == nil {
			return txn.Del(dbi, key, nil)
		}
		return nil
	})
}

func (s *store) Put(bucketName string, key []byte, value interface{}) error {
	buf := new(bytes.Buffer)
	if v, ok := value.(msgp.Encodable); ok {
		if err := msgp.Encode(buf, v); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("value is not msgp.Encodable: %T", value)
	}
	return s.env.Update(func(txn *lmdb.Txn) (err error) {
		defer checkErr(&err)

		s.mux.RLock()
		dbi, ok := s.buckets[bucketName]
		s.mux.RUnlock()
		if !ok {
			s.mux.Lock()
			bucket, err := txn.OpenDBI(bucketName, lmdb.Create)
			if err == nil {
				s.buckets[bucketName] = bucket
				dbi = bucket
			}
			s.mux.Unlock()
		}

		return txn.Put(dbi, key, buf.Bytes(), 0)
	})
}

func (s *store) Close() error {
	if s.env != nil {
		err := s.env.Close()
		s.env = nil
		return err
	}
	return nil
}

func orPanic(err error, finalizers ...func()) {
	if err != nil {
		for _, fn := range finalizers {
			fn()
		}
		panic(err)
	}
}

func checkErr(err *error) {
	if v := recover(); v != nil {
		*err = fmt.Errorf("%+v", v)
	}
}

func checkErrStack(err *error) {
	if v := recover(); v != nil {
		stack := make([]byte, 32*1024)
		n := runtime.Stack(stack, false)
		switch event := v.(type) {
		case error:
			*err = fmt.Errorf("%s\n%s", event.Error(), stack[:n])
		default:
			*err = fmt.Errorf("%+v %s", v, stack[:n])
		}
	}
}
