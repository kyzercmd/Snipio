package main

import (
	"testing"
	"time"

	"github.com/kyzercmd/snipio/internal/assert"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		to   time.Time
		want string
	}{
		{
			name: "UTC",
			to:   time.Date(2000, 1, 1, 1, 1, 0, 0, time.UTC),
			want: "01 Jan 2000 at 01:01",
		},
		{
			name: "empty",
			to:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			to:   time.Date(2000, 1, 1, 1, 1, 0, 0, time.FixedZone("CET", int(1*time.Hour.Seconds()))),
			want: "01 Jan 2000 at 00:01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.to)
			assert.Equal(t, hd, tt.want)
		})
	}
}
