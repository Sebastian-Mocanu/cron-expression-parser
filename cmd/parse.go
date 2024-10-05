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
		fmt.Println("Error: Invalid cron expression. Expected at least 6 fields.")
		fmt.Println("Correct format: minute hour day_of_month month day_of_week command")
		fmt.Println("Example: */15 0 1,15 * 1-5 /usr/bin/find")
		return
	}

	hasError := false
	for i, field := range cronFields {
		fmt.Printf("%-14s", field.Name)
		expanded, err := expandField(fields[i], field)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			hasError = true
		} else {
			fmt.Println(expanded)
		}
	}

	if !hasError {
		fmt.Printf("%-14s%s\n", "command", strings.Join(fields[5:], " "))
	}
}

func expandField(field string, cronField CronField) (string, error) {
	if field == "*" {
		return formatExpanded(expandRange(cronField.Min, cronField.Max)), nil
	}

	var expanded []int
	parts := strings.Split(field, ",")
	for _, part := range parts {
		if strings.Contains(part, "/") {
			values, err := expandStep(part, cronField)
			if err != nil {
				return "", err
			}
			expanded = append(expanded, values...)
		} else if strings.Contains(part, "-") {
			start, end, err := parseRange(part, cronField)
			if err != nil {
				return "", err
			}
			expanded = append(expanded, expandRange(start, end)...)
		} else {
			num, err := strconv.Atoi(part)
			if err != nil {
				return "", fmt.Errorf("invalid value %s: %v", part, err)
			}
			if num < cronField.Min || num > cronField.Max {
				return "", fmt.Errorf("value %d out of range (allowed range: %d-%d)", num, cronField.Min, cronField.Max)
			}
			expanded = append(expanded, num)
		}
	}

	return formatExpanded(uniqueSort(expanded)), nil
}

func expandStep(field string, cronField CronField) ([]int, error) {
	parts := strings.Split(field, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid step format: %s", field)
	}

	var start, end int
	var err error

	if parts[0] == "*" {
		start, end = cronField.Min, cronField.Max
	} else if strings.Contains(parts[0], "-") {
		start, end, err = parseRange(parts[0], cronField)
		if err != nil {
			return nil, err
		}
	} else {
		start, err = strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid start value in step expression: %v", err)
		}
		if start < cronField.Min || start > cronField.Max {
			return nil, fmt.Errorf("start value %d out of range (allowed range: %d-%d)", start, cronField.Min, cronField.Max)
		}
		end = cronField.Max
	}

	step, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid step value: %v", err)
	}

	if step <= 0 {
		return nil, fmt.Errorf("step value must be positive")
	}

	if step > cronField.Max-cronField.Min+1 {
		return nil, fmt.Errorf("step value %d is too large for range %d-%d", step, cronField.Min, cronField.Max)
	}

	return expandRange(start, end, step), nil
}

func parseRange(field string, cronField CronField) (int, int, error) {
	parts := strings.Split(field, "-")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid range format: %s", field)
	}

	start, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid start of range: %v", err)
	}

	end, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid end of range: %v", err)
	}

	if start > end {
		return 0, 0, fmt.Errorf("invalid range: start (%d) is greater than end (%d)", start, end)
	}

	if start < cronField.Min || end > cronField.Max {
		return 0, 0, fmt.Errorf("range %d-%d out of bounds (allowed range: %d-%d)", start, end, cronField.Min, cronField.Max)
	}

	return start, end, nil
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

func formatExpanded(expanded []int) string {
	var result strings.Builder
	for i, num := range expanded {
		if i > 0 {
			result.WriteString(" ")
		}
		result.WriteString(strconv.Itoa(num))
	}
	return result.String()
}
