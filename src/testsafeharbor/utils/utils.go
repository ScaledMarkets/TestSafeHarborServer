/*******************************************************************************
 * Functions requried by main.go.
 */

package utils

import (
	"fmt"
	//"net/http"
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
	//"testsafeharbor/rest"
)

var TestStatus map[string]string = make(map[string]string)
var NoOfTests int = 0
var NoOfTestsThatFailed int = 0


/*******************************************************************************
 * 
 */
func (testContext *TestContext) StartTest(name string) {
	
	testContext.testName = name
	TestStatus[name] = ""
	NoOfTests++
	fmt.Println()
	fmt.Println(NoOfTests, "Begin Test", name, "-------------------------------------------")
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) PassTest() {
	TestStatus[testContext.testName] = "Pass"
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) FailTest() {
	NoOfTestsThatFailed++
	TestStatus[testContext.testName] = "Fail"
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
	fmt.Print(msg)
	if testContext.StopOnFirstError { os.Exit(1) }
	return false
}

/*******************************************************************************
 * 
 */
func BoolToString(b bool) string {
	if b { return "true" } else { return "false" }	
}
