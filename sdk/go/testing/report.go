// Copyright 2026 PulumiCost/FinFocus Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package testing provides a comprehensive testing framework for PulumiCost plugins.
// This file implements reporting functionality for conformance test results.
package testing

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// Report format constants.
const (
	// ReportVersion is the current version of the conformance report schema.
	ReportVersion = "1.0.0"

	// Report column widths for formatted output.
	reportPluginNameWidth = 56
	reportMethodWidth     = 60
	reportErrorMsgWidth   = 58
)

// ToJSON serializes the conformance result to JSON format.
func (r *ConformanceResult) ToJSON() ([]byte, error) {
	// Ensure string representations are set
	r.LevelAchievedStr = r.LevelAchieved.String()
	r.DurationStr = r.Duration.String()

	return json.MarshalIndent(r, "", "  ")
}

// PrintReport prints a human-readable summary of the conformance result to stdout.
func PrintReport(result *ConformanceResult) {
	PrintReportTo(result, os.Stdout)
}

// PrintReportTo prints a human-readable summary of the conformance result to the writer.
//
//nolint:gocognit // Report printing inherently requires sequential output
func PrintReportTo(result *ConformanceResult, w io.Writer) {
	fmt.Fprintf(w, "
")
	fmt.Fprintf(w, "╔══════════════════════════════════════════════════════════════════╗
")
	fmt.Fprintf(w, "║               Plugin Conformance Test Report                     ║
")
	fmt.Fprintf(w, "╠══════════════════════════════════════════════════════════════════╣
")
	fmt.Fprintf(w, "║ Plugin: %-56s ║
", truncate(result.PluginName, reportPluginNameWidth))
	fmt.Fprintf(w, "║ Level Achieved: %-48s ║
", result.LevelAchieved.String())
	fmt.Fprintf(w, "║ Duration: %-54s ║
", result.Duration.String())
	fmt.Fprintf(w, "╠══════════════════════════════════════════════════════════════════╣
")
	fmt.Fprintf(w, "║ Summary                                                          ║
")
	fmt.Fprintf(w, "╠──────────────────────────────────────────────────────────────────╣
")
	fmt.Fprintf(w, "║   Total:   %-53d ║
", result.Summary.Total)
	fmt.Fprintf(w, "║   Passed:  %-53d ║
", result.Summary.Passed)
	fmt.Fprintf(w, "║   Failed:  %-53d ║
", result.Summary.Failed)
	fmt.Fprintf(w, "║   Skipped: %-53d ║
", result.Summary.Skipped)
	fmt.Fprintf(w, "╠══════════════════════════════════════════════════════════════════╣
")
	fmt.Fprintf(w, "║ Categories                                                       ║
")
	fmt.Fprintf(w, "╠──────────────────────────────────────────────────────────────────╣
")

	// Print category results
	categories := []TestCategory{
		CategorySpecValidation,
		CategoryRPCCorrectness,
		CategoryPerformance,
		CategoryConcurrency,
	}

	for _, cat := range categories {
		if catResult, ok := result.Categories[cat]; ok {
			status := "✓"
			if catResult.Failed > 0 {
				status = "✗"
			}
			fmt.Fprintf(w, "║   %s %-16s  Passed: %2d  Failed: %2d  Skipped: %2d      ║
",
				status, cat.String(), catResult.Passed, catResult.Failed, catResult.Skipped)
		}
	}

	fmt.Fprintf(w, "╠══════════════════════════════════════════════════════════════════╣
")

	// Print failed tests if any
	if result.Summary.Failed > 0 {
		fmt.Fprintf(w, "║ Failed Tests                                                     ║
")
		fmt.Fprintf(w, "╠──────────────────────────────────────────────────────────────────╣
")
		for _, catResult := range result.Categories {
			for _, testResult := range catResult.Results {
				if !testResult.Success {
					fmt.Fprintf(w, "║   • %-60s ║
", truncate(testResult.Method, reportMethodWidth))
					if testResult.Error != nil {
						errMsg := truncate(testResult.Error.Error(), reportErrorMsgWidth)
						fmt.Fprintf(w, "║     Error: %-53s ║
", errMsg)
					}
				}
			}
		}
		fmt.Fprintf(w, "╠══════════════════════════════════════════════════════════════════╣
")
	}

	// Final verdict
	if result.Passed() {
		fmt.Fprintf(w, "║                         ✓ ALL TESTS PASSED                       ║
")
	} else {
		fmt.Fprintf(w, "║                         ✗ SOME TESTS FAILED                      ║
")
	}
	fmt.Fprintf(w, "╚══════════════════════════════════════════════════════════════════╝
")
	fmt.Fprintf(w, "
")
}

// truncate truncates a string to the specified length.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// FormatCategoryResults formats category results as a multi-line string.
func FormatCategoryResults(categories map[TestCategory]*CategoryResult) string {
	var sb strings.Builder

	for cat, result := range categories {
		sb.WriteString(fmt.Sprintf("%s: %d passed, %d failed, %d skipped
",
			cat.String(), result.Passed, result.Failed, result.Skipped))
	}

	return sb.String()
}
