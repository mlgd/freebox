package freebox

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

// Device :
type Device struct {
	Name           string
	BoxModel       string
	BoxModelName   string
	Host           string
	IP             string
	IPv6           string
	PortHTTP       int
	PortHTTPS      int
	HTTPSAvailable bool
	APIVersion     string
	APIDomain      string
	APIBaseURL     string
	UID            string
}

func (d Device) majorAPIVersion() (int, error) {
	re := regexp.MustCompile(`([\d])+\.([\d])+`)
	res := re.FindAllStringSubmatch(d.APIVersion, 1)
	if len(res) == 1 {
		major, _ := strconv.Atoi(res[0][1])
		return major, nil
	}
	return 0, errors.New("Major API version not found")
}

// URL :
func (d Device) url() string {
	version := ""
	if major, err := d.majorAPIVersion(); err == nil {
		version = fmt.Sprintf("v%d/", major)
	}
	if d.HTTPSAvailable {
		return fmt.Sprintf("https://%s:%d%s%s", d.APIDomain, d.PortHTTPS, d.APIBaseURL, version)
	}
	return fmt.Sprintf("http://%s:%d%s%s", d.APIDomain, d.PortHTTP, d.APIBaseURL, version)
}
