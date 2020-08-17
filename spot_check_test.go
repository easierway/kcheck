package main

import (
	"fmt"
	"testing"
)

func TestRunningOnSpotCheck(t *testing.T) {
	ros := &RunningOnDifferentNodes{}
	data, err := loadDataFromFile("example_deployment_1.yaml")
	if err != nil {
		t.Error(err)
	}
	hints, err := ros.Check(data)
	if err != nil {
		t.Error(err)
	}
	t.Log(hints)
}

func TestCorrectionForRunningOnSpotCheck(t *testing.T) {
	ros := &RunningOnDifferentNodes{}
	data, err := loadDataFromFile("example_deployment_1.yaml")
	if err != nil {
		t.Error(err)
	}
	corrected, err := ros.Correct(data)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(corrected))
}
