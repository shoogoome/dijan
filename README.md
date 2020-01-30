# dijan分布式缓存系统

## 简介
采用rocksdb作为缓存主体，使用异步存储、批量操作、pipe、tcp手段提高缓存性能；  
一致性哈希、gossip协议为分布式提供支持，针对kubernetes的StatefulSet进行适配，自动进行集群部署、宕机重连接；  
缓存使用ttl规定过期时间，在获取时清理或定时清理  
系统具有功能完备的ui界面，有对数据库对curd、清空、节点再平衡等功能。

## 使用说明

**前置条件: 系统配置docker环境**

1. 拉取docker镜像
```
docker pull registry.cn-shenzhen.aliyuncs.com/shoogoome/dijan:v1.0
```

2. 替换默认系统配置(/etc/dijan/config.yaml)
```
#  /etc/dijan/config.yaml
#   --------------------
#  |   Dijan配置文件     |
#  |  依照需求修改配置     |
#  |  请勿改动配置结构     |
#   --------------------


# 系统配置
server:
    # ui管理界面登录密码
    system_password: "12345678"
    # api交互token
    system_api_token: "12345678"
    # tcp监听端口
    tcp_listen_port: ":2375"
    # http监听端口，也是api接口端口和ui管理界面端口
    http_listen_port: ":8080"
    # cookie过期时间
    cookie_expires: 86400
    # 系统加密盐（建议修改随机字符串）
    hash_salt: "edfghjuy6tredfghtrerfgh"

# 集群配置
cluster:
    # 集群节点数（包含虚拟节点数）
    circle_number_of_node: 256

# 缓存数据库配置
rocksdb:
    # 数据持久化根节点
    root_dir: /mnt/rocksdb
    # 异步处理最大同时处理数（功能可替代客户端的连接池）
    asynchronous_number: 5000
    # 批处理功能:
    # 若batch_handle_time秒内累计set请求达到batch_size则进行set操作并重置计数器
    # 若到batch_handle_time秒则将执行累计的set操作
    # 应搭配好这两个配置项。在默认配置下，即便系统宕机也只会丢失1s内且不多于100条数据（作为缓存系统，允许数据丢失）
    batch_handle_time: 1
    batch_size: 100
```

3. (可选)使用Kubernetes的StatefulSet，添加环境变量'HEADLESS_SERVICE'为无头服务名称。系统将自动完成集群。

4. 访问可通过sdk、api、ui界面  
* sdk
    * GO: ```https://github.com/shoogoome/godijan```
* api(请求头添加'token'参数对应配置文件'server.system_api_token'参数)
    * /api/storage/get/:key  获取缓存
        * method: get
    * /api/storage/set       缓存
        * method: post
        * content-type: application/json
        * json: {"key": 缓存键, "value": 缓存值, "ttl": 缓存有效期(-1不清除)
    * /api/storage/delete    删除缓存
        * method: delete
        * content-type: application/json
        * josn: {"key": 缓存键}
    * /api/node              获取节点信息
        * method: get
    * /api/rebalance         节点再平衡
        * method: get
* ui界面端口同配置文件server.http_listen_port配置

## 缓存清理策略

* 在获取时检查ttl，若过期则删除缓存
* 在每天凌晨启动清理缓存，实时监控cpu使用情况。超过20%则停止清理

## PS
1. 当集群节点数发生改变时才有必要进行节点再平衡，但这不是必须的操作。集群缓存较大的时候节点在平衡速度可能偏慢，尽量在请求不频繁时使用。
2. 与rocksdb交互的核心代码学习自胡世杰老师的《分布式缓存--原理、架构及GO语言实现》一书
