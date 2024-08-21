package httpclient

import (
	"PaymentGateway/pkg/protocols/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResilientHttpClient_Post_Success(t *testing.T) {
	// Set up a test server to mock the external API
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message": "success"}`))
	}))
	defer ts.Close()

	client := NewResilientHttpClient()

	headers := map[string]string{"Content-Type": "application/json"}
	data := []byte(`{"data": "test"}`)
	resiliencyParams := models.ResiliencyParameters{
		UniqueCommandName:            "test_command",
		TimeoutInMilliSec:            1000,
		RetryTimes:                   3,
		WaitBetweenRetriesInMilliSec: 100,
	}

	resp, err := client.Post(ts.URL, headers, data, resiliencyParams)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"message": "success"}`, string(resp))
}

// Todo other tests
