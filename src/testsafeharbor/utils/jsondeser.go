package utils

import (
	"fmt"
	"reflect"
	"strings"
	"errors"
)

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryJsonDeserSimple() {
	testContext.StartTest("TryJsonDeserialization")

	var client Client = &InMemClient{}
	var abc = &InMemABC{ 123, "this is a string", []string{"alpha", "beta"}, true }
	var jsonString = abc.toJSON()
	fmt.Println("jsonString=" + jsonString)  // debug

	var typeName string
	var remainder string
	var err error
	typeName, remainder, err = retrieveTypeName(jsonString)
	testContext.AssertErrIsNil(err, "when retrieving type name")
	fmt.Println("typeName=" + typeName)  // debug
	fmt.Println("remainder=" + remainder)  // debug
	
	var methodName = "New" + typeName
	var method = reflect.ValueOf(client).MethodByName(methodName)
	if ! testContext.AssertThat(method.IsValid(),
		"Method with name " + methodName + " not found") { return }
	var argAr []reflect.Value = parseJSON(remainder)
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
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryJsonDeserNestedType() {
	testContext.StartTest("TryJsonDeserNestedType")

	var client Client = &InMemClient{}
	var def = &InMemDEF{
		ABC: client.NewABC(123, "this is a string", []string{"alpha", "beta"}, true),
		xyz: 456,
	}
	var jsonString = def.toJSON()
	fmt.Println("jsonString=" + jsonString)  // debug

	var typeName string
	var remainder string
	var err error
	typeName, remainder, err = retrieveTypeName(jsonString)
	testContext.AssertErrIsNil(err, "when retrieving type name")
	fmt.Println("typeName=" + typeName)  // debug
	fmt.Println("remainder=" + remainder)  // debug
	
	var methodName = "New" + typeName
	var method = reflect.ValueOf(client).MethodByName(methodName)
	if ! testContext.AssertThat(method.IsValid(),
		"Method with name " + methodName + " not found") { return }
	var argAr []reflect.Value = parseJSON(remainder)
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
}

type Client interface {
	NewABC(a int, bs string, car []string, db bool) ABC
}

type ABC interface {
	getA() int
	getBs() string
	getCar() []string
	getDb() bool
	toJSON() string
}

type InMemClient struct {
}

type InMemABC struct {
	a int
	bs string
	car []string
	db bool
}

func (client *InMemClient) NewABC(a int, bs string, car []string, db bool) ABC {
	fmt.Println("Entered NewABC...")  // debug
	var abc *InMemABC = &InMemABC{a, bs, car, db}
	fmt.Println("\tconstructed an ABC...")  // debug
	fmt.Println(fmt.Sprintf("\tabc.a=%d", abc.a))  // debug
	fmt.Println(fmt.Sprintf("\tabc.getA()=%d", abc.getA()))
	return abc
}

func (abc *InMemABC) getA() int {
	return abc.a
}

func (abc *InMemABC)  getBs() string {
	return abc.bs
}

func (abc *InMemABC)  getCar() []string {
	return abc.car
}

func (abc *InMemABC)  getDb() bool {
	return abc.db
}

func (abc *InMemABC) toJSON() string {
	var res = fmt.Sprintf("\"ABC\": {\"a\": %d, \"bs\": \"%s\", \"car\": [", abc.a, abc.bs)
		// Note - need to replace any quotes in abc.bs
	for i, s := range abc.car {
		if i > 0 { res = res + ", " }
		res = res + "\"" + s + "\""  // Note - need to replace any quotes in s
	}
	res = res + fmt.Sprintf("], %s}", BoolToString(abc.db))
	return res
}

type DEF interface {
	ABC
	getXyz() int
}

type InMemDEF struct {
	ABC
	xyz int
}

func (client *InMemClient) NewDEF(a int, bs string, car []string, db bool, x int) DEF {
	var def = &InMemDEF{
		ABC: client.NewABC(a, bs, car, db),
		xyz: x,
	}
	return def
}

func (def *InMemDEF) getXyz() int {
	return def.xyz
}

func (def *InMemDEF) toJSON() string {
	var res = fmt.Sprintf("\"DEF\": {\"a\": %d, \"bs\": \"%s\", \"car\": [",
		def.getA(), def.getBs())
		// Note - need to replace any quotes in abc.bs
	for i, s := range def.getCar() {
		if i > 0 { res = res + ", " }
		res = res + "\"" + s + "\""  // Note - need to replace any quotes in s
	}
	res = res + fmt.Sprintf("], %s, %d}", BoolToString(def.getDb()), def.xyz)
	return res
}


/*******************************************************************************
 * Retrieve the type name that precedes the JSON string. It is in the format,
 *	"type-name" : <json-string>
 * or, as BNF,
	line			::= '"' (no spaces) type_name (no spaces) '"'  ':'  string
	type_name		::= string
 * 
 * Returns the type name, and then the remainder of the string that follows the
 * first colon.
 */
func retrieveTypeName(json string) (typeName string, remainder string, err error) {
	
	var i = strings.Index(json, "\"")
	var s2 = json[i+1:]
	var j = strings.Index(s2, "\"")
	if j == -1 { return "", "", errors.New(
		fmt.Sprintf("Ill-formatted json: no \" found after pos %d", i)) }
	var s3 = s2[:j]
	
	var k = strings.Index(s2[j:], ":")
	if k == -1 { return "", "", errors.New(
		fmt.Sprintf("Ill-formatted json: no : found after position %d", j)) }
	
	return s3, s2[k:], nil
}

/*******************************************************************************
 * Parse each json field and return a Value for each.
 * Only built-in types are allowed in the JSON fields, including byte and time.Time.
 * Arrays of these are allowed as well. Recursive.
 * A nil result indicates that no object was found, but no syntax error either.
 * An empty array means that an object was found but it contained no fields.
 * Note: This function is not intended to be a general purpose JSON parser -
 * its use is limited to the needs of this application.
 * BNF of JSON syntax:
	obj_value		::= '{'  field...    // ellipsis indicates one or more
	field			::=	'"' (no spaces) field_name (no spaces) '"'  ':'  value
	field_name		::= <char_seq>
	value			::= array_value | simple_value
	array_value		::= '['  value  [ comma_value ]  ']'
	comma_value		::= ','  value  [ comma_value ]
	simple_value	::= number | string_value | bool_value | time_value
	string_value	::= '"' (no spaces assumed) <char_seq> (no spaces assumed) '"'
	bool_value		::= 'true' | 'false'
	time_value		::= 'time ' '"' <char_seq in time format> '"'
 */
func parseJSON(json string) (map[string]reflect.Value, error) {
	
	/*
	// For now, just parse an int value.
	var aval int = 10
	//fmt.Println("Parsing json: " + json)
	//fmt.Sscanf(json, "{%d}", &aval)
	fmt.Println(fmt.Sprintf("Parsed value of %d for aval", aval))  // debug
	
	var numArgs = 4  // normally we will obtain this from the number of JSON fields.
	var valAr []reflect.Value = make([]reflect.Value, numArgs)
	valAr[0] = reflect.ValueOf(aval)
	valAr[1] = reflect.ValueOf("")
	valAr[2] = reflect.ValueOf([]string{})
	valAr[3] = reflect.ValueOf(false)
	
	return valAr
	*/
	
	return parseJSON_obj_value(json, 0)
}

/*******************************************************************************
 * Parse a JSON object, delimited by { and }. If none found due to EOF, or if
 * the first token does not match the possible values for an obj_value, return nil.
 * Otherwise, return an array of the field values. If there is a syntax error,
 * return an error.
 */
func parseJSON_obj_value(json string, pos *int) (map[string]reflect.Value, error) {
	
	var token string = parseJSON_findNextToken(json, pos)
	if token == "" { return nil, nil }
	
	if token != "{" { return nil, nil }
	
	var values map[string]reflect.Value = make(map[string]reflect.Value)
	var int noOfFields = 0
	for {
		var fieldName string
		var value reflect.Value
		var err error
		fieldName, value, err = parseJSON_field(json, pos)
		if err != nil { return values, parseJSON_syntaxError(pos) }
		if value == nil { break } // no more fields
		noOfFields++
		values[fieldName] = value
	}
	
	token = parseJSON_findNextToken(json, pos)
	if token != "}" { return values, parseJSON_syntaxError(*pos) }
	
	return values, nil
}

func parseJSON_field(json string, pos *int) (string, reflect.Value, error) {
	
	var token = parseJSON_findNextToken(json, pos)
	if token == "" { return nil, nil }
	
	if token != "\"" { return nil, nil }
	
	var value reflect.Value
	
	var fieldName string
	var err error
	fieldName, err = parseJSON_field_name(json, pos)
	if err != nil { return "", nil, err }
	
	token = parseJSON_findNextToken(json, pos)
	if token != "\"" { return fieldName, value, parseJSON_syntaxError(*pos) }
	
	token = parseJSON_findNextToken(json, pos)
	if token != ":" { return fieldName, value, parseJSON_syntaxError(*pos) }
	
	var value reflect.Value
	value, err = parseJSON_value(json, pos)
	if err != nil { return fieldName, nil, err }
	if value == nil return fieldName, nil, parseJSON_syntaxError(*pos)
	
	return fieldName, value, nil
}

func parseJSON_field_name(json string, pos *int) (string, error) {
	var dblQuotePos = strings.Index(json[*pos:], "\"")
	if dblQuotePos == -1 { return "", errors.New(
		fmt.Sprintf("Terminating double quote not found for field name, after pos %d", *pos)) }
	// ....to do: recognize escapes, etc.
	var startPos = *pos
	*pos = dblQuotePos
	var res = json[startPos:dblQuotePos]
	fmt.Printf("parseJSON_field_name returning " + res)  // debug
	return res, nil
}

func parseJSON_value(json string, pos *int) (reflect.Value, error) {
	
	var token = parseJSON_findNextToken(json, pos)
	if token == "" { return nil, nil }
	
	if token == "[" {
		parseJSON_pushTokenBack(token, pos)
		return parseJSON_array_value(json, pos)
	} else {
		return parseJSON_simple_value(json, pos)
	}
}

func parseJSON_array_value(json string, pos *int) ([]reflect.Value, error) {
	
	var token = parseJSON_findNextToken(json, pos)
	if token == "" { return nil, nil }
	
	if token != "[" { return nil, parseJSON_syntaxError(pos) }
	
	var values = make([]reflect.Value, 0)
	
	var value reflect.Value
	var err error
	value, err = parseJSON_value(json, pos)
	if err != nil { return value, err }
	
	values = append(values, value)
	
	var commaValues []reflect.Value
	commaValues, err = parseJSON_comma_value(json, pos)
	if commaValues != nil {
		values = append(values, commaValues...)
	}
	
	token = parseJSON_findNextToken(json, pos)
	if token != "]" { return value, parseJSON_syntaxError(*pos) }
	
	return values, nil
}

func parseJSON_simple_value(json string, pos *int) (reflect.Value, error ) {
	
	var value reflect.Value
	var err error
	value, err = parseJSON_number(json, pos)
	if err != nil { return value, err }
	if value != nil { return value, nil }
	
	value, err = parseJSON_string_value(json, pos)
	if err != nil { return value, err }
	if value != nil { return value, nil }
	
	value, err = parseJSON_bool_value(json, pos)
	if err != nil { return value, err }
	if value != nil { return value, nil }
	
	value, err = parseJSON_time_value(json, pos)
	if err != nil { return value, err }
	if value != nil { return value, nil }
	
	return nil, parseJSON_syntaxError(pos)
}

func parseJSON_comma_value(json string, pos *int) ([]reflect.Value, error) {
	
	var token = parseJSON_findNextToken(json, pos)
	if token == "" { return nil, nil }
	
	if token != "," { return nil, parseJSON_syntaxError(pos) }
	
	var values = make([]reflect.Value, 0)
	
	var value reflect.Value
	var err error
	value, err = parseJSON_value(json, pos)
	if err != nil { return value, err }
	
	values = append(values, value)
	
	var commaValues []reflect.Value
	commaValues, err = parseJSON_comma_value(json, pos)
	if commaValues != nil {
		values = append(values, commaValues...)
	}
	
	return values, nil
}

func parseJSON_number(json string, pos *int) (reflect.Value, error) {
	
	var token = parseJSON_findNextToken(json, pos)
	if token == "" { return nil, nil }
	
	var number int
	_, err := fmt.Sscanf(token, "%d", &number)
	if err != nil { return nil, parseJSON_syntaxError(pos) }
	return reflect.ValueOf(number), err
}

func parseJSON_string_value(json string, pos *int) (reflect.Value, error) {
	
	var token = parseJSON_findNextToken(json, pos)
	if token == "" { return nil, nil }
	
	if token != "\"" { return nil, parseJSON_syntaxError(pos) }
	
	posOfNextDblQuote = strings.Index(json[*pos:], "\"")
	if posOfNextDblQuote == -1 { return nil, parseJSON_syntaxError(pos) }
	
	var startPos = *pos
	var strval = json[startPos: posOfNextDblQuote]
	
	var value = reflect.ValueOf(strval)
	*pos = posOfNextDblQuote+1
	return value, nil
}

func parseJSON_bool_value(json string, pos *int) (reflect.Value, error) {
	
	var token = parseJSON_findNextToken(json, pos)
	if token == "true" { return reflect.ValueOf(true), nil }
	if token == "false" { return reflect.ValueOf(false), nil }
	parseJSON_pushTokenBack(token, pos)
	return nil, nil
}

func parseJSON_time_value(json string, pos *int) (reflect.Value, error) {
	
	var token = parseJSON_findNextToken(json, pos)
	if token != "time" {
		parseJSON_pushTokenBack(token, pos)
		return nil, nil
	}
	
	posOfNextDblQuote = strings.Index(json[*pos:], "\"")
	if posOfNextDblQuote == -1 { return nil, parseJSON_syntaxError(pos) }
	
	var startPos = *pos
	var strval = json[startPos: *pos]
	
	var t time.Time
	var err error
	err = t.UnmarshalJSON([]byte(strval)) // scan the time value from strval
	if err != nil { return nil, parseJSON_syntaxError(pos) }
	var value = reflect.ValueOf(t)

	*pos = posOfNextDblQuote+1
	return value, nil
}

/*******************************************************************************
 * If EOF, return "". Update the position value to point to the start of the first
 * position following the token (which might be one past the end of the json, if
 * there is not more content in the string.
 */
func parseJSON_findNextToken(json string, pos *int) string {
	var restOfJson = json[*pos:]
	var posOfNextWhitespace = strings.IndexAny(restOfJson, " \t\r\n")
	if posOfNextWhitespace == -1 { return "" }
	*pos = *pos + posOfNextWhitespace
	return restOfJson[:posOfNextWhitespace]
}

func parseJSON_pushTokenBack(token string, pos *int) {
	*pos = *pos - len(token)
}

func parseJSON_syntaxError(pos *int) error {
	returns errors.New(fmt.Sprintf("Syntax error at position %d", *pos))
}
