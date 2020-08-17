package main

// Checkers are all the checkers being ready to use.
var Checkers map[string]Checker

func init() {
	// The checkers which are ready to use should be initialized here.
	Checkers = make(map[string]Checker)
	Checkers["RunningOnDifferentNodes"] = &RunningOnDifferentNodes{}
	Checkers["WithGracefulTermination"] = &WithGracefulTermination{}
	Checkers["WithHealthCheck"] = &WithHealthCheck{}
	Checkers["WithResourceRequestAndLimit"] = &WithResourceRequestAndLimit{}
	Checkers["WithReadiness"] = &WithReadiness{}

}

var spotCheckSet = &Rule{
	Name: "spot",
	Checkers: []Checker{
		&RunningOnDifferentNodes{},
		&WithGracefulTermination{},
		&WithHealthCheck{},
		&WithResourceRequestAndLimit{},
	},
}

var normalCheckSet = &Rule{
	Name: "normal",
	Checkers: []Checker{
		&WithHealthCheck{},
		&WithResourceRequestAndLimit{},
		&WithReadiness{},
	},
}

// initialize the default rule set
var ruleSet = []*Rule{
	spotCheckSet,
	normalCheckSet,
}
