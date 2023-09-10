package storage

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type LinkStorageMock struct {
	mock.Mock
}

func (m *LinkStorageMock) Create(ctx context.Context, link Link) error {
	args := m.Called(link)
	return args.Error(0)
}

func (m *LinkStorageMock) GetLinkByHash(ctx context.Context, hash string) (Link, error) {
	args := m.Called(hash)
	return args.Get(0).(Link), args.Error(1)
}

func (m *LinkStorageMock) CreateTable(ctx context.Context) error {
	args := m.Called()
	return args.Error(0)
}
