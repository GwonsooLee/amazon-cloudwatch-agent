// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package procstat

import (
	"github.com/aws/amazon-cloudwatch-agent/translator"
)

type PidFinder struct{}

const keyPidFinder = "pid_finder"

func (t *PidFinder) ApplyRule(input interface{}) (returnKey string, returnVal interface{}) {
	returnKey, returnVal = translator.DefaultCase(keyPidFinder, "native", input)
	return
}

func init() {
	e := new(PidFinder)
	RegisterRule(keyPidFinder, e)
}
