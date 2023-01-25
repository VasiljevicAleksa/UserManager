package mocks

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type NotificationServiceMock struct {
	mock.Mock
}

// this mock will only record that the method was called
func (n *NotificationServiceMock) NotifyAboutUserChange(userId uuid.UUID) {
	_ = n.Called(userId)
}
