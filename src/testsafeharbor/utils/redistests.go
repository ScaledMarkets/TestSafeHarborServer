package utils

import (
	"redis"
	"fmt"
)

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryRedisPing(client redis.Client) {
	testContext.StartTest("TryRedisPing")

	err := client.Ping()
	if ! testContext.AssertErrIsNil(err, "failed to create the client") { return }
	
	fmt.Println("Ping succeeded")
	testContext.PassTestIfNoFailures()
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryRedisSetGetString(client redis.Client) {
	testContext.StartTest("TryRedis")

	var err error
	err = client.Set("abc", []byte("12345"))
	if ! testContext.AssertErrIsNil(err, "When setting value") { return }

	var value []byte
	value, err = client.Get("abc")
	if ! testContext.AssertErrIsNil(err, "When getting value") { return }
	fmt.Println("Retrieved value", string(value))
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryRedisSetGetStringArray(client *redis.Client) {
	testContext.StartTest("TryRedis")


}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryRedisGetReleaseLock(client *redis.Client) {
	testContext.StartTest("TryRedis")


}
