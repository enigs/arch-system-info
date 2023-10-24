package main

import (
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"os"
)

func main() {

	yellow := color.New(color.FgHiYellow).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()

	var sysInfo = [][]string{
		[]string{" ", " ", " "},
		[]string{yellow("███████╗███╗   ██╗██╗ ██████╗ ███████╗"), blue("Distro:"), GetDistro()},
		[]string{yellow("██╔════╝████╗  ██║██║██╔════╝ ██╔════╝"), blue("Kernel:"), GetKernelInfo()},
		[]string{yellow("█████╗  ██╔██╗ ██║██║██║  ███╗███████╗"), blue("User:"), GetCurrentUser()},
		[]string{yellow("██╔══╝  ██║╚██╗██║██║██║   ██║╚════██║"), blue("CPU:"), GetCpuInfo()},
		[]string{yellow("███████╗██║ ╚████║██║╚██████╔╝███████║"), blue("GPU:"), GetGpuInfo()},
		[]string{yellow("╚══════╝╚═╝  ╚═══╝╚═╝ ╚═════╝ ╚══════╝"), blue("RAM"), GetRamInfo()},
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
