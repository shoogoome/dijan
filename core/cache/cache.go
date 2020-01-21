package cache

type Cache interface {
	Set(string, []byte, ...int) error
	Get(string) ([]byte, error)
	Del(string) error
	GetStat() Stat
	NewScanner() Scanner
}
