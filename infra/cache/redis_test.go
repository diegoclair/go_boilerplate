package cache

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/diegoclair/go_boilerplate/infra/config"
	"github.com/diegoclair/go_boilerplate/mocks"
	"github.com/diegoclair/go_utils/logger"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func getRedisCacheMock(ctrl *gomock.Controller) (*redisCache, *mocks.MockIRedisCache) {
	m := mocks.NewMockIRedisCache(ctrl)
	return &redisCache{
		redis: m,
	}, m
}

func TestNewRedisCache(t *testing.T) {
	ctx := context.Background()

	//this config will be the same from main_test.go because of singleton pattern from GetConfigEnvironment
	cfg, err := config.GetConfigEnvironment(config.ProfileTest)
	require.NoError(t, err)

	client, err := newRedisCache(ctx, cfg, logger.NewNoop())
	require.NoError(t, err)
	require.NotNil(t, client)
}

func TestRedisCache_GetItem(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedRedis, redisMock := getRedisCacheMock(ctrl)

	type args struct {
		key   string
		cache *redisCache
	}
	tests := []struct {
		name      string
		args      args
		setupTest func(args args, m *mocks.MockIRedisCache)
		want      []byte
		wantErr   error
	}{
		{
			name: "Cache hit",
			args: args{key: "getItem_key", cache: testRedis.(*redisCache)},
			setupTest: func(args args, _ *mocks.MockIRedisCache) {
				err := testRedis.SetItem(ctx, args.key, []byte("test_value"))
				require.NoError(t, err)
			},
			want:    []byte("test_value"),
			wantErr: nil,
		},
		{
			name:    "Cache miss",
			args:    args{key: "getItem_other_key", cache: testRedis.(*redisCache)},
			want:    nil,
			wantErr: ErrCacheMiss,
		},
		{
			name: "Return error when get item from cache fails",
			args: args{key: "test_key", cache: mockedRedis},
			setupTest: func(args args, m *mocks.MockIRedisCache) {
				m.EXPECT().Get(gomock.Any(), args.key).Return(redis.NewStringResult("", errors.New("some error")))
			},
			want:    nil,
			wantErr: errors.New("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupTest != nil {
				tt.setupTest(tt.args, redisMock)
			}

			got, err := tt.args.cache.GetItem(ctx, tt.args.key)
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestRedisCache_SetItem(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedRedis, redisMock := getRedisCacheMock(ctrl)

	type args struct {
		key   string
		data  []byte
		cache *redisCache
	}
	tests := []struct {
		name       string
		args       args
		setupCache func(args args, m *mocks.MockIRedisCache)
		wantErr    error
	}{
		{
			name:    "Success",
			args:    args{key: "setItem_test_key", data: []byte("test_value"), cache: testRedis.(*redisCache)},
			wantErr: nil,
		},
		{
			name: "Error",
			args: args{key: "test_key", data: []byte("test_value"), cache: mockedRedis},
			setupCache: func(args args, m *mocks.MockIRedisCache) {
				cmd := redis.NewStatusCmd(ctx, args.key, "OK")
				cmd.SetErr(errors.New("some error"))
				m.EXPECT().Set(gomock.Any(), args.key, args.data, gomock.Any()).
					Return(cmd)
			},
			wantErr: errors.New("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupCache != nil {
				tt.setupCache(tt.args, redisMock)
			}

			err := tt.args.cache.SetItem(ctx, tt.args.key, tt.args.data)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestRedisCache_GetInt(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedRedis, redisMock := getRedisCacheMock(ctrl)

	type args struct {
		key string
	}
	tests := []struct {
		name       string
		args       args
		setupCache func(args args, m *mocks.MockIRedisCache)
		want       int64
		wantErr    error
	}{
		{
			name: "Cache hit",
			args: args{key: "getItem_test_key"},
			setupCache: func(args args, m *mocks.MockIRedisCache) {
				m.EXPECT().Get(gomock.Any(), args.key).Return(redis.NewStringResult("123", nil))
			},
			want:    123,
			wantErr: nil,
		},
		{
			name: "Cache miss",
			args: args{key: "getItem_other_key"},
			setupCache: func(args args, m *mocks.MockIRedisCache) {
				m.EXPECT().Get(gomock.Any(), args.key).Return(redis.NewStringResult("", redis.Nil))
			},
			want:    0,
			wantErr: ErrCacheMiss,
		},
		{
			name: "Error",
			args: args{key: "test_key"},
			setupCache: func(args args, m *mocks.MockIRedisCache) {
				m.EXPECT().Get(gomock.Any(), args.key).Return(redis.NewStringResult("", errors.New("some error")))
			},
			want:    0,
			wantErr: errors.New("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupCache != nil {
				tt.setupCache(tt.args, redisMock)
			}

			got, err := mockedRedis.GetInt(ctx, tt.args.key)
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestRedisCache_GetString(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedRedis, redisMock := getRedisCacheMock(ctrl)

	type args struct {
		key   string
		cache *redisCache
	}
	tests := []struct {
		name       string
		args       args
		setupCache func(args args, m *mocks.MockIRedisCache)
		want       string
		wantErr    error
	}{
		{
			name: "Cache hit",
			args: args{key: "getString_test_key", cache: testRedis.(*redisCache)},
			setupCache: func(args args, _ *mocks.MockIRedisCache) {
				err := testRedis.SetItem(ctx, args.key, []byte("test_value"))
				require.NoError(t, err)
			},
			want:    "test_value",
			wantErr: nil,
		},
		{
			name:    "Cache miss",
			args:    args{key: "getString_other_key", cache: testRedis.(*redisCache)},
			want:    "",
			wantErr: ErrCacheMiss,
		},
		{
			name: "Error",
			args: args{key: "test_key", cache: mockedRedis},
			setupCache: func(args args, m *mocks.MockIRedisCache) {
				m.EXPECT().Get(gomock.Any(), args.key).Return(redis.NewStringResult("", errors.New("some error")))
			},
			want:    "",
			wantErr: errors.New("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupCache != nil {
				tt.setupCache(tt.args, redisMock)
			}

			got, err := tt.args.cache.GetString(ctx, tt.args.key)
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestRedisCache_SetString(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedRedis, redisMock := getRedisCacheMock(ctrl)

	type args struct {
		key   string
		data  string
		cache *redisCache
	}
	tests := []struct {
		name       string
		args       args
		setupCache func(args args, m *mocks.MockIRedisCache)
		wantErr    error
	}{
		{
			name:    "Success",
			args:    args{key: "setString_test_key", data: "test_value", cache: testRedis.(*redisCache)},
			wantErr: nil,
		},
		{
			name: "Error",
			args: args{key: "test_key", data: "test_value", cache: mockedRedis},
			setupCache: func(args args, m *mocks.MockIRedisCache) {
				cmd := redis.NewStatusCmd(ctx, args.key, "OK")
				cmd.SetErr(errors.New("some error"))
				m.EXPECT().Set(ctx, args.key, []byte(args.data), gomock.Any()).
					Return(cmd)
			},
			wantErr: errors.New("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupCache != nil {
				tt.setupCache(tt.args, redisMock)
			}

			err := tt.args.cache.SetString(ctx, tt.args.key, tt.args.data)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestRedisCache_SetStringWithExpiration(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedRedis, redisMock := getRedisCacheMock(ctrl)

	type args struct {
		key        string
		data       string
		expiration time.Duration
		cache      *redisCache
	}
	tests := []struct {
		name       string
		args       args
		setupCache func(args args, m *mocks.MockIRedisCache)
		wantErr    error
	}{
		{
			name:    "Success",
			args:    args{key: "setStringWithExpiration_test_key", data: "test_value", expiration: 10, cache: testRedis.(*redisCache)},
			wantErr: nil,
		},
		{
			name: "Error",
			args: args{key: "test_key", data: "test_value", expiration: 10, cache: mockedRedis},
			setupCache: func(args args, m *mocks.MockIRedisCache) {
				cmd := redis.NewStatusCmd(ctx, args.key, "OK")
				cmd.SetErr(errors.New("some error"))
				m.EXPECT().Set(gomock.Any(), args.key, []byte(args.data), gomock.Any()).
					Return(cmd)
			},
			wantErr: errors.New("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupCache != nil {
				tt.setupCache(tt.args, redisMock)
			}

			err := tt.args.cache.SetStringWithExpiration(ctx, tt.args.key, tt.args.data, tt.args.expiration)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestRedisCache_SetItemWithExpiration(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedRedis, redisMock := getRedisCacheMock(ctrl)

	type args struct {
		key        string
		data       []byte
		expiration time.Duration
		cache      *redisCache
	}
	tests := []struct {
		name       string
		args       args
		setupCache func(args args, m *mocks.MockIRedisCache)
		wantErr    error
	}{
		{
			name: "Success",
			args: args{
				key: "setItemWithExpiration_test_key", data: []byte("test_value"),
				expiration: time.Minute, cache: testRedis.(*redisCache),
			},
			wantErr: nil,
		},
		{
			name: "Error",
			args: args{key: "test_key", data: []byte("test_value"), expiration: time.Minute, cache: mockedRedis},
			setupCache: func(args args, m *mocks.MockIRedisCache) {
				cmd := redis.NewStatusCmd(ctx, args.key, "OK")
				cmd.SetErr(errors.New("some error"))
				m.EXPECT().Set(gomock.Any(), args.key, args.data, gomock.Any()).
					Return(cmd)
			},
			wantErr: errors.New("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupCache != nil {
				tt.setupCache(tt.args, redisMock)
			}

			err := tt.args.cache.SetItemWithExpiration(ctx, tt.args.key, tt.args.data, tt.args.expiration)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestRedisCache_Increase(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedRedis, redisMock := getRedisCacheMock(ctrl)

	type args struct {
		key   string
		cache *redisCache
	}
	tests := []struct {
		name       string
		args       args
		setupCache func(args args, m *mocks.MockIRedisCache)
		wantErr    error
	}{
		{
			name:    "Success",
			args:    args{key: "increase_key", cache: testRedis.(*redisCache)},
			wantErr: nil,
		},
		{
			name: "Error",
			args: args{key: "increase_key", cache: mockedRedis},
			setupCache: func(args args, m *mocks.MockIRedisCache) {
				m.EXPECT().Incr(gomock.Any(), args.key).Return(redis.NewIntResult(0, errors.New("some error")))
			},
			wantErr: errors.New("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupCache != nil {
				tt.setupCache(tt.args, redisMock)
			}

			err := tt.args.cache.Increase(ctx, tt.args.key)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestRedisCache_SetStruct(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedRedis, redisMock := getRedisCacheMock(ctrl)

	type args struct {
		key   string
		data  any
		cache *redisCache
	}
	tests := []struct {
		name       string
		args       args
		setupCache func(args args, m *mocks.MockIRedisCache)
		wantErr    error
	}{
		{
			name:    "Success",
			args:    args{key: "setStruct_test_key", data: "test_value", cache: testRedis.(*redisCache)},
			wantErr: nil,
		},
		{
			name: "Error",
			args: args{key: "test_key", data: "test_value", cache: mockedRedis},
			setupCache: func(args args, m *mocks.MockIRedisCache) {
				dataString, err := json.Marshal(args.data)
				require.NoError(t, err)

				cmd := redis.NewStatusCmd(ctx, args.key, "OK")
				cmd.SetErr(errors.New("some error"))
				m.EXPECT().Set(ctx, args.key, dataString, gomock.Any()).Return(cmd)
			},
			wantErr: errors.New("some error"),
		},
		{
			name:    "Should set with default expiration",
			args:    args{key: "setStruct_default_key", data: "test_value", cache: testRedis.(*redisCache)},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupCache != nil {
				tt.setupCache(tt.args, redisMock)
			}

			err := tt.args.cache.SetStruct(ctx, tt.args.key, tt.args.data)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestRedisCache_SetStructWithExpiration(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedRedis, redisMock := getRedisCacheMock(ctrl)

	type args struct {
		key        string
		data       interface{}
		expiration time.Duration
		cache      *redisCache
	}
	tests := []struct {
		name       string
		args       args
		setupCache func(args args, m *mocks.MockIRedisCache)
		wantErr    error
	}{
		{
			name:    "Success",
			args:    args{key: "setStructWithExpiration_test_key", data: "test_value", expiration: time.Minute, cache: testRedis.(*redisCache)},
			wantErr: nil,
		},
		{
			name: "Error",
			args: args{key: "test_key", data: "test_value", expiration: time.Minute, cache: mockedRedis},
			setupCache: func(args args, m *mocks.MockIRedisCache) {
				dataString, err := json.Marshal(args.data)
				require.NoError(t, err)

				cmd := redis.NewStatusCmd(ctx, args.key, "OK")
				cmd.SetErr(errors.New("some error"))
				m.EXPECT().Set(ctx, args.key, dataString, args.expiration).Return(cmd)
			},
			wantErr: errors.New("some error"),
		},
		{
			name:    "Should set with default expiration",
			args:    args{key: "setStructWithExpiration_default_key", data: "test_value", expiration: 0, cache: testRedis.(*redisCache)},
			wantErr: nil,
		},
		{
			name: "Should return error when fail to marshal data",
			args: args{
				key: "setStructWithExpiration_error_key", data: make(chan int),
				expiration: 0, cache: testRedis.(*redisCache),
			},
			wantErr: errors.New("json: unsupported type: chan int"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupCache != nil {
				tt.setupCache(tt.args, redisMock)
			}

			err := tt.args.cache.SetStructWithExpiration(ctx, tt.args.key, tt.args.data, tt.args.expiration)
			if tt.wantErr != nil {
				require.Equal(t, tt.wantErr.Error(), err.Error())
				return
			}
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestRedisCache_GetStruct(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedRedis, redisMock := getRedisCacheMock(ctrl)

	type someStruct struct {
		Name string
	}

	type args struct {
		key   string
		cache *redisCache
	}
	tests := []struct {
		name       string
		args       args
		setupCache func(args args, m *mocks.MockIRedisCache)
		want       interface{}
		wantErr    error
	}{
		{
			name: "Cache hit",
			args: args{key: "getStruct_test_key", cache: testRedis.(*redisCache)},
			setupCache: func(args args, _ *mocks.MockIRedisCache) {
				data := someStruct{Name: "test_value"}
				err := testRedis.SetStruct(ctx, args.key, data)
				require.NoError(t, err)
			},
			want:    someStruct{Name: "test_value"},
			wantErr: nil,
		},
		{
			name:    "Cache miss",
			args:    args{key: "getStruct_other_key", cache: testRedis.(*redisCache)},
			want:    someStruct{},
			wantErr: ErrCacheMiss,
		},
		{
			name: "Error",
			args: args{key: "test_key", cache: mockedRedis},
			setupCache: func(args args, m *mocks.MockIRedisCache) {
				m.EXPECT().Get(ctx, args.key).Return(redis.NewStringResult("", errors.New("some error")))
			},
			want:    someStruct{},
			wantErr: errors.New("some error"),
		},
		{
			name: "Should return error when fail to unmarshal data",
			args: args{key: "unmarshal_key", cache: mockedRedis},
			setupCache: func(args args, m *mocks.MockIRedisCache) {
				m.EXPECT().Get(ctx, args.key).Return(redis.NewStringResult("invalid_json", nil))
			},
			want:    someStruct{},
			wantErr: errors.New("invalid character 'i' looking for beginning of value"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupCache != nil {
				tt.setupCache(tt.args, redisMock)
			}

			data := someStruct{}
			err := tt.args.cache.GetStruct(ctx, tt.args.key, &data)
			if tt.wantErr != nil {
				require.Equal(t, tt.wantErr.Error(), err.Error())
				return
			}
			require.Equal(t, tt.want, data)
		})
	}
}

func TestRedisCache_GetExpiration(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedRedis, redisMock := getRedisCacheMock(ctrl)

	type args struct {
		key   string
		cache *redisCache
	}
	tests := []struct {
		name       string
		args       args
		setupCache func(args args, m *mocks.MockIRedisCache)
		want       time.Duration
		wantErr    error
	}{
		{
			name: "Success",
			args: args{key: "getExpiration_test_key", cache: testRedis.(*redisCache)},
			setupCache: func(args args, _ *mocks.MockIRedisCache) {
				err := testRedis.SetItem(ctx, args.key, []byte("test_value"))
				require.NoError(t, err)
			},
			want: time.Hour * 24,
		},
		{
			name: "Error",
			args: args{key: "test_key", cache: mockedRedis},
			setupCache: func(args args, m *mocks.MockIRedisCache) {
				cmd := redis.NewDurationCmd(ctx, time.Hour*24)
				cmd.SetErr(errors.New("some error"))
				m.EXPECT().TTL(ctx, args.key).Return(cmd)
			},
			want:    0,
			wantErr: errors.New("some error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupCache != nil {
				tt.setupCache(tt.args, redisMock)
			}

			got, err := tt.args.cache.GetExpiration(ctx, tt.args.key)
			require.Equal(t, tt.wantErr, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestRedisCache_SetExpiration(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedRedis, redisMock := getRedisCacheMock(ctrl)

	type args struct {
		key        string
		expiration time.Duration
		cache      *redisCache
	}
	tests := []struct {
		name       string
		args       args
		setupCache func(args args, m *mocks.MockIRedisCache)
		wantErr    error
	}{
		{
			name:    "Success",
			args:    args{key: "setExpiration_test_key", expiration: time.Hour * 24, cache: testRedis.(*redisCache)},
			wantErr: nil,
		},
		{
			name: "Error",
			args: args{key: "test_key", expiration: time.Hour * 24, cache: mockedRedis},
			setupCache: func(args args, m *mocks.MockIRedisCache) {
				cmd := redis.NewBoolCmd(ctx, time.Hour*24)
				cmd.SetErr(errors.New("some error"))
				m.EXPECT().Expire(ctx, args.key, args.expiration).Return(cmd)
			},
			wantErr: errors.New("some error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupCache != nil {
				tt.setupCache(tt.args, redisMock)
			}

			err := tt.args.cache.SetExpiration(ctx, tt.args.key, tt.args.expiration)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestRedisCache_Delete(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedRedis, redisMock := getRedisCacheMock(ctrl)

	type args struct {
		key   string
		cache *redisCache
	}
	tests := []struct {
		name       string
		args       args
		setupCache func(args args, m *mocks.MockIRedisCache)
		wantErr    error
	}{
		{
			name:    "Success",
			args:    args{key: "delete_test_key", cache: testRedis.(*redisCache)},
			wantErr: nil,
		},
		{
			name: "Error",
			args: args{key: "test_key", cache: mockedRedis},
			setupCache: func(args args, m *mocks.MockIRedisCache) {
				cmd := redis.NewIntCmd(ctx, 1)
				cmd.SetErr(errors.New("some error"))
				m.EXPECT().Del(ctx, args.key).Return(cmd)
			},
			wantErr: errors.New("some error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupCache != nil {
				tt.setupCache(tt.args, redisMock)
			}

			err := tt.args.cache.Delete(ctx, tt.args.key)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestRedisCache_CleanAll(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedRedis, redisMock := getRedisCacheMock(ctrl)

	type args struct {
		key   string
		cache *redisCache
	}

	tests := []struct {
		name       string
		args       args
		setupCache func(args args, m *mocks.MockIRedisCache)
		validate   func(args args)
		wantErr    error
	}{
		{
			name: "Success",
			args: args{key: "cleanAll_test_key", cache: testRedis.(*redisCache)},
			setupCache: func(args args, _ *mocks.MockIRedisCache) {
				err := testRedis.SetItem(ctx, args.key, []byte("test_value"))
				require.NoError(t, err)
			},
			validate: func(args args) {
				data, err := args.cache.GetItem(ctx, args.key)
				require.Equal(t, ErrCacheMiss, err)
				require.Nil(t, data)
			},
			wantErr: nil,
		},
		{
			name:    "Success with no keys",
			args:    args{cache: testRedis.(*redisCache)},
			wantErr: nil,
		},
		{
			name: "Error getting keys",
			args: args{cache: mockedRedis},
			setupCache: func(args args, m *mocks.MockIRedisCache) {
				m.EXPECT().Keys(ctx, "*").Return(redis.NewStringSliceResult([]string{}, errors.New("some error")))
			},
			wantErr: errors.New("some error"),
		},
		{
			name: "Error deleting keys",
			args: args{cache: mockedRedis},
			setupCache: func(args args, m *mocks.MockIRedisCache) {
				m.EXPECT().Keys(ctx, "*").Return(redis.NewStringSliceResult([]string{"test_key"}, nil))
				m.EXPECT().Del(ctx, "test_key").Return(redis.NewIntResult(0, errors.New("some error")))
			},
			wantErr: errors.New("some error"),
		},
		{
			name: "Cache miss deleting keys",
			args: args{cache: mockedRedis},
			setupCache: func(args args, m *mocks.MockIRedisCache) {
				m.EXPECT().Keys(ctx, "*").Return(redis.NewStringSliceResult([]string{"test_key"}, nil))
				m.EXPECT().Del(ctx, "test_key").Return(redis.NewIntResult(0, redis.Nil))
			},
			wantErr: ErrCacheMiss,
		},
		{
			name: "Cache miss getting keys",
			args: args{cache: mockedRedis},
			setupCache: func(args args, m *mocks.MockIRedisCache) {
				m.EXPECT().Keys(ctx, "*").Return(redis.NewStringSliceResult([]string{"test_key"}, redis.Nil))
			},
			wantErr: ErrCacheMiss,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupCache != nil {
				tt.setupCache(tt.args, redisMock)
			}

			err := tt.args.cache.CleanAll(ctx)
			require.Equal(t, tt.wantErr, err)

			if tt.validate != nil {
				tt.validate(tt.args)
			}
		})
	}
}
