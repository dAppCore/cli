//go:build !(darwin && arm64 && mlx)

package ml

import "forge.lthn.ai/core/go/pkg/ml"

func createServeBackend() (ml.Backend, error) {
	return ml.NewHTTPBackend(apiURL, modelName), nil
}
