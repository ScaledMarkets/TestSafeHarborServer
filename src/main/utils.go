/*******************************************************************************
 * Functions requried by main.go.
 */

package main

import (
	"fmt"
	"net/http"
	"net/url"
	"io"
	"bufio"
	"strings"
	"errors"
)


/*******************************************************************************
 * Send an HTTP POST formatted according to what is required by the SafeHarborServer
 * REST API, as defined in the slides "SafeHarbor REST API" of the design,
 * https://drive.google.com/open?id=1r6Xnfg-XwKvmF4YppEZBcxzLbuqXGAA2YCIiPb_9Wfo
 */
func sendPost(reqName string, names []string, values []string) *http.Response {
	// Send REST POST request to server.
	var urlstr string = fmt.Sprintf(
		"http://%s:%s/%s",
		"127.0.0.1", "6000", reqName)
	
	var data url.Values = url.Values{}
	for index, each := range names {
		data[each] = []string{values[index]}
	}
	var reader io.Reader = strings.NewReader(data.Encode())
	var request *http.Request
	var err error
	request, err = http.NewRequest("POST", urlstr, reader)
		if err != nil { panic(err) }
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	var resp *http.Response
	var tr *http.Transport = &http.Transport{
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	resp, err = client.Do(request)
	if err != nil { panic(err) }
	return resp
}

/*******************************************************************************
 * Retrieve name=value pairs from the HTTP response body.
 * See slide "API REST Binding" in
 * https://drive.google.com/open?id=1r6Xnfg-XwKvmF4YppEZBcxzLbuqXGAA2YCIiPb_9Wfo
 */
func parseResponseBody(body io.ReadCloser) map[string]string {
	//var reader *bufio.Reader = bufio.NewReader(body)
	//if reader == nil { panic(errors.New("reader is nil")) }
	var responseMap map[string]string = map[string]string{}
	scanner := bufio.NewScanner(body)
	for scanner.Scan() {
		var line string = scanner.Text()
		var tokens []string = strings.Split(line, "=")
		if len(tokens) != 2 { panic(errors.New(fmt.Sprintf("Ill-formatted response: %s", line))) }
		var name string = strings.Trim(tokens[0], " ")
		var value string = strings.Trim(tokens[1], " ")
		responseMap[name] = value
	}
	return responseMap
}


/*******************************************************************************
 * If the response is not 200, then throw an exception.
 */
func verify200Response(resp *http.Response) {
	assertThat(resp.StatusCode == 200, fmt.Sprintf("Response code %s", resp.StatusCode))
	fmt.Println("Response code ", resp.StatusCode)
}

/*******************************************************************************
 * If the specified condition is not true, then thrown an exception with the message.
 */
func assertThat(condition bool, msg string) {
	if ! condition { panic(errors.New(fmt.Sprintf("ERROR: %s", msg))) }
}

/*******************************************************************************
 * Write the specified map to stdout.
 */
func printMap(m map[string]string) {
	fmt.Println("Map:")
	for k, v := range m {
		fmt.Println(k, v)
	}
}