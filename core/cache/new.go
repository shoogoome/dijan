package cache

import "log"

var Conn Cache

func New(typ string, ttl int) Cache {
	var c Cache
	if typ == "rocksdb" {
		c = newRocksdbCache(ttl)
	}
	if c == nil {
		panic("unknown cache type " + typ)
	}
	log.Println(typ, "ready to serve")
	Conn = c
	return c
}
