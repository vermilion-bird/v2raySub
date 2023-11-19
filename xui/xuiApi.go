package xui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type ListInstanceResponse struct {
	Success bool            `json:"success"`
	Msg     string          `json:"msg"`
	Obj     []V2rayInstance `json:"obj"`
}

type V2rayInstance struct {
	ID             int    `json:"id"`
	Up             int64  `json:"up"`
	Down           int64  `json:"down"`
	Total          int    `json:"total"`
	Remark         string `json:"remark"`
	Enable         bool   `json:"enable"`
	ExpiryTime     int    `json:"expiryTime"`
	Listen         string `json:"listen"`
	Port           int    `json:"port"`
	Protocol       string `json:"protocol"`
	Settings       string `json:"settings"`
	StreamSettings string `json:"streamSettings"`
	Tag            string `json:"tag"`
	Sniffing       string `json:"sniffing"`
}

type V2raySettings struct {
	Clients                   []V2rayClient `json:"clients"`
	DisableInsecureEncryption bool          `json:"disableInsecureEncryption"`
}

type V2rayClient struct {
	ID      string `json:"id"`
	AlterID int    `json:"alterId"`
}

const (
	acceptHeader      = "application/json, text/plain, */*"
	contentTypeHeader = "application/x-www-form-urlencoded; charset=UTF-8"
	acceptLanguage    = "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7"
	connectionHeader  = "keep-alive"
	xRequestedWith    = "XMLHttpRequest"
	userAgentHeader   = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"
)

func buildRequest(method, url string, payload []byte, cookie string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", acceptHeader)
	req.Header.Set("Accept-Language", acceptLanguage)
	req.Header.Set("Connection", connectionHeader)
	req.Header.Set("Content-Type", contentTypeHeader)
	req.Header.Set("Cookie", "session="+cookie)
	req.Header.Set("X-Requested-With", xRequestedWith)
	req.Header.Set("User-Agent", userAgentHeader)

	return req, nil
}

func makeRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	return client.Do(req)
}

func Login(host, user, pwd string) (string, error) {
	url := fmt.Sprintf("%s/login", host)
	payload := []byte(fmt.Sprintf("username=%s&password=%s", user, pwd))
	req, err := buildRequest("POST", url, payload, "")
	if err != nil {
		return "", err
	}

	resp, err := makeRequest(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Process the response body here if needed.

	return resp.Cookies()[0].Value, nil
}

func AddInstance(host, name, port, cookie string) error {
	url := fmt.Sprintf("%s/xui/inbound/add", host)
	config := fmt.Sprintf("remark=%s&listen=&port=%s&", name, port)
	payload := []byte(config + `up=0&down=0&total=0&enable=true&expiryTime=0&protocol=vmess&settings=%7B%0A%20%20%22clients%22%3A%20%5B%0A%20%20%20%20%7B%0A%20%20%20%20%20%20%22id%22%3A%20%22c3700e59-b55b-4db7-c836-c7c9b4c7d607%22%2C%0A%20%20%20%20%20%20%22alterId%22%3A%200%0A%20%20%20%20%7D%0A%20%20%5D%2C%0A%20%20%22disableInsecureEncryption%22%3A%20false%0A%7D&streamSettings=%7B%0A%20%20%22network%22%3A%20%22ws%22%2C%0A%20%20%22security%22%3A%20%22none%22%2C%0A%20%20%22wsSettings%22%3A%20%7B%0A%20%20%20%20%22path%22%3A%20%22%2F%22%2C%0A%20%20%20%20%22headers%22%3A%20%7B%7D%0A%20%20%7D%0A%7D&sniffing=%7B%0A%20%20%22enabled%22%3A%20true%2C%0A%20%20%22destOverride%22%3A%20%5B%0A%20%20%20%20%22http%22%2C%0A%20%20%20%20%22tls%22%0A%20%20%5D%0A%7D`)

	req, err := buildRequest("POST", url, payload, cookie)
	if err != nil {
		return err
	}

	resp, err := makeRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Process the response body if needed.

	return nil
}

func DelInstance(host string, id int, cookie string) error {
	idStr := strconv.Itoa(id)
	url := host + fmt.Sprintf("/xui/inbound/del/%s", idStr)
	method := "POST"
	// payload := bytes.NewBuffer([]byte{})
	req, err := buildRequest(method, url, nil, cookie)
	if err != nil {
		return err
	}

	res, err := makeRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Process the response body if needed.

	return nil
}

func ListInstance(host, cookie string) ([]V2rayInstance, error) {
	url := host + "/xui/inbound/list"
	method := "POST"
	req, err := buildRequest(method, url, nil, cookie)
	if err != nil {
		return nil, err
	}

	res, err := makeRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	fmt.Println(res.Body)
	// Process the response body if needed.
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %s", err)
	}
	// fmt.Println(string(body))
	var instances ListInstanceResponse
	err = json.Unmarshal(body, &instances)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil, fmt.Errorf("Error unmarshalling JSON: %s", err)
	}
	// fmt.Println(instances.ID)
	// fmt.Println(instances.Remark)

	// for _, v := range instances.Obj {
	// 	fmt.Println(v.Remark, v.ID)
	// }
	return instances.Obj, err
}
