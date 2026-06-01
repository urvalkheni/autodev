package cmd

import (
	"context"
	"net/http"
	"os"
	"time"
)

// trackCLIMetric sends a non-blocking background ping to a public counter to track activations.
// Users can opt-out by setting DO_NOT_TRACK=1 or AUTODEV_TELEMETRY_DISABLED=1.
func trackCLIMetric(event string) {
	if os.Getenv("AUTODEV_TELEMETRY_DISABLED") == "1" || os.Getenv("DO_NOT_TRACK") == "1" {
		return
	}

	go func() {
		// Set a tight timeout to ensure we never block the command completion
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		url := "https://api.counterapi.dev/v1/heetmehta18-autodev/cli_" + event + "/up"
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return
		}

		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			_ = resp.Body.Close()
		}
	}()
}
