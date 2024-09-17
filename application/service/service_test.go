package service

import (
	"testing"
	"time"

	"github.com/diegoclair/go_boilerplate/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("Valid infrastructure", func(t *testing.T) {
		m, ctrl := newServiceTestMock(t)
		defer ctrl.Finish()

		apps, err := New(m.mockDomain, time.Hour)
		assert.NoError(t, err)
		assert.NotNil(t, apps)
	})

	t.Run("Invalid infrastructure", func(t *testing.T) {
		m, ctrl := newServiceTestMock(t)
		m.mockDomain = mocks.NewMockInfrastructure(ctrl)
		m.mockDomain.EXPECT().Logger().Return(nil)
		defer ctrl.Finish()

		apps, err := New(m.mockDomain, time.Hour)
		assert.Error(t, err)
		assert.Nil(t, apps)
	})
}

func TestValidateInfrastructure(t *testing.T) {
	m, ctrl := newServiceTestMock(t)
	defer ctrl.Finish()
	m.mockDomain = mocks.NewMockInfrastructure(ctrl)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		setup   func(allMocks)
		wantErr string
	}{
		{
			name: "Valid infrastructure",
			setup: func(m allMocks) {
				m.mockDomain.EXPECT().Logger().Return(m.mockLogger)
				m.mockDomain.EXPECT().DataManager().Return(m.mockDataManager)
				m.mockDomain.EXPECT().CacheManager().Return(m.mockCacheManager)
				m.mockDomain.EXPECT().Crypto().Return(m.mockCrypto)
				m.mockDomain.EXPECT().Validator().Return(m.mockValidator)
			},
			wantErr: "",
		},
		{
			name: "Missing logger",
			setup: func(m allMocks) {
				m.mockDomain.EXPECT().Logger().Return(nil)
			},
			wantErr: "logger is required",
		},
		{
			name: "Missing data manager",
			setup: func(m allMocks) {
				m.mockDomain.EXPECT().Logger().Return(m.mockLogger)
				m.mockDomain.EXPECT().DataManager().Return(nil)
			},
			wantErr: "data manager is required",
		},
		{
			name: "Missing cache manager",
			setup: func(m allMocks) {
				m.mockDomain.EXPECT().Logger().Return(m.mockLogger)
				m.mockDomain.EXPECT().DataManager().Return(m.mockDataManager)
				m.mockDomain.EXPECT().CacheManager().Return(nil)
			},
			wantErr: "cache manager is required",
		},
		{
			name: "Missing crypto",
			setup: func(m allMocks) {
				m.mockDomain.EXPECT().Logger().Return(m.mockLogger)
				m.mockDomain.EXPECT().DataManager().Return(m.mockDataManager)
				m.mockDomain.EXPECT().CacheManager().Return(m.mockCacheManager)
				m.mockDomain.EXPECT().Crypto().Return(nil)
			},
			wantErr: "crypto is required",
		},
		{
			name: "Missing validator",
			setup: func(m allMocks) {
				m.mockDomain.EXPECT().Logger().Return(m.mockLogger)
				m.mockDomain.EXPECT().DataManager().Return(m.mockDataManager)
				m.mockDomain.EXPECT().CacheManager().Return(m.mockCacheManager)
				m.mockDomain.EXPECT().Crypto().Return(m.mockCrypto)
				m.mockDomain.EXPECT().Validator().Return(nil)
			},
			wantErr: "validator is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.setup != nil {
				tt.setup(m)
			}

			err := validateInfrastructure(m.mockDomain)

			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}
