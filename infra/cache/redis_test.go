package cache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/diegoclair/go_boilerplate/infra/mocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func getRedisCacheMock(ctrl *gomock.Controller) (*CacheManager, *mocks.MockIRedisCache) {
	m := mocks.NewMockIRedisCache(ctrl)
	return &CacheManager{
		redis: m,
	}, m
}

func TestNewRedisCache(t *testing.T) {
	ctx := context.Background()

	cacheManager, client, err := NewRedisCache(ctx,
		cfg.Redis.Host, cfg.Redis.Password, cfg.Redis.DB, cfg.Redis.DefaultExpiration,
		cfg.GetLogger(),
	)
	require.NoError(t, err)
	require.NotNil(t, cacheManager)
	require.NotNil(t, client)
}

func TestRedisCache_Get(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedRedis, redisMock := getRedisCacheMock(ctrl)

	type args struct {
		key   string
		cache *CacheManager
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
			args: args{key: "get_key", cache: testRedis},
			setupTest: func(args args, _ *mocks.MockIRedisCache) {
				err := testRedis.Set(ctx, args.key, []byte("test_value"))
				require.NoError(t, err)
			},
			want:    []byte("test_value"),
			wantErr: nil,
		},
		{
			name:    "Cache miss",
			args:    args{key: "get_other_key", cache: testRedis},
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

			got, err := tt.args.cache.Get(ctx, tt.args.key)
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestRedisCache_Set(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedRedis, redisMock := getRedisCacheMock(ctrl)

	type args struct {
		key        string
		data       any
		expiration []time.Duration
		cache      *CacheManager
	}
	tests := []struct {
		name       string
		args       args
		setupCache func(args args, m *mocks.MockIRedisCache)
		wantErr    error
	}{
		{
			name:    "Set string with default expiration",
			args:    args{key: "set_string_key", data: "test_value", cache: testRedis},
			wantErr: nil,
		},
		{
			name:    "Set string with custom expiration",
			args:    args{key: "set_string_exp_key", data: "test_value", expiration: []time.Duration{10 * time.Second}, cache: testRedis},
			wantErr: nil,
		},
		{
			name:    "Set bytes with default expiration",
			args:    args{key: "set_bytes_key", data: []byte("test_value"), cache: testRedis},
			wantErr: nil,
		},
		{
			name:    "Set struct via JSON",
			args:    args{key: "set_struct_key", data: struct{ Name string }{Name: "test"}, cache: testRedis},
			wantErr: nil,
		},
		{
			name:    "Set struct with custom expiration",
			args:    args{key: "set_struct_exp_key", data: struct{ Name string }{Name: "test"}, expiration: []time.Duration{time.Minute}, cache: testRedis},
			wantErr: nil,
		},
		{
			name:    "Should return error when fail to marshal data",
			args:    args{key: "set_marshal_error_key", data: make(chan int), cache: testRedis},
			wantErr: errors.New("json: unsupported type: chan int"),
		},
		{
			name: "Should return error when redis fails",
			args: args{key: "test_key", data: "test_value", cache: mockedRedis},
			setupCache: func(args args, m *mocks.MockIRedisCache) {
				cmd := redis.NewStatusCmd(ctx, args.key, "OK")
				cmd.SetErr(errors.New("some error"))
				m.EXPECT().Set(gomock.Any(), args.key, []byte(args.data.(string)), gomock.Any()).
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

			err := tt.args.cache.Set(ctx, tt.args.key, tt.args.data, tt.args.expiration...)
			if tt.wantErr != nil {
				require.EqualError(t, err, tt.wantErr.Error())
				return
			}
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
		cache *CacheManager
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
			args: args{key: "getString_test_key", cache: testRedis},
			setupCache: func(args args, _ *mocks.MockIRedisCache) {
				err := testRedis.Set(ctx, args.key, "test_value")
				require.NoError(t, err)
			},
			want:    "test_value",
			wantErr: nil,
		},
		{
			name:    "Cache miss",
			args:    args{key: "getString_other_key", cache: testRedis},
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

func TestRedisCache_Increase(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedRedis, redisMock := getRedisCacheMock(ctrl)

	type args struct {
		key   string
		cache *CacheManager
	}
	tests := []struct {
		name       string
		args       args
		setupCache func(args args, m *mocks.MockIRedisCache)
		wantErr    error
	}{
		{
			name:    "Success",
			args:    args{key: "increase_key", cache: testRedis},
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

func Test_redisCache_Increase(t *testing.T) {
	ctx := context.Background()
	cacheRegis := testRedis

	t.Run("Success", func(t *testing.T) {
		key := "increase_key_1"
		err := cacheRegis.Increase(ctx, key)
		require.NoError(t, err)
	})

	t.Run("success_2", func(t *testing.T) {
		key := "increase_key_2"
		err := cacheRegis.Increase(ctx, key)
		require.NoError(t, err)

		err = cacheRegis.Increase(ctx, key)
		require.NoError(t, err)

		value, err := cacheRegis.GetInt(ctx, key)
		require.NoError(t, err)

		require.Equal(t, int64(2), value)
	})
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
		cache *CacheManager
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
			args: args{key: "getStruct_test_key", cache: testRedis},
			setupCache: func(args args, _ *mocks.MockIRedisCache) {
				data := someStruct{Name: "test_value"}
				err := testRedis.Set(ctx, args.key, data)
				require.NoError(t, err)
			},
			want:    someStruct{Name: "test_value"},
			wantErr: nil,
		},
		{
			name:    "Cache miss",
			args:    args{key: "getStruct_other_key", cache: testRedis},
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
		cache *CacheManager
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
			args: args{key: "getExpiration_test_key", cache: testRedis},
			setupCache: func(args args, _ *mocks.MockIRedisCache) {
				err := testRedis.Set(ctx, args.key, "test_value")
				require.NoError(t, err)
			},
			want: time.Minute,
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
		cache      *CacheManager
	}
	tests := []struct {
		name       string
		args       args
		setupCache func(args args, m *mocks.MockIRedisCache)
		wantErr    error
	}{
		{
			name:    "Success",
			args:    args{key: "setExpiration_test_key", expiration: time.Hour * 24, cache: testRedis},
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
		cache *CacheManager
	}
	tests := []struct {
		name       string
		args       args
		setupCache func(args args, m *mocks.MockIRedisCache)
		wantErr    error
	}{
		{
			name:    "Success",
			args:    args{key: "delete_test_key", cache: testRedis},
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
		cache *CacheManager
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
			args: args{key: "cleanAll_test_key", cache: testRedis},
			setupCache: func(args args, _ *mocks.MockIRedisCache) {
				err := testRedis.Set(ctx, args.key, "test_value")
				require.NoError(t, err)
			},
			validate: func(args args) {
				data, err := args.cache.Get(ctx, args.key)
				require.Equal(t, ErrCacheMiss, err)
				require.Nil(t, data)
			},
			wantErr: nil,
		},
		{
			name:    "Success with no keys",
			args:    args{cache: testRedis},
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

func TestRedisCache_GetAllKeys(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedRedis, redisMock := getRedisCacheMock(ctrl)

	type args struct {
		pattern string
		cache   *CacheManager
	}
	tests := []struct {
		name      string
		args      args
		setupTest func(args args, m *mocks.MockIRedisCache)
		want      []string
		wantErr   error
	}{
		{
			name: "Success",
			args: args{
				pattern: "test*",
				cache:   mockedRedis,
			},
			setupTest: func(args args, m *mocks.MockIRedisCache) {
				m.EXPECT().Keys(ctx, "test*").Return(redis.NewStringSliceResult([]string{"test1", "test2"}, nil))
			},
			want: []string{"test1", "test2"},
		},
		{
			name: "Error getting keys",
			args: args{
				pattern: "test*",
				cache:   mockedRedis,
			},
			setupTest: func(args args, m *mocks.MockIRedisCache) {
				m.EXPECT().Keys(ctx, "test*").Return(redis.NewStringSliceResult(nil, errors.New("some error")))
			},
			wantErr: errors.New("failed to get keys: some error"),
		},
		{
			name: "No keys found",
			args: args{
				pattern: "nonexistent*",
				cache:   mockedRedis,
			},
			setupTest: func(args args, m *mocks.MockIRedisCache) {
				m.EXPECT().Keys(ctx, "nonexistent*").Return(redis.NewStringSliceResult([]string{}, nil))
			},
			want: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupTest != nil {
				tt.setupTest(tt.args, redisMock)
			}

			got, err := tt.args.cache.GetAllKeys(ctx, tt.args.pattern)
			if tt.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tt.wantErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}
