package config

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var globalConfig *Config

// ================================
// Config Struct
// ================================
type Config struct {
	Environment string `yaml:"environment"`
	// crawler config
	CrawlerParallelism int    `yaml:"crawler_parallelism"`
	CrawlerDelay       string `yaml:"crawler_delay"`
	CrawlerRandomDelay string `yaml:"crawler_randomDelay"`
	// log config
	LogLevel string `yaml:"log_level"`
}

// ================================
// Load Config
// ================================
func Load() {
	abs := ""
	data := []byte{}
	abs1, err1 := filepath.Abs("env/config.yaml")
	if err1 != nil {
		log.Fatalf("无法解析配置文件路径: %v", err1)
	}

	abs2, err1 := filepath.Abs("configs/config.yaml")
	if err1 != nil {
		log.Fatalf("无法解析配置文件路径: %v", err1)
	}

	data1, err := ioutil.ReadFile(abs1)
	if err != nil {
		log.Printf("配置文件 %s 不存在, 尝试读取 %s", abs1, abs2)
		data2, err := ioutil.ReadFile(abs2)
		if err != nil {
			log.Fatalf("配置文件 %s 不存在: %v", abs2, err)
		}
		abs = abs2
		data = data2
	} else {
		abs = abs1
		data = data1
	}

	cfg := &Config{}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		log.Fatalf("配置文件解析失败: %v", err)
	}

	globalConfig = cfg
	log.Printf("配置加载完成: %s", abs)
}

// ================================
// Get Global Config
// ================================
func Get() *Config {
	if globalConfig == nil {
		Load()
	}
	return globalConfig
}
