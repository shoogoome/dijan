package cache

// #include <stdlib.h>
// #include "rocksdb/c.h"
// #cgo CFLAGS: -I/rocksdb/include
// #cgo LDFLAGS: -L/rocksdb -lrocksdb -lz -lpthread -lsnappy -lstdc++ -lm -O3
import "C"
import (
	"dijan/models"
	"encoding/json"
	"errors"
	"time"
	"unsafe"
)

func (c *rocksdbCache) Get(key string) ([]byte, error) {
	k := C.CString(key)
	defer C.free(unsafe.Pointer(k))
	var length C.size_t
	v := C.rocksdb_get(c.db, c.ro, k, C.size_t(len(key)), &length, &c.e)
	if c.e != nil {
		return nil, errors.New(C.GoString(c.e))
	}
	defer C.free(unsafe.Pointer(v))

	var storageValue models.StorageValue
	value := C.GoBytes(unsafe.Pointer(v), C.int(length))
	// 转码出错直接return原本值
	if err := json.Unmarshal(value, &storageValue); err != nil {
		return value, nil
	}
	if storageValue.TTL != -1 && time.Now().Unix() >= storageValue.TTL {
		c.Del(key)
		return []byte{}, nil
	}
	return storageValue.Value, nil
}
