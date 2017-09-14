/*******************************************************************************
 * Functions requried by main.go.
 */

package helpers

import (
	"flag"
	"fmt"
	"net/http"
	//"net/url"
	//"io"
	"io/ioutil"
	"os"
	//"path/filepath"
	//"mime/multipart"
	//"bufio"
	//"bytes"
	//"strings"
	//"errors"
	//"encoding/json"
	//"reflect"
	"crypto/sha256"
	"crypto/sha512"
	"runtime/debug"	
	
	// My packages:
	"utilities/rest"
	"utilities/utils"
)

type TestContext struct {
	rest.RestContext
	SessionId string
	IsAdmin bool
	testName string
	StopOnFirstError bool
	TestStatus map[string]string
	CurrentTestPassed bool
	NoOfTests int
	NoOfTestsThatFailed int
	RedisPswd string
	NoLargeFileTransfers bool
}

func NewTestContext(scheme, hostname string, port int,
	setSessionId func(req *http.Request, sessionId string),
	stopOnFirstError bool, redisPswd string, nolargefiles bool) *TestContext {

	return &TestContext{
		RestContext: *rest.CreateTCPRestContext(scheme, hostname, port, "", "", setSessionId),
		SessionId: "",
		StopOnFirstError: stopOnFirstError,
		TestStatus:  make(map[string]string),
		NoOfTests:  0,
		NoOfTestsThatFailed: 0,
		RedisPswd: redisPswd,
		NoLargeFileTransfers: nolargefiles,
	}
}

func (testContext *TestContext) Print() {
	testContext.RestContext.Print()
	fmt.Println("TestContext:")
	fmt.Println(fmt.Sprintf("\tSessionId: %s", testContext.SessionId))
	fmt.Println(fmt.Sprintf("\tIsAdmin: %v", testContext.IsAdmin))
	fmt.Println(fmt.Sprintf("\ttestName: %s", testContext.testName))
	fmt.Println(fmt.Sprintf("\tStopOnFirstError: %v", testContext.StopOnFirstError))
	fmt.Println(fmt.Sprintf("\tCurrentTestPassed: %v", testContext.CurrentTestPassed))
	fmt.Println(fmt.Sprintf("\tNoOfTests: %d", testContext.NoOfTests))
	fmt.Println(fmt.Sprintf("\tNoOfTestsThatFailed: %d", testContext.NoOfTestsThatFailed))
	fmt.Println(fmt.Sprintf("\tRedisPswd: %s", testContext.RedisPswd))
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
func (testContext *TestContext) GetCurrentTestName() string {
	return testContext.testName
}

/*******************************************************************************
 * Write this line to the server''s stdout at the start of each test.
 */
func (testContext *TestContext) TestDemarcation() string {
	return "\n\n" + testContext.GetCurrentTestName() + "<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) StartTest(name string) {
	
	if testContext.StopOnFirstError && (testContext.NoOfTestsThatFailed > 0) {
		testContext.AbortAllTests(fmt.Sprintf("After test number %d, before test %s",
			testContext.NoOfTests, name))
	}
	testContext.NoOfTests++
	var testNumber = testContext.NoOfTests
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
func (testContext *TestContext) PassTestIfNoFailures() bool {
	if testContext.TestStatus[testContext.testName] == "" {
		testContext.CurrentTestPassed = true
		testContext.TestStatus[testContext.testName] = "Pass"
	}
	return testContext.CurrentTestPassed
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) FailTest() {
	if testContext.TestStatus[testContext.testName] == "Fail" { return }
	testContext.NoOfTestsThatFailed++
	testContext.TestStatus[testContext.testName] = "Fail"
	fmt.Println("Failed test", testContext.testName)
	fmt.Println("Stack trace:")
	debug.PrintStack()
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TestHasFailed() bool {
	return (testContext.TestStatus[testContext.testName] == "Fail")
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) AbortAllTests(msg string) {
	fmt.Println("Aborting tests: " + msg)
	os.Exit(1)
}

/*******************************************************************************
 * If the specified condition is not true, then print an error message.
 */
func (testContext *TestContext) AssertThat(condition bool, msg string) bool {
	if ! condition {
		testContext.FailTest()
		fmt.Println(fmt.Sprintf("ERROR: %s", msg))
	}
	return condition
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) AssertOKResponse(resp *http.Response) {
	if ! testContext.Verify200Response(resp) {
		testContext.FailTest()
		fmt.Println("Response status: " + resp.Status)
	}
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) AssertErrIsNil(err error, msg string) bool {
	if err == nil { return true }
	fmt.Println("Original error message:", err.Error())
	fmt.Println("Supplemental message:", msg)
	testContext.FailTest()
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
func ComputeSHA256FileDigest(filepath string) ([]byte, error) {
	
	return utils.ComputeFileDigest(sha256.New(), filepath)
}

/*******************************************************************************
 * 
 */
func ComputeSHA512FileDigest(filepath string) ([]byte, error) {
	
	return utils.ComputeFileDigest(sha512.New(), filepath)
}

/*******************************************************************************
 * Create a temporary directory.
 */
func CreateTempDir() (string, error) {
	var path string
	var err error
	path, err = ioutil.TempDir("", "")
	if err != nil { return "", err }
	fmt.Println("Creating test directory " + path)
	return path, err
}

/*******************************************************************************
 * Create a temporary file in the specified directory, with the given name,
 * write the given content to it, and return the path to the file.
 */
func CreateTempFile(dir, name string, content string) (string, error) {
	var file *os.File
	var err error
	file, err = os.Create(dir + "/" + name)
	if err != nil { return "", err }
	var bytes []byte = []byte(content)
	var mode os.FileMode = os.ModeTemporary | os.ModePerm
	err = ioutil.WriteFile(file.Name(), bytes, mode)
	return file.Name(), err
}

/*******************************************************************************
 * Retrieve a file at the specified URL and save it to the specified path.
 * If the file already exists at that path, and useCachedFile is true, then
 * merely return.
 */
func DownloadFile(url string, finalPath string, useCachedFile bool) error {
	
	var err error
	
	if useCachedFile {
		var file *os.File
		file, err = os.Open(finalPath)
		if err != nil {
			_, err = file.Stat()
			if err == nil { return nil }  // file exists
		}
	}
	
	var resp *http.Response
	resp, err = http.Get(url)
	if err != nil { return err }
	
	var bytes []byte
	bytes, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil { return err }
	
	err = ioutil.WriteFile(finalPath, bytes, os.ModePerm)
	return err
}

/*******************************************************************************
 * Set the session Id as a cookie.
 */
func SetSessionId(req *http.Request, sessionId string) {
	
	// Set cookie containing the session Id.
	var cookie = &http.Cookie{
		Name: "SessionId",
		Value: sessionId,
		//Path: 
		//Domain: 
		//Expires: 
		//RawExpires: 
		MaxAge: 86400,  // 24 hrs
		Secure: false,  //....change to true later.
		HttpOnly: true,
		//Raw: 
		//Unparsed: 
	}
	
	req.AddCookie(cookie)
}

/*******************************************************************************
 * 
 */
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
	flag.PrintDefaults()
}
