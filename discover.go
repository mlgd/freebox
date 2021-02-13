package freebox

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/hashicorp/mdns"
)

const (
	// ServiceName :
	ServiceName = "_fbx-api._tcp"
	// DiscoverProtocolMDNS :
	DiscoverProtocolMDNS DiscoverProtocol = "mdns"
	// DiscoverProtocolHTTP :
	DiscoverProtocolHTTP DiscoverProtocol = "http"
	// DiscoverProtocolHTTPS :
	DiscoverProtocolHTTPS DiscoverProtocol = "https"

	urlMaFreebox = "mafreebox.freebox.fr"
)

// DiscoverProtocol :
type DiscoverProtocol string

type apiVersionResponse struct {
	UID            string `json:"uid"`
	DeviceType     string `json:"device_type"`
	DeviceName     string `json:"device_name"`
	BoxModel       string `json:"box_model"`
	BoxModelName   string `json:"box_model_name"`
	APIVersion     string `json:"api_version"`
	APIBaseURL     string `json:"api_base_url"`
	APIDomain      string `json:"api_domaine"`
	HTTPSAvailable bool   `json:"https_available"`
	HTTPSPort      int    `json:"https_port"`
}

// Discover :
func Discover(protocol DiscoverProtocol) (devices []Device, err error) {
	switch protocol {
	case DiscoverProtocolMDNS:
		return discoverMDNS()
	case DiscoverProtocolHTTP:
		return discoverHTTP(DiscoverProtocolHTTP)
	case DiscoverProtocolHTTPS:
		return discoverHTTP(DiscoverProtocolHTTPS)
	}
	return nil, nil
}

func discoverMDNS() (devices []Device, err error) {
	// Make a channel for results and start listening
	entriesCh := make(chan *mdns.ServiceEntry, 10)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			entry, ok := <-entriesCh
			if !ok {
				break
			}
			device := Device{
				Name:     entry.Name,
				Host:     entry.Host,
				IP:       entry.AddrV4.String(),
				IPv6:     entry.AddrV6.String(),
				PortHTTP: entry.Port,
			}
			for _, info := range entry.InfoFields {
				infoKV := strings.Split(info, "=")
				switch infoKV[0] {
				case "api_domain":
					device.APIDomain = infoKV[1]
				case "api_version":
					device.APIVersion = infoKV[1]
				case "api_base_url":
					device.APIBaseURL = infoKV[1]
				case "box_model":
					device.BoxModel = infoKV[1]
				case "box_model_name":
					device.BoxModelName = infoKV[1]
				case "https_port":
					device.PortHTTPS, _ = strconv.Atoi(infoKV[1])
				case "https_available":
					device.HTTPSAvailable, _ = strconv.ParseBool(infoKV[1])
				case "uid":
					device.UID = infoKV[1]
				}
			}
			devices = append(devices, device)
		}
		wg.Done()
	}()

	// Start the lookup
	if err := mdns.Lookup(ServiceName, entriesCh); err != nil {
		close(entriesCh)
		wg.Wait()
		return nil, err
	}
	close(entriesCh)
	wg.Wait()

	return devices, nil
}

func discoverHTTP(protocol DiscoverProtocol) (devices []Device, err error) {
	if protocol == "" {
		protocol = "http"
	}
	api := fmt.Sprintf("%s://%s/api_version", protocol, urlMaFreebox)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, api, nil)
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var apiVersion *apiVersionResponse
	if err := json.Unmarshal(body, &apiVersion); err != nil {
		return nil, err
	}

	addr, err := net.LookupIP(urlMaFreebox)
	if err != nil {
		return nil, err
	}
	var addrV4 net.IP
	if len(addr) > 0 {
		addrV4 = addr[0]
	}

	device := Device{
		Name:           apiVersion.DeviceName,
		Host:           urlMaFreebox,
		IP:             addrV4.String(),
		PortHTTP:       80,
		APIDomain:      apiVersion.APIDomain,
		APIVersion:     apiVersion.APIVersion,
		APIBaseURL:     apiVersion.APIBaseURL,
		BoxModel:       apiVersion.BoxModel,
		BoxModelName:   apiVersion.BoxModelName,
		PortHTTPS:      apiVersion.HTTPSPort,
		HTTPSAvailable: apiVersion.HTTPSAvailable,
		UID:            apiVersion.UID,
	}
	devices = append(devices, device)

	return devices, nil
}
