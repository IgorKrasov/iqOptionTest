package app

import (
	"testing"
	"time"
)

func TestCache_Get_Set(t *testing.T) {
	tc := NewCache(0)
	a, err := tc.get("a")
	if a != nil || err == nil {
		t.Error("Getting A found value that shouldn't exist:", a)
	}

	b, err := tc.get("b")
	if b != nil || err == nil {
		t.Error("Getting B found value that shouldn't exist:", b)
	}

	c, err := tc.get("c")
	if c != nil || err == nil {
		t.Error("Getting C found value that shouldn't exist:", c)
	}

	tc.set("a", 1, 1)
	tc.set("b", "b", 1)
	tc.set("c", 3.5, 1)

	a, err = tc.get("a")
	if a == nil || err != nil {
		t.Error("Can't find  A value that should exist", a)
	}

	ai, ok := a.(int)
	if !ok {
		t.Error("A don't integer", a)
	}

	if ai != 1 {
		t.Error("Geting A value don't equals 1", ai)
	}

	b, err = tc.get("b")
	if b == nil || err != nil {
		t.Error("Can't find  B value that should exist", b)
	}

	bs, ok := b.(string)

	if !ok {
		t.Error("B don't string", b)
	}

	if bs != b {
		t.Error("Geting B value don't equals b", bs)
	}

	c, err = tc.get("c")
	if c == nil || err != nil {
		t.Error("Can't find  C value that should exist", c)
	}

	cf, ok := c.(float64)
	if !ok {
		t.Error("C don't float", c)
	}

	if cf != 3.5 {
		t.Error("Geting C value don't equals 3.5", cf)
	}

}

func TestCache_Keys(t *testing.T) {
	tc := NewCache(0)
	keys := tc.keys()
	if len(keys) != 0 {
		t.Error("Geting not empty slice, when keys doesn't exist", keys)
	}

	tc.set("a", 1, 0)
	tc.set("b", "b", 0)
	tc.set("c", 3.5, 0)

	keys = tc.keys()

	if len(keys) != 3 {
		t.Error("Length keys doesn't equals 3", keys)
	}

	ab := false
	bb := false
	cb := false
	for _, v := range keys {
		if v == "a" {
			ab = true
		}
		if v == "b" {
			bb = true
		}
		if v == "c" {
			cb = true
		}
	}

	if ab == false {
		t.Error("Key A doesn't found", keys)
	}
	if bb == false {
		t.Error("Key B doesn't found", keys)
	}
	if cb == false {
		t.Error("Key C doesn't found", keys)
	}
}

func TestCache_List(t *testing.T) {
	tc := NewCache(0)
	key := "l"
	list, err := tc.lgetall(key)
	if list != nil || err == nil {
		t.Error("Find l value that shouldn't exist", list)
	}

	tc.rpush(key, 1, 0)
	tc.rpush(key, "b", 0)
	tc.rpush(key, 3.5, 0)

	list, err = tc.lgetall(key)
	if list == nil || err != nil {
		t.Error("Cant find list L that should be exist", list)
	}
	if len(list) != 3 {
		t.Error("Length L doesn't equals 3", list)
	}
	f, err := tc.lget(key, 0)
	if f == nil || err != nil {
		t.Error("Can't find 0 element in list that should be exist", f)
	}

	s, err := tc.lget(key, 1)
	if s == nil || err != nil {
		t.Error("Can't find 1 element in list that should be exist", s)
	}

	th, err := tc.lget(key, 2)
	if th == nil || err != nil {
		t.Error("Can't find 2 element in list that should be exist", th)
	}

	lp, err := tc.pop(key)
	if lp == nil || err != nil {
		t.Error("Can't pop last element from list that should be exist", lp)
	}

	lpf, ok := lp.(float64)
	if !ok {
		t.Error("Element lp doesn't float")
	}
	if lpf != 3.5 {
		t.Error("Geting lp value don't equals 3.5", lpf)

	}
	list2, _ := tc.lgetall(key)
	if len(list2) != 2 {
		t.Error("length of List after pop doesn't equals 2", list2)
	}
}

func TestCache_Hash(t *testing.T) {
	tc := NewCache(0)
	key := "р"
	dict, err := tc.hgetall(key)
	if dict != nil || err == nil {
		t.Error("Find dict value that shouldn't exist", dict)
	}

	tc.hset(key, map[string]interface{}{"name": "Igor", "age": 28}, 0)
	dict, err = tc.hgetall(key)
	if dict == nil || err != nil {
		t.Error("Cant find dict that should be exist", dict)
	}

	name, ok := dict["name"]
	if !ok {
		t.Error("Cant find dict name that should be exist", dict)
	}

	nameS, ok := name.(string)
	if !ok {
		t.Error("dict name that doesn't string", dict["name"])
	}

	if nameS != "Igor" {
		t.Error("Name doesn't equals Igor", dict["name"])
	}

	age, ok := dict["age"]
	if !ok {
		t.Error("Cant find dict age that should be exist", dict)
	}

	ageI, ok := age.(int)
	if !ok {
		t.Error("dict age doesn't int", dict["age"])
	}

	if ageI != 28 {
		t.Error("Age doesn't equals 28", dict["name"])
	}

	test, ok := dict["test"]
	if test != nil || ok {
		t.Error("Find dict value test that shouldn't exist", dict)
	}

	tc.hset(key, map[string]interface{}{"test": 3.5}, 0)

	test, err = tc.hget(key, "test")
	if test == nil || err != nil {
		t.Error("Cant find dict test that should be exist", dict)
	}

	name2, err := tc.hget(key, "name")
	if name2 == nil || err != nil {
		t.Error("Cant find dict name that should be exist", dict)
	}

	age2, err := tc.hget(key, "age")
	if age2 == nil || err != nil {
		t.Error("Cant find dict age that should be exist", dict)
	}
}

func TestCache_Unset(t *testing.T) {
	tc := NewCache(0)
	tc.set("a", "a", 0)
	a, err := tc.get("a")
	if a == nil || err != nil {
		t.Error("Can't find value A that should be exist", a)
	}
	tc.deleteItem("a")

	a2, err2 := tc.get("a")
	if a2 != nil || err2 == nil {
		t.Error("Find value A that shouldn't exist", a2)
	}
}

func TestCache_Expired(t *testing.T) {
	tc := NewCache(time.Duration(5 * time.Millisecond))
	tc.ExpiredTimeMultiplier = time.Millisecond
	tc.set("a", 1, 10)
	tc.set("b", 2, 0)
	tc.set("c", 3, 30)
	tc.set("d", 4, 70)

	<-time.After(20 * time.Millisecond)
	_, err := tc.get("a")
	if err == nil {
		t.Error("Found a when it should have been automatically deleted")
	}

	<-time.After(40 * time.Millisecond)
	_, err2 := tc.get("c")
	if err2 == nil {
		t.Error("Found c when it should have been automatically deleted")
	}

	_, err3 := tc.get("b")
	if err3 != nil {
		t.Error("Did not find b even though it was set to never expire")
	}

	_, err4 := tc.get("d")
	if err4 != nil {
		t.Error("Did not find d even though it was set to expire later than the default")
	}

	<-time.After(80 * time.Millisecond)
	_, err5 := tc.get("d")
	if err5 == nil {
		t.Error("Found d when it should have been automatically deleted (later than the default)")
	}
}
