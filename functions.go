package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"strconv"
	"strings"
)

func executeCommand(command string, args ...string) string {
	cmd := exec.Command(command, args...)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return color.New(color.FgHiRed).Sprint(err.Error())
	}

	return color.New(color.FgHiWhite).Sprint(strings.TrimSpace(out.String()))
}

func GetCpuInfo() string {
	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return color.New(color.FgHiRed).Sprint("Unknown")
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			println("Unable to close file")
		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "model name") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				return color.New(color.FgHiWhite).Sprint(strings.TrimSpace(parts[1]))
			}
		}
	}

	return color.New(color.FgHiRed).Sprint("Unknown")
}

func GetCurrentUser() string {
	currentUser, err := user.Current()
	if err != nil {
		return color.New(color.FgHiRed).Sprint("Unknown")
	}

	return color.New(color.FgHiWhite).Sprint(currentUser.Username)
}

func GetDistro() string {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return color.New(color.FgHiRed).Sprint("Unknown")
	}

	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			return color.New(color.FgHiWhite).Sprint(strings.Trim(line[len("PRETTY_NAME="):], `"`))
		}
	}

	return color.New(color.FgHiRed).Sprint("Unknown")
}

func GetKernelInfo() string {
	return color.New(color.FgHiWhite).Sprint(executeCommand("uname", "-r"))
}

func GetGpuInfo() string {
	cmd := exec.Command("lspci")
	output, err := cmd.Output()
	if err != nil {
		return color.New(color.FgHiRed).Sprint("Unknown")
	}

	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "VGA compatible controller") {
			r := regexp.MustCompile(`GeForce.*?Ti`)
			match := r.FindString(line)
			if match != "" {
				return color.New(color.FgHiWhite).Sprint(match)
			}
		}
	}

	return color.New(color.FgHiRed).Sprint("Unknown")
}

func GetRamInfo() string {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return color.New(color.FgHiRed).Sprint("Unknown")
	}

	var memTotal, memFree int
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal") {
			memTotalStr := strings.TrimSpace(strings.Split(line, ":")[1])
			memTotalStr = strings.TrimRight(memTotalStr, " kB")
			memTotal, _ = strconv.Atoi(memTotalStr)
		} else if strings.HasPrefix(line, "MemFree") {
			memFreeStr := strings.TrimSpace(strings.Split(line, ":")[1])
			memFreeStr = strings.TrimRight(memFreeStr, " kB")
			memFree, _ = strconv.Atoi(memFreeStr)
		}
	}

	freePercentage := float64(memFree) / float64(memTotal)

	// Define thresholds. Adjust these based on what you consider acceptable and low.
	const (
		acceptableThreshold = 0.6 // 60%
		lowThreshold        = 0.3 // 30%
	)

	memFreeStr := fmt.Sprintf("%d MiB", memFree/1024)
	totalStr := color.New(color.FgHiWhite).Sprint(fmt.Sprintf("/ %d MiB", memTotal/1024))

	switch {
	case freePercentage > acceptableThreshold:
		return color.New(color.FgHiRed).Sprint(memFreeStr) + " " + color.New(color.FgHiWhite).Sprint(totalStr)
	case freePercentage <= acceptableThreshold && freePercentage > lowThreshold:
		return color.New(color.FgYellow).Sprint(memFreeStr) + " " + color.New(color.FgHiWhite).Sprint(totalStr)
	default:
		return color.New(color.FgHiGreen).Sprint(memFreeStr) + " " + color.New(color.FgHiWhite).Sprint(totalStr)
	}
}
