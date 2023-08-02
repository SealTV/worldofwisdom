package pow

import (
	"context"
	"testing"
	"time"
)

func TestIsValidPoW(t *testing.T) {
	tests := []struct {
		name       string
		hash       string
		difficulty int
		want       bool
	}{
		{
			name:       "valid",
			hash:       "00000000001231231231",
			difficulty: 10,
			want:       true,
		},
		{
			name:       "invalid",
			hash:       "somehash",
			difficulty: 10,
			want:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidPoW(tt.hash, tt.difficulty); got != tt.want {
				t.Errorf("IsValidPoW() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHashSHA256(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "valid",
			input: "someinput",
			want:  "8a6e2baea04dbec0023a2050b6d5244fcbd4bc1972a2db2061f4d11585015ad1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HashSHA256(tt.input); got != tt.want {
				t.Errorf("HashSHA256() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProofOfWork(t *testing.T) {
	tests := []struct {
		name       string
		challenge  string
		difficulty int
		timeout    time.Duration
		wantErr    bool
	}{
		{
			name:       "valid",
			challenge:  "somechallenge",
			difficulty: 1,
			timeout:    1 * time.Second,
			wantErr:    false,
		},
		{
			name:       "timeout",
			challenge:  "somechallenge",
			difficulty: 10,
			timeout:    100 * time.Microsecond,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()

			_, err := ProofOfWork(ctx, tt.challenge, tt.difficulty)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProofOfWork() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
