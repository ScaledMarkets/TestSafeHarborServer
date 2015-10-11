/*******************************************************************************
 * Functions requried by main.go.
 */

package main

import (
	"fmt"
	"net/http"
	"net/url"
	"io"
	//"io/ioutil"
	"os"
	//"path/filepath"
	"mime/multipart"
	//"bufio"
	"bytes"
	"strings"
	//"errors"
	"encoding/json"
	//"reflect"
)

var testStatus map[string]string = make(map[string]string)
var noOfTests int = 0
var noOfTestsThatFailed int = 0

/*******************************************************************************
 * 
 */
func (testContext *TestContext) StartTest(name string) {
	testContext.testName = name
	testStatus[name] = ""
	noOfTests++
	fmt.Println()
	fmt.Println("Begin Test", name, "-------------------------------------------")
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) PassTest() {
	testStatus[testContext.testName] = "Pass"
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) FailTest() {
	noOfTestsThatFailed++
	testStatus[testContext.testName] = "Fail"
}

/*******************************************************************************
 * Send a GET request to the SafeHarborServer, at the specified REST endpoint method
 * (reqName), with the specified query parameters.
 */
func (testContext *TestContext) sendGet(sessionId string, reqName string, names []string,
	values []string) *http.Response {

	return testContext.sendReq(sessionId, "GET", reqName, names, values)
}

/*******************************************************************************
 * Send an HTTP POST formatted according to what is required by the SafeHarborServer
 * REST API, as defined in the slides "SafeHarbor REST API" of the design,
 * https://drive.google.com/open?id=1r6Xnfg-XwKvmF4YppEZBcxzLbuqXGAA2YCIiPb_9Wfo
 */
func (testContext *TestContext) sendPost(sessionId string, reqName string, names []string,
	values []string) *http.Response {

	return testContext.sendReq(sessionId, "POST", reqName, names, values)
}

/*******************************************************************************
 * Send an HTTP POST formatted according to what is required by the SafeHarborServer
 * REST API, as defined in the slides "SafeHarbor REST API" of the design,
 * https://drive.google.com/open?id=1r6Xnfg-XwKvmF4YppEZBcxzLbuqXGAA2YCIiPb_9Wfo
 */
func (testContext *TestContext) sendReq(sessionId string, reqMethod string,
	reqName string, names []string, values []string) *http.Response {

	// Send REST POST request to server.
	var urlstr string = fmt.Sprintf(
		"http://%s:%s/%s",
		testContext.hostname, testContext.port, reqName)
	
	var data url.Values = url.Values{}
	for index, each := range names {
		data[each] = []string{values[index]}
	}
	var reader io.Reader = strings.NewReader(data.Encode())
	var request *http.Request
	var err error
	request, err = http.NewRequest(reqMethod, urlstr, reader)
		testContext.assertErrIsNil(err, "")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if sessionId != "" { request.Header.Set("Session-Id", sessionId) }
	
	var resp *http.Response
	var tr *http.Transport = &http.Transport{
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	resp, err = client.Do(request)
	testContext.assertErrIsNil(err, "")
	return resp
}

/*******************************************************************************
 * Similar to sendPost, but send as a multi-part so that a file can be attached.
 */
func (testContext *TestContext) sendFilePost(sessionId string, reqName string, names []string,
	values []string, path string) *http.Response {

	var urlstr string = fmt.Sprintf(
		"http://%s:%s/%s",
		testContext.hostname, testContext.port, reqName)

	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	
	// Add file
	f, err := os.Open(path)
	testContext.assertErrIsNil(err, "Cannot open path: " + path)
	fw, err := w.CreateFormFile("filename", path)
	testContext.assertErrIsNil(err, "Cannot create form file: " + path)
	_, err = io.Copy(fw, f)
	testContext.assertErrIsNil(err, "Could not copy file")
	
	// Add the other fields
	for index, each := range names {
		fw, err = w.CreateFormField(each)
		testContext.assertErrIsNil(err, "Could not create form field, " + each)
		_, err = fw.Write([]byte(values[index]))
		testContext.assertErrIsNil(err, "Could not write to file; index=" + string(index))
	}
	
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", urlstr, &b)
	testContext.assertErrIsNil(err, "When creating a POST request for " + urlstr)
	
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())
	if sessionId != "" { req.Header.Set("Session-Id", sessionId) }

	// Submit the request
	client := &http.Client{}
	res, err := client.Do(req)
	testContext.assertErrIsNil(err, "When doing request")

	return res
}

/*******************************************************************************
 * Parse an HTTP JSON response that can be converted to a map.
 */
func parseResponseBodyToMap(body io.ReadCloser) map[string]interface{} {
	var value []byte = parseResponseBody(body)
	var obj map[string]interface{}
	err := json.Unmarshal(value, &obj)
	//var dec *json.Decoder = json.NewDecoder(body)
	//err := dec.Decode(&obj)
	if err != nil { fmt.Println(err.Error()) }
	//assertErrIsNil(err, "When unmarshalling obj")
	
	
	//var result map[string]interface{}
	//var isType bool
	
	//result, _ = obj.(map[string]interface{})
	//assertThat(isType, "Wrong type: obj is not a map[string]interface{}")
	return obj
}

/*******************************************************************************
 * Parse an HTTP JSON response that can be converted to an array of maps.
 */
func parseResponseBodyToMaps(body io.ReadCloser) []map[string]interface{} {
	var value []byte = parseResponseBody(body)
	var obj []map[string]interface{}
	err := json.Unmarshal(value, &obj)
	if err != nil { fmt.Println(err.Error()) }
	//var isType bool
	//result, isType = obj.([]map[string]interface{})
	//if ! isType {
	//	fmt.Println("This is what was returned:")
	//	fmt.Println(result)
	//}
	//assertThat(isType, "Wrong type: obj is not a []map[string]interface{} - it is a " + 
	//	fmt.Sprintf("%s", reflect.TypeOf(obj)))
	return obj
}

/*******************************************************************************
 * Parse an arbitrary HTTP JSON response.
 */
func parseResponseBody(body io.ReadCloser) []byte {
	
	var value []byte = make([]byte, 0)
	//var s string = ""
	for {
		var buf []byte = make([]byte, 100)
		n, err := body.Read(buf)
		if n > 0 { value = append(value, buf[0:n]...) }
		//fmt.Println("Read", string(buf))
		if err != nil { break }
		if n < len(buf) { break }
	}
	fmt.Println("Read this:")
	fmt.Println(string(value))
	fmt.Println("--")
	fmt.Println()
	
	return value
}

/*******************************************************************************
 * Write the specified map to stdout.
 */
func printMap(m map[string]interface{}) {
	fmt.Println("Map:")
	for k, v := range m {
		fmt.Println(k, v)
	}
}

/*******************************************************************************
 * If the response is not 200, then throw an exception.
 */
func (testContext *TestContext) verify200Response(resp *http.Response) {
	if resp.StatusCode != 200 {
		fmt.Sprintf("Response code %d", resp.StatusCode)
		var responseMap map[string]interface{}
		responseMap  = parseResponseBodyToMap(resp.Body)
		printMap(responseMap)
		//if testContext.stopOnFirstError { os.Exit(1) }
	}
	fmt.Println("Response code ", resp.StatusCode)
}

/*******************************************************************************
 * If the specified condition is not true, then print an error message.
 */
func (testContext *TestContext) assertThat(condition bool, msg string) {
	if ! condition {
		testContext.FailTest()
		fmt.Println(fmt.Sprintf("ERROR: %s", msg))
		if testContext.stopOnFirstError { os.Exit(1) }
	}
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) assertErrIsNil(err error, msg string) {
	if err == nil { return }
	testContext.FailTest()
	fmt.Print(msg)
	if testContext.stopOnFirstError { os.Exit(1) }
}
