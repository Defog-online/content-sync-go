package main

import (
	"testing"
	"time"
)

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Duration
		wantErr bool
	}{
		{
			name:    "valid simple seconds",
			input:   "PT10S",
			want:    10 * time.Second,
			wantErr: false,
		},
		{
			name:    "valid minutes and seconds",
			input:   "PT1M30S",
			want:    1*time.Minute + 30*time.Second,
			wantErr: false,
		},
		{
			name:    "valid hours minutes and seconds",
			input:   "PT1H2M10S",
			want:    1*time.Hour + 2*time.Minute + 10*time.Second,
			wantErr: false,
		},
		{
			name:    "valid only hours",
			input:   "PT2H",
			want:    2 * time.Hour,
			wantErr: false,
		},
		{
			name:    "valid only minutes",
			input:   "PT45M",
			want:    45 * time.Minute,
			wantErr: false,
		},
		{
			name:    "lowercase input",
			input:   "pt5s",
			want:    5 * time.Second,
			wantErr: false,
		},
		{
			name:    "zero duration",
			input:   "PT0S",
			want:    0,
			wantErr: false,
		},
		{
			name:    "invalid format",
			input:   "INVALID",
			want:    0,
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			want:    0,
			wantErr: true,
		},
		{
			name:    "only PT",
			input:   "PT",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDuration(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDuration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseDurationFilterLongForm(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool // true if > 60 seconds
	}{
		{
			name:     "60 seconds exactly (not long form)",
			input:    "PT60S",
			expected: false,
		},
		{
			name:     "61 seconds (long form)",
			input:    "PT61S",
			expected: true,
		},
		{
			name:     "2 minutes (long form)",
			input:    "PT2M",
			expected: true,
		},
		{
			name:     "30 seconds (short form)",
			input:    "PT30S",
			expected: false,
		},
		{
			name:     "1 hour (long form)",
			input:    "PT1H",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			duration, err := parseDuration(tt.input)
			if err != nil {
				t.Fatalf("parseDuration() error = %v", err)
			}
			isLongForm := duration.Seconds() > 60
			if isLongForm != tt.expected {
				t.Errorf("isLongForm = %v, want %v", isLongForm, tt.expected)
			}
		})
	}
}
