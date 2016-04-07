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
	response, err = registry.sendSessionReq("", "GET", uri, make{[]string, 0}, make{[]string, 0})
	if err != nil {
		return err
	}
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
	
	
	if 200 {
		Content-Length: <length of manifest>
		Docker-Content-Digest: <digest>
		return true, nil
	} else if 404 Not Found {
		return false, nil
	} else {
		return false, new error
	}
}

/*******************************************************************************
 * 
 */
func (registry *DockerRegistry) GetImage() error {
	
	GET /v2/<name>/manifests/<reference>
	GET /v2/<name>/blobs/<digest>
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
