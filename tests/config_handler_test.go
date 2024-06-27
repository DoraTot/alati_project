// In your tests package (e.g., tests/config_handler_test.go)
package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"projekat/handlers"
	"projekat/model"
)

// Ensure MockConfigService implements services.ConfigServiceInterface
var _ model.ConfigRepository = (*MockConfigService)(nil)

// TestGetHandler tests the Get handler
func TestGetHandler(t *testing.T) {
	tracer := otel.Tracer("test-tracer")
	mockService := new(MockConfigService)
	handler := &handlers.ConfigHandler{
		Service: mockService,
		Tracer:  tracer,
	}

	router := mux.NewRouter()
	router.Handle("/config/{name}/{version}/", http.HandlerFunc(handler.Get)).Methods("GET")

	name := "exampleConfig"
	version := "1.0"
	versionFloat32, _ := strconv.ParseFloat(version, 32)
	version32 := float32(versionFloat32)

	config := &model.Config{
		Name:    name,
		Version: version32,
	}

	// Mock the service response
	mockService.On("GetConfig", name, version32, mock.Anything).Return(config, nil)

	req, err := http.NewRequest("GET", "/config/"+name+"/"+version+"/", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Verify the response
	require.Equal(t, http.StatusOK, rr.Code)
	var responseConfig model.Config
	err = json.Unmarshal(rr.Body.Bytes(), &responseConfig)
	require.NoError(t, err)
	require.Equal(t, config, &responseConfig)

	// Verify that the service's GetConfig method was called with the expected parameters
	mockService.AssertCalled(t, "GetConfig", name, version32, mock.Anything)
}
