package testutil

import (
	"github.com/diegoclair/go_boilerplate/mocks"
	"github.com/golang/mock/gomock"
)

type Mocks struct {
	AccountServiceMock  *mocks.MockAccountService
	AuthServiceMock     *mocks.MockAuthService
	TransferServiceMock *mocks.MockTransferService
}

func NewServiceManagerTest(ctrl *gomock.Controller) Mocks {
	return Mocks{
		AccountServiceMock:  mocks.NewMockAccountService(ctrl),
		AuthServiceMock:     mocks.NewMockAuthService(ctrl),
		TransferServiceMock: mocks.NewMockTransferService(ctrl),
	}
}
