package checker

import "regexp"

var ipReg = regexp.MustCompile(`[\d.:]+`)
