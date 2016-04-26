/*******************************************************************************
 * Interface for accessing a docker engine via its REST API.
 * Engine API:
 * https://github.com/docker/docker/blob/master/docs/reference/api/docker_remote_api.md
 */

package docker

import (
	"fmt"
	"io"
	"os"
	"io/ioutil"
	"net"
	"net/http"
	"archive/tar"
	//"errors"
	"path/filepath"
	"encoding/base64"
	
	"testsafeharbor/utils"
	"testsafeharbor/rest"
)

type DockerEngine struct {
	rest.RestContext
}

/*******************************************************************************
 * 
 */
func OpenDockerEngineConnection() (*DockerEngine, error) {

	var engine *DockerEngine = &DockerEngine{
		// https://docs.docker.com/engine/quickstart/#bind-docker-to-another-host-port-or-a-unix-socket
		// Note: When the SafeHarborServer container is run, it must mount the
		// /var/run/docker.sock unix socket in the container:
		//		-v /var/run/docker.sock:/var/run/docker.sock
		RestContext: *rest.CreateUnixRestContext(
			unixDial,
			"", "",
			func (req *http.Request, s string) {}),
	}
	
	fmt.Println("Attempting to ping the engine...")
	var err error = engine.Ping()
	if err != nil {
		return nil, err
	}
	
	return engine, nil
}

/*******************************************************************************
 * For connecting to docker''s unix domain socket.
 */
func unixDial(network, addr string) (conn net.Conn, err error) {
	return net.Dial("unix", "/var/run/docker.sock")
}

/*******************************************************************************
 * 
 */
func (engine *DockerEngine) Ping() error {
	
	var uri = "_ping"
	var response *http.Response
	var err error
	response, err = engine.SendBasicGet(uri)
	if err != nil { return err }
	if response.StatusCode != 200 {
		return utils.ConstructError(fmt.Sprintf("Ping returned status: %s", response.Status))
	}
	return nil
}

/*******************************************************************************
 * Retrieve a list of the images that the docker engine has.
 */
func (engine *DockerEngine) GetImages() ([]map[string]interface{}, error) {
	
	var uri = "/images/json?all=1"
	var response *http.Response
	var err error
	response, err = engine.SendBasicGet(uri)
	if err != nil { return nil, err }
	if response.StatusCode != 200 {
		return nil, utils.ConstructError(fmt.Sprintf("GetImages returned status: %s", response.Status))
	}
	var imageMaps []map[string]interface{}
	imageMaps, err = rest.ParseResponseBodyToMaps(response.Body)
	if err != nil { return nil, err }
	return imageMaps, nil
}

/*******************************************************************************
 * Retrieve info on the specified docker image. Return an error if the image
 * is not found.
 */
func (engine *DockerEngine) GetImage(imageName string) (map[string]interface{}, error) {
	
	var uri = fmt.Sprintf("/images/%s/json", imageName)
	var response *http.Response
	var err error
	response, err = engine.SendBasicGet(uri)
	if err != nil { return nil, err }
	if response.StatusCode != 200 {
		return nil, utils.ConstructError(fmt.Sprintf("GetImage returned status: %s", response.Status))
	}
	var imageMap map[string]interface{}
	imageMap, err = rest.ParseResponseBodyToMap(response.Body)
	if err != nil { return nil, err }
	return imageMap, nil
}

/*******************************************************************************
 * Invoke the docker engine to build the image defined by the specified contents
 * of the build directory, which presumably contains a dockerfile. The textual
 * response from the docker engine is returned.
 */
func (engine *DockerEngine) BuildImage(buildDirPath, imageFullName string,
	dockerfileName string) (string, error) {

	// https://docs.docker.com/engine/reference/api/docker_remote_api_v1.23/#build-image-from-a-dockerfile
	// POST /build HTTP/1.1
	//
	// {{ TAR STREAM }} (this is the contents of the "build context")
	
	// See also the docker command line code, in docker/vendor/src/github.com/docker/engine-api/client/image_build.go:
	// https://github.com/docker/docker/blob/7fd53f7c711474791ce4292326e0b1dc7d4d6b0f/vendor/src/github.com/docker/engine-api/client/image_build.go
	
	// Create a temporary tar file of the build directory contents.
	var tarFile *os.File
	var err error
	var tempDirPath string
	tempDirPath, err = ioutil.TempDir("", "")
	if err != nil { return "", err }
	defer os.RemoveAll(tempDirPath)
	tarFile, err = ioutil.TempFile(tempDirPath, "")
	if err != nil { return "", utils.ConstructError(fmt.Sprintf(
		"When creating temp file '%s': %s", tarFile.Name(), err.Error()))
	}
	
	// Walk the build directory and add each file to the tar.
	var tarWriter = tar.NewWriter(tarFile)
	err = filepath.Walk(buildDirPath,
		func(path string, info os.FileInfo, err error) error {
		
			// Open the file to be written to the tar.
			if info.Mode().IsDir() { return nil }
			var new_path = path[len(buildDirPath):]
			if len(new_path) == 0 { return nil }
			var file *os.File
			file, err = os.Open(path)
			if err != nil { return err }
			defer file.Close()
			
			// Write tar header for the file.
			var header *tar.Header
			header, err = tar.FileInfoHeader(info, new_path)
			if err != nil { return err }
			header.Name = new_path
			err = tarWriter.WriteHeader(header)
			if err != nil { return err }
			
			// Write the file contents to the tar.
			_, err = io.Copy(tarWriter, file)
			if err != nil { return err }
			
			return nil  // success - file was written to tar.
		})
	
	if err != nil { return "", err }
	tarWriter.Close()
	
	// Send the request to the docker engine, with the tar file as the body content.
	var tarReader io.ReadCloser
	tarReader, err = os.Open(tarFile.Name())
	defer tarReader.Close()
	if err != nil { return "", err }
	var headers = make(map[string]string)
	headers["Content-Type"] = "application/tar"
	headers["X-Registry-Config"] = base64.URLEncoding.EncodeToString([]byte("{}"))
	var response *http.Response
	response, err = engine.SendBasicStreamPost(
		fmt.Sprintf("build?t=%s&dockerfile=%s", imageFullName, dockerfileName), headers, tarReader)
	defer response.Body.Close()
	if err != nil { return "", err }
	if response.StatusCode != 200 {
		fmt.Println("Response message: " + response.Status)
		return "", utils.ConstructError(response.Status)
	}
	
	var bytes []byte
	bytes, err = ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil { return "", err }
	var responseStr = string(bytes)
	
	return responseStr, nil
}

/*******************************************************************************
 * 
 */
func (engine *DockerEngine) TagImage(imageName, hostAndRepoName, tag string) error {
	
	var uri = fmt.Sprintf("images/%s/tag", imageName)
	var response *http.Response
	var err error
	var names = []string{ "repo", "tag" }
	var values = []string{ hostAndRepoName, tag }
	response, err = engine.SendBasicFormPost(uri, names, values)
	if err != nil { return err }
	if response.StatusCode != 201 {
		return utils.ConstructError(response.Status)
	}
	return nil
}


/*******************************************************************************
 * The imageFullName must be the full registry host:port/repo:tag name.
 */
func (engine *DockerEngine) PushImage(imageFullName, regUserId, regPass, regEmail string) error {
	
	var uri = fmt.Sprintf("images/%s/push", imageFullName)
	
	var parmNames = make([]string, 0)
	var parmValues = make([]string, 0)
	var headers = map[string]string{
		"X-Registry-Auth": fmt.Sprintf(
			"{\"username\": \"%s\", \"password\": \"%s\", \"email\": \"%s\"}",
			regUserId, regPass, regEmail),
	}
	
	var response *http.Response
	var err error
	response, err = engine.SendBasicFormPostWithHeaders(uri, parmNames, parmValues, headers)
	if err != nil { return err }
	if response.StatusCode != 200 {
		return utils.ConstructError(response.Status)
	}
	
	// Apr 25 20:46:25 ip-172-31-41-187.us-west-2.compute.internal docker[1092]:
	// time="2016-04-25T20:46:25.066856155Z" level=error
	// msg="Handler for POST /images/:0/localhost:5000/myimage:alpha/push returned error:
	// Error parsing reference: ":0/localhost:5000/myimage:alpha"
	// is not a valid repository/tag"

	return nil
}

/*******************************************************************************
 * 
 */
func (engine *DockerEngine) DeleteImage(imageName string) error {
	
	var uri = "/images/" + imageName
	var response *http.Response
	var err error
	response, err = engine.SendBasicDelete(uri)
	if err != nil { return err }
	if response.StatusCode != 200 {
		return utils.ConstructError(response.Status)
	}
	return nil
}
