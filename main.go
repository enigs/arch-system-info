package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
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
		return ""
	}
	return strings.TrimSpace(out.String())
}

func getCPUInfo() string {
	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return "Unknown"
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
			// Split the line on colon and return the right-hand side
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	if scanner.Err() != nil {
		return "Unknown"
	}
	return "Unknown"
}

func getDistro() string {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return "Unknown"
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			return strings.Trim(line[len("PRETTY_NAME="):], `"`)
		}
	}
	return "Unknown"
}

func getGPUInfo() string {
	cmd := exec.Command("lspci")
	output, err := cmd.Output()
	if err != nil {
		return "Unknown"
	}

	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "VGA compatible controller") {
			r := regexp.MustCompile(`GeForce.*?Ti`)
			match := r.FindString(line)
			if match != "" {
				return match
			}
		}
	}

	return "Unknown"
}

func getRAMInfo() string {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return "Unknown / Unknown"
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
		acceptableThreshold = 0.2 // 20%
		lowThreshold        = 0.1 // 10%
	)

	memFreeStr := fmt.Sprintf("%d MiB", memFree/1024)
	totalStr := color.New(color.FgHiWhite).Sprint(fmt.Sprintf("/ %d MiB", memTotal/1024))

	switch {
	case freePercentage > acceptableThreshold:
		return color.New(color.FgHiGreen).Sprint(memFreeStr) + " " + totalStr
	case freePercentage <= acceptableThreshold && freePercentage > lowThreshold:
		return memFreeStr + " " + totalStr // Default white for used memory
	default:
		return color.New(color.FgHiRed).Sprint(memFreeStr) + " " + totalStr
	}
}

func main() {

	currentUser, _ := user.Current()
	distro := getDistro()
	kernel := executeCommand("uname", "-r")

	yellow := color.New(color.FgYellow).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()

	var sysInfo = [][]string{
		[]string{" ", " ", " "},
		[]string{yellow("███████╗███╗   ██╗██╗ ██████╗ ███████╗"), blue("Distro:"), distro},
		[]string{yellow("██╔════╝████╗  ██║██║██╔════╝ ██╔════╝"), blue("Kernel:"), kernel},
		[]string{yellow("█████╗  ██╔██╗ ██║██║██║  ███╗███████╗"), blue("User:"), currentUser.Username},
		[]string{yellow("██╔══╝  ██║╚██╗██║██║██║   ██║╚════██║"), blue("CPU"), getCPUInfo()},
		[]string{yellow("███████╗██║ ╚████║██║╚██████╔╝███████║"), blue("GPU"), getGPUInfo()},
		[]string{yellow("╚══════╝╚═╝  ╚═══╝╚═╝ ╚═════╝ ╚══════╝"), blue("RAM"), getRAMInfo()},
		[]string{" ", " ", " "},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetColWidth(120)

	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(false)
	table.SetColumnSeparator("  ")
	table.SetAutoWrapText(false)
	table.AppendBulk(sysInfo)
	table.Render()
}
