package service

import (
	"testing"

	"github.com/diegoclair/go_boilerplate/infra/config"
	"github.com/diegoclair/go_boilerplate/mocks"
	"github.com/diegoclair/go_utils/logger"
	"github.com/diegoclair/go_utils/validator"
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
	dm.EXPECT().Account().Return(accountRepo).AnyTimes()

	authRepo := mocks.NewMockAuthRepo(ctrl)
	dm.EXPECT().Auth().Return(authRepo).AnyTimes()

	cm := mocks.NewMockCacheManager(ctrl)
	crypto := mocks.NewMockCrypto(ctrl)

	m = allMocks{
		mockDataManager:  dm,
		mockAccountRepo:  accountRepo,
		mockCacheManager: cm,
		mockAuthRepo:     authRepo,
		mockCrypto:       crypto,
	}

	v, err := validator.NewValidator()
	require.NoError(t, err)

	svc = &service{}
	WithDataManager(dm)(svc)
	WithConfig(cfg)(svc)
	WithCacheManager(cm)(svc)
	WithLogger(log)(svc)
	WithCrypto(crypto)(svc)
	WithValidator(v)(svc)

	// validate func New
	s, err := New(WithDataManager(dm))
	require.NoError(t, err)
	require.NotNil(t, s)

	return
}
