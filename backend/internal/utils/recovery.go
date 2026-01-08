package utils

import (
	"context"
	"crypto/rand"
	"fmt"
	"sync"

	"github.com/eeritvan/calendar/internal/sqlc"
	"golang.org/x/crypto/bcrypt"
)

const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

func GenerateRecoveryCode() string {
	b := make([]byte, 12)
	rand.Read(b)

	for i := range len(b) {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return fmt.Sprintf("%s-%s-%s", b[0:4], b[4:8], b[8:12])
}

func HashRecoveryCodes(codes []string) ([]string, error) {
	var wg sync.WaitGroup
	hashes := make([]string, len(codes))
	errors := make(chan error, len(codes))

	for i, code := range codes {
		wg.Add(1)

		go func(index int, plainText string) {
			defer wg.Done()

			hash, err := bcrypt.GenerateFromPassword([]byte(plainText), 12)
			if err != nil {
				errors <- err
				return
			}
			hashes[index] = string(hash)
		}(i, code)
	}

	wg.Wait()
	close(errors)

	if len(errors) > 0 {
		return nil, <-errors
	}

	return hashes, nil
}

func VerifyRecoveryCode(input string, codes []sqlc.GetUnusedRecoveryCodesRow) (int32, error) {
	var wg sync.WaitGroup
	resultChan := make(chan int32, 1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, code := range codes {
		wg.Add(1)

		go func(c sqlc.GetUnusedRecoveryCodesRow) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				return
			default:
				err := bcrypt.CompareHashAndPassword([]byte(c.CodeHash), []byte(input))
				if err == nil {
					select {
					case resultChan <- c.ID:
						cancel()
					default:
					}
				}
			}
		}(code)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	matchedID, ok := <-resultChan
	if !ok {
		return 0, fmt.Errorf("invalid code")
	}

	return matchedID, nil
}
