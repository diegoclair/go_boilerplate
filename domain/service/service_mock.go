package service

import (
	"testing"

	"github.com/diegoclair/go_boilerplate/infra/config"
	"github.com/diegoclair/go_boilerplate/infra/logger"
	"github.com/diegoclair/go_boilerplate/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type allMocks struct {
	mockDataManager *mocks.MockDataManager

	mockAuthRepo    *mocks.MockAuthRepo
	mockAccountRepo *mocks.MockAccountRepo

	mockCacheManager *mocks.MockCacheManager
	mockCrypto       *mocks.MockCrypto
}

func newServiceTestMock(t *testing.T) (m allMocks, svc *service, ctrl *gomock.Controller) {

	cfg, err := config.GetConfigEnvironment("../../" + config.ProfileTest)
	require.NoError(t, err)

	ctrl = gomock.NewController(t)
	log := logger.NewNoop()
	dm := mocks.NewMockDataManager(ctrl)
	accountRepo := mocks.NewMockAccountRepo(ctrl)
	cm := mocks.NewMockCacheManager(ctrl)
	crypto := mocks.NewMockCrypto(ctrl)
	authRepo := mocks.NewMockAuthRepo(ctrl)

	m = allMocks{
		mockDataManager:  dm,
		mockAccountRepo:  accountRepo,
		mockCacheManager: cm,
		mockAuthRepo:     authRepo,
		mockCrypto:       crypto,
	}

	dm.EXPECT().Account().Return(accountRepo).AnyTimes()
	dm.EXPECT().Auth().Return(authRepo).AnyTimes()

	svc = newService(dm, cfg, cm, crypto, log)

	return
}
