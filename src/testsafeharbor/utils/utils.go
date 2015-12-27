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
	"crypto/sha512"
	"hash"
	
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
	CurrentTestPassed bool
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
	testContext.CurrentTestPassed = false
	testContext.TestStatus[hashKey] = ""
	fmt.Println()
	fmt.Println(testNumber, "Begin Test", name, "-------------------------------------------")
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) PassTest() {
	testContext.TestStatus[testContext.testName] = "Pass"
	testContext.CurrentTestPassed = true
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
func (testContext *TestContext) AssertErrIsNil(err error, msg string) bool {
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

/*******************************************************************************
 * Utility to determine if an array contains a specified value.
 */
func ContainsString(ar []string, val string) bool {
	for _, v := range ar {
		if v == val { return true }
	}
	return false
}

/*******************************************************************************
 * 
 */
func ComputeFileSignature(filepath string) ([]byte, error) {
	
	var file *os.File
	var err error
	file, err = os.Open(filepath)
	if err != nil { return nil, err }
	var numBytesRead int
	var buf []byte = make([]byte, 100000)
	var hash hash.Hash = sha512.New()
	for {
		numBytesRead, err = file.Read(buf)
		if numBytesRead == 0 { break }
		hash.Write(buf)
		if err != nil { break }
		if numBytesRead < 100000 { break }
	}
	
	return hash.Sum(make([]byte, 20)), nil
}
