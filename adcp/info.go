package adcp

import (
	"context"
	"encoding/json"
	"strings"
)

// HardwareInfo contains the common information for device hardware information
type HardwareInfo struct {
	ModelName     string           `json:"model_name,omitempty"`
	SerialNumber  string           `json:"serial_number,omitempty"`
	NetworkInfo   NetworkInfo      `json:"network_information,omitempty"`
	FilterStatus  string           `json:"filter_status,omitempty"`
	WarningStatus []string         `json:"warning_status,omitempty"`
	ErrorStatus   []string         `json:"error_status,omitempty"`
	PowerStatus   string           `json:"power_status,omitempty"`
	TimerInfo     []map[string]int `json:"timer_info,omitempty"`
}

// NetworkInfo contains the network information for the device
type NetworkInfo struct {
	IPAddress  string   `json:"ip_address,omitempty"`
	MACAddress string   `json:"mac_address,omitempty"`
	Gateway    string   `json:"gateway,omitempty"`
	DNS        []string `json:"dns,omitempty"`
}

var (
	modelName   = []byte("modelname ?\r\n")
	serialNum   = []byte("serialnum ?\r\n")
	ipAddr      = []byte("ipv4_ip_address ?\r\n")
	gateway     = []byte("ipv4_default_gateway ?\r\n")
	dns         = []byte("ipv4_dns_server1 ?\r\n")
	dns2        = []byte("ipv4_dns_server2 ?\r\n")
	macAddr     = []byte("mac_address ?\r\n")
	filter      = []byte("filter_status ?\r\n")
	warnings    = []byte("warning ?\r\n")
	errors      = []byte("error ?\r\n")
	powerStatus = []byte("power_status ?\r\n")
	timer       = []byte("timer ?\r\n")
)

// Info returns the hardware information of the projector
func (p *Projector) Info(ctx context.Context) (interface{}, error) {
	var info HardwareInfo

	// model name
	resp, err := p.SendCommand(ctx, p.Address, modelName)
	if err != nil {
		return info, err
	}

	info.ModelName = strings.Trim(resp, "\"")

	// ip address
	resp, err = p.SendCommand(ctx, p.Address, ipAddr)
	if err != nil {
		return info, err
	}

	info.NetworkInfo.IPAddress = strings.Trim(resp, "\"")

	// gateway
	resp, err = p.SendCommand(ctx, p.Address, gateway)
	if err != nil {
		return info, err
	}

	info.NetworkInfo.Gateway = strings.Trim(resp, "\"")

	// dns
	resp, err = p.SendCommand(ctx, p.Address, dns)
	if err != nil {
		return info, err
	}

	info.NetworkInfo.DNS = append(info.NetworkInfo.DNS, strings.Trim(resp, "\""))

	resp, err = p.SendCommand(ctx, p.Address, dns2)
	if err != nil {
		return info, err
	}

	info.NetworkInfo.DNS = append(info.NetworkInfo.DNS, strings.Trim(resp, "\""))

	// mac address
	resp, err = p.SendCommand(ctx, p.Address, macAddr)
	if err != nil {
		return info, err
	}

	info.NetworkInfo.MACAddress = strings.Trim(resp, "\"")

	// serial number
	resp, err = p.SendCommand(ctx, p.Address, serialNum)
	if err != nil {
		return info, err
	}

	info.SerialNumber = strings.Trim(resp, "\"")

	// filter status
	resp, err = p.SendCommand(ctx, p.Address, filter)
	if err != nil {
		return info, err
	}

	info.FilterStatus = strings.Trim(resp, "\"")

	// power status
	resp, err = p.SendCommand(ctx, p.Address, powerStatus)
	if err != nil {
		return info, err
	}

	info.PowerStatus = strings.Trim(resp, "\"")

	// warnings
	resp, err = p.SendCommand(ctx, p.Address, warnings)
	if err != nil {
		return info, err
	}

	err = json.Unmarshal([]byte(resp), &info.WarningStatus)
	if err != nil {
		return info, err
	}

	// errors
	resp, err = p.SendCommand(ctx, p.Address, errors)
	if err != nil {
		return info, err
	}

	err = json.Unmarshal([]byte(resp), &info.ErrorStatus)
	if err != nil {
		return info, err
	}

	// timer info
	resp, err = p.SendCommand(ctx, p.Address, timer)
	if err != nil {
		return info, err
	}

	err = json.Unmarshal([]byte(resp), &info.TimerInfo)
	if err != nil {
		return info, err
	}

	return info, nil
}
