package utils

import (
	"goredis"
	"fmt"
)

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryGoRedisPing(redis *goredis.Redis) {
	testContext.StartTest("TryRedisPing")

	err := redis.Ping()
	if ! testContext.AssertErrIsNil(err, "failed to create the client") { return }
	
	fmt.Println("Ping succeeded")
	testContext.PassTestIfNoFailures()
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryGoRedisSetGetString(redis *goredis.Redis) {
	testContext.StartTest("TryRedis")

	var err error
	err = redis.Set("abc", "12345", 0, 0, false, false)
	// args: key, value string, seconds, milliseconds int, mustExists, mustNotExists bool
	if ! testContext.AssertErrIsNil(err, "When setting value") { return }

	var value []byte
	value, err = redis.Get("abc")
	if ! testContext.AssertErrIsNil(err, "When getting value") { return }
	fmt.Println("Retrieved value", string(value))
	testContext.PassTestIfNoFailures()
}

/*******************************************************************************
 * Test SIsmember, SAdd.
 */
func (testContext *TestContext) TryGoRedisSet(redis *goredis.Redis) {
	testContext.StartTest("TryRedisGetJSONObject")
	
	var numElementsAdded int64
	var err error
	numElementsAdded, err = redis.SAdd("set1", "I belong to number one!")
	// SAdd(key string, members ...string) (int64, error)
	testContext.AssertErrIsNil(err, "after SAdd")
	testContext.AssertThat(numElementsAdded == 1,
		fmt.Sprintf("Returned %d for numElementsAdded", numElementsAdded))
	
	var b bool
	b, err = redis.SIsMember("set1", "I belong to number one!")
	// SIsMember(key, member string) (bool, error)
	testContext.AssertThat(b, "Oh no - I don't belong to number one!!!")
	testContext.PassTestIfNoFailures()
}
