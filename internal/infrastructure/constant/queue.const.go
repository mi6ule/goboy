package constants

const (
	DefaultQueue     = "default"
	FirstEmailQueue  = "first-email"
	SecondEmailQueue = "second-email"
	ImageResizeQueue = "image"
)

type Queue struct {
	Name     string
	Priority int
}

var Queues []Queue = []Queue{
	{Name: DefaultQueue, Priority: 1},
	{Name: FirstEmailQueue, Priority: 2},
	{Name: SecondEmailQueue, Priority: 2},
	{Name: ImageResizeQueue, Priority: 2},
}
