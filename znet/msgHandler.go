package znet

import (
	"fmt"
	"log"
	"strconv"

	"github.com/fatih/color"
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
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
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

// 启动一个Worker工作池(开启工作池的动作只能发生一次, 一个zinx框架只能有一个工作池)
func (m *MsgHandle) StartWorkerPool() {
	// 根据workerPoolSize开启Worker, 每个Worker用一个go承载
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		// 一个worker被启动
		// 1. 给当前的worker对应的channel消息队列, 开辟空间 第0个worker就用第0个channel
		m.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskSize)
		// 2. 启动当前的worker, 阻塞等待消息从channel传递进来
		go m.StartOneWorker(i, m.TaskQueue[i])
	}
}

// 启动一个Worker工作流程
func (m *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	color.Cyan("Worker ID = %d, is started...", workerID)

	// 不断的阻塞等待对应消息队列的消息
	for {
		select {
		// 如果有消息过来, 出列的就是一个客户端的Request, 执行当前Request所绑定的业务
		case request := <-taskQueue:
			m.DoMsgHandler(request)
		}
	}
}

// 将消息交给TaskQueue, 由worker进行处理
func (m *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	// 1. 将消息平均分配给不同的worker
	// 根据客户端建立的ConnID来进行分配
	workerID := request.GetConnection().GetConnID() % m.WorkerPoolSize
	color.Cyan("Add ConnID = %d, request MsgID = %d, to WorkerID = %d\n",
		request.GetConnection().GetConnID(),
		request.GetMsgID(),
		workerID,
	)
	// 2. 将消息发送给对应的worker的TaskQueue即可
	m.TaskQueue[workerID] <- request
}
