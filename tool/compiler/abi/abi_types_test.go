package abi

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestNewMethod(t *testing.T) {
	v1 := InputType{
		Name: "testFunc",
		Type: "string",
	}
	v2 := OutputType{
		Type: "string",
	}

	method := &Method{
		Name:    "testFunc",
		Inputs:  []*InputType{&v1, &v1},
		Outputs: &v2,
		Type:    []string{"action", "event"},
	}

	jsonByte, err := json.MarshalIndent(method, "", "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(jsonByte))

	//unmarshal
	var m2 = new(Method)
	err = json.Unmarshal(jsonByte, m2)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(m2)
}

func TestNewActionEntry(t *testing.T) {
	v1 := InputType{
		Name: "testFunc",
		Type: "string",
	}

	v2 := OutputType{
		Type: "string",
	}
	method := &Method{
		Name:    "testFunc",
		Inputs:  []*InputType{&v1, &v1},
		Outputs: &v2,
		Type:    []string{"action", "event"},
	}
	actionEntry := NewActionEntry()
	actionEntry = append(actionEntry, method, method)
	jsonByte, err := json.MarshalIndent(actionEntry, "", "\t")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(string(jsonByte))
}