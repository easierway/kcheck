package main

import (
	//"encoding/json"
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	rs, cs, err := ParserRuleSetConfig("my_rules.yaml")
	if err != nil {
		t.Error(err)
	}
	for _, r := range rs {
		fmt.Println(*r)
	}
	fmt.Println(cs)
	// rSet := RuleSet{*rs}
	// str, e := json.Marshal(rSet)
	// fmt.Println(rSet)
	// fmt.Println(e, string(str))

}
