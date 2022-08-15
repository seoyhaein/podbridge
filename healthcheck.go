package podbridge

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/containers/image/v5/manifest"
	"github.com/containers/podman/v4/pkg/specgen"
	"github.com/pkg/errors"
)

func SetHealthChecker(s *specgen.SpecGenerator, inCmd, interval string, retries uint, timeout, startPeriod string) (*specgen.SpecGenerator, error) {

	if s == nil {
		return nil, fmt.Errorf("specgen.SpecGenerator is nil")
	}
	healthConfig, err := makeHealthCheckFromCli(inCmd, interval, retries, timeout, startPeriod)
	if err != nil {
		return nil, err
	}
	s.HealthConfig = healthConfig
	return s, nil
}

func makeHealthCheckFromCli(inCmd, interval string, retries uint, timeout, startPeriod string) (*manifest.Schema2HealthConfig, error) {
	cmdArr := []string{}
	isArr := true
	err := json.Unmarshal([]byte(inCmd), &cmdArr) // array unmarshalling
	if err != nil {
		cmdArr = strings.SplitN(inCmd, " ", 2) // default for compat
		isArr = false
	}
	// Every healthcheck requires a command
	if len(cmdArr) == 0 {
		return nil, errors.New("Must define a healthcheck command for all healthchecks")
	}

	var concat string
	if cmdArr[0] == "CMD" || cmdArr[0] == "none" { // this is for compat, we are already split properly for most compat cases
		cmdArr = strings.Fields(inCmd)
	} else if cmdArr[0] != "CMD-SHELL" { // this is for podman side of things, won't contain the keywords
		if isArr && len(cmdArr) > 1 { // an array of consecutive commands
			cmdArr = append([]string{"CMD"}, cmdArr...)
		} else { // one singular command
			if len(cmdArr) == 1 {
				concat = cmdArr[0]
			} else {
				concat = strings.Join(cmdArr[0:], " ")
			}
			cmdArr = append([]string{"CMD-SHELL"}, concat)
		}
	}

	if cmdArr[0] == "none" { // if specified to remove healtcheck
		cmdArr = []string{"NONE"}
	}

	// healthcheck is by default an array, so we simply pass the user input
	hc := manifest.Schema2HealthConfig{
		Test: cmdArr,
	}

	if interval == "disable" {
		interval = "0"
	}
	intervalDuration, err := time.ParseDuration(interval)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid healthcheck-interval")
	}

	hc.Interval = intervalDuration

	if retries < 1 {
		return nil, errors.New("healthcheck-retries must be greater than 0")
	}
	hc.Retries = int(retries)
	timeoutDuration, err := time.ParseDuration(timeout)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid healthcheck-timeout")
	}
	if timeoutDuration < time.Duration(1) {
		return nil, errors.New("healthcheck-timeout must be at least 1 second")
	}
	hc.Timeout = timeoutDuration

	startPeriodDuration, err := time.ParseDuration(startPeriod)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid healthcheck-start-period")
	}
	if startPeriodDuration < time.Duration(0) {
		return nil, errors.New("healthcheck-start-period must be 0 seconds or greater")
	}
	hc.StartPeriod = startPeriodDuration

	return &hc, nil
}
