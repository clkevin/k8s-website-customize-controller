package controller

import (
	"k8s.io/client-go/util/workqueue"
	"k8s.io/apimachinery/pkg/util/wait"
	"time"
)

// Controller is the controller implementation for Student resources
type WebsiteController struct {
	workqueue workqueue.RateLimitingInterface//保存资源处理事件
}

// 判断Controller
func NewController() *WebsiteController {
	controller := &WebsiteController{
		workqueue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Websites"),
	}
	return controller
}

//多线程循环监控controller中队列
func (c *WebsiteController) Run(threadiness int, stopCh <-chan struct{}) {
	defer c.workqueue.ShutDown()

	for i := 0; i < threadiness; i++ { //启动两个线程
		go wait.Until(c.runWorker, time.Second*10, stopCh) //在上一个执行结束的情况下，每10秒执行一次runWorker
	}
    <- stopCh
}

//监控controller队列并执行业务逻辑
func (c *WebsiteController) runWorker() {
	println("runWorker")
	obj, shutdown := c.workqueue.Get()
	println("get obj")
	if shutdown {
		println("workqueue is shutdown")
		return
	}

	defer c.workqueue.Done(obj)

	action := obj.(EventAction)

	switch action.action {
	case "OnAdd":
		println("RunWorker:OnAdd")
	case "OnDelete":
		println("RunWorker:OnDelete")
	default:
		println("RunWorker:Default")
	}
}
