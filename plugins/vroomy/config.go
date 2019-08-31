package main

import "encoding/json"

type PluginConfig struct {
	ListenHost   string   `json:"meshrpc_host"`
	ListenPort   string   `json:"meshrpc_port"`
	ClusterName  string   `json:"meshrpc_cluster_name"`
	ClusterNodes []string `json:"meshrpc_cluster_nodes"`
	Debug        bool     `json:"meshrpc_debug"`
}

func checkPluginConfig(cfg *PluginConfig) *PluginConfig {
	if cfg == nil {
		cfg = &PluginConfig{}
	}
	if len(cfg.ClusterName) == 0 {
		cfg.ClusterName = "vroomy-meshrpc"
	}
	if len(cfg.ClusterNodes) == 0 {
		// detect from env maybe
	}
	return cfg
}

func configFromEnv(env map[string]string) (cfg *PluginConfig) {
	data, _ := json.Marshal(env)
	_ = json.Unmarshal(data, &cfg)
	return cfg
}
