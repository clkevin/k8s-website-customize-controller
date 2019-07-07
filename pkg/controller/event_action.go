package controller

//事件行为，如OnAdd OnDelete OnUpdate
type EventAction struct {
	action string
	obj    interface{}
}
