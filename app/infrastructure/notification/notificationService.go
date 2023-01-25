package notif

import (
	"usermanager/app/infrastructure/rabbit"

	"github.com/google/uuid"
)

type NotificationService interface {
	NotifyAboutUserChange(userId uuid.UUID)
}

type notificationService struct {
	rmq *rabbit.RMQ
}

func NewNotificationService(rmqPublisher *rabbit.RMQ) *notificationService {
	return &notificationService{
		rmq: rmqPublisher,
	}
}

// This function will push user id to channel where on other side will
// be a running goroutine that publish messages to a rabbit queue.
// On the other side of the queue are subscribed listeners (services)
// that are interested about updated users.
func (n *notificationService) NotifyAboutUserChange(userId uuid.UUID) {
	n.rmq.PublishChannel <- userId
}
