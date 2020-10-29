package bravia

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"go.uber.org/zap"
)

type NetworkInfo struct {
	IPAddress  string
	MACAddress string
	Gateway    string
	DNS        []string
}

type Info struct {
	Hostname        string
	ModelName       string
	SerialNumber    string
	FirmwareVersion string
	NetworkInfo     NetworkInfo
	PowerStatus     bool
}

// Info returns the hardware information for the device
func (t *TV) Info(ctx context.Context) (interface{}, error) {
	var toReturn Info

	// get the hostname
	addr, e := net.LookupAddr(t.Address)
	if e != nil {
		toReturn.Hostname = t.Address
	} else {
		toReturn.Hostname = strings.Trim(addr[0], ".")
	}

	// get Sony TV system information
	systemInfo, err := t.getSystemInfo(ctx)
	if err != nil {
		return toReturn, fmt.Errorf("could not get system info from %s: %s", t.Address, err)
	}

	toReturn.ModelName = systemInfo.Model
	toReturn.SerialNumber = systemInfo.Serial
	toReturn.FirmwareVersion = systemInfo.Generation

	// get Sony TV network settings
	networkInfo, err := t.getNetworkInfo(ctx)
	if err != nil {
		return toReturn, fmt.Errorf("could not get network info from %s: %s", t.Address, err)
	}

	toReturn.NetworkInfo = NetworkInfo{
		IPAddress:  networkInfo.IPv4,
		MACAddress: networkInfo.HardwareAddress,
		Gateway:    networkInfo.Gateway,
		DNS:        networkInfo.DNS,
	}

	t.Log.Info("network info", zap.Any("networkInfo", toReturn))

	// get power status
	powerStatus, err := t.Power(context.TODO())
	if err != nil {
		return toReturn, fmt.Errorf("could not get power status from %s: %s", t.Address, err)
	}

	toReturn.PowerStatus = powerStatus

	return toReturn, nil
}

func (t *TV) getSystemInfo(ctx context.Context) (SonySystemInformation, error) {
	var system SonyTVSystemResponse

	payload := SonyTVRequest{
		Params: []map[string]interface{}{},
		Method: "getSystemInformation", Version: "1.0",
		ID: 1,
	}

	response, err := t.PostHTTPWithContext(ctx, "system", payload)
	if err != nil {
		return SonySystemInformation{}, err
	}

	err = json.Unmarshal(response, &system)
	if err != nil {
		return SonySystemInformation{}, err
	}

	return system.Result[0], nil
}

func (t *TV) getNetworkInfo(ctx context.Context) (SonyTVNetworkInformation, error) {
	var network SonyNetworkResponse

	payload := SonyTVRequest{
		ID:      2,
		Method:  "getNetworkSettings",
		Version: "1.0",
		Params: []map[string]interface{}{
			map[string]interface{}{
				"netif": "eth0",
			},
		},
	}

	response, err := t.PostHTTPWithContext(ctx, "system", payload)
	if err != nil {
		return SonyTVNetworkInformation{}, err
	}

	err = json.Unmarshal(response, &network)
	if err != nil {
		return SonyTVNetworkInformation{}, err
	}

	return network.Result[0][0], nil
}

func (t *TV) Healthy(ctx context.Context) error {
	_, err := t.Power(ctx)
	if err != nil {
		return fmt.Errorf("failed health check: %s", err)
	}

	return nil
}
