package cache

// #include "rocksdb/c.h"
// #cgo CFLAGS: -I/rocksdb/include
// #cgo LDFLAGS: -L/rocksdb -lrocksdb -lz -lpthread -lsnappy -lstdc++ -lm -O3
import "C"
import (
	"dijan/utils"
	"fmt"
	"os/exec"
	"runtime"
	"time"
)

func newRocksdbCache(ttl int) *rocksdbCache {
	options := C.rocksdb_options_create()
	C.rocksdb_options_increase_parallelism(options, C.int(runtime.NumCPU()))
	C.rocksdb_options_set_create_if_missing(options, 1)
	var e *C.char
Conn:
	db := C.rocksdb_open_with_ttl(options, C.CString(utils.GlobalSystemConfig.Rocksdb.RootDir), C.int(ttl), &e)
	if e != nil {
		fmt.Println(fmt.Sprintf("[!] rocksdb资源连接失败: %s。尝试删除lock文件，将于5s后尝试重连", C.GoString(e)))
		exec.Command("rm", "-rf", "/mnt/rocksdb/LOCK").Run()
		time.Sleep(time.Second * 5)
		goto Conn
	}
	C.rocksdb_options_destroy(options)
	c := make(chan *pair, utils.GlobalSystemConfig.Rocksdb.BatchSize * 10)
	wo := C.rocksdb_writeoptions_create()
	go write_func(db, c, wo)
	return &rocksdbCache{db, C.rocksdb_readoptions_create(), wo, e, c}
}
