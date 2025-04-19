package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

// DomainConfig represents the configuration for a domain and its targets
type DomainConfig struct {
	Domain  string
	Targets []string
}

// TargetIndices keeps track of the last used target index for each domain
var targetIndices = make(map[string]int)
var mutex = &sync.Mutex{}

// parseDomainMapping parses a domain mapping string into a DomainConfig
func parseDomainMapping(mappingStr string) DomainConfig {
	parts := strings.Split(mappingStr, "->")
	if len(parts) < 2 {
		return DomainConfig{}
	}
	domain := parts[0]
	targets := strings.Split(parts[1], ",")
	return DomainConfig{
		Domain:  domain,
		Targets: targets,
	}
}

// getDomainConfigs retrieves all domain configurations from environment variables
func getDomainConfigs() []DomainConfig {
	var configs []DomainConfig
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "DOMAIN_MAPPING_") {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 {
				configs = append(configs, parseDomainMapping(parts[1]))
			}
		}
	}
	return configs
}

// findMatchingConfig finds the configuration that matches the given host
func findMatchingConfig(host string, configs []DomainConfig) *DomainConfig {
	for _, config := range configs {
		if host == config.Domain || host == fmt.Sprintf("%s:%s", config.Domain, os.Getenv("PORT")) {
			return &config
		}
	}
	return nil
}

// getNextTarget gets the next target in round-robin fashion
func getNextTarget(domain string, targets []string) string {
	mutex.Lock()
	defer mutex.Unlock()
	
	currentIndex := targetIndices[domain]
	nextIndex := (currentIndex + 1) % len(targets)
	targetIndices[domain] = nextIndex
	
	return targets[currentIndex]
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		host := r.Host
		path := r.URL.Path
		
		// Get domain configurations
		configs := getDomainConfigs()
		
		// Find matching configuration
		config := findMatchingConfig(host, configs)
		if config != nil {
			// Get next target
			target := getNextTarget(config.Domain, config.Targets)
			
			// Construct target URL
			targetURL, err := url.Parse(target)
			if err != nil {
				http.Error(w, "Invalid target URL", http.StatusInternalServerError)
				return
			}
			
			// Add path to target URL if preserve_path is enabled
			if os.Getenv("PRESERVE_PATH") == "true" {
				targetURL.Path = path
			}
			
			// Add cache prevention parameter if enabled
			query := targetURL.Query()
			if os.Getenv("ENABLE_TIMESTAMP") == "true" {
				query.Set("_t", fmt.Sprintf("%d", time.Now().UnixNano()))
			}
			
			// Add referral information if enabled
			if os.Getenv("INCLUDE_REFERRAL") == "true" {
				query.Set("ref", r.Host)
			}
			
			targetURL.RawQuery = query.Encode()
			
			// Redirect
			http.Redirect(w, r, targetURL.String(), http.StatusFound)
			return
		}
		
		// Default response
		fmt.Fprintf(w, "Domain redirect service is running")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	fmt.Printf("Server starting on port %s\n", port)
	http.ListenAndServe(":"+port, nil)
}