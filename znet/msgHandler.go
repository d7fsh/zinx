package znet

import (
	"fmt"
	"log"
	"strconv"

	"zinx_demo/utils"
	"zinx_demo/ziface"
)

/*
消息处理模块的实现
*/
type MsgHandle struct {
	// 存放每个MsgID所对应的处理方法
	Apis map[uint32]ziface.IRouter
	// 业务工作WorkerPool的worker数量
	TaskQueue []chan ziface.IRequest
	// 负责Worker读取任务的消息队列
	WorkerPoolSize uint32
}

// 初始化/创建MsgHandle的方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize, // 从全局配置中获取
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskSize),
	}
}

// 调度/执行对应的Router消息处理方法
func (m *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	// 1. 从request中找到msgID
	handler, ok := m.Apis[request.GetMsgID()]
	if !ok {
		fmt.Printf("api msgID = %d, is NOT FOUND! Need register", request.GetMsgID())
	}

	// 2. 根据msgID调度对应的router业务即可
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 为消息添加具体的处理逻辑
func (m *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	// 1. 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := m.Apis[msgID]; ok {
		// id已经注册了
		log.Panic("repeat api, msgID =" + strconv.Itoa(int(msgID)))
	}
	// 2. 添加msg与API的绑定关系
	m.Apis[msgID] = router
	fmt.Printf("Add api MsgID = %d, success", msgID)
}
