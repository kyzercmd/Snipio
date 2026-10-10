package models

import (
	"testing"

	"github.com/kyzercmd/snipio/internal/assert"
)

func TestUserModelExists(t *testing.T) {

	if testing.Short() {
		t.Skip()
	}

	tests := []struct {
		name   string
		userID int
		want   bool
	}{
		{
			name:   "Exists",
			userID: 1,
			want:   true,
		},
		{
			name:   "Doesnt exist",
			userID: 2,
			want:   false,
		},
		{
			name:   "Zero ID",
			userID: 0,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)

			m := UserModel{db}

			exists, err := m.Exists(tt.userID)

			assert.Equal(t, exists, tt.want)
			assert.NilErr(t, err)
		})
	}
}
