package seed

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"os"

	authsqlc "github.com/DeRuina/KUHA-REST-API/internal/db/auth"
	"github.com/DeRuina/KUHA-REST-API/internal/store"
)

// Map of clients to their assigned roles
// {clientName: {roles}}
var clientsToSeed = map[string][]string{}

// Seed inserts predefined clients with hashed client_tokens into the database
func Seed(store store.Auth, db *sql.DB) {
	ctx := context.Background()
	queries := store.Queries()

	file, err := os.Create("client_tokens.txt")
	if err != nil {
		log.Fatalf("Failed to create client_tokens.txt: %v", err)
	}
	defer file.Close()

	for clientName, roles := range clientsToSeed {
		// Generate a secure token and its SHA-256 hash
		rawToken, hashedToken, err := generateSecureTokenPair(32)
		if err != nil {
			log.Printf("Failed to generate token for %s: %v", clientName, err)
			continue
		}

		// Insert client into the database
		err = queries.CreateClient(ctx, authsqlc.CreateClientParams{
			ClientName:  clientName,
			ClientToken: hashedToken,
			Role:        roles,
		})

		if err != nil {
			log.Printf("Error inserting client %s: %v", clientName, err)
			continue
		}

		// Save raw token to file for reference
		line := fmt.Sprintf("Client: %-15s Token: %s\n", clientName, rawToken)
		if _, err := file.WriteString(line); err != nil {
			log.Printf("Failed to write token to file for %s: %v", clientName, err)
		}

		fmt.Print(line)
	}

	log.Println("Seeding complete. Tokens written to client_tokens.txt")
}

// Create a random token and its SHA-256 hash.
func generateSecureTokenPair(length int) (rawToken string, hashedToken string, err error) {
	bytes := make([]byte, length)
	_, err = rand.Read(bytes)
	if err != nil {
		return "", "", err
	}

	raw := hex.EncodeToString(bytes)      // Original token
	hash := sha256.Sum256([]byte(raw))    // SHA-256 hash
	hashed := hex.EncodeToString(hash[:]) // Convert hash to hex

	return raw, hashed, nil
}
