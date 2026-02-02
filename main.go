package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/its-the-vibe/RediFire/config"
	"github.com/redis/go-redis/v9"
	"google.golang.org/api/option"
)

// Record represents a message with payload and timestamp
type Record struct {
	Payload   map[string]interface{} `firestore:"payload"`
	Timestamp time.Time              `firestore:"timestamp"`
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}

func run() error {
	// Load configuration
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	log.Printf("Starting RediFire service with %d mappings", len(cfg.Mappings))

	// Create Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	defer redisClient.Close()

	// Test Redis connection
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}
	log.Println("Successfully connected to Redis")

	// Create Firestore client
	var firestoreOpts []option.ClientOption
	if cfg.Firestore.CredentialsFile != "" {
		firestoreOpts = append(firestoreOpts, option.WithCredentialsFile(cfg.Firestore.CredentialsFile))
	}

	firestoreClient, err := firestore.NewClient(ctx, cfg.Firestore.ProjectID, firestoreOpts...)
	if err != nil {
		return fmt.Errorf("failed to create Firestore client: %w", err)
	}
	defer firestoreClient.Close()
	log.Println("Successfully connected to Firestore")

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Create context for workers
	workerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Start workers for each mapping
	var wg sync.WaitGroup
	for _, mapping := range cfg.Mappings {
		wg.Add(1)
		go func(m config.Mapping) {
			defer wg.Done()
			worker(workerCtx, redisClient, firestoreClient, m)
		}(mapping)
		log.Printf("Started worker for %s -> %s", mapping.Source, mapping.Target)
	}

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutdown signal received, stopping workers...")
	cancel()
	wg.Wait()
	log.Println("All workers stopped, exiting")

	return nil
}

// computeSHA256 computes the SHA-256 hash of the given data
func computeSHA256(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func worker(ctx context.Context, redisClient *redis.Client, firestoreClient *firestore.Client, mapping config.Mapping) {
	log.Printf("[%s] Worker started", mapping.Source)

	for {
		select {
		case <-ctx.Done():
			log.Printf("[%s] Worker stopping", mapping.Source)
			return
		default:
			// Pop message from Redis list (blocking with timeout)
			result, err := redisClient.BLPop(ctx, 5*time.Second, mapping.Source).Result()
			if err != nil {
				if err == redis.Nil {
					// Timeout, continue loop
					continue
				}
				if ctx.Err() != nil {
					// Context cancelled
					return
				}
				log.Printf("[%s] Error popping from Redis: %v", mapping.Source, err)
				time.Sleep(time.Second)
				continue
			}

			// result[0] is the list name, result[1] is the value
			if len(result) < 2 {
				log.Printf("[%s] Invalid result from Redis: %v", mapping.Source, result)
				continue
			}

			message := result[1]

			// Compute SHA-256 hash of the message to use as document ID
			docID := computeSHA256(message)

			// Parse JSON payload
			var payload map[string]interface{}
			if err := json.Unmarshal([]byte(message), &payload); err != nil {
				log.Printf("[%s] Error parsing JSON: %v, message: %s", mapping.Source, err, message)
				continue
			}

			// Create record with timestamp
			record := Record{
				Payload:   payload,
				Timestamp: time.Now().UTC(),
			}

			// Write to Firestore with SHA-256 hash as document ID
			docRef := firestoreClient.Collection(mapping.Target).Doc(docID)
			if _, err := docRef.Set(ctx, record); err != nil {
				log.Printf("[%s] Error writing to Firestore: %v", mapping.Source, err)
				// In production, you might want to push back to Redis or a dead-letter queue
				continue
			}

			log.Printf("[%s] Successfully transferred message to %s (doc ID: %s)", mapping.Source, mapping.Target, docRef.ID)
		}
	}
}
