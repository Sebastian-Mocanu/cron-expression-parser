package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestParseCronExpression(t *testing.T) {
	type testCase struct {
		name           string
		cronExpression string
		expectedOutput string
	}

	tests := []testCase{
		{
			name:           "Standard cron expression",
			cronExpression: "*/15 0 1,15 * 1-5 /usr/bin/find",
			expectedOutput: `minute        0 15 30 45
hour          0
day of month  1 15
month         1 2 3 4 5 6 7 8 9 10 11 12
day of week   1 2 3 4 5
command       /usr/bin/find
`,
		},
		{
			name:           "All stars",
			cronExpression: "* * * * * /usr/bin/find",
			expectedOutput: `minute        0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31 32 33 34 35 36 37 38 39 40 41 42 43 44 45 46 47 48 49 50 51 52 53 54 55 56 57 58 59
hour          0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23
day of month  1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31
month         1 2 3 4 5 6 7 8 9 10 11 12
day of week   0 1 2 3 4 5 6
command       /usr/bin/find
`,
		},
		{
			name:           "Complex cron expression",
			cronExpression: "0-5,10-59/5 1-4,22,23 1,15 1-6/2 * /usr/local/bin/complex_command",
			expectedOutput: `minute        0 1 2 3 4 5 10 15 20 25 30 35 40 45 50 55
hour          1 2 3 4 22 23
day of month  1 15
month         1 3 5
day of week   0 1 2 3 4 5 6
command       /usr/local/bin/complex_command
`,
		},
		{
			name:           "Every 5 minutes",
			cronExpression: "*/5 * * * * /scripts/every-five-minutes.sh",
			expectedOutput: `minute        0 5 10 15 20 25 30 35 40 45 50 55
hour          0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23
day of month  1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31
month         1 2 3 4 5 6 7 8 9 10 11 12
day of week   0 1 2 3 4 5 6
command       /scripts/every-five-minutes.sh
`,
		},
		{
			name:           "Every hour at 30 minutes",
			cronExpression: "30 * * * * /scripts/hourly-half-past.sh",
			expectedOutput: `minute        30
hour          0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23
day of month  1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31
month         1 2 3 4 5 6 7 8 9 10 11 12
day of week   0 1 2 3 4 5 6
command       /scripts/hourly-half-past.sh
`,
		},
		{
			name:           "Every day at midnight",
			cronExpression: "0 0 * * * /scripts/daily-midnight.sh",
			expectedOutput: `minute        0
hour          0
day of month  1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31
month         1 2 3 4 5 6 7 8 9 10 11 12
day of week   0 1 2 3 4 5 6
command       /scripts/daily-midnight.sh
`,
		},
		{
			name:           "Every Sunday at 6:30 PM",
			cronExpression: "30 18 * * 0 /scripts/sunday-evening.sh",
			expectedOutput: `minute        30
hour          18
day of month  1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31
month         1 2 3 4 5 6 7 8 9 10 11 12
day of week   0
command       /scripts/sunday-evening.sh
`,
		},
		{
			name:           "Every 15 minutes during work hours",
			cronExpression: "*/15 9-17 * * 1-5 /scripts/work-hours-check.sh",
			expectedOutput: `minute        0 15 30 45
hour          9 10 11 12 13 14 15 16 17
day of month  1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31
month         1 2 3 4 5 6 7 8 9 10 11 12
day of week   1 2 3 4 5
command       /scripts/work-hours-check.sh
`,
		},
		{
			name:           "First day of every month at noon",
			cronExpression: "0 12 1 * * /scripts/monthly-report.sh",
			expectedOutput: `minute        0
hour          12
day of month  1
month         1 2 3 4 5 6 7 8 9 10 11 12
day of week   0 1 2 3 4 5 6
command       /scripts/monthly-report.sh
`,
		},
		{
			name:           "Every quarter hour",
			cronExpression: "0,15,30,45 * * * * /scripts/quarter-hourly.sh",
			expectedOutput: `minute        0 15 30 45
hour          0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23
day of month  1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31
month         1 2 3 4 5 6 7 8 9 10 11 12
day of week   0 1 2 3 4 5 6
command       /scripts/quarter-hourly.sh
`,
		},
	}

	passCount := 0
	failCount := 0

	for _, test := range tests {

		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		parseCronExpression(test.cronExpression)

		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		io.Copy(&buf, r)
		actualOutput := buf.String()

		if actualOutput != test.expectedOutput {
			failCount++
			fmt.Printf(`
---------------------------------
Test Failed: %s
 cron expression: %s
 expected output:
%s
 actual output:
%s`, test.name, test.cronExpression, test.expectedOutput, actualOutput)
		} else {
			passCount++
			fmt.Printf(`
---------------------------------
Test Passed: %s
 cron expression: %s
 expected output:
%s
 actual output:
%s`, test.name, test.cronExpression, test.expectedOutput, actualOutput)
		}
	}

	fmt.Println("---------------------------------")
	fmt.Printf("%d passed, %d failed\n", passCount, failCount)

	if failCount > 0 {
		t.Errorf("%d tests failed", failCount)
	}
}
