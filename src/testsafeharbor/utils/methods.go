package utils

import (
	"fmt"
	"net/http"
	"os"
	"io"
	//"flag"
	"testsafeharbor/rest"
	//"testsafeharbor/utils"
)

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryGetGroupDesc(groupId string) {
	
	testContext.StartTest("getGroupDesc")
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"getGroupDesc",
		[]string{"GroupId"},
		[]string{groupId})
	
	defer resp.Body.Close()
	
	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if testContext.AssertErrIsNil(err, "at ParseResponseBodyToMap") { return }
	
	// Expect a GroupDesc
	var retGroupId string = responseMap["GroupId"].(string)
	var retRealmId string = responseMap["RealmId"].(string)
	var retGroupName string = responseMap["GroupName"].(string)
	var retCreationDate string = responseMap["CreationDate"].(string)
	var retDescription string = responseMap["Description"].(string)
	
	testContext.AssertThat(retGroupId != "", "retGroupId is empty")
	testContext.AssertThat(retRealmId != "", "retRealmId is empty")
	testContext.AssertThat(retGroupName != "", "retGroupName is empty")
	testContext.AssertThat(retCreationDate != "", "retCreationDate is empty")
	testContext.AssertThat(retDescription != "", "retDescription is empty")
}
	
/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryGetRepoDesc(repoId string) {
	
	testContext.StartTest("getRepoDesc")
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"getRepoDesc",
		[]string{"RepoId"},
		[]string{repoId})
	
	defer resp.Body.Close()
	
	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if testContext.AssertErrIsNil(err, "at ParseResponseBodyToMap") { return }
	
	// Expect a RepoDesc
	var retId string = responseMap["Id"].(string)
	var retRealmId string = responseMap["RealmId"].(string)
	var retRepoName string = responseMap["RepoName"].(string)
	var retDescription string = responseMap["Description"].(string)
	var retCreationDate string = responseMap["CreationDate"].(string)
	var retDockerfileIds []string = responseMap["DockerfileIds"].([]string)
	
	testContext.AssertThat(retId != "", "retId is empty")
	testContext.AssertThat(retRealmId != "", "retRealmId is empty")
	testContext.AssertThat(retRepoName != "", "retRepoName is empty")
	testContext.AssertThat(retDescription != "", "retDescription is empty")
	testContext.AssertThat(retCreationDate != "", "retCreationDate is empty")
	testContext.AssertThat(retDockerfileIds != nil, "retDockerfileIds is nil")
}
	
/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryGetDockerImageDesc(dockerImageId string,
	expectSuccess bool) map[string]interface{} {
	
	testContext.StartTest("getDockerImageDesc")
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"getDockerImageDesc",
		[]string{"DockerImageId"},
		[]string{dockerImageId})
	
	defer resp.Body.Close()
	
	if expectSuccess {
		if ! testContext.Verify200Response(resp) {
			testContext.FailTest()
			return nil
		}
	} else {
		if resp.StatusCode == 200 {
			testContext.FailTest()
			return nil
		} else {
			testContext.PassTest()
			return nil
		}	
	}
	
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if testContext.AssertErrIsNil(err, "at ParseResponseBodyToMap") { return nil }
	
	// Expect a DockerImageDesc
	var retObjId string = responseMap["ObjId"].(string)
	var retRepoId string = responseMap["RepoId"].(string)
	var retName string = responseMap["Name"].(string)
	var retDescription string = responseMap["Description"].(string)
	var retCreationDate string = responseMap["CreationDate"].(string)
	
	testContext.AssertThat(retObjId != "", "retObjId is empty")
	testContext.AssertThat(retRepoId != "", "retRepoId is empty")
	testContext.AssertThat(retName != "", "retName is empty")
	testContext.AssertThat(retDescription != "", "retDescription is empty")
	testContext.AssertThat(retCreationDate != "", "retCreationDate is empty")
	
	return responseMap
}
	
/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryGetDockerfileDesc(dockerfileId string) {
	
	testContext.StartTest("getDockerfileDesc")
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"getDockerfileDesc",
		[]string{"DockerfileId"},
		[]string{dockerfileId})
	
	defer resp.Body.Close()
	
	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if testContext.AssertErrIsNil(err, "at ParseResponseBodyToMap") { return }
	
	// Expect a DockerfileDesc
	var retId string = responseMap["Id"].(string)
	var retRepoId string = responseMap["RepoId"].(string)
	var retDescription string = responseMap["Description"].(string)
	var retDockerfileName string = responseMap["DockerfileName"].(string)
	
	testContext.AssertThat(retId != "", "retId is empty")
	testContext.AssertThat(retRepoId != "", "retRepoId is empty")
	testContext.AssertThat(retDescription != "", "retDescription is empty")
	testContext.AssertThat(retDockerfileName != "", "retDockerfileName is empty")
}

/*******************************************************************************
 * Verify that we can create a new realm.
 */
func (testContext *TestContext) TryCreateRealm(realmName, orgFullName,
	adminUserId string) string {
	
	testContext.StartTest("TryCreateRealm")
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"createRealm",
		[]string{"RealmName", "OrgFullName", "AdminUserId"},
		[]string{realmName, orgFullName, adminUserId})
	
	defer resp.Body.Close()
	
	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	// Get the realm Id that is returned in the response body.
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if err != nil { fmt.Println(err.Error()); return "" }
	var retId string = responseMap["Id"].(string)
	var retName string = responseMap["RealmName"].(string)
	var retOrgFullName string = responseMap["OrgFullName"].(string)
	var retAdminUserId string = responseMap["AdminUserId"].(string)
	rest.PrintMap(responseMap)
	testContext.AssertThat(retId != "", "Realm Id not found in response body")
	testContext.AssertThat(retName != "", "Realm Name not found in response body")
	testContext.AssertThat(retOrgFullName != "", "Realm OrgFullName not found in response body")
	testContext.AssertThat(retAdminUserId != "", "Realm AdminUserId not found in response body")
	
	return retId
}

/*******************************************************************************
 * Return the object Id of the new user.
 */
func (testContext *TestContext) TryCreateUser(userId string, userName string,
	email string, pswd string, realmId string) (string, []interface{}) {
	testContext.StartTest("TryCreateUser")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"createUser",
		[]string{"UserId", "UserName", "EmailAddress", "Password", "RealmId"},
		[]string{userId, userName, email, pswd, realmId})
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if err != nil { fmt.Println(err.Error()); return "", nil }
	var retUserObjId string = responseMap["Id"].(string)
	var retUserId string = responseMap["UserId"].(string)
	var retUserName string = responseMap["UserName"].(string)
	var retRealmId string = responseMap["RealmId"].(string)
	var retCanModifyTheseRealms []interface{} = responseMap["CanModifyTheseRealms"].([]interface{})
	rest.PrintMap(responseMap)
	
	testContext.AssertThat(retUserObjId != "", "User obj Id not returned")
	testContext.AssertThat(retUserId == userId, "Returned user id, " + retUserId +
		" does not match the original user id")
	testContext.AssertThat(retUserName == userName, "Returned user name, " + retUserName +
		" does not match the original user name")
	testContext.AssertThat(retRealmId == realmId, "Returned realm Id, " + retRealmId +
		" does not match the original realm Id")
	testContext.AssertThat(retCanModifyTheseRealms != nil, "No realms returned")
	
	return retUserObjId, retCanModifyTheseRealms
}

/*******************************************************************************
 * Returns session Id, and isAdmin.
 */
func (testContext *TestContext) TryAuthenticate(userId string, pswd string,
	expectSuccess bool) (string, bool) {
	testContext.StartTest("TryAuthenticate")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"authenticate",
		[]string{"UserId", "Password"},
		[]string{userId, pswd})
	
	defer resp.Body.Close()

	if expectSuccess {
		if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	} else {
		if resp.StatusCode == 200 {
			testContext.FailTest()
			return "", false
		} else {
			testContext.PassTest()
			return "", true
		}	
	}
	
	// Get the repo Id that is returned in the response body.
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if err != nil { fmt.Println(err.Error()); return "", false }
	rest.PrintMap(responseMap)
	var retReason interface{} = responseMap["Reason"]
	if retReason != nil { return "", false }
	var retSessionId string = responseMap["UniqueSessionId"].(string)
	var retUserId string = responseMap["AuthenticatedUserid"].(string)
	var retIsAdmin bool = responseMap["IsAdmin"].(bool)
	testContext.AssertThat(retSessionId != "", "Session id is empty string")
	testContext.AssertThat(retUserId == userId, "Returned user id '" + retUserId +
		"' does not match user id")
	return retSessionId, retIsAdmin
}

/*******************************************************************************
 * Return true if successful.
 */
func (testContext *TestContext) TryDisableUser(userObjId string) bool {
	testContext.StartTest("TryDisableUser")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"disableUser",
		[]string{"UserObjId"},
		[]string{userObjId})
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if err != nil { fmt.Println(err.Error()); return false }
	rest.PrintMap(responseMap)
	var retStatus string = responseMap["Status"].(string)
	//var retMessage string = responseMap["Message"].(string)
	if retStatus != "200" { return false }
	return testContext.CurrentTestPassed
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryDeleteGroup(groupId string) bool {
	testContext.StartTest("TryDeleteGroup")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"deleteGroup",
		[]string{"GroupId"},
		[]string{groupId})
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if err != nil { fmt.Println(err.Error()); return false }
	rest.PrintMap(responseMap)
	var retStatus string = responseMap["Status"].(string)
	//var retMessage string = responseMap["Message"].(string)
	if retStatus != "200" { return false }
	return testContext.CurrentTestPassed
}

/*******************************************************************************
 * If successful, return true.
 */
func (testContext *TestContext) TryLogout() bool {
	testContext.StartTest("TryLogout")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"logout",
		[]string{},
		[]string{})
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if err != nil { fmt.Println(err.Error()); return false }
	rest.PrintMap(responseMap)
	var retStatus string = responseMap["Status"].(string)
	//var retMessage string = responseMap["Message"].(string)
	if retStatus != "200" { return false }
	return testContext.CurrentTestPassed
}

/*******************************************************************************
 * Verify that we can create a new repo. This requires that we first created
 * a realm that the repo can belong to.
 */
func (testContext *TestContext) TryCreateRepo(realmId string, name string,
	desc string, optDockerfilePath string) string {
	testContext.StartTest("TryCreateRepo")
	
	var resp *http.Response
	var err error
	
	if optDockerfilePath == "" {
		fmt.Println("Using SendPost")
		resp, err = testContext.SendPost(testContext.SessionId,
			"createRepo",
			[]string{"RealmId", "Name", "Description"},
			[]string{realmId, name, desc})
	} else {
		fmt.Println("Using SendFilePost")
		resp, err = testContext.SendFilePost(testContext.SessionId,
			"createRepo",
			[]string{"RealmId", "Name", "Description"},
			[]string{realmId, name, desc},
			optDockerfilePath)
	}
	if err != nil { fmt.Println(err.Error()); return "" }
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	// Get the repo Id that is returned in the response body.
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if err != nil { fmt.Println(err.Error()); return "" }
	var repoId string = responseMap["Id"].(string)
	var repoName string = responseMap["RepoName"].(string)
	rest.PrintMap(responseMap)
	testContext.AssertThat(repoId != "", "Repo Id not found in response body")
	testContext.AssertThat(repoName != "", "Repo Name not found in response body")
	
	return repoId
}

/*******************************************************************************
 * Verify that we can upload a dockerfile. This requries that we first created
 * a repo to uplaod it into.
 */
func (testContext *TestContext) TryAddDockerfile(repoId string, dockerfilePath string,
	desc string) string {
	
	testContext.StartTest("TryAddDockerfile")
	fmt.Println("\t", dockerfilePath)
	var resp *http.Response
	var err error
	resp, err = testContext.SendFilePost(testContext.SessionId,
		"addDockerfile",
		[]string{"RepoId", "Description"},
		[]string{repoId, desc},
		dockerfilePath)
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	// Get the DockerfileDesc that is returned.
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if err != nil { fmt.Println(err.Error()); return "" }
	var dockerfileId string = responseMap["Id"].(string)
	//var dockerfileName string = responseMap["Name"]
	rest.PrintMap(responseMap)
	//AssertThat(dockerfileId != "", "Dockerfile Id not found in response body")
	//AssertThat(dockerfileName != "", "Dockerfile Name not found in response body")
	
	return dockerfileId
}

/*******************************************************************************
 * Verify that we can obtain the names of the dockerfiles owned by the specified
 * repo. The result is an array of dockerfile names.
 */
func (testContext *TestContext) TryGetDockerfiles(repoId string) []string {
	testContext.StartTest("TryGetDockerfiles")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"getDockerfiles",
		[]string{"RepoId"},
		[]string{repoId})
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMaps []map[string]interface{}
	responseMaps, err = rest.ParseResponseBodyToMaps(resp.Body)
	if err != nil { fmt.Println(err.Error()); return nil }
	var result []string = make([]string, 0)
	for _, responseMap := range responseMaps {
		var dockerfileId string = responseMap["Id"].(string)
		var repoId string = responseMap["RepoId"].(string)
		var dockerfileName string = responseMap["DockerfileName"].(string)

		rest.PrintMap(responseMap)
		testContext.AssertThat(dockerfileId != "", "Dockerfile Id not found in response body")
		testContext.AssertThat(repoId != "", "Repo Id not found in response body")
		testContext.AssertThat(dockerfileName != "", "Dockerfile Name not found in response body")
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
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"execDockerfile",
		[]string{"RepoId", "DockerfileId", "ImageName"},
		[]string{repoId, dockerfileId, imageName})
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	// Get the repo Id that is returned in the response body.
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if err != nil { fmt.Println(err.Error()); return "", "" }
	var retObjId string = responseMap["ObjId"].(string)
	var retDockerImageTag string = responseMap["Name"].(string)
	var retDesc string = responseMap["Description"].(string)
	var retCreationDate = responseMap["CreationDate"].(string)
	rest.PrintMap(responseMap)
	
	testContext.AssertThat(retObjId != "", "ObjId is empty")
	testContext.AssertThat(retDockerImageTag != "", "Name is empty")
	testContext.AssertThat(retDesc != "", "Description is empty")
	testContext.AssertThat(retCreationDate != "", "CreationDate is empty")
	return retObjId, retDockerImageTag
}

/*******************************************************************************
 * Verify that we can upload a dockerfile and build an image from it.
 * The result is the object id and docker id of the image that was built.
 */
func (testContext *TestContext) TryAddAndExecDockerfile(repoId string, desc string,
	imageName string, dockerfilePath string) (string, string) {
	testContext.StartTest("TryAddAndExecDockerfile")
	
	var resp *http.Response
	var err error
	//resp, err = testContext.SendFilePost(testContext.SessionId,
	resp, err = testContext.SendFilePost("",
		"addAndExecDockerfile",
		[]string{"RepoId", "Description", "ImageName", "SessionId"},
		[]string{repoId, desc, imageName, testContext.SessionId},
		dockerfilePath)

	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	// Get the repo Id that is returned in the response body.
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if err != nil { fmt.Println(err.Error()); return "", "" }
	var retObjId string = responseMap["ObjId"].(string)
	var retDockerImageTag string = responseMap["Name"].(string)
	var retDesc string = responseMap["Description"].(string)
	var retCreationDate = responseMap["CreationDate"].(string)
	rest.PrintMap(responseMap)
	
	testContext.AssertThat(retObjId != "", "ObjId is empty")
	testContext.AssertThat(retDockerImageTag != "", "Name is empty")
	testContext.AssertThat(retDesc != "", "Description is empty")
	testContext.AssertThat(retCreationDate != "", "CreationDate is empty")
	return retObjId, retDockerImageTag
}

/*******************************************************************************
 * Result is an array of the names of the images owned by the specified repo.
 */
func (testContext *TestContext) TryGetImages(repoId string) []string {
	testContext.StartTest("TryGetImages")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"getImages",
		[]string{"RepoId"},
		[]string{repoId})
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMaps []map[string]interface{}
	responseMaps, err = rest.ParseResponseBodyToMaps(resp.Body)
	if err != nil { fmt.Println(err.Error()); return nil }
	var result []string = make([]string, 0)
	for _, responseMap := range responseMaps {
		var objId string = responseMap["ObjId"].(string)
		var dockerImageTag string = responseMap["Name"].(string)

		rest.PrintMap(responseMap)
		testContext.AssertThat(objId != "", "ObjId not found in response body")
		testContext.AssertThat(dockerImageTag != "", "DockerImageTag not found in response body")
		fmt.Println()

		result = append(result, dockerImageTag)
	}
	
	return result
}

/*******************************************************************************
 * Return the object Id of the specified user, and a list of the realms that
 * the user can modify.
 */
func (testContext *TestContext) TryGetUserDesc(realmId, userId string) map[string]interface{} {
	testContext.StartTest("TryGetUserDesc")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"getUserDesc",
		[]string{"RealmId", "UserId"},
		[]string{realmId, userId})
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if ! testContext.AssertErrIsNil(err, "") { return nil }
	var retUserObjId string = responseMap["Id"].(string)
	var retUserId string = responseMap["UserId"].(string)
	var retUserName string = responseMap["UserName"].(string)
	var retRealmId string = responseMap["RealmId"].(string)
	var retCanModifyTheseRealms []interface{} = responseMap["CanModifyTheseRealms"].([]interface{})
	rest.PrintMap(responseMap)
	
	testContext.AssertThat(retUserObjId != "", "User obj Id not returned")
	testContext.AssertThat(retUserId == userId, "Returned user id, " + retUserId +
		" does not match the original user id")
	testContext.AssertThat(retUserName != "", "Returned user name is blank")
	testContext.AssertThat(retRealmId == realmId, "Returned realm Id, " + retRealmId +
		" does not match the original realm Id")
	testContext.AssertThat(retCanModifyTheseRealms != nil, "No realms returned")
	
	return responseMap
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryCreateGroup(realmId, name, description string,
	addMe bool) string {
	testContext.StartTest("TryCreateGroup")
	
	var addMeStr = "false"
	if addMe { addMeStr = "true" }
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"createGroup",
		[]string{"RealmId", "Name", "Description", "AddMe"},
		[]string{realmId, name, description, addMeStr})
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if err != nil { fmt.Println(err.Error()); return "" } // returns GroupDesc
	// Id
	// Name
	// Description
	var retGroupId string = responseMap["GroupId"].(string)
	var retRealmId string = responseMap["RealmId"].(string)
	var retName string = responseMap["GroupName"].(string)
	var retCreationDate string = responseMap["CreationDate"].(string)
	var retDescription string = responseMap["Description"].(string)
	rest.PrintMap(responseMap)
	
	testContext.AssertThat(retGroupId != "", "Returned GroupId is empty")
	testContext.AssertThat(retRealmId != "", "Returned RealmId is empty")
	testContext.AssertThat(retName != "", "Returned Name is empty")
	testContext.AssertThat(retCreationDate != "", "Returned CreationDate is empty")
	testContext.AssertThat(retDescription != "", "Returned Description is empty")
	
	return retGroupId
}

/*******************************************************************************
 * Return an array of the user object ids.
 */
func (testContext *TestContext) TryGetGroupUsers(groupId string) []string {
	testContext.StartTest("TryGetGroupUsers")

	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"getGroupUsers",
		[]string{"GroupId"},
		[]string{groupId})
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMaps []map[string]interface{}
	responseMaps, err = rest.ParseResponseBodyToMaps(resp.Body)  // returns [UserDesc]
	if err != nil { fmt.Println(err.Error()); return nil }
	var result []string = make([]string, 0)
	for _, responseMap := range responseMaps {
		rest.PrintMap(responseMap)
		var retId string = responseMap["Id"].(string)
		var retUserId string = responseMap["UserId"].(string)
		var retUserName string = responseMap["UserName"].(string)
		var retRealmId string = responseMap["RealmId"].(string)
		var retCanModifyTheseRealms []interface{} = responseMap["CanModifyTheseRealms"].([]interface{})
	
		testContext.AssertThat(retId != "", "Returned Id is empty")
		testContext.AssertThat(retUserId != "", "Returned UserId is empty")
		testContext.AssertThat(retUserName != "", "Returned UserName is empty")
		testContext.AssertThat(retRealmId != "", "Returned RealmId is empty")
		testContext.AssertThat(retCanModifyTheseRealms != nil, "No realms returned")
		result = append(result, retId)
	}
	
	return result
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryAddGroupUser(groupId, userId string) bool {
	testContext.StartTest("TryAddGroupUser")

	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"addGroupUser",
		[]string{"GroupId", "UserObjId"},
		[]string{groupId, userId})
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if err != nil { fmt.Println(err.Error()); return false }  // returns Result
	// Status - A value of “0” indicates success.
	// Message
	var retStatus string = responseMap["Status"].(string)
	var retMessage string = responseMap["Message"].(string)
	rest.PrintMap(responseMap)
	
	testContext.AssertThat(retStatus == "200", "Returned Status is empty")
	testContext.AssertThat(retMessage != "", "Returned Message is empty")
	
	return testContext.CurrentTestPassed
}

/*******************************************************************************
 * Returns result.
 */
func (testContext *TestContext) TryMoveUserToRealm(realmId string, userObjId string) bool {
	testContext.StartTest("TryMoveUserToRealm")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"moveUserToRealm",
		[]string{"UserObjId", "RealmId"},
		[]string{userObjId, realmId})
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if ! testContext.AssertErrIsNil(err, "") { return false }
	var retStatus string = responseMap["Status"].(string)
	var retMsg string = responseMap["Message"].(string)
	rest.PrintMap(responseMap)
	testContext.AssertThat(retStatus != "", "Empty return status")
	testContext.AssertThat(retMsg != "", "Empty return message")
	return testContext.CurrentTestPassed
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryGetRealmGroups(realmId string) []string {
	testContext.StartTest("TryGetRealmGroups")

	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"getRealmGroups",
		[]string{"RealmId"},
		[]string{realmId})
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMaps []map[string]interface{}
	responseMaps, err = rest.ParseResponseBodyToMaps(resp.Body)  // returns [GroupDesc]
	if err != nil { fmt.Println(err.Error()); return nil }
	// Id
	// Name
	// Description
	var result []string = make([]string, 0)
	for _, responseMap := range responseMaps {
		rest.PrintMap(responseMap)
		var retGroupId string = responseMap["GroupId"].(string)
		var retRealmId string = responseMap["RealmId"].(string)
		var retName string = responseMap["GroupName"].(string)
		var retCreationDate string = responseMap["CreationDate"].(string)
		var retDescription string = responseMap["Description"].(string)
	
		testContext.AssertThat(retGroupId != "", "Returned GroupId is empty")
		testContext.AssertThat(retRealmId != "", "Returned RealmId is empty")
		testContext.AssertThat(retName != "", "Returned group Name is empty")
		testContext.AssertThat(retCreationDate != "", "Returned CreationDate is empty")
		testContext.AssertThat(retDescription != "", "Returned group Description is empty")
		result = append(result, retGroupId)
	}
	
	return result
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryGetRealmRepos(realmId string) []string {
	testContext.StartTest("TryGetRealmRepos")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"getRealmRepos",
		[]string{"RealmId"},
		[]string{realmId})
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMaps []map[string]interface{}
	responseMaps, err = rest.ParseResponseBodyToMaps(resp.Body)
	if err != nil { fmt.Println(err.Error()); return nil }
	var result []string = make([]string, 0)
	for _, responseMap := range responseMaps {
		rest.PrintMap(responseMap)
		var retRepoId string = responseMap["Id"].(string)
		var retRealmId string = responseMap["RealmId"].(string)
		var retName string = responseMap["RepoName"].(string)
	
		testContext.AssertThat(retRepoId != "", "No repo Id returned")
		testContext.AssertThat(retRealmId == realmId, "returned realm Id is nil")
		testContext.AssertThat(retName != "", "Empty returned Name")
		
		result = append(result, retRepoId)
	}
	return result
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryGetAllRealms() []string {
	testContext.StartTest("TryGetAllRealms")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"getAllRealms",
		[]string{},
		[]string{})
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMaps []map[string]interface{}
	responseMaps, err = rest.ParseResponseBodyToMaps(resp.Body)
	if err != nil { fmt.Println(err.Error()); return nil }
	var result []string = make([]string, 0)
	for _, responseMap := range responseMaps {
		rest.PrintMap(responseMap)
		var retRealmId string = responseMap["Id"].(string)
		var retName string = responseMap["RealmName"].(string)
	
		testContext.AssertThat(retRealmId != "", "Returned realm Id is empty string")
		testContext.AssertThat(retName != "", "Empty returned Name")
		
		result = append(result, retRealmId)
	}
	return result
}

/*******************************************************************************
 * Returns the Ids of the dockerfiles.
 */
func (testContext *TestContext) TryGetMyDockerfiles() []string {
	testContext.StartTest("TryGetMyDockerfiles")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"getMyDockerfiles",
		[]string{},
		[]string{})
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMaps []map[string]interface{}
	responseMaps, err = rest.ParseResponseBodyToMaps(resp.Body)
	if err != nil { fmt.Println(err.Error()); return nil }
	var result []string = make([]string, 0)
	for _, responseMap := range responseMaps {
		rest.PrintMap(responseMap)
		var retId string = responseMap["Id"].(string)
		var retName string = responseMap["DockerfileName"].(string)
	
		testContext.AssertThat(retId != "", "Returned Id is empty string")
		testContext.AssertThat(retName != "", "Returned Name is empty string")
		
		result = append(result, retId)
	}
	return result
}

/*******************************************************************************
 * Returns the Ids of the image objects.
 */
func (testContext *TestContext) TryGetMyDockerImages() []string {
	testContext.StartTest("TryGetMyDockerImages")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"getMyDockerImages",
		[]string{},
		[]string{})
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMaps []map[string]interface{}
	responseMaps, err = rest.ParseResponseBodyToMaps(resp.Body)
	if err != nil { fmt.Println(err.Error()); return nil }
	var result []string = make([]string, 0)
	for _, responseMap := range responseMaps {
		rest.PrintMap(responseMap)
		var retObjId string = responseMap["ObjId"].(string)
		var retDockerImageTag string = responseMap["Name"].(string)
	
		testContext.AssertThat(retObjId != "", "Returned ObjId is empty string")
		testContext.AssertThat(retDockerImageTag != "", "Returned DockerImageTag is empty string")
		
		result = append(result, retObjId)
	}
	return result
}

/*******************************************************************************
 * Returns the obj Ids of the realm's users.
 */
func (testContext *TestContext) TryGetRealmUsers(realmId string) []string {
	testContext.StartTest("TryGetRealmUsers")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"getRealmUsers",
		[]string{"RealmId"},
		[]string{realmId})
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMaps []map[string]interface{}
	responseMaps, err = rest.ParseResponseBodyToMaps(resp.Body)
	if err != nil { fmt.Println(err.Error()); return nil }
	var result []string = make([]string, 0)
	for _, responseMap := range responseMaps {
		var retId string = responseMap["Id"].(string)
		var retGroupId string = responseMap["UserId"].(string)
		var retUserName string = responseMap["UserName"].(string)
		var retRealmId string = responseMap["RealmId"].(string)
		var retCanModifyTheseRealms []interface{} = responseMap["CanModifyTheseRealms"].([]interface{})
		rest.PrintMap(responseMap)
		testContext.AssertThat(retId != "", "Empty Id returned")
		testContext.AssertThat(retUserName != "", "Empty UserName returned")
		testContext.AssertThat(retGroupId != "", "Empty GroupId returned")
		testContext.AssertThat(retRealmId != "", "Empty RealmId returned")
		testContext.AssertThat(retCanModifyTheseRealms != nil, "No realms returned")
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
	
	var resp1 *http.Response
	var err error
	resp1, err = testContext.SendPost(testContext.SessionId,
		"createRealmAnon",
		[]string{"UserId", "UserName", "EmailAddress", "Password", "RealmName", "OrgFullName"},
		[]string{adminUserId, adminUserFullName, adminEmailAddr, adminPassword,
			realmName, orgFullName})
	
		// Returns UserDesc, which contains:
		// Id string
		// UserId string
		// UserName string
		// RealmId string
		
	if err != nil { fmt.Println(err.Error()); return "", "", nil }

	defer resp1.Body.Close()

	testContext.Verify200Response(resp1)
	
	var response1Map map[string]interface{}
	response1Map, err = rest.ParseResponseBodyToMap(resp1.Body)
	if err != nil { fmt.Println(err.Error()); return "", "", nil }
	rest.PrintMap(response1Map)

	var retId string = response1Map["Id"].(string)
	var retUserId string = response1Map["UserId"].(string)
	var retUserName string = response1Map["UserName"].(string)
	var retRealmId string = response1Map["RealmId"].(string)
	var retCanModifyTheseRealms []interface{} = response1Map["CanModifyTheseRealms"].([]interface{})
	testContext.AssertThat(retId != "", "Empty return Id")
	testContext.AssertThat(retUserId != "", "Empty return UserId")
	testContext.AssertThat(retUserName != "", "Empty return UserName")
	testContext.AssertThat(retRealmId != "", "Empty return RealmId")
	testContext.AssertThat(retCanModifyTheseRealms != nil, "No realms returned")
	
	// Authenticate as the admin user that was just created.
	var resp2 *http.Response
	resp2, err = testContext.SendPost(testContext.SessionId,
		"authenticate",
		[]string{"UserId", "Password"},
		[]string{adminUserId, adminPassword})
	if err != nil { fmt.Println(err.Error()); return "", "", nil }
	
	defer resp2.Body.Close()
	var response2Map map[string]interface{}
	response2Map, err = rest.ParseResponseBodyToMap(resp2.Body)
	if err != nil { fmt.Println(err.Error()); return "", "", nil }
	rest.PrintMap(response2Map)
	var ret2SessionId string = response2Map["UniqueSessionId"].(string)
	var ret2UserId string = response2Map["AuthenticatedUserid"].(string)
	testContext.AssertThat(ret2SessionId != "", "Session id is empty string")
	testContext.AssertThat(ret2UserId == adminUserId, "Returned user id '" + ret2UserId + "' does not match user id")

	testContext.Verify200Response(resp2)	
	testContext.SessionId = ret2SessionId
	
	// Now retrieve the description of the realm that we just created.
	var resp3 *http.Response
	resp3, err = testContext.SendPost(testContext.SessionId,
		"getRealmDesc",
		[]string{"RealmId"},
		[]string{retRealmId})
	
		// Returns RealmDesc, which contains:
		// Id
		// Name
		// OrgFullName
	
	if err != nil { fmt.Println(err.Error()); return "", "", nil }

	defer resp3.Body.Close()

	testContext.Verify200Response(resp3)
	
	var response3Map map[string]interface{}
	response3Map, err = rest.ParseResponseBodyToMap(resp3.Body)
	if err != nil { fmt.Println(err.Error()); return "", "", nil }
	var ret3Id string = response3Map["Id"].(string)
	var ret3Name string = response3Map["RealmName"].(string)
	var ret3OrgFullName string = response3Map["OrgFullName"].(string)
	rest.PrintMap(response3Map)
	testContext.AssertThat(ret3Id != "", "Empty return Id")
	testContext.AssertThat(ret3Name != "", "Empty return Name")
	testContext.AssertThat(ret3OrgFullName != "", "Empty return Org Full Name")
	
	return ret3Id, retId, retCanModifyTheseRealms
}

/*******************************************************************************
 * Returns the permissions that resulted.
 */
func (testContext *TestContext) TrySetPermission(partyId, resourceId string,
	permissions []bool) []bool {

	testContext.StartTest("TrySetPermission")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"setPermission",
		[]string{"PartyId", "ResourceId", "CanCreateIn", "CanRead", "CanWrite", "CanExecute", "CanDelete"},
		[]string{partyId, resourceId, BoolToString(permissions[0]),
			BoolToString(permissions[1]), BoolToString(permissions[2]),
			BoolToString(permissions[3]), BoolToString(permissions[4])})
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if err != nil { fmt.Println(err.Error()); return nil }
	rest.PrintMap(responseMap)

	var retACLEntryId string = responseMap["ACLEntryId"].(string)
	var retPartyId string = responseMap["PartyId"].(string)
	var retResourceId string = responseMap["ResourceId"].(string)
	var retMask []bool = make([]bool, 5)
	retMask[0] = responseMap["CanCreateIn"].(bool)
	retMask[1] = responseMap["CanRead"].(bool)
	retMask[2] = responseMap["CanWrite"].(bool)
	retMask[3] = responseMap["CanExecute"].(bool)
	retMask[4] = responseMap["CanDelete"].(bool)
	testContext.AssertThat(retACLEntryId != "", "Empty return retACLEntryId")
	testContext.AssertThat(retPartyId != "", "Empty return retPartyId")
	testContext.AssertThat(retResourceId != "", "Empty return retResourceId")
	
	return retMask
}

/*******************************************************************************
 * Returns the permissions that resulted.
 */
func (testContext *TestContext) TryAddPermission(partyId, resourceId string,
	permissions []bool) []bool {

	testContext.StartTest("TryAddPermission")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"addPermission",
		[]string{"PartyId", "ResourceId", "CanCreateIn", "CanRead", "CanWrite", "CanExecute", "CanDelete"},
		[]string{partyId, resourceId, BoolToString(permissions[0]),
			BoolToString(permissions[1]), BoolToString(permissions[2]),
			BoolToString(permissions[3]), BoolToString(permissions[4])})
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if err != nil { fmt.Println(err.Error()); return nil }
	rest.PrintMap(responseMap)

	var retACLEntryId string = responseMap["ACLEntryId"].(string)
	var retPartyId string = responseMap["PartyId"].(string)
	var retResourceId string = responseMap["ResourceId"].(string)
	var retMask []bool = make([]bool, 5)
	retMask[0] = responseMap["CanCreateIn"].(bool)
	retMask[1] = responseMap["CanRead"].(bool)
	retMask[2] = responseMap["CanWrite"].(bool)
	retMask[3] = responseMap["CanExecute"].(bool)
	retMask[4] = responseMap["CanDelete"].(bool)
	testContext.AssertThat(retACLEntryId != "", "Empty return retACLEntryId")
	testContext.AssertThat(retPartyId != "", "Empty return retPartyId")
	testContext.AssertThat(retResourceId != "", "Empty return retResourceId")
	
	return retMask
}

/*******************************************************************************
 * Return an array of string representing the values for the permission mask.
 */
func (testContext *TestContext) TryGetPermission(partyId, resourceId string) []bool {

	testContext.StartTest("TryGetPermission")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"getPermission",
		[]string{"PartyId", "ResourceId"},
		[]string{partyId, resourceId})
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if ! testContext.AssertErrIsNil(err, "while parsing response body to map") { return nil }
	rest.PrintMap(responseMap)
	
	//var retACLEntryId string = responseMap["ACLEntryId"].(string)
	var retPartyId string = responseMap["PartyId"].(string)
	var retResourceId string = responseMap["ResourceId"].(string)
	var retCreate bool = responseMap["CanCreateIn"].(bool)
	var retRead bool = responseMap["CanRead"].(bool)
	var retWrite bool = responseMap["CanWrite"].(bool)
	var retExecute bool = responseMap["CanExecute"].(bool)
	var retDelete bool = responseMap["CanDelete"].(bool)
	//testContext.AssertThat(retACLEntryId != "", "Empty return retACLEntryId")
	testContext.AssertThat(retPartyId != "", "Empty return retPartyId")
	testContext.AssertThat(retResourceId != "", "Empty return retResourceId")
	
	return []bool{retCreate, retRead, retWrite, retExecute, retDelete}
}

/*******************************************************************************
 * Return an array of the names of the available providers.
 */
func (testContext *TestContext) TryGetScanProviders() {
	testContext.StartTest("TryGetScanProviders")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"getScanProviders",
		[]string{},
		[]string{})
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMaps []map[string]interface{}
	responseMaps, err = rest.ParseResponseBodyToMaps(resp.Body)
	if err != nil { fmt.Println(err.Error()); return }
	var result []string = make([]string, 0)
	for _, responseMap := range responseMaps {
		rest.PrintMap(responseMap)
		var retProviderName string = responseMap["ProviderName"].(string)
		var retParameters []interface{} = responseMap["Parameters"].([]interface{})
		testContext.AssertThat(retProviderName != "", "Returned ProviderName is empty string")
		testContext.AssertThat(retParameters != nil, "Returned Parameters is nil")
		result = append(result, retProviderName)
	}
}

/*******************************************************************************
 * Returns the Id of the ScanConfig that gets created.
 */
func (testContext *TestContext) TryDefineScanConfig(name, desc, repoId, providerName,
	successExpr, successGraphicFilePath string, providerParamNames []string,
	providerParamValues []string) string {

	testContext.StartTest("TryDefineScanConfig")
	
	var paramNames []string = []string{"Name", "Description", "RepoId", "ProviderName"}
	var paramValues []string = []string{name, desc, repoId, providerName}
	paramNames = append(paramNames, providerParamNames...)
	paramValues = append(paramValues, providerParamValues...)
	
	fmt.Println("Param names:")
	for _, n := range paramNames { fmt.Println("\t" + n) }
	fmt.Println("Param values:")
	for _, v := range paramValues { fmt.Println("\t" + v) }
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendFilePost(testContext.SessionId,
		"defineScanConfig", paramNames, paramValues, successGraphicFilePath)
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if err != nil { fmt.Println(err.Error()); return "" }
	rest.PrintMap(responseMap)
	
	var retId string = responseMap["Id"].(string)
	var obj interface{} = responseMap["ProviderName"]
	testContext.AssertThat(obj != nil, "No ProviderName returned")
	var retProvName string = obj.(string)
	testContext.AssertThat(retId != "", "Returned Id is empty")
	testContext.AssertThat(retProvName != "", "Returned ProviderName is empty")
	// ParamValueDescs []*ParameterValueDesc
	var retParamValueDescs []interface{} = responseMap["ParameterValueDescs"].([]interface{})
	for _, desc := range retParamValueDescs {
		descMap, isType := desc.(map[string]interface{})
		if ! testContext.AssertThat(isType, "param value is not a map[string]interface{}") { continue }
		var retParamName string
		retParamName, isType = descMap["Name"].(string)
		if testContext.AssertThat(isType, "ParameterValueDesc field 'Name' is not a string") {
			testContext.AssertThat(retParamName != "", "ParameterValueDesc missing Name field")
		}
		var retParamVal string
		retParamVal, isType = descMap["Value"].(string)
		if testContext.AssertThat(isType, "ParameterValueDesc field 'Value' is not a string") {
			testContext.AssertThat(retParamVal != "", "ParameterValueDesc missing Value field")
		}
	}
	
	return retId
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryUpdateScanConfig(scanConfigId, name, desc, providerName,
	successExpr, successGraphicFilePath string, providerParamNames []string,
	providerParamValues []string) bool {
	
	testContext.StartTest("TryUpdateScanConfig")
	
	var paramNames []string = []string{"ScanConfigId", "Name", "Description", "ProviderName"}
	var paramValues []string = []string{scanConfigId, name, desc, providerName}
	paramNames = append(paramNames, providerParamNames...)
	paramValues = append(paramValues, providerParamValues...)
	
	fmt.Println("Param names:")
	for _, n := range paramNames { fmt.Println("\t" + n) }
	fmt.Println("Param values:")
	for _, v := range paramValues { fmt.Println("\t" + v) }
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendFilePost(testContext.SessionId,
		"updateScanConfig", paramNames, paramValues, successGraphicFilePath)
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if ! testContext.AssertErrIsNil(err, "") { return false }
	rest.PrintMap(responseMap)
	
	// Returns ScanConfigDesc
	var retId string
	var retProviderName string
	var retSuccessExpression string
	var retFlagId string
	var retParameterValueDescs []map[string]interface{}
	
	var isType bool
	
	retId, isType = responseMap["Id"].(string)
	if testContext.AssertThat(isType, "Id") {
		testContext.AssertThat(retId != "", "Returned Id is empty")
	}
	
	retProviderName, isType = responseMap["ProviderName"].(string)
	if testContext.AssertThat(isType, "ProviderName") {
		testContext.AssertThat(retProviderName != "", "Returned ProviderName is empty")
	}
	
	retSuccessExpression, isType = responseMap["SuccessExpression"].(string)
	if testContext.AssertThat(isType, "SuccessExpression") {
		testContext.AssertThat(retSuccessExpression != "", "Returned SuccessExpression is empty")
	}
	
	retFlagId, isType = responseMap["FlagId"].(string)
	if testContext.AssertThat(isType, "FlagId") {
		testContext.AssertThat(retFlagId != "", "Returned FlagId is empty")
	}
	
	retParameterValueDescs, isType = responseMap["ParameterValueDescs"].([]map[string]interface{})
	if testContext.AssertThat(isType, "ParameterValueDescs") {
		if testContext.AssertThat(len(retParameterValueDescs) == len(providerParamNames),
			"Wrong number of parameter descriptions returned") {
			for i, _ := range providerParamNames {
				testContext.AssertThat(providerParamNames[i] == retParameterValueDescs[i]["Name"],
					fmt.Sprintf("Parameter name %d mismatch", i))
				testContext.AssertThat(providerParamValues[i] == retParameterValueDescs[i]["StringValue"],
					fmt.Sprintf("Parameter value %d mismatch", i))
			}
		}
	}
	
	return testContext.CurrentTestPassed
}

/*******************************************************************************
 * Returns output message.
 */
func (testContext *TestContext) TryScanImage(scriptId, imageObjId string) string {
	testContext.StartTest("TryScanImage")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"scanImage",
		[]string{"ScanConfigId", "ImageObjId"},
		[]string{scriptId, imageObjId})
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if err != nil { fmt.Println(err.Error()); return "" }
	rest.PrintMap(responseMap)
	
	var retId string = responseMap["Id"].(string)
	var retWhen string = responseMap["When"].(string)
	var retUserId string = responseMap["UserId"].(string)
	var retScanConfigId string = responseMap["ScanConfigId"].(string)
	var retScore string = responseMap["Score"].(string)
	
	testContext.AssertThat(retId != "", "Returned Id is empty")
	testContext.AssertThat(retWhen != "", "Returned When is empty")
	testContext.AssertThat(retUserId != "", "Returned UserId is empty")
	testContext.AssertThat(retScanConfigId != "", "Returned ScanConfigId is empty")
	testContext.AssertThat(retScore != "", "Returned Score is empty")
	
	return retScore
}

/*******************************************************************************
 * Return the object Id of the current authenticated user.
 */
func (testContext *TestContext) TryGetMyDesc(expectSuccess bool) (string, []interface{}) {
	testContext.StartTest("TryGetMyDesc")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"getMyDesc",
		[]string{},
		[]string{})
	
	defer resp.Body.Close()

	if expectSuccess {
		if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	} else {
		if resp.StatusCode == 200 {
			testContext.FailTest()
		} else {
			testContext.PassTest()
		}
		return "", nil
	}
	
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if err != nil { fmt.Println(err.Error()); return "", nil }
	rest.PrintMap(responseMap)
	var retId string = responseMap["Id"].(string)
	var retUserId string = responseMap["UserId"].(string)
	var retUserName string = responseMap["UserName"].(string)
	var retRealmId string = responseMap["RealmId"].(string)
	var retCanModifyTheseRealms []interface{} = responseMap["CanModifyTheseRealms"].([]interface{})

	testContext.AssertThat(retId != "", "Returned Id is empty string")
	testContext.AssertThat(retUserId != "", "Returned UserId is empty string")
	testContext.AssertThat(retUserName != "", "Returned UserName is empty string")
	testContext.AssertThat(retRealmId != "", "Returned RealmId is empty string")
	testContext.AssertThat(retCanModifyTheseRealms != nil, "No realms returned")
	
	return retId, retCanModifyTheseRealms
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryGetMyGroups() []string {
	testContext.StartTest("TryGetMyGroups")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"getMyGroups",
		[]string{},
		[]string{})
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMaps []map[string]interface{}
	responseMaps, err = rest.ParseResponseBodyToMaps(resp.Body)
	if err != nil { fmt.Println(err.Error()); return nil }
	var result []string = make([]string, 0)
	for _, responseMap := range responseMaps {
		rest.PrintMap(responseMap)
		var retGroupId string = responseMap["GroupId"].(string)
		var retRealmId string = responseMap["RealmId"].(string)
		var retName string = responseMap["GroupName"].(string)
		var retCreationDate string = responseMap["CreationDate"].(string)
		var retDescription string = responseMap["Description"].(string)
		testContext.AssertThat(retGroupId != "", "Returned GroupId is empty string")
		testContext.AssertThat(retRealmId != "", "Empty returned RealmId")
		testContext.AssertThat(retName != "", "Empty returned Name")
		testContext.AssertThat(retCreationDate != "", "Empty CreationDate returned")
		testContext.AssertThat(retDescription != "", "Empty returned Description")
		
		result = append(result, retGroupId)
	}
	return result
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryGetMyRealms() []string {
	testContext.StartTest("TryGetMyRealms")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"getMyRealms",
		[]string{},
		[]string{})
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMaps []map[string]interface{}
	responseMaps, err = rest.ParseResponseBodyToMaps(resp.Body)
	if err != nil { fmt.Println(err.Error()); return nil }
	var result []string = make([]string, 0)
	for _, responseMap := range responseMaps {
		rest.PrintMap(responseMap)
		var retId string = responseMap["Id"].(string)
		var retName string = responseMap["RealmName"].(string)
	
		testContext.AssertThat(retId != "", "Returned Id is empty string")
		testContext.AssertThat(retName != "", "Empty returned Name")
		
		result = append(result, retId)
	}
	return result
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryGetMyRepos() []string {
	testContext.StartTest("TryGetMyRepos")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"getMyRepos",
		[]string{},
		[]string{})
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMaps []map[string]interface{}
	responseMaps, err = rest.ParseResponseBodyToMaps(resp.Body)
	if err != nil { fmt.Println(err.Error()); return nil }
	var result []string = make([]string, 0)
	for _, responseMap := range responseMaps {
		rest.PrintMap(responseMap)
		var retId string = responseMap["Id"].(string)
		var retRealmId string = responseMap["RealmId"].(string)
		var retName string = responseMap["RepoName"].(string)
	
		testContext.AssertThat(retId != "", "Returned Id is empty string")
		testContext.AssertThat(retRealmId != "", "Returned realm Id is empty string")
		testContext.AssertThat(retName != "", "Empty returned Name")
		
		result = append(result, retId)
	}
	return result
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryReplaceDockerfile(dockerfileId, dockerfilePath,
	desc string) {

	testContext.StartTest("TryReplaceDockerfile")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendFilePost(testContext.SessionId,
		"replaceDockerfile",
		[]string{"DockerfileId", "Description"},
		[]string{dockerfileId, desc},
		dockerfilePath)
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if ! testContext.AssertErrIsNil(err, "") { return }
	var retStatus string = responseMap["Status"].(string)
	var retMessage string = responseMap["Message"].(string)
	rest.PrintMap(responseMap)
	
	testContext.AssertThat(retStatus == "200", "Returned Status is empty")
	testContext.AssertThat(retMessage != "", "Returned Message is empty")
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryDownloadImage(imageId, filename string) {

	testContext.StartTest("TryDownloadImage")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"downloadImage",
		[]string{"ImageObjId"},
		[]string{imageId})
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	// Check that the server actual sent compressed data
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
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryRemGroupUser(groupId, userObjId string) bool {

	testContext.StartTest("TryRemGroupUser")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"remGroupUser",
		[]string{"GroupId", "UserObjId"},
		[]string{groupId, userObjId})
	
	defer resp.Body.Close()
	if ! testContext.AssertErrIsNil(err, "") { return false }

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	return testContext.CurrentTestPassed
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryReenableUser(userObjId string) bool {
	testContext.StartTest("TryReenableUser")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"reenableUser",
		[]string{"UserObjId"},
		[]string{userObjId})
	
	defer resp.Body.Close()
	if ! testContext.AssertErrIsNil(err, "") { return false }

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	return testContext.CurrentTestPassed
}
	
/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryRemRealmUser(realmId, userObjId string) bool {
	testContext.StartTest("TryRemRealmUser")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"remRealmUser",
		[]string{"RealmId", "UserObjId"},
		[]string{realmId, userObjId})
	
	defer resp.Body.Close()
	if ! testContext.AssertErrIsNil(err, "") { return false }

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	return testContext.CurrentTestPassed
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryDeactivateRealm(realmId string) bool {
	testContext.StartTest("TryDeactivateRealm")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"deactivateRealm",
		[]string{"RealmId"},
		[]string{realmId})
	
	defer resp.Body.Close()
	if ! testContext.AssertErrIsNil(err, "") { return false }

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	return testContext.CurrentTestPassed
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryDeleteRepo() {
	testContext.StartTest("TryDeleteRepo")
	testContext.FailTest()
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryRemPermission(partyId, resourceId string) bool {
	
	testContext.StartTest("TryRemPermission")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"remPermission",
		[]string{"PartyId", "ResourceId"},
		[]string{partyId, resourceId})
	
	defer resp.Body.Close()
	if ! testContext.AssertErrIsNil(err, "") { return false }

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	return testContext.CurrentTestPassed
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryGetUserEvents(userId string) []string {
	
	testContext.StartTest("TryGetUserEvents")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"getUserEvents",
		[]string{"UserId"},
		[]string{userId})
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMaps []map[string]interface{}
	responseMaps, err = rest.ParseResponseBodyToMaps(resp.Body)
	if err != nil { fmt.Println(err.Error()); return nil }
	var result []string = make([]string, 0)
	for _, responseMap := range responseMaps {
		rest.PrintMap(responseMap)
		var retId string = responseMap["Id"].(string)
		testContext.AssertThat(retId != "", "Returned Id is empty string")
		result = append(result, retId)
	}
	return result
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryGetDockerImageEvents(imageId string) []string {
	
	testContext.StartTest("TryGetDockerImageEvents")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"getDockerImageEvents",
		[]string{"ImageId"},
		[]string{imageId})
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMaps []map[string]interface{}
	responseMaps, err = rest.ParseResponseBodyToMaps(resp.Body)
	if err != nil { fmt.Println(err.Error()); return nil }
	var result []string = make([]string, 0)
	for _, responseMap := range responseMaps {
		rest.PrintMap(responseMap)
		var retId string = responseMap["Id"].(string)
		testContext.AssertThat(retId != "", "Returned Id is empty string")
		result = append(result, retId)
	}
	return result
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryGetDockerImageStatus(imageObjId string) map[string]interface{} {
	
	testContext.StartTest("TryGetImageStatus")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"getDockerImageStatus",
		[]string{"ImageObjId"},
		[]string{imageObjId},
		)
	if ! testContext.AssertErrIsNil(err, "") { return nil }
	
	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if ! testContext.AssertErrIsNil(err, "") { return nil }

	//EventId string
	//When time.Time
	//UserObjId string
	//EventDescBase
	//ScanConfigId string
	//ProviderName string
    //ParameterValueDescs []*ParameterValueDesc
	//Score string

	return responseMap
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryGetDockerfileEvents(dockerfileId string) []string {

	testContext.StartTest("TryGetDockerfileEvents")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"getDockerfileEvents",
		[]string{"DockerfileId"},
		[]string{dockerfileId})
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMaps []map[string]interface{}
	responseMaps, err = rest.ParseResponseBodyToMaps(resp.Body)
	if err != nil { fmt.Println(err.Error()); return nil }
	var result []string = make([]string, 0)
	for _, responseMap := range responseMaps {
		rest.PrintMap(responseMap)
		var retId string = responseMap["Id"].(string)
		testContext.AssertThat(retId != "", "Returned Id is empty string")
		result = append(result, retId)
	}
	return result
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryDefineFlag(repoId, flagName, desc,
	imageFilePath string) map[string]interface{} {
	
	testContext.StartTest("TryDefineFlag")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendFilePost(testContext.SessionId,
		"defineFlag",
		[]string{"RepoId", "Name", "Description"},
		[]string{repoId, flagName, desc},
		imageFilePath)
	if ! testContext.AssertErrIsNil(err, "") { return nil}
	
	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	testContext.AssertErrIsNil(err, "")

	return responseMap
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryGetScanConfigDesc(scanConfigId string,
	expectToFindIt bool) map[string]interface{} {
	
	testContext.StartTest("TryGetScanConfigDesc")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"getScanConfigDesc",
		[]string{"ScanConfigId"},
		[]string{scanConfigId})
	if ! testContext.AssertErrIsNil(err, "") { return nil }
	
	if expectToFindIt {
		if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	} else {
		if resp.StatusCode == 200 {
			testContext.FailTest()
		} else {
			testContext.PassTest()
		}	
		return nil
	}
	
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	testContext.AssertErrIsNil(err, "")

	var retScanConfigId string = ""
	var scanConfigIdIsType bool
	if retScanConfigId, scanConfigIdIsType = responseMap["Id"].(string); (! scanConfigIdIsType) || (retScanConfigId == "") { testContext.FailTest() }
	if retProviderName, isType := responseMap["ProviderName"].(string); (! isType) || (retProviderName == "") { testContext.FailTest() }
	if retSuccessExpression, isType := responseMap["SuccessExpression"].(string); (! isType) || (retSuccessExpression == "") { testContext.FailTest() }
	if retFlagId, isType := responseMap["FlagId"].(string); (! isType) || (retFlagId == "") { testContext.FailTest() }
	if retParameterValueDescs, isType := responseMap["ParameterValueDescs"].(string); (! isType) || (retParameterValueDescs == "") { testContext.FailTest() }
	
	return responseMap
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryChangePassword(userId, oldPswd, newPswd string) bool {
	
	
	testContext.StartTest("TryChangePassword")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"changePassword",
		[]string{"UserId", "OldPassword", "NewPassword"},
		[]string{userId, oldPswd, newPswd})
	if ! testContext.AssertErrIsNil(err, "") { return false }
	
	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	//var responseMap map[string]interface{}
	_, err = rest.ParseResponseBodyToMap(resp.Body)
	if ! testContext.AssertErrIsNil(err, "") { return false }

	return testContext.CurrentTestPassed
}

/*******************************************************************************
 * Returns the name of the flag.
 */
func (testContext *TestContext) TryGetFlagDesc(flagId string, expectToFindIt bool) string {
	
	testContext.StartTest("TryGetFlagDesc")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"getFlagDesc",
		[]string{"FlagId"},
		[]string{flagId})
	if ! testContext.AssertErrIsNil(err, "") { return ""}
	
	if expectToFindIt {
		if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	} else {
		if resp.StatusCode == 200 {
			testContext.FailTest()
		} else {
			testContext.PassTest()
		}	
		return ""
	}

	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if ! testContext.AssertErrIsNil(err, "") { return "" }

	var retNameIsType bool
	var retName string = ""
	if retFlagId, isType := responseMap["FlagId"].(string); (! isType) || (retFlagId == "") { testContext.FailTest() }
	if retRepoId, isType := responseMap["RepoId"].(string); (! isType) || (retRepoId == "") { testContext.FailTest() }
	if retName, retNameIsType = responseMap["Name"].(string); (! retNameIsType) || (retName == "") { testContext.FailTest() }
	if retImageURL, isType := responseMap["ImageURL"].(string); (! isType) || (retImageURL == "") { testContext.FailTest() }

	return retName
}

/*******************************************************************************
 * Returns the size of the file that was downloaded.
 */
func (testContext *TestContext) TryGetFlagImage(flagId string, filename string) int64 {
	
	testContext.StartTest("TryGetFlagImage")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"getFlagImage",
		[]string{"FlagId"},
		[]string{flagId})
	if ! testContext.AssertErrIsNil(err, "") { return 0 }
	
	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	//var responseMap map[string]interface{}
	_, err = rest.ParseResponseBodyToMap(resp.Body)
	if ! testContext.AssertErrIsNil(err, "") { return 0 }

	var reader io.ReadCloser = resp.Body
	var file *os.File
	file, err = os.Create(filename)
	testContext.AssertErrIsNil(err, "")
	_, err = io.Copy(file, reader)
	testContext.AssertErrIsNil(err, "")
	var fileInfo os.FileInfo
	fileInfo, err = file.Stat()
	if ! testContext.AssertErrIsNil(err, "") { return 0 }
	testContext.AssertThat(fileInfo.Size() > 0, "File has zero size")
	
	return fileInfo.Size()
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryGetMyScanConfigs() []string {
	testContext.StartTest("TryGetMyScanConfigs")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"getMyScanConfigs",
		[]string{},
		[]string{})
	if ! testContext.AssertErrIsNil(err, "") { return nil }
	
	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMaps []map[string]interface{}
	responseMaps, err = rest.ParseResponseBodyToMaps(resp.Body)
	if err != nil { fmt.Println(err.Error()); return nil }
	var retConfigIds []string = make([]string, 0)
	for _, responseMap := range responseMaps {
		rest.PrintMap(responseMap)
		
		if retId, isType := responseMap["Id"].(string); (! isType) || (retId == "") {
			testContext.FailTest()
		} else {
			retConfigIds = append(retConfigIds, retId)
		}
		if retProviderName, isType := responseMap["ProviderName"].(string); (! isType) || (retProviderName == "") { testContext.FailTest() }
		if retSuccessExpression, isType := responseMap["SuccessExpression"].(string); (! isType) || (retSuccessExpression == "") { testContext.FailTest() }
		if retFlagId, isType := responseMap["FlagId"].(string); (! isType) || (retFlagId == "") { testContext.FailTest() }
		if retParameterValueDescs, isType := responseMap["ParameterValueDescs"].(string); (! isType) || (retParameterValueDescs == "") { testContext.FailTest() }
	}

	return retConfigIds
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryGetScanConfigDescByName(repoId, scanConfigName string) string {
	testContext.StartTest("TryGetScanConfigDescByName")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"getScanConfigDescByName",
		[]string{"RepoId", "ScanConfigName"},
		[]string{repoId, scanConfigName})
	if ! testContext.AssertErrIsNil(err, "") { return "" }
	
	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if ! testContext.AssertErrIsNil(err, "") { return "" }

	var retScanConfigId string = ""
	var scanConfigIdIsType bool
	if retScanConfigId, scanConfigIdIsType = responseMap["Id"].(string); (! scanConfigIdIsType) || (retScanConfigId == "") { testContext.FailTest() }
	if retProviderName, isType := responseMap["ProviderName"].(string); (! isType) || (retProviderName == "") { testContext.FailTest() }
	if retSuccessExpression, isType := responseMap["SuccessExpression"].(string); (! isType) || (retSuccessExpression == "") { testContext.FailTest() }
	if retFlagId, isType := responseMap["FlagId"].(string); (! isType) || (retFlagId == "") { testContext.FailTest() }
	if retParameterValueDescs, isType := responseMap["ParameterValueDescs"].(string); (! isType) || (retParameterValueDescs == "") { testContext.FailTest() }
	return retScanConfigId
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryRemScanConfig(scanConfigId string,
	expectSuccess bool) bool {
	testContext.StartTest("TryRemScanConfig")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"remScanConfig",
		[]string{"ScanConfigId"},
		[]string{scanConfigId})
	if testContext.AssertErrIsNil(err, "") { return false }
	
	if expectSuccess {
		if ! testContext.Verify200Response(resp) {
			testContext.FailTest()
			return false
		}
	} else {
		if resp.StatusCode == 200 {
			testContext.FailTest()
			return false
		} else {
			testContext.PassTest()
			return true
		}	
	}
		
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if ! testContext.AssertErrIsNil(err, "") { return false }

	if retStatus, isType := responseMap["Status"].(string); (! isType) || (retStatus == "") { testContext.FailTest() }
	if retMessage, isType := responseMap["Message"].(string); (! isType) || (retMessage == "") { testContext.FailTest() }

	return testContext.CurrentTestPassed
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryGetMyFlags() []string {
	testContext.StartTest("TryGetMyFlags")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"getMyFlags",
		[]string{},
		[]string{})
	if testContext.AssertErrIsNil(err, "") { return nil }
	
	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMaps []map[string]interface{}
	responseMaps, err = rest.ParseResponseBodyToMaps(resp.Body)
	if err != nil { fmt.Println(err.Error()); return nil }
	var retFlagIds []string = make([]string, 0)
	for _, responseMap := range responseMaps {
		rest.PrintMap(responseMap)
		
		if retFlagId, isType := responseMap["FlagId"].(string); (! isType) || (retFlagId == "") {
			testContext.FailTest()
		} else {
			retFlagIds = append(retFlagIds, retFlagId)
		}
		if retRepoId, isType := responseMap["RepoId"].(string); (! isType) || (retRepoId == "") { testContext.FailTest() }
		if retName, isType := responseMap["Name"].(string); (! isType) || (retName == "") { testContext.FailTest() }
		if retImageURL, isType := responseMap["ImageURL"].(string); (! isType) || (retImageURL == "") { testContext.FailTest() }
	}

	return retFlagIds
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryGetFlagDescByName(repoId, flagName string) string {
	testContext.StartTest("TryGetFlagDescByName")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"getFlagDescByName",
		[]string{"RepoId", "FlagName"},
		[]string{repoId, flagName})
	if err != nil { fmt.Println(err.Error()); return "" }
	
	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if ! testContext.AssertErrIsNil(err, "") { return "" }

	var retFlagId string = ""
	var flagIdIsType bool
	if retFlagId, flagIdIsType = responseMap["FlagId"].(string); (! flagIdIsType) || (retFlagId == "") { testContext.FailTest() }
	if retRepoId, isType := responseMap["RepoId"].(string); (! isType) || (retRepoId == "") { testContext.FailTest() }
	if retName, isType := responseMap["Name"].(string); (! isType) || (retName == "") { testContext.FailTest() }
	if retImageURL, isType := responseMap["ImageURL"].(string); (! isType) || (retImageURL == "") { testContext.FailTest() }
	return retFlagId
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryRemFlag(flagId string) bool {
	testContext.StartTest("TryRemFlag")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"remFlag",
		[]string{"FlagId"},
		[]string{flagId})
	if ! testContext.AssertErrIsNil(err, "") { return false }
	
	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if ! testContext.AssertErrIsNil(err, "") { return false }

	if retStatus, isType := responseMap["Status"].(string); (! isType) || (retStatus == "") { testContext.FailTest() }
	if retMessage, isType := responseMap["Message"].(string); (! isType) || (retMessage == "") { testContext.FailTest() }

	return testContext.CurrentTestPassed
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryRemDockerImage(imageId string) bool {
	testContext.StartTest("TryRemDockerImage")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"remDockerImage",
		[]string{"ImageId"},
		[]string{imageId})
	if ! testContext.AssertErrIsNil(err, "") { return false }
	
	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if ! testContext.AssertErrIsNil(err, "") { return false }

	if retStatus, isType := responseMap["Status"].(string); (! isType) || (retStatus == "") { testContext.FailTest() }
	if retMessage, isType := responseMap["Message"].(string); (! isType) || (retMessage == "") { testContext.FailTest() }

	return testContext.CurrentTestPassed
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryClearAll() {
	testContext.StartTest("TryClearAll")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendGet("",
		"clearAll",
		[]string{},
		[]string{},
		)
	if ! testContext.AssertErrIsNil(err, "") { return }
	
	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
}
