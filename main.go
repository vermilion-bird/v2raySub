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
)

const (
	//
	tickerInterval     = 1 * 24 * 60 * time.Minute
	keepDays           = 14
	serverPort         = ":8889"
	configFilePath     = "./config/config.yaml"
	outputFilePath     = "./v2rayConfig.txt"
	autoOutputFilePath = "./autoV2rayConfig.txt"
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

	writeFile(records)
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
	http.HandleFunc("/me/rVMhVnCboe75XPMxVw9aVAN1u6wHZ", handler)
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

// isInSlice and extractIP functions remain unchanged...
