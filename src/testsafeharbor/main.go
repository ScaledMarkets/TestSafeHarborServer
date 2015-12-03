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
	"testsafeharbor/rest"
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
	
	var testContext *utils.TestContext = &utils.TestContext{
		RestContext: *rest.CreateRestContext(*hostname, *port, setSessionId),
		SessionId: "",
		StopOnFirstError: *stopOnFirstError,
		PerformDockerTests: ! (*doNotPerformDockerTests),
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
	testContext.AssertThat(realm4Id != "", "Realm Id is empty")
	testContext.AssertThat(user4Id != "", "User Id is empty")
	testContext.AssertThat(len(user4AdminRealms) == 1, "Wrong number of admin realms")
	
	// Verify that we can log in as the admin user that we just created.
	var sessionId string
	var IsAdmin bool
	sessionId, IsAdmin = testContext.TryAuthenticate("realm4admin", "realm4adminpswd")
	testContext.SessionId = sessionId
	testContext.IsAdmin = IsAdmin
	fmt.Println("sessionId =", sessionId)
	
	// Verify that the authenticated user is an admin user.
	testContext.AssertThat(testContext.IsAdmin, "User is not flagged as admin")
	
	// Log in so that we can do stuff.
	sessionId, IsAdmin = testContext.TryAuthenticate("testuser1", "password1")
	testContext.SessionId = sessionId
	testContext.IsAdmin = IsAdmin
	fmt.Println("sessionId =", sessionId)
	
	// Verify that the authenticated user is NOT an admin user.
	testContext.AssertThat(! testContext.IsAdmin, "User is flagged as admin")
	
	// Test ability to create a realm.
	var realmId string = testContext.TryCreateRealm("MyRealm", "A Big Company", "bigshotadmin")
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
	sessionId, IsAdmin = testContext.TryAuthenticate(userId, "password1")
	testContext.SessionId = sessionId
	testContext.IsAdmin = IsAdmin
	
	// Test ability to create a realm.
	var jrealm1Id string = testContext.TryCreateRealm("Johns First Realm",
		"Johns Little Outfit", "john")
	testContext.AssertThat(jrealm1Id != "", "TryCreateRealm failed")
	
	// Test ability to create a realm.
	var jrealm2Id string = testContext.TryCreateRealm("Johns Second Realm",
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
	var repoId string = testContext.TryCreateRepo(realmId, "Johns Repo",
		"A very fine repo", "")
	testContext.AssertThat(repoId != "", "TryCreateRepo failed")
		
	// Test ability create another repo.
	var repo2Id string = testContext.TryCreateRepo(realmId, "Susans Repo",
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
	userObjId, userAdminRealms = testContext.TryGetRealmUser(realmId, userId)
	testContext.AssertThat(userObjId == johnDoeUserObjId, "Looking up user by user id failed")
	testContext.AssertThat(len(userAdminRealms) == 2, "Wrong number of admin realms")
	
	//var msg string = testContext.TryAddRealmUser(....realmId, userObjId)
	
	var repoIds []string = testContext.TryGetRealmRepos(realmId)
	testContext.AssertThat(len(repoIds) == 2, "Number of repo Ids returned was " +
		string(len(repoIds)) + ", expected 2")
	
	var realmIds []string = testContext.TryGetAllRealms()
	// Assumes that server is in debug mode, which creates test data.
	testContext.AssertThat(len(realmIds) == 5, "Wrong number of realms found")
	
	userObjId, userAdminRealms = testContext.TryGetMyDesc()
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
	var saveSessionId string = testContext.SessionId
	realm3Id, user3Id, user3AdminRealms = testContext.TryCreateRealmAnon("Realm3", "Realm 3 Org",
		"realm3admin", "Realm 3 Admin Full Name", "realm3admin@gmail.com", "realm3adminpswd")
	testContext.AssertThat(realm3Id != "", "Realm Id is empty")
	testContext.AssertThat(user3Id != "", "User Id is empty")
	testContext.AssertThat(len(user3AdminRealms) == 1, "Wrong number of admin realms")
	
	// Restore user context to what it was before we called TryCreateRealmAnon.
	testContext.SessionId = saveSessionId
	
	if testContext.PerformDockerTests {
		// Test ability to build image from a dockerfile.
		var dockerImageObjId string
		var imageId string
		dockerImageObjId, imageId = testContext.TryExecDockerfile(repoId, dockerfileId, "myimage")
		testContext.AssertThat(dockerImageObjId != "", "TryExecDockerfile failed - obj id is nil")
		testContext.AssertThat(imageId != "", "TryExecDockerfile failed - docker image id is nil")
	
		// Test ability to list the images in a repo.
		var imageNames []string = testContext.TryGetImages(repoId)
		testContext.AssertThat(len(imageNames) == 1, "Wrong number of images")
	
		var myDockerImageIds []string = testContext.TryGetMyDockerImages()
		testContext.AssertThat(len(myDockerImageIds) == 1, "Wrong number of docker images")
	}
	
	var myDockerfileIds []string = testContext.TryGetMyDockerfiles()
	testContext.AssertThat(len(myDockerfileIds) == 1, "Wrong number of dockerfiles")
		
	var realmUsers []string = testContext.TryGetRealmUsers(jrealm2Id)
	testContext.AssertThat(len(realmUsers) == 2, "Wrong number of realm users")
	
	var group1Id string = testContext.TryCreateGroup(jrealm2Id, "MyGroup",
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
	myObjId, myAdminRealms = testContext.TryGetMyDesc()
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
	for i, p := range perms1 {
		testContext.AssertThat(p == perms2[i], "Returned permission does not match")
	}
		
	var perms3 []bool = []bool{false, false, true, true, true}
	var retPerms3 []bool = testContext.TryAddPermission(user3Id, dockerfileId, perms3)
	var expectedPerms3 []bool = []bool{false, true, true, true, true}
	for i, p := range retPerms3 {
		testContext.AssertThat(p == expectedPerms3[i], "Returned permission does not match")
	}
	
	var group2Id string = testContext.TryCreateGroup(jrealm2Id, "MySecondGroup",
		"For Overthrowning Skynet Again", true)
	testContext.AssertThat(group2Id != "", "Empty group Id returned")
	var myGroups []string = testContext.TryGetMyGroups()
	testContext.AssertThat(len(myGroups) == 2, "Wrong number of groups")
	
	// Test ability create a repo and upload a dockerfile at the same time.
	var repo5Id string = testContext.TryCreateRepo(realmId, "Zippys Repo",
		"A super smart repo", "dockerfile")
	testContext.AssertThat(repo5Id != "", "TryCreateRepo failed")
		
	var dockerImageObjId string
	if testContext.PerformDockerTests {
		var imageId string
		dockerImageObjId, imageId = testContext.TryAddAndExecDockerfile(repoId,
			"My second image", "myimage2", "Dockerfile")
		testContext.AssertThat(dockerImageObjId != "", "TryExecDockerfile failed - obj id is nil")
		testContext.AssertThat(imageId != "", "TryExecDockerfile failed - docker image id is nil")
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
	
	testContext.TryGetScanProviders()

	if testContext.PerformDockerTests {
		var config1Id string = testContext.TryDefineScanConfig("My Config 1",
			"A very find config", repoId, "clair",
			"http://someimage.com", "http://someimage.com", []string{}, []string{})
		testContext.AssertThat(config1Id != "", "No ScanConfig Id was returned")
	
		var scanScore string = testContext.TryScanImage(config1Id, dockerImageObjId)
		testContext.AssertThat(scanScore != "", "Empty scan score")
	}

	// Test that permissions work.
	
	
	
	
	// Test ability to receive progress while a Dockerfile is processed.
	
	// Test ability to make a private image available to the SafeHarbor closed community.
	
	// Test ability to make a private image available to another user.
	
	// Test ability to clear the entire database and docker repository.
	testContext.TryClearAll()
	
	
	fmt.Println()
	fmt.Println(fmt.Sprintf("%d tests failed out of %d", utils.NoOfTestsThatFailed,
		utils.NoOfTests))
}
