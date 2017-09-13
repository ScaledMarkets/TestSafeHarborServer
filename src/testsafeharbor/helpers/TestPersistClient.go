package helpers

import (
	"fmt"
	
)

type TestPersistClient interface {
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
