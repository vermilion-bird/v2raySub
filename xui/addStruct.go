package xui

import (
	"fmt"
	"net/url"
)

type VmessConfig struct {
	Up             int    `json:"up"`
	Down           int    `json:"down"`
	Total          int    `json:"total"`
	Remark         string `json:"remark"`
	Enable         bool   `json:"enable"`
	ExpiryTime     int    `json:"expiryTime"`
	Listen         string `json:"listen"`
	Port           int    `json:"port"`
	Protocol       string `json:"protocol"`
	Settings       Settings
	StreamSettings StreamSettings `json:"streamSettings"`
	Sniffing       Sniffing       `json:"sniffing"`
}

type Settings struct {
	Clients                   []Client `json:"clients"`
	DisableInsecureEncryption bool     `json:"disableInsecureEncryption"`
}

type Client struct {
	ID      string `json:"id"`
	AlterID int    `json:"alterId"`
}

type StreamSettings struct {
	Network    string `json:"network"`
	Security   string `json:"security"`
	WSSettings struct {
		Path    string   `json:"path"`
		Headers struct{} `json:"headers"`
	} `json:"wsSettings"`
}

type Sniffing struct {
	Enabled      bool     `json:"enabled"`
	DestOverride []string `json:"destOverride"`
}

func (c VmessConfig) ToString() string {
	query := url.Values{}
	query.Set("up", fmt.Sprintf("%d", c.Up))
	query.Set("down", fmt.Sprintf("%d", c.Down))
	query.Set("total", fmt.Sprintf("%d", c.Total))
	query.Set("remark", c.Remark)
	query.Set("enable", fmt.Sprintf("%t", c.Enable))
	query.Set("expiryTime", fmt.Sprintf("%d", c.ExpiryTime))
	query.Set("listen", c.Listen)
	query.Set("port", fmt.Sprintf("%d", c.Port))
	query.Set("protocol", c.Protocol)

	settingsStr := fmt.Sprintf(`{"clients": [{"id": "%s", "alterId": %d}], "disableInsecureEncryption": %t}`,
		c.Settings.Clients[0].ID, c.Settings.Clients[0].AlterID, c.Settings.DisableInsecureEncryption)
	query.Set("settings", url.QueryEscape(settingsStr))

	streamSettingsStr := fmt.Sprintf(`{"network": "%s", "security": "%s", "wsSettings": {"path": "%s", "headers": {}}}`,
		c.StreamSettings.Network, c.StreamSettings.Security, c.StreamSettings.WSSettings.Path)
	query.Set("streamSettings", url.QueryEscape(streamSettingsStr))

	sniffingStr := fmt.Sprintf(`{"enabled": %t, "destOverride": %s}`,
		c.Sniffing.Enabled, arrayToString(c.Sniffing.DestOverride))
	query.Set("sniffing", url.QueryEscape(sniffingStr))

	return query.Encode()
}

func arrayToString(arr []string) string {
	return fmt.Sprintf(`["%s"]`, arr[0]) // Assuming only one element in the array for simplicity
}
