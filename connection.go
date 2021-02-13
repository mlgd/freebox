package freebox

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	endPointConnectionStatus = "connection"
	endPointConnectionLogs   = "connection/logs"
)

// ConnectionStatusResponse :
type ConnectionStatusResponse struct {
	Success   bool                            `json:"success"`
	ErrorCode string                          `json:"error_code"`
	Message   string                          `json:"msg"`
	Result    *ConnectionStatusResponseResult `json:"result"`
}

// ConnectionStatusResponseResult :
type ConnectionStatusResponseResult struct {
	Type          string `json:"type"`
	State         string `json:"state"`
	Media         string `json:"media"`
	IPv4          string `json:"ipv4"`
	IPv4PortRange []int  `json:"ipv4_port_range"`
	RateDown      int64  `json:"rate_down"`
	RateUp        int64  `json:"rate_up"`
	BytesUp       int64  `json:"bytes_up"`
	BytesDown     int64  `json:"bytes_down"`
	BandwidthUp   int64  `json:"bandwidth_up"`
	BandwidthDown int64  `json:"bandwidth_down"`
}

// ConnectionLogsResponse :
type ConnectionLogsResponse struct {
	Success   bool            `json:"success"`
	ErrorCode string          `json:"error_code"`
	Message   string          `json:"msg"`
	Result    []ConnectionLog `json:"result"`
}

// ConnectionLog :
type ConnectionLog struct {
	ID           int    `json:"id"`
	Date         int    `json:"date"`
	State        string `json:"state"`
	Type         string `json:"type"`
	Conn         string `json:"conn"`
	Link         string `json:"link"`
	BandwithDown int64  `json:"bw_down"`
	BandwithUp   int64  `json:"bw_up"`
}

// ConnectionStatus :
func (device Device) ConnectionStatus(sessionToken string) (response *ConnectionStatusResponse, err error) {
	api := fmt.Sprintf("%s%s/", device.url(), endPointConnectionStatus)

	client := &http.Client{}
	req, err := http.NewRequest("GET", api, nil)
	req.Header.Add("X-Fbx-App-Auth", sessionToken)
	httpResponse, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	body, err := ioutil.ReadAll(httpResponse.Body)
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response, nil
}

// ConnectionLogs :
func (device Device) ConnectionLogs(sessionToken string) (response *ConnectionLogsResponse, err error) {
	api := fmt.Sprintf("%s%s/", device.url(), endPointConnectionLogs)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, api, nil)
	req.Header.Add("X-Fbx-App-Auth", sessionToken)
	httpResponse, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	body, err := ioutil.ReadAll(httpResponse.Body)
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response, nil
}
