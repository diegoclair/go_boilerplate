package routeutils_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/diegoclair/go_boilerplate/infra"
	"github.com/diegoclair/go_boilerplate/internal/transport/rest/routeutils"
	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func setupEchoContext(queryParams map[string]string) echo.Context {
	e := echo.New()

	// Build URL with proper encoding
	u := &url.URL{Path: "/"}
	if len(queryParams) > 0 {
		values := url.Values{}
		for key, value := range queryParams {
			values.Set(key, value)
		}
		u.RawQuery = values.Encode()
	}

	req := httptest.NewRequest(http.MethodGet, u.String(), nil)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec)
}

func TestGetRequiredParam(t *testing.T) {
	tests := []struct {
		name         string
		rawValue     string
		converter    routeutils.ArrayConverter[string]
		errorMessage string
		wantValue    string
		wantError    bool
	}{
		{
			name:         "Valid string value",
			rawValue:     "hello",
			converter:    routeutils.StringConverter,
			errorMessage: "Invalid parameter",
			wantValue:    "hello",
			wantError:    false,
		},
		{
			name:         "Empty value should return error",
			rawValue:     "",
			converter:    routeutils.StringConverter,
			errorMessage: "Invalid parameter",
			wantValue:    "",
			wantError:    true,
		},
		{
			name:         "Whitespace only value should return error",
			rawValue:     "   ",
			converter:    routeutils.StringConverter,
			errorMessage: "Invalid parameter",
			wantValue:    "",
			wantError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := routeutils.GetRequiredParam(tt.rawValue, tt.converter, tt.errorMessage)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantValue, got)
			}
		})
	}
}

func TestGetRequiredParamWithInt64(t *testing.T) {
	tests := []struct {
		name      string
		rawValue  string
		wantValue int64
		wantError bool
	}{
		{
			name:      "Valid positive int64",
			rawValue:  "123",
			wantValue: 123,
			wantError: false,
		},
		{
			name:      "Valid negative int64",
			rawValue:  "-456",
			wantValue: -456,
			wantError: false,
		},
		{
			name:      "Zero value should return error (zero value check)",
			rawValue:  "0",
			wantValue: 0,
			wantError: true,
		},
		{
			name:      "Invalid int64",
			rawValue:  "invalid",
			wantValue: 0,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := routeutils.GetRequiredParam(tt.rawValue, routeutils.Int64Converter, "Invalid parameter")

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantValue, got)
			}
		})
	}
}

func TestGetContext(t *testing.T) {
	tests := []struct {
		name         string
		userUUID     string
		sessionValue string
	}{
		{
			name:         "Context with user UUID and session",
			userUUID:     "user-uuid-123",
			sessionValue: "session-value-456",
		},
		{
			name:         "Context with empty values",
			userUUID:     "",
			sessionValue: "",
		},
		{
			name:         "Context with only user UUID",
			userUUID:     "user-uuid-789",
			sessionValue: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Set values in echo context
			if tt.userUUID != "" {
				c.Set(infra.AccountUUIDKey.String(), tt.userUUID)
			}
			if tt.sessionValue != "" {
				c.Set(infra.SessionKey.String(), tt.sessionValue)
			}

			ctx := routeutils.GetContext(c)

			// Verify the context contains the expected values
			assert.NotNil(t, ctx)

			userUUIDFromCtx := ctx.Value(infra.AccountUUIDKey)
			sessionFromCtx := ctx.Value(infra.SessionKey)

			if tt.userUUID != "" {
				assert.Equal(t, tt.userUUID, userUUIDFromCtx)
			} else {
				assert.Nil(t, userUUIDFromCtx)
			}

			if tt.sessionValue != "" {
				assert.Equal(t, tt.sessionValue, sessionFromCtx)
			} else {
				assert.Nil(t, sessionFromCtx)
			}
		})
	}
}

func TestGetRequiredStringQueryParam(t *testing.T) {
	tests := []struct {
		name       string
		paramName  string
		paramValue string
		wantValue  string
		wantError  bool
	}{
		{
			name:       "Valid string query parameter",
			paramName:  "search",
			paramValue: "test-value",
			wantValue:  "test-value",
			wantError:  false,
		},
		{
			name:       "Empty parameter should return error",
			paramName:  "search",
			paramValue: "",
			wantValue:  "",
			wantError:  true,
		},
		{
			name:       "Whitespace only parameter should return error",
			paramName:  "search",
			paramValue: "   ",
			wantValue:  "   ",
			wantError:  true,
		},
		{
			name:       "Parameter with valid content surrounded by whitespace",
			paramName:  "search",
			paramValue: "  valid-content  ",
			wantValue:  "  valid-content  ",
			wantError:  false, // GetRequiredStringQueryParam only trims for empty check, but returns original value
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queryParams := map[string]string{}
			if tt.paramValue != "" {
				queryParams[tt.paramName] = tt.paramValue
			}

			c := setupEchoContext(queryParams)

			got, err := routeutils.GetRequiredStringQueryParam(c, tt.paramName, "Invalid "+tt.paramName)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantValue, got)
			}
		})
	}
}

func TestGetRequiredInt64PathParam(t *testing.T) {
	tests := []struct {
		name       string
		paramName  string
		paramValue string
		wantValue  int64
		wantError  bool
	}{
		{
			name:       "Valid int64 parameter",
			paramName:  "id",
			paramValue: "123",
			wantValue:  123,
			wantError:  false,
		},
		{
			name:       "Zero value should return error",
			paramName:  "id",
			paramValue: "0",
			wantValue:  0,
			wantError:  true,
		},
		{
			name:       "Invalid int64 parameter",
			paramName:  "id",
			paramValue: "invalid",
			wantValue:  0,
			wantError:  true,
		},
		{
			name:       "Empty parameter",
			paramName:  "id",
			paramValue: "",
			wantValue:  0,
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames(tt.paramName)
			c.SetParamValues(tt.paramValue)

			got, err := routeutils.GetRequiredInt64PathParam(c, tt.paramName, "Invalid "+tt.paramName)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantValue, got)
			}
		})
	}
}

func TestGetRequiredStringPathParam(t *testing.T) {
	tests := []struct {
		name       string
		paramName  string
		paramValue string
		wantValue  string
		wantError  bool
	}{
		{
			name:       "Valid string parameter",
			paramName:  "uuid",
			paramValue: "abc-123",
			wantValue:  "abc-123",
			wantError:  false,
		},
		{
			name:       "Empty parameter should return error",
			paramName:  "uuid",
			paramValue: "",
			wantValue:  "",
			wantError:  true,
		},
		{
			name:       "Whitespace only parameter should return error",
			paramName:  "uuid",
			paramValue: "   ",
			wantValue:  "   ",
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames(tt.paramName)
			c.SetParamValues(tt.paramValue)

			got, err := routeutils.GetRequiredStringPathParam(c, tt.paramName, "Invalid "+tt.paramName)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantValue, got)
			}
		})
	}
}

func TestGetPagingParams(t *testing.T) {
	tests := []struct {
		name     string
		page     string
		quantity string
		wantTake int64
		wantSkip int64
	}{
		{
			name:     "Valid pagination parameters",
			page:     "2",
			quantity: "20",
			wantTake: 20,
			wantSkip: 20,
		},
		{
			name:     "Default values for empty parameters",
			page:     "",
			quantity: "",
			wantTake: 10,
			wantSkip: 0,
		},
		{
			name:     "Page below 1 should default to 1",
			page:     "0",
			quantity: "15",
			wantTake: 15,
			wantSkip: 0,
		},
		{
			name:     "Quantity above 1000 should default to 10",
			page:     "1",
			quantity: "2000",
			wantTake: 10,
			wantSkip: 0,
		},
		{
			name:     "Invalid page number",
			page:     "invalid",
			quantity: "10",
			wantTake: 10,
			wantSkip: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queryParams := map[string]string{}
			if tt.page != "" {
				queryParams["page"] = tt.page
			}
			if tt.quantity != "" {
				queryParams["quantity"] = tt.quantity
			}

			c := setupEchoContext(queryParams)

			gotTake, gotSkip := routeutils.GetPagingParams(c, "", "")

			assert.Equal(t, tt.wantTake, gotTake)
			assert.Equal(t, tt.wantSkip, gotSkip)
		})
	}
}

func TestStringConverter(t *testing.T) {
	result, err := routeutils.StringConverter("test")
	assert.NoError(t, err)
	assert.Equal(t, "test", result)
}

func TestInt64Converter(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantValue int64
		wantError bool
	}{
		{
			name:      "Valid int64",
			input:     "123",
			wantValue: 123,
			wantError: false,
		},
		{
			name:      "Negative int64",
			input:     "-456",
			wantValue: -456,
			wantError: false,
		},
		{
			name:      "Invalid int64",
			input:     "invalid",
			wantValue: 0,
			wantError: true,
		},
		{
			name:      "Float value",
			input:     "123.45",
			wantValue: 0,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := routeutils.Int64Converter(tt.input)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantValue, got)
			}
		})
	}
}

func TestIntConverter(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantValue int
		wantError bool
	}{
		{
			name:      "Valid int",
			input:     "123",
			wantValue: 123,
			wantError: false,
		},
		{
			name:      "Negative int",
			input:     "-456",
			wantValue: -456,
			wantError: false,
		},
		{
			name:      "Invalid int",
			input:     "invalid",
			wantValue: 0,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := routeutils.IntConverter(tt.input)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantValue, got)
			}
		})
	}
}

func TestGetStringArrayQueryParam(t *testing.T) {
	tests := []struct {
		name       string
		paramValue string
		separator  string
		wantResult []string
	}{
		{
			name:       "Valid string array",
			paramValue: "a,b,c",
			separator:  ",",
			wantResult: []string{"a", "b", "c"},
		},
		{
			name:       "Empty parameter",
			paramValue: "",
			separator:  ",",
			wantResult: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queryParams := map[string]string{}
			if tt.paramValue != "" {
				queryParams["test"] = tt.paramValue
			}

			c := setupEchoContext(queryParams)

			got := routeutils.GetStringArrayQueryParam(c, "test", tt.separator)
			assert.Equal(t, tt.wantResult, got)
		})
	}
}

func TestGetInt64ArrayQueryParam(t *testing.T) {
	tests := []struct {
		name       string
		paramValue string
		separator  string
		wantResult []int64
		wantError  bool
	}{
		{
			name:       "Valid int64 array",
			paramValue: "1,2,3",
			separator:  ",",
			wantResult: []int64{1, 2, 3},
			wantError:  false,
		},
		{
			name:       "Invalid int64 in array",
			paramValue: "1,invalid,3",
			separator:  ",",
			wantResult: nil,
			wantError:  true,
		},
		{
			name:       "Empty parameter",
			paramValue: "",
			separator:  ",",
			wantResult: []int64{},
			wantError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queryParams := map[string]string{}
			if tt.paramValue != "" {
				queryParams["test"] = tt.paramValue
			}

			c := setupEchoContext(queryParams)

			got, err := routeutils.GetInt64ArrayQueryParam(c, "test", tt.separator)

			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantResult, got)
			}
		})
	}
}

func TestGetIntArrayQueryParam(t *testing.T) {
	tests := []struct {
		name       string
		paramValue string
		separator  string
		wantResult []int
		wantError  bool
	}{
		{
			name:       "Valid int array",
			paramValue: "1,2,3",
			separator:  ",",
			wantResult: []int{1, 2, 3},
			wantError:  false,
		},
		{
			name:       "Invalid int in array",
			paramValue: "1,invalid,3",
			separator:  ",",
			wantResult: nil,
			wantError:  true,
		},
		{
			name:       "Empty parameter",
			paramValue: "",
			separator:  ",",
			wantResult: []int{},
			wantError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queryParams := map[string]string{}
			if tt.paramValue != "" {
				queryParams["test"] = tt.paramValue
			}

			c := setupEchoContext(queryParams)

			got, err := routeutils.GetIntArrayQueryParam(c, "test", tt.separator)

			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantResult, got)
			}
		})
	}
}

func TestGetBoolQueryParam(t *testing.T) {
	tests := []struct {
		name       string
		paramValue string
		wantResult bool
	}{
		{
			name:       "Valid true",
			paramValue: "true",
			wantResult: true,
		},
		{
			name:       "Valid false",
			paramValue: "false",
			wantResult: false,
		},
		{
			name:       "Valid 1 (true)",
			paramValue: "1",
			wantResult: true,
		},
		{
			name:       "Valid 0 (false)",
			paramValue: "0",
			wantResult: false,
		},
		{
			name:       "Invalid value defaults to false",
			paramValue: "invalid",
			wantResult: false,
		},
		{
			name:       "Empty parameter defaults to false",
			paramValue: "",
			wantResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queryParams := map[string]string{}
			if tt.paramValue != "" {
				queryParams["test"] = tt.paramValue
			}

			c := setupEchoContext(queryParams)

			got := routeutils.GetBoolQueryParam(c, "test")
			assert.Equal(t, tt.wantResult, got)
		})
	}
}

func TestGetArrayParam(t *testing.T) {
	// String converter for testing
	stringConverter := func(s string) (string, error) {
		if s == "error" {
			return "", errors.New("conversion error")
		}
		return s, nil
	}

	// Int64 converter for testing
	int64Converter := func(s string) (int64, error) {
		if s == "invalid" {
			return 0, errors.New("invalid int64")
		}
		return routeutils.Int64Converter(s)
	}

	tests := []struct {
		name       string
		rawValue   string
		separator  string
		converter  interface{} // We'll cast this in the test
		wantResult interface{}
		wantError  bool
	}{
		{
			name:       "Valid comma-separated strings",
			rawValue:   "feedback,one_on_one,observation",
			separator:  ",",
			converter:  stringConverter,
			wantResult: []string{"feedback", "one_on_one", "observation"},
			wantError:  false,
		},
		{
			name:       "Valid semicolon-separated strings",
			rawValue:   "1;2;3",
			separator:  ";",
			converter:  stringConverter,
			wantResult: []string{"1", "2", "3"},
			wantError:  false,
		},
		{
			name:       "Values with whitespace should be trimmed",
			rawValue:   "  feedback  , one_on_one  ,  observation  ",
			separator:  ",",
			converter:  stringConverter,
			wantResult: []string{"feedback", "one_on_one", "observation"},
			wantError:  false,
		},
		{
			name:       "Empty string should return empty slice",
			rawValue:   "",
			separator:  ",",
			converter:  stringConverter,
			wantResult: []string{},
			wantError:  false,
		},
		{
			name:       "Whitespace only should return empty slice",
			rawValue:   "   ",
			separator:  ",",
			converter:  stringConverter,
			wantResult: []string{},
			wantError:  false,
		},
		{
			name:       "Single value",
			rawValue:   "feedback",
			separator:  ",",
			converter:  stringConverter,
			wantResult: []string{"feedback"},
			wantError:  false,
		},
		{
			name:       "Empty values should be filtered out",
			rawValue:   "feedback,,observation,",
			separator:  ",",
			converter:  stringConverter,
			wantResult: []string{"feedback", "observation"},
			wantError:  false,
		},
		{
			name:       "Converter error should return error",
			rawValue:   "feedback,error,observation",
			separator:  ",",
			converter:  stringConverter,
			wantResult: nil,
			wantError:  true,
		},
		{
			name:       "Valid int64 array",
			rawValue:   "1,2,3",
			separator:  ",",
			converter:  int64Converter,
			wantResult: []int64{1, 2, 3},
			wantError:  false,
		},
		{
			name:       "Invalid int64 in array should return error",
			rawValue:   "1,invalid,3",
			separator:  ",",
			converter:  int64Converter,
			wantResult: nil,
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch conv := tt.converter.(type) {
			case func(string) (string, error):
				got, err := routeutils.GetArrayParam(tt.rawValue, tt.separator, conv)
				if tt.wantError {
					assert.Error(t, err)
					assert.Nil(t, got)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tt.wantResult, got)
				}
			case func(string) (int64, error):
				got, err := routeutils.GetArrayParam(tt.rawValue, tt.separator, conv)
				if tt.wantError {
					assert.Error(t, err)
					assert.Nil(t, got)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tt.wantResult, got)
				}
			}
		})
	}
}

func TestGetTakeSkipFromPageQuantity(t *testing.T) {
	tests := []struct {
		name     string
		page     int64
		quantity int64
		wantTake int64
		wantSkip int64
	}{
		{
			name:     "First page",
			page:     1,
			quantity: 10,
			wantTake: 10,
			wantSkip: 0,
		},
		{
			name:     "Second page",
			page:     2,
			quantity: 10,
			wantTake: 10,
			wantSkip: 10,
		},
		{
			name:     "Page below 1 should default to 1",
			page:     0,
			quantity: 10,
			wantTake: 10,
			wantSkip: 0,
		},
		{
			name:     "Quantity below 1 should default to 10",
			page:     1,
			quantity: 0,
			wantTake: 10,
			wantSkip: 0,
		},
		{
			name:     "Quantity above 1000 should default to 10",
			page:     1,
			quantity: 2000,
			wantTake: 10,
			wantSkip: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTake, gotSkip := routeutils.GetTakeSkipFromPageQuantity(tt.page, tt.quantity)
			assert.Equal(t, tt.wantTake, gotTake)
			assert.Equal(t, tt.wantSkip, gotSkip)
		})
	}
}
