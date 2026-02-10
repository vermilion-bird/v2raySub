package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"v2raysub/config"
	"v2raysub/joke"
	"v2raysub/xui"

	"gopkg.in/yaml.v3"
)

const (
	//
	tickerInterval = 1 * 24 * 60 * time.Minute
	// 60 * time.Second
	keepDays               = 7
	serverPort             = ":8889"
	configFilePath         = "./config/config.yaml"
	outputFilePath         = "./v2rayConfig.txt"
	autoOutputFilePath     = "./autoV2rayConfig.txt"
	clashOutputFilePath    = "./clashConfig.yaml"
	autoClashOutputFilePath = "./autoClashConfig.yaml"
)

func main() {
	go startWebServer()
	scheduleTask()
	select {}
}

func scheduleTask() {
	ticker := time.NewTicker(tickerInterval)
	var wg sync.WaitGroup

	go func() {
		defer wg.Done()
		for range ticker.C {
			log.Println("定时任务执行:", time.Now())
			task()
		}
	}()

	// Add a WaitGroup to wait for goroutines to finish
	wg.Add(1)
	wg.Wait()
}

func reverseStringSlice(slice []string) []string {
	// Get the length of the slice
	length := len(slice)

	// Create a new slice to store the reversed elements
	reversedSlice := make([]string, length)

	// Iterate over the original slice in reverse order and populate the new slice
	for i, value := range slice {
		reversedSlice[length-1-i] = value
	}

	return reversedSlice
}

func task() {
	configs, err := config.LoadConfigs(configFilePath)
	if err != nil {
		log.Printf("Error loading configs: %v", err)
		return
	}

	currentTime := time.Now()
	boundaryDate := currentTime.Add(-time.Duration(keepDays) * 24 * time.Hour)
	dateString := boundaryDate.Format("20060102")
	currentString := currentTime.Format("20060102")

	var records []string
	for _, v := range configs {
		ips := extractIP(v.Host)[0]
		cookie, err := xui.Login(v.Host, v.User, v.Passwd)
		if err != nil {
			log.Printf("Error logging in to %s: %v", v.Host, err)
			continue
		}

		instances, err := xui.ListInstance(v.Host, cookie)
		if err != nil {
			log.Printf("Error listing instances for %s: %v", v.Host, err)
			continue
		}

		ports := getValidPorts(instances, dateString, v.Host, cookie)

		name := v.Country + joke.GenerateMythicalName() + ":" + currentString
		newPort := generateRandomPort(ports)
		xui.AddInstance(v.Host, name, strconv.Itoa(newPort), cookie)

		instancesWrite, _ := xui.ListInstance(v.Host, cookie)
		records = append(records, generateConfigRecords(instancesWrite, ips)...)
	}
	records = reverseStringSlice(records)
	writeFile(records)
	
	// Generate Clash configs
	generateAndWriteClashConfigs(configs)
}

func getValidPorts(instances []xui.V2rayInstance, dateString, host, cookie string) []int {
	var ports []int
	for _, v := range instances {
		dates := strings.Split(v.Remark, ":")
		date := dates[len(dates)-1]
		_, err := time.Parse("20060102", date)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			continue
		}
		if date < dateString {
			xui.DelInstance(host, v.ID, cookie)
		} else {
			ports = append(ports, v.Port)
		}
	}
	return ports
}

func generateRandomPort(ports []int) int {
	newPort := rand.Intn(60000-10000) + 10000
	for isInSlice(newPort, ports) {
		newPort = rand.Intn(60000-10000) + 10000
	}
	return newPort
}

func generateConfigRecords(instances []xui.V2rayInstance, ips string) []string {
	var records []string
	for _, v := range instances {
		config := VPNConfig{
			V:    "2",
			PS:   v.Remark,
			Add:  ips,
			Port: v.Port,
			ID:   "c3700e59-b55b-4db7-c836-c7c9b4c7d607",
			AID:  0,
			Net:  "ws",
			Type: "none",
			Host: "",
			Path: "/",
			TLS:  "none",
		}
		jsonString, _ := json.Marshal(config)
		encodedString := base64.StdEncoding.EncodeToString([]byte(jsonString))
		encodedString = "vmess://" + encodedString
		records = append(records, encodedString)
	}
	return records
}

func writeFile(records []string) {
	file, err := os.OpenFile(autoOutputFilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return
	}
	defer file.Close()

	err = file.Truncate(0)
	if err != nil {
		log.Printf("Error truncating file: %v", err)
		return
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		log.Printf("Error seeking file: %v", err)
		return
	}

	log.Println("File content cleared successfully.")

	writer := bufio.NewWriter(file)

	for _, line := range records {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			log.Printf("Error writing to file: %v", err)
			return
		}
	}

	err = writer.Flush()
	if err != nil {
		log.Printf("Error flushing buffer: %v", err)
		return
	}

	log.Println("File written successfully.")
}

func startWebServer() {
	http.HandleFunc("/me/ZzooG**oyqYVdrx", handler)
	http.HandleFunc("/me/ZzooGdW**qYVdrx/clash", clashHandler)
	err := http.ListenAndServe(serverPort, nil)
	if err != nil {
		log.Fatal("服务器启动失败:", err)
	}
	
}

func handler(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile(outputFilePath)
	if err != nil {
		log.Println("读取文件出错:", err)
		return
	}
	autoContent, err := ioutil.ReadFile(autoOutputFilePath)
	if err != nil {
		log.Println("读取文件出错:", err)
		return
	}

	text := string(content) + "\n" + string(autoContent)
	encoded := base64.StdEncoding.EncodeToString([]byte(text))

	fmt.Fprint(w, encoded)
}

func clashHandler(w http.ResponseWriter, r *http.Request) {
	// Set response header to indicate YAML content
	w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
	
	// Merge both Clash configs
	var config1 ClashConfig
	var config2 ClashConfig
	
	// Read manual config file (optional)
	content, err := ioutil.ReadFile(clashOutputFilePath)
	if err != nil {
		log.Println("手动 Clash 配置文件不存在，跳过:", err)
	} else if len(content) > 0 {
		err = yaml.Unmarshal(content, &config1)
		if err != nil {
			log.Println("解析 Clash 配置出错:", err)
		}
	}
	
	// Read auto-generated config file (optional)
	autoContent, err := ioutil.ReadFile(autoClashOutputFilePath)
	if err != nil {
		log.Println("自动 Clash 配置文件不存在，跳过:", err)
	} else if len(autoContent) > 0 {
		err = yaml.Unmarshal(autoContent, &config2)
		if err != nil {
			log.Println("解析自动 Clash 配置出错:", err)
		}
	}

	// Merge proxies
	mergedConfig := ClashConfig{
		Proxies: append(config1.Proxies, config2.Proxies...),
	}

	// If no proxies available, return empty config with message
	if len(mergedConfig.Proxies) == 0 {
		log.Println("警告: 没有可用的 Clash 代理配置")
		mergedConfig = ClashConfig{
			Proxies: []ClashProxy{},
		}
	}

	// Marshal to YAML
	yamlData, err := yaml.Marshal(mergedConfig)
	if err != nil {
		log.Println("生成 YAML 出错:", err)
		http.Error(w, "Failed to generate clash config", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(yamlData))
}

func isInSlice(target int, slice []int) bool {
	for _, value := range slice {
		if value == target {
			return true
		}
	}
	return false
}

type VPNConfig struct {
	V    string `json:"v"`
	PS   string `json:"ps"`
	Add  string `json:"add"`
	Port int    `json:"port"`
	ID   string `json:"id"`
	AID  int    `json:"aid"`
	Net  string `json:"net"`
	Type string `json:"type"`
	Host string `json:"host"`
	Path string `json:"path"`
	TLS  string `json:"tls"`
}

// ClashProxy represents a Clash proxy configuration
type ClashProxy struct {
	Name     string            `yaml:"name"`
	Type     string            `yaml:"type"`
	Server   string            `yaml:"server"`
	Port     int               `yaml:"port"`
	UUID     string            `yaml:"uuid"`
	AlterID  int               `yaml:"alterId"`
	Cipher   string            `yaml:"cipher"`
	TLS      bool              `yaml:"tls,omitempty"`
	Network  string            `yaml:"network,omitempty"`
	WSPath   string            `yaml:"ws-path,omitempty"`
	WSHeaders map[string]string `yaml:"ws-headers,omitempty"`
}

// ClashConfig represents a complete Clash configuration
type ClashConfig struct {
	Proxies []ClashProxy `yaml:"proxies"`
}

func extractIP(input string) []string {
	// 定义IPv4和IPv6正则表达式
	ipv4Pattern := `(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`
	ipv6Pattern := `([0-9a-fA-F:]+)`

	// 将正则表达式编译
	ipv4Regex := regexp.MustCompile(ipv4Pattern)
	ipv6Regex := regexp.MustCompile(ipv6Pattern)

	// 使用正则表达式查找匹配项
	ipv4Matches := ipv4Regex.FindAllString(input, -1)
	ipv6Matches := ipv6Regex.FindAllString(input, -1)

	// 合并结果
	matches := append(ipv4Matches, ipv6Matches...)

	return matches
}

// generateClashProxies converts V2Ray instances to Clash proxy format
func generateClashProxies(instances []xui.V2rayInstance, ips string) []ClashProxy {
	var proxies []ClashProxy
	nameCount := make(map[string]int)
	
	for _, v := range instances {
		// Generate unique name by checking for duplicates
		baseName := v.Remark
		uniqueName := baseName
		
		// If name already exists, append a counter
		if count, exists := nameCount[baseName]; exists {
			nameCount[baseName] = count + 1
			uniqueName = fmt.Sprintf("%s-%d", baseName, count+1)
		} else {
			nameCount[baseName] = 1
		}
		
		proxy := ClashProxy{
			Name:    uniqueName,
			Type:    "vmess",
			Server:  ips,
			Port:    v.Port,
			UUID:    "c3700e59-b55b-4db7-c836-c7c9b4c7d607",
			AlterID: 0,
			Cipher:  "auto",
			TLS:     false,
			Network: "ws",
			WSPath:  "/",
			WSHeaders: map[string]string{
				"Host": ips,
			},
		}
		proxies = append(proxies, proxy)
	}
	return proxies
}

// ensureUniqueProxyNames ensures all proxy names are globally unique
func ensureUniqueProxyNames(proxies []ClashProxy) []ClashProxy {
	nameCount := make(map[string]int)
	result := make([]ClashProxy, 0, len(proxies))
	
	for _, proxy := range proxies {
		baseName := proxy.Name
		uniqueName := baseName
		
		// Check if name already exists
		if count, exists := nameCount[baseName]; exists {
			nameCount[baseName] = count + 1
			// Append counter to make it unique
			uniqueName = fmt.Sprintf("%s-%d", baseName, count+1)
		} else {
			nameCount[baseName] = 1
		}
		
		// Create new proxy with unique name
		proxy.Name = uniqueName
		result = append(result, proxy)
	}
	
	return result
}

// generateAndWriteClashConfigs generates Clash format configs and writes to file
func generateAndWriteClashConfigs(configs []config.V2rayInstance) {
	currentTime := time.Now()
	boundaryDate := currentTime.Add(-time.Duration(keepDays) * 24 * time.Hour)
	dateString := boundaryDate.Format("20060102")

	var allProxies []ClashProxy
	
	for _, v := range configs {
		ips := extractIP(v.Host)[0]
		cookie, err := xui.Login(v.Host, v.User, v.Passwd)
		if err != nil {
			log.Printf("Error logging in to %s: %v", v.Host, err)
			continue
		}

		instances, err := xui.ListInstance(v.Host, cookie)
		if err != nil {
			log.Printf("Error listing instances for %s: %v", v.Host, err)
			continue
		}

		// Filter valid instances
		var validInstances []xui.V2rayInstance
		for _, inst := range instances {
			dates := strings.Split(inst.Remark, ":")
			date := dates[len(dates)-1]
			_, err := time.Parse("20060102", date)
			if err == nil && date >= dateString {
				validInstances = append(validInstances, inst)
			}
		}

		proxies := generateClashProxies(validInstances, ips)
		allProxies = append(allProxies, proxies...)
	}

	// Ensure all proxy names are unique globally
	allProxies = ensureUniqueProxyNames(allProxies)

	// Reverse the order to match the original behavior
	for i, j := 0, len(allProxies)-1; i < j; i, j = i+1, j-1 {
		allProxies[i], allProxies[j] = allProxies[j], allProxies[i]
	}

	clashConfig := ClashConfig{
		Proxies: allProxies,
	}

	writeClashFile(clashConfig, autoClashOutputFilePath)
}

// writeClashFile writes Clash config to YAML file
func writeClashFile(config ClashConfig, filepath string) {
	yamlData, err := yaml.Marshal(config)
	if err != nil {
		log.Printf("Error marshalling Clash config to YAML: %v", err)
		return
	}

	err = ioutil.WriteFile(filepath, yamlData, 0666)
	if err != nil {
		log.Printf("Error writing Clash config file: %v", err)
		return
	}

	log.Printf("Clash config file written successfully: %s", filepath)
}

// isInSlice and extractIP functions remain unchanged...
