package main

import (
	"admin-panel/redis_database"
	"log"
)

func onAppUp() {

	keys, err := redis_database.GetDataRedis().Scan("chain-client-status-report-*")
	if err != nil {
		log.Printf("Error scanning keys: %v", err)
		return
	}

	if len(keys) > 0 {
		log.Printf("Found %d keys matching pattern 'chain-client-status-report-*', deleting...", len(keys))

		for _, key := range keys {
			deleted, delErr := redis_database.GetDataRedis().Del(key)
			if delErr != nil {
				log.Printf("Error deleting key %s: %v", key, delErr)
			} else if deleted > 0 {
				log.Printf("Successfully deleted key: %s (confirmed)", key)
			} else {
				log.Printf("Failed to delete key: %s (Del returned 0)", key)
			}
		}

		log.Printf("Finished cleaning up keys with pattern 'chain-client-status-report-*'")
	} else {
		log.Printf("No keys found matching pattern 'chain-client-status-report-*'")
	}
}
