package cache

var Conn Cache

func New() Cache {
	c := newRocksdbCache()
	Conn = c
	return c
}
