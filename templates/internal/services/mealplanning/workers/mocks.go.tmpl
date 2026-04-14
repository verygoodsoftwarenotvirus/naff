package workers

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type (
	MockWorker struct {
		mock.Mock
	}

	MockWorkerCounter struct {
		mock.Mock
	}
)

// Work satisfies the Worker interface.
func (m *MockWorker) Work(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

// Work satisfies the WorkerCounter interface.
func (m *MockWorkerCounter) Work(ctx context.Context) (int64, error) {
	returnValues := m.Called(ctx)
	return returnValues.Get(0).(int64), returnValues.Error(1)
}
