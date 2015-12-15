/*******************************************************************************
 * Functions requried by main.go.
 */

package utils

import (
	"fmt"
	"net/http"
	//"net/url"
	//"io"
	//"io/ioutil"
	"os"
	//"path/filepath"
	//"mime/multipart"
	//"bufio"
	//"bytes"
	//"strings"
	//"errors"
	//"encoding/json"
	//"reflect"
	
	// My packages:
	"testsafeharbor/rest"
)

type TestContext struct {
	rest.RestContext
	SessionId string
	IsAdmin bool
	testName string
	StopOnFirstError bool
	PerformDockerTests bool
	TestStatus map[string]string
	NoOfTests int
	NoOfTestsThatFailed int
}

func NewTestContext(hostname, port string,
	setSessionId func(req *http.Request, sessionId string),
	stopOnFirstError, doNotPerformDockerTests bool) *TestContext {

	return &TestContext{
		RestContext: *rest.CreateRestContext(hostname, port, setSessionId),
		SessionId: "",
		StopOnFirstError: stopOnFirstError,
		PerformDockerTests: ! doNotPerformDockerTests,
		TestStatus:  make(map[string]string),
		NoOfTests:  0,
		NoOfTestsThatFailed: 0,
	}
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) GetTestsThatFailed() []string {
	var testsThatFailed = []string{}
	for test, status := range testContext.TestStatus {
		if status  == "Fail" { testsThatFailed = append(testsThatFailed, test) }
	}
	return testsThatFailed
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) StartTest(name string) {
	
	testContext.NoOfTests++
	var testNumber =testContext.NoOfTests
	var hashKey = fmt.Sprintf("%d: %s", testNumber, name)
	testContext.testName = hashKey
	testContext.TestStatus[hashKey] = ""
	fmt.Println()
	fmt.Println(testNumber, "Begin Test", name, "-------------------------------------------")
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) PassTest() {
	testContext.TestStatus[testContext.testName] = "Pass"
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) FailTest() {
	if testContext.TestStatus[testContext.testName] == "Fail" { return }
	testContext.NoOfTestsThatFailed++
	testContext.TestStatus[testContext.testName] = "Fail"
	fmt.Println("Failed test", testContext.testName)
}

/*******************************************************************************
 * If the specified condition is not true, then print an error message.
 */
func (testContext *TestContext) AssertThat(condition bool, msg string) bool {
	if ! condition {
		testContext.FailTest()
		fmt.Println(fmt.Sprintf("ERROR: %s", msg))
		if testContext.StopOnFirstError { os.Exit(1) }
	}
	return condition
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) assertErrIsNil(err error, msg string) bool {
	if err == nil { return true }
	testContext.FailTest()
	fmt.Print(msg, err.Error())
	if testContext.StopOnFirstError { os.Exit(1) }
	return false
}

/*******************************************************************************
 * 
 */
func BoolToString(b bool) string {
	if b { return "true" } else { return "false" }	
}
