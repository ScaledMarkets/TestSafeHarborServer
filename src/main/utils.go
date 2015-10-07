/*******************************************************************************
 * Functions requried by main.go.
 */

package main

import (
	"fmt"
	"net/http"
	"net/url"
	"io"
	"os"
	//"path/filepath"
	"mime/multipart"
	//"bufio"
	"bytes"
	"strings"
	"errors"
	"encoding/json"
)

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
		assertErrIsNil(err, "")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if sessionId != "" { request.Header.Set("Session-Id", sessionId) }
	
	var resp *http.Response
	var tr *http.Transport = &http.Transport{
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	resp, err = client.Do(request)
	assertErrIsNil(err, "")
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
	assertErrIsNil(err, "Cannot open path: " + path)
	fw, err := w.CreateFormFile("filename", path)
	assertErrIsNil(err, "Cannot create form file: " + path)
	_, err = io.Copy(fw, f)
	assertErrIsNil(err, "Could not copy file")
	
	// Add the other fields
	for index, each := range names {
		fw, err = w.CreateFormField(each)
		assertErrIsNil(err, "Could not create form field, " + each)
		_, err = fw.Write([]byte(values[index]))
		assertErrIsNil(err, "Could not write to file; index=" + string(index))
	}
	
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", urlstr, &b)
	assertErrIsNil(err, "When creating a POST request for " + urlstr)
	
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())
	if sessionId != "" { req.Header.Set("Session-Id", sessionId) }

	// Submit the request
	client := &http.Client{}
	res, err := client.Do(req)
	assertErrIsNil(err, "When doing request")

	return res
}

/*******************************************************************************
 * Parse an HTTP JSON response that can be converted to a map.
 */
func parseResponseBodyToMap(body io.ReadCloser) map[string]interface{} {
	var obj interface{} = parseResponseBody(body)
	var result map[string]interface{}
	var isType bool
	
	result, isType = obj.(map[string]interface{})
	assertThat(isType, "Wrong type: obj is not a map[string]interface{}")
	return result
}

/*******************************************************************************
 * Parse an HTTP JSON response that can be converted to an array of maps.
 */
func parseResponseBodyToMaps(body io.ReadCloser) []map[string]interface{} {
	var obj interface{} = parseResponseBody(body)
	var result []map[string]interface{}
	var isType bool
	
	result, isType = obj.([]map[string]interface{})
	assertThat(isType, "Wrong type: obj is not a map[string]interface{}")
	return result
}

/*******************************************************************************
 * Parse an arbitrary HTTP JSON response.
 */
func parseResponseBody(body io.ReadCloser) interface{} {
	
	var dec *json.Decoder = json.NewDecoder(body)
	var obj interface{}
	err := dec.Decode(&obj)
	assertErrIsNil(err, "When unmarshalling obj")
	return obj
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
func verify200Response(resp *http.Response) {
	assertThat(resp.StatusCode == 200, fmt.Sprintf("Response code %d", resp.StatusCode))
	fmt.Println("Response code ", resp.StatusCode)
}

/*******************************************************************************
 * If the specified condition is not true, then thrown an exception with the message.
 */
func assertThat(condition bool, msg string) {
	if ! condition { panic(errors.New(fmt.Sprintf("ERROR: %s", msg))) }
}

/*******************************************************************************
 * 
 */
func assertErrIsNil(err error, msg string) {
	if err == nil { return }
	fmt.Print(msg)
	panic(err)
}
