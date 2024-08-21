package models

import "time"

type ResiliencyParameters struct {
	// UniqueCommandName use standard format - ServiceName_MethodName
	UniqueCommandName string
	// TimeoutInMilliSec is how long to wait for command to complete, in milliseconds
	TimeoutInMilliSec int
	// MaxConcurrentRequests is how many commands of the same type can run at the same time
	MaxConcurrentRequests int
	// SleepWindowInMilliSec is how long, in milliseconds, to wait after a circuit opens before testing for recovery
	SleepWindowInMilliSec int
	// ErrorPercentThreshold causes circuits to open once the rolling measure of errors exceeds this percent of requests
	ErrorPercentThreshold int
	// RequestVolumeThreshold is the minimum number of requests needed before a circuit can be tripped due to health
	RequestVolumeThreshold int
	// RetryTimes is number of calls to the api before handing over the error to the circuit breaker
	RetryTimes int
	// WaitBetweenRetriesInMilliSec is wait time between 2 retries in case of an error from the api
	WaitBetweenRetriesInMilliSec time.Duration
}

func DefaultResiliencyParameters(uniqueCommandName string) ResiliencyParameters {
	return ResiliencyParameters{
		UniqueCommandName:            uniqueCommandName,
		TimeoutInMilliSec:            2000,
		MaxConcurrentRequests:        10,
		SleepWindowInMilliSec:        5000,
		ErrorPercentThreshold:        50,
		RequestVolumeThreshold:       20,
		RetryTimes:                   3,
		WaitBetweenRetriesInMilliSec: 200 * time.Millisecond,
	}
}
