package mcp

import (
    "encoding/json"
    "os"
    "path/filepath"
    "runtime"
)

const (
    // DefaultPort is the default port for MCP servers
    DefaultPort = 49200
    // DefaultHost is the default host for MCP servers
    DefaultHost = "localhost"
)

// GetDefaultConfigPath returns the default path for Claude Desktop config
func GetDefaultConfigPath() string {
    var configDir string
    
    if runtime.GOOS == "windows" {
        configDir = filepath.Join(os.Getenv("APPDATA"), "claude-desktop")
    } else {
        homeDir, err := os.UserHomeDir()
        if err != nil {
            homeDir = "."
        }
        configDir = filepath.Join(homeDir, ".config", "claude-desktop")
    }
    
    return filepath.Join(configDir, "config.json")
}

// EnsureConfigDir ensures the config directory exists
func EnsureConfigDir() error {
    configPath := GetDefaultConfigPath()
    return os.MkdirAll(filepath.Dir(configPath), 0755)
}

// ServerConfig represents the MCP server configuration
type ServerConfig struct {
    Name         string            `json:"name"`
    Type         string            `json:"type"`
    Host         string            `json:"host"`
    Port         int              `json:"port"`
    Capabilities map[string]bool   `json:"capabilities"`
}

// ClaudeConfig represents the Claude Desktop configuration file
type ClaudeConfig struct {
    Servers []ServerConfig `json:"servers"`
}

// LoadClaudeConfig loads the existing Claude Desktop configuration
func LoadClaudeConfig() (*ClaudeConfig, error) {
    configPath := GetDefaultConfigPath()
    
    data, err := os.ReadFile(configPath)
    if err != nil {
        if os.IsNotExist(err) {
            return &ClaudeConfig{
                Servers: make([]ServerConfig, 0),
            }, nil
        }
        return nil, err
    }

    var config ClaudeConfig
    if err := json.Unmarshal(data, &config); err != nil {
        return nil, err
    }

    return &config, nil
}

// RegisterServer registers or updates a server in the Claude Desktop configuration
func RegisterServer(server ServerConfig) error {
    if err := EnsureConfigDir(); err != nil {
        return err
    }

    config, err := LoadClaudeConfig()
    if err != nil {
        return err
    }

    // Update existing server or add new one
    found := false
    for i, s := range config.Servers {
        if s.Name == server.Name {
            config.Servers[i] = server
            found = true
            break
        }
    }

    if !found {
        config.Servers = append(config.Servers, server)
    }

    // Save updated configuration
    data, err := json.MarshalIndent(config, "", "    ")
    if err != nil {
        return err
    }

    return os.WriteFile(GetDefaultConfigPath(), data, 0644)
}
