package utils

import (
	"fmt"
	"net/http"
	"os"
	"io"
	
	"testsafeharbor/rest"
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
		[]string{"Log", "GroupId"},
		[]string{testContext.TestDemarcation(), groupId})
	
	defer resp.Body.Close()
	
	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if ! testContext.AssertErrIsNil(err, "at ParseResponseBodyToMap") { return }
	
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
	testContext.PassTestIfNoFailures()
}
	
/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryGetRepoDesc(repoId string) {
	
	testContext.StartTest("TryGetRepoDesc")
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"getRepoDesc",
		[]string{"Log", "RepoId"},
		[]string{testContext.TestDemarcation(), repoId})
	
	defer resp.Body.Close()
	
	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if ! testContext.AssertErrIsNil(err, "at ParseResponseBodyToMap") { return }
	rest.PrintMap(responseMap)
	
	// Expect a RepoDesc
	var retId string = responseMap["Id"].(string)
	var retRealmId string = responseMap["RealmId"].(string)
	var retRepoName string = responseMap["RepoName"].(string)
	var retDescription string = responseMap["Description"].(string)
	var retCreationDate string = responseMap["CreationDate"].(string)
	if retDockerfileIds, isType := responseMap["DockerfileIds"].([]interface{}); (! isType) ||
		(retDockerfileIds == nil) {
		testContext.FailTest()
	}
	testContext.AssertThat(retId != "", "retId is empty")
	testContext.AssertThat(retRealmId != "", "retRealmId is empty")
	testContext.AssertThat(retRepoName != "", "retRepoName is empty")
	testContext.AssertThat(retDescription != "", "retDescription is empty")
	testContext.AssertThat(retCreationDate != "", "retCreationDate is empty")
	testContext.PassTestIfNoFailures()
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
		[]string{"Log", "DockerImageId"},
		[]string{testContext.TestDemarcation(), dockerImageId})
	
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
			testContext.PassTestIfNoFailures()
			return nil
		}	
	}
	
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if ! testContext.AssertErrIsNil(err, "at ParseResponseBodyToMap") { return nil }
	
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
	
	testContext.PassTestIfNoFailures()
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
		[]string{"Log", "DockerfileId"},
		[]string{testContext.TestDemarcation(), dockerfileId})
	
	defer resp.Body.Close()
	
	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if ! testContext.AssertErrIsNil(err, "at ParseResponseBodyToMap") { return }
	
	// Expect a DockerfileDesc
	var retId string = responseMap["Id"].(string)
	var retRepoId string = responseMap["RepoId"].(string)
	var retDescription string = responseMap["Description"].(string)
	var retDockerfileName string = responseMap["DockerfileName"].(string)
	
	testContext.AssertThat(retId != "", "retId is empty")
	testContext.AssertThat(retRepoId != "", "retRepoId is empty")
	testContext.AssertThat(retDescription != "", "retDescription is empty")
	testContext.AssertThat(retDockerfileName != "", "retDockerfileName is empty")
	testContext.PassTestIfNoFailures()
}

/*******************************************************************************
 * Verify that we can create a new realm.
 */
func (testContext *TestContext) TryCreateRealm(realmName, orgFullName,
	desc string) string {
	
	testContext.StartTest("TryCreateRealm")
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"createRealm",
		[]string{"Log", "RealmName", "OrgFullName", "Description"},
		[]string{testContext.TestDemarcation(), realmName, orgFullName, desc})
	
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
	
	testContext.PassTestIfNoFailures()
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
		[]string{"Log", "UserId", "UserName", "EmailAddress", "Password", "RealmId"},
		[]string{testContext.TestDemarcation(), userId, userName, email, pswd, realmId})
	
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
	
	testContext.PassTestIfNoFailures()
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
		[]string{"Log", "UserId", "Password"},
		[]string{testContext.TestDemarcation(), userId, pswd})
	
	defer resp.Body.Close()

	if expectSuccess {
		if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	} else {
		if resp.StatusCode == 200 {
			testContext.FailTest()
			return "", false
		} else {
			testContext.PassTestIfNoFailures()
			return "", false
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
	testContext.PassTestIfNoFailures()
	testContext.SessionId = retSessionId
	testContext.IsAdmin = retIsAdmin
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
		[]string{"Log", "UserObjId"},
		[]string{testContext.TestDemarcation(), userObjId})
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if err != nil { fmt.Println(err.Error()); return false }
	rest.PrintMap(responseMap)
	var retStatus string = responseMap["Status"].(string)
	//var retMessage string = responseMap["Message"].(string)
	if retStatus != "200" { return false }
	testContext.PassTestIfNoFailures()
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
		[]string{"Log", "GroupId"},
		[]string{testContext.TestDemarcation(), groupId})
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if err != nil { fmt.Println(err.Error()); return false }
	rest.PrintMap(responseMap)
	var retStatus string = responseMap["Status"].(string)
	//var retMessage string = responseMap["Message"].(string)
	if retStatus != "200" { return false }
	testContext.PassTestIfNoFailures()
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
		[]string{"Log"},
		[]string{testContext.TestDemarcation()})
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if err != nil { fmt.Println(err.Error()); return false }
	rest.PrintMap(responseMap)
	var retStatus string = responseMap["Status"].(string)
	//var retMessage string = responseMap["Message"].(string)
	if retStatus != "200" { return false }
	testContext.PassTestIfNoFailures()
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
			[]string{"Log", "RealmId", "Name", "Description"},
			[]string{testContext.TestDemarcation(), realmId, name, desc})
	} else {
		fmt.Println("Using SendFilePost")
		resp, err = testContext.SendFilePost(testContext.SessionId,
			"createRepo",
			[]string{"Log", "RealmId", "Name", "Description"},
			[]string{testContext.TestDemarcation(), realmId, name, desc},
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
	
	testContext.PassTestIfNoFailures()
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
		[]string{"Log", "RepoId", "Description"},
		[]string{testContext.TestDemarcation(), repoId, desc},
		dockerfilePath)
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	// Get the DockerfileDesc that is returned.
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if err != nil { fmt.Println(err.Error()); return "" }
	var dockerfileId string = responseMap["Id"].(string)
	var dockerfileName string = responseMap["DockerfileName"].(string)
	rest.PrintMap(responseMap)
	testContext.AssertThat(dockerfileId != "", "Dockerfile Id not found in response body")
	testContext.AssertThat(dockerfileName != "", "Dockerfile Name not found in response body")
	
	testContext.PassTestIfNoFailures()
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
		[]string{"Log", "RepoId"},
		[]string{testContext.TestDemarcation(), repoId})
	
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
		
	testContext.PassTestIfNoFailures()
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
		[]string{"Log", "RepoId", "DockerfileId", "ImageName"},
		[]string{testContext.TestDemarcation(), repoId, dockerfileId, imageName})
	
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
	testContext.PassTestIfNoFailures()
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
		[]string{"Log", "RepoId", "Description", "ImageName", "SessionId"},
		[]string{testContext.TestDemarcation(), repoId, desc, imageName, testContext.SessionId},
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
	testContext.PassTestIfNoFailures()
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
		[]string{"Log", "RepoId"},
		[]string{testContext.TestDemarcation(), repoId})
	
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
	
	testContext.PassTestIfNoFailures()
	return result
}

/*******************************************************************************
 * Return the object Id of the specified user, and a list of the realms that
 * the user can modify.
 */
func (testContext *TestContext) TryGetUserDesc(userId string) map[string]interface{} {
	testContext.StartTest("TryGetUserDesc")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"getUserDesc",
		[]string{"Log", "UserId"},
		[]string{testContext.TestDemarcation(), userId})
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if ! testContext.AssertErrIsNil(err, "") { return nil }
	var retUserObjId string = responseMap["Id"].(string)
	var retUserId string = responseMap["UserId"].(string)
	var retUserName string = responseMap["UserName"].(string)
	var retCanModifyTheseRealms []interface{} = responseMap["CanModifyTheseRealms"].([]interface{})
	rest.PrintMap(responseMap)
	
	testContext.AssertThat(retUserObjId != "", "User obj Id not returned")
	testContext.AssertThat(retUserId == userId, "Returned user id, " + retUserId +
		" does not match the original user id")
	testContext.AssertThat(retUserName != "", "Returned user name is blank")
	testContext.AssertThat(retCanModifyTheseRealms != nil, "No realms returned")
	
	testContext.PassTestIfNoFailures()
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
		[]string{"Log", "RealmId", "Name", "Description", "AddMe"},
		[]string{testContext.TestDemarcation(), realmId, name, description, addMeStr})
	
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
	
	testContext.PassTestIfNoFailures()
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
		[]string{"Log", "GroupId"},
		[]string{testContext.TestDemarcation(), groupId})
	
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
	
	testContext.PassTestIfNoFailures()
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
		[]string{"Log", "GroupId", "UserObjId"},
		[]string{testContext.TestDemarcation(), groupId, userId})
	
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
	
	testContext.PassTestIfNoFailures()
	return testContext.CurrentTestPassed
}

/*******************************************************************************
 * Returns result.
 */
func (testContext *TestContext) TryMoveUserToRealm(userObjId, realmId string) bool {
	testContext.StartTest("TryMoveUserToRealm")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"moveUserToRealm",
		[]string{"Log", "UserObjId", "RealmId"},
		[]string{testContext.TestDemarcation(), userObjId, realmId})
	
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
	testContext.PassTestIfNoFailures()
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
		[]string{"Log", "RealmId"},
		[]string{testContext.TestDemarcation(), realmId})
	
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
	
	testContext.PassTestIfNoFailures()
	return result
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryGetRealmRepos(realmId string, expectSuccess bool) []string {
	testContext.StartTest("TryGetRealmRepos")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"getRealmRepos",
		[]string{"Log", "RealmId"},
		[]string{testContext.TestDemarcation(), realmId})
	
	if expectSuccess {
		if ! testContext.Verify200Response(resp) {
			testContext.FailTest()
		}
	} else {
		if resp.StatusCode == 200 {
			testContext.FailTest()
		} else {
			testContext.PassTestIfNoFailures()
		}
		return nil
	}
	
	defer resp.Body.Close()
	
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
	testContext.PassTestIfNoFailures()
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
		[]string{"Log"},
		[]string{testContext.TestDemarcation()})
	
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
	testContext.PassTestIfNoFailures()
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
		[]string{"Log"},
		[]string{testContext.TestDemarcation()})
	
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
	testContext.PassTestIfNoFailures()
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
		[]string{"Log"},
		[]string{testContext.TestDemarcation()})
	
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
	testContext.PassTestIfNoFailures()
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
		[]string{"Log", "RealmId"},
		[]string{testContext.TestDemarcation(), realmId})
	
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
	testContext.PassTestIfNoFailures()
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
		[]string{"Log", "UserId", "UserName", "EmailAddress", "Password", "RealmName", "OrgFullName"},
		[]string{testContext.TestDemarcation(), adminUserId, adminUserFullName, adminEmailAddr, adminPassword,
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
	
	testContext.PassTestIfNoFailures()
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
		[]string{"Log", "PartyId", "ResourceId", "CanCreateIn", "CanRead", "CanWrite", "CanExecute", "CanDelete"},
		[]string{testContext.TestDemarcation(), partyId, resourceId, BoolToString(permissions[0]),
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
	
	testContext.PassTestIfNoFailures()
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
		[]string{"Log", "PartyId", "ResourceId", "CanCreateIn", "CanRead", "CanWrite", "CanExecute", "CanDelete"},
		[]string{testContext.TestDemarcation(), partyId, resourceId, BoolToString(permissions[0]),
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
	
	testContext.PassTestIfNoFailures()
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
		[]string{"Log", "PartyId", "ResourceId"},
		[]string{testContext.TestDemarcation(), partyId, resourceId})
	
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
	
	testContext.PassTestIfNoFailures()
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
		[]string{"Log"},
		[]string{testContext.TestDemarcation()})
	
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
	testContext.PassTestIfNoFailures()
}

/*******************************************************************************
 * Returns the Id of the ScanConfig that gets created.
 */
func (testContext *TestContext) TryDefineScanConfig(name, desc, repoId, providerName,
	successExpr, successGraphicFilePath string, providerParamNames []string,
	providerParamValues []string) string {

	testContext.StartTest("TryDefineScanConfig")
	
	var paramNames []string = []string{"Log", "Name", "Description", "RepoId", "ProviderName"}
	var paramValues []string = []string{testContext.TestDemarcation(), name, desc, repoId, providerName}
	paramNames = append(paramNames, providerParamNames...)
	paramValues = append(paramValues, providerParamValues...)
	
	fmt.Println("Param names:")
	for _, n := range paramNames { fmt.Println("\t" + n) }
	fmt.Println("Param values:")
	for _, v := range paramValues { fmt.Println("\t" + v) }
	
	var resp *http.Response
	var err error
	if successGraphicFilePath == "" {
		resp, err = testContext.SendPost(testContext.SessionId,
			"defineScanConfig", paramNames, paramValues)
	} else {
		resp, err = testContext.SendFilePost(testContext.SessionId,
			"defineScanConfig",
			paramNames,
			paramValues,
			successGraphicFilePath)
	}
	testContext.AssertErrIsNil(err, "at the POST")
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	testContext.AssertErrIsNil(err, "at ParseResponseBodyToMap")
	rest.PrintMap(responseMap)
	
	var retId string = responseMap["Id"].(string)
	var obj interface{} = responseMap["ProviderName"]
	testContext.AssertThat(obj != nil, "No ProviderName returned")
	var retProvName string = obj.(string)
	testContext.AssertThat(retId != "", "Returned Id is empty")
	testContext.AssertThat(retProvName != "", "Returned ProviderName is empty")
	if successGraphicFilePath != "" {
		obj = responseMap["FlagId"]
		var retFlagId string
		var isType bool
		retFlagId, isType = obj.(string)
		testContext.AssertThat(isType && (retFlagId != ""), "Returned FlagId is empty")
	}
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
	
	testContext.PassTestIfNoFailures()
	return retId
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryUpdateScanConfig(scanConfigId, name, desc, providerName,
	successExpr, successGraphicFilePath string, providerParamNames []string,
	providerParamValues []string) map[string]interface{} {
	
	testContext.StartTest("TryUpdateScanConfig")
	
	var paramNames []string = []string{"Log", "ScanConfigId", "Name", "Description", "ProviderName"}
	var paramValues []string = []string{testContext.TestDemarcation(), scanConfigId, name, desc, providerName}
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
	if ! testContext.AssertErrIsNil(err, "") { return nil }
	rest.PrintMap(responseMap)
	
	// Returns ScanConfigDesc
	var retId string
	var retProviderName string
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
	
	retFlagId, isType = responseMap["FlagId"].(string)
	if testContext.AssertThat(isType, "FlagId") {
		testContext.AssertThat(retFlagId != "", "Returned FlagId is empty")
	}
	
	if len(providerParamNames) > 0 {
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
	}
	
	testContext.PassTestIfNoFailures()
	return responseMap
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
		[]string{"Log", "ScanConfigId", "ImageObjId"},
		[]string{testContext.TestDemarcation(), scriptId, imageObjId})
	
	defer resp.Body.Close()

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if err != nil { fmt.Println(err.Error()); return "" }
	rest.PrintMap(responseMap)
	
	var retId string = responseMap["Id"].(string)
	var retWhen string = responseMap["When"].(string)
	var retUserId string = responseMap["UserObjId"].(string)
	var retScanConfigId string = responseMap["ScanConfigId"].(string)
	var retScore string = responseMap["Score"].(string)
	
	testContext.AssertThat(retId != "", "Returned Id is empty")
	testContext.AssertThat(retWhen != "", "Returned When is empty")
	testContext.AssertThat(retUserId != "", "Returned UserId is empty")
	testContext.AssertThat(retScanConfigId != "", "Returned ScanConfigId is empty")
	testContext.AssertThat(retScore != "", "Returned Score is empty")
	
	testContext.PassTestIfNoFailures()
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
		[]string{"Log"},
		[]string{testContext.TestDemarcation()})
	
	defer resp.Body.Close()

	if expectSuccess {
		if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	} else {
		if resp.StatusCode == 200 {
			testContext.FailTest()
		} else {
			testContext.PassTestIfNoFailures()
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
	
	testContext.PassTestIfNoFailures()
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
		[]string{"Log"},
		[]string{testContext.TestDemarcation()})
	
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
	testContext.PassTestIfNoFailures()
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
		[]string{"Log"},
		[]string{testContext.TestDemarcation()})
	
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
	testContext.PassTestIfNoFailures()
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
		[]string{"Log"},
		[]string{testContext.TestDemarcation()})
	
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
	testContext.PassTestIfNoFailures()
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
		[]string{"Log", "DockerfileId", "Description"},
		[]string{testContext.TestDemarcation(), dockerfileId, desc},
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
	testContext.PassTestIfNoFailures()
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
		[]string{"Log", "ImageObjId"},
		[]string{testContext.TestDemarcation(), imageId})
	
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
	testContext.PassTestIfNoFailures()
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
		[]string{"Log", "GroupId", "UserObjId"},
		[]string{testContext.TestDemarcation(), groupId, userObjId})
	
	defer resp.Body.Close()
	if ! testContext.AssertErrIsNil(err, "") { return false }

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	testContext.PassTestIfNoFailures()
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
		[]string{"Log", "UserObjId"},
		[]string{testContext.TestDemarcation(), userObjId})
	
	defer resp.Body.Close()
	if ! testContext.AssertErrIsNil(err, "") { return false }

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	testContext.PassTestIfNoFailures()
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
		[]string{"Log", "RealmId", "UserObjId"},
		[]string{testContext.TestDemarcation(), realmId, userObjId})
	
	defer resp.Body.Close()
	if ! testContext.AssertErrIsNil(err, "") { return false }

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	testContext.PassTestIfNoFailures()
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
		[]string{"Log", "RealmId"},
		[]string{testContext.TestDemarcation(), realmId})
	
	defer resp.Body.Close()
	if ! testContext.AssertErrIsNil(err, "") { return false }

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	testContext.PassTestIfNoFailures()
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
		[]string{"Log", "PartyId", "ResourceId"},
		[]string{testContext.TestDemarcation(), partyId, resourceId})
	
	defer resp.Body.Close()
	if ! testContext.AssertErrIsNil(err, "") { return false }

	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	testContext.PassTestIfNoFailures()
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
		[]string{"Log", "UserId"},
		[]string{testContext.TestDemarcation(), userId})
	
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
	testContext.PassTestIfNoFailures()
	return result
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryGetDockerImageEvents(imageObjId string) []string {
	
	testContext.StartTest("TryGetDockerImageEvents")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"getDockerImageEvents",
		[]string{"Log", "ImageObjId"},
		[]string{testContext.TestDemarcation(), imageObjId})
	
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
	testContext.PassTestIfNoFailures()
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
		[]string{"Log", "ImageObjId"},
		[]string{testContext.TestDemarcation(), imageObjId},
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

	testContext.PassTestIfNoFailures()
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
		[]string{"Log", "DockerfileId"},
		[]string{testContext.TestDemarcation(), dockerfileId})
	
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
	testContext.PassTestIfNoFailures()
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
		[]string{"Log", "RepoId", "Name", "Description"},
		[]string{testContext.TestDemarcation(), repoId, flagName, desc},
		imageFilePath)
	if ! testContext.AssertErrIsNil(err, "") { return nil}
	
	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	testContext.AssertErrIsNil(err, "")

	testContext.PassTestIfNoFailures()
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
		[]string{"Log", "ScanConfigId"},
		[]string{testContext.TestDemarcation(), scanConfigId})
	if ! testContext.AssertErrIsNil(err, "") { return nil }
	
	if expectToFindIt {
		if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	} else {
		if resp.StatusCode == 200 {
			testContext.FailTest()
		} else {
			testContext.PassTestIfNoFailures()
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
	if retFlagId, isType := responseMap["FlagId"].(string); (! isType) || (retFlagId == "") { testContext.FailTest() }
	if retParameterValueDescs, isType := responseMap["ParameterValueDescs"].([]interface{}); (! isType) || (retParameterValueDescs == nil) { testContext.FailTest() }
	
	testContext.PassTestIfNoFailures()
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
		[]string{"Log", "UserId", "OldPassword", "NewPassword"},
		[]string{testContext.TestDemarcation(), userId, oldPswd, newPswd})
	if ! testContext.AssertErrIsNil(err, "") { return false }
	
	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	//var responseMap map[string]interface{}
	_, err = rest.ParseResponseBodyToMap(resp.Body)
	if ! testContext.AssertErrIsNil(err, "") { return false }

	testContext.PassTestIfNoFailures()
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
		[]string{"Log", "FlagId"},
		[]string{testContext.TestDemarcation(), flagId})
	if ! testContext.AssertErrIsNil(err, "") { return ""}
	
	if expectToFindIt {
		if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	} else {
		if resp.StatusCode == 200 {
			testContext.FailTest()
		} else {
			testContext.PassTestIfNoFailures()
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

	testContext.PassTestIfNoFailures()
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
		[]string{"Log", "FlagId"},
		[]string{testContext.TestDemarcation(), flagId})
	if ! testContext.AssertErrIsNil(err, "") { return 0 }
	
	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
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
	
	testContext.PassTestIfNoFailures()
	return fileInfo.Size()
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryGetMyScanConfigs() ([]map[string]interface{}, []string) {
	testContext.StartTest("TryGetMyScanConfigs")
	
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"getMyScanConfigs",
		[]string{"Log"},
		[]string{testContext.TestDemarcation()})
	if ! testContext.AssertErrIsNil(err, "") { return nil, nil }
	
	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMaps []map[string]interface{}
	responseMaps, err = rest.ParseResponseBodyToMaps(resp.Body)
	if err != nil { fmt.Println(err.Error()); return nil, nil }
	var retConfigIds []string = make([]string, 0)
	for _, responseMap := range responseMaps {
		rest.PrintMap(responseMap)
		
		if retId, isType := responseMap["Id"].(string); (! isType) || (retId == "") {
			testContext.FailTest()
		} else {
			retConfigIds = append(retConfigIds, retId)
		}
		if retProviderName, isType := responseMap["ProviderName"].(string); (! isType) || (retProviderName == "") { testContext.FailTest() }
	}

	testContext.PassTestIfNoFailures()
	return responseMaps, retConfigIds
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
		[]string{"Log", "RepoId", "ScanConfigName"},
		[]string{testContext.TestDemarcation(), repoId, scanConfigName})
	if ! testContext.AssertErrIsNil(err, "") { return "" }
	
	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if ! testContext.AssertErrIsNil(err, "") { return "" }

	var retScanConfigId string = ""
	var scanConfigIdIsType bool
	if retScanConfigId, scanConfigIdIsType = responseMap["Id"].(string); (! scanConfigIdIsType) || (retScanConfigId == "") { testContext.FailTest() }
	if retProviderName, isType := responseMap["ProviderName"].(string); (! isType) || (retProviderName == "") { testContext.FailTest() }
	if retFlagId, isType := responseMap["FlagId"].(string); (! isType) || (retFlagId == "") { testContext.FailTest() }
	testContext.PassTestIfNoFailures()
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
		[]string{"Log", "ScanConfigId"},
		[]string{testContext.TestDemarcation(), scanConfigId})
	if ! testContext.AssertErrIsNil(err, "") { return false }
	
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
			testContext.PassTestIfNoFailures()
			return true
		}	
	}
		
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if ! testContext.AssertErrIsNil(err, "") { return false }

	if retStatus, isType := responseMap["Status"].(string); (! isType) || (retStatus == "") { testContext.FailTest() }
	if retMessage, isType := responseMap["Message"].(string); (! isType) || (retMessage == "") { testContext.FailTest() }

	testContext.PassTestIfNoFailures()
	return testContext.CurrentTestPassed
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryGetMyFlags() []string {
	testContext.StartTest("TryGetMyFlags")
	
	fmt.Println("TryGetMyFlags:1")  // debug
	var resp *http.Response
	var err error
	resp, err = testContext.SendPost(testContext.SessionId,
		"getMyFlags",
		[]string{"Log"},
		[]string{testContext.TestDemarcation()})
	fmt.Println("TryGetMyFlags:2")  // debug
	if ! testContext.AssertErrIsNil(err, "while performing SendPost") { return nil }
	fmt.Println("TryGetMyFlags:3")  // debug
	
	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	fmt.Println("TryGetMyFlags:A")  // debug
	
	var responseMaps []map[string]interface{}
	responseMaps, err = rest.ParseResponseBodyToMaps(resp.Body)
	if err != nil { fmt.Println(err.Error()); return nil }
	fmt.Println("TryGetMyFlags:B")  // debug
	var retFlagIds []string = make([]string, 0)
	for _, responseMap := range responseMaps {
		fmt.Println("TryGetMyFlags:C")  // debug
		rest.PrintMap(responseMap)
		
		fmt.Println("TryGetMyFlags:D")  // debug
		if retFlagId, isType := responseMap["FlagId"].(string); (! isType) || (retFlagId == "") {
			testContext.FailTest()
		} else {
			retFlagIds = append(retFlagIds, retFlagId)
		}
		fmt.Println("TryGetMyFlags:E")  // debug
		if retRepoId, isType := responseMap["RepoId"].(string); (! isType) || (retRepoId == "") { testContext.FailTest() }
		if retName, isType := responseMap["Name"].(string); (! isType) || (retName == "") { testContext.FailTest() }
		if retImageURL, isType := responseMap["ImageURL"].(string); (! isType) || (retImageURL == "") { testContext.FailTest() }
		fmt.Println("TryGetMyFlags:F")  // debug
	}

	fmt.Println("TryGetMyFlags:G")  // debug
	fmt.Println(fmt.Sprintf("Returning %d flag ids", len(retFlagIds)))
	testContext.PassTestIfNoFailures()
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
		[]string{"Log", "RepoId", "FlagName"},
		[]string{testContext.TestDemarcation(), repoId, flagName})
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
	testContext.PassTestIfNoFailures()
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
		[]string{"Log", "FlagId"},
		[]string{testContext.TestDemarcation(), flagId})
	if ! testContext.AssertErrIsNil(err, "") { return false }
	
	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if ! testContext.AssertErrIsNil(err, "") { return false }

	if retStatus, isType := responseMap["Status"].(string); (! isType) || (retStatus == "") { testContext.FailTest() }
	if retMessage, isType := responseMap["Message"].(string); (! isType) || (retMessage == "") { testContext.FailTest() }

	testContext.PassTestIfNoFailures()
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
		[]string{"Log", "ImageId"},
		[]string{testContext.TestDemarcation(), imageId})
	if ! testContext.AssertErrIsNil(err, "") { return false }
	
	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
	
	var responseMap map[string]interface{}
	responseMap, err = rest.ParseResponseBodyToMap(resp.Body)
	if ! testContext.AssertErrIsNil(err, "") { return false }

	if retStatus, isType := responseMap["Status"].(string); (! isType) || (retStatus == "") { testContext.FailTest() }
	if retMessage, isType := responseMap["Message"].(string); (! isType) || (retMessage == "") { testContext.FailTest() }

	testContext.PassTestIfNoFailures()
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
		[]string{"Log"},
		[]string{testContext.TestDemarcation()})
	if ! testContext.AssertErrIsNil(err, "") { return }
	
	testContext.PassTestIfNoFailures()
	if ! testContext.Verify200Response(resp) { testContext.FailTest() }
}
