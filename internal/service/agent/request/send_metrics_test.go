package request

import "testing"

func TestSendMetricsHandler(t *testing.T) {
	type args struct {
		url         string
		contentType string
		body        []byte
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SendMetricsHandler(tt.args.url, tt.args.contentType, tt.args.body)
		})
	}
}
