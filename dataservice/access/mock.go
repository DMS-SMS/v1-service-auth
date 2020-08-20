package access

import (
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock *mock.Mock
}

func NewMock(mock *mock.Mock) Mock {
	return Mock{
		mock: mock,
	}
}
