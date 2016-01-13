package utils

import (
	"fmt"
	"reflect"
	"strings"
	"errors"
	"time"
	"runtime/debug"
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
	testContext.PassTestIfNoFailures()
}

/*******************************************************************************
 * 
 */
func (testContext *TestContext) TryJsonDeserSimple() {
	testContext.StartTest("TryJsonDeserialization")

	var client Client = &InMemClient{}
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
	var abc *InMemABC = &InMemABC{a, bs, car, db}
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
	res = res + fmt.Sprintf("], \"db\": %s}", BoolToString(abc.db))
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
	res = res + fmt.Sprintf("], \"db\": %s, \"xyz\": %d}",
		BoolToString(def.getDb()), def.xyz)
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
	
	return s3, s2[j+k+1:], nil
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
	obj_value		::= '{'  field  [ comma_field ]  '}'
	field			::=	'"' (no spaces) field_name (no spaces) '"'  ':'  value
	field_name		::= <char_seq>
	value			::= array_value | simple_value
	comma_field		::= ','  field  [ comma_field ]
	array_value		::= '['  value  [ comma_value ]  ']'
	comma_value		::= ','  value  [ comma_value ]
	simple_value	::= number | string_value | bool_value | time_value
	string_value	::= '"' (no spaces assumed) <char_seq> (no spaces assumed) '"'
	bool_value		::= 'true' | 'false'
	time_value		::= 'time ' '"' <char_seq in time format> '"'
 */
func parseJSON(json string) ([]reflect.Value, error) {
	
	var pos int = 0
	return parseJSON_obj_value(json, &pos)
}

/*******************************************************************************
 * Parse a JSON object, delimited by { and }. If none found due to EOF, or if
 * the first token does not match the possible values for an obj_value, return nil.
 * Otherwise, return an array of the field values. If there is a syntax error,
 * return an error.
 */
func parseJSON_obj_value(json string, pos *int) ([]reflect.Value, error) {
	
	var token string = parseJSON_findNextToken(json, pos)
	if token == "" { return nil, nil }
	
	if token != "{" {
		parseJSON_pushTokenBack(token, pos)
		return nil, nil
	}
	
	var err error
	var values []reflect.Value = make([]reflect.Value, 0)
	var fieldValue reflect.Value
	_, fieldValue, err = parseJSON_field(json, pos)
	if err != nil { return values, err }
	if ! fieldValue.IsValid() { return values, nil } // no fields
	
	values = append(values, fieldValue)
	
	// Add additional fields, if any.
	var addlFieldValues []reflect.Value
	addlFieldValues, err = parseJSON_comma_field(json, pos)
	if err != nil { return values, err }
	if addlFieldValues != nil {
		values = append(values, addlFieldValues...)
	}
	
	token = parseJSON_findNextToken(json, pos)
	if token != "}" { return values, parseJSON_tokenError(token, pos,
		"while looking for object terminator") }
	
	return values, nil
}

func parseJSON_field(json string, pos *int) (string, reflect.Value, error) {
	
	var value reflect.Value

	var token = parseJSON_findNextToken(json, pos)
	if token == "" { return "", value, nil }
	
	if token != "\"" {
		parseJSON_pushTokenBack(token, pos)
		return "", value, nil
	}
	
	var fieldName string
	var err error
	fieldName, err = parseJSON_field_name(json, pos)
	if err != nil { return "", value, err }
	if fieldName == "" { return "", value, parseJSON_syntaxError(pos,
		"Did not find field name") }
	
	token = parseJSON_findNextToken(json, pos)
	if token != "\"" { return fieldName, value, parseJSON_tokenError(
		token, pos, "while looking for \" following a field name") }
	
	token = parseJSON_findNextToken(json, pos)
	if token != ":" { return fieldName, value, parseJSON_tokenError(
		token, pos, "while looking for colon following a field name") }
	
	value, err = parseJSON_value(json, pos)
	if err != nil { return fieldName, value, err }
	if ! value.IsValid() { return fieldName, value, parseJSON_tokenError(
		token, pos, "while looking for object field value") }
	
	return fieldName, value, nil
}

func parseJSON_comma_field(json string, pos *int) ([]reflect.Value, error) {
	
	var token = parseJSON_findNextToken(json, pos)
	if token == "" { return nil, nil }
	
	if token != "," {
		parseJSON_pushTokenBack(token, pos)
		return nil, nil
	}
	
	var err error
	var value reflect.Value
	var values []reflect.Value = make([]reflect.Value, 0)
	_, value, err = parseJSON_field(json, pos)
	if err != nil { return nil, err }
	values = append(values, value)
	
	var addlValues []reflect.Value
	addlValues, err = parseJSON_comma_field(json, pos)
	if err != nil { return values, err }
	if addlValues != nil {
		values = append(values, addlValues...)
	}
	
	return values, nil
}

func parseJSON_field_name(json string, pos *int) (string, error) {

	// Find trailing double-quote.
	var dblQuotePos = strings.Index(json[*pos:], "\"")
	if dblQuotePos == -1 { return "", errors.New(
		fmt.Sprintf("Terminating double quote not found for field name, after pos %d", *pos)) }
	
	// ....to do: recognize escapes, etc.
	
	var startPos = *pos  // beginning of field name
	*pos += dblQuotePos  // update json pos to point to the trailing double-quote
	var fieldName = json[startPos:*pos]
	return fieldName, nil
}

func parseJSON_value(json string, pos *int) (reflect.Value, error) {
	
	var value reflect.Value
	var token = parseJSON_findNextToken(json, pos)
	if token == "" { return value, nil }
	
	if token == "[" {
		parseJSON_pushTokenBack(token, pos)
		return parseJSON_array_value(json, pos)
	} else {
		parseJSON_pushTokenBack(token, pos)
		return parseJSON_simple_value(json, pos)
	}
}

func parseJSON_array_value(json string, pos *int) (reflect.Value, error) {
	
	var value reflect.Value
	var token = parseJSON_findNextToken(json, pos)
	if token == "" { return value, nil }
	
	if token != "[" {
		parseJSON_pushTokenBack(token, pos)
		return value, nil
	}
	
	var err error
	value, err = parseJSON_value(json, pos)
	if err != nil { return value, err }
	if ! value.IsValid() { return value, nil } // no fields
	
	// If there are more elements in the array, they must be of the same underlying
	// simple type as the above value.
	var elementType = value.Type()
	
	var sliceType = reflect.SliceOf(elementType)
	var slice = reflect.MakeSlice(sliceType, 0, 1)
	slice = reflect.Append(slice, value)
	
	var commaValues reflect.Value
	commaValues, err = parseJSON_comma_value(json, pos, elementType)
	if err != nil { return slice, err }
	if commaValues.IsValid() {
		slice = reflect.AppendSlice(slice, commaValues)
	}
	
	token = parseJSON_findNextToken(json, pos)
	if token != "]" { return slice, parseJSON_tokenError(token, pos,
		"while looking for array value") }
	
	return slice, nil
}

func parseJSON_simple_value(json string, pos *int) (reflect.Value, error ) {
	
	var value reflect.Value
	var err error
	value, err = parseJSON_number(json, pos)
	if err != nil { return value, err }
	if value.IsValid() { return value, nil }
	
	value, err = parseJSON_string_value(json, pos)
	if err != nil { return value, err }
	if value.IsValid() { return value, nil }
	
	value, err = parseJSON_bool_value(json, pos)
	if err != nil { return value, err }
	if value.IsValid() { return value, nil }
	
	value, err = parseJSON_time_value(json, pos)
	if err != nil { return value, err }
	if value.IsValid() { return value, nil }
	
	return value, parseJSON_syntaxError(pos, "While looking for simple value")
}

func parseJSON_comma_value(json string, pos *int,
	elementType reflect.Type) (reflect.Value, error) {
	
	var value reflect.Value
	var token = parseJSON_findNextToken(json, pos)
	if token == "" { return value, nil }
	
	if token != "," {
		parseJSON_pushTokenBack(token, pos)
		return value, nil
	}
	
	var err error
	value, err = parseJSON_value(json, pos)
	if err != nil { return value, err }
	
	var sliceType = reflect.SliceOf(elementType)
	var slice = reflect.MakeSlice(sliceType, 0, 1)
	slice = reflect.Append(slice, value)
	
	var commaValues reflect.Value
	commaValues, err = parseJSON_comma_value(json, pos, elementType)
	if err != nil { return slice, err }
	if commaValues.IsValid() {
		slice = reflect.AppendSlice(slice, commaValues)
	}
	
	return slice, nil
}

func parseJSON_number(json string, pos *int) (reflect.Value, error) {
	
	var value reflect.Value
	var token = parseJSON_findNextToken(json, pos)
	if token == "" { return value, nil }
	
	var number int
	_, err := fmt.Sscanf(token, "%d", &number)
	if err != nil {
		parseJSON_pushTokenBack(token, pos)
		return value, nil
	}
	return reflect.ValueOf(number), err
}

func parseJSON_string_value(json string, pos *int) (reflect.Value, error) {
	
	var value reflect.Value
	var token = parseJSON_findNextToken(json, pos)
	if token == "" { return value, nil }
	
	if token != "\"" {
		parseJSON_pushTokenBack(token, pos)
		return value, nil
	}
	
	var posOfNextDblQuote = strings.Index(json[*pos:], "\"")
	if posOfNextDblQuote == -1 { return value, parseJSON_syntaxError(pos,
		"While looking for a string value") }
	
	var startPos = *pos
	*pos += posOfNextDblQuote
	var strval = json[startPos: *pos]
	*pos++  // advance one past the trailing double quote
	
	value = reflect.ValueOf(strval)
	return value, nil
}

func parseJSON_bool_value(json string, pos *int) (reflect.Value, error) {
	
	var value reflect.Value
	var token = parseJSON_findNextToken(json, pos)
	if token == "true" { return reflect.ValueOf(true), nil }
	if token == "false" { return reflect.ValueOf(false), nil }
	parseJSON_pushTokenBack(token, pos)
	return value, nil
}

func parseJSON_time_value(json string, pos *int) (reflect.Value, error) {
	
	var value reflect.Value
	var token = parseJSON_findNextToken(json, pos)
	if token != "time" {
		parseJSON_pushTokenBack(token, pos)
		return value, nil
	}
	
	var posOfNextDblQuote = strings.Index(json[*pos:], "\"")
	if posOfNextDblQuote == -1 { return value, parseJSON_syntaxError(pos,
		"While looking for a time value") }
	
	var startPos = *pos
	var strval = json[startPos: *pos]
	
	var t time.Time
	var err error
	err = t.UnmarshalJSON([]byte(strval)) // scan the time value from strval
	if err != nil { return value, parseJSON_syntaxError(pos,
		"While looking for a time value") }
	value = reflect.ValueOf(t)

	*pos = posOfNextDblQuote+1
	return value, nil
}

/*******************************************************************************
 * If EOF, return "". Update the position value to point to the start of the first
 * position following the token (which might be one past the end of the json, if
 * there is not more content in the string.
 */
func parseJSON_findNextToken(json string, pos *int) (token string) {

	var whitespace = " \t\r\n"
	var trimmedJson = strings.TrimLeft(json[*pos:], whitespace)
	if len(trimmedJson) == 0 {
		*pos = len(json)
		token = ""  // found no token
		return
	}
	
	*pos += (len(json[*pos:]) - len(trimmedJson))  // advance to start of token
	
	var specialJsonChars = "\":'[]{},"
	
	var posAfterToken int
	if strings.IndexAny(trimmedJson[:1], specialJsonChars) == 0 { // token is a special char
		posAfterToken = 1
	} else {
		posAfterToken = strings.IndexAny(trimmedJson[1:], whitespace + specialJsonChars)
		if posAfterToken == -1 {  // token goes through end of line
			posAfterToken = len(trimmedJson)
		} else {
			posAfterToken++  // account for fact that we started counting from position 1
		}
	}
	
	*pos += posAfterToken
	token = trimmedJson[:posAfterToken]
	return
}

func parseJSON_pushTokenBack(token string, pos *int) {
	*pos = *pos - len(token)
}

func parseJSON_syntaxError(pos *int, msg string) error {
	var err = errors.New(fmt.Sprintf("%s: at char no. %d", msg, (*pos + 1)))
	fmt.Println(err.Error())
	debug.PrintStack()
	return err
}

func parseJSON_tokenError(token string, pos *int, msg string) error {
	parseJSON_pushTokenBack(token, pos)
	var err = errors.New(fmt.Sprintf("Syntax error at char no. %d: %s %s",
		(*pos + 1), token, msg))
	fmt.Println(err.Error())
	debug.PrintStack()
	return err
}
