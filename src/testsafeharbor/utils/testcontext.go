package utils

import (
	"testsafeharbor/rest"
)

type TestContext struct {
	rest.RestContext
	SessionId string
	IsAdmin bool
	testName string
	StopOnFirstError bool
	PerformDockerTests bool
}

