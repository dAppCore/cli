//go:build !(darwin && arm64)

package ml

import "forge.lthn.ai/core/go-ai/ml"

func createServeBackend() (ml.Backend, error) {
	return ml.NewHTTPBackend(apiURL, modelName), nil
}
