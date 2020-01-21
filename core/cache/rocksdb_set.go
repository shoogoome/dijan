package cache

// #include <stdlib.h>
// #include "rocksdb/c.h"
// #cgo CFLAGS: -I/rocksdb/include
// #cgo LDFLAGS: -L/rocksdb -lrocksdb -lz -lpthread -lsnappy -lstdc++ -lm -O3
import "C"
import (
	"dijan/models"
	"dijan/utils"
	"encoding/json"
	"time"
	"unsafe"
)


func flush_batch(db *C.rocksdb_t, b *C.rocksdb_writebatch_t, o *C.rocksdb_writeoptions_t) {
	var e *C.char
	C.rocksdb_write(db, o, b, &e)
	if e != nil {
		panic(C.GoString(e))
	}
	C.rocksdb_writebatch_clear(b)
}

func write_func(db *C.rocksdb_t, c chan *pair, o *C.rocksdb_writeoptions_t) {
	count := 0
	barchHandleTime := time.Second * time.Duration(utils.GlobalSystemConfig.Rocksdb.BatchHandleTime)
	t := time.NewTimer(barchHandleTime)
	b := C.rocksdb_writebatch_create()
	for {
		select {
		case p := <-c:
			count++
			key := C.CString(p.k)
			value := C.CBytes(p.v)
			C.rocksdb_writebatch_put(b, key, C.size_t(len(p.k)), (*C.char)(value), C.size_t(len(p.v)))
			C.free(unsafe.Pointer(key))
			C.free(value)
			if count == utils.GlobalSystemConfig.Rocksdb.BatchSize {
				flush_batch(db, b, o)
				count = 0
			}
			if !t.Stop() {
				<-t.C
			}
			t.Reset(barchHandleTime)
		case <-t.C:
			if count != 0 {
				flush_batch(db, b, o)
				count = 0
			}
			t.Reset(barchHandleTime)
		}
	}
}

func (c *rocksdbCache) Set(key string, value []byte, ttl ...int) error {

	// 为
	storageValue := models.StorageValue {
		Value: value,
		TTL: -1,
	}
	if len(ttl) > 0 && ttl[0] > 0 {
		storageValue = models.StorageValue{
			Value: value,
			TTL: time.Now().Unix() + int64(ttl[0]),
		}
	}
	if nValue, err := json.Marshal(storageValue); err != nil {
		c.ch <- &pair{key, value}
	} else {
		c.ch <- &pair{key, nValue}
	}
	return nil
}
