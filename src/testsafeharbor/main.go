/*******************************************************************************
 * Perform independent end-to-end ("behavioral") tests on the SafeHarbor server.
 */

package main

import (
	"fmt"
	//"net/http"
	"os"
	"flag"
	"time"
	"strings"
	"reflect"
	
	"redis"
	"goredis"
	
	// SafeHarbor packages:
	"testsafeharbor/utils"
	"testsafeharbor/rest"
)

const (
	SealURL = "https://itsonlywords55.files.wordpress.com/2010/01/seal-of-approval.jpg"
	Seal2URL = "http://thumb10.shutterstock.com/display_pic_with_logo/681547/140365213/stock-photo-seal-of-approval-quality-check-grunge-vector-on-white-background-this-graphic-illustration-140365213.jpg"
)

func main() {
	
	var testSuite = map[string]func(*utils.TestContext) {
		"Registry": TestDockerRegistry,
		"json": TestJSONDeserialization,
		"goredis": TestGoRedis,
		"redis": TestRedis,
		"CreateRealmsAndUsers": TestCreateRealmsAndUsers,
		"CreateResources": TestCreateResources,
		"CreateGroups": TestCreateGroups,
		"GetMy": TestGetMy,
		"AccessControl": TestAccessControl,
		"UpdateAndReplace": TestUpdateAndReplace,
		"Delete": TestDelete,
		"DockerFunctions": TestDockerFunctions,
	}

	var help *bool = flag.Bool("help", false, "Provide help instructions.")
	var hostname *string = flag.String("h", "localhost", "Internet address of server.")
	var port *int = flag.Int("p", 80, "Port server is on.")
	var stopOnFirstError *bool = flag.Bool("stop", false, "Stop after the first error.")
	var redisPswd *string = flag.String("redispswd", "ahdal8934k383898&*kdu&^", "Redis password")
	
	var keys []reflect.Value = reflect.ValueOf(testSuite).MapKeys()
	var allTestNames string
	for i, key := range keys {
		if i > 0 { allTestNames = allTestNames + "," }
		allTestNames = allTestNames + key.String()
	}
	var tests *string = flag.String("tests", allTestNames,
		"Perform the tests listed, comma-separated.")

	flag.Parse()

	if *help {
		fmt.Println("Help:")
		utils.Usage()
		os.Exit(0)
	}
	
	// Parse the 'tests' option to determine which tests to run.
	var testsToRun []string = strings.Split(*tests, ",")
	var testFunctionsToRun = make([]func(*utils.TestContext), 0)
	for _, testName := range testsToRun {
		var testFunction = testSuite[testName]
		if testFunction == nil {
			fmt.Println("Test '" + testName + "' not recognized")
			os.Exit(0)
		}
		testFunctionsToRun = append(testFunctionsToRun, testFunction)
	}
	
	// Prepare to run tests.
	var testContext = utils.NewTestContext(*hostname, *port, utils.SetSessionId,
		*stopOnFirstError, *redisPswd)
	testContext.Print()
	if strings.Contains(*tests, "DockerFunctions") {
		fmt.Println("Note: Ensure that the docker daemon is running on the server.",
			"To start the docker daemon, run 'sudo service docker start'.")
	}
	fmt.Println()
	
	// Run the tests, one by one.
	for _, testFunctionToRun := range testFunctionsToRun {
		testFunctionToRun(testContext)
	}
	
	// Print result summary.
	fmt.Println()
	fmt.Println(fmt.Sprintf("%d tests failed out of %d:", testContext.NoOfTestsThatFailed,
		testContext.NoOfTests))
	for i, testName := range testContext.GetTestsThatFailed() {
		if i > 0 { fmt.Print(", ") }
		fmt.Print(testName)
	}
	fmt.Println()
}

/*******************************************************************************
 * 
 */
func TestDockerRegistry(testContext *utils.TestContext) {
	
	fmt.Println("\nTest suite TestDockerRegistry------------------\n")

	// -------------------------------------
	// Test setup:
	
	// Auth:
	// https://github.com/docker/distribution/blob/master/docs/deploying.md
	var registryHost = "localhost"
	var registryPort = 5000
	var registryUserId = os.Getenv("registryUser")
	var registryPassword = os.Getenv("registryPassword")
	var testImageName = os.Getenv("TestImageName")
	var testImageTag = os.Getenv("TestImageTag")
	var downloadedImageFilePath = "DownloadedImage.tar"
	
	var registry *utils.DockerRegistry
	var err error
	
	// Test connecting to Registry.
	{
		testContext.StartTest("Open Registry connection")
		registry, err = utils.OpenDockerRegistryConnection(registryHost, registryPort,
			registryUserId, registryPassword)
		testContext.AssertErrIsNil(err, "In opening connection to docker registry")
		testContext.PassTestIfNoFailures()
	}
	
	// Test Inspect.
	{
		testContext.StartTest("Test Inspect")
		var exists bool
		exists, err = registry.ImageExists(testImageName, testImageTag)
		testContext.AssertErrIsNil(err, "While calling ImageExists")
		testContext.AssertThat(exists, "Did not find image")
		testContext.PassTestIfNoFailures()
	}
	
	// Test getting image.
	{
		testContext.StartTest("Test Getting Image")
		
		os.Remove(downloadedImageFilePath)
		var err error
		
		// Contact Registry to get image.
		err = registry.GetImage(testImageName, testImageTag, downloadedImageFilePath)
		
		// Verify that the image was retrieved properly.
		testContext.AssertErrIsNil(err, "While calling GetImage")
		var downloadedImageFile *os.File
		downloadedImageFile, err = os.OpenFile(downloadedImageFilePath, os.O_WRONLY, 0600)
		testContext.AssertErrIsNil(err, fmt.Sprintf(
			"When opening image file '%s': %s", downloadedImageFilePath, err.Error()))
		var fileInfo os.FileInfo
		fileInfo, err = downloadedImageFile.Stat()
		testContext.AssertErrIsNil(err, fmt.Sprintf(
			"When getting status of image file '%s': %s", downloadedImageFilePath, err.Error()))
		testContext.AssertThat(fileInfo.Size() > 0, "Downloaded file is size 0")
		
		testContext.PassTestIfNoFailures()
	}
	
	// Test deleting image.
	{
		//testContext.StartTest("Test Deleting Image")
		//registry.DeleteImage
		//testContext.AssertErrIsNil(err, "DeleteImage")
		//testContext.AssertThat()
		//testContext.PassTestIfNoFailures()
	}
}

/*******************************************************************************
 * Test the goredis API to verify understanding of it.
 */
func TestGoRedis(testContext *utils.TestContext) {
	
	fmt.Println("\nTest suite TestGoRedis------------------\n")

	// -------------------------------------
	// Test setup:
	
	var redis *goredis.Redis
	var err error
	
	{
		var network		= "tcp"
		var host string = testContext.GetHostname()
		var port int	= 6379
		var db			= 1
		var password	= testContext.RedisPswd
		var timeout		= 5 * time.Second
		var maxidle		= 1
		
		redis, err = goredis.Dial(&goredis.DialConfig{
			network, (host + ":" + fmt.Sprintf("%d", port)), db, password, timeout, maxidle})
		testContext.AssertErrIsNil(err, "In test setup, after Dial")
	}
	
	// -------------------------------------
	// Tests
	//
	
	{
		testContext.TryGoRedisPing(redis)
	}
	
	{
		testContext.TryGoRedisSetGetString(redis)
	}
	
	{
		testContext.TryGoRedisSet(redis)
	}
}

/*******************************************************************************
 * Test the redis API to verify understanding of it.
 * Redis bindings for go: http://redis.io/clients#go
 * Chosen binding: https://github.com/alphazero/Go-Redis
 * Alternative binding: https://github.com/hoisie/redis
 */
func TestRedis(testContext *utils.TestContext) {

	fmt.Println("\nTest suite TestRedis------------------\n")

	// -------------------------------------
	// Test setup:
	
	var client redis.Client
	
	{
		var spec *redis.ConnectionSpec =
			redis.DefaultSpec().Host(testContext.GetHostname()).Port(6379).Password(
				testContext.RedisPswd)
		var err error
		client, err = redis.NewSynchClientWithSpec(spec);
		testContext.AssertErrIsNil(err, "failed to create the client")
		if err != nil { return }
	}
	
	// -------------------------------------
	// Tests
	//
	
	{
		testContext.TryRedisPing(client)
	}
	
	{
		testContext.TryRedisSetGetString(client)
	}
	
	{
		testContext.TryRedisGetJSONObject(client)
	}
}

/*******************************************************************************
 * 
 */
func TestJSONDeserialization(testContext *utils.TestContext) {

	fmt.Println("\nTest suite TestJSONDeserialization------------------\n")

	{
		var json = "{\"abc\": 123, \"bs\": \"this_is_a_string\", " +
			"\"car\": [\"alpha\", \"beta\"], true}"
		var expected = []string{
			"{",
			"\"",
			"abc",
			"\"",
			":",
			"123",
			",",
			"\"",
			"bs",
			"\"",
			":",
			"\"",
			"this_is_a_string",
			"\"",
			",",
			"\"",
			"car",
			"\"",
			":",
			"[",
			"\"",
			"alpha",
			"\"",
			",",
			"\"",
			"beta",
			"\"",
			"]",
			",",
			"true",
			"}",
		}
		testContext.TryJsonDeserTokenizer(json, expected)
	}
	
	{
		var json = "\"this is a string\""
		var expected = "this is a string"
		testContext.TryJsonDeserString(json, expected)
		
		json = "   \"\""
		expected = ""
		testContext.TryJsonDeserString(json, expected)
		
		json = "\"1\""
		expected = "1"
		testContext.TryJsonDeserString(json, expected)
	}
	
	{
		var json = "[244, 26, 234, 221, 169, 129, 22, 245, 25, 151, 124, 137, 22, 44, 202, 205, 84, 206, 21, 99, 170, 55, 200, 12, 100, 137, 211, 73, 140, 41, 63, 10, 244, 166, 51, 24, 160, 2, 53, 171, 231, 244, 254, 58, 56, 140, 54, 4, 253, 195, 221, 75, 172, 173, 175, 10, 12, 11, 107, 0, 64, 64, 207, 187]"
		var expected []int64 = []int64{244, 26, 234, 221, 169, 129, 22, 245, 25, 151, 124, 137, 22, 44, 202, 205, 84, 206, 21, 99, 170, 55, 200, 12, 100, 137, 211, 73, 140, 41, 63, 10, 244, 166, 51, 24, 160, 2, 53, 171, 231, 244, 254, 58, 56, 140, 54, 4, 253, 195, 221, 75, 172, 173, 175, 10, 12, 11, 107, 0, 64, 64, 207, 187}
		testContext.TryJsonDeserByteArray(json, expected)
	}
	
	{
		var json = "time \"2016-01-18T15:10:03.984179856Z\""
		var expected time.Time
		var err = expected.UnmarshalJSON([]byte("\"2016-01-18T15:10:03.984179856Z\""))
		if err != nil {
			fmt.Println("test setup error")
			panic(err)
		}
		testContext.TryJsonDeserTime(json, expected)
	}
	
	{
		testContext.TryJsonDeserSimple()
	}
	
	{
		testContext.TryJsonDeserNestedType()
	}
	
	{
		var json = "{\"Id\": \"\"}"
		testContext.TryJsonDeserComplex(json)
	}
	
	{
		var json = "{\"Id\": \"\", \"ACLEntryIds\": [], \"Name\": \"testrealm\", " +
			"\"Description\": \"For Testing\", \"ParentId\": \"\"}"
		testContext.TryJsonDeserComplex(json)
	}
	
	{
		var json = "{\"Id\": \"100000006\", \"ACLEntryIds\": [], \"Name\": \"testrealm\", " +
			"\"Description\": \"For Testing\", \"ParentId\": \"\"}"
		testContext.TryJsonDeserComplex(json)
	}
	
	{
		var json = "{\"Id\": \"100000006\", \"ACLEntryIds\": [], \"Name\": \"testrealm\", " +
			"\"Description\": \"For Testing\", \"ParentId\": \"\", " +
			"\"CreationTime\": time \"2016-01-18T16:16:30.289421913Z\"}"
		testContext.TryJsonDeserComplex(json)
	}
	
	{
		var json = "{\"Id\": \"100000006\", \"ACLEntryIds\": [], \"Name\": \"testrealm\", " +
			"\"Description\": \"For Testing\", \"ParentId\": \"\", " +
			"\"CreationTime\": time \"2016-01-18T16:16:30.289421913Z\", " +
			"\"AdminUserId\": \"testuser1\", \"OrgFullName\": \"Test Org\", " +
			"\"UserObjIds\": [], \"GroupIds\": [], \"RepoIds\": [], " +
			"\"FileDirectory\": \"Repositories/100000006\"}"
		testContext.TryJsonDeserComplex(json)
	}
	
	{
		var json = "{\"Id\": \"100000007\", \"IsActive\": true, " +
			"\"Name\": \"realm 4 Admin Full Name\", " +
			"\"CreationTime\": time \"2016-02-21T18:49:08.576404647Z\", " +
			"\"RealmId\": \"100000006\", \"ACLEntryIds\": [\"100000008\"], " +
			"\"UserId\": \"realm4admin\", \"EmailAddress\": \"realm4admin@gmail.com\", " +
			"\"PasswordHash\": [244, 26, 234, 221, 169, 129, 22, 245, 25, 151, 124, 137, 22, 44, 202, 205, 84, 206, 21, 99, 170, 55, 200, 12, 100, 137, 211, 73, 140, 41, 63, 10, 244, 166, 51, 24, 160, 2, 53, 171, 231, 244, 254, 58, 56, 140, 54, 4, 253, 195, 221, 75, 172, 173, 175, 10, 12, 11, 107, 0, 64, 64, 207, 187], " +
			"\"GroupIds\": [], \"MostRecentLoginAttempts\": [], \"EventIds\": []}"
		testContext.TryJsonDeserComplex(json)
	}
}

/*******************************************************************************
 * Test ability to create realms and users within those realms.
 * Creates/uses the following:
 *	realm4
 *	realm4admin
 */
func TestCreateRealmsAndUsers(testContext *utils.TestContext) {
	
	fmt.Println("\nTest suite TestCreateRealmsAndUsers------------------\n")

	defer testContext.TryClearAll()
	
	// -------------------------------------
	// Test setup:
	
	var realm4AdminUserId = "realm4admin"
	var realm4AdminPswd = "RealmPswd"
	var joeUserId = "jdoe"
	var joePswd = "weakpswd"
	
	// -------------------------------------
	// Tests
	//
	
	var realm4Id string
	//var realm4AdminObjId string
	//var defaultUserObjId string
	
	// Verify that we can create a realm without being logged in first.
	{
		var user4AdminRealms []interface{}
		realm4Id, _, user4AdminRealms = testContext.TryCreateRealmAnon(
			"realm4", "realm 4 Org",
			realm4AdminUserId, "realm 4 Admin Full Name", "realm4admin@gmail.com", realm4AdminPswd)
		testContext.AssertThat(len(user4AdminRealms) == 1, "Wrong number of admin realms")
	}
	
	// Verify that we can log in as the admin user that we just created.
	{
		testContext.TryAuthenticate(realm4AdminUserId, realm4AdminPswd, true)
	}
	
	// -------------------------------------
	// User id realm4admin (of realm4) is authenticated.
	//
	
	// Verify that the authenticated user is an admin user.
	{
		testContext.AssertThat(testContext.IsAdmin, "User is not flagged as admin")
	}
	
	// Check that we can retrieve the users of a realm.
	{
		var realmUsers []string = testContext.TryGetRealmUsers(realm4Id)
		testContext.AssertThat(len(realmUsers) == 1, "Wrong number of realm users")
	}
	
	// Test ability to create a realm while logged in.
	{
		testContext.TryCreateRealm("my2ndrealm", "A Big Company", "bigshotadmin")
	}
	
	var johnDoeUserObjId string
	
	// Test ability to create a user for a realm.
	{
		var johnDoeAdminRealms []interface{}
		johnDoeUserObjId, johnDoeAdminRealms = testContext.TryCreateUser(
			joeUserId, "John Doe", "johnd@gmail.com", joePswd, realm4Id)
		testContext.AssertThat(len(johnDoeAdminRealms) == 0, "Wrong number of admin realms")
		fmt.Println(johnDoeUserObjId)
	}
	
	// Login as the user that we just created.
	{
		testContext.TryAuthenticate("jdoe", "weakpswd", true)
	}
	
	// -------------------------------------
	// User id jdoe is authenticated
	//
	
	// Verify that the authenticated user is not an admin user.
	{
		testContext.AssertThat(! testContext.IsAdmin, "User is flagged as admin")
	}
	
	{
		var realmIds []string = testContext.TryGetAllRealms()
		// Assumes that server is in debug mode, which creates a test realm.
		testContext.AssertThat(len(realmIds) == 2, "Wrong number of realms found")
	}
	
	// Test ability to retrieve user by user id from realm.
	{
		testContext.TryAuthenticate(realm4AdminUserId, realm4AdminPswd, true)
		var userObjId string
		var userAdminRealms []interface{}
		var responseMap = testContext.TryGetUserDesc("jdoe")
		var obj = responseMap["Id"]
		var isType bool
		userObjId, isType = obj.(string)
		testContext.AssertThat(isType, "Wrong type for Id")
		obj = responseMap["CanModifyTheseRealms"]
		userAdminRealms, isType = obj.([]interface{})
		testContext.AssertThat(isType, "Wrong type for CanModifyTheseRealms")
		testContext.AssertThat(userObjId == johnDoeUserObjId, "Looking up user by user id failed")
		testContext.AssertThat(len(userAdminRealms) == 0, "Wrong number of admin realms")
	}
}


/*******************************************************************************
 * Test ability to create resources within a realm, and retrieve info about them.
 * Creates/uses the following:
 */
func TestCreateResources(testContext *utils.TestContext) {
	
	fmt.Println("\nTest suite TestCreateResources------------------\n")

	defer testContext.TryClearAll()
	
	// -------------------------------------
	// Test setup:
	// Create a realm and an admin user for the realm, and then log in as that user.
	// Write a dockerfile to a new temp directory.
	//
	
	var realm4Id string
	//var user4Id string
	var dockerfilePath string
	var flagImagePath = "Seal.png"
	var flag2ImagePath = "Seal2.png"
	
	{
		realm4Id, _, _ = testContext.TryCreateRealmAnon(
			"realm4", "realm 4 Org", "realm4admin", "realm 4 Admin Full Name",
			"realm4admin@gmail.com", "realm4adminpswd")
		
		testContext.TryAuthenticate("realm4admin", "realm4adminpswd", true)
		
		var err error
		dockerfilePath, err = utils.CreateTempFile("Dockerfile", "FROM centos\nRUN echo moo > oink")
		if err != nil { testContext.AbortAllTests(err.Error()) }
		defer os.Remove(dockerfilePath)
		
		err = utils.DownloadFile(SealURL, flagImagePath, true)
		if err != nil { testContext.AbortAllTests(err.Error()) }
		
		err = utils.DownloadFile(Seal2URL, flag2ImagePath, true)
		if err != nil { testContext.AbortAllTests(err.Error()) }
	}
	
	// -------------------------------------
	// Tests
	//
	
	var johnsRepoId string
	//var johnsDockerfileId string
	
	// Test ability create a repo.
	{
		johnsRepoId = testContext.TryCreateRepo(realm4Id, "johnsrepo", "A very fine repo", "")
	}
		
	// Test ability to upload a Dockerfile.
	{
		testContext.TryAddDockerfile(johnsRepoId, dockerfilePath, "A fine dockerfile")
	}
	
	// Test ability to list the Dockerfiles in a repo.
	{
		var dockerfileNames []string = testContext.TryGetDockerfiles(johnsRepoId)
		testContext.AssertThat(len(dockerfileNames) == 1, "Wrong number of dockerfiles")
	}
	
	// Test ability create a repo and upload a dockerfile at the same time.
	{
		var zippysRepoId string = testContext.TryCreateRepo(realm4Id, "zippysrepo",
			"A super smart repo", "dockerfile")
		var dockerfileNames []string = testContext.TryGetDockerfiles(zippysRepoId)
		testContext.AssertThat(len(dockerfileNames) == 1, "Wrong number of dockerfiles")
	}
	
	// Test ability to list the repos in a realm.
	{
		var repoIds []string = testContext.TryGetRealmRepos(realm4Id, true)
		testContext.AssertThat(len(repoIds) == 2, "Number of repo Ids returned was " +
			string(len(repoIds)) + ", expected 2")
	}
	
	// Test ability to define a Flag and then retrieve info about it.
	{
		var responseMap = testContext.TryDefineFlag(
			johnsRepoId, "myflag", "A really boss flag", flag2ImagePath)
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
				
				var fId string = testContext.TryGetFlagDescByName(johnsRepoId, "myflag")
				testContext.AssertThat(fId == flagId, "Flag not found by name")
			}
		}
	}

	// Test ability to define a scan config and then get info about it.
	{
		testContext.TryGetScanProviders()
		var config1Id string = testContext.TryDefineScanConfig("My Config 1",
			"A very find config", johnsRepoId, "clair", "", flagImagePath, []string{}, []string{})
		
		var responseMap = testContext.TryGetScanConfigDesc(config1Id, true)
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
				fileInfo, err = os.Stat(flagImagePath)
				if testContext.AssertErrIsNil(err, "") {
					testContext.AssertThat(fileInfo.Size() == size, "File has wrong size")
				}
			}
		}
		
		var configId string = testContext.TryGetScanConfigDescByName(johnsRepoId, "My Config 1")
		testContext.AssertThat(configId == config1Id, "Did not find scan config")
	}
}
	
/*******************************************************************************
 * Test ability to create groups, and use them.
 * Creates/uses the following:
 */
func TestCreateGroups(testContext *utils.TestContext) {
	
	fmt.Println("\nTest suite TestCreateGroups------------------\n")

	defer testContext.TryClearAll()
	
	// -------------------------------------
	// Test setup:
	// Create a realm and an admin user for the realm, and then log in as that user.
	// Create some users to add to groups.
	//
	
	var realm4Id string
	//var user4Id string
	var johnConnorUserId = "jconnor"
	var johnConnorPswd = "Cameron loves me"
	var johnConnorUserObjId string
	var sarahConnorUserId = "sconnor"
	var sarahConnorPswd = "pancakes"
	var sarahConnorUserObjId string
	
	{
		realm4Id, _, _ = testContext.TryCreateRealmAnon(
			"realm4", "realm 4 Org", "realm4admin", "realm 4 Admin Full Name",
			"realm4admin@gmail.com", "realm4adminpswd")
		
		testContext.TryAuthenticate("realm4admin", "realm4adminpswd", true)

		johnConnorUserObjId, _ = testContext.TryCreateUser(johnConnorUserId, "John Connor",
			"johnc@gmail.com", johnConnorPswd, realm4Id)

		sarahConnorUserObjId, _ = testContext.TryCreateUser(sarahConnorUserId, "Sarah Connor",
			"sarahc@gmail.com", sarahConnorPswd, realm4Id)
	}
	
	// -------------------------------------
	// Tests
	//
	
	var group1Id string
	
	// Test ability to create a group.
	{
		group1Id = testContext.TryCreateGroup(realm4Id, "mygroup",
			"For Overthrowning Skynet", false)
	}
	
	// Test ability to retrieve info about a group.
	{
		testContext.TryGetGroupDesc(group1Id)
	}
	
	// Test ability to add users to the group.
	{
		testContext.TryAddGroupUser(group1Id, johnConnorUserObjId)
		testContext.TryAddGroupUser(group1Id, sarahConnorUserObjId)
	}
	
	// Test ability to retrieve the users of a group.
	{
		var myGroupUsers []string = testContext.TryGetGroupUsers(group1Id)
		testContext.AssertThat(len(myGroupUsers) == 2, "Wrong number of group users")
	}
	
	// Test ability to retrieve the groups in a realm.
	{
		var realm4IdGroupIds []string = testContext.TryGetRealmGroups(realm4Id)
		testContext.AssertThat(len(realm4IdGroupIds) == 1, "Wrong number of realm groups")
	}
	
	// Test ability to remove a user from a group.
	{
		testContext.TryRemGroupUser(group1Id, sarahConnorUserObjId)
		var userIdsAfterRemoval []string = testContext.TryGetGroupUsers(group1Id)
		testContext.AssertThat(len(userIdsAfterRemoval) == 1, "Wrong number of users")
	}
	
	// Test ability of a user to to retrieve the user's groups.
	{
		testContext.TryAuthenticate(johnConnorUserId, johnConnorPswd, true)
		var myGroupIds []string = testContext.TryGetMyGroups()
		testContext.AssertThat(len(myGroupIds) == 1, "Wrong number of groups")
	}
}

/*******************************************************************************
 * Test the getMy... functions.
 * Creates/uses the following:
 */
func TestGetMy(testContext *utils.TestContext) {
		
	fmt.Println("\nTest suite TestGetMy------------------\n")

	defer testContext.TryClearAll()
	
	// -------------------------------------
	// Test setup:
	// 1. Create a realm X and an admin user for the realm, and then log in as that user.
	// 2. Create a non-admin user in realm X.
	// 3. Create a second realm Y and give the non-admin user access to it.
	// 4. Create a third realm Z, and a repo within that realm, and give the user access
	// to the repo.
	// 5.a. Write a dockerfile to a new temp directory.
	// 5.b. Create a second repo within the above realm, create a dockerfile within the repo,
	// and give the user access to that dockerfile.
	// 6. Same as above, for for a scan config and a flag.
	//
	
	var realmXId string
	var realmXAdminUserId = "realm4admin"
	var realmXAdminPswd = "Realm4Pswd"
	//var realmXAdminObjId string
	var realmXJohnUserId = "jconnor"
	var realmXJohnPswd = "ILoveCameron"
	var realmXJohnObjId string
	var realmYId string
	var realmZId string
	//var realmZRepo1Id string
	var realmZRepo2Id string
	var dockerfilePath string
	var realmZRepo2DockerfileId string
	var realmZRepo2ScanConfigId string
	var realmZRepo2FlagId string
	var flagImagePath = "Seal.png"
	
	{
		realmXId, _, _ = testContext.TryCreateRealmAnon(
			"realm4", "realm 4 Org", realmXAdminUserId, "realm 4 Admin Full Name",
			"realm4admin@gmail.com", realmXAdminPswd)
		
		testContext.TryAuthenticate(realmXAdminUserId, realmXAdminPswd, true)
		
		realmXJohnObjId, _ = testContext.TryCreateUser(realmXJohnUserId, "John Connor",
			"johnc@gmail.com", realmXJohnPswd, realmXId)
		
		realmYId = testContext.TryCreateRealm("sarahrealm", "Sarahs_Realm", "Escape into here")
		// Give john access:
		var permissions = []bool{true, false, false, false, false}
		testContext.TryAddPermission(realmXJohnObjId, realmYId, permissions)
		
		realmZId = testContext.TryCreateRealm("cromardirealm", "Cromardis_Realm", "Beware in here")
		testContext.TryCreateRepo(realmZId, "repo1", "A first repo", "")
		
		var err error
		dockerfilePath, err = utils.CreateTempFile("Dockerfile", "FROM centos\nRUN echo moo > oink")
		if err != nil { testContext.AbortAllTests(err.Error()) }
		defer os.Remove(dockerfilePath)
		
		realmZRepo2Id = testContext.TryCreateRepo(realmZId, "repo2", "Repo in realm z", "")
		testContext.TryAddPermission(realmXJohnObjId, realmZRepo2Id, permissions)
		
		realmZRepo2DockerfileId = testContext.TryAddDockerfile(realmZRepo2Id, dockerfilePath,
			"A dockerfile")
		testContext.TryAddPermission(realmXJohnObjId, realmZRepo2DockerfileId, permissions)
		
		realmZRepo2ScanConfigId = testContext.TryDefineScanConfig("Security Scan",
			"Show that scans passed", realmZRepo2Id, "clair", "", "", nil, nil)
		testContext.TryAddPermission(realmXJohnObjId, realmZRepo2ScanConfigId, permissions)
		
		err = utils.DownloadFile(SealURL, flagImagePath, true)
		if err != nil { testContext.AbortAllTests(err.Error()) }
		
		var responseMap = testContext.TryDefineFlag(realmZRepo2Id, "SuperSuccessFlag",
			"Show much better", flagImagePath)
		realmZRepo2FlagId = responseMap["FlagId"].(string)
		testContext.TryAddPermission(realmXJohnObjId, realmZRepo2FlagId, permissions)
	}
	
	// -------------------------------------
	// Tests
	//
	
	// Test ability of a user to retrieve information about the user's account.
	{
		var myAdminRealms []interface{}

		testContext.TryAuthenticate(realmXAdminUserId, realmXAdminPswd, true)
		_, myAdminRealms = testContext.TryGetMyDesc(true)
		testContext.AssertThat(len(myAdminRealms) == 3, "Wrong number of admin realms")
		
		testContext.TryAuthenticate(realmXJohnUserId, realmXJohnPswd, true)
		_, myAdminRealms = testContext.TryGetMyDesc(true)
		testContext.AssertThat(len(myAdminRealms) == 0, "Wrong number of admin realms")
	}
		
	testContext.TryAuthenticate(realmXJohnUserId, realmXJohnPswd, true)
	
	// Test ability of a user to to retrieve the user's realms.
	{
		var realmIds []string = testContext.TryGetMyRealms()
		testContext.AssertThat(len(realmIds) == 1, "Wrong number of realms")
	}
	
	// Test ability of a user to to retrieve the user's repos.
	{
		var myRepos []string = testContext.TryGetMyRepos()
		testContext.AssertThat(len(myRepos) == 1, fmt.Sprintf(
			"Only returned %d repos", len(myRepos)))
	}
	
	// Test ability of a user to to retrieve the user's dockerfiles.
	{
		var myDockerfileIds []string = testContext.TryGetMyDockerfiles()
		testContext.AssertThat(len(myDockerfileIds) == 1, "Wrong number of dockerfiles")
	}
	
	// Test ability of a user to to retrieve the user's scan configs.
	{
		var configIds []string
		_, configIds = testContext.TryGetMyScanConfigs()
		testContext.AssertThat(utils.ContainsString(configIds, realmZRepo2ScanConfigId),
			"Scan config not found")
	}
}

/*******************************************************************************
 * Test access control.
 * Creates/uses the following:
 */
func TestAccessControl(testContext *utils.TestContext) {
	
	fmt.Println("\nTest suite TestAccessControl------------------\n")

	defer testContext.TryClearAll()
	
	// -------------------------------------
	// Test setup:
	// 1. Create a realm X and an admin user for the realm, and then log in as that user.
	// 2. Create a non-admin user in realm X.
	// 3. Create a repo.
	// 4.a. Write a dockerfile to a new temp directory.
	// 4.b. Create a dockerfile within the repo,
	//
	
	var realmXId string
	var realmXAdminUserId = "realmXadmin"
	var realmXAdminPswd = "fluffy"
	//var realmXAdminObjId string
	var realmXJohnUserId = "jconnor"
	var realmXJohnPswd = "I am never safe"
	var realmXJohnObjId string
	var realmXRepo1Id string
	var dockerfileId string
	var dockerfilePath string
	
	{
		realmXId, _, _ = testContext.TryCreateRealmAnon(
			"realm4", "realm 4 Org", realmXAdminUserId, "realm 4 Admin Full Name",
			"realm4admin@gmail.com", realmXAdminPswd)
		
		testContext.TryAuthenticate(realmXAdminUserId, realmXAdminPswd, true)
		
		realmXJohnObjId, _ = testContext.TryCreateUser(realmXJohnUserId, "John Connor",
			"johnc@gmail.com", realmXJohnPswd, realmXId)
		
		realmXRepo1Id = testContext.TryCreateRepo(realmXId, "repo1", "Repo in realm x", "")
		
		var err error
		dockerfilePath, err = utils.CreateTempFile("Dockerfile", "FROM centos\nRUN echo moo > oink")
		if err != nil { testContext.AbortAllTests(err.Error()) }
		defer os.Remove(dockerfilePath)
		
		dockerfileId = testContext.TryAddDockerfile(realmXRepo1Id, dockerfilePath,
			"A first dockerfile")
	}
	
	// -------------------------------------
	// Tests
	//
	
	// Test ability to set permission.
	
	var perms1 []bool = []bool{false, true, false, true, true}
	
	{
		var retPerms1 []bool = testContext.TrySetPermission(realmXJohnObjId, dockerfileId, perms1)
		var expectedPerms1 []bool = []bool{false, true, false, true, true}
		for i, p := range retPerms1 {
			testContext.AssertThat(p == expectedPerms1[i], "Returned permission does not match")
		}
	}
	
	// Test ability to get permission.
	{
		var perms2 []bool = testContext.TryGetPermission(realmXJohnObjId, dockerfileId)
		if perms2 != nil {
			for i, p := range perms1 {
				testContext.AssertThat(p == perms2[i], "Returned permission does not match")
			}
		}
	}
		
	// Test ability to add permission.
	{
		var perms3 []bool = []bool{false, false, true, true, true}
		var retPerms3 []bool = testContext.TryAddPermission(realmXJohnObjId, dockerfileId, perms3)
		if retPerms3 != nil {
			var expectedPerms3 []bool = []bool{false, true, true, true, true}
			for i, p := range retPerms3 {
				testContext.AssertThat(p == expectedPerms3[i], "Returned permission does not match")
			}
		}
	}
	
	// Test ability to remove permission.
	{
		if testContext.TryRemPermission(realmXJohnObjId, dockerfileId) {
			var retPerms4 []bool = testContext.TryGetPermission(realmXJohnObjId, dockerfileId)
			var expectedPerms4 []bool = []bool{false, false, false, false, false}
			for i, p := range retPerms4 {
				fmt.Println(fmt.Sprintf("\tret perm[%d]: %#v; exp perm[%d]: %#v", i, p, i, expectedPerms4[i]))
				testContext.AssertThat(p == expectedPerms4[i], "Returned permission does not match")
			}
		}
	}
}

/*******************************************************************************
 * Test update/replace functions.
 * Creates/uses the following:
 */
func TestUpdateAndReplace(testContext *utils.TestContext) {
	
	fmt.Println("\nTest suite TestUpdateAndReplace------------------\n")

	defer testContext.TryClearAll()
	
	// -------------------------------------
	// Test setup:
	// 1. Create a realm and an admin user for the realm, and then log in as that user.
	// 2. Create a repo.
	// 3. Create a scan config.
	// 4. Create a non-admin user.
	// 5.a. Write a dockerfile to a new temp directory.
	// 5.b. Create a dockerfile within the repo,
	// 6. Create another realm.
	//
	
	var realmXId string
	var realmYId string
	var realmXYAdminUserId = "bigboss"
	var realmXYAdminPswd = "fluffy"
	//var realmXYAdminObjId string
	var realmXJohnUserId = "johnc"
	var realmXJohnPswd = "Ilovecam"
	var realmXJohnObjId string
	var realmXRepo1Id string
	var dockerfilePath string
	var dockerfileId string
	var scanConfigId string
	//var flagId string
	var flagImagePath = "Seal.png"
	var flag2ImagePath = "Seal2.png"
	
	{
		realmXId, _, _ = testContext.TryCreateRealmAnon(
			"realm4", "realm 4 Org", realmXYAdminUserId, "realm 4 Admin Full Name",
			"realm4admin@gmail.com", realmXYAdminPswd)
		
		testContext.TryAuthenticate(realmXYAdminUserId, realmXYAdminPswd, true)
		
		realmXRepo1Id = testContext.TryCreateRepo(realmXId, "repo1", "Repo in realm x", "")
		
		scanConfigId = testContext.TryDefineScanConfig("My Config 1",
			"A very find config", realmXRepo1Id, "clair", "", flagImagePath, []string{}, []string{})

		realmXJohnObjId, _ = testContext.TryCreateUser(realmXJohnUserId, "John Connor",
			"johnc@gmail.com", realmXJohnPswd, realmXId)
		
		var err error
		dockerfilePath, err = utils.CreateTempFile("Dockerfile", "FROM centos\nRUN echo moo > oink")
		if err != nil { testContext.AbortAllTests(err.Error()) }
		defer os.Remove(dockerfilePath)
		
		dockerfileId = testContext.TryAddDockerfile(realmXRepo1Id, dockerfilePath,
			"A first dockerfile")
		
		err = utils.DownloadFile(SealURL, flagImagePath, true)
		if err != nil { testContext.AbortAllTests(err.Error()) }
		
		err = utils.DownloadFile(Seal2URL, flag2ImagePath, true)
		if err != nil { testContext.AbortAllTests(err.Error()) }

		realmYId = testContext.TryCreateRealm(
			"realmq", "realm_q_org", "realm Q realm for fluffy things")
	}
	
	// -------------------------------------
	// Tests
	//
	
	// Test ability to replace a dockerfile.
	{
		dockerfileId = testContext.TryAddDockerfile(realmXRepo1Id, dockerfilePath, "A fine dockerfile")
		testContext.TryReplaceDockerfile(dockerfileId, "Dockerfile2", "The boo/ploink one")
	}
	
	// Test ability to substitute a scan config's flag with a different flag.
	{
		testContext.TryUpdateScanConfig(scanConfigId, "", "", "", "", flag2ImagePath,
			[]string{}, []string{})
		var scanConfig1Map map[string]interface{}
		scanConfig1Map = testContext.TryGetScanConfigDesc(scanConfigId, true)
		if testContext.CurrentTestPassed {
			// Id string
			// ProviderName string
			// SuccessExpression string
			// FlagId string
			// ParameterValueDescs []*ParameterValueDesc
			var newFlagId = scanConfig1Map["FlagId"]
			testContext.AssertThat(newFlagId != "", "FlagId returned empty")
		}
	}

	// Test ability to update one's own password.
	{
		testContext.TryAuthenticate(realmXJohnUserId, realmXJohnPswd, true)
		if testContext.TryChangePassword(realmXJohnUserId, realmXJohnPswd, "password2") {
			testContext.TryLogout()
			testContext.TryAuthenticate(realmXJohnUserId, realmXJohnPswd, false)
			testContext.TryAuthenticate(realmXJohnUserId, "password2", true)
		}
	}
	
	// Note: the password for realmXJohnUserId has now been changed.
	
	// Test ability to move a user from one realm to another.
	{
		testContext.TryAuthenticate(realmXYAdminUserId, realmXYAdminPswd, true)
		if testContext.TryMoveUserToRealm(realmXJohnObjId, realmYId) {
			// Verify that John is no longer in her realm.
			var responseMap = testContext.TryGetUserDesc(realmXJohnUserId)
			if testContext.CurrentTestPassed {
				// Verify that John is in realm Y.
				if ! testContext.AssertThat(responseMap["RealmId"] == realmYId,
					"Error: Realm move failed") {
					fmt.Println("Reponse map:")
					rest.PrintMap(responseMap)
				}
			}
		}
	}
}

/*******************************************************************************
 * Test deletion, diabling, etc.
 * Creates/uses the following:
 */
func TestDelete(testContext *utils.TestContext) {

	fmt.Println("\nTest suite TestDelete------------------\n")

	defer testContext.TryClearAll()
	
	// -------------------------------------
	// Test setup:
	// Create a realm and an admin user for the realm, and then log in as that user.
	// Create a repo.
	// Create a non-admin user.
	// Create a ScanConfig, with a flag.
	// Create a group.
	//
	
	var realmXId string
	var realmXAdminUserId = "bigcheese"
	var realmXAdminPswd = "I am a lumberjack"
	var realmXJohnUserId = "jconnor"
	var realmXJohnPswd = "bullets"
	var realmXJohnObjId string
	var realmXRepo1Id string
	var realmXScanConfigId string
	var realmXGroupId string
	var realmXFlagId string
	var flagImagePath = "Seal.png"

	{
		realmXId, _, _ = testContext.TryCreateRealmAnon(
			"realm4", "realm 4 Org", realmXAdminUserId, "realm 4 Admin Full Name",
			"realm4admin@gmail.com", realmXAdminPswd)
		
		testContext.TryAuthenticate(realmXAdminUserId, realmXAdminPswd, true)
		
		realmXRepo1Id = testContext.TryCreateRepo(realmXId, "repo1", "Repo in realm x", "")
		
		realmXJohnObjId, _ = testContext.TryCreateUser(realmXJohnUserId, "John Connor",
			"johnc@gmail.com", realmXJohnPswd, realmXId)
		
		realmXScanConfigId = testContext.TryDefineScanConfig("My Config 1",
			"A very fine config", realmXRepo1Id, "clair", "", flagImagePath, []string{}, []string{})

		realmXGroupId = testContext.TryCreateGroup(realmXId, "mygroup",
			"For Overthrowning Skynet", false)
		
		var err = utils.DownloadFile(SealURL, flagImagePath, true)
		if err != nil { testContext.AbortAllTests(err.Error()) }
	}
	
	// -------------------------------------
	// Tests
	//
	
	// Test ability to disable a user.
	{
		testContext.TryAuthenticate(realmXAdminUserId, realmXAdminPswd, true)
		if testContext.TryDisableUser(realmXJohnObjId) {
			// Now see if that user can authenticate - expect no.
			testContext.TryAuthenticate(realmXJohnUserId, realmXJohnPswd, false)
			if testContext.TryReenableUser(realmXJohnObjId) {
				testContext.TryAuthenticate(realmXJohnUserId, realmXJohnPswd, true)
			}
		}
	}
	
	// Test ability to delete a group.
	{
		testContext.TryAuthenticate(realmXAdminUserId, realmXAdminPswd, true)
		testContext.TryDeleteGroup(realmXGroupId)
	}
	
	// Test abilty to delete a scan config.
	{
		var responseMap = testContext.TryGetScanConfigDesc(realmXScanConfigId, true)
		var obj = responseMap["FlagId"]
		var isType bool
		realmXFlagId, isType = obj.(string)
		if ! isType { testContext.FailTest() } else {
			if testContext.TryRemScanConfig(realmXScanConfigId, true) {
				testContext.TryGetScanConfigDesc(realmXScanConfigId, false)
			}
		}
	}
	
	// Test ability to delete a flag.
	if realmXFlagId != "" {
		if testContext.TryRemFlag(realmXFlagId) {
			testContext.TryGetFlagDesc(realmXFlagId, false)
		}
	}
	
	// Test ability to log out.
	{
		if testContext.TryLogout() {
			testContext.TryGetMyDesc(false)
		}
	}
	
	// Test ability to deactivate a realm.
	{
		testContext.TryAuthenticate(realmXAdminUserId, realmXAdminPswd, true)
		if testContext.TryDeactivateRealm(realmXId) {
			testContext.TryGetRealmRepos(realmXId, false)
		}
	}
}
	
/*******************************************************************************
 * Test docker functions.
 * Creates/uses the following:
 */
func TestDockerFunctions(testContext *utils.TestContext) {

	fmt.Println("\nTest suite TestDockerFunctions------------------\n")

	defer testContext.TryClearAll()
	
	// -------------------------------------
	// Test setup:
	// Create a realm and an admin user for the realm, and then log in as that user.
	// Create a repo.
	// Create a ScanConfig.
	// Write a dockerfile to a new temp directory.
	// Create a dockerfile within the repo.
	// Write another dockerfile to the temp directory.
	// Create a dockerfile object for the new file.

	var realmXId string
	var realmXAdminUserId = "admin"
	var realmXAdminPswd = "fluffy"
	var realmXAdminObjId string
	var realmXRepo1Id string
	var dockerImage1ObjId string
	var scanConfigId string
	var dockerfilePath string
	var dockerfileId string
	var dockerfile2Path string
	//var dockerfile2Id string
	var flagImagePath = "Seal.png"
	
	{
		realmXId, realmXAdminObjId, _ = testContext.TryCreateRealmAnon(
			"realm4", "realm 4 Org", realmXAdminUserId, "realm 4 Admin Full Name",
			"realm4admin@gmail.com", realmXAdminPswd)
		
		testContext.TryAuthenticate(realmXAdminUserId, realmXAdminPswd, true)
		
		realmXRepo1Id = testContext.TryCreateRepo(realmXId, "repo1", "Repo in realm x", "")
		
		scanConfigId = testContext.TryDefineScanConfig("My Config 1",
			"A very fine config", realmXRepo1Id, "clair", "", flagImagePath, []string{}, []string{})

		var err error
		dockerfilePath, err = utils.CreateTempFile("Dockerfile", "FROM centos\nRUN echo moo > oink")
		if err != nil { testContext.AbortAllTests(err.Error()) }
		defer os.Remove(dockerfilePath)
		dockerfileId = testContext.TryAddDockerfile(realmXRepo1Id, dockerfilePath, "A fine dockerfile")
		
		dockerfile2Path, err = utils.CreateTempFile("Dockerfile2", "FROM centos\nRUN echo boo > ploink")
		if err != nil { testContext.AbortAllTests(err.Error()) }
		defer os.Remove(dockerfile2Path)
		testContext.TryAddDockerfile(realmXRepo1Id, dockerfile2Path, "A finer dockerfile")
		
		err = utils.DownloadFile(SealURL, flagImagePath, true)
		if err != nil { testContext.AbortAllTests(err.Error()) }
	}
	
	// -------------------------------------
	// Tests
	//
	
	// Test ability to build image from a dockerfile.
	{
		dockerImage1ObjId, _ = testContext.TryExecDockerfile(realmXRepo1Id, dockerfileId, "myimage")
	}
	
	// Test ability to list the images in a repo.
	{
		var imageNames []string = testContext.TryGetDockerImages(realmXRepo1Id)
		testContext.AssertThat(len(imageNames) == 1, "Wrong number of images")
	}
	
	// Test abilty to get the current logged in user's docker images.
	{
		var myDockerImageIds []string = testContext.TryGetMyDockerImages()
		testContext.AssertThat(len(myDockerImageIds) == 1, "Wrong number of docker images")
	}
	
	// Test ability to scan a docker image.
	{
		testContext.TryScanImage(scanConfigId, dockerImage1ObjId)
	}
	
	// Test ability to upload and exec a dockerfile in one command.
	{
		var dockerImage2ObjId string
		dockerImage2ObjId, _ = testContext.TryAddAndExecDockerfile(realmXRepo1Id,
			"My second image", "myimage2", dockerfile2Path)
		fmt.Println(dockerImage2ObjId)
	
		/*
		testContext.TryDownloadImage(dockerImage2ObjId, "BooPloinkImage")
		var responseMap = testContext.TryGetDockerImageDesc(dockerImage2ObjId, true)
		if testContext.CurrentTestPassed {
			// Check image signature.
			var image2Signature []byte
			var err error
			image2Signature, err = utils.ComputeFileSignature("BooPloinkImage")
			if testContext.AssertErrIsNil(err, "Unable to compute signature") {
				var obj interface{} = responseMap["Signature"]
				var sig, isType = obj.([]interface{})
				if testContext.AssertThat(isType, "Wrong type: " + reflect.TypeOf(sig).String()) {
					for i, sigi := range sig {
						var b = uint8(sigi.(float64))
						if ! testContext.AssertThat(
							b == image2Signature[i],
							fmt.Sprintf("Wrong signature: %d != %d", b, image2Signature[i])) { break }
					}
				}
			}
		}
		*/
	}
	
	// Test ability of a user to to retrieve the user's docker images.
	{
		var imageIds []string = testContext.TryGetMyDockerImages()
		testContext.AssertThat(len(imageIds) == 2, "Wrong number of docker images")
	}

	// Test ability to get the events for a specified user, including docker build events.
	{
		var eventIds []string = testContext.TryGetUserEvents(realmXAdminObjId)
		testContext.AssertThat(len(eventIds) == 3, "Wrong number of user events")
			// Should be one scan event and two dockerfile exec events.
	}
	
	// Test ability to get the events for a specified docker image.
	{
		var eventIds []string = testContext.TryGetDockerImageEvents(dockerImage1ObjId)
		testContext.AssertThat(len(eventIds) == 1, "Wrong number of image events")
			// Should be one scan event.
	}
	
	// Test ability to get the scan status of a docker image.
	{
		var responseMap = testContext.TryGetDockerImageStatus(dockerImage1ObjId)
		if testContext.CurrentTestPassed {
			testContext.AssertThat(responseMap["EventId"] != "", "No image status")
			testContext.AssertThat(responseMap["ScanConfigId"] == scanConfigId,
				"Wrong scan config Id")
			testContext.AssertThat(responseMap["ProviderName"] == "clair",
				"Wrong provider")
		}
	}
	
	// Test ability to get the events for a specified docker file.
	{
		var eventIds []string = testContext.TryGetDockerfileEvents(dockerfileId)
		testContext.AssertThat(len(eventIds) == 1, "Wrong number of image events")
			// Should be one dockerfile exec event.
	}
	
	// Test abilit to delete a specified docker image.
	{
		if testContext.TryRemDockerImage(dockerImage1ObjId) {
			testContext.TryGetDockerImageDesc(dockerImage1ObjId, false)
		}
	}
}
