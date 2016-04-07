package utils

import (
	"fmt"
	"io"
	"net/http"
	
	"testsafeharbor/rest"
)

type DockerRegistry struct {
	RestContext
}

/*******************************************************************************
 * 
 */
func OpenDockerRegistryConnection(host string, port int, userId string,
	password string) (*DockerRegistry, error) {
	
	var registry *DockerRegistry = &DockerRegistry{
		RestContext: CreateRestContext(false, host, port, userId, password, noop),
	}
	
	var err error = registry.Ping()
	if err != nil {
		return nil, err
	}
	
	return registry, nil
}

/*******************************************************************************
 * 
 */
func (registry *DockerRegistry) Close() {
	
	registry.Host = ""
	registry.Port = 0
	registry.UserId = ""
	registry.Password = ""
}

/*******************************************************************************
 * 
 */
func (registry *DockerRegistry) Ping() error {
	
	var uri = "v2/"
	var response *http.Response
	var err error
	response, err = registry.sendBasicGet(uri)
	if err != nil { return err }
	if response.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Ping returned status: %s", response.Status))
	}
	return nil
}

/*******************************************************************************
 * If the specified image exists, return true. The image name is the image path
 * of the image namespace and registry repository name, separated by a "/".
 */
func (registry *DockerRegistry) ImageExists(name string, tag string) (bool, error) {
	
	// https://github.com/docker/distribution/blob/master/docs/spec/api.md
	// https://docs.docker.com/apidocs/v1.4.0/#!/repositories/GetRepository
	var uri = "/v2/" + name + "/manifests/" + tag
	//v0: GET /api/v0/repositories/{namespace}/{reponame}
	// Make HEAD request to registry.
	var response *http.Response
	var err error
	response, err = registry.sendBasicHead(uri)
	if err != nil { return err }
	if response.StatusCode == 200 {
		return true, nil
	} else if response.StatusCode == 404 { // Not Found
		return false, nil
	} else {
		return false, errors.New(fmt.Sprintf("ImageExists returned status: %s", response.Status))
	}
}

/*******************************************************************************
 * 
 */
func (registry *DockerRegistry) GetImage(name string, tag string, filepath string) error {
	
	// GET /v2/<name>/manifests/<reference>
	// GET /v2/<name>/blobs/<digest>
	
	// Retrieve manifest.
	var uri = "/v2/" + name + "/manifests/" + tag
	var resp *http.Response
	var err error
	resp, err = registry.sendBasicGet(uri)
	if err != nil { return err }
	resp.Body.Close()
	if resp.StatusCode == 200 {

		var responseMap map[string]interface{}
		responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
		resp.Body.Close()
		var layersObj = responseMap["fsLayers"]
		if layersObj == nil {
			return errors.New("Did not find fsLayers field in response")
		}
		var bool isType
		var layerAr []map[string]interface{}
		layerAr, isType = layersObj.([]map[string]interface{})
		if ! isType { return errors.New(
			"Type of layer description is " + reflect.TypeOf(layersObj).String())
		}
		
		// Retrieve layers.
		for _, layerDesc := range layerAr {
			var layerDigest = layerDesc["blobSum"]
			if layerDigest == nil {
				return errors.New("Did not find blobSum field in response for layer")
			}
			var digest string
			digest, isType = layerDigest.(string)
			if ! isType { return errors.New("blogSum field is not a string - it is a " +
				reflect.TypeOf(layerDigest)
			}
			uri = "/v2/" + name + "/blobs/" + digest
			response, err = registry.sendBasicGet(uri)
			
			
			
			....
			
		
			if ! testContext.Verify200Response(resp) { testContext.FailTest() }
			
			var reader io.ReadCloser = resp.Body
			var file *os.File
			file, err = os.Create(filename)
			testContext.AssertErrIsNil(err, "")
			_, err = io.Copy(file, reader)
			testContext.AssertErrIsNil(err, "")
			var fileInfo os.FileInfo
			fileInfo, err = file.Stat()
			if ! testContext.AssertErrIsNil(err, "") { return }
			testContext.AssertThat(fileInfo.Size() > 0, "File has zero size")
			resp.Body.Close()
		}
		
		
		
	} else if response.StatusCode == 404 { // Not Found
		return false, nil
	} else {
		return false, errors.New(fmt.Sprintf("ImageExists returned status: %s", response.Status))
	}
}

/*******************************************************************************
 * 
 */
func (registry *DockerRegistry) DeleteImage() error {
	
	v2: DELETE /v2/<name>/blobs/<digest>
		DELETE /v2/<name>/manifests/<reference>
	v1: DELETE /api/v0/repositories/{namespace}/{reponame}
}

/*******************************************************************************
 * 
 */
func noop(req *http.Request, s string) {
}
