package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type CronField struct {
	Name  string
	Min   int
	Max   int
	Index int
}

var cronFields = []CronField{
	{Name: "minute", Min: 0, Max: 59, Index: 0},
	{Name: "hour", Min: 0, Max: 23, Index: 1},
	{Name: "day of month", Min: 1, Max: 31, Index: 2},
	{Name: "month", Min: 1, Max: 12, Index: 3},
	{Name: "day of week", Min: 0, Max: 6, Index: 4},
}

var parseCmd = &cobra.Command{
	Use:   "parse [cron expression]",
	Short: "Parse a cron expression",
	Long: `Parse a cron expression and expand each field to show the times at which it will run.
The cron expression should be in the standard format with five time fields
(minute, hour, day of month, month, and day of week) plus a command.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		parseCronExpression(args[0])
	},
}

func init() {
	rootCmd.AddCommand(parseCmd)
}

func parseCronExpression(cronExpr string) {
	fields := strings.Fields(cronExpr)
	if len(fields) < 6 {
		fmt.Println("Invalid cron expression. Expected at least 6 fields.")
		return
	}

	for i, field := range cronFields {
		fmt.Printf("%-14s", field.Name)
		expandField(fields[i], field)
		fmt.Println()
	}

	fmt.Printf("%-14s%s\n", "command", strings.Join(fields[5:], " "))
}

func expandField(field string, cronField CronField) {
	var expanded []int

	if field == "*" {
		expanded = expandRange(cronField.Min, cronField.Max)
	} else {
		parts := strings.Split(field, ",")
		for _, part := range parts {
			if strings.Contains(part, "/") {
				expanded = append(expanded, expandStep(part, cronField)...)
			} else if strings.Contains(part, "-") {
				expanded = append(expanded, expandRange(parseRange(part))...)
			} else {
				num, _ := strconv.Atoi(part)
				if num >= cronField.Min && num <= cronField.Max {
					expanded = append(expanded, num)
				}
			}
		}
	}

	expanded = uniqueSort(expanded)

	for i, num := range expanded {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(num)
	}
}

func expandStep(field string, cronField CronField) []int {
	parts := strings.Split(field, "/")
	var start, end int

	if parts[0] == "*" {
		start, end = cronField.Min, cronField.Max
	} else if strings.Contains(parts[0], "-") {
		start, end = parseRange(parts[0])
	} else {
		start, _ = strconv.Atoi(parts[0])
		end = cronField.Max
	}

	step, _ := strconv.Atoi(parts[1])
	return expandRange(start, end, step)
}

func parseRange(field string) (int, int) {
	parts := strings.Split(field, "-")
	start, _ := strconv.Atoi(parts[0])
	end, _ := strconv.Atoi(parts[1])
	return start, end
}

func expandRange(start, end int, step ...int) []int {
	var expanded []int
	stepValue := 1
	if len(step) > 0 {
		stepValue = step[0]
	}
	for i := start; i <= end; i += stepValue {
		expanded = append(expanded, i)
	}
	return expanded
}

func uniqueSort(nums []int) []int {
	if len(nums) == 0 {
		return nums
	}

	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i] > nums[j] {
				nums[i], nums[j] = nums[j], nums[i]
			}
		}
	}

	unique := nums[:1]
	for i := 1; i < len(nums); i++ {
		if nums[i] != nums[i-1] {
			unique = append(unique, nums[i])
		}
	}

	return unique
}
