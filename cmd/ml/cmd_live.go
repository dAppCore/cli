package ml

import (
	"fmt"
	"os"

	"forge.lthn.ai/core/go/pkg/cli"
	"forge.lthn.ai/core/go-ai/ml"
)

const targetTotal = 15000

var liveCmd = &cli.Command{
	Use:   "live",
	Short: "Show live generation progress from InfluxDB",
	Long:  "Queries InfluxDB for real-time generation progress, worker breakdown, and domain/voice counts.",
	RunE:  runLive,
}

func runLive(cmd *cli.Command, args []string) error {
	influx := ml.NewInfluxClient(influxURL, influxDB)

	// Total completed generations
	total, err := influx.QueryScalar("SELECT count(DISTINCT i) AS n FROM gold_gen")
	if err != nil {
		return fmt.Errorf("live: query total: %w", err)
	}

	// Distinct domains and voices
	domains, err := influx.QueryScalar("SELECT count(DISTINCT d) AS n FROM gold_gen")
	if err != nil {
		return fmt.Errorf("live: query domains: %w", err)
	}

	voices, err := influx.QueryScalar("SELECT count(DISTINCT v) AS n FROM gold_gen")
	if err != nil {
		return fmt.Errorf("live: query voices: %w", err)
	}

	// Per-worker breakdown
	workers, err := influx.QueryRows("SELECT w, count(DISTINCT i) AS n FROM gold_gen GROUP BY w ORDER BY n DESC")
	if err != nil {
		return fmt.Errorf("live: query workers: %w", err)
	}

	pct := float64(total) / float64(targetTotal) * 100
	remaining := targetTotal - total

	fmt.Fprintln(os.Stdout, "Golden Set Live Status (from InfluxDB)")
	fmt.Fprintln(os.Stdout, "─────────────────────────────────────────────")
	fmt.Fprintf(os.Stdout, "  Total:     %d / %d (%.1f%%)\n", total, targetTotal, pct)
	fmt.Fprintf(os.Stdout, "  Remaining: %d\n", remaining)
	fmt.Fprintf(os.Stdout, "  Domains:   %d\n", domains)
	fmt.Fprintf(os.Stdout, "  Voices:    %d\n", voices)
	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, "  Workers:")
	for _, w := range workers {
		name := w["w"]
		n := w["n"]
		marker := ""
		if name == "migration" {
			marker = " (seed data)"
		}
		fmt.Fprintf(os.Stdout, "    %-20s %6s generations%s\n", name, n, marker)
	}

	return nil
}
