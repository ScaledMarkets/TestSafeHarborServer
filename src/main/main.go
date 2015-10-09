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
	sessionId string
	testName string
}

func main() {
	
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <hostname> <port>\n", os.Args[0])
		os.Exit(2)
	}
	
	var testContext *TestContext = &TestContext{
		hostname: os.Args[1],
		port: os.Args[2],
		sessionId: "",
	}
	
	fmt.Println("Note: Ensure that the docker daemon is running on the server,",
		"and that python 2 is installed on the server. To start the docker daemon",
		"run 'sudo service docker start'")
	fmt.Println()
	
	// Log in so that we can do stuff.
	var sessionId string = testContext.TryAuthenticate("testuser1", "password1")
	testContext.sessionId = sessionId
	fmt.Println("sessionId =", sessionId)
	
	// Test ability to create a realm.
	var realmId string = testContext.TryCreateRealm("MyRealm")
	testContext.assertThat(realmId != "", "TryCreateRealm failed")
	
	// Test ability to create a user for the realm.
	var userId string = "jdoe"
	var userName string = "John Doe"
	var johnDoeUserObjId string = testContext.TryCreateUser(userId, userName, realmId)
	testContext.assertThat(johnDoeUserObjId != "", "TryCreateUser failed")
	
	// Login as the user that we just created.
	sessionId = testContext.TryAuthenticate(userId, "password1")
	testContext.sessionId = sessionId
	
	// Test ability to create a realm.
	var jrealm1Id string = testContext.TryCreateRealm("John's First Realm")
	testContext.assertThat(jrealm1Id != "", "TryCreateRealm failed")
	
	// Test ability to create a realm.
	var jrealm2Id string = testContext.TryCreateRealm("John's Second Realm")
	testContext.assertThat(jrealm2Id != "", "TryCreateRealm failed")
	
	// Test ability create a repo.
	var repoId string = testContext.TryCreateRepo(realmId, "John's Repo")
	testContext.assertThat(repoId != "", "TryCreateRepo failed")
		
	// Test ability create another repo.
	var repo2Id string = testContext.TryCreateRepo(realmId, "Susan's Repo")
	testContext.assertThat(repo2Id != "", "TryCreateRepo failed")
		
	// Test ability to upload a Dockerfile.
	var dockerfileId string = testContext.TryAddDockerfile(repoId, "Dockerfile")
	testContext.assertThat(dockerfileId != "", "TryAddDockerfile failed")
	
	// Test ability to list the Dockerfiles in a repo.
	var dockerfileNames []string = testContext.TryGetDockerfiles(repoId)
	testContext.assertThat(len(dockerfileNames) == 1, "Wrong number of dockerfiles")
	
	// Test ability to build image from a dockerfile.
	var dockerImageObjId string
	var imageId string
	dockerImageObjId, imageId = testContext.TryExecDockerfile(repoId, dockerfileId, "myimage")
	testContext.assertThat(dockerImageObjId != "", "TryExecDockerfile failed - obj id is nil")
	testContext.assertThat(imageId != "", "TryExecDockerfile failed - docker image id is nil")
	
	// Test ability to list the images in a repo.
	var imageNames []string = testContext.TryGetImages(repoId)
	testContext.assertThat(len(imageNames) == 1, "Wrong number of images")
	
	// Test ability to retrieve user by user id from realm.
	var userObjId = testContext.TryGetRealmUser(realmId, userId)
	testContext.assertThat(userObjId == johnDoeUserObjId, "Looking up user by user id failed")
	
	//var msg string = testContext.TryAddRealmUser(....realmId, userObjId)
	
	var repoIds []string = testContext.TryGetRealmRepos(realmId)
	testContext.assertThat(len(repoIds) == 2, "Number of repo Ids returned was " +
		string(len(repoIds)) + ", expected 2")
	
	var realmIds []string = testContext.TryGetAllRealms()
	// Assumes that server is in debug mode, which creates test data.
	testContext.assertThat(len(realmIds) == 4, "Wrong number of realms found")
	
	
	
	testContext.TryCreateGroup()
	
	
	testContext.TryGetGroupUsers()
	
	
	testContext.TryAddGroupUser()
	
	
	testContext.TryGetRealmGroups()
	
	
	testContext.TryReplaceDockerfile()
	
	
	testContext.TryDownloadImage()
	
	
	testContext.TrySetPermission()
	
	
	testContext.TryAddPermission()
	
	
	testContext.TryScanImage()
	
	
	userObjId = testContext.TryGetMyInfo()
	testContext.assertThat(userObjId == johnDoeUserObjId,
		"Returned user obj id was " + userObjId)
	
	
	testContext.TryGetMyGroups()
	
	
	var myRealms []string = testContext.TryGetMyRealms()
	testContext.assertThat(len(myRealms) == 2, fmt.Sprintf(
		"Only returned %d realms", len(myRealms)))
	
	
	var myRepos []string = testContext.TryGetMyRepos()
	testContext.assertThat(len(myRepos) == 2, fmt.Sprintf(
		"Only returned %d repos", len(myRepos)))
	
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
	
	// Test ability to clear the entire database and docker repository.
	testContext.TryClearAll()
	
	
	fmt.Println()
	fmt.Println(fmt.Sprintf("%d tests failed out of %d", noOfTestsThatFailed, noOfTests))
}


/*******************************************************************************
 * Verify that we can create a new realm.
 */
func (testContext *TestContext) TryCreateRealm(realmName string) string {
	
	testContext.StartTest("TryCreateRealm")
	var resp *http.Response = testContext.sendPost(testContext.sessionId,
		"createRealm",
		[]string{"Name"},
		[]string{realmName})
	
	defer resp.Body.Close()
	
	testContext.verify200Response(resp)
	
	// Get the realm Id that is returned in the response body.
	var responseMap map[string]interface{}
	responseMap  = parseResponseBodyToMap(resp.Body)
	var realmId string = responseMap["Id"].(string)
	printMap(responseMap)
	testContext.assertThat(realmId != "", "Realm Id not found in response body")
	
	return realmId
}

/*******************************************************************************
 * Return the object Id of the new user.
 */
func (testContext *TestContext) TryCreateUser(userId string, userName string,
	realmId string) string {
	testContext.StartTest("TryCreateUser")
	
	var resp *http.Response = testContext.sendPost(testContext.sessionId,
		"createUser",
		[]string{"UserId", "UserName", "RealmId"},
		[]string{userId, userName, realmId})
	
	defer resp.Body.Close()

	testContext.verify200Response(resp)
	
	var responseMap map[string]interface{}
	responseMap  = parseResponseBodyToMap(resp.Body)
	var retUserObjId string = responseMap["Id"].(string)
	var retUserId string = responseMap["UserId"].(string)
	var retUserName string = responseMap["UserName"].(string)
	var retRealmId string = responseMap["RealmId"].(string)
	printMap(responseMap)
	
	testContext.assertThat(retUserObjId != "", "User obj Id not returned")
	testContext.assertThat(retUserId == userId, "Returned user id, " + retUserId +
		" does not match the original user id")
	testContext.assertThat(retUserName == userName, "Returned user name, " + retUserName +
		" does not match the original user name")
	testContext.assertThat(retRealmId == realmId, "Returned realm Id, " + retRealmId +
		" does not match the original realm Id")
	
	return retUserObjId
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryAuthenticate(userId string, pswd string) string {
	testContext.StartTest("TryAuthenticate")
	
	var resp *http.Response = testContext.sendPost(testContext.sessionId,
		"authenticate",
		[]string{"UserId", "Password"},
		[]string{userId, pswd})
	
	defer resp.Body.Close()

	testContext.verify200Response(resp)
	
	// Get the repo Id that is returned in the response body.
	var responseMap map[string]interface{}
	responseMap  = parseResponseBodyToMap(resp.Body)
	var retSessionId string = responseMap["UniqueSessionId"].(string)
	var retUserId string = responseMap["AuthenticatedUserid"].(string)
	printMap(responseMap)
	
	testContext.assertThat(retSessionId != "", "Session id is empty string")
	testContext.assertThat(retUserId == userId, "Returned user id '" + retUserId + "' does not match user id")
	return retSessionId
}

/*******************************************************************************
 * Verify that we can create a new repo. This requires that we first created
 * a realm that the repo can belong to.
 */
func (testContext *TestContext) TryCreateRepo(realmId string, name string) string {
	testContext.StartTest("TryCreateRepo")
	var resp *http.Response = testContext.sendPost(testContext.sessionId,
		"createRepo",
		[]string{"RealmId", "Name"},
		[]string{realmId, name})
	
	defer resp.Body.Close()

	testContext.verify200Response(resp)
	
	// Get the repo Id that is returned in the response body.
	var responseMap map[string]interface{}
	responseMap  = parseResponseBodyToMap(resp.Body)
	var repoId string = responseMap["Id"].(string)
	var repoName string = responseMap["Name"].(string)
	printMap(responseMap)
	testContext.assertThat(repoId != "", "Repo Id not found in response body")
	testContext.assertThat(repoName != "", "Repo Name not found in response body")
	
	return repoId
}

/*******************************************************************************
 * Verify that we can upload a dockerfile. This requries that we first created
 * a repo to uplaod it into.
 */
func (testContext *TestContext) TryAddDockerfile(repoId string, dockerfilePath string) string {
	
	testContext.StartTest("TryAddDockerfile")
	fmt.Println("\t", dockerfilePath)
	var resp *http.Response = testContext.sendFilePost(testContext.sessionId,
		"addDockerfile",
		[]string{"RepoId"},
		[]string{repoId},
		dockerfilePath)
	
	defer resp.Body.Close()

	testContext.verify200Response(resp)
	
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
	testContext.StartTest("TryGetDockerfiles")
	
	var resp *http.Response = testContext.sendPost(testContext.sessionId,
		"getDockerfiles",
		[]string{"RepoId"},
		[]string{repoId})
	
	defer resp.Body.Close()

	testContext.verify200Response(resp)
	
	var responseMaps []map[string]interface{} = parseResponseBodyToMaps(resp.Body)
	var result []string = make([]string, 0)
	for _, responseMap := range responseMaps {
		var dockerfileId string = responseMap["Id"].(string)
		var repoId string = responseMap["RepoId"].(string)
		var dockerfileName string = responseMap["Name"].(string)

		printMap(responseMap)
		testContext.assertThat(dockerfileId != "", "Dockerfile Id not found in response body")
		testContext.assertThat(repoId != "", "Repo Id not found in response body")
		testContext.assertThat(dockerfileName != "", "Dockerfile Name not found in response body")
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
	testContext.StartTest("TryExecDockerfile")
	
	var resp *http.Response = testContext.sendPost(testContext.sessionId,
		"execDockerfile",
		[]string{"RepoId", "DockerfileId", "ImageName"},
		[]string{repoId, dockerfileId, imageName})
	
	defer resp.Body.Close()

	testContext.verify200Response(resp)
	
	// Get the repo Id that is returned in the response body.
	var responseMap map[string]interface{}
	responseMap  = parseResponseBodyToMap(resp.Body)
	var objId string = responseMap["ObjId"].(string)
	var dockerImageTag string = responseMap["DockerImageTag"].(string)
	printMap(responseMap)
	return objId, dockerImageTag
}

/*******************************************************************************
 * Result is an array of the names of the images owned by the specified repo.
 */
func (testContext *TestContext) TryGetImages(repoId string) []string {
	testContext.StartTest("TryGetImages")
	
	var resp *http.Response = testContext.sendPost(testContext.sessionId,
		"getImages",
		[]string{"RepoId"},
		[]string{repoId})
	
	defer resp.Body.Close()

	testContext.verify200Response(resp)
	
	var responseMaps []map[string]interface{} = parseResponseBodyToMaps(resp.Body)
	var result []string = make([]string, 0)
	for _, responseMap := range responseMaps {
		var objId string = responseMap["ObjId"].(string)
		var dockerImageTag string = responseMap["DockerImageTag"].(string)

		printMap(responseMap)
		testContext.assertThat(objId != "", "ObjId not found in response body")
		testContext.assertThat(dockerImageTag != "", "DockerImageTag not found in response body")
		fmt.Println()

		result = append(result, dockerImageTag)
	}
	
	return result
}

/*******************************************************************************
 * Return the object Id of the specified user.
 */
func (testContext *TestContext) TryGetRealmUser(realmId, userId string) string {
	testContext.StartTest("TryGetRealmUser")
	
	var resp *http.Response = testContext.sendPost(testContext.sessionId,
		"getRealmUser",
		[]string{"RealmId", "UserId"},
		[]string{realmId, userId})
	
	defer resp.Body.Close()

	testContext.verify200Response(resp)
	
	var responseMap map[string]interface{}
	responseMap  = parseResponseBodyToMap(resp.Body)
	var retUserObjId string = responseMap["Id"].(string)
	var retUserId string = responseMap["UserId"].(string)
	var retUserName string = responseMap["UserName"].(string)
	var retRealmId string = responseMap["RealmId"].(string)
	printMap(responseMap)
	
	testContext.assertThat(retUserObjId != "", "User obj Id not returned")
	testContext.assertThat(retUserId == userId, "Returned user id, " + retUserId +
		" does not match the original user id")
	testContext.assertThat(retUserName != "", "Returned user name is blank")
	testContext.assertThat(retRealmId == realmId, "Returned realm Id, " + retRealmId +
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
	testContext.StartTest("TryAddRealmUser")
	
	var resp *http.Response = testContext.sendPost(testContext.sessionId,
		"addRealmUser",
		[]string{"RealmId", "UserObjId"},
		[]string{realmId, userObjId})
	
	defer resp.Body.Close()

	testContext.verify200Response(resp)
	
	var responseMap map[string]interface{}
	responseMap  = parseResponseBodyToMap(resp.Body)
	var retStatus string = responseMap["Status"].(string)
	var retMsg string = responseMap["Message"].(string)
	printMap(responseMap)
	testContext.assertThat(retStatus != "", "Empty return status")
	testContext.assertThat(retMsg != "", "Empty return message")
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
	testContext.StartTest("TryGetRealmRepos")
	
	var resp *http.Response = testContext.sendPost(testContext.sessionId,
		"getRealmRepos",
		[]string{"RealmId"},
		[]string{realmId})
	
	defer resp.Body.Close()

	testContext.verify200Response(resp)
	
	var responseMaps []map[string]interface{} = parseResponseBodyToMaps(resp.Body)
	var result []string = make([]string, 0)
	for _, responseMap := range responseMaps {
		printMap(responseMap)
		var retRepoId string = responseMap["Id"].(string)
		var retRealmId string = responseMap["RealmId"].(string)
		var retName string = responseMap["Name"].(string)
	
		testContext.assertThat(retRepoId != "", "No repo Id returned")
		testContext.assertThat(retRealmId == realmId, "returned realm Id is nil")
		testContext.assertThat(retName != "", "Empty returned Name")
		
		result = append(result, retRepoId)
	}
	return result
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryGetAllRealms() []string {
	testContext.StartTest("TryGetAllRealms")
	
	var resp *http.Response = testContext.sendPost(testContext.sessionId,
		"getAllRealms",
		[]string{},
		[]string{})
	
	defer resp.Body.Close()

	testContext.verify200Response(resp)
	
	var responseMaps []map[string]interface{} = parseResponseBodyToMaps(resp.Body)
	var result []string = make([]string, 0)
	for _, responseMap := range responseMaps {
		printMap(responseMap)
		var retRealmId string = responseMap["Id"].(string)
		var retName string = responseMap["Name"].(string)
	
		testContext.assertThat(retRealmId != "", "Returned realm Id is empty string")
		testContext.assertThat(retName != "", "Empty returned Name")
		
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
 * Return the object Id of the current authenticated user.
 */
func (testContext *TestContext) TryGetMyInfo() string {
	testContext.StartTest("TryGetMyInfo")
	
	var resp *http.Response = testContext.sendPost(testContext.sessionId,
		"getMyInfo",
		[]string{},
		[]string{})
	
	defer resp.Body.Close()

	testContext.verify200Response(resp)
	
	var responseMap map[string]interface{} = parseResponseBodyToMap(resp.Body)
	printMap(responseMap)
	var retId string = responseMap["Id"].(string)
	var retUserId string = responseMap["UserId"].(string)
	var retUserName string = responseMap["UserName"].(string)
	var retRealmId string = responseMap["RealmId"].(string)

	testContext.assertThat(retId != "", "Returned Id is empty string")
	testContext.assertThat(retUserId != "", "Returned UserId is empty string")
	testContext.assertThat(retUserName != "", "Returned UserName is empty string")
	testContext.assertThat(retRealmId != "", "Returned RealmId is empty string")
	
	return retId
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryGetMyGroups() {
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryGetMyRealms() []string {
	testContext.StartTest("TryGetMyRealms")
	
	var resp *http.Response = testContext.sendPost(testContext.sessionId,
		"getMyRealms",
		[]string{},
		[]string{})
	
	defer resp.Body.Close()

	testContext.verify200Response(resp)
	
	var responseMaps []map[string]interface{} = parseResponseBodyToMaps(resp.Body)
	var result []string = make([]string, 0)
	for _, responseMap := range responseMaps {
		printMap(responseMap)
		var retId string = responseMap["Id"].(string)
		var retName string = responseMap["Name"].(string)
	
		testContext.assertThat(retId != "", "Returned Id is empty string")
		testContext.assertThat(retName != "", "Empty returned Name")
		
		result = append(result, retId)
	}
	return result
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryGetMyRepos() []string {
	testContext.StartTest("TryGetMyRepos")
	
	var resp *http.Response = testContext.sendPost(testContext.sessionId,
		"getMyRepos",
		[]string{},
		[]string{})
	
	defer resp.Body.Close()

	testContext.verify200Response(resp)
	
	var responseMaps []map[string]interface{} = parseResponseBodyToMaps(resp.Body)
	var result []string = make([]string, 0)
	for _, responseMap := range responseMaps {
		printMap(responseMap)
		var retId string = responseMap["Id"].(string)
		var retRealmId string = responseMap["RealmId"].(string)
		var retName string = responseMap["Name"].(string)
	
		testContext.assertThat(retId != "", "Returned Id is empty string")
		testContext.assertThat(retRealmId != "", "Returned realm Id is empty string")
		testContext.assertThat(retName != "", "Empty returned Name")
		
		result = append(result, retId)
	}
	return result
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
	testContext.StartTest("TryClearAll")
	
	var resp *http.Response = testContext.sendGet("",
		"clearAll",
		[]string{},
		[]string{},
		)
	
	testContext.verify200Response(resp)
}
