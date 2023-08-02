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

type PoW struct {
	difficulty int
}

func NewPoW(difficulty int) *PoW {
	return &PoW{
		difficulty: difficulty,
	}
}

func (p *PoW) GetChallenge() string {
	return getChallenge()
}

func (p *PoW) IsValid(input string) bool {
	hash := hashSHA256(input)
	return isValidPoW(hash, p.difficulty)
}
func (p *PoW) ProofOfWork(ctx context.Context, challenge string) (string, error) {
	return proofOfWork(ctx, challenge, p.difficulty)
}

// Proof of Work (PoW) function using SHA-256 as the hashing algorithm.
func proofOfWork(ctx context.Context, challenge string, difficulty int) (string, error) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	result := make(chan string)
	defer close(result)

	for {
		select {
		case <-ctx.Done():
			return "", fmt.Errorf("timeout")
		default:
			nonce := strconv.Itoa(random.Int()) // Generate a random nonce
			hash := hashSHA256(challenge + nonce)

			if isValidPoW(hash, difficulty) {
				return nonce, nil
			}
		}
	}
}

func getChallenge() string {
	return strconv.Itoa(rand.Int())
}

// Simple SHA-256 hashing function
func hashSHA256(input string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(input)))
}

// Validate if the hash has the required number of leading zeros.
func isValidPoW(hash string, difficulty int) bool {
	return strings.HasPrefix(hash, strings.Repeat("0", difficulty))
}
