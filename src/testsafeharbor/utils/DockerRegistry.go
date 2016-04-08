package utils

/* Interface for interacting with a Docker Registry version 2.

	What a docker "name" is:

		(From: https://github.com/docker/distribution/blob/master/docs/spec/api.md)
		
		All endpoints will be prefixed by the API version and the repository name:
		
		/v2/<name>/
		
		For example, an API endpoint that will work with the library/ubuntu repository,
		the URI prefix will be:
		
		/v2/library/ubuntu/
		
		This scheme provides rich access control over various operations and methods
		using the URI prefix and http methods that can be controlled in variety of ways.
		
		Classically, repository names have always been two path components where each
		path component is less than 30 characters. The V2 registry API does not enforce this.
		The rules for a repository name are as follows:
		
			1. A repository name is broken up into path components. A component of a
				repository name must be at least one lowercase, alpha-numeric characters,
				optionally separated by periods, dashes or underscores. More strictly,
				it must match the regular expression [a-z0-9]+(?:[._-][a-z0-9]+)*.
				
			2. If a repository name has two or more path components, they must be
				separated by a forward slash ("/").
				
			3. The total length of a repository name, including slashes, must be
				less the 256 characters.
		
		These name requirements only apply to the registry API and should accept
		a superset of what is supported by other docker ecosystem components.
*/



import (
	"fmt"
	"io"
	"net/http"
	"archive/tar"
	
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
	if response.StatusCode == 404 {
		return errors.New("Not found")
	} else if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("ImageExists returned status: %s", response.Status))
	}
	
	// Parse description of each layer.
	var layerAr []map[string]interface{}
	layerAr, err = parseLayerDescriptions(resp.Body)
	if err != nil { return err }
	
	// Retrieve layers, and add each to a tar archive.
	var tarFile *os.File
	tarFile, err = os.Create(filepath)
	if err != nil { return errors.New(fmt.Sprintf(
		"When creating image file '%s': %s", filepath, err.Error()))
	}
	var tarWriter = tar.NewWriter(tarFile)
	var tempDirPath
	tempDirPath, err = ioutil.TempDir("", "")
	if err != nil { return errors.New(fmt.Sprintf(
		"When creating temp directory for writing layer files: %s", err.Error()))
	}
	defer os.RemoveAll(tempDirPath)
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
		if err != nil { return errors.New(fmt.Sprintf(
			"When requesting uri: '%s' - %s", uri, error.Error()))
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 { return errors.New(fmt.Sprintf(
			"Response code %d, when requesting uri: '%s' - %s", 
				resp.StatusCode, uri, error.Error()))
		}

		// Create temporary file in which to write layer.
		var layerfilePath *os.File
		layerfilePath, err = ioutil.TempFile(tempDirPath, digest)
		if err != nil { return errors.New(fmt.Sprintf(
			"When creating layer file '%s': %s", layerfilePath, err.Error()))
		}
		
		var reader io.ReadCloser = resp.Body
		var layerFile *os.File
		layerFile, err = os.OpenFile(layerfilePath, os.O_WRONLY, 0600)
		if err != nil { return errors.New(fmt.Sprintf(
			"When opening layer file '%s': %s", layerfilePath, err.Error()))
		}
		_, err = io.Copy(layerFile, reader)
		if err != nil { return errors.New(fmt.Sprintf(
			"When writing layer file '%s': %s", layerfilePath, err.Error()))
		}
		var fileInfo os.FileInfo
		fileInfo, err = layerFile.Stat()
		if err != nil { return errors.New(fmt.Sprintf(
			"When getting status of layer file '%s': %s", layerfilePath, err.Error()))
		}
		if fileInfo.Size() == 0 { return errors.New(fmt.Sprintf(
			"Layer file that was written, '%s', has zero size", layerfilePath))
		}
		
		// Add file to tar archive.
		var tarHeader = &tar.Header{
			Name: fileInfo.Name(),
			Mode: 0600,
			Size: int64(fileInfo.Size()),
		}
		err = tarWriter.WriteHeader(tarHeader)
		if err != nil {	return errors.New(fmt.Sprintf(
			"While writing layer header to tar archive: , %s", err.Error()))
		}
		
		layerFile, err = os.Open(layerfilePath)
		if err != nil {	return errors.New(fmt.Sprintf(
			"While opening layer file '%s': , %s", layerfilePath, err.Error()))
		_, err := io.Copy(tarWriter, layerFile)
		if err != nil {	return errors.New(fmt.Sprintf(
			"While writing layer content to tar archive: , %s", err.Error()))
		}
	}
	
	err = tarWriter.Close()
	if err != nil {	return errors.New(fmt.Sprintf(
		"While closing tar archive: , %s", err.Error()))
	}
	
	return nil
}

/*******************************************************************************
 * 
 */
func (registry *DockerRegistry) DeleteImage(name string) error {
	
	//v2: DELETE /v2/<name>/blobs/<digest>
	//	DELETE /v2/<name>/manifests/<reference>
	//v1: DELETE /api/v0/repositories/{namespace}/{reponame}
	
	// Retrieve manifest.
	var uri = "/v2/" + name + "/manifests/" + tag
	var resp *http.Response
	var err error
	resp, err = registry.sendBasicGet(uri)
	if err != nil { return err }
	resp.Body.Close()
	if response.StatusCode == 404 {
		return errors.New("Not found")
	} else if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("DeleteImage returned status: %s", response.Status))
	}
	
	// Parse description of each layer.
	var layerAr []map[string]interface{}
	layerAr, err = parseLayerDescriptions(resp.Body)
	if err != nil { return err }
	
	// Delete each layer.
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
		
		uri = fmt.Sprintf("/v2/%s/blobs/%s", name, digest)
		var response *http.Response
		var err error
		response, err = registry.SendBasicDelete(uri)
		if err != nil { return err }
		if response.StatusCode == 200 {
			return true, nil
		} else if response.StatusCode == 404 { // Not Found
			return false, nil
		} else {
			return false, errors.New(fmt.Sprintf("DeleteImage returned status: %s", response.Status))
		}
	}
	
	// Delete manifest.
	uri = "/v2/" + name + "/manifests/" + tag
	response, err = registry.SendBasicDelete(uri)
	if err != nil { return err }
	
	return nil
}

/*******************************************************************************
 * 
 */
func parseLayerDescriptions(body io.ReadCloser) ([]map[string]interface{}, error) {
	
	var responseMap map[string]interface{}
	var err error
	responseMap, err = rest.ParseResponseBodyToMap(body)
	body.Close()
	var layersObj = responseMap["fsLayers"]
	if layersObj == nil {
		return errors.New("Did not find fsLayers field in body")
	}
	var bool isType
	var layerAr []map[string]interface{}
	layerAr, isType = layersObj.([]map[string]interface{})
	if ! isType { return errors.New(
		"Type of layer description is " + reflect.TypeOf(layersObj).String())
	}
	return layerAr
}

/*******************************************************************************
 * 
 */
func noop(req *http.Request, s string) {
}
