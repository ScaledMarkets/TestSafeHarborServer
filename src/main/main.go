/*******************************************************************************
 * Perform independent end-to-end tests on the SafeHarbor server.
 * It is assumed that the SafeHarbor server is running on localhost:6000.
 */

package main

import (
	"fmt"
	"net/http"
)

func main() {
	
	// Test ability to create a realm.
	realmId := TryCreateRealm()
	assertThat(realmId != "", "TryCreateRealm failed")
	
	// Test ability create a repo.
	repoId := TryCreateRepo(realmId)
	assertThat(repoId != "", "TryCreateRepo failed")
		
	// Test ability to upload a Dockerfile.
	//dockerfileId := TryUploadDockerfile(repoId, dockerfilePath)
	//assertThat(dockerfileId != "", "TryUploadDockerfile failed")
	
	// Test ability to build image from a dockerfile.
	//imageId := TryBuildImage(dockerfileId)
	//assertThat(imageId != "", "TryBuildImage failed")
}


/*******************************************************************************
 * Verify that we can create a new realm.
 */
func TryCreateRealm() string {
	
	fmt.Println("TryCreateRealm")
	var resp *http.Response = sendPost(
		"createRealm",
		[]string{"Name"},
		[]string{"My Realm"})
	
	defer resp.Body.Close()
	
	verify200Response(resp)
	
	// Get the realm Id that is returned in the response body.
	
	var responseMap map[string]string = parseResponseBody(resp.Body)
	var realmId string = responseMap["Id"]
	if realmId == "" { printMap(responseMap) }
	assertThat(realmId != "", "Realm Id not found in response body")
	
	return realmId
}

/*******************************************************************************
 * Verify that we can create a new repo. This requires that we first created
 * a realm that the repo can belong to.
 */
func TryCreateRepo(realmId string) string {
	fmt.Println("TryCreateRepo")
	var resp *http.Response = sendPost(
		"createRepo",
		[]string{"RealmId", "Name"},
		[]string{realmId, "John's Repo"})
	
	defer resp.Body.Close()

	verify200Response(resp)
	
	// Get the repo Id that is returned in the response body.
	
	var responseMap map[string]string = parseResponseBody(resp.Body)
	var repoId string = responseMap["Id"]
	if repoId == "" { printMap(responseMap) }
	//var repoName string = responseMap["Name"]
	assertThat(repoId != "", "Repo Id not found in response body")
	
	return repoId
}

/*******************************************************************************
 * Verify that we can upload a dockerfile. This requries that we first created
 * a repo to uplaod it into.
 */
func TryUploadDockerfile(repoId string, dockerfilePath string) string {
	return ""
}

/*******************************************************************************
 * Verify that we can build an image, from a dockerfile that has already been
 * uploaded into a repo and for which we have the SafeHarborServer image id.
 */
func TryBuildImage(dockerfileId string) string {
	return ""
}
