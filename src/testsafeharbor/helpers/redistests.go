package helpers

import (
	"redis"
	"fmt"
	"reflect"
)

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryRedisPing(redisClient redis.Client) {
	testContext.StartTest("TryRedisPing")

	err := redisClient.Ping()
	if ! testContext.AssertErrIsNil(err, "failed to create the client") { return }
	
	fmt.Println("Ping succeeded")
	testContext.PassTestIfNoFailures()
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryRedisSetGetString(redisClient redis.Client) {
	testContext.StartTest("TryRedis")

	var err error
	err = redisClient.Set("abc", []byte("12345"))
	if ! testContext.AssertErrIsNil(err, "When setting value") { return }

	var value []byte
	value, err = redisClient.Get("abc")
	if ! testContext.AssertErrIsNil(err, "When getting value") { return }
	fmt.Println("Retrieved value", string(value))
}

/*******************************************************************************
 * If this passes, then our redis/JSON database strategy works.
 */
func (testContext *TestContext) TryRedisGetJSONObject(redisClient redis.Client) {
	testContext.StartTest("TryRedisGetJSONObject")
	
	// Create a target object, on which there is a New method.
	var target TestPersistClient = &InMemClient{}

	// Create an object.
	var abc = &InMemABC{ 123, "this is a string", []string{"alpha", "beta"}, true }
		// ....to do: Add a time field.
		// ....to do: Add a []bool field.
	
	// Write it as JSON.
	var jsonString = abc.toJSON()
	
	// Store in redis.
	var err error
	err = redisClient.Set("abc", []byte(jsonString))
	
	// Retrieve from redis.
	var bytes []byte
	bytes, err = redisClient.Get("abc")
	if ! testContext.AssertErrIsNil(err, "When getting value") { return }
	var stringFromRedis = string(bytes)
	
	// Reconstitute as object.
	var typeName string
	var obj interface{}
	typeName, obj, err = ReconstituteObject(target, stringFromRedis)
	if ! testContext.AssertErrIsNil(err, "When reconstituting object") { return }
	var reconstitutedABC ABC
	var isType bool
	testContext.AssertThat(obj != nil, "Object is nil")
	reconstitutedABC, isType = obj.(ABC)
	if ! testContext.AssertThat(isType, "Object is not an ABC") {
		fmt.Println("Object is a " + reflect.TypeOf(obj).String())
	}
	if ! testContext.AssertThat(typeName == "ABC", "Object is not an InMemABC") {
		fmt.Println("typeName=" + typeName)
	}
	if ! testContext.AssertThat(reconstitutedABC.getA() == 123, "Wrong value for abc") {
		fmt.Println("obj.getABC()=" + string(reconstitutedABC.getA()))
	}
	
	testContext.PassTestIfNoFailures()
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryRedisGetReleaseLock(redisClient redis.Client) {
	testContext.StartTest("TryRedis")

}
