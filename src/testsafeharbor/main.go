/*******************************************************************************
 * Perform independent end-to-end ("behavioral") tests on the SafeHarbor server.
 * It is assumed that the SafeHarbor server is running on localhost:6000.
 */

package main

import (
	"fmt"
	"net/http"
	"os"
	"flag"
	
	// My packages:
	//"testsafeharbor/rest"
	"testsafeharbor/utils"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
	flag.PrintDefaults()
}

/*******************************************************************************
 * Set the session Id as a cookie.
 */
func setSessionId(req *http.Request, sessionId string) {
	
	// Set cookie containing the session Id.
	var cookie = &http.Cookie{
		Name: "SessionId",
		Value: sessionId,
		//Path: 
		//Domain: 
		//Expires: 
		//RawExpires: 
		MaxAge: 86400,  // 24 hrs
		Secure: false,  //....change to true later.
		HttpOnly: true,
		//Raw: 
		//Unparsed: 
	}
	
	req.AddCookie(cookie)
}

func main() {
	
	var help *bool = flag.Bool("help", false, "Provide help instructions.")
	var hostname *string = flag.String("h", "localhost", "Internet address of server.")
	var port *string = flag.String("p", "80", "Port server is on.")
	var stopOnFirstError *bool = flag.Bool("stop", false, "Provide help instructions.")
	var doNotPerformDockerTests *bool = flag.Bool("n", false, "Do not perform docker tests.")

	flag.Parse()

	if flag.NArg() > 0 {
		usage()
		os.Exit(2)
	}
	
	if *help {
		usage()
		os.Exit(0)
	}
	
	var testContext = utils.NewTestContext(*hostname, *port, setSessionId,
		*stopOnFirstError, *doNotPerformDockerTests)
		
	fmt.Println("Note: Ensure that the docker daemon is running on the server,",
		"and that python 2 is installed on the server. To start the docker daemon",
		"run 'sudo service docker start'")
	fmt.Println()
	
	var responseMap map[string]interface{}
	
	// Verify that we can create a realm without being logged in first.
	var realm4Id string
	var user4Id string
	var user4AdminRealms []interface{}
	realm4Id, user4Id, user4AdminRealms = testContext.TryCreateRealmAnon("realm4", "realm 4 Org",
		"realm4admin", "realm 4 Admin Full Name", "realm4admin@gmail.com", "realm4adminpswd")
	testContext.AssertThat(realm4Id != "", "Realm Id is empty")
	testContext.AssertThat(user4Id != "", "User Id is empty")
	testContext.AssertThat(len(user4AdminRealms) == 1, "Wrong number of admin realms")
	
	// Verify that we can log in as the admin user that we just created.
	var sessionId string
	var IsAdmin bool
	sessionId, IsAdmin = testContext.TryAuthenticate("realm4admin", "realm4adminpswd", true)
	testContext.SessionId = sessionId
	testContext.IsAdmin = IsAdmin
	fmt.Println("sessionId =", sessionId)
	
	// -------------------------------
	// User id realm4admin is authenticated.
	//
	
	// Verify that the authenticated user is an admin user.
	testContext.AssertThat(testContext.IsAdmin, "User is not flagged as admin")
	
	// Log in so that we can do stuff.
	sessionId, IsAdmin = testContext.TryAuthenticate("testuser1", "Password1", true)
	testContext.SessionId = sessionId
	testContext.IsAdmin = IsAdmin
	fmt.Println("sessionId =", sessionId)
	
	// -------------------------------
	// User id testuser1 is authenticated.
	//
	
	// Verify that the authenticated user is NOT an admin user.
	testContext.AssertThat(! testContext.IsAdmin, "User is flagged as admin")
	
	// Test ability to create a realm.
	var realmId string = testContext.TryCreateRealm("myrealm", "A Big Company", "bigshotadmin")
	testContext.AssertThat(realmId != "", "TryCreateRealm failed")
	
	// Test ability to create a user for the realm.
	var userId string = "jdoe"
	var userName string = "John Doe"
	var johnDoeUserObjId string
	var johnDoeAdminRealms []interface{}
	johnDoeUserObjId, johnDoeAdminRealms = testContext.TryCreateUser(userId, userName,
		"johnd@gmail.com", "weakpswd", realmId)
	testContext.AssertThat(johnDoeUserObjId != "", "TryCreateUser failed")
	testContext.AssertThat(len(johnDoeAdminRealms) == 0, "Wrong number of admin realms")
	
	// Login as the user that we just created.
	sessionId, IsAdmin = testContext.TryAuthenticate(userId, "weakpswd", true)
	testContext.SessionId = sessionId
	testContext.IsAdmin = IsAdmin
	
	// -------------------------------
	// User id jdoe is authenticated
	//
	
	// Test ability to create a realm.
	var jrealm1Id string = testContext.TryCreateRealm("johnsfirstrealm",
		"Johns Little Outfit", "john")
	testContext.AssertThat(jrealm1Id != "", "TryCreateRealm failed")
	
	// Test ability to create a realm.
	var jrealm2Id string = testContext.TryCreateRealm("johnssecondrealm",
		"Johns Next Venture", "admin")
	testContext.AssertThat(jrealm2Id != "", "TryCreateRealm failed")
	
	var sarahConnorUserObjId string
	var sarahConnorAdminRealms []interface{}
	sarahConnorUserObjId, sarahConnorAdminRealms = testContext.TryCreateUser("sconnor", 
		"Sarah Connor", "sarahc@sky.net", "IllMakePancakes", jrealm2Id)
	testContext.AssertThat(sarahConnorUserObjId != "", "TryCreateUser failed")
	testContext.AssertThat(len(sarahConnorAdminRealms) == 0, "Wrong number of admin realms")
	
	var johnConnorUserObjId string
	var johnConnorAdminRealms []interface{}
	johnConnorUserObjId, johnConnorAdminRealms = testContext.TryCreateUser("jconnor",
		"John Connor", "johnc@sky.net", "ILoveCameron", jrealm2Id)
	testContext.AssertThat(johnConnorUserObjId != "", "TryCreateUser failed")
	testContext.AssertThat(len(johnConnorAdminRealms) == 0, "Wrong number of admin realms")
	
	// Test ability create a repo.
	var repoId string = testContext.TryCreateRepo(realmId, "johnsrepo",
		"A very fine repo", "")
	testContext.AssertThat(repoId != "", "TryCreateRepo failed")
		
	// Test ability create another repo.
	var repo2Id string = testContext.TryCreateRepo(realmId, "susansrepo",
		"A super fine repo", "")
	testContext.AssertThat(repo2Id != "", "TryCreateRepo failed")
		
	// Test ability to upload a Dockerfile.
	var dockerfileId string = testContext.TryAddDockerfile(repoId, "Dockerfile", "A fine dockerfile")
	testContext.AssertThat(dockerfileId != "", "TryAddDockerfile failed")
	
	// Test ability to list the Dockerfiles in a repo.
	var dockerfileNames []string = testContext.TryGetDockerfiles(repoId)
	testContext.AssertThat(len(dockerfileNames) == 1, "Wrong number of dockerfiles")
	
	// Test ability to retrieve user by user id from realm.
	var userObjId string
	var userAdminRealms []interface{}
	responseMap = testContext.TryGetUserDesc(userId)
	var obj = responseMap["Id"]
	var isType bool
	userObjId, isType = obj.(string)
	testContext.AssertThat(isType, "Wrong type for Id")
	obj = responseMap["CanModifyTheseRealms"]
	userAdminRealms, isType = obj.([]interface{})
	testContext.AssertThat(isType, "Wrong type for CanModifyTheseRealms")
	testContext.AssertThat(userObjId == johnDoeUserObjId, "Looking up user by user id failed")
	testContext.AssertThat(len(userAdminRealms) == 2, "Wrong number of admin realms")
	
	var repoIds []string = testContext.TryGetRealmRepos(realmId)
	testContext.AssertThat(len(repoIds) == 2, "Number of repo Ids returned was " +
		string(len(repoIds)) + ", expected 2")
	
	var realmIds []string = testContext.TryGetAllRealms()
	// Assumes that server is in debug mode, which creates test data.
	testContext.AssertThat(len(realmIds) == 5, "Wrong number of realms found")
	
	userObjId, userAdminRealms = testContext.TryGetMyDesc(true)
	testContext.AssertThat(userObjId == johnDoeUserObjId,
		"Returned user obj id was " + userObjId)
	testContext.AssertThat(len(userAdminRealms) == 2, "Wrong number of admin realms")
	
	
	var myRealms []string = testContext.TryGetMyRealms()
	testContext.AssertThat(len(myRealms) == 2, fmt.Sprintf(
		"Only returned %d realms", len(myRealms)))
	
	
	var myRepos []string = testContext.TryGetMyRepos()
	testContext.AssertThat(len(myRepos) == 2, fmt.Sprintf(
		"Only returned %d repos", len(myRepos)))
	
	var realm3Id string
	var user3Id string
	var user3AdminRealms []interface{}
	realm3Id, user3Id, user3AdminRealms = testContext.TryCreateRealmAnon("realm3", "Realm 3 Org",
		"realm3admin", "Realm 3 Admin Full Name", "realm3admin@gmail.com", "realm3adminpswd")
	testContext.AssertThat(realm3Id != "", "Realm Id is empty")
	testContext.AssertThat(user3Id != "", "User Id is empty")
	testContext.AssertThat(len(user3AdminRealms) == 1, "Wrong number of admin realms")
	
	// Restore user context to what it was before we called TryCreateRealmAnon.
	sessionId, IsAdmin = testContext.TryAuthenticate(userId, "weakpswd", true)
	testContext.SessionId = sessionId
	testContext.IsAdmin = IsAdmin
	
	var myDockerfileIds []string = testContext.TryGetMyDockerfiles()
	testContext.AssertThat(len(myDockerfileIds) == 1, "Wrong number of dockerfiles")
		
	var realmUsers []string = testContext.TryGetRealmUsers(jrealm2Id)
	testContext.AssertThat(len(realmUsers) == 2, "Wrong number of realm users")
	
	var group1Id string = testContext.TryCreateGroup(jrealm2Id, "mygroup",
		"For Overthrowning Skynet", false)
	testContext.AssertThat(group1Id != "", "Empty group Id returned")
	
	var success bool = testContext.TryAddGroupUser(group1Id, johnConnorUserObjId)
	testContext.AssertThat(success, "TryAddGroupUser failed")
	
	success = testContext.TryAddGroupUser(group1Id, sarahConnorUserObjId)
	testContext.AssertThat(success, "TryAddGroupUser failed")
	
	var myGroupUsers []string = testContext.TryGetGroupUsers(group1Id)
	testContext.AssertThat(len(myGroupUsers) == 2, "Wrong number of group users")
	
	var jrealm2GroupIds []string = testContext.TryGetRealmGroups(jrealm2Id)
	testContext.AssertThat(len(jrealm2GroupIds) == 1, "Wrong number of realm groups")
	
	var myObjId string
	var myAdminRealms []interface{}
	myObjId, myAdminRealms = testContext.TryGetMyDesc(true)
	testContext.AssertThat(len(myAdminRealms) == 3, "Wrong number of admin realms")
	
	success = testContext.TryAddGroupUser(group1Id, myObjId)
	testContext.AssertThat(success, "TryAddGroupUser failed")
	
	var myGroupIds []string = testContext.TryGetMyGroups()
	testContext.AssertThat(len(myGroupIds) == 1, "Wrong number of groups")
	
	var perms1 []bool = []bool{false, true, false, true, true}
	var retPerms1 []bool = testContext.TrySetPermission(user3Id, dockerfileId, perms1)
	var expectedPerms1 []bool = []bool{false, true, false, true, true}
	for i, p := range retPerms1 {
		testContext.AssertThat(p == expectedPerms1[i], "Returned permission does not match")
	}
	
	var perms2 []bool = testContext.TryGetPermission(user3Id, dockerfileId)
	if perms2 != nil {
		for i, p := range perms1 {
			testContext.AssertThat(p == perms2[i], "Returned permission does not match")
		}
	}
		
	var perms3 []bool = []bool{false, false, true, true, true}
	var retPerms3 []bool = testContext.TryAddPermission(user3Id, dockerfileId, perms3)
	if retPerms3 != nil {
		var expectedPerms3 []bool = []bool{false, true, true, true, true}
		for i, p := range retPerms3 {
			testContext.AssertThat(p == expectedPerms3[i], "Returned permission does not match")
		}
	}
	
	if testContext.TryRemPermission(user3Id, dockerfileId) {
		var retPerms4 []bool = testContext.TryGetPermission(user3Id, dockerfileId)
		var expectedPerms4 []bool = []bool{false, false, false, false, false}
		for i, p := range retPerms4 {
			fmt.Println(fmt.Sprintf("\tret perm[%d]: %#v; exp perm[%d]: %#v", i, p, i, expectedPerms4[i]))
			testContext.AssertThat(p == expectedPerms4[i], "Returned permission does not match")
		}
	}

	var group2Id string = testContext.TryCreateGroup(jrealm2Id, "MySecondGroup",
		"For Overthrowning Skynet Again", true)
	testContext.AssertThat(group2Id != "", "Empty group Id returned")
	var myGroups []string = testContext.TryGetMyGroups()
	testContext.AssertThat(len(myGroups) == 2, "Wrong number of groups")
	
	// Test ability create a repo and upload a dockerfile at the same time.
	var repo5Id string = testContext.TryCreateRepo(realmId, "zippysrepo",
		"A super smart repo", "dockerfile")
	testContext.AssertThat(repo5Id != "", "TryCreateRepo failed")
		
	testContext.TryGetGroupDesc(group2Id)
	
	testContext.TryGetRepoDesc(repoId)
	
	testContext.TryGetDockerfileDesc(dockerfileId)
	
	testContext.TryReplaceDockerfile(dockerfileId, "Dockerfile2", "The boo/ploink one")
		
	if testContext.TryAddGroupUser(group2Id, sarahConnorUserObjId) {
		if testContext.TryAddGroupUser(group2Id, johnConnorUserObjId) {
			var userIdsBeforeRemoval []string = testContext.TryGetGroupUsers(group2Id)
			testContext.TryRemGroupUser(group2Id, sarahConnorUserObjId)
			var userIdsAfterRemoval []string = testContext.TryGetGroupUsers(group2Id)
			testContext.AssertThat(len(userIdsBeforeRemoval)-len(userIdsAfterRemoval)==1,
				fmt.Sprintf("Before: %d users, after: %d users", userIdsBeforeRemoval, userIdsAfterRemoval))
		}
	}
	
	testContext.TryGetScanProviders()

	var config1Id string = testContext.TryDefineScanConfig("My Config 1",
		"A very find config", repoId, "clair", "", "Seal.png", []string{}, []string{})
	testContext.AssertThat(config1Id != "", "No ScanConfig Id was returned")
	
	responseMap = testContext.TryDefineFlag(repoId, "myflag", "A really boss flag", "Seal2.png")
	if testContext.CurrentTestPassed {
		var obj interface{} = responseMap["FlagId"]
		var flagId string
		var isType bool
		flagId, isType = obj.(string)
		testContext.AssertThat(isType, "Returned FlagId is not a string")
		if flagId == "" { testContext.FailTest() } else {
			var flagName string = testContext.TryGetFlagDesc(flagId, true)
			if flagName != "myflag" { testContext.FailTest() }
			
			var flagIds []string = testContext.TryGetMyFlags()
			testContext.AssertThat(utils.ContainsString(flagIds, flagId),
				"Flag Id " + flagId + " not returned")
			
			var fId string = testContext.TryGetFlagDescByName(repoId, "myflag")
			testContext.AssertThat(fId == flagId, "Flag not found by name")
		}
	}
	
	responseMap = testContext.TryGetScanConfigDesc(config1Id, true)
	var flag1Id string
	if testContext.CurrentTestPassed {
		var obj = responseMap["FlagId"]
		var isType bool
		flag1Id, isType = obj.(string)
		testContext.AssertThat(isType, "Wrong type for returned FlagId")
		if flag1Id == "" { testContext.FailTest() } else {
			var flagName string = testContext.TryGetFlagDesc(flag1Id, true)
			if flagName == "" { testContext.FailTest() } else {
				if flagName != "My Config 1" { testContext.FailTest() }
			}
			var size int64 = testContext.TryGetFlagImage(flag1Id, "ShouldBeIdenticalToSeal2.png")
			var fileInfo os.FileInfo
			var err error
			fileInfo, err = os.Stat("Seal.png")
			if testContext.AssertErrIsNil(err, "") {
				testContext.AssertThat(fileInfo.Size() == size, "File has wrong size")
			}
		}
	}

	// Replace the Scan Config's flag with a new flag.
	testContext.TryUpdateScanConfig(config1Id, "", "", "", "", "Seal2.png",
		[]string{}, []string{})
	var scanConfig1Map map[string]interface{}
	scanConfig1Map = testContext.TryGetScanConfigDesc(config1Id, true)
	if testContext.CurrentTestPassed {
		// Id string
		// ProviderName string
		// SuccessExpression string
		// FlagId string
		// ParameterValueDescs []*ParameterValueDesc
		var newFlagId = scanConfig1Map["FlagId"]
		testContext.AssertThat(newFlagId != "", "FlagId returned empty")
	}
	
	var configIds []string
	_, configIds = testContext.TryGetMyScanConfigs()
	testContext.AssertThat(utils.ContainsString(configIds, config1Id),
		"Scan config not found")
	
	var configId string = testContext.TryGetScanConfigDescByName(repoId, "My Config 1")
	testContext.AssertThat(configId == config1Id, "Did not find scan config")
	
	// ....Test that permissions work.
	
	
	// Verify that we can update our password.
	if testContext.TryChangePassword(userId, "weakpswd", "password2") {
		testContext.TryLogout()
		testContext.TryAuthenticate(userId, "weakpswd", false)
		testContext.SessionId, testContext.IsAdmin = testContext.TryAuthenticate(userId, "password2", true)
	}
	
	// -------------------------------
	// User id jdoe is authenticated.
	//
	
	// Test ability to make a private image available to the SafeHarbor closed community.
	
	// Test ability to make a private image available to another user.
	

	if testContext.PerformDockerTests {
		// Test ability to build image from a dockerfile.
		var dockerImage1ObjId string
		var image1Id string
		dockerImage1ObjId, image1Id = testContext.TryExecDockerfile(repoId, dockerfileId, "myimage")
		testContext.AssertThat(dockerImage1ObjId != "", "TryExecDockerfile failed - obj id is nil")
		testContext.AssertThat(image1Id != "", "TryExecDockerfile failed - docker image id is nil")
	
		// Test ability to list the images in a repo.
		var imageNames []string = testContext.TryGetImages(repoId)
		testContext.AssertThat(len(imageNames) == 1, "Wrong number of images")
	
		var myDockerImageIds []string = testContext.TryGetMyDockerImages()
		testContext.AssertThat(len(myDockerImageIds) == 1, "Wrong number of docker images")
	
		var scanScore string = testContext.TryScanImage(config1Id, dockerImage1ObjId)
		testContext.AssertThat(scanScore != "", "Empty scan score")

		var dockerImage2ObjId string
		var image2Id string
		dockerImage2ObjId, image2Id = testContext.TryAddAndExecDockerfile(repoId,
			"My second image", "myimage2", "Dockerfile")
		testContext.AssertThat(dockerImage2ObjId != "", "TryExecDockerfile failed - obj id is nil")
		testContext.AssertThat(image2Id != "", "TryExecDockerfile failed - docker image id is nil")
	
		testContext.TryDownloadImage(dockerImage2ObjId, "MooOinkImage")
		responseMap = testContext.TryGetDockerImageDesc(dockerImage2ObjId, true)
		if testContext.CurrentTestPassed {
			// Check image signature.
			var image2Signature []byte
			var err error
			image2Signature, err = utils.ComputeFileSignature("MooOinkImage")
			if testContext.AssertErrIsNil(err, "Unable to compute signature") {
				var obj interface{} = responseMap["Signature"]
				var isType bool
				var sig []byte
				sig, isType = obj.([]byte)
				if testContext.AssertThat(isType, "Signature is not an array of byte") {
					for i, b := range sig {
						if ! testContext.AssertThat(b == image2Signature[i], "Wrong signature") { break }
					}
				}
			}
		}
	
		var eventIds []string = testContext.TryGetUserEvents(userId)
		testContext.AssertThat(len(eventIds) == 2, "Wrong number of user events")
	
		eventIds = testContext.TryGetDockerImageEvents(dockerImage1ObjId)
		testContext.AssertThat(len(eventIds) == 1, "Wrong number of image events")
	
		responseMap = testContext.TryGetDockerImageStatus(dockerImage1ObjId)
		if testContext.CurrentTestPassed {
			testContext.AssertThat(responseMap["EventId"] != "", "No image status")
			testContext.AssertThat(responseMap["ScanConfigId"] == config1Id,
				"Wrong scan config Id")
			testContext.AssertThat(responseMap["ProviderName"] == "clair",
				"Wrong provider")
		}
	
		eventIds = testContext.TryGetDockerfileEvents(dockerfileId)
		testContext.AssertThat(len(eventIds) == 1, "Wrong number of image events")
		
		if testContext.TryRemDockerImage(dockerImage1ObjId) {
			testContext.TryGetDockerImageDesc(dockerImage1ObjId, false)
		}
	}

	// Test that we can disable a user.
	if testContext.TryDisableUser(johnConnorUserObjId) {
		// Now see if that user can authenticate.
		var isAdmin bool
		_, isAdmin = testContext.TryAuthenticate("jconnor", "ILoveCameron", false)
		if testContext.CurrentTestPassed {
			testContext.AssertThat(!isAdmin, "Error: user is an admin but should not be")		
		}
		if testContext.TryReenableUser(johnConnorUserObjId) {
			testContext.SessionId, testContext.IsAdmin = testContext.TryAuthenticate("jconnor", "ILoveCameron", true)
		}
	}
	
	testContext.TryDeleteGroup(group2Id)
	
	// Try moving a user from one realm to another.
	if testContext.TryMoveUserToRealm(sarahConnorUserObjId, realm3Id) {
		// Verify that Sarah is no longer in her realm.
		responseMap = testContext.TryGetUserDesc("sconnor")
		if testContext.CurrentTestPassed {
			// Verify that Sarah is in John's realm.
			testContext.AssertThat(responseMap["RealmId"] == realm3Id,
				"Error: Sarah Connor does not belong to John's realm")
		}
	}
	
	if testContext.TryRemScanConfig(config1Id, true) {
		testContext.TryGetScanConfigDesc(config1Id, false)
	}
	
	if testContext.TryRemFlag(flag1Id) {
		testContext.TryGetFlagDesc(flag1Id, false)
	}
	
	// Test ability to log out.
	if testContext.TryLogout() {
		testContext.TryGetMyDesc(false)
	} else {
		testContext.AssertThat(false, "Unable to log out")
	}
	
	if testContext.TryDeactivateRealm(jrealm2Id) {
		testContext.TryGetMyDesc(false)
	}
	
	// Test ability to clear the entire database and docker repository.
	testContext.TryClearAll()
	
	
	
	fmt.Println()
	fmt.Println(fmt.Sprintf("%d tests failed out of %d:", testContext.NoOfTestsThatFailed,
		testContext.NoOfTests))
	for i, testName := range testContext.GetTestsThatFailed() {
		if i > 0 { fmt.Print(", ") }
		fmt.Print(testName)
	}
	fmt.Println()
}
