package cache

// #include "rocksdb/c.h"
// #cgo CFLAGS: -I/rocksdb/include
// #cgo LDFLAGS: -L/rocksdb -lrocksdb -lz -lpthread -lsnappy -lstdc++ -lm -O3
import "C"
import (
	"dijan/utils"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func newRocksdbCache() *rocksdbCache {
	options := C.rocksdb_options_create()
	C.rocksdb_options_increase_parallelism(options, C.int(runtime.NumCPU()))
	C.rocksdb_options_set_create_if_missing(options, 1)
	var e *C.char

	db := C.rocksdb_open(options, C.CString(utils.GlobalSystemConfig.Rocksdb.RootDir), &e)
	if e != nil {
		fmt.Println(fmt.Sprintf("[!] rocksdb资源连接失败: %s。已删除lock文件，将尝试重启程序", C.GoString(e)))
		exec.Command("rm", "-rf", fmt.Sprintf("%s/LOCK", utils.GlobalSystemConfig.Rocksdb.RootDir)).Run()
		exec.Command("kill", "-9", fmt.Sprintf("%d", os.Getpid())).Run()
	}
	C.rocksdb_options_destroy(options)
	wo := C.rocksdb_writeoptions_create()
	var c chan *pair
	// 开启批量写入
	if utils.GlobalSystemConfig.Rocksdb.BatchHandleTime > 0 {
		c = make(chan *pair, utils.GlobalSystemConfig.Rocksdb.BatchSize * 10)
		go write_func(db, c, wo)
	}
	return &rocksdbCache{db, C.rocksdb_readoptions_create(), wo, e, c}
}
