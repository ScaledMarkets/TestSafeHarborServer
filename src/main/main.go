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
	httpClient *http.Client
	hostname string
	port string
	sessionId string
	testName string
	stopOnFirstError bool
}

func main() {
	
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <hostname> <port> [stop]\n", os.Args[0])
		os.Exit(2)
	}
	
	var stopOnFirstError bool = false
	if len(os.Args) > 3 { if os.Args[3] == "stop" { stopOnFirstError = true } }
	
	var testContext *TestContext = &TestContext{
		hostname: os.Args[1],
		port: os.Args[2],
		sessionId: "",
		stopOnFirstError: stopOnFirstError,
	}
	
	fmt.Println("Note: Ensure that the docker daemon is running on the server,",
		"and that python 2 is installed on the server. To start the docker daemon",
		"run 'sudo service docker start'")
	fmt.Println()
	
	// Verify that we can create a realm without being logged in first.
	var realm4Id string
	var user4Id string
	var user4AdminRealms []interface{}
	realm4Id, user4Id, user4AdminRealms = testContext.TryCreateRealmAnon("Realm4", "Realm 4 Org",
		"realm4admin", "Realm 4 Admin Full Name", "realm4admin@gmail.com", "realm4adminpswd")
	testContext.assertThat(realm4Id != "", "Realm Id is empty")
	testContext.assertThat(user4Id != "", "User Id is empty")
	testContext.assertThat(len(user4AdminRealms) == 1, "Wrong number of admin realms")
	
	// Verify that we can log in as the admin user that we just created.
	var sessionId string = testContext.TryAuthenticate("realm4admin", "realm4adminpswd")
	testContext.sessionId = sessionId
	fmt.Println("sessionId =", sessionId)
	
	// Log in so that we can do stuff.
	sessionId = testContext.TryAuthenticate("testuser1", "password1")
	testContext.sessionId = sessionId
	fmt.Println("sessionId =", sessionId)
	
	// Test ability to create a realm.
	var realmId string = testContext.TryCreateRealm("MyRealm", "A Big Company", "bigshotadmin")
	testContext.assertThat(realmId != "", "TryCreateRealm failed")
	
	// Test ability to create a user for the realm.
	var userId string = "jdoe"
	var userName string = "John Doe"
	var johnDoeUserObjId string
	var johnDoeAdminRealms []interface{}
	johnDoeUserObjId, johnDoeAdminRealms = testContext.TryCreateUser(userId, userName,
		"johnd@gmail.com", "weakpswd", realmId)
	testContext.assertThat(johnDoeUserObjId != "", "TryCreateUser failed")
	testContext.assertThat(len(johnDoeAdminRealms) == 0, "Wrong number of admin realms")
	
	// Login as the user that we just created.
	sessionId = testContext.TryAuthenticate(userId, "password1")
	testContext.sessionId = sessionId
	
	// Test ability to create a realm.
	var jrealm1Id string = testContext.TryCreateRealm("John's First Realm",
		"John's Little Outfit", "john")
	testContext.assertThat(jrealm1Id != "", "TryCreateRealm failed")
	
	// Test ability to create a realm.
	var jrealm2Id string = testContext.TryCreateRealm("John's Second Realm",
		"John's Next Venture", "admin")
	testContext.assertThat(jrealm2Id != "", "TryCreateRealm failed")
	
	var sarahConnorUserObjId string
	var sarahConnorAdminRealms []interface{}
	sarahConnorUserObjId, sarahConnorAdminRealms = testContext.TryCreateUser("sconnor", 
		"Sarah Connor", "sarahc@sky.net", "I'llMakePancakes", jrealm2Id)
	testContext.assertThat(sarahConnorUserObjId != "", "TryCreateUser failed")
	testContext.assertThat(len(sarahConnorAdminRealms) == 0, "Wrong number of admin realms")
	
	var johnConnorUserObjId string
	var johnConnorAdminRealms []interface{}
	johnConnorUserObjId, johnConnorAdminRealms = testContext.TryCreateUser("jconnor",
		"John Connor", "johnc@sky.net", "ILoveCameron", jrealm2Id)
	testContext.assertThat(johnConnorUserObjId != "", "TryCreateUser failed")
	testContext.assertThat(len(johnConnorAdminRealms) == 0, "Wrong number of admin realms")
	
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
	
	// Test ability to retrieve user by user id from realm.
	var userObjId string
	var userAdminRealms []interface{}
	userObjId, userAdminRealms = testContext.TryGetRealmUser(realmId, userId)
	testContext.assertThat(userObjId == johnDoeUserObjId, "Looking up user by user id failed")
	testContext.assertThat(len(userAdminRealms) == 2, "Wrong number of admin realms")
	
	//var msg string = testContext.TryAddRealmUser(....realmId, userObjId)
	
	var repoIds []string = testContext.TryGetRealmRepos(realmId)
	testContext.assertThat(len(repoIds) == 2, "Number of repo Ids returned was " +
		string(len(repoIds)) + ", expected 2")
	
	var realmIds []string = testContext.TryGetAllRealms()
	// Assumes that server is in debug mode, which creates test data.
	testContext.assertThat(len(realmIds) == 5, "Wrong number of realms found")
	
	userObjId, userAdminRealms = testContext.TryGetMyDesc()
	testContext.assertThat(userObjId == johnDoeUserObjId,
		"Returned user obj id was " + userObjId)
	testContext.assertThat(len(userAdminRealms) == 2, "Wrong number of admin realms")
	
	
	var myRealms []string = testContext.TryGetMyRealms()
	testContext.assertThat(len(myRealms) == 2, fmt.Sprintf(
		"Only returned %d realms", len(myRealms)))
	
	
	var myRepos []string = testContext.TryGetMyRepos()
	testContext.assertThat(len(myRepos) == 2, fmt.Sprintf(
		"Only returned %d repos", len(myRepos)))
	
	var realm3Id string
	var user3Id string
	var user3AdminRealms []interface{}
	var saveSessionId string = testContext.sessionId
	realm3Id, user3Id, user3AdminRealms = testContext.TryCreateRealmAnon("Realm3", "Realm 3 Org",
		"realm3admin", "Realm 3 Admin Full Name", "realm3admin@gmail.com", "realm3adminpswd")
	testContext.assertThat(realm3Id != "", "Realm Id is empty")
	testContext.assertThat(user3Id != "", "User Id is empty")
	testContext.assertThat(len(user3AdminRealms) == 1, "Wrong number of admin realms")
	
	// Restore user context to what it was before we called TryCreateRealmAnon.
	testContext.sessionId = saveSessionId
	
	// Test ability to build image from a dockerfile.
	var dockerImageObjId string
	var imageId string
	dockerImageObjId, imageId = testContext.TryExecDockerfile(repoId, dockerfileId, "myimage")
	testContext.assertThat(dockerImageObjId != "", "TryExecDockerfile failed - obj id is nil")
	testContext.assertThat(imageId != "", "TryExecDockerfile failed - docker image id is nil")
	
	// Test ability to list the images in a repo.
	var imageNames []string = testContext.TryGetImages(repoId)
	testContext.assertThat(len(imageNames) == 1, "Wrong number of images")
	
	var myDockerfileIds []string = testContext.TryGetMyDockerfiles()
	testContext.assertThat(len(myDockerfileIds) == 1, "Wrong number of dockerfiles")
	
	var myDockerImageIds []string = testContext.TryGetMyDockerImages()
	testContext.assertThat(len(myDockerImageIds) == 1, "Wrong number of docker images")
	
	var realmUsers []string = testContext.TryGetRealmUsers(jrealm2Id)
	testContext.assertThat(len(realmUsers) == 2, "Wrong number of realm users")
	
	var group1Id string = testContext.TryCreateGroup(jrealm2Id, "MyGroup", "For Overthrowning Skynet")
	testContext.assertThat(group1Id != "", "Empty group Id returned")
	
	var success bool = testContext.TryAddGroupUser(group1Id, johnConnorUserObjId)
	testContext.assertThat(success, "TryAddGroupUser failed")
	
	success = testContext.TryAddGroupUser(group1Id, sarahConnorUserObjId)
	testContext.assertThat(success, "TryAddGroupUser failed")
	
	var myGroupUsers []string = testContext.TryGetGroupUsers(group1Id)
	testContext.assertThat(len(myGroupUsers) == 2, "Wrong number of group users")
	
	var jrealm2GroupIds []string = testContext.TryGetRealmGroups(jrealm2Id)
	testContext.assertThat(len(jrealm2GroupIds) == 1, "Wrong number of realm groups")
	
	var myObjId string
	var myAdminRealms []interface{}
	myObjId, myAdminRealms = testContext.TryGetMyDesc()
	testContext.assertThat(len(myAdminRealms) == 2, "Wrong number of admin realms")
	
	success = testContext.TryAddGroupUser(group1Id, myObjId)
	testContext.assertThat(success, "TryAddGroupUser failed")
	
	var myGroupIds []string = testContext.TryGetMyGroups()
	testContext.assertThat(len(myGroupIds) == 1, "Wrong number of groups")
	
	var perms1 []bool = []bool{false, true, false, true, true}
	var retPerms1 []bool = testContext.TrySetPermission(user3Id, dockerfileId, perms1)
	var expectedPerms1 []bool = []bool{false, true, false, true, true}
	for i, p := range retPerms1 {
		testContext.assertThat(p == expectedPerms1[i], "Returned permission does not match")
	}
	
	var perms2 []bool = testContext.TryGetPermission(user3Id, dockerfileId)
	for i, p := range perms1 {
		testContext.assertThat(p == perms2[i], "Returned permission does not match")
	}
		
	var perms3 []bool = []bool{false, false, true, true, true}
	var retPerms3 []bool = testContext.TryAddPermission(user3Id, dockerfileId, perms3)
	var expectedPerms3 []bool = []bool{false, true, true, true, true}
	for i, p := range retPerms3 {
		testContext.assertThat(p == expectedPerms3[i], "Returned permission does not match")
	}
	
	testContext.TryReplaceDockerfile()
	
	
	testContext.TryDownloadImage()
	
	
	testContext.TryDeleteUser()
	
	
	testContext.TryDeleteGroup()
	
	
	testContext.TryRemGroupUser()
	
	
	testContext.TryDeleteRealm()
	
	
	testContext.TryRemRealmUser()
	
	
	testContext.TryDeleteRepo()
	
	
	testContext.TryRemPermission()

	//....testContext.TryDefineQualityScan(Script, SuccessGraphicImageURL, FailureGraphicImageURL)

	//var scriptId string = ""
	//var scanMessage string = testContext.TryScanImage(scriptId, dockerImageObjId)
	//testContext.assertThat(scanMessage != "", "Empty scan result message")

	// Test that permissions work.
	
	
	
	
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
func (testContext *TestContext) TryCreateRealm(realmName, orgFullName,
	adminUserId string) string {
	
	testContext.StartTest("TryCreateRealm")
	var resp *http.Response = testContext.sendPost(testContext.sessionId,
		"createRealm",
		[]string{"RealmName", "OrgFullName", "AdminUserId"},
		[]string{realmName, orgFullName, adminUserId})
	
	defer resp.Body.Close()
	
	testContext.verify200Response(resp)
	
	// Get the realm Id that is returned in the response body.
	var responseMap map[string]interface{}
	responseMap  = parseResponseBodyToMap(resp.Body)
	var retId string = responseMap["Id"].(string)
	var retName string = responseMap["RealmName"].(string)
	var retOrgFullName string = responseMap["OrgFullName"].(string)
	var retAdminUserId string = responseMap["AdminUserId"].(string)
	printMap(responseMap)
	testContext.assertThat(retId != "", "Realm Id not found in response body")
	testContext.assertThat(retName != "", "Realm Name not found in response body")
	testContext.assertThat(retOrgFullName != "", "Realm OrgFullName not found in response body")
	testContext.assertThat(retAdminUserId != "", "Realm AdminUserId not found in response body")
	
	return retId
}

/*******************************************************************************
 * Return the object Id of the new user.
 */
func (testContext *TestContext) TryCreateUser(userId string, userName string,
	email string, pswd string, realmId string) (string, []interface{}) {
	testContext.StartTest("TryCreateUser")
	
	var resp *http.Response = testContext.sendPost(testContext.sessionId,
		"createUser",
		[]string{"UserId", "UserName", "EmailAddress", "Password", "RealmId"},
		[]string{userId, userName, email, pswd, realmId})
	
	defer resp.Body.Close()

	testContext.verify200Response(resp)
	
	var responseMap map[string]interface{}
	responseMap  = parseResponseBodyToMap(resp.Body)
	var retUserObjId string = responseMap["Id"].(string)
	var retUserId string = responseMap["UserId"].(string)
	var retUserName string = responseMap["UserName"].(string)
	var retRealmId string = responseMap["RealmId"].(string)
	var retCanModifyTheseRealms []interface{} = responseMap["CanModifyTheseRealms"].([]interface{})
	printMap(responseMap)
	
	testContext.assertThat(retUserObjId != "", "User obj Id not returned")
	testContext.assertThat(retUserId == userId, "Returned user id, " + retUserId +
		" does not match the original user id")
	testContext.assertThat(retUserName == userName, "Returned user name, " + retUserName +
		" does not match the original user name")
	testContext.assertThat(retRealmId == realmId, "Returned realm Id, " + retRealmId +
		" does not match the original realm Id")
	testContext.assertThat(retCanModifyTheseRealms != nil, "No realms returned")
	
	return retUserObjId, retCanModifyTheseRealms
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
	printMap(responseMap)
	var retSessionId string = responseMap["UniqueSessionId"].(string)
	var retUserId string = responseMap["AuthenticatedUserid"].(string)
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
	var repoName string = responseMap["RepoName"].(string)
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
		var dockerfileName string = responseMap["DockerfileName"].(string)

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
func (testContext *TestContext) TryGetRealmUser(realmId, userId string) (string, []interface{}) {
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
	var retCanModifyTheseRealms []interface{} = responseMap["CanModifyTheseRealms"].([]interface{})
	printMap(responseMap)
	
	testContext.assertThat(retUserObjId != "", "User obj Id not returned")
	testContext.assertThat(retUserId == userId, "Returned user id, " + retUserId +
		" does not match the original user id")
	testContext.assertThat(retUserName != "", "Returned user name is blank")
	testContext.assertThat(retRealmId == realmId, "Returned realm Id, " + retRealmId +
		" does not match the original realm Id")
	testContext.assertThat(retCanModifyTheseRealms != nil, "No realms returned")
	
	return retUserObjId, retCanModifyTheseRealms
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryCreateGroup(realmId, name, description string) string {
	testContext.StartTest("TryCreateGroup")
	
	var resp *http.Response = testContext.sendPost(testContext.sessionId,
		"createGroup",
		[]string{"RealmId", "Name", "Description"},
		[]string{realmId, name, description})
	
	defer resp.Body.Close()

	testContext.verify200Response(resp)
	
	var responseMap map[string]interface{}
	responseMap  = parseResponseBodyToMap(resp.Body) // returns GroupDesc
	// Id
	// Name
	// Description
	var retGroupId string = responseMap["GroupId"].(string)
	var retRealmId string = responseMap["RealmId"].(string)
	var retName string = responseMap["GroupName"].(string)
	var retCreationDate string = responseMap["CreationDate"].(string)
	var retDescription string = responseMap["Description"].(string)
	printMap(responseMap)
	
	testContext.assertThat(retGroupId != "", "Returned GroupId is empty")
	testContext.assertThat(retRealmId != "", "Returned RealmId is empty")
	testContext.assertThat(retName != "", "Returned Name is empty")
	testContext.assertThat(retCreationDate != "", "Returned CreationDate is empty")
	testContext.assertThat(retDescription != "", "Returned Description is empty")
	
	return retGroupId
}

/*******************************************************************************
 * Return an array of the user object ids.
 */
func (testContext *TestContext) TryGetGroupUsers(groupId string) []string {
	testContext.StartTest("TryGetGroupUsers")

	var resp *http.Response = testContext.sendPost(testContext.sessionId,
		"getGroupUsers",
		[]string{"GroupId"},
		[]string{groupId})
	
	defer resp.Body.Close()

	testContext.verify200Response(resp)
	
	var responseMaps []map[string]interface{}
	responseMaps  = parseResponseBodyToMaps(resp.Body)  // returns [UserDesc]
	var result []string = make([]string, 0)
	for _, responseMap := range responseMaps {
		printMap(responseMap)
		var retId string = responseMap["Id"].(string)
		var retUserId string = responseMap["UserId"].(string)
		var retUserName string = responseMap["UserName"].(string)
		var retRealmId string = responseMap["RealmId"].(string)
		var retCanModifyTheseRealms []interface{} = responseMap["CanModifyTheseRealms"].([]interface{})
	
		testContext.assertThat(retId != "", "Returned Id is empty")
		testContext.assertThat(retUserId != "", "Returned UserId is empty")
		testContext.assertThat(retUserName != "", "Returned UserName is empty")
		testContext.assertThat(retRealmId != "", "Returned RealmId is empty")
		testContext.assertThat(retCanModifyTheseRealms != nil, "No realms returned")
		result = append(result, retId)
	}
	
	return result
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryAddGroupUser(groupId, userId string) bool {
	testContext.StartTest("TryAddGroupUser")

	var resp *http.Response = testContext.sendPost(testContext.sessionId,
		"addGroupUser",
		[]string{"GroupId", "UserObjId"},
		[]string{groupId, userId})
	
	defer resp.Body.Close()

	testContext.verify200Response(resp)
	
	var responseMap map[string]interface{}
	responseMap  = parseResponseBodyToMap(resp.Body)  // returns Result
	// Status - A value of “0” indicates success.
	// Message
	var retStatus string = responseMap["Status"].(string)
	var retMessage string = responseMap["Message"].(string)
	printMap(responseMap)
	
	testContext.assertThat(retStatus != "", "Returned Status is empty")
	testContext.assertThat(retMessage != "", "Returned Message is empty")
	
	return true
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
func (testContext *TestContext) TryGetRealmGroups(realmId string) []string {
	testContext.StartTest("TryGetRealmGroups")

	var resp *http.Response = testContext.sendPost(testContext.sessionId,
		"getRealmGroups",
		[]string{"RealmId"},
		[]string{realmId})
	
	defer resp.Body.Close()

	testContext.verify200Response(resp)
	
	var responseMaps []map[string]interface{}
	responseMaps  = parseResponseBodyToMaps(resp.Body)  // returns [GroupDesc]
	// Id
	// Name
	// Description
	var result []string = make([]string, 0)
	for _, responseMap := range responseMaps {
		printMap(responseMap)
		var retGroupId string = responseMap["GroupId"].(string)
		var retRealmId string = responseMap["RealmId"].(string)
		var retName string = responseMap["GroupName"].(string)
		var retCreationDate string = responseMap["CreationDate"].(string)
		var retDescription string = responseMap["Description"].(string)
	
		testContext.assertThat(retGroupId != "", "Returned GroupId is empty")
		testContext.assertThat(retRealmId != "", "Returned RealmId is empty")
		testContext.assertThat(retName != "", "Returned group Name is empty")
		testContext.assertThat(retCreationDate != "", "Returned CreationDate is empty")
		testContext.assertThat(retDescription != "", "Returned group Description is empty")
		result = append(result, retGroupId)
	}
	
	return result
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
		var retName string = responseMap["RepoName"].(string)
	
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
		var retName string = responseMap["RealmName"].(string)
	
		testContext.assertThat(retRealmId != "", "Returned realm Id is empty string")
		testContext.assertThat(retName != "", "Empty returned Name")
		
		result = append(result, retRealmId)
	}
	return result
}

/*******************************************************************************
 * Returns the Ids of the dockerfiles.
 */
func (testContext *TestContext) TryGetMyDockerfiles() []string {
	testContext.StartTest("TryGetMyDockerfiles")
	
	var resp *http.Response = testContext.sendPost(testContext.sessionId,
		"getMyDockerfiles",
		[]string{},
		[]string{})
	
	defer resp.Body.Close()

	testContext.verify200Response(resp)
	
	var responseMaps []map[string]interface{} = parseResponseBodyToMaps(resp.Body)
	var result []string = make([]string, 0)
	for _, responseMap := range responseMaps {
		printMap(responseMap)
		var retId string = responseMap["Id"].(string)
		var retName string = responseMap["DockerfileName"].(string)
	
		testContext.assertThat(retId != "", "Returned Id is empty string")
		testContext.assertThat(retName != "", "Returned Name is empty string")
		
		result = append(result, retId)
	}
	return result
}

/*******************************************************************************
 * Returns the Ids of the image objects.
 */
func (testContext *TestContext) TryGetMyDockerImages() []string {
	testContext.StartTest("TryGetMyDockerImages")
	
	var resp *http.Response = testContext.sendPost(testContext.sessionId,
		"getMyDockerImages",
		[]string{},
		[]string{})
	
	defer resp.Body.Close()

	testContext.verify200Response(resp)
	
	var responseMaps []map[string]interface{} = parseResponseBodyToMaps(resp.Body)
	var result []string = make([]string, 0)
	for _, responseMap := range responseMaps {
		printMap(responseMap)
		var retObjId string = responseMap["ObjId"].(string)
		var retDockerImageTag string = responseMap["DockerImageTag"].(string)
	
		testContext.assertThat(retObjId != "", "Returned ObjId is empty string")
		testContext.assertThat(retDockerImageTag != "", "Returned DockerImageTag is empty string")
		
		result = append(result, retObjId)
	}
	return result
}

/*******************************************************************************
 * Returns the obj Ids of the realm's users.
 */
func (testContext *TestContext) TryGetRealmUsers(realmId string) []string {
	testContext.StartTest("TryGetRealmUsers")
	
	var resp *http.Response = testContext.sendPost(testContext.sessionId,
		"getRealmUsers",
		[]string{"RealmId"},
		[]string{realmId})
	
	defer resp.Body.Close()

	testContext.verify200Response(resp)
	
	var responseMaps []map[string]interface{}
	responseMaps  = parseResponseBodyToMaps(resp.Body)
	var result []string = make([]string, 0)
	for _, responseMap := range responseMaps {
		var retId string = responseMap["Id"].(string)
		var retGroupId string = responseMap["UserId"].(string)
		var retUserName string = responseMap["UserName"].(string)
		var retRealmId string = responseMap["RealmId"].(string)
		var retCanModifyTheseRealms []interface{} = responseMap["CanModifyTheseRealms"].([]interface{})
		printMap(responseMap)
		testContext.assertThat(retId != "", "Empty Id returned")
		testContext.assertThat(retUserName != "", "Empty UserName returned")
		testContext.assertThat(retGroupId != "", "Empty GroupId returned")
		testContext.assertThat(retRealmId != "", "Empty RealmId returned")
		testContext.assertThat(retCanModifyTheseRealms != nil, "No realms returned")
		result = append(result, retId)
	}
	return result
}

/*******************************************************************************
 * Returns the (Id, Id) of the created realm and user, respectively
 */
func (testContext *TestContext) TryCreateRealmAnon(realmName, orgFullName, adminUserId,
	adminUserFullName, adminEmailAddr, adminPassword string) (string, string, []interface{}) {
	testContext.StartTest("TryCreateRealmAnon")
	
	var resp1 *http.Response = testContext.sendPost(testContext.sessionId,
		"createRealmAnon",
		[]string{"UserId", "UserName", "EmailAddress", "Password", "RealmName", "OrgFullName"},
		[]string{adminUserId, adminUserFullName, adminEmailAddr, adminPassword,
			realmName, orgFullName})
	
		// Returns UserDesc, which contains:
		// Id string
		// UserId string
		// UserName string
		// RealmId string

	defer resp1.Body.Close()

	testContext.verify200Response(resp1)
	
	var response1Map map[string]interface{}
	response1Map  = parseResponseBodyToMap(resp1.Body)
	printMap(response1Map)

	var retId string = response1Map["Id"].(string)
	var retUserId string = response1Map["UserId"].(string)
	var retUserName string = response1Map["UserName"].(string)
	var retRealmId string = response1Map["RealmId"].(string)
	var retCanModifyTheseRealms []interface{} = response1Map["CanModifyTheseRealms"].([]interface{})
	testContext.assertThat(retId != "", "Empty return Id")
	testContext.assertThat(retUserId != "", "Empty return UserId")
	testContext.assertThat(retUserName != "", "Empty return UserName")
	testContext.assertThat(retRealmId != "", "Empty return RealmId")
	testContext.assertThat(retCanModifyTheseRealms != nil, "No realms returned")
	
	// Authenticate as the admin user that was just created.
	var resp2 *http.Response = testContext.sendPost(testContext.sessionId,
		"authenticate",
		[]string{"UserId", "Password"},
		[]string{adminUserId, adminPassword})
	
	defer resp2.Body.Close()
	var response2Map map[string]interface{}
	response2Map  = parseResponseBodyToMap(resp2.Body)
	printMap(response2Map)
	var ret2SessionId string = response2Map["UniqueSessionId"].(string)
	var ret2UserId string = response2Map["AuthenticatedUserid"].(string)
	testContext.assertThat(ret2SessionId != "", "Session id is empty string")
	testContext.assertThat(ret2UserId == adminUserId, "Returned user id '" + ret2UserId + "' does not match user id")

	testContext.verify200Response(resp2)	
	testContext.sessionId = ret2SessionId
	
	// Now retrieve the description of the realm that we just created.
	var resp3 *http.Response = testContext.sendPost(testContext.sessionId,
		"getRealmDesc",
		[]string{"RealmId"},
		[]string{retRealmId})
	
		// Returns RealmDesc, which contains:
		// Id
		// Name
		// OrgFullName
	
	defer resp3.Body.Close()

	testContext.verify200Response(resp3)
	
	var response3Map map[string]interface{}
	response3Map  = parseResponseBodyToMap(resp3.Body)
	var ret3Id string = response3Map["Id"].(string)
	var ret3Name string = response3Map["RealmName"].(string)
	var ret3OrgFullName string = response3Map["OrgFullName"].(string)
	printMap(response3Map)
	testContext.assertThat(ret3Id != "", "Empty return Id")
	testContext.assertThat(ret3Name != "", "Empty return Name")
	testContext.assertThat(ret3OrgFullName != "", "Empty return Org Full Name")
	
	return ret3Id, retId, retCanModifyTheseRealms
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
 * Returns the permissions that resulted.
 */
func (testContext *TestContext) TrySetPermission(partyId, resourceId string,
	permissions []bool) []bool {

	testContext.StartTest("TrySetPermission")
	
	var resp *http.Response = testContext.sendPost(testContext.sessionId,
		"setPermission",
		[]string{"PartyId", "ResourceId", "Create", "Read", "Write", "Execute", "Delete"},
		[]string{partyId, resourceId, boolToString(permissions[0]),
			boolToString(permissions[1]), boolToString(permissions[2]),
			boolToString(permissions[3]), boolToString(permissions[4])})
	
	defer resp.Body.Close()

	testContext.verify200Response(resp)
	
	var responseMap map[string]interface{}
	responseMap  = parseResponseBodyToMap(resp.Body)
	printMap(responseMap)

	var retACLEntryId string = responseMap["ACLEntryId"].(string)
	var retPartyId string = responseMap["PartyId"].(string)
	var retResourceId string = responseMap["ResourceId"].(string)
	var retMask []bool = make([]bool, 5)
	retMask[0] = responseMap["Create"].(bool)
	retMask[1] = responseMap["Read"].(bool)
	retMask[2] = responseMap["Write"].(bool)
	retMask[3] = responseMap["Execute"].(bool)
	retMask[4] = responseMap["Delete"].(bool)
	testContext.assertThat(retACLEntryId != "", "Empty return retACLEntryId")
	testContext.assertThat(retPartyId != "", "Empty return retPartyId")
	testContext.assertThat(retResourceId != "", "Empty return retResourceId")
	
	return retMask
}

/*******************************************************************************
 * Returns the permissions that resulted.
 */
func (testContext *TestContext) TryAddPermission(partyId, resourceId string,
	permissions []bool) []bool {

	testContext.StartTest("TryAddPermission")
	
	var resp *http.Response = testContext.sendPost(testContext.sessionId,
		"addPermission",
		[]string{"PartyId", "ResourceId", "Create", "Read", "Write", "Execute", "Delete"},
		[]string{partyId, resourceId, boolToString(permissions[0]),
			boolToString(permissions[1]), boolToString(permissions[2]),
			boolToString(permissions[3]), boolToString(permissions[4])})
	
	defer resp.Body.Close()

	testContext.verify200Response(resp)
	
	var responseMap map[string]interface{}
	responseMap  = parseResponseBodyToMap(resp.Body)
	printMap(responseMap)

	var retACLEntryId string = responseMap["ACLEntryId"].(string)
	var retPartyId string = responseMap["PartyId"].(string)
	var retResourceId string = responseMap["ResourceId"].(string)
	var retMask []bool = make([]bool, 5)
	retMask[0] = responseMap["Create"].(bool)
	retMask[1] = responseMap["Read"].(bool)
	retMask[2] = responseMap["Write"].(bool)
	retMask[3] = responseMap["Execute"].(bool)
	retMask[4] = responseMap["Delete"].(bool)
	testContext.assertThat(retACLEntryId != "", "Empty return retACLEntryId")
	testContext.assertThat(retPartyId != "", "Empty return retPartyId")
	testContext.assertThat(retResourceId != "", "Empty return retResourceId")
	
	return retMask
}

/*******************************************************************************
 * Return an array of string representing the values for the permission mask.
 */
func (testContext *TestContext) TryGetPermission(partyId, resourceId string) []bool {

	testContext.StartTest("TryGetPermission")
	
	var resp *http.Response = testContext.sendPost(testContext.sessionId,
		"getPermission",
		[]string{"PartyId", "ResourceId"},
		[]string{partyId, resourceId})
	
	defer resp.Body.Close()

	testContext.verify200Response(resp)
	
	var responseMap map[string]interface{}
	responseMap  = parseResponseBodyToMap(resp.Body)
	printMap(responseMap)
	
	var retACLEntryId string = responseMap["ACLEntryId"].(string)
	var retPartyId string = responseMap["PartyId"].(string)
	var retResourceId string = responseMap["ResourceId"].(string)
	var retCreate bool = responseMap["Create"].(bool)
	var retRead bool = responseMap["Read"].(bool)
	var retWrite bool = responseMap["Write"].(bool)
	var retExecute bool = responseMap["Execute"].(bool)
	var retDelete bool = responseMap["Delete"].(bool)
	testContext.assertThat(retACLEntryId != "", "Empty return retACLEntryId")
	testContext.assertThat(retPartyId != "", "Empty return retPartyId")
	testContext.assertThat(retResourceId != "", "Empty return retResourceId")
	
	return []bool{retCreate, retRead, retWrite, retExecute, retDelete}
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryDefineQualityScan(Script,
	SuccessGraphicImageURL, FailureGraphicImageURL string) {
	testContext.StartTest("TryDefineQualityScan")
}

/*******************************************************************************
 * Returns output message.
 */
func (testContext *TestContext) TryScanImage(scriptId, imageObjId string) string {
	testContext.StartTest("TryScanImage")
	
	var resp *http.Response = testContext.sendPost(testContext.sessionId,
		"scanImage",
		[]string{"ScriptId", "ImageObjId"},
		[]string{scriptId, imageObjId})
	
	defer resp.Body.Close()

	testContext.verify200Response(resp)
	
	var responseMap map[string]interface{}
	responseMap  = parseResponseBodyToMap(resp.Body)
	printMap(responseMap)
	
	var msg string = responseMap["Message"].(string)
	return msg
}

/*******************************************************************************
 * Return the object Id of the current authenticated user.
 */
func (testContext *TestContext) TryGetMyDesc() (string, []interface{}) {
	testContext.StartTest("TryGetMyDesc")
	
	var resp *http.Response = testContext.sendPost(testContext.sessionId,
		"getMyDesc",
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
	var retCanModifyTheseRealms []interface{} = responseMap["CanModifyTheseRealms"].([]interface{})

	testContext.assertThat(retId != "", "Returned Id is empty string")
	testContext.assertThat(retUserId != "", "Returned UserId is empty string")
	testContext.assertThat(retUserName != "", "Returned UserName is empty string")
	testContext.assertThat(retRealmId != "", "Returned RealmId is empty string")
	testContext.assertThat(retCanModifyTheseRealms != nil, "No realms returned")
	
	return retId, retCanModifyTheseRealms
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryGetMyGroups() []string {
	testContext.StartTest("TryGetMyGroups")
	
	var resp *http.Response = testContext.sendPost(testContext.sessionId,
		"getMyGroups",
		[]string{},
		[]string{})
	
	defer resp.Body.Close()

	testContext.verify200Response(resp)
	
	var responseMaps []map[string]interface{} = parseResponseBodyToMaps(resp.Body)
	var result []string = make([]string, 0)
	for _, responseMap := range responseMaps {
		printMap(responseMap)
		var retGroupId string = responseMap["GroupId"].(string)
		var retRealmId string = responseMap["RealmId"].(string)
		var retName string = responseMap["GroupName"].(string)
		var retCreationDate string = responseMap["CreationDate"].(string)
		var retDescription string = responseMap["Description"].(string)
		testContext.assertThat(retGroupId != "", "Returned GroupId is empty string")
		testContext.assertThat(retRealmId != "", "Empty returned RealmId")
		testContext.assertThat(retName != "", "Empty returned Name")
		testContext.assertThat(retCreationDate != "", "Empty CreationDate returned")
		testContext.assertThat(retDescription != "", "Empty returned Description")
		
		result = append(result, retGroupId)
	}
	return result
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
		var retName string = responseMap["RealmName"].(string)
	
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
		var retName string = responseMap["RepoName"].(string)
	
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
