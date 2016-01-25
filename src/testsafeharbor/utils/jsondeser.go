package utils

import (
	"fmt"
	"reflect"
	//"strings"
	//"errors"
	"time"
	//"runtime/debug"
	
	//"testsafeharbor/rest"
)

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryJsonDeserTokenizer(json string, expected []string) {

	testContext.StartTest("TryJsonDeserTokenizer")

	var pos int = 0
	for i, expect := range expected {
		var token string = parseJSON_findNextToken(json, &pos)
		if ! testContext.AssertThat(token == expect,
			fmt.Sprintf("Token #%d, was %s, expected %s", (i+1), token, expect)) { break }
	}
	testContext.PassTestIfNoFailures()
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryJsonDeserString(json, expected string) {
	
	testContext.StartTest("TryJsonDeserString")

	var value reflect.Value
	var err error
	var pos int = 0
	value, err = parseJSON_string_value(json, &pos)
	testContext.AssertErrIsNil(err, "")
	testContext.AssertThat(value.IsValid(), "Value is not valid")
	if testContext.AssertThat(value.String() == expected,
		"value: " + value.String() + ", expected: " + expected) {
		fmt.Println("Success: value=" + value.String())
	}
	testContext.PassTestIfNoFailures()
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryJsonDeserTime(json string, expected time.Time) {
	
	testContext.StartTest("TryJsonDeserTime")

	var value reflect.Value
	var err error
	var pos int = 0
	value, err = parseJSON_time_value(json, &pos)
	testContext.AssertErrIsNil(err, "")
	testContext.AssertThat(value.IsValid(), "Value is not valid")
	if testContext.AssertThat(value.Interface() == expected,
		"value: " + value.String() + ", expected: " + expected.String()) {
		fmt.Println("Success: value=" + value.String())
	}
	testContext.PassTestIfNoFailures()
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryJsonDeserSimple() {
	testContext.StartTest("TryJsonDeserialization")

	var client TestPersistClient = &InMemClient{}
	var abc = &InMemABC{ 123, "this is a string", []string{"alpha", "beta"}, true }
	var jsonString = abc.toJSON()

	var typeName string
	var remainder string
	var err error
	typeName, remainder, err = retrieveTypeName(jsonString)
	testContext.AssertErrIsNil(err, "when retrieving type name")
	
	var methodName = "New" + typeName
	var method = reflect.ValueOf(client).MethodByName(methodName)
	if ! testContext.AssertThat(method.IsValid(),
		"Method with name " + methodName + " not found") { return }
	var argAr []reflect.Value
	argAr, err = parseJSON(remainder)
	fmt.Println("argAr has " + fmt.Sprintf("%d", len(argAr)) + " elements")

	var retValues []reflect.Value = method.Call(argAr)
	var retValue0 interface{} = retValues[0].Interface()
	var abc2 ABC
	var isType bool
	abc2, isType = retValue0.(ABC)
	if !isType { fmt.Println("abc2 is NOT an ABC") } else {
		fmt.Println("abc2 IS an ABC")
		fmt.Println(fmt.Sprintf("\tabc2.a=%d", abc2.getA()))
		fmt.Println("\tabc2.bs=" + abc2.getBs())
		//fmt.Println("\tabc2.Car=" + string(abc2.getCar()))
		//fmt.Println("\tabc2.db=" + string(abc2.getDb()))
	}
	testContext.PassTestIfNoFailures()
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryJsonDeserNestedType() {
	testContext.StartTest("TryJsonDeserNestedType")

	var client TestPersistClient = &InMemClient{}
	var def = &InMemDEF{
		ABC: client.NewABC(123, "this is a string", []string{"alpha", "beta"}, true),
		xyz: 456,
	}
	var jsonString = def.toJSON()

	var typeName string
	var remainder string
	var err error
	typeName, remainder, err = retrieveTypeName(jsonString)
	testContext.AssertErrIsNil(err, "when retrieving type name")
	
	var methodName = "New" + typeName
	var method = reflect.ValueOf(client).MethodByName(methodName)
	if ! testContext.AssertThat(method.IsValid(),
		"Method with name " + methodName + " not found") { return }
	var argAr []reflect.Value
	argAr, err = parseJSON(remainder)
	testContext.AssertErrIsNil(err, "when calling parseJSON")
	fmt.Println("argAr has " + fmt.Sprintf("%d", len(argAr)) + " elements")

	var retValues []reflect.Value = method.Call(argAr)
	var retValue0 interface{} = retValues[0].Interface()
	var def2 DEF
	var isType bool
	def2, isType = retValue0.(DEF)
	if !isType { fmt.Println("def2 is NOT an DEF") } else {
		fmt.Println("def2 IS a DEF")
		fmt.Println(fmt.Sprintf("\tdef2.a=%d", def2.getA()))
		fmt.Println("\tdef2.bs=" + def2.getBs())
		//fmt.Println("\tabc2.Car=" + string(abc2.getCar()))
		//fmt.Println("\tabc2.db=" + string(abc2.getDb()))
		fmt.Println(fmt.Sprintf("\tdef2.xyz=%d", def2.getXyz()))
	}
	testContext.PassTestIfNoFailures()
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryJsonDeserComplex(json string) {
	testContext.StartTest("TryJsonDeserComplex")
	
	
	var values []reflect.Value
	var err error
	values, err = parseJSON(json)
	testContext.AssertErrIsNil(err, "in parseJSON")
	if ! testContext.AssertThat(values != nil, "Nil returned for values") { return }
	if ! testContext.AssertThat(len(values) != 0, "Zero values returned") { return }
	
	
	
	testContext.PassTestIfNoFailures()
}
