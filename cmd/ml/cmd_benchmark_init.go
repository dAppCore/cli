//go:build darwin && arm64

package ml

func init() {
	mlCmd.AddCommand(benchmarkCmd)
}
