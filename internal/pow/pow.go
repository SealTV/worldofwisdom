package pow

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const POW_DIFFICULTY = 3

type PoW struct{}

func NewPoW() *PoW {
	return &PoW{}
}

func (p *PoW) GetChallenge() string {
	return GetChallenge()
}

func (p *PoW) IsValid(input string) bool {
	hash := HashSHA256(input)
	return IsValidPoW(hash, POW_DIFFICULTY)
}

// Proof of Work (PoW) function using SHA-256 as the hashing algorithm.
func ProofOfWork(ctx context.Context, challenge string, difficulty int) (string, error) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	result := make(chan string)
	defer close(result)

	for {
		select {
		case <-ctx.Done():
			return "", fmt.Errorf("timeout")
		default:
			nonce := strconv.Itoa(random.Int()) // Generate a random nonce
			hash := HashSHA256(challenge + nonce)

			if IsValidPoW(hash, difficulty) {
				return nonce, nil
			}
		}
	}
}

func GetChallenge() string {
	return strconv.Itoa(rand.Int())
}

// Simple SHA-256 hashing function
func HashSHA256(input string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(input)))
}

// Validate if the hash has the required number of leading zeros.
func IsValidPoW(hash string, difficulty int) bool {
	return strings.HasPrefix(hash, strings.Repeat("0", difficulty))
}
