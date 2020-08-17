package main

import (
	yaml "github.com/ghodss/yaml"
)

// RuleDef is a rule definition. A rule is composed of a set of checkItems.
type RuleDef struct {
	// Name is the name of the rule
	Name string `json:"name,omitempty"`
	// CheckItems is a list of checkItems
	CheckItems []string `json:"checkItems,omitempty"`
}

// RuleSet is created by the rule definition
type RuleSet struct {
	// Rule is the checking rules
	Rules []RuleDef `json:"rules,omitempty"`
	// CorrectorNames the of the correctors
	CorrectorNames []string `json:"correctors,omitempty"`
}

// CorrectorSetDef is
type CorrectorSetDef struct {
	CorrectorNames []string `json:"correctors,omitempty"`
}

// ParserRuleSetConfig is to parese the rule definition file and create the RuleSet
func ParserRuleSetConfig(configFile string) ([]*Rule, []Corrector, error) {
	data, err := loadDataFromFile(configFile)
	if err != nil {
		return nil, nil, err
	}
	ruleSet := &RuleSet{}
	err = yaml.Unmarshal(data, ruleSet)
	if err != nil {
		return nil, nil, err
	}
	rules := []*Rule{}
	correctors := []Corrector{}
	for _, ruleDef := range ruleSet.Rules {
		rule := &Rule{}
		rule.Name = ruleDef.Name
		checkers := []Checker{}
		for _, checkName := range ruleDef.CheckItems {
			checker := Checkers[checkName]
			checkers = append(checkers, checker)
			corrector, ok := checker.(Corrector)
			if ok {
				correctors = append(correctors, corrector)
			}
		}
		rule.Checkers = checkers
		rules = append(rules, rule)
	}
	return rules, correctors, nil
}
