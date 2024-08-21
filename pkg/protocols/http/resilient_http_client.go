package httpclient

import (
	"PaymentGateway/pkg/protocols/models"
	"bytes"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/eapache/go-resiliency/retrier"
	"io"
	"log"
	"net/http"
)

type ResilientHttpClient struct {
	client *http.Client
}

func NewResilientHttpClient() *ResilientHttpClient {
	return &ResilientHttpClient{client: http.DefaultClient}
}

func (h *ResilientHttpClient) Post(endpoint string, headers map[string]string, data []byte, resiliencyParameters models.ResiliencyParameters) ([]byte, error) {
	return h.callUsingCircuitBreakerWithRetries("POST", endpoint, headers, data, resiliencyParameters)
}

func (h *ResilientHttpClient) callUsingCircuitBreakerWithRetries(method string, url string, headers map[string]string,
	body []byte, resiliencyParameters models.ResiliencyParameters) ([]byte, error) {

	var responseBody []byte
	h.configureHystrix(resiliencyParameters)

	err := hystrix.Do(resiliencyParameters.UniqueCommandName, func() error {
		req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
		if err != nil {
			return err
		}

		for key, value := range headers {
			req.Header.Set(key, value)
		}

		err = h.retry(req, &responseBody, resiliencyParameters)
		return err
	}, func(err error) error {
		log.Printf("Fallback triggered for %s: %v", resiliencyParameters.UniqueCommandName, err)
		return err
	})

	if err != nil {
		return nil, err
	}
	return responseBody, nil
}

func (h *ResilientHttpClient) retry(req *http.Request, responseBody *[]byte, resiliencyParameters models.ResiliencyParameters) error {
	r := retrier.New(retrier.ConstantBackoff(resiliencyParameters.RetryTimes, resiliencyParameters.WaitBetweenRetriesInMilliSec), nil)

	attempt := 0
	err := r.Run(func() error {
		attempt++
		resp, err := h.client.Do(req)
		if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			defer resp.Body.Close()
			*responseBody, err = io.ReadAll(resp.Body)
			return nil
		} else if err == nil {
			err = fmt.Errorf("status code was %v %d", resp.Status, attempt)
		}
		return err
	})
	return err
}

func (h *ResilientHttpClient) configureHystrix(resiliencyParameters models.ResiliencyParameters) {
	hystrix.ConfigureCommand(resiliencyParameters.UniqueCommandName, hystrix.CommandConfig{
		Timeout:                resiliencyParameters.TimeoutInMilliSec,
		MaxConcurrentRequests:  resiliencyParameters.MaxConcurrentRequests,
		RequestVolumeThreshold: resiliencyParameters.RequestVolumeThreshold,
		SleepWindow:            resiliencyParameters.SleepWindowInMilliSec,
		ErrorPercentThreshold:  resiliencyParameters.ErrorPercentThreshold,
	})
}
