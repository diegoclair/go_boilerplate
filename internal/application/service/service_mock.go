package service

import (
	"testing"
	"time"

	"github.com/diegoclair/go_boilerplate/infra/configmock"
	"github.com/diegoclair/go_boilerplate/mocks"
	"github.com/diegoclair/go_utils/logger"
	"github.com/diegoclair/go_utils/validator"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type allMocks struct {
	mockDataManager *mocks.MockDataManager

	mockAuthRepo    *mocks.MockAuthRepo
	mockAccountRepo *mocks.MockAccountRepo

	mockCacheManager *mocks.MockCacheManager
	mockCrypto       *mocks.MockCrypto
	mockValidator    validator.Validator
	mockLogger       logger.Logger

	mockAccountSvc *mocks.MockAccountApp

	mockDomain *mocks.MockInfrastructure
}

func newServiceTestMock(t *testing.T) (m allMocks, ctrl *gomock.Controller) {
	t.Helper()
	cfg := configmock.New()

	ctrl = gomock.NewController(t)

	dm := mocks.NewMockDataManager(ctrl)

	accountRepo := mocks.NewMockAccountRepo(ctrl)
	dm.EXPECT().Account().Return(accountRepo).AnyTimes()

	authRepo := mocks.NewMockAuthRepo(ctrl)
	dm.EXPECT().Auth().Return(authRepo).AnyTimes()

	cm := cfg.GetCacheManager(ctrl)
	crypto := cfg.GetCrypto(ctrl)
	log := cfg.GetLogger()
	v := cfg.GetValidator(t)

	accountSvc := mocks.NewMockAccountApp(ctrl)

	domainMock := mocks.NewMockInfrastructure(ctrl)
	domainMock.EXPECT().DataManager().Return(dm).AnyTimes()
	domainMock.EXPECT().Logger().Return(log).AnyTimes()
	domainMock.EXPECT().CacheManager().Return(cm).AnyTimes()
	domainMock.EXPECT().Crypto().Return(crypto).AnyTimes()
	domainMock.EXPECT().Validator().Return(v).AnyTimes()

	m = allMocks{
		mockDataManager:  dm,
		mockAccountRepo:  accountRepo,
		mockCacheManager: cm,
		mockAuthRepo:     authRepo,
		mockCrypto:       crypto,
		mockAccountSvc:   accountSvc,
		mockDomain:       domainMock,
		mockValidator:    v,
		mockLogger:       log,
	}

	// validate func New
	s, err := New(domainMock, time.Minute)
	require.NoError(t, err)
	require.NotNil(t, s)

	return
}
