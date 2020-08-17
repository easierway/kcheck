package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

// Checker is to check the K8S deployment/configuration with the different rules
type Checker interface {
	// Check is to check the file and return the suggestions
	Check(data []byte) (string, error)
}

// Corrector is to correct the K8S deployment/configuration with the different rules
// If the checker implementation also realizes the interface, the configuration can be auto-corrected when it breaks the check item.
type Corrector interface {
	// Correct is to correct the config
	Correct(org []byte) ([]byte, error)
}

func loadDataFromFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func writeDataToFile(data []byte, filePath string) error {
	return ioutil.WriteFile(filePath, data, 0777)
}

// Rule is composed of a set of checkers
type Rule struct {
	// Name is the name of the rule
	Name string
	// Checkers is the checkers relating to the rule's check items
	Checkers []Checker
}

func isStringParamValid(param *string, prompt string) bool {
	if param == nil || *param == "" {
		fmt.Println(prompt)
		return false
	}
	return true
}

//kcheck -d [rule definition file] -f [filePath] -r [checking set] -c
func main() {
	correctors := []Corrector{}
	//checkingRules := make(map[string]Checker)
	srcFile := flag.String("f", "", "the kubernetes deployment/configuration file")
	ruleName := flag.String("r", "", "the name of the checking rule")
	ruleConfig := flag.String("d", "", "the rule definition file")
	isCorrected := flag.Bool("c", false, "[Optional] try to correct the files ")
	needHelp := flag.Bool("help", false, "get the help")
	flag.Parse()

	if *needHelp {
		flag.Usage()
	}

	isOk := isStringParamValid(srcFile, "Set the file needing to check with '-f'. To get help with '-help' ")
	isOk = isStringParamValid(ruleName, "Set the rule for checking with '-r'")

	if !isOk {
		os.Exit(-1)
	}

	isOk = isStringParamValid(ruleConfig, "The default rule definition would be used."+
		" For setting custom rule definition, set the rule definition file by '-d'.")
	if isOk {
		var err error
		ruleSet, correctors, err = ParserRuleSetConfig(*ruleConfig)
		if err != nil {
			panic(err)
		}
	}
	srcData, err := loadDataFromFile(*srcFile)
	if err != nil {
		fmt.Printf("Failed to load the file '%s'.", *srcFile)
		os.Exit(-2)
	}
	var rule *Rule
	for _, r := range ruleSet {
		if r.Name == *ruleName {
			rule = r
		}
	}
	if rule == nil {
		fmt.Printf("Could not find the checking rule ‘%s’\n.", *ruleName)
		os.Exit(-3)
	}

	for _, check := range rule.Checkers {
		hint, err := check.Check(srcData)
		if err != nil {
			fmt.Printf("checking error: %v", err)
			os.Exit(-4)
		}
		fmt.Println(hint)
	}

	if !*isCorrected {
		return
	}

	correctedData := srcData
	for _, corrector := range correctors {
		cdata, err := corrector.Correct(correctedData)
		if err == nil {
			correctedData = cdata
		} else {
			fmt.Println("Corrector erorr", err)
		}

	}
	writeDataToFile(correctedData, "corrected.yaml")

}
