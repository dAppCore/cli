package ml

import (
	"fmt"
	"os"

	"forge.lthn.ai/core/go/pkg/cli"
	"forge.lthn.ai/core/go-ai/ml"
)

var expandStatusCmd = &cli.Command{
	Use:   "expand-status",
	Short: "Show expansion pipeline progress",
	Long:  "Queries DuckDB for expansion prompts, generated responses, scoring status, and overall pipeline progress.",
	RunE:  runExpandStatus,
}

func runExpandStatus(cmd *cli.Command, args []string) error {
	path := dbPath
	if path == "" {
		path = os.Getenv("LEM_DB")
	}
	if path == "" {
		return fmt.Errorf("--db or LEM_DB required")
	}

	db, err := ml.OpenDB(path)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	fmt.Fprintln(os.Stdout, "LEM Expansion Pipeline Status")
	fmt.Fprintln(os.Stdout, "==================================================")

	// Expansion prompts
	total, pending, err := db.ExpansionPromptCounts()
	if err != nil {
		fmt.Fprintln(os.Stdout, "  Expansion prompts:  not created (run: normalize)")
		return nil
	}
	fmt.Fprintf(os.Stdout, "  Expansion prompts:  %d total, %d pending\n", total, pending)

	// Generated responses
	generated, models, err := db.ExpansionRawCounts()
	if err != nil {
		generated = 0
		fmt.Fprintln(os.Stdout, "  Generated:          0 (run: core ml expand)")
	} else if len(models) > 0 {
		modelStr := ""
		for i, m := range models {
			if i > 0 {
				modelStr += ", "
			}
			modelStr += fmt.Sprintf("%s: %d", m.Name, m.Count)
		}
		fmt.Fprintf(os.Stdout, "  Generated:          %d (%s)\n", generated, modelStr)
	} else {
		fmt.Fprintf(os.Stdout, "  Generated:          %d\n", generated)
	}

	// Scored
	scored, hPassed, jScored, jPassed, err := db.ExpansionScoreCounts()
	if err != nil {
		fmt.Fprintln(os.Stdout, "  Scored:             0 (run: score --tier 1)")
	} else {
		fmt.Fprintf(os.Stdout, "  Heuristic scored:   %d (%d passed)\n", scored, hPassed)
		if jScored > 0 {
			fmt.Fprintf(os.Stdout, "  Judge scored:       %d (%d passed)\n", jScored, jPassed)
		}
	}

	// Pipeline progress
	if total > 0 && generated > 0 {
		genPct := float64(generated) / float64(total) * 100
		fmt.Fprintf(os.Stdout, "\n  Progress:           %.1f%% generated\n", genPct)
	}

	// Golden set context
	golden, err := db.GoldenSetCount()
	if err == nil && golden > 0 {
		fmt.Fprintf(os.Stdout, "\n  Golden set:         %d / %d\n", golden, targetTotal)
		if generated > 0 {
			fmt.Fprintf(os.Stdout, "  Combined:           %d total examples\n", golden+generated)
		}
	}

	return nil
}
