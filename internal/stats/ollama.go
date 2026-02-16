package stats

import (
	"encoding/json"
	"net/http"
	"time"
)

// OllamaProcess matches the JSON response from the Ollama API.
// We are mapping specifically to the /api/ps endpoint to get running models.
type OllamaProcess struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	Details struct {
		ParameterSize string `json:"parameter_size"`     // e.g. "8B"
		Quantization  string `json:"quantization_level"` // e.g. "Q4_0"
	} `json:"details"`
	SizeVRAM int64 `json:"size_vram"` // VRAM usage in bytes
}

type OllamaResponse struct {
	Models []OllamaProcess `json:"models"`
}

// ModelInfo is our clean, internal struct for passing data to the UI.
type ModelInfo struct {
	Name         string
	ParamSize    string
	Quantization string
	VRAM         int64
}

func GetRunningOllamaModels() ([]ModelInfo, error) {
	// Set a short timeout so we don't block the UI if Ollama is slow/down
	client := http.Client{
		Timeout: 200 * time.Millisecond,
	}

	// NOTE: If you have a different device you want to log for whatever reason, this will need to be changed.
	resp, err := client.Get("http://127.0.0.1:11434/api/ps")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Map the API response to our internal struct
	var models []ModelInfo
	for _, m := range result.Models {
		models = append(models, ModelInfo{
			Name:         m.Name,
			ParamSize:    m.Details.ParameterSize,
			Quantization: m.Details.Quantization,
			VRAM:         m.SizeVRAM,
		})
	}
	return models, nil
}
