package service

import (
	"testing"

	infraMocks "github.com/diegoclair/go_boilerplate/infra/mocks"
	"github.com/stretchr/testify/assert"
)

func TestValidateInfrastructure(t *testing.T) {
	m, ctrl := newServiceTestMock(t)
	defer ctrl.Finish()
	m.mockInfra = infraMocks.NewMockInfrastructure(ctrl)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		setup   func(allMocks)
		wantErr string
	}{
		{
			name: "Valid infrastructure",
			setup: func(m allMocks) {
				m.mockInfra.EXPECT().Logger().Return(m.mockLogger)
				m.mockInfra.EXPECT().DataManager().Return(m.mockDataManager)
				m.mockInfra.EXPECT().CacheManager().Return(m.mockCacheManager)
				m.mockInfra.EXPECT().Crypto().Return(m.mockCrypto)
				m.mockInfra.EXPECT().Validator().Return(m.mockValidator)
			},
			wantErr: "",
		},
		{
			name: "Missing logger",
			setup: func(m allMocks) {
				m.mockInfra.EXPECT().Logger().Return(nil)
			},
			wantErr: "logger is required",
		},
		{
			name: "Missing data manager",
			setup: func(m allMocks) {
				m.mockInfra.EXPECT().Logger().Return(m.mockLogger)
				m.mockInfra.EXPECT().DataManager().Return(nil)
			},
			wantErr: "data manager is required",
		},
		{
			name: "Missing cache manager",
			setup: func(m allMocks) {
				m.mockInfra.EXPECT().Logger().Return(m.mockLogger)
				m.mockInfra.EXPECT().DataManager().Return(m.mockDataManager)
				m.mockInfra.EXPECT().CacheManager().Return(nil)
			},
			wantErr: "cache manager is required",
		},
		{
			name: "Missing crypto",
			setup: func(m allMocks) {
				m.mockInfra.EXPECT().Logger().Return(m.mockLogger)
				m.mockInfra.EXPECT().DataManager().Return(m.mockDataManager)
				m.mockInfra.EXPECT().CacheManager().Return(m.mockCacheManager)
				m.mockInfra.EXPECT().Crypto().Return(nil)
			},
			wantErr: "crypto is required",
		},
		{
			name: "Missing validator",
			setup: func(m allMocks) {
				m.mockInfra.EXPECT().Logger().Return(m.mockLogger)
				m.mockInfra.EXPECT().DataManager().Return(m.mockDataManager)
				m.mockInfra.EXPECT().CacheManager().Return(m.mockCacheManager)
				m.mockInfra.EXPECT().Crypto().Return(m.mockCrypto)
				m.mockInfra.EXPECT().Validator().Return(nil)
			},
			wantErr: "validator is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.setup != nil {
				tt.setup(m)
			}

			err := validateInfrastructure(m.mockInfra)

			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}
