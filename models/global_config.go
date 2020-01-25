package models

type GlobalConfig struct {

	// 系统配置
	Server server `json:"server" yaml:"server"`

	// 集群配置
	Cluster cluster `json:"cluster" yaml:"cluster"`

	// 缓存数据库配置
	Rocksdb rocksdb `json:"rocksdb" yaml:"rocksdb"`
}


type cluster struct {
	// 集群节点数（包含虚拟节点）
	CircleNumberOfNode int `json:"circle_number_of_node" yaml:"circle_number_of_node"`
}

type rocksdb struct {
	// 存储数据根目录
	RootDir string `json:"root_dir" yaml:"root_dir"`
	// 异步处理最大同时处理数
	AsynchronousNumber int `json:"asynchronous_number" yaml:"asynchronous_number"`
	// 批量处理最大延时
	BatchHandleTime int `json:"batch_handle_time" yaml:"batch_handle_time"`
	// 批处理数量
	BatchSize int `json:"batch_size" yaml:"batch_size"`

}

type server struct {
	// ui管理界面登录密码
	SystemPassword string `json:"system_password" yaml:"system_password"`

	// api交互token
	SystemApiToken string `json:"system_api_token" yaml:"system_api_token"`

	// tcp监听端口
	TcpListenPort string `json:"tcp_listen_port" yaml:"tcp_listen_port"`

	// http监听端口
	HttpListenPort string `json:"http_listen_port" yaml:"http_listen_port"`

	// 主机名
	HostName string `json:"host_name"`

	// 本机地址
	Address string `json:"address" yaml:"address"`

	// 无头服务名
	HeadlessService string `json:"headless_service"`

	// 集群地址
	ClusterAddress string `json:"cluster_address" yaml:"cluster_address"`

	// cookie过期时间
	CookieExpires int64 `json:"cookie_expires" yaml:"cookie_expires"`

	// 系统加密盐
	HashSalt string `json:"hash_salt" yaml:"hash_salt"`
}