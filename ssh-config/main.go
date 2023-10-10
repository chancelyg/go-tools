package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func parseSSHConfig(configPath string) ([]map[string]string, error) {
	configFile, err := os.Open(configPath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer configFile.Close()

	configDicts := []map[string]string{}
	currentDict := map[string]string{}

	scanner := bufio.NewScanner(configFile)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "Host ") {
			if len(currentDict) > 0 {
				configDicts = append(configDicts, currentDict)
			}
			currentDict = map[string]string{}
			currentDict["Host"] = strings.TrimSpace(strings.TrimPrefix(line, "Host "))
		} else {
			parts := strings.SplitN(line, " ", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				currentDict[key] = value
			}
		}
	}

	if len(currentDict) > 0 {
		configDicts = append(configDicts, currentDict)
	}

	return configDicts, nil
}

func printSSHConfig(configs []map[string]string) {
	fmt.Println("************************ Hi, Welcome to use Go-SSH Tool *****************************")
	fmt.Println()
	fmt.Println("+-----+------------------------------+-------------------------+------------------------------------------+")
	fmt.Println("| id  | Host                         | username                | address                                  |")
	fmt.Println("+-----+------------------------------+-------------------------+------------------------------------------+")
	for i, config := range configs {
		fmt.Printf("| %-3d | %-28s | %-23s | %-40s |\n", i+1, config["Host"], config["User"], config["HostName"])
	}
	fmt.Println("+-----+------------------------------+-------------------------+------------------------------------------+")
	fmt.Println()
	fmt.Println("Tips: Press a number between 1 and", len(configs)-1, "to select the host to connect, or \"q\" to quit.")
	fmt.Println()
}

func main() {
	configPath := flag.String("c", "$HOME/.ssh/config", "Path to SSH config file")
	flag.Parse()

	expandedPath := os.ExpandEnv(*configPath)
	configs, err := parseSSHConfig(expandedPath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	printSSHConfig(configs)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("# ")
		scanner.Scan()
		input := scanner.Text()

		if input == "q" {
			return
		}

		id := -1
		fmt.Sscanf(input, "%d", &id)
		if id >= 1 && id <= len(configs) {
			sshConfig := configs[id-1]
			var args []string
			for key, value := range sshConfig {
				if key == "Host" || key == "User" || key == "HostName" {
					continue
				}
				arg := fmt.Sprintf("-o %s=%s", key, value)
				args = append(args, arg)
			}
			cmdArgs := append(args, fmt.Sprintf("%s@%s", sshConfig["User"], sshConfig["HostName"]))
			sshCmd := exec.Command("ssh", cmdArgs...)
			sshCmd.Stdout = os.Stdout
			sshCmd.Stdin = os.Stdin
			sshCmd.Stderr = os.Stderr
			err := sshCmd.Run()
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(0)
			}
		} else {
			fmt.Println("Error: Invalid input")
		}
		os.Exit(1)

		printSSHConfig(configs)
	}
}
