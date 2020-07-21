// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

// +build linux

package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/aws/amazon-cloudwatch-agent/translator/cmdutil"
	"github.com/aws/amazon-cloudwatch-agent/translator/config"
	"github.com/aws/amazon-cloudwatch-agent/translator/context"
)

const (
	AGENT_DIR_LINUX = "/opt/aws/amazon-cloudwatch-agent"

	JSON_DIR_LINUX = "amazon-cloudwatch-agent.d"

	TRANSLATOR_BINARY_LINUX = "config-translator"
	AGENT_BINARY_LINUX      = "amazon-cloudwatch-agent"
)

func startAgent(writer io.WriteCloser) error {
	if os.Getenv(config.RUN_IN_CONTAINER) == config.RUN_IN_CONTAINER_TRUE {
		cmd := exec.Command(agentBinaryPath, "-config", tomlConfigPath, "-envconfig", envConfigPath,
			"-pidfile", AGENT_DIR_LINUX+"/var/amazon-cloudwatch-agent.pid")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		return errors.New("agent unexpectedly exited")
	}

	mergedJsonConfigMap, err := generateMergedJsonConfigMap()
	if err != nil {
		log.Printf("E! Failed to generate merged json config: %v ", err)
		return err
	}

	_, err = cmdutil.ChangeUser(mergedJsonConfigMap)
	if err != nil {
		log.Printf("E! Failed to ChangeUser: %v ", err)
		return err
	}

	name, err := exec.LookPath(agentBinaryPath)
	if err != nil {
		log.Printf("E! Faield to lookpath: %v ", err)
		return err
	}

	if err := writer.Close(); err != nil {
		log.Printf("E! Cannot close the log file, ERROR is %v ", err)
		return err
	}

	// linux command has pid passed while windows does not
	agentCmd := []string{agentBinaryPath, "-config", tomlConfigPath, "-envconfig", envConfigPath,
		"-pidfile", AGENT_DIR_LINUX + "/var/amazon-cloudwatch-agent.pid"}
	if err = syscall.Exec(name, agentCmd, os.Environ()); err != nil {
		// log file is closed, so use fmt here
		fmt.Printf("E! Exec failed: %v \n", err)
		return err
	}

	return nil
}

func generateMergedJsonConfigMap() (map[string]interface{}, error) {
	ctx := context.CurrentContext()
	ctx.SetOs(config.OS_TYPE_LINUX)
	ctx.SetInputJsonFilePath(jsonConfigPath)
	ctx.SetInputJsonDirPath(jsonDirPath)
	ctx.SetMultiConfig("remove")
	return cmdutil.GenerateMergedJsonConfigMap(ctx)
}

func init() {
	jsonConfigPath = AGENT_DIR_LINUX + "/etc/" + JSON
	jsonDirPath = AGENT_DIR_LINUX + "/etc/" + JSON_DIR_LINUX
	envConfigPath = AGENT_DIR_LINUX + "/etc/" + ENV
	tomlConfigPath = AGENT_DIR_LINUX + "/etc/" + TOML
	commonConfigPath = AGENT_DIR_LINUX + "/etc/" + COMMON_CONFIG

	agentLogFilePath = AGENT_DIR_LINUX + "/logs/" + AGENT_LOG_FILE

	translatorBinaryPath = AGENT_DIR_LINUX + "/bin/" + TRANSLATOR_BINARY_LINUX
	agentBinaryPath = AGENT_DIR_LINUX + "/bin/" + AGENT_BINARY_LINUX
}
