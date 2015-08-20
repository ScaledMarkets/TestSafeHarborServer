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
 * 
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
 * 
 */
func verify200Response(resp *http.Response) {
	if resp.StatusCode != 200 {
		fmt.Println(fmt.Sprintf("Response code %s", resp.StatusCode))
		return
	}
	
	fmt.Println("Response code ", resp.StatusCode)
}

/*******************************************************************************
 * 
 */
func assertThat(condition bool, msg string) {
	if ! condition { panic(errors.New(fmt.Sprintf("ERROR: %s", msg))) }
}

/*******************************************************************************
 * 
 */
func printMap(m map[string]string) {
	fmt.Println("Map:")
	for k, v := range m {
		fmt.Println(k, v)
	}
}
