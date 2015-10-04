/*******************************************************************************
 * Perform independent end-to-end ("behavioral") tests on the SafeHarbor server.
 * It is assumed that the SafeHarbor server is running on localhost:6000.
 */

package main

import (
	"fmt"
	"net/http"
	"os"
	//"bufio"
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
	
	// Log in so that we can do stuff.
	_ = testContext.TryAuthenticate("testuser1", "password1")
	
	// Test ability to create a realm.
	var realmId string = testContext.TryCreateRealm()
	assertThat(realmId != "", "TryCreateRealm failed")
	
	// Test ability to create a user for the realm.
	var userId string = "jdoe"
	var userName string = "John Doe"
	var johnDoeUserObjId string = testContext.TryCreateUser(userId, userName, realmId)
	assertThat(johnDoeUserObjId != "", "TryCreateUser failed")
	
	// Login as the user that we just created.
	_ = testContext.TryAuthenticate(userId, "password1")
	
	// Test ability create a repo.
	var repoId string = testContext.TryCreateRepo(realmId, "John's Repo")
	assertThat(repoId != "", "TryCreateRepo failed")
		
	// Test ability create another repo.
	var repo2Id string = testContext.TryCreateRepo(realmId, "Susan's Repo")
	assertThat(repo2Id != "", "TryCreateRepo failed")
		
	// Test ability to upload a Dockerfile.
	var dockerfileId string = testContext.TryAddDockerfile(repoId, "Dockerfile")
	assertThat(dockerfileId != "", "TryAddDockerfile failed")
	
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
	var userObjId = testContext.TryGetRealmUser(realmId, userId)
	assertThat(userObjId == johnDoeUserObjId, "Looking up user by user id failed")
	
	//var msg string = testContext.TryAddRealmUser(....realmId, userObjId)
	
	var repoIds []string = testContext.TryGetRealmRepos(realmId)
	assertThat(len(repoIds) == 2, "Number of repo Ids returned was " +
		string(len(repoIds)) + ", expected 2")
	
	var realmIds []string = testContext.TryGetAllRealms()
	// Assumes that server is in debug mode, which creates test data.
	assertThat(len(realmIds) == 2, "Wrong number of realms found")
	
	
	
	// Test ability to clear the entire database and docker repository.
	testContext.TryClearAll()
	
	
	testContext.TryCreateGroup()
	
	
	testContext.TryGetGroupUsers()
	
	
	testContext.TryAddGroupUser()
	
	
	testContext.TryGetRealmGroups()
	
	
	testContext.TryReplaceDockerfile()
	
	
	testContext.TryDownloadImage()
	
	
	testContext.TrySetPermission()
	
	
	testContext.TryAddPermission()
	
	
	testContext.TryScanImage()
	
	
	testContext.TryGetMyGroups()
	
	
	testContext.TryGetMyRealms()
	
	
	testContext.TryGetMyRepos()
	
	
	testContext.TryDeleteUser()
	
	
	testContext.TryDeleteGroup()
	
	
	testContext.TryRemGroupUser()
	
	
	testContext.TryDeleteRealm()
	
	
	testContext.TryRemRealmUser()
	
	
	testContext.TryDeleteRepo()
	
	
	testContext.TryRemPermission()


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
	var responseMap map[string]interface{}
	responseMap  = parseResponseBodyToMap(resp.Body)
	var realmId string = responseMap["Id"].(string)
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
	
	var responseMap map[string]interface{}
	responseMap  = parseResponseBodyToMap(resp.Body)
	var retUserObjId string = responseMap["Id"].(string)
	var retUserId string = responseMap["UserId"].(string)
	var retUserName string = responseMap["UserName"].(string)
	var retRealmId string = responseMap["RealmId"].(string)
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
 * 
 */
func (testContext *TestContext) TryAuthenticate(userId string, pswd string) string {
	fmt.Println("TryAuthenticate")
	
	var resp *http.Response = testContext.sendPost(
		"authenticate",
		[]string{"UserUd", "Password"},
		[]string{userId, pswd})
	
	defer resp.Body.Close()

	verify200Response(resp)
	
	// Get the repo Id that is returned in the response body.
	var responseMap map[string]interface{}
	responseMap  = parseResponseBodyToMap(resp.Body)
	var retSessionId string = responseMap["UniqueSessionId"].(string)
	var retUserId string = responseMap["AuthenticatedUserid"].(string)
	printMap(responseMap)
	
	assertThat(retSessionId != "", "Session id is empty string")
	assertThat(retUserId == userId, "Returned user id '" + retUserId + "' does not match user id")
	return retSessionId
}

/*******************************************************************************
 * Verify that we can create a new repo. This requires that we first created
 * a realm that the repo can belong to.
 */
func (testContext *TestContext) TryCreateRepo(realmId string, name string) string {
	fmt.Println("TryCreateRepo")
	var resp *http.Response = testContext.sendPost(
		"createRepo",
		[]string{"RealmId", "Name"},
		[]string{realmId, name})
	
	defer resp.Body.Close()

	verify200Response(resp)
	
	// Get the repo Id that is returned in the response body.
	var responseMap map[string]interface{}
	responseMap  = parseResponseBodyToMap(resp.Body)
	var repoId string = responseMap["Id"].(string)
	var repoName string = responseMap["Name"].(string)
	printMap(responseMap)
	assertThat(repoId != "", "Repo Id not found in response body")
	assertThat(repoName != "", "Repo Name not found in response body")
	
	return repoId
}

/*******************************************************************************
 * Verify that we can upload a dockerfile. This requries that we first created
 * a repo to uplaod it into.
 */
func (testContext *TestContext) TryAddDockerfile(repoId string, dockerfilePath string) string {
	
	fmt.Println("TryAddDockerfile")
	fmt.Println("\t", dockerfilePath)
	var resp *http.Response = testContext.sendFilePost(
		"addDockerfile",
		[]string{"RepoId"},
		[]string{repoId},
		dockerfilePath)
	
	defer resp.Body.Close()

	verify200Response(resp)
	
	// Get the DockerfileDesc that is returned.
	var responseMap map[string]interface{}
	responseMap  = parseResponseBodyToMap(resp.Body)
	var dockerfileId string = responseMap["Id"].(string)
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
	
	var responseMaps []map[string]interface{} = parseResponseBodyToMaps(resp.Body)
	var result []string = make([]string, 0)
	for _, responseMap := range responseMaps {
		var dockerfileId string = responseMap["Id"].(string)
		var repoId string = responseMap["RepoId"].(string)
		var dockerfileName string = responseMap["Name"].(string)

		printMap(responseMap)
		assertThat(dockerfileId != "", "Dockerfile Id not found in response body")
		assertThat(repoId != "", "Repo Id not found in response body")
		assertThat(dockerfileName != "", "Dockerfile Name not found in response body")
		fmt.Println()

		result = append(result, dockerfileName)
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
	var responseMap map[string]interface{}
	responseMap  = parseResponseBodyToMap(resp.Body)
	var objId string = responseMap["ObjId"].(string)
	var dockerImageId string = responseMap["DockerImageId"].(string)
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
	
	var responseMaps []map[string]interface{} = parseResponseBodyToMaps(resp.Body)
	var result []string = make([]string, 0)
	for _, responseMap := range responseMaps {
		var objId string = responseMap["ObjId"].(string)
		var dockerImageId string = responseMap["DockerImageId"].(string)

		printMap(responseMap)
		assertThat(objId != "", "ObjId not found in response body")
		assertThat(dockerImageId != "", "DockerImageId not found in response body")
		fmt.Println()

		result = append(result, dockerImageId)
	}
	
	return result
}

/*******************************************************************************
 * Return the object Id of the specified user.
 */
func (testContext *TestContext) TryGetRealmUser(realmId, userId string) string {
	fmt.Println("TryGetUserById")
	
	var resp *http.Response = testContext.sendPost(
		"getRealmUser",
		[]string{"RealmId", "UserId"},
		[]string{realmId, userId})
	
	defer resp.Body.Close()

	verify200Response(resp)
	
	var responseMap map[string]interface{}
	responseMap  = parseResponseBodyToMap(resp.Body)
	var retUserObjId string = responseMap["Id"].(string)
	var retUserId string = responseMap["UserId"].(string)
	var retUserName string = responseMap["UserName"].(string)
	var retRealmId string = responseMap["RealmId"].(string)
	printMap(responseMap)
	
	assertThat(retUserObjId != "", "User obj Id not returned")
	assertThat(retUserId == userId, "Returned user id, " + retUserId +
		" does not match the original user id")
	assertThat(retUserName != "", "Returned user name is blank")
	assertThat(retRealmId == realmId, "Returned realm Id, " + retRealmId +
		" does not match the original realm Id")
	
	return retUserObjId
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryCreateGroup() {
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryGetGroupUsers() {
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryAddGroupUser() {
}

/*******************************************************************************
 * Returns result.
 */
func (testContext *TestContext) TryAddRealmUser(realmId string, userObjId string) string {
	fmt.Println("TryAddRealmUser")
	
	var resp *http.Response = testContext.sendPost(
		"addRealmUser",
		[]string{"RealmId", "UserObjId"},
		[]string{realmId, userObjId})
	
	defer resp.Body.Close()

	verify200Response(resp)
	
	var responseMap map[string]interface{}
	responseMap  = parseResponseBodyToMap(resp.Body)
	var retStatus string = responseMap["Status"].(string)
	var retMsg string = responseMap["Message"].(string)
	printMap(responseMap)
	assertThat(retStatus != "", "Empty return status")
	assertThat(retMsg != "", "Empty return message")
	return retMsg
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryGetRealmGroups() {
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryGetRealmRepos(realmId string) []string {
	fmt.Println("TryGetRealmRepos")
	
	var resp *http.Response = testContext.sendPost(
		"getRealmRepos",
		[]string{"RealmId"},
		[]string{realmId})
	
	defer resp.Body.Close()

	verify200Response(resp)
	
	var responseMaps []map[string]interface{} = parseResponseBodyToMaps(resp.Body)
	var result []string = make([]string, 0)
	for _, responseMap := range responseMaps {
		printMap(responseMap)
		var retRepoId string = responseMap["Id"].(string)
		var retRealmId string = responseMap["RealmId"].(string)
		var retName string = responseMap["Name"].(string)
	
		assertThat(retRepoId != "", "No repo Id returned")
		assertThat(retRealmId == realmId, "returned realm Id is nil")
		assertThat(retName != "", "Empty returned Name")
		
		result = append(result, retRepoId)
	}
	return result
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryGetAllRealms() []string {
	fmt.Println("TryGetAllRealms")
	
	var resp *http.Response = testContext.sendPost(
		"getRealmRepos",
		[]string{},
		[]string{})
	
	defer resp.Body.Close()

	verify200Response(resp)
	
	var responseMaps []map[string]interface{} = parseResponseBodyToMaps(resp.Body)
	var result []string = make([]string, 0)
	for _, responseMap := range responseMaps {
		printMap(responseMap)
		var retRealmId string = responseMap["Id"].(string)
		var retName string = responseMap["Name"].(string)
	
		assertThat(retRealmId != "", "Returned realm Id is empty string")
		assertThat(retName != "", "Empty returned Name")
		
		result = append(result, retRealmId)
	}
	return result
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryReplaceDockerfile() {
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryDownloadImage() {
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TrySetPermission() {
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryAddPermission() {
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryScanImage() {
}




/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryGetMyGroups() {
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryGetMyRealms() {
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryGetMyRepos() {
}




/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryDeleteUser() {
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryDeleteGroup() {
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryRemGroupUser() {
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryDeleteRealm() {
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryRemRealmUser() {
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryDeleteRepo() {
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryRemPermission() {
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryClearAll() {
	fmt.Println("TryClearAll")
	
	var resp *http.Response = testContext.sendGet(
		"clearAll",
		[]string{},
		[]string{},
		)
	
	verify200Response(resp)
}
