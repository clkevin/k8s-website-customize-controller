package controller

import (
	"k8s-website-customize-controller/pkg/apis/kevin/v1"
)

// interface cache.ResourceEventHandlerFuncs
type EventHandler struct {
	Controller WebsiteController
}

//资源对象添加时执行该方法
func (e *EventHandler) OnAdd(obj interface{}) {
	println("OnAdd")
	event := EventAction{action: "OnAdd", obj: obj}
	e.Controller.workqueue.Add(event)
}
//资源对象修改时执行该方法
func (e *EventHandler) OnUpdate(oldObj, newObj interface{}) {
	println("OnUpdate")
	old := oldObj.(*v1.Website)
	new := newObj.(*v1.Website)
	if old.ResourceVersion == new.ResourceVersion {
		//版本一致，就表示没有实际更新的操作，立即返回
		return
	}

	event := EventAction{action: "OnUpdate", obj: newObj}
	e.Controller.workqueue.Add(event)
}
//资源对象删除时执行该方法
func (e *EventHandler) OnDelete(obj interface{}) {
	println("OnDelete")
	event := EventAction{action: "OnDelete", obj: obj}
	e.Controller.workqueue.Add(event)
}
