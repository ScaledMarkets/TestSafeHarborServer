package rest

import (
	"net/http"
	"mime/multipart"
	"fmt"
	"net/url"
	"io"
	"os"
	"strings"
	"bytes"
	"encoding/json"
)

type RestContext struct {
	httpClient *http.Client
	ssl bool
	hostname string
	port int
	UserId string
	Password string
	setSessionId func(request *http.Request, id string)
}

/*******************************************************************************
 * userId and password are optional.
 */
func CreateRestContext(ssl bool, hostname string, port int, userId string, password string,
	sessionIdSetter func(*http.Request, string)) *RestContext {
	return &RestContext{
		httpClient: &http.Client{
			Transport: &http.Transport{
				DisableCompression: true,
			},
		},
		ssl: ssl,
		hostname: hostname,
		port: port,
		UserId: userId,
		Password: password,
		setSessionId: sessionIdSetter,
	}
}

func (restContext *RestContext) Print() {
	fmt.Println("RestContext:")
	fmt.Println(fmt.Sprintf("\thostname: %s", restContext.hostname))
	fmt.Println(fmt.Sprintf("\tport: %d", restContext.port))
}

func (restContext *RestContext) GetScheme() string {
	if restContext.ssl { return "https" } else { return "http" }
}

func (restContext *RestContext) GetHostname() string { return restContext.hostname }

func (restContext *RestContext) GetPort() int { return restContext.port }

func (restContext *RestContext) GetUserId() string { return restContext.UserId }

func (restContext *RestContext) getPassword() string { return restContext.Password }

/*******************************************************************************
 * Send a GET request to the SafeHarborServer, at the specified REST endpoint method
 * (reqName), with the specified query parameters, using basic authentication.
 */
func (restContext *RestContext) SendBasicGet(reqName string) (*http.Response, error) {
	
	var urlstr string = restContext.getURL(reqName)
	var resp *http.Response
	resp, err = restContext.httpClient.Get(urlstr)
	if err != nil { return nil, err }
	return resp, nil
}

/*******************************************************************************
 * Send a HEAD request to the SafeHarborServer, at the specified REST endpoint method
 * (reqName), with the specified query parameters, using basic authentication.
 */
func (restContext *RestContext) SendBasicHead(reqName string) (*http.Response, error) {
	
	var urlstr string = restContext.getURL(reqName)
	var resp *http.Response
	resp, err = restContext.httpClient.Head(urlstr)
	if err != nil { return nil, err }
	return resp, nil
}

/*******************************************************************************
 * Send a DELETE request to the SafeHarborServer, at the specified REST endpoint method
 * (reqName), with the specified query parameters, using basic authentication.
 */
func (restContext *RestContext) SendBasicDelete(reqName string) (*http.Response, error) {
	
	var urlstr string = restContext.getURL(reqName)
	var resp *http.Response
	var request *http.Request
	var err error
	request, err = http.NewRequest(reqMethod, urlstr)
	if err != nil { return nil, err }
	resp, err = restContext.httpClient.Do(request)
	if err != nil { return nil, err }
	return resp, nil
}

/*******************************************************************************
 * Send a POST request to the SafeHarborServer, at the specified REST endpoint method
 * (reqName), with the specified query parameters, using basic authentication.
 */
func (restContext *RestContext) SendBasicPost(reqName string, names []string,
	values []string) (*http.Response, error) {
	
	var urlstr string = restContext.getURL(reqName)
	var data = make(map[string][]string)
	for i, value := range values { data[names[i]] = make {string[], value } }
	var resp *http.Response
	resp, err = restContext.httpClient.PostForm(urlstr, data)
	if err != nil { return nil, err }
	return resp, nil
}

/*******************************************************************************
 * Send request as a multi-part so that a file can be attached. Use basic authentication.
 */
func (restContext *RestContext) SendBasicFilePost(reqName string, names []string,
	values []string, path string) (*http.Response, error) {

	var urlstr string = restContext.getURL(reqName)
	var resp *http.Response
	resp, err = restContext.httpClient.Post(url string, bodyType string, body io.Reader)
	if err != nil { return nil, err }
	return resp, nil
}

/*******************************************************************************
 * Send a GET request to the SafeHarborServer, at the specified REST endpoint method
 * (reqName), with the specified query parameters, using the specified session Id.
 */
func (restContext *RestContext) SendSessionGet(sessionId string, reqName string, names []string,
	values []string) (*http.Response, error) {

	return restContext.sendReq(sessionId, "GET", reqName, names, values)
}

/*******************************************************************************
 * Send an HTTP POST formatted according to what is required by the SafeHarborServer
 * REST API, as defined in the slides "SafeHarbor REST API" of the design,
 * https://drive.google.com/open?id=1r6Xnfg-XwKvmF4YppEZBcxzLbuqXGAA2YCIiPb_9Wfo
 * Use the specified session Id.
 */
func (restContext *RestContext) SendSessionPost(sessionId string, reqName string, names []string,
	values []string) (*http.Response, error) {

	return restContext.sendReq(sessionId, "POST", reqName, names, values)
}

/*******************************************************************************
 * Send an HTTP POST formatted according to what is required by the SafeHarborServer
 * REST API, as defined in the slides "SafeHarbor REST API" of the design,
 * https://drive.google.com/open?id=1r6Xnfg-XwKvmF4YppEZBcxzLbuqXGAA2YCIiPb_9Wfo
 */
func (restContext *RestContext) sendSessionReq(sessionId string, reqMethod string,
	reqName string, names []string, values []string) (*http.Response, error) {

	// Send REST POST request to server.
	var urlstr string = restContext.getURL(reqName)
	var data url.Values = url.Values{}
	for index, each := range names {
		data[each] = []string{values[index]}
	}
	var reader io.Reader = strings.NewReader(data.Encode())
	var request *http.Request
	var err error
	request, err = http.NewRequest(reqMethod, urlstr, reader)
	if err != nil { return nil, err }
	if reqMethod == "POST" {
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	if sessionId != "" {
		restContext.setSessionId(request, sessionId)
	}
	
	var resp *http.Response
	resp, err = restContext.httpClient.Do(request)
	if err != nil { return nil, err }
	return resp, nil
}

/*******************************************************************************
 * Send request as a multi-part so that a file can be attached.
 */
func (restContext *RestContext) SendSessionFilePost(sessionId string, reqName string, names []string,
	values []string, path string) (*http.Response, error) {

	var urlstr string = restContext.getURL(reqName)

	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	
	// Add file
	f, err := os.Open(path)
	if err != nil { return nil, err }
	var fileInfo os.FileInfo
	fileInfo, err = f.Stat()
	if err != nil { return nil, err }
	fw, err := w.CreateFormFile("filename", fileInfo.Name())
	if err != nil { return nil, err }
	_, err = io.Copy(fw, f)
	if err != nil { return nil, err }
	
	// Add the other fields
	if names != nil {
		for index, each := range names {
			fw, err = w.CreateFormField(each)
			if err != nil { return nil, err }
			_, err = fw.Write([]byte(values[index]))
			if err != nil { return nil, err }
		}
	}
	
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", urlstr, &b)
	if err != nil { return nil, err }
	
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())
	if sessionId != "" { restContext.setSessionId(req, sessionId) }
	//if sessionId != "" { req.Header.Set("Session-Id", sessionId) }

	// Submit the request
	res, err := restContext.httpClient.Do(req)
	if err != nil { return nil, err }

	return res, nil
}

/*******************************************************************************
 * Parse an HTTP JSON response that can be converted to a map.
 */
func ParseResponseBodyToMap(body io.ReadCloser) (map[string]interface{}, error) {
	var value []byte = ReadResponseBody(body)
	var obj map[string]interface{}
	err := json.Unmarshal(value, &obj)
	//var dec *json.Decoder = json.NewDecoder(body)
	//err := dec.Decode(&obj)
	if err != nil { return nil, err }
	//AssertErrIsNil(err, "When unmarshalling obj")
	
	
	//var result map[string]interface{}
	//var isType bool
	
	//result, _ = obj.(map[string]interface{})
	//AssertThat(isType, "Wrong type: obj is not a map[string]interface{}")
	return obj, nil
}

/*******************************************************************************
 * Parse an HTTP JSON response that can be converted to an array of maps.
 */
func ParseResponseBodyToMaps(body io.ReadCloser) ([]map[string]interface{}, error) {
	var value []byte = ReadResponseBody(body)
	var obj []map[string]interface{}
	err := json.Unmarshal(value, &obj)
	if err != nil { return nil, err }
	//var isType bool
	//result, isType = obj.([]map[string]interface{})
	//if ! isType {
	//	fmt.Println("This is what was returned:")
	//	fmt.Println(result)
	//}
	//AssertThat(isType, "Wrong type: obj is not a []map[string]interface{} - it is a " + 
	//	fmt.Sprintf("%s", reflect.TypeOf(obj)))
	return obj, nil
}

/*******************************************************************************
 * Parse an arbitrary HTTP JSON response.
 */
func ReadResponseBody(body io.ReadCloser) []byte {
	
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
func PrintMap(m map[string]interface{}) {
	fmt.Println("Map:")
	for k, v := range m {
		fmt.Println(k, v)
	}
}

/*******************************************************************************
 * If the response is not 200, then throw an exception.
 */
func (restContext *RestContext) Verify200Response(resp *http.Response) bool {
	var is200 bool = true
	if resp.StatusCode != 200 {
		is200 = false
		fmt.Sprintf("Response code %d", resp.StatusCode)
		var responseMap map[string]interface{}
		var err error
		responseMap, err = ParseResponseBodyToMap(resp.Body)
		if err == nil { PrintMap(responseMap) }
		//if restContext.stopOnFirstError { os.Exit(1) }
	}
	fmt.Println("Response code ", resp.StatusCode)
	return is200
}

/*******************************************************************************
 * 
 */
func (restContext *RestContext) getURL(reqName string) {
	var basicAuthCreds = ""
	if restContext.UserId != "" {
		basicAuthCreds = fmt.Sprintf("%s:%s@", restContext.UserId, restContext.Password)
	}
	return fmt.Sprintf(
		"%s://%s%s:%d/%s",
		restContext.GetScheme(), basicAuthCreds, restContext.hostname, restContext.port, reqName)
}

/*******************************************************************************
 * 
 * Utility to encode an arbitrary string value, which might contain quotes and other
 * characters, so that it can be safely and securely transported as a JSON string value,
 * delimited by double quotes. Ref. http://json.org/.
 */
func EncodeStringForJSON(value string) string {
	// Replace each occurrence of double-quote and backslash with backslash double-quote
	// or backslash backslash, respectively.
	
	var encodedValue = value
	encodedValue = strings.Replace(encodedValue, "\"", "\\\"", -1)
	encodedValue = strings.Replace(encodedValue, "\\", "\\\\", -1)
	encodedValue = strings.Replace(encodedValue, "/", "\\/", -1)
	encodedValue = strings.Replace(encodedValue, "\b", "\\b", -1)
	encodedValue = strings.Replace(encodedValue, "\f", "\\f", -1)
	encodedValue = strings.Replace(encodedValue, "\n", "\\n", -1)
	encodedValue = strings.Replace(encodedValue, "\r", "\\r", -1)
	encodedValue = strings.Replace(encodedValue, "\t", "\\t", -1)
	return encodedValue
}

/*******************************************************************************
 * Reverse the encoding that is performed by EncodeStringForJSON.
 */
func DecodeStringFromJSON(encodedValue string) string {
	var decodedValue = encodedValue
	decodedValue = strings.Replace(decodedValue, "\\t", "\t", -1)
	decodedValue = strings.Replace(decodedValue, "\\r", "\r", -1)
	decodedValue = strings.Replace(decodedValue, "\\n", "\n", -1)
	decodedValue = strings.Replace(decodedValue, "\\f", "\f", -1)
	decodedValue = strings.Replace(decodedValue, "\\b", "\b", -1)
	decodedValue = strings.Replace(decodedValue, "\\/", "/", -1)
	decodedValue = strings.Replace(decodedValue, "\\\\", "\\", -1)
	decodedValue = strings.Replace(decodedValue, "\\\"", "\"", -1)
	return decodedValue
}
