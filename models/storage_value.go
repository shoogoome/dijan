package models

type StorageValue struct {

	// 实际存储值
	Value []byte

	// 过期时间
	TTL int64
}
