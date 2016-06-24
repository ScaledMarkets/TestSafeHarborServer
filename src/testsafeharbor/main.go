/*******************************************************************************
 * Perform independent end-to-end ("behavioral") tests on the SafeHarbor server.
 */

package main

import (
	"fmt"
	//"net/http"
	"net/url"
	"os"
	"flag"
	"time"
	"strings"
	"reflect"
	"strconv"
	
	"redis"
	"goredis"
	
	// SafeHarbor packages:
	"testsafeharbor/docker"
	"testsafeharbor/utils"
	"testsafeharbor/rest"
)

const (
	SealURL = "https://itsonlywords55.files.wordpress.com/2010/01/seal-of-approval.jpg"
	Seal2URL = "http://thumb10.shutterstock.com/display_pic_with_logo/681547/140365213/stock-photo-seal-of-approval-quality-check-grunge-vector-on-white-background-this-graphic-illustration-140365213.jpg"
)

func main() {
	
	var testSuite = map[string]func(*utils.TestContext) {
		"Email": TestEmail,
		"DockSvcs": TestDockerServices,
		"Registry": TestDockerRegistry,
		"Engine": TestDockerEngine,
		"json": TestJSONDeserialization,
		"goredis": TestGoRedis,
		"redis": TestRedis,
		"CreateRealmsAndUsers": TestCreateRealmsAndUsers,
		"CreateResources": TestCreateResources,
		"OptionalParams": TestOptionalParams,
		"CreateGroups": TestCreateGroups,
		"ScanConfigs": TestScanConfigs,
		"GetMy": TestGetMy,
		"AccessControl": TestAccessControl,
		"EmailVerificationStep1": TestEmailIdentityVerificationStep1, 
		"EmailVerificationStep2": TestEmailIdentityVerificationStep2, 
		"UpdateAndReplace": TestUpdateAndReplace,
		"Delete": TestDelete,
		"DockerFunctions": TestDockerFunctions,
	}

	var help *bool = flag.Bool("help", false, "Provide help instructions.")
	var scheme *string = flag.String("s", "http", "Protocol scheme (one of http, https, unix)")
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
	fmt.Println("tests: " + *tests)
	fmt.Println("Test suites that will be run:")
	for _, testName := range testsToRun {
		var testFunction = testSuite[testName]
		if testFunction == nil {
			fmt.Println("Test '" + testName + "' not recognized")
			os.Exit(0)
		}
		testFunctionsToRun = append(testFunctionsToRun, testFunction)
		fmt.Println("\t" + testName)
	}
	
	// Prepare to run tests.
	var testContext = utils.NewTestContext(*scheme, *hostname, *port, utils.SetSessionId,
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
func TestEmail(testContext *utils.TestContext) {
	
	fmt.Println("\nTest suite TestEmail------------------\n")

	// -------------------------------------
	// Test setup:
	
	var emailService *utils.EmailService
	var err error
	var emailConfigMap = map[string]interface{}{
		"SES_SMTP_hostname": "email-smtp.us-west-2.amazonaws.com",
		"SES_SMTP_Port": 25.0,
		"SenderAddress": "cliff_test@cliffberg.com",
		"SenderUserId": "AKIAI2FOYVEKGEZXKX6A",
		"SenderPassword": "Amcjxs1E9+mFH06zM38SoyeOMfmG5sy77OC3y6ifhSJ3",
	}
	
	{
		fmt.Println("Creating EmailService...")
		emailService, err = utils.CreateEmailService(emailConfigMap)
		testContext.AssertErrIsNil(err, "When instantiating email service")
	}
	
	// Tests
	{
		testContext.StartTest("Calling SendEmail...")
		fmt.Println("Sending message...")
		err = emailService.SendEmail("cliff_cromarti@cliffberg.com",
			"testing email service", "This is a test of the email service")
		fmt.Println("...message sent.")
		testContext.AssertErrIsNil(err, "When calling SendMail")
	}
	
	fmt.Println("Done")
}

/*******************************************************************************
 * 
 */
func TestDockerServices(testContext *utils.TestContext) {
	
	fmt.Println("\nTest suite TestDockerServices------------------\n")

	// -------------------------------------
	// Test setup:
	
	{
		//....
	}
	
	// Test BuildDockerfile
	{
		//....
	}
}

/*******************************************************************************
 * 
 */
func TestDockerEngine(testContext *utils.TestContext) {
	
	fmt.Println("\nTest suite TestDockerEngine------------------\n")

	// -------------------------------------
	// Test setup:
	
	var engine docker.DockerEngine
	var err error
	var buildDirPath string
	var imageFullName = "testimage:5"
	var dockerfileContent = "FROM centos\nRUN touch newfile"

	var registryHost = os.Getenv("RegistryHost")
	var registryPort int
	registryPort, err = strconv.Atoi(os.Getenv("RegistryPort"))
	if err != nil { testContext.AbortAllTests(err.Error()) }
	var registryRepo = "greatimage"
	var registryUserId = os.Getenv("registryUser")
	var registryPassword = os.Getenv("registryPassword")

	var tag = "alpha"
	
	{
		// Create a build directory.
		buildDirPath, err = utils.CreateTempDir()
		if err != nil { testContext.AbortAllTests(err.Error()) }
		
		// Create a dockerfile.
		_, err = utils.CreateTempFile(buildDirPath, "Dockerfile", dockerfileContent)
		if err != nil { testContext.AbortAllTests(err.Error()) }
	}
	
	// Test connecting to Engine.
	{
		testContext.StartTest("Open Engine connection...")
		engine, err = docker.OpenDockerEngineConnection()
		testContext.AssertErrIsNil(err, "In opening connection to docker engine")
		testContext.PassTestIfNoFailures()
	}
	
	// Test GetImages().
	{
		testContext.StartTest("GetImages")
		var imageMaps []map[string]interface{}
		imageMaps, err = engine.GetImages()
		testContext.AssertErrIsNil(err, "In getting images")

		// debug
		fmt.Println("Images:")  
		for _, imageMap := range imageMaps { rest.PrintMap(imageMap) }
		// end debug
		
		testContext.PassTestIfNoFailures()
	}
	
	// Test BuildImage.
	{
		testContext.StartTest("BuildImage")
		fmt.Println("Building image '" + imageFullName + "' in directory '" + buildDirPath + "'")
		var responseStr string
		responseStr, err = engine.BuildImage(buildDirPath, imageFullName, "Dockerfile")
		testContext.AssertErrIsNil(err, "In building image")
		fmt.Println("Response from BuildImage:")
		fmt.Println(responseStr)
		fmt.Println()
		
		// Check that the image was actually created.
		fmt.Println("Images:")  // debug
		var imageMaps []map[string]interface{}
		imageMaps, err = engine.GetImages()
		var found bool = false
		for _, imageMap := range imageMaps {
			var obj interface{} = imageMap["RepoTags"]
			var isType bool
			var tags []interface{}
			tags, isType = obj.([]interface{})
			if ! testContext.AssertThat(isType, "RepoTags is not an interface array") { break }
			if ! testContext.AssertThat(tags != nil, "No RepoTags found") { break }
			for _, tagi := range tags {
				var tag string
				tag, isType = tagi.(string)
				if ! testContext.AssertThat(isType,
					"tag is not a string - it is a " + reflect.TypeOf(tagi).String()) { break }
				if tag == imageFullName {
					found = true
					break
				}
			}
			if testContext.TestHasFailed() { break }
			if found { break }
		}
		testContext.AssertThat(found, "Image not found")
		
		testContext.PassTestIfNoFailures()
	}
	
	// Test TagImage.
	{
		testContext.StartTest("TagImage")
		
		var regHostAndRepoName = fmt.Sprintf("%s:%d/%s", registryHost, registryPort, registryRepo)
		err = engine.TagImage(imageFullName, regHostAndRepoName, tag)
		testContext.AssertErrIsNil(err, "In tagging image " + imageFullName)
		
		// Verify that the engine now contains an image with the host/repo:tag name.
		_, err = engine.GetImageInfo(regHostAndRepoName + ":" + tag)
		testContext.AssertErrIsNil(err, "In getting image")
		
		testContext.PassTestIfNoFailures()
	}
	
	// Test PushImage.
	{
		testContext.StartTest("PushImage")
		
		var regHostAndRepoName = fmt.Sprintf("%s:%d/%s", registryHost, registryPort, registryRepo)
		err = engine.PushImage(regHostAndRepoName, tag,
			registryUserId, registryPassword, "noone@nowhere.com")
		testContext.AssertErrIsNil(err, "In pushing image")
		fmt.Println("Image pushed")
		
		// Verify that the registry now contains an image with the full name.
		fmt.Println("Now verifying that the image is in the registry")
		var registry docker.DockerRegistry
		registry, err = docker.OpenDockerRegistryConnection(registryHost, registryPort,
			registryUserId, registryPassword)
		testContext.AssertErrIsNil(err, "In opening connection to docker registry")
		var exists bool
		exists, err = registry.ImageExists(registryRepo, tag)
		testContext.AssertErrIsNil(err, "While calling ImageExists")
		testContext.AssertThat(exists, "Did not find image")
		
		testContext.PassTestIfNoFailures()
	}
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
	var registryHost = os.Getenv("RegistryHost")
	var registryPort int
	var err error
	registryPort, err = strconv.Atoi(os.Getenv("RegistryPort"))
	var registryUserId = os.Getenv("registryUser")
	var registryPassword = os.Getenv("registryPassword")
	var testImageRepoName = os.Getenv("TestImageRepoName")
	var testImageTag = os.Getenv("TestImageTag")
	var imageToUploadPath = os.Getenv("ImageToUploadPath")
	var imageToUploadDigest = os.Getenv("ImageToUploadDigest")
	var downloadedImageFilePath = "DownloadedImage.tar"
	
	var registry docker.DockerRegistry
	
	{
		testContext.StartTest("Initialization")
		testContext.AssertThat(registryUserId != "", "registryUserId is empty")
		testContext.AssertThat(registryPassword != "", "registryPassword is empty")
		testContext.AssertThat(testImageRepoName != "", "TestImageRepoName is empty")
		testContext.AssertThat(testImageTag != "", "TestImageTag is empty")
		testContext.AssertThat(imageToUploadPath != "", "ImageToUploadPath is empty")
		testContext.AssertThat(imageToUploadDigest != "", "ImageToUploadDigest is empty")
	}
	
	// Test connecting to Registry.
	{
		testContext.StartTest(fmt.Sprintf(
			"Open Registry connection, using %s:%s...", registryUserId, registryPassword))
		registry, err = docker.OpenDockerRegistryConnection(registryHost, registryPort,
			registryUserId, registryPassword)
		testContext.AssertErrIsNil(err, "In opening connection to docker registry")
		testContext.PassTestIfNoFailures()
	}
	
	// Test PushImage.
	{
		testContext.StartTest("PushImage")
		
		err = registry.PushImage(testImageRepoName, testImageTag, imageToUploadPath)
		testContext.AssertErrIsNil(err, "While calling PushImage")
		
		testContext.PassTestIfNoFailures()
	}
	
	// Test ImageExists.
	{
		testContext.StartTest("ImageExists")
		
		var exists bool
		exists, err = registry.ImageExists(testImageRepoName, testImageTag)
		testContext.AssertErrIsNil(err, "While calling ImageExists")
		testContext.AssertThat(exists, "Did not find image")
		
		testContext.PassTestIfNoFailures()
	}
	
	// Test GetImageInfo.
	{
		testContext.StartTest("GetImageInfo")
		
		var testDigest string
		var layerAr []map[string]interface{}
		testDigest, layerAr, err = registry.GetImageInfo(testImageRepoName, testImageTag)
		testContext.AssertThat(testDigest == imageToUploadDigest,
			"Incorrect digest for " + testImageRepoName + ":" + testImageTag)
		testContext.AssertThat(len(layerAr) > 0, "No layer descriptions")
		
		testContext.PassTestIfNoFailures()
	}
	
	// Test GetImage.
	{
		testContext.StartTest("GetImage")
		
		os.Remove(downloadedImageFilePath)  // ignore error, if any.
		
		// Contact Registry to get image.
		err = registry.GetImage(testImageRepoName, testImageTag, downloadedImageFilePath)
		
		// Verify that the image was retrieved properly.
		testContext.AssertErrIsNil(err, "While calling GetImage")
		var downloadedImageFile *os.File
		downloadedImageFile, err = os.OpenFile(downloadedImageFilePath, os.O_WRONLY, 0600)
		testContext.AssertErrIsNil(err, fmt.Sprintf(
			"When opening image file '%s'", downloadedImageFilePath))
		var fileInfo os.FileInfo
		fileInfo, err = downloadedImageFile.Stat()
		testContext.AssertErrIsNil(err, fmt.Sprintf(
			"When getting status of image file '%s'", downloadedImageFilePath))
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
	//var highTrustClientUserId = "HighTrustClient"
	//var highTrustClientPswd = "trustme"
	
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
		testContext.AssertThat(len(user4AdminRealms) == 1,
			fmt.Sprintf("Wrong number of admin realms: %d", len(user4AdminRealms)))
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
		testContext.TryCreateRealm("my2ndrealm", "A Big Company",
			"A second realm for a really big company")
	}
	
	// Test ability to look up a realm by its name.
	{
		testContext.TestGetRealmByName("my2ndrealm")
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
	
	// Test ability to check if a user exists.
	{
		//testContext.TryAuthenticate(highTrustClientUserId, highTrustClientPswd, true)
		testContext.TryUserExists(true, joeUserId)
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
	var tempdir string
	
	{
		var err error
		tempdir, err = utils.CreateTempDir()
		if err != nil { testContext.AbortAllTests(err.Error()) }
		
		realm4Id, _, _ = testContext.TryCreateRealmAnon(
			"realm4", "realm 4 Org", "realm4admin", "realm 4 Admin Full Name",
			"realm4admin@gmail.com", "realm4adminpswd")
		
		testContext.TryAuthenticate("realm4admin", "realm4adminpswd", true)
		
		dockerfilePath, err = utils.CreateTempFile(tempdir, "Dockerfile", "FROM centos\nRUN echo moo > oink")
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
 * Test the ability to not specify parameters that are optional:
 * RepoId may be omitted for addDockerfile, addAndExecDockerfile, defineScanConfig,
 * and defineFlag.
 * ImageName and Desc may be omitted for ....
 */
func TestOptionalParams(testContext *utils.TestContext) {
	
	fmt.Println("\nTest suite TestOptionalParams------------------\n")

	defer testContext.TryClearAll()
	
	// -------------------------------------
	// Test setup:
	
	....
	// -------------------------------------
	// Tests
	//
	.....
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
	var tempdir string
	{
		var err error
		tempdir, err = utils.CreateTempDir()
		if err != nil { testContext.AbortAllTests(err.Error()) }
		
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
		
		dockerfilePath, err = utils.CreateTempFile(tempdir, "Dockerfile", "FROM centos\nRUN echo moo > oink")
		if err != nil { testContext.AbortAllTests(err.Error()) }
		defer os.Remove(dockerfilePath)
		
		realmZRepo2Id = testContext.TryCreateRepo(realmZId, "repo2", "Repo in realm z", "")
		testContext.TryAddPermission(realmXJohnObjId, realmZRepo2Id, permissions)
		
		realmZRepo2DockerfileId, _ = testContext.TryAddDockerfile(realmZRepo2Id, dockerfilePath,
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
	var tempdir string
	
	{
		var err error
		tempdir, err = utils.CreateTempDir()
		if err != nil { testContext.AbortAllTests(err.Error()) }
		
		realmXId, _, _ = testContext.TryCreateRealmAnon(
			"realm4", "realm 4 Org", realmXAdminUserId, "realm 4 Admin Full Name",
			"realm4admin@gmail.com", realmXAdminPswd)
		
		testContext.TryAuthenticate(realmXAdminUserId, realmXAdminPswd, true)
		
		realmXJohnObjId, _ = testContext.TryCreateUser(realmXJohnUserId, "John Connor",
			"johnc@gmail.com", realmXJohnPswd, realmXId)
		
		realmXRepo1Id = testContext.TryCreateRepo(realmXId, "repo1", "Repo in realm x", "")
		
		dockerfilePath, err = utils.CreateTempFile(tempdir, "Dockerfile", "FROM centos\nRUN echo moo > oink")
		if err != nil { testContext.AbortAllTests(err.Error()) }
		defer os.Remove(dockerfilePath)
		
		dockerfileId, _ = testContext.TryAddDockerfile(realmXRepo1Id, dockerfilePath,
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
	
	// Test can one cannot modify the user info for another user.
	{
		testContext.TryUpdateUserInfo(false, realmXJohnUserId, "John Baum", "jbaum@gmail.com")
	}
	
	// Test that only the special user "HighTrustClient" can call userExists.
	{
		//testContext.TryUserExists(false, realmXAdminUserId)
	}
}

/*******************************************************************************
 * Test email based identity verification - step 1.
 */
var TestEmailIdentityVerificationStep1Explanation = `
Test email based identity verification.
This test must be completed manually, via the user checking email and
clicking on the link in the email. To run this test,
	1. Make sure that server has been started with the -toggleemail option and
		without the -noauthorization option.
	2. Run the test.
	3. Check the email account "cromarti_verifrealm@cliffberg.com". (If no email, then fail.)
	4. Click on the link in the email.
	5. Perform test TestEmailIdentityVerificationStep2. (Do not restart the server
		inbetween steps 1 and 2.)

`
func TestEmailIdentityVerificationStep1(testContext *utils.TestContext) {

	fmt.Println("\nTest suite TestEmailIdentityVerificationStep1------------------\n")
	fmt.Println(TestEmailIdentityVerificationStep1Explanation)

	// -------------------------------------
	// Test setup:

	var realmXId string
	var realmXAdminUserId = "realmXadmin"
	var realmXAdminPswd = "fluffy"
	
	{
		realmXId, _, _ = testContext.TryCreateRealmAnon(
			"verifrealm", "Email Verification Realm", realmXAdminUserId,
			"verifrealm Admin Full Name", "admin_verifrealm@cliffberg.com", realmXAdminPswd)
		
		testContext.TryAuthenticate(realmXAdminUserId, realmXAdminPswd, true)
		
		// Give 
	}
	
	// Test email based identity verification.
	{
		testContext.TryEnableEmailVerification(true)
		defer testContext.TryEnableEmailVerification(false)
		
		var userObjId string
		userObjId, _ = testContext.TryCreateUser("cromarti", "Cromarti",
			"cliff@cliffberg.com", "cromartiPswd", realmXId)
			//"cromarti_verifrealm@cliffberg.com", "cromartiPswd", realmXId)
		
		// Give the user permission to modify the realm.
		var perms []bool = []bool{true, true, true, true, true}
		var retPerms []bool = testContext.TryAddPermission(userObjId, realmXId, perms)
		if retPerms != nil {
			var expectedPerms []bool = []bool{true, true, true, true, true}
			for i, p := range retPerms {
				testContext.AssertThat(p == expectedPerms[i], "Returned permission does not match")
			}
		}
	}
}

/*******************************************************************************
 * Test email based identity verification - step 2.
 */
func TestEmailIdentityVerificationStep2(testContext *utils.TestContext) {
	
	fmt.Println("\nTest suite TestEmailIdentityVerificationStep2------------------\n")
	
	defer testContext.TryClearAll()
	
	// -------------------------------------
	// Test setup:

	var realmXAdminUserId = "realmXadmin"
	var realmXAdminPswd = "fluffy"
	var realmXId string

	{
		testContext.TryAuthenticate(realmXAdminUserId, realmXAdminPswd, true)

		// Identify the realm.
		var realmDescMap map[string]interface{}
		realmDescMap = testContext.TryGetRealmByName("verifrealm")
		var isType bool
		realmXId, isType = realmDescMap["Id"].(string)
		testContext.AssertThat(isType, "Id is not a string")
	}
	
	// Test that the user created by TestEmailIdentityVerificationStep1 can perform
	// actions that only a verified user can perform.
	{
		testContext.TryAuthenticate("cromarti", "cromartiPswd", true)
		
		var userDesc map[string]interface{}
		userDesc = testContext.TryGetUserDesc("cromarti")
		var obj interface{} = userDesc["EmailIsVerified"]
		var isVerified bool
		var isType bool
		isVerified, isType = obj.(bool)
		if testContext.AssertThat(isType, "Field EmailIsVerified is not bool") {
			testContext.AssertThat(isVerified, "EmailIsVerified is false")
		}
		
		testContext.TryCreateRepo(realmXId, "arepo",
			"a fine repo for email verification", "")
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
	var tempdir string
	
	{
		var err error
		tempdir, err = utils.CreateTempDir()
		if err != nil { testContext.AbortAllTests(err.Error()) }
		
		realmXId, _, _ = testContext.TryCreateRealmAnon(
			"realm4", "realm 4 Org", realmXYAdminUserId, "realm 4 Admin Full Name",
			"realm4admin@gmail.com", realmXYAdminPswd)
		
		testContext.TryAuthenticate(realmXYAdminUserId, realmXYAdminPswd, true)
		
		realmXRepo1Id = testContext.TryCreateRepo(realmXId, "repo1", "Repo in realm x", "")
		
		scanConfigId = testContext.TryDefineScanConfig("My Config 1",
			"A very find config", realmXRepo1Id, "clair", "", flagImagePath, []string{}, []string{})

		realmXJohnObjId, _ = testContext.TryCreateUser(realmXJohnUserId, "John Connor",
			"johnc@gmail.com", realmXJohnPswd, realmXId)
		
		dockerfilePath, err = utils.CreateTempFile(tempdir, "Dockerfile", "FROM centos\nRUN echo moo > oink")
		if err != nil { testContext.AbortAllTests(err.Error()) }
		defer os.Remove(dockerfilePath)
		
		dockerfileId, _ = testContext.TryAddDockerfile(realmXRepo1Id, dockerfilePath,
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
		//dockerfileId = testContext.TryAddDockerfile(realmXRepo1Id, dockerfilePath, "A fine dockerfile")
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
	
	// Test ability to update one's info.
	{
		testContext.TryUpdateUserInfo(true, realmXYAdminUserId, "realm 4 Admin",
			"realm4admin@mycorp.com")
	}
}

/*******************************************************************************
 * Test functions that link and unlink scan configs to DockerImages.
 */
func TestScanConfigs(testContext *utils.TestContext) {
	
	fmt.Println("\nTest suite TestScanConfigs------------------\n")

	defer testContext.TryClearAll()

	// -------------------------------------
	// Test setup:
	//
	var dockerImage1Id, dockerImage2Id string
	var scanConfigAId, scanConfigBId, scanConfigCId string
	
	{
		var realmId string
		var repoId string
		var dockerfile1Path, dockerfile2Path string
		var dockerfile1Id, dockerfile2Id string
		var mrscanneruserid string = "mrscanner"
		var mrscannerpswd string = "abc"
		
		realmId, _, _ = testContext.TryCreateRealmAnon(
			"securerealm", "SecureRealm Org", mrscanneruserid, "Mr. Scanner",
			"mrscanner@gmail.com", mrscannerpswd)
		
		testContext.TryAuthenticate(mrscanneruserid, mrscannerpswd, true)
		
		repoId = testContext.TryCreateRepo(realmId, "repo1", "Repo in SecureRealm", "")
		
		var tempdir string
		var err error
		tempdir, err = utils.CreateTempDir()
		if err != nil { testContext.AbortAllTests(err.Error()) }
		dockerfile1Path, err = utils.CreateTempFile(tempdir, "Dockerfile1", "FROM centos\nRUN echo goo > oink")
		if err != nil { testContext.AbortAllTests(err.Error()) }
		defer os.Remove(dockerfile1Path)
		dockerfile1Id, _ = testContext.TryAddDockerfile(repoId, dockerfile1Path, "A gooey dockerfile")
		
		dockerfile2Path, err = utils.CreateTempFile(tempdir, "Dockerfile2", "FROM centos\nRUN echo shoo > oink")
		if err != nil { testContext.AbortAllTests(err.Error()) }
		defer os.Remove(dockerfile2Path)
		dockerfile2Id, _ = testContext.TryAddDockerfile(repoId, dockerfile2Path, "A shooey dockerfile")
		
		// Create two docker images.
		_, dockerImage1Id = testContext.TryExecDockerfile(repoId,
			dockerfile1Id, "image1", []string{}, []string{})
		testContext.AssertThat(dockerImage1Id != "", "No image obj Id returned")
		
		_, dockerImage2Id = testContext.TryExecDockerfile(repoId,
			dockerfile2Id, "image2", []string{}, []string{})
		testContext.AssertThat(dockerImage1Id != "", "No image obj Id returned")
		
		// Create three scan configs.
		scanConfigAId = testContext.TryDefineScanConfig("Config A",
			"For scanning all images", repoId, "clair", "",
			"", []string{}, []string{})

		scanConfigBId = testContext.TryDefineScanConfig("Config B",
			"For scanning image 1", repoId, "clair", "",
			"", []string{}, []string{})
		
		scanConfigCId = testContext.TryDefineScanConfig("Config C",
			"For scanning image 2", repoId, "clair", "",
			"", []string{}, []string{})
	}
	
	// -------------------------------------
	// Tests
	//
	
	// Test that one can link DockerImages and ScanConfigs, and unlink them.
	{
		// Link Image1 to A and B.
		testContext.TryUseScanConfigForImage(dockerImage1Id, scanConfigAId)
		testContext.TryUseScanConfigForImage(dockerImage1Id, scanConfigBId)
		
		// Link Image2 to A and C.
		testContext.TryUseScanConfigForImage(dockerImage2Id, scanConfigAId)
		testContext.TryUseScanConfigForImage(dockerImage2Id, scanConfigCId)
		
		// Unlink C from Image2: now 1 uses A, B and 2 uses only A.
		testContext.TryStopUsingScanConfigForImage(dockerImage2Id, scanConfigCId)
		
		var scanConfigDescMap map[string]interface{}
		scanConfigDescMap = testContext.TryGetScanConfigDesc(scanConfigAId, true)  // should be 1, 2.
		var obj interface{} = scanConfigDescMap["DockerImagesIdsThatUse"]
		var objAr []interface{}
		var isType bool
		objAr, isType = obj.([]interface{})
		testContext.AssertThat(isType, "DockerImagesIdsThatUse is not an array")
		var ids []string = make([]string, len(objAr))
		for i, obj := range objAr {
			var isType bool
			ids[i], isType = obj.(string)
			testContext.AssertThat(isType, "DockerImagesIdsThatUse contains non-strings")
		}
		if testContext.AssertThat(len(ids) == 2, fmt.Sprintf(
			"Wrong number of image Ids returned: %d", len(ids))) {
		
			testContext.AssertThat(utils.Contains(dockerImage1Id, ids), "Id not in image Id list")
			testContext.AssertThat(utils.Contains(dockerImage2Id, ids), "Id not in image Id list")
		}
		
		scanConfigDescMap = testContext.TryGetScanConfigDesc(scanConfigBId, true)  // should be 1.
		obj = scanConfigDescMap["DockerImagesIdsThatUse"]
		objAr, isType = obj.([]interface{})
		testContext.AssertThat(isType, "DockerImagesIdsThatUse is not an array")
		ids = make([]string, len(objAr))
		for i, obj := range objAr {
			var isType bool
			ids[i], isType = obj.(string)
			testContext.AssertThat(isType, "DockerImagesIdsThatUse contains non-strings")
		}
		if testContext.AssertThat(len(ids) == 1, fmt.Sprintf(
			"Wrong number of image Ids returned: %d", len(ids))) {
		
			testContext.AssertThat(utils.Contains(dockerImage1Id, ids), "Id not in image Id list")
		}
	}
	
	// Test that one can omit specifying the ScanConfigId.
	{
		var scanEventDescs []map[string]interface{}
		scanEventDescs = testContext.TryScanImage("", dockerImage1Id)
		testContext.AssertThat(len(scanEventDescs) == 2, "Wrong number of scan events")
	}
	
	// Test that one can specify multiple ScanConfigIds.
	{
		var scanEventDescs []map[string]interface{}
		scanEventDescs = testContext.TryScanImage(
			url.QueryEscape(scanConfigAId + "," + scanConfigBId), dockerImage2Id)
		testContext.AssertThat(len(scanEventDescs) == 2, "Wrong number of scan events")
	}
}

/*******************************************************************************
 * Test deletion, diabling, etc.
 */
func TestDelete(testContext *utils.TestContext) {

	fmt.Println("\nTest suite TestDelete------------------\n")

	defer testContext.TryClearAll()
	
	// -------------------------------------
	// Test setup:
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
	var dockerfile1Path string
	var dockerfile1Id string
	var imageVersion1ObjId string
	var imageVersion2ObjId string
	var execEvent1Id string

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
		
		var tempdir string
		tempdir, err = utils.CreateTempDir()
		if err != nil { testContext.AbortAllTests(err.Error()) }
		dockerfile1Path, err = utils.CreateTempFile(
			tempdir, "Dockerfile", "FROM centos\nRUN echo moo > oink")
		if err != nil { testContext.AbortAllTests(err.Error()) }
		defer os.Remove(dockerfile1Path)

		imageVersion1ObjId, _, execEvent1Id = testContext.TryAddAndExecDockerfile(realmXRepo1Id,
			"My first image", "myimage1", dockerfile1Path, []string{}, []string{})
		testContext.AssertThat(imageVersion1ObjId != "", "Failed to create image")
		
		var event1Map map[string]interface{}
		event1Map = testContext.TryGetEventDesc(execEvent1Id)
		testContext.AssertThat(event1Map != nil, "Unable to get EventDesc for image")
		var obj = event1Map["DockerfileId"]
		testContext.AssertThat(obj != nil, "nil value for DockerfileId")
		var isType bool
		dockerfile1Id, isType = obj.(string)
		testContext.AssertThat(isType, "DockerfileId is not a string")

		imageVersion2ObjId, _ = testContext.TryExecDockerfile(realmXRepo1Id,
			dockerfile1Id, "myimage1", []string{}, []string{})
		testContext.AssertThat(imageVersion2ObjId != "", "Failed to create image")
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
	
	// Test ability to delete a docker image version
	{
		// Obtain the Id of image1.
		var image1ObjId string
		var dockerImageVersionDescMap map[string]interface{}
		dockerImageVersionDescMap = testContext.TryGetDockerImageDesc(imageVersion1ObjId, true)
		var obj interface{}
		obj = dockerImageVersionDescMap["ImageObjId"]
		var isType bool
		image1ObjId, isType = obj.(string)
		testContext.AssertThat(isType, "ImageObjId is not a string")
			
		// Delete an image version.
		if testContext.TryRemImageVersion(imageVersion1ObjId) {
			
			// Verify that the image now has the expected set of image versions.
			var eltFieldMaps = make([]map[string]interface{}, 0)
			eltFieldMaps = testContext.TryGetDockerImageVersions(image1ObjId)
			if testContext.AssertThat(eltFieldMaps != nil, "In TryGetDockerImageVersions") {
				
				/* Fields expected in each elt:
				HTTPStatusCode int
				HTTPReasonPhrase string
				ObjId string
				Version string
				ImageObjId string
				CreationDate string
				Digest []byte
				Signature []byte
				ScanEventIds []string
				DockerBuildOutput string
				*/
				
				testContext.AssertThat(len(eltFieldMaps) == 1,
					"TryGetDockerImageVersions: wrong number of elements returned")
			}
			
			// Verify that events had their image version references nullified.
			var eventIds []string = testContext.TryGetUserEvents(image1ObjId)
			var imageVersionEmptyCount = 0
			var imageVersion1ObjIdCount = 0
			var imageVersion2ObjIdCount = 0
			if ! testContext.TestHasFailed() {
				testContext.AssertThat(len(eventIds) == 2, "Wrong number of events")
				for _, eventId := range eventIds {
					
					var eventMap map[string]interface{}
					eventMap = testContext.TryGetEventDesc(eventId)
					if testContext.AssertThat(eventMap != nil, "Nil event map") {
						var obj = eventMap["ObjectType"]
						var isType bool
						var objectType string
						objectType, isType = obj.(string)
						if testContext.AssertThat(isType, "ObjectType is not a string") {
							if objectType == "DockerfileExecEventDesc" {
								obj = eventMap["ImageVersionObjId"]
								var versionObjId string
								versionObjId, isType = obj.(string)
								if testContext.AssertThat(isType, "ImageVersionObjId is not a string") {
									if versionObjId == imageVersion1ObjId {
										imageVersion1ObjIdCount++
									}
									if versionObjId == imageVersion2ObjId {
										imageVersion2ObjIdCount++
									}
									if versionObjId == "" {
										imageVersionEmptyCount++
									}
								}
							}
						}
					}
				}
			}
			testContext.AssertThat(imageVersionEmptyCount == 1, "Empty count is not 1")
			testContext.AssertThat(imageVersion1ObjIdCount == 0, "Version 1 count is not 0")
			testContext.AssertThat(imageVersion2ObjIdCount == 1, "Version 2 count is not 1")
		}
	}

	var image1ObjId string

	// Test ability to delete a docker image.
	{
		// Obtain the Id of the image that owns the remaining image version that we created.
		var dockerImageVersionDescMap map[string]interface{}
		dockerImageVersionDescMap = testContext.TryGetDockerImageDesc(imageVersion2ObjId, true)
		var obj interface{}
		obj = dockerImageVersionDescMap["ImageObjId"]
		var isType bool
		image1ObjId, isType = obj.(string)
		testContext.AssertThat(isType, "ImageObjId is not a string")
	}
	
	// Test ability to delete a dockerfile.
	{
		testContext.TryRemDockerfile(dockerfile1Id)
		if ! testContext.TestHasFailed() {
		
			// Verify that image creation event's reference to the dockerfile has
			// been nullified.
			var eventIds []string = testContext.TryGetDockerImageEvents(imageVersion2ObjId)
			if testContext.AssertThat(eventIds != nil, "") {
				testContext.AssertThat(len(eventIds) == 1, fmt.Sprintf(
					"Wrong number of event Ids (%d)", len(eventIds)))
				var noOfDockerfileExecEvents int = 0
				for _, eventId := range eventIds {
					var eventMap map[string]interface{}
					eventMap = testContext.TryGetEventDesc(eventId)
					var obj interface{}
					obj = eventMap["ObjectType"]
					if testContext.AssertThat(obj != nil, "No ObjectType field") {
						var objectType string
						var isType bool
						objectType, isType = obj.(string)
						if testContext.AssertThat(isType, "ObjectType is not a string") {
							if objectType == "DockerfileExecEventDesc" {
								if testContext.AssertThat(noOfDockerfileExecEvents == 0,
										"More than one DockerfileExecEventDesc for image version") {
									noOfDockerfileExecEvents++
									obj = eventMap["DockerfileId"]
									if testContext.AssertThat(obj != nil, "No DockerfileId field") {
										var dockerfileId string
										dockerfileId, isType = obj.(string)
										if testContext.AssertThat(isType, "DockerfileId is not a string") {
											testContext.AssertThat(dockerfileId == "", "DockerfileId was not nullified")
										}
									}
								}
							}
						}
					}
				}
				testContext.AssertThat(noOfDockerfileExecEvents == 1, "Wrong number of nofOfDockerfileExecEvents")
			}
		}
	}
	
	// Attempt to delete the Image object.
	{
		testContext.AssertThat(testContext.TryRemDockerImage(image1ObjId),
				"Unable to remove docker image")
	}
	
	// Test ability to delete a repo.
	{
		testContext.TryDeleteRepo(realmXRepo1Id)
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

	var err error
	var realmXId string
	var realmXAdminUserId = "admin"
	var realmXAdminPswd = "fluffy"
	var realmXAdminObjId string
	var realmXRepo1Id string
	var dockerImage1ObjId string
	var dockerImage1Version1ObjId string
	var scanConfigId string
	var dockerfilePath string
	var dockerfileParamPath string
	var dockerfile2ParamPath string
	var dockerfileId string
	var dockerfileParamId string
	var dockerfile2ParamId string
	var dockerfile2Path string
	var dockerfile3Path string
	//var dockerfile2Id string
	var flagImagePath = "Seal.png"
	var tempdir string
	
	{
		tempdir, err = utils.CreateTempDir()
		if err != nil { testContext.AbortAllTests(err.Error()) }
		
		realmXId, realmXAdminObjId, _ = testContext.TryCreateRealmAnon(
			"realm4", "realm 4 Org", realmXAdminUserId, "realm 4 Admin Full Name",
			"realm4admin@gmail.com", realmXAdminPswd)
		
		testContext.TryAuthenticate(realmXAdminUserId, realmXAdminPswd, true)
		
		realmXRepo1Id = testContext.TryCreateRepo(realmXId, "repo1", "Repo in realm x", "")
		
		scanConfigId = testContext.TryDefineScanConfig("My Config 1",
			"A very fine config", realmXRepo1Id, "clair", "", flagImagePath, []string{}, []string{})

		dockerfilePath, err = utils.CreateTempFile(tempdir, "Dockerfile", "FROM centos\nRUN echo moo > oink")
		if err != nil { testContext.AbortAllTests(err.Error()) }
		defer os.Remove(dockerfilePath)
		dockerfileId, _ = testContext.TryAddDockerfile(realmXRepo1Id, dockerfilePath, "A fine dockerfile")
		
		dockerfileParamPath, err = utils.CreateTempFile(tempdir, "DockerfileP",
			"FROM centos\nARG param1=\"good doggy\"\nRUN echo $param1 > doggy.txt")
		if err != nil { testContext.AbortAllTests(err.Error()) }
		defer os.Remove(dockerfileParamPath)
		//var dockerfileParamDescMap map[string]interface{}
		dockerfileParamId, _ = testContext.TryAddDockerfile(
			realmXRepo1Id, dockerfileParamPath, "A parameterized dockerfile")
		
		dockerfile2ParamPath, err = utils.CreateTempFile(tempdir, "Dockerfile2P",
			"FROM centos\nARG param1\nARG param2=\"abc def\"\nRUN echo $param2 > $param1")
		if err != nil { testContext.AbortAllTests(err.Error()) }
		defer os.Remove(dockerfile2ParamPath)
		var dockerfile2ParamDescMap map[string]interface{}
		dockerfile2ParamId, dockerfile2ParamDescMap = testContext.TryAddDockerfile(
			realmXRepo1Id, dockerfile2ParamPath, "A dockerfile with two params")
		
		// Check params that were returned.
		// Should be an array of objects, each containing a Name and Value string field.
		var objAr []interface{}
		var isType bool
		objAr, isType = dockerfile2ParamDescMap["ParameterValueDescs"].([]interface{})
		testContext.AssertThat(isType,
			"ParameterValueDescs is not an array of interface: it is a " +
			reflect.TypeOf(dockerfile2ParamDescMap["ParameterValueDescs"]).String())
		var params = make(map[string]string)
		for i, obj := range objAr {
			var param = obj.(map[string]interface{})
			if ! testContext.AssertThat(isType, fmt.Sprintf("Element %d is not a map[string]interface", i)) { continue }
			var name = param["Name"].(string)
			if testContext.AssertThat(name != "", "No Name for parameter") {
				var value = param["Value"]
				params[name] = value.(string)
			}
		}
		var contains bool
		_, contains = params["param1"]
		testContext.AssertThat(contains, "Parameter param1 was not returned")
		var param2Value string
		param2Value, contains = params["param2"]
		testContext.AssertThat(contains, "Parameter param2 was not returned")
		testContext.AssertThat(param2Value == "\"abc def\"", "Parameter param2 had wrong value: '" + param2Value + "'")
		
		dockerfile2Path, err = utils.CreateTempFile(tempdir, "Dockerfile2", "FROM centos\nRUN echo boo > ploink")
		if err != nil { testContext.AbortAllTests(err.Error()) }
		defer os.Remove(dockerfile2Path)
		testContext.TryAddDockerfile(realmXRepo1Id, dockerfile2Path, "A finer dockerfile")
		
		dockerfile3Path, err = utils.CreateTempFile(tempdir, "Dockerfile3", "FROM centos\nRUN echo split > splat")
		if err != nil { testContext.AbortAllTests(err.Error()) }
		defer os.Remove(dockerfile3Path)

		err = utils.DownloadFile(SealURL, flagImagePath, true)
		if err != nil { testContext.AbortAllTests(err.Error()) }
	}
	
	// -------------------------------------
	// Tests
	//
	
	// Test ability to build image from a dockerfile.
	{
		dockerImage1Version1ObjId, dockerImage1ObjId = testContext.TryExecDockerfile(realmXRepo1Id,
			dockerfileId, "myimage", []string{}, []string{})
		testContext.AssertThat(dockerImage1ObjId != "", "No image obj Id returned")
	}
	
	// Test ability to build image from a dockerfile that takes one parameter.
	{
		var imageObjId string
		_, imageObjId = testContext.TryExecDockerfile(realmXRepo1Id,
			dockerfileParamId, "myparamimage", []string{ "param1" },
			[]string{ "abc" })
		if testContext.AssertThat(imageObjId != "", "No image obj Id returned") {
			var imageInfo map[string]interface{}
			imageInfo = testContext.TryGetDockerImageDesc(imageObjId, true)
			if testContext.AssertThat(imageInfo != nil, "No image info") {
				var obj = imageInfo["Name"]
				var name string
				var isType bool
				name, isType = obj.(string)
				if testContext.AssertThat(isType, "No string value for Name in image desc") {
					testContext.AssertThat(name == "myparamimage",
						"Image has wrong name: " + name)
				}
			}
		}
	}
	
	// Test ability to build image from a dockerfile that takes two parameters.
	{
		var imageObjId string
		_, imageObjId = testContext.TryExecDockerfile(realmXRepo1Id,
			dockerfile2ParamId, "my2paramimage", []string{ "param1", "param2" },
			[]string{ "abc", "def" })
		if testContext.AssertThat(imageObjId != "", "No image obj Id returned") {
			var imageInfo map[string]interface{}
			imageInfo = testContext.TryGetDockerImageDesc(imageObjId, true)
			if testContext.AssertThat(imageInfo != nil, "No image info") {
				var obj = imageInfo["Name"]
				var name string
				var isType bool
				name, isType = obj.(string)
				if testContext.AssertThat(isType, "No string value for Name in image desc") {
					testContext.AssertThat(name == "my2paramimage",
						"Image has wrong name: " + name)
				}
			}
		}
	}
	
	// Test ability to list the images in a repo.
	{
		var imageNames []string = testContext.TryGetDockerImages(realmXRepo1Id)
		testContext.AssertThat(len(imageNames) == 3, fmt.Sprintf(
			"Wrong number of images: %d", len(imageNames)))
	}
	
	// Test abilty to get the current logged in user's docker images.
	{
		var myDockerImageIds []string = testContext.TryGetMyDockerImages()
		testContext.AssertThat(len(myDockerImageIds) == 3, "Wrong number of docker images")
	}
	
	// Test ability to scan a docker image.
	{
		var scanEventDescs []map[string]interface{}
		scanEventDescs = testContext.TryScanImage(scanConfigId, dockerImage1Version1ObjId)
		testContext.AssertThat(len(scanEventDescs) == 1, "Wrong number of scan events")
	}
	
	// Test ability to upload and exec a dockerfile in one command.
	{
		var dockerImage3ObjId string
		_, dockerImage3ObjId, _ = testContext.TryAddAndExecDockerfile(realmXRepo1Id,
			"My third image", "myimage3", dockerfile3Path, []string{}, []string{})
		fmt.Println(dockerImage3ObjId)
	
		/*
		testContext.TryDownloadImage(dockerImage3ObjId, "BooPloinkImage")
		var responseMap = testContext.TryGetDockerImageDesc(dockerImage3ObjId, true)
		if testContext.CurrentTestPassed {
			// Check image signature.
			var image2Signature []byte
			var err error
			image2Signature, err = utils.ComputeSHA512FileSignature("BooPloinkImage")
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
		testContext.AssertThat(len(imageIds) == 4, "Wrong number of docker images")
	}

	// Test ability to get the events for a specified user, including docker build events.
	{
		var eventIds []string = testContext.TryGetUserEvents(realmXAdminObjId)
		testContext.AssertThat(len(eventIds) == 5, "Wrong number of user events")
			// Should be one scan event and two dockerfile exec events.
	}
	
	// Test ability to get the events for a specified docker image.
	{
		var eventIds []string = testContext.TryGetDockerImageEvents(dockerImage1ObjId)
		testContext.AssertThat(len(eventIds) == 2,
			fmt.Sprintf("Wrong number of image events: it is %d", len(eventIds)))
			// Should be one scan event.
		
		// Try for an image version.
		eventIds = testContext.TryGetDockerImageEvents(dockerImage1Version1ObjId)
		testContext.AssertThat(len(eventIds) == 2,
			fmt.Sprintf("Wrong number of image events: it is %d", len(eventIds)))
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
		var eventIds []string
		var paramMap map[string]string
		eventIds, paramMap = testContext.TryGetDockerfileEvents(dockerfileId, dockerfilePath)
		testContext.AssertThat(len(eventIds) == 1, "Wrong number of image events")
		testContext.AssertThat(len(paramMap) == 0, "Wrong number of build parameters")
		
		eventIds, paramMap = testContext.TryGetDockerfileEvents(dockerfileParamId, dockerfileParamPath)
		testContext.AssertThat(len(eventIds) == 1, "Wrong number of image events")
		testContext.AssertThat(len(paramMap) == 1, "Wrong number of build parameters")
	}
	
	// Test abilit to delete a specified docker image version.
	{
		if testContext.TryRemImageVersion(dockerImage1Version1ObjId) {
			testContext.TryGetDockerImageDesc(dockerImage1Version1ObjId, false)
		}
	}
}
