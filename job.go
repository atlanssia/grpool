package grpool

type Job interface {
	Process(interface{})
}