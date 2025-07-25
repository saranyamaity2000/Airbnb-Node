package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ProxyObject defines the structure for each service in the proxy.yaml
type ProxyObject struct {
	ServiceName string `yaml:"serviceName"`
	Alias       string `yaml:"alias"`
	BaseURL     string `yaml:"baseUrl"`
}

// Define structure to match YAML format
type ProxyConfig struct {
	Services []ProxyObject `yaml:"services"`
}

// AvailableProxyServers holds the list of proxy configurations
var AvailableProxyServers []ProxyObject
var proxyFileName = "proxy.yaml"

func init() {
	// Load from proxy.yaml file
	fmt.Printf("reading Proxy.yaml\n")
	data, err := os.ReadFile(proxyFileName)
	if err != nil {
		fmt.Printf("Error Occured while reading %s\n", proxyFileName)
		return
	}

	var config ProxyConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		fmt.Printf("Error Occured while parsing the YAML content in Go Structs\n")
		return
	}

	for _, proxyObj := range config.Services {
		if proxyObj.ServiceName == "" || proxyObj.Alias == "" || proxyObj.BaseURL == "" {
			fmt.Printf("ServiceName, Alias, and BaseURL cannot be empty\n")
			continue
		}
		AvailableProxyServers = append(AvailableProxyServers, proxyObj)
	}

	if len(AvailableProxyServers) == 0 {
		fmt.Printf("No Proxies Found in %s\n", proxyFileName)
	} else {
		fmt.Printf("%d Proxies Found in %s\n", len(AvailableProxyServers), proxyFileName)
	}
}
