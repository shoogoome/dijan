# dijan分布式缓存系统

## 简介
使用rocksdb作为缓存主体，使用异步存储、批量操作、pipe、tcp手段提高缓存性能；  
针对kubernetes的StatefulSet进行适配，自动进行集群部署；  
系统具有功能完备的ui界面，有对数据库对curd、清空、节点再平衡等功能。

## 使用说明
