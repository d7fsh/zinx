package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"zinx_demo/ziface"
)

/*
存储一些有关Zinx框架的全局参数, 供其他模块使用
一些参数是可以通过zinx.json由用户进行配置的
*/
type GlobalObj struct {
	// Server配置
	Name      string         `json:"name"` // 当前服务器的名称
	TcpServer ziface.IServer // 当前zinx全局的server对象
	Host      string         `json:"host"` // 当前服务器主机监听的IP
	TcpPort   int            `json:"port"` // 当前服务器主机监听的端口号
	// Zinx配置
	Version           string // 当前zinx的版本号
	MaxConn           int    `json:"max_conn"` // 当前服务器主机允许最大连接数
	MaxPackageSize    uint32 // 当前zinx框架数据包的最大值
	WorkerPoolSize    uint32 `json:"worker_pool_size"`     // 当前业务工作WorkerPool的goroutine数量
	MaxWorkerTaskSize uint32 `json:"max_worker_task_size"` // zinx框架允许用户最多开辟多少个Worker(限定条件)
}

/*
	定义一个全局的对外Globalobj对象
*/
var GlobalObject *GlobalObj

// 从zinx.json加载用户自定义的参数
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		log.Panic(err)
	}
	// 将json文件数据解析到struct中
	err = json.Unmarshal(data, GlobalObject)
	if err != nil {
		log.Panic(err)
	}
}

// 提供一个init方法, 初始化当前的GlobalObject对象
func init() {
	// 如果配置文件没有加载, 默认值
	GlobalObject = &GlobalObj{
		Name:              "ZinxServerApp",
		TcpServer:         nil,
		Host:              "0.0.0.0",
		TcpPort:           8999,
		Version:           "V0.4",
		MaxConn:           100,
		MaxPackageSize:    4096,
		WorkerPoolSize:    10,   // Worker工作池的队列的个数
		MaxWorkerTaskSize: 1000, // 每个Worker对应的消息队列的任务的数量最大值
	}

	// 应该尝试从conf/zinx.json中加载用户自定义配置
	GlobalObject.Reload()
}
