package utils

import (
	"fmt"
	"io"
	"os"
	"io/ioutil"
	"net/http"
	"archive/tar"
	"errors"
	"path/filepath"
	
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
		RestContext: *rest.CreateRestContext("unix",
			"/var/run/docker.sock", 0, "", "", func (req *http.Request, s string) {}),
	}
	
	var err error = engine.Ping()
	if err != nil {
		return nil, err
	}
	
	return engine, nil
}

/*******************************************************************************
 * 
 */
func (registry *DockerEngine) Ping() error {
	
	var uri = "_ping"
	var response *http.Response
	var err error
	response, err = registry.SendBasicGet(uri)
	if err != nil { return err }
	if response.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Ping returned status: %s", response.Status))
	}
	return nil
}


/*******************************************************************************
 * 
 */
func (engine *DockerEngine) BuildImage(buildDirPath, imageFullName string) error {

	// https://docs.docker.com/engine/reference/api/docker_remote_api_v1.23/#build-image-from-a-dockerfile
	// POST /build HTTP/1.1
	//
	// {{ TAR STREAM }} (this is the contents of the "build context")
	
	// Create a temporary tar file of the build directory contents.
	var tarFile *os.File
	var err error
	var tempDirPath string
	tempDirPath, err = ioutil.TempDir("", "")
	if err != nil { return err }
	defer os.RemoveAll(tempDirPath)
	tarFile, err = ioutil.TempFile(tempDirPath, "")
	if err != nil { return errors.New(fmt.Sprintf(
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
	
	if err != nil { return err }
	tarWriter.Close()
	
	// Send the request to the docker engine, with the tar file as the body content.
	var tarReader io.ReadCloser
	tarReader, err = os.Open(tarFile.Name())
	defer tarReader.Close()
	if err != nil { return err }
	var response *http.Response
	response, err = engine.SendBasicStreamPost("build", "application/tar", tarReader)
	if err != nil { return err }
	if response.StatusCode != 200 { return errors.New(response.Status) }
	return nil
}
