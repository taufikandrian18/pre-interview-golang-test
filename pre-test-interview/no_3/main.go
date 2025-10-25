package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "server" {
		runServer()
	} else {
		runDemo()
	}
}

func runDemo() {
	fmt.Println("=== Cache Implementation Demo ===")
	fmt.Println("💡 Tip: Run with 'go run . server' to start the interactive HTTP API server")
	fmt.Println()

	fmt.Println("1. Simple In-Memory Cache:")
	simpleCache := NewSimpleCache()

	err := simpleCache.Set("user_123", "John Doe")
	if err != nil {
		fmt.Printf("   Error setting user: %v\n", err)
		return
	}
	err = simpleCache.Set("age_123", 30)
	if err != nil {
		fmt.Printf("   Error setting age: %v\n", err)
		return
	}

	if name, exists, err := simpleCache.Get("user_123"); err == nil && exists {
		fmt.Printf("   User name: %s\n", name)
	}

	if age, exists, err := simpleCache.Get("age_123"); err == nil && exists {
		fmt.Printf("   User age: %d\n", age)
	}

	err = simpleCache.Delete("age_123")
	if err != nil {
		fmt.Printf("   Error deleting age: %v\n", err)
	} else if _, exists, err := simpleCache.Get("age_123"); err == nil && !exists {
		fmt.Println("   Age deleted successfully")
	}

	fmt.Println()

	fmt.Println("2. TTL Cache (with 2 second expiration):")
	ttlCache, err := NewTTLCache(2 * time.Second)
	if err != nil {
		fmt.Printf("   Error creating TTLCache: %v\n", err)
		return
	}
	defer ttlCache.Close()

	err = ttlCache.Set("temp_data", "This will expire")
	if err != nil {
		fmt.Printf("   Error setting temp_data: %v\n", err)
		return
	}
	fmt.Println("   Set temp_data with 2s TTL")

	if value, exists, err := ttlCache.Get("temp_data"); err == nil && exists {
		fmt.Printf("   Immediately after set: %s\n", value)
	}

	fmt.Println("   Waiting 3 seconds for expiration...")
	time.Sleep(3 * time.Second)

	if _, exists, err := ttlCache.Get("temp_data"); err == nil && !exists {
		fmt.Println("   After expiration: data no longer exists")
	}

	fmt.Println("\n3. TTL Reset on Update:")
	err = ttlCache.Set("reset_demo", "Original value")
	if err != nil {
		fmt.Printf("   Error setting reset_demo: %v\n", err)
		return
	}
	fmt.Println("   Set reset_demo with 2s TTL")

	time.Sleep(1 * time.Second)
	fmt.Println("   After 1 second, updating value...")
	err = ttlCache.Set("reset_demo", "Updated value")
	if err != nil {
		fmt.Printf("   Error updating reset_demo: %v\n", err)
		return
	}

	time.Sleep(1 * time.Second)
	if value, exists, err := ttlCache.Get("reset_demo"); err == nil && exists {
		fmt.Printf("   After another second: %s (TTL was reset)\n", value)
	}

	fmt.Println("\n=== Demo Complete ===")
	fmt.Println("💡 To test interactively, run: go run . server")
}
