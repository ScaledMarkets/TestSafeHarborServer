/*******************************************************************************
 * Perform independent end-to-end tests on the SafeHarbor server.
 * It is assumed that the SafeHarbor server is running on localhost:6000.
 */

package main

import (
	"fmt"
	"net/http"
	"os"
	"bufio"
)

type TestContext struct {
	hostname string
	port string
}

func main() {
	
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <hostname> <port>\n", os.Args[0])
		os.Exit(2)
	}
	
	var testContext *TestContext = &TestContext{
		hostname: os.Args[1],
		port: os.Args[2],
	}
	
	fmt.Println("Note: Ensure that the docker daemon is running on the server,",
		"and that python 2 is installed on the server. To start the docker daemon",
		"run 'sudo service docker start'")
	fmt.Println()
	
	// Test ability to create a realm.
	var realmId string = testContext.TryCreateRealm()
	assertThat(realmId != "", "TryCreateRealm failed")
	
	// Test ability to create a user for the realm.
	var userId string = "jdoe"
	var userName string = "John Doe"
	var johnDoeUserObjId string = testContext.TryCreateUser(userId, userName, realmId)
	assertThat(johnDoeUserObjId != "", "TryCreateUser failed")
	
	// Test ability create a repo.
	var repoId string = testContext.TryCreateRepo(realmId)
	assertThat(repoId != "", "TryCreateRepo failed")
		
	// Test ability to upload a Dockerfile.
	var dockerfileId string = testContext.TryUploadDockerfile(repoId, "Dockerfile")
	assertThat(dockerfileId != "", "TryUploadDockerfile failed")
	
	// Test ability to list the Dockerfiles in a repo.
	var dockerfileNames []string = testContext.TryGetDockerfiles(repoId)
	assertThat(len(dockerfileNames) == 1, "Wrong number of dockerfiles")
	
	// Test ability to build image from a dockerfile.
	var dockerImageObjId string
	var imageId string
	dockerImageObjId, imageId = testContext.TryExecDockerfile(repoId, dockerfileId, "myimage")
	assertThat(dockerImageObjId != "", "TryExecDockerfile failed - obj id is nil")
	assertThat(imageId != "", "TryExecDockerfile failed - docker image id is nil")
	
	// Test ability to list the images in a repo.
	var imageNames []string = testContext.TryGetImages(repoId)
	assertThat(len(imageNames) == 1, "Wrong number of images")
	
	// Test ability to retrieve user by user id from realm.
	var userObjId = testContext.TryGetUserByUserId(realmId, userId)
	assertThat(userObjId == johnDoeUserObjId, "Looking up user by user id failed")
	
	// Test ability to receive progress while a Dockerfile is processed.
	
	// Test ability to make a private image available to the SafeHarbor closed community.
	
	// Test ability to make a private image available to another user.
	
}


/*******************************************************************************
 * Verify that we can create a new realm.
 */
func (testContext *TestContext) TryCreateRealm() string {
	
	fmt.Println("TryCreateRealm")
	var resp *http.Response = testContext.sendPost(
		"createRealm",
		[]string{"Name"},
		[]string{"My Realm"})
	
	defer resp.Body.Close()
	
	verify200Response(resp)
	
	// Get the realm Id that is returned in the response body.
	var responseMap map[string]string
	responseMap, _ = parseResponseBody(resp.Body)
	var realmId string = responseMap["Id"]
	printMap(responseMap)
	assertThat(realmId != "", "Realm Id not found in response body")
	
	return realmId
}

/*******************************************************************************
 * Return the object Id of the new user.
 */
func (testContext *TestContext) TryCreateUser(userId string, userName string,
	realmId string) string {
	fmt.Println("TryCreateUser")
	
	var resp *http.Response = testContext.sendPost(
		"createUser",
		[]string{"UserId", "UserName", "RealmId"},
		[]string{userId, userName, realmId})
	
	defer resp.Body.Close()

	verify200Response(resp)
	
	var responseMap map[string]string
	responseMap, _ = parseResponseBody(resp.Body)
	var retUserObjId string = responseMap["Id"]
	var retUserId string = responseMap["UserId"]
	var retUserName string = responseMap["UserName"]
	var retRealmId string = responseMap["RealmId"]
	printMap(responseMap)
	
	assertThat(retUserObjId != "", "User obj Id not returned")
	assertThat(retUserId == userId, "Returned user id, " + retUserId +
		" does not match the original user id")
	assertThat(retUserName == userName, "Returned user name, " + retUserName +
		" does not match the original user name")
	assertThat(retRealmId == realmId, "Returned realm Id, " + retRealmId +
		" does not match the original realm Id")
	
	return retUserObjId
}

/*******************************************************************************
 * Verify that we can create a new repo. This requires that we first created
 * a realm that the repo can belong to.
 */
func (testContext *TestContext) TryCreateRepo(realmId string) string {
	fmt.Println("TryCreateRepo")
	var resp *http.Response = testContext.sendPost(
		"createRepo",
		[]string{"RealmId", "Name"},
		[]string{realmId, "John's Repo"})
	
	defer resp.Body.Close()

	verify200Response(resp)
	
	// Get the repo Id that is returned in the response body.
	var responseMap map[string]string
	responseMap, _ = parseResponseBody(resp.Body)
	var repoId string = responseMap["Id"]
	var repoName string = responseMap["Name"]
	printMap(responseMap)
	assertThat(repoId != "", "Repo Id not found in response body")
	assertThat(repoName != "", "Repo Name not found in response body")
	
	return repoId
}

/*******************************************************************************
 * Verify that we can upload a dockerfile. This requries that we first created
 * a repo to uplaod it into.
 */
func (testContext *TestContext) TryUploadDockerfile(repoId string, dockerfilePath string) string {
	
	fmt.Println("TryUploadDockerfile")
	fmt.Println("\t", dockerfilePath)
	var resp *http.Response = testContext.sendFilePost(
		"addDockerfile",
		[]string{"RepoId"},
		[]string{repoId},
		dockerfilePath)
	
	defer resp.Body.Close()

	verify200Response(resp)
	
	// Get the DockerfileDesc that is returned.
	var responseMap map[string]string
	responseMap, _ = parseResponseBody(resp.Body)
	var dockerfileId string = responseMap["Id"]
	//var dockerfileName string = responseMap["Name"]
	printMap(responseMap)
	//assertThat(dockerfileId != "", "Dockerfile Id not found in response body")
	//assertThat(dockerfileName != "", "Dockerfile Name not found in response body")
	
	return dockerfileId
}

/*******************************************************************************
 * Verify that we can obtain the names of the dockerfiles owned by the specified
 * repo. The result is an array of dockerfile names.
 */
func (testContext *TestContext) TryGetDockerfiles(repoId string) []string {
	fmt.Println("TryGetDockerfiles")
	
	var resp *http.Response = testContext.sendPost(
		"getDockerfiles",
		[]string{"RepoId"},
		[]string{repoId})
	
	defer resp.Body.Close()

	verify200Response(resp)
	
	var responseMap map[string]string
	var scanner *bufio.Scanner
	var result []string = make([]string, 1)
	responseMap, scanner = parseResponseBody(resp.Body)
	for responseMap != nil {
		var dockerfileId string = responseMap["Id"]
		var repoId string = responseMap["RepoId"]
		var dockerfileName string = responseMap["Name"]

		printMap(responseMap)
		assertThat(dockerfileId != "", "Dockerfile Id not found in response body")
		assertThat(repoId != "", "Repo Id not found in response body")
		assertThat(dockerfileName != "", "Dockerfile Name not found in response body")
		fmt.Println()

		result = append(result, dockerfileName)
		responseMap = parseNextBodyPart(scanner)
	}
		
	return result
}

/*******************************************************************************
 * Verify that we can build an image, from a dockerfile that has already been
 * uploaded into a repo and for which we have the SafeHarborServer image id.
 * The result is the object id and docker id of the image that was built.
 */
func (testContext *TestContext) TryExecDockerfile(repoId string, dockerfileId string,
	imageName string) (string, string) {
	fmt.Println("TryExecDockerfile")
	
	var resp *http.Response = testContext.sendPost(
		"execDockerfile",
		[]string{"RepoId", "DockerfileId", "ImageName"},
		[]string{repoId, dockerfileId, imageName})
	
	defer resp.Body.Close()

	verify200Response(resp)
	
	// Get the repo Id that is returned in the response body.
	var responseMap map[string]string
	responseMap, _ = parseResponseBody(resp.Body)
	var objId string = responseMap["ObjId"]
	var dockerImageId string = responseMap["DockerImageId"]
	printMap(responseMap)
	return objId, dockerImageId
}

/*******************************************************************************
 * Result is an array of the names of the images owned by the specified repo.
 */
func (testContext *TestContext) TryGetImages(repoId string) []string {
	fmt.Println("TryGetImages")
	
	var resp *http.Response = testContext.sendPost(
		"getImages",
		[]string{"RepoId"},
		[]string{repoId})
	
	defer resp.Body.Close()

	verify200Response(resp)
	
	var responseMap map[string]string
	var scanner *bufio.Scanner
	var result []string = make([]string, 1)
	responseMap, scanner = parseResponseBody(resp.Body)
	for responseMap != nil {
		var objId string = responseMap["ObjId"]
		var dockerImageId string = responseMap["DockerImageId"]

		printMap(responseMap)
		assertThat(objId != "", "ObjId not found in response body")
		assertThat(dockerImageId != "", "DockerImageId not found in response body")
		fmt.Println()

		result = append(result, dockerImageId)
		responseMap = parseNextBodyPart(scanner)
	}
	
	return result
}

/*******************************************************************************
 * Return the object Id of the specified user.
 */
func (testContext *TestContext) TryGetUserByUserId(realmId, userId string) string {
	fmt.Println("TryGetUserById")
	
	var resp *http.Response = testContext.sendPost(
		"getRealmUser",
		[]string{"RealmId", "UserId"},
		[]string{realmId, userId})
	
	defer resp.Body.Close()

	verify200Response(resp)
	
	var responseMap map[string]string
	responseMap, _ = parseResponseBody(resp.Body)
	var retUserObjId string = responseMap["Id"]
	var retUserId string = responseMap["UserId"]
	var retUserName string = responseMap["UserName"]
	var retRealmId string = responseMap["RealmId"]
	printMap(responseMap)
	
	assertThat(retUserObjId != "", "User obj Id not returned")
	assertThat(retUserId == userId, "Returned user id, " + retUserId +
		" does not match the original user id")
	assertThat(retUserName != "", "Returned user name is blank")
	assertThat(retRealmId == realmId, "Returned realm Id, " + retRealmId +
		" does not match the original realm Id")
	
	return retUserObjId
}
