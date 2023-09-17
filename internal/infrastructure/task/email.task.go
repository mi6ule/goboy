package task

import (
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/queue"
)

type EmailDeliveryPayload struct {
	UserID     int
	TemplateID string
}

func NewEmailDeliveryTask(userID int, tmplID string) (*queue.AsynqTask, error) {
	// do something here
	// return task
	return &queue.AsynqTask{TypeName: EmailDeliveryTask, Payload: EmailDeliveryPayload{UserID: userID, TemplateID: tmplID}}, nil
}
