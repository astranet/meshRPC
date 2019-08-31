package main

import (
	"strconv"
	"strings"
)

type PluginConfig struct {
	ListenHost   string
	ListenPort   int
	ClusterName  string
	ClusterNodes []string
	Debug        bool
}

func checkPluginConfig(cfg *PluginConfig) *PluginConfig {
	if cfg == nil {
		cfg = &PluginConfig{}
	}
	if len(cfg.ClusterName) == 0 {
		cfg.ClusterName = "vroomy-meshrpc"
	}
	return cfg
}

func configFromEnv(env map[string]string) (cfg *PluginConfig) {
	cfg = &PluginConfig{
		ListenHost:   env["meshrpc_host"],
		ListenPort:   toInt(env["meshrpc_port"], 11999),
		ClusterName:  env["meshrpc_cluster_name"],
		ClusterNodes: strings.Split(env["meshrpc_cluster_nodes"], ","),
		Debug:        toBool(env["meshrpc_debug"]),
	}
	return cfg
}

func toBool(s string) bool {
	switch strings.ToLower(s) {
	case "true", "1", "t", "yes":
		return true
	default:
		return false
	}
}

func toInt(s string, fallback int) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		return fallback
	}
	return v
}
