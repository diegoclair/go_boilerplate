package service

import (
	"testing"
	"time"

	"github.com/diegoclair/go_boilerplate/infra/configmock"
	infraMocks "github.com/diegoclair/go_boilerplate/infra/mocks"
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

	mockCacheManager *infraMocks.MockCacheManager
	mockCrypto       *infraMocks.MockCrypto
	mockValidator    validator.Validator
	mockLogger       logger.Logger

	mockAccountSvc *mocks.MockAccountApp

	mockInfra *infraMocks.MockInfrastructure
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

	infraMock := infraMocks.NewMockInfrastructure(ctrl)
	infraMock.EXPECT().DataManager().Return(dm).AnyTimes()
	infraMock.EXPECT().Logger().Return(log).AnyTimes()
	infraMock.EXPECT().CacheManager().Return(cm).AnyTimes()
	infraMock.EXPECT().Crypto().Return(crypto).AnyTimes()
	infraMock.EXPECT().Validator().Return(v).AnyTimes()

	m = allMocks{
		mockDataManager:  dm,
		mockAccountRepo:  accountRepo,
		mockCacheManager: cm,
		mockAuthRepo:     authRepo,
		mockCrypto:       crypto,
		mockAccountSvc:   accountSvc,
		mockInfra:        infraMock,
		mockValidator:    v,
		mockLogger:       log,
	}

	// validate func New
	s, err := New(infraMock, time.Minute)
	require.NoError(t, err)
	require.NotNil(t, s)

	return
}
