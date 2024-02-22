package cache

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/diegoclair/go_boilerplate/mocks/infra"
	"github.com/golang/mock/gomock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func getRedisCacheMock(ctrl *gomock.Controller) (*redisCache, *infra.MockIRedisCache) {
	m := infra.NewMockIRedisCache(ctrl)
	return &redisCache{
		redis: m,
	}, m
}

func TestRedisCache_GetItem(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	c, cm := getRedisCacheMock(ctrl)

	type args struct {
		key string
	}
	tests := []struct {
		name       string
		args       args
		setupCache func(args args, m *infra.MockIRedisCache)
		want       []byte
		wantErr    error
	}{
		{
			name: "Cache hit",
			args: args{key: "test_key"},
			setupCache: func(args args, m *infra.MockIRedisCache) {
				m.EXPECT().Get(gomock.Any(), args.key).Return(redis.NewStringResult("test_value", nil))
			},
			want:    []byte("test_value"),
			wantErr: nil,
		},
		{
			name: "Cache miss",
			args: args{key: "test_key"},
			setupCache: func(args args, m *infra.MockIRedisCache) {
				m.EXPECT().Get(gomock.Any(), args.key).Return(redis.NewStringResult("", redis.Nil))
			},
			want:    nil,
			wantErr: ErrCacheMiss,
		},
		{
			name: "Error",
			args: args{key: "test_key"},
			setupCache: func(args args, m *infra.MockIRedisCache) {
				m.EXPECT().Get(gomock.Any(), args.key).Return(redis.NewStringResult("", errors.New("some error")))
			},
			want:    nil,
			wantErr: errors.New("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupCache != nil {
				tt.setupCache(tt.args, cm)
			}

			got, err := c.GetItem(ctx, tt.args.key)
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestRedisCache_SetItem(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	c, cm := getRedisCacheMock(ctrl)

	type args struct {
		key  string
		data []byte
	}
	tests := []struct {
		name       string
		args       args
		setupCache func(args args, m *infra.MockIRedisCache)
		wantErr    error
	}{
		{
			name: "Success",
			args: args{key: "test_key", data: []byte("test_value")},
			setupCache: func(args args, m *infra.MockIRedisCache) {
				m.EXPECT().Set(gomock.Any(), args.key, args.data, gomock.Any()).Return(redis.NewStatusCmd(ctx, args.key, "OK"))
			},
			wantErr: nil,
		},
		{
			name: "Error",
			args: args{key: "test_key", data: []byte("test_value")},
			setupCache: func(args args, m *infra.MockIRedisCache) {
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
				tt.setupCache(tt.args, cm)
			}

			err := c.SetItem(ctx, tt.args.key, tt.args.data)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestRedisCache_GetInt(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	c, cm := getRedisCacheMock(ctrl)

	type args struct {
		key string
	}
	tests := []struct {
		name       string
		args       args
		setupCache func(args args, m *infra.MockIRedisCache)
		want       int64
		wantErr    error
	}{
		{
			name: "Cache hit",
			args: args{key: "test_key"},
			setupCache: func(args args, m *infra.MockIRedisCache) {
				m.EXPECT().Get(gomock.Any(), args.key).Return(redis.NewStringResult("123", nil))
			},
			want:    123,
			wantErr: nil,
		},
		{
			name: "Cache miss",
			args: args{key: "test_key"},
			setupCache: func(args args, m *infra.MockIRedisCache) {
				m.EXPECT().Get(gomock.Any(), args.key).Return(redis.NewStringResult("", redis.Nil))
			},
			want:    0,
			wantErr: ErrCacheMiss,
		},
		{
			name: "Error",
			args: args{key: "test_key"},
			setupCache: func(args args, m *infra.MockIRedisCache) {
				m.EXPECT().Get(gomock.Any(), args.key).Return(redis.NewStringResult("", errors.New("some error")))
			},
			want:    0,
			wantErr: errors.New("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupCache != nil {
				tt.setupCache(tt.args, cm)
			}

			got, err := c.GetInt(ctx, tt.args.key)
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestRedisCache_GetString(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	c, cm := getRedisCacheMock(ctrl)

	type args struct {
		key string
	}
	tests := []struct {
		name       string
		args       args
		setupCache func(args args, m *infra.MockIRedisCache)
		want       string
		wantErr    error
	}{
		{
			name: "Cache hit",
			args: args{key: "test_key"},
			setupCache: func(args args, m *infra.MockIRedisCache) {
				m.EXPECT().Get(gomock.Any(), args.key).Return(redis.NewStringResult("test_value", nil))
			},
			want:    "test_value",
			wantErr: nil,
		},
		{
			name: "Cache miss",
			args: args{key: "test_key"},
			setupCache: func(args args, m *infra.MockIRedisCache) {
				m.EXPECT().Get(gomock.Any(), args.key).Return(redis.NewStringResult("", redis.Nil))
			},
			want:    "",
			wantErr: ErrCacheMiss,
		},
		{
			name: "Error",
			args: args{key: "test_key"},
			setupCache: func(args args, m *infra.MockIRedisCache) {
				m.EXPECT().Get(gomock.Any(), args.key).Return(redis.NewStringResult("", errors.New("some error")))
			},
			want:    "",
			wantErr: errors.New("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupCache != nil {
				tt.setupCache(tt.args, cm)
			}

			got, err := c.GetString(ctx, tt.args.key)
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestRedisCache_SetString(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	c, cm := getRedisCacheMock(ctrl)

	type args struct {
		key  string
		data string
	}
	tests := []struct {
		name       string
		args       args
		setupCache func(args args, m *infra.MockIRedisCache)
		wantErr    error
	}{
		{
			name: "Success",
			args: args{key: "test_key", data: "test_value"},
			setupCache: func(args args, m *infra.MockIRedisCache) {
				m.EXPECT().Set(ctx, args.key, []byte(args.data), gomock.Any()).Return(redis.NewStatusCmd(ctx, args.key, "OK"))
			},
			wantErr: nil,
		},
		{
			name: "Error",
			args: args{key: "test_key", data: "test_value"},
			setupCache: func(args args, m *infra.MockIRedisCache) {
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
				tt.setupCache(tt.args, cm)
			}

			err := c.SetString(ctx, tt.args.key, tt.args.data)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestRedisCache_SetStringWithExpiration(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	c, cm := getRedisCacheMock(ctrl)

	type args struct {
		key        string
		data       string
		expiration time.Duration
	}
	tests := []struct {
		name       string
		args       args
		setupCache func(args args, m *infra.MockIRedisCache)
		wantErr    error
	}{
		{
			name: "Success",
			args: args{key: "test_key", data: "test_value", expiration: 10},
			setupCache: func(args args, m *infra.MockIRedisCache) {
				m.EXPECT().Set(gomock.Any(), args.key, []byte(args.data), gomock.Any()).Return(redis.NewStatusCmd(ctx, args.key, "OK"))
			},
			wantErr: nil,
		},
		{
			name: "Error",
			args: args{key: "test_key", data: "test_value", expiration: 10},
			setupCache: func(args args, m *infra.MockIRedisCache) {
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
				tt.setupCache(tt.args, cm)
			}

			err := c.SetStringWithExpiration(ctx, tt.args.key, tt.args.data, tt.args.expiration)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestRedisCache_SetItemWithExpiration(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	c, cm := getRedisCacheMock(ctrl)

	type args struct {
		key        string
		data       []byte
		expiration time.Duration
	}
	tests := []struct {
		name       string
		args       args
		setupCache func(args args, m *infra.MockIRedisCache)
		wantErr    error
	}{
		{
			name: "Success",
			args: args{key: "test_key", data: []byte("test_value"), expiration: 10},
			setupCache: func(args args, m *infra.MockIRedisCache) {
				m.EXPECT().Set(gomock.Any(), args.key, args.data, gomock.Any()).Return(redis.NewStatusCmd(ctx, args.key, "OK"))
			},
			wantErr: nil,
		},
		{
			name: "Error",
			args: args{key: "test_key", data: []byte("test_value"), expiration: 10},
			setupCache: func(args args, m *infra.MockIRedisCache) {
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
				tt.setupCache(tt.args, cm)
			}

			err := c.SetItemWithExpiration(ctx, tt.args.key, tt.args.data, tt.args.expiration)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestRedisCache_Increase(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	c, cm := getRedisCacheMock(ctrl)

	type args struct {
		key string
	}
	tests := []struct {
		name       string
		args       args
		setupCache func(args args, m *infra.MockIRedisCache)
		wantErr    error
	}{
		{
			name: "Success",
			args: args{key: "test_key"},
			setupCache: func(args args, m *infra.MockIRedisCache) {
				m.EXPECT().Incr(gomock.Any(), args.key).Return(redis.NewIntResult(1, nil))
			},
			wantErr: nil,
		},
		{
			name: "Error",
			args: args{key: "test_key"},
			setupCache: func(args args, m *infra.MockIRedisCache) {
				m.EXPECT().Incr(gomock.Any(), args.key).Return(redis.NewIntResult(0, errors.New("some error")))
			},
			wantErr: errors.New("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupCache != nil {
				tt.setupCache(tt.args, cm)
			}

			err := c.Increase(ctx, tt.args.key)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestRedisCache_SetStruct(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	c, cm := getRedisCacheMock(ctrl)

	type args struct {
		key        string
		data       interface{}
		expiration time.Duration
	}
	tests := []struct {
		name       string
		args       args
		setupCache func(args args, m *infra.MockIRedisCache)
		wantErr    error
	}{
		{
			name: "Success",
			args: args{key: "test_key", data: "test_value", expiration: 10},
			setupCache: func(args args, m *infra.MockIRedisCache) {
				dataString, err := json.Marshal(args.data)
				require.NoError(t, err)
				m.EXPECT().Set(ctx, args.key, dataString, args.expiration).Return(redis.NewStatusCmd(ctx, args.key, "OK"))
			},
			wantErr: nil,
		},
		{
			name: "Error",
			args: args{key: "test_key", data: "test_value", expiration: 10},
			setupCache: func(args args, m *infra.MockIRedisCache) {
				dataString, err := json.Marshal(args.data)
				require.NoError(t, err)

				cmd := redis.NewStatusCmd(ctx, args.key, "OK")
				cmd.SetErr(errors.New("some error"))
				m.EXPECT().Set(ctx, args.key, dataString, args.expiration).Return(cmd)
			},
			wantErr: errors.New("some error"),
		},
		{
			name: "Should set with expiration",
			args: args{key: "test_key", data: "test_value", expiration: 0},
			setupCache: func(args args, m *infra.MockIRedisCache) {
				dataString, err := json.Marshal(args.data)
				require.NoError(t, err)
				m.EXPECT().Set(ctx, args.key, dataString, gomock.Any()).Return(redis.NewStatusCmd(ctx, args.key, "OK"))
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupCache != nil {
				tt.setupCache(tt.args, cm)
			}

			err := c.SetStruct(ctx, tt.args.key, tt.args.data, tt.args.expiration)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestRedisCache_SetStructWithExpiration(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	c, cm := getRedisCacheMock(ctrl)

	type args struct {
		key        string
		data       interface{}
		expiration time.Duration
	}
	tests := []struct {
		name       string
		args       args
		setupCache func(args args, m *infra.MockIRedisCache)
		wantErr    error
	}{
		{
			name: "Success",
			args: args{key: "test_key", data: "test_value", expiration: 10},
			setupCache: func(args args, m *infra.MockIRedisCache) {
				dataString, err := json.Marshal(args.data)
				require.NoError(t, err)
				m.EXPECT().Set(ctx, args.key, dataString, args.expiration).Return(redis.NewStatusCmd(ctx, args.key, "OK"))
			},
			wantErr: nil,
		},
		{
			name: "Error",
			args: args{key: "test_key", data: "test_value", expiration: 10},
			setupCache: func(args args, m *infra.MockIRedisCache) {
				dataString, err := json.Marshal(args.data)
				require.NoError(t, err)

				cmd := redis.NewStatusCmd(ctx, args.key, "OK")
				cmd.SetErr(errors.New("some error"))
				m.EXPECT().Set(ctx, args.key, dataString, args.expiration).Return(cmd)
			},
			wantErr: errors.New("some error"),
		},
		{
			name: "Error marshalling",
			args: args{key: "test_key", data: make(chan int), expiration: 10},
			setupCache: func(args args, m *infra.MockIRedisCache) {
				m.EXPECT().Set(gomock.Any(), args.key, gomock.Any(), gomock.Any()).Times(0)
			},
			wantErr: errors.New("json: unsupported type: chan int"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupCache != nil {
				tt.setupCache(tt.args, cm)
			}

			err := c.SetStructWithExpiration(ctx, tt.args.key, tt.args.data, tt.args.expiration)
			if tt.wantErr != nil {
				require.Equal(t, tt.wantErr.Error(), err.Error())
			}
		})
	}
}

func TestRedisCache_GetStruct(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	c, cm := getRedisCacheMock(ctrl)

	type testStruct struct {
		Test string
	}

	type args struct {
		key  string
		data testStruct
	}
	tests := []struct {
		name       string
		args       args
		setupCache func(args args, m *infra.MockIRedisCache)
		wantErr    error
	}{
		{
			name: "Success",
			args: args{key: "test_key", data: testStruct{Test: "test_value"}},
			setupCache: func(args args, m *infra.MockIRedisCache) {
				dataString, err := json.Marshal(args.data)
				require.NoError(t, err)
				m.EXPECT().Get(gomock.Any(), args.key).Return(redis.NewStringResult(string(dataString), nil))
			},
			wantErr: nil,
		},
		{
			name: "Cache miss",
			args: args{key: "test_key", data: testStruct{Test: "test_value"}},
			setupCache: func(args args, m *infra.MockIRedisCache) {
				m.EXPECT().Get(gomock.Any(), args.key).Return(redis.NewStringResult("", redis.Nil))
			},
			wantErr: ErrCacheMiss,
		},
		{
			name: "Error",
			args: args{key: "test_key", data: testStruct{Test: "test_value"}},
			setupCache: func(args args, m *infra.MockIRedisCache) {
				m.EXPECT().Get(gomock.Any(), args.key).Return(redis.NewStringResult("", errors.New("some error")))
			},
			wantErr: errors.New("some error"),
		},
		{
			name: "Invalid data",
			args: args{key: "test_key", data: testStruct{Test: "test_value"}},
			setupCache: func(args args, m *infra.MockIRedisCache) {
				m.EXPECT().Get(gomock.Any(), args.key).Return(redis.NewStringResult("invalid_data", nil))
			},
			wantErr: errors.New("invalid character 'i' looking for beginning of value"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupCache != nil {
				tt.setupCache(tt.args, cm)
			}

			err := c.GetStruct(ctx, tt.args.key, &tt.args.data)
			if tt.wantErr != nil {
				require.Equal(t, tt.wantErr.Error(), err.Error())
			}
		})
	}
}

func TestRedisCache_GetExpiration(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	c, cm := getRedisCacheMock(ctrl)

	type args struct {
		key string
	}
	tests := []struct {
		name       string
		args       args
		setupCache func(args args, m *infra.MockIRedisCache)
		want       time.Duration
		wantErr    error
	}{
		{
			name: "Success",
			args: args{key: "test_key"},
			setupCache: func(args args, m *infra.MockIRedisCache) {
				m.EXPECT().TTL(gomock.Any(), args.key).Return(redis.NewDurationResult(10*time.Second, nil))
			},
			want:    10 * time.Second,
			wantErr: nil,
		},
		{
			name: "Error",
			args: args{key: "test_key"},
			setupCache: func(args args, m *infra.MockIRedisCache) {
				m.EXPECT().TTL(gomock.Any(), args.key).Return(redis.NewDurationResult(0, errors.New("some error")))
			},
			want:    0,
			wantErr: errors.New("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupCache != nil {
				tt.setupCache(tt.args, cm)
			}

			got, err := c.GetExpiration(ctx, tt.args.key)
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestRedisCache_SetExpiration(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	c, cm := getRedisCacheMock(ctrl)

	type args struct {
		key        string
		expiration time.Duration
	}
	tests := []struct {
		name       string
		args       args
		setupCache func(args args, m *infra.MockIRedisCache)
		wantErr    error
	}{
		{
			name: "Success",
			args: args{key: "test_key", expiration: 10},
			setupCache: func(args args, m *infra.MockIRedisCache) {
				m.EXPECT().Expire(gomock.Any(), args.key, args.expiration).Return(redis.NewBoolResult(true, nil))
			},
			wantErr: nil,
		},
		{
			name: "Error",
			args: args{key: "test_key", expiration: 10},
			setupCache: func(args args, m *infra.MockIRedisCache) {
				m.EXPECT().Expire(gomock.Any(), args.key, args.expiration).Return(redis.NewBoolResult(false, errors.New("some error")))
			},
			wantErr: errors.New("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupCache != nil {
				tt.setupCache(tt.args, cm)
			}

			err := c.SetExpiration(ctx, tt.args.key, tt.args.expiration)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestRedisCache_Delete(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	c, cm := getRedisCacheMock(ctrl)

	type args struct {
		keys string
	}
	tests := []struct {
		name       string
		args       args
		setupCache func(args args, m *infra.MockIRedisCache)
		wantErr    error
	}{
		{
			name: "Success",
			args: args{keys: "test_key"},
			setupCache: func(args args, m *infra.MockIRedisCache) {
				m.EXPECT().Del(ctx, args.keys).Return(redis.NewIntResult(1, nil))
			},
			wantErr: nil,
		},
		{
			name: "Error",
			args: args{keys: "test_key"},
			setupCache: func(args args, m *infra.MockIRedisCache) {
				m.EXPECT().Del(ctx, args.keys).Return(redis.NewIntResult(0, errors.New("some error")))
			},
			wantErr: errors.New("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupCache != nil {
				tt.setupCache(tt.args, cm)
			}

			err := c.Delete(ctx, tt.args.keys)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestRedisCache_CleanAll(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	c, cm := getRedisCacheMock(ctrl)

	tests := []struct {
		name       string
		setupCache func(m *infra.MockIRedisCache)
		wantErr    error
	}{
		{
			name: "Success",
			setupCache: func(m *infra.MockIRedisCache) {
				m.EXPECT().Keys(ctx, "*").Return(redis.NewStringSliceResult([]string{"test_key"}, nil))
				m.EXPECT().Del(ctx, "test_key").Return(redis.NewIntResult(1, nil))
			},
			wantErr: nil,
		},
		{
			name: "Success with no keys",
			setupCache: func(m *infra.MockIRedisCache) {
				m.EXPECT().Keys(ctx, "*").Return(redis.NewStringSliceResult([]string{}, nil))
			},
			wantErr: nil,
		},
		{
			name: "Error getting keys",
			setupCache: func(m *infra.MockIRedisCache) {
				m.EXPECT().Keys(ctx, "*").Return(redis.NewStringSliceResult([]string{}, errors.New("some error")))
			},
			wantErr: errors.New("some error"),
		},
		{
			name: "Error deleting keys",
			setupCache: func(m *infra.MockIRedisCache) {
				m.EXPECT().Keys(ctx, "*").Return(redis.NewStringSliceResult([]string{"test_key"}, nil))
				m.EXPECT().Del(ctx, "test_key").Return(redis.NewIntResult(0, errors.New("some error")))
			},
			wantErr: errors.New("some error"),
		},
		{
			name: "Cache miss deleting keys",
			setupCache: func(m *infra.MockIRedisCache) {
				m.EXPECT().Keys(ctx, "*").Return(redis.NewStringSliceResult([]string{"test_key"}, nil))
				m.EXPECT().Del(ctx, "test_key").Return(redis.NewIntResult(0, redis.Nil))
			},
			wantErr: ErrCacheMiss,
		},
		{
			name: "Cache miss getting keys",
			setupCache: func(m *infra.MockIRedisCache) {
				m.EXPECT().Keys(ctx, "*").Return(redis.NewStringSliceResult([]string{"test_key"}, redis.Nil))
			},
			wantErr: ErrCacheMiss,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupCache != nil {
				tt.setupCache(cm)
			}

			err := c.CleanAll(ctx)
			require.Equal(t, tt.wantErr, err)
		})
	}
}
