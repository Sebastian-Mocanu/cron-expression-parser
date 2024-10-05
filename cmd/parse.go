package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

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
		os.Exit(1)
	}

	fieldNames := []string{"minute", "hour", "day of month", "month", "day of week", "command"}
	for i, field := range fields {
		if i < 5 {
			fmt.Printf("%-14s", fieldNames[i])
			expandField(field, i)
		} else {
			fmt.Printf("%-14s%s\n", fieldNames[5], strings.Join(fields[5:], " "))
			break
		}
	}
}

func expandField(field string, fieldIndex int) {
	var expanded []string
	fieldRange := getFieldRange(fieldIndex)

	if field == "*" {
		expanded = expandRange(fieldRange[0], fieldRange[1])
	} else if strings.Contains(field, "/") {
		parts := strings.Split(field, "/")
		step, _ := strconv.Atoi(parts[1])
		for i := fieldRange[0]; i <= fieldRange[1]; i += step {
			expanded = append(expanded, strconv.Itoa(i))
		}
	} else if strings.Contains(field, "-") {
		parts := strings.Split(field, "-")
		start, _ := strconv.Atoi(parts[0])
		end, _ := strconv.Atoi(parts[1])
		expanded = expandRange(start, end)
	} else if strings.Contains(field, ",") {
		expanded = strings.Split(field, ",")
	} else {
		expanded = []string{field}
	}

	fmt.Println(strings.Join(expanded, " "))
}

func expandRange(start, end int) []string {
	var expanded []string
	for i := start; i <= end; i++ {
		expanded = append(expanded, strconv.Itoa(i))
	}
	return expanded
}

func getFieldRange(fieldIndex int) []int {
	switch fieldIndex {
	case 0: // minute
		return []int{0, 59}
	case 1: // hour
		return []int{0, 23}
	case 2: // day of month
		return []int{1, 31}
	case 3: // month
		return []int{1, 12}
	case 4: // day of week
		return []int{0, 6}
	default:
		return []int{0, 0}
	}
}
