package bravia

//SonyAudioResponse is the parent struct returned when we query audio state
type SonyAudioResponse struct {
	Result [][]SonyAudioSettings `json:"result"`
	ID     int                   `json:"id"`
}

//SonyAudioSettings is the child struct returned
type SonyAudioSettings struct {
	Target    string `json:"target"`
	Volume    int    `json:"volume"`
	Mute      bool   `json:"mute"`
	MaxVolume int    `json:"maxVolume"`
	MinVolume int    `json:"minVolume"`
}

type SonyAVContentSettings struct {
	URI        string `json:"uri"`
	Source     string `json:"source"`
	Title      string `json:"title"`
	Status     string `json:"status"`
	Connection bool   `json:"connection"`
}

type SonyAVContentResponse struct {
	Result []SonyAVContentSettings `json:"result"`
	ID     int                     `json:"id"`
}

type SonyMultiAVContentResponse struct {
	Result [][]SonyAVContentSettings `json:"result"`
	ID     int                       `json:"id"`
}

//SonyTVRequest represents the struct we need to send.
type SonyTVRequest struct {
	Method  string                   `json:"method"`
	Version string                   `json:"version"`
	ID      int                      `json:"id"`
	Params  []map[string]interface{} `json:"params"`
}

type SonyTVSystemResponse struct {
	ID     int `json:"id"`
	Result []SonySystemInformation
}

type SonySystemInformation struct {
	Product    string `json:"product"`
	Region     string `json:"region,omitempty"`
	Language   string `json:"language,omitempty"`
	Model      string `json:"model"`
	Serial     string `json:"serial,omitempty"`
	MAC        string `json:"macAddr,omitempty"`
	Name       string `json:"name"`
	Generation string `json:"generation,omitempty"`
	Area       string `json:"area,omitempty"`
	CID        string `json:"cid,omitempty"`
}

type SonyNetworkResponse struct {
	ID     int `json:"id"`
	Result [][]SonyTVNetworkInformation
}

type SonyTVNetworkInformation struct {
	NetworkInterface string   `json:"netif"`
	HardwareAddress  string   `json:"hwAddr"`
	IPv4             string   `json:"ipAddrV4"`
	IPv6             string   `json:"ipAddrV6"`
	Netmask          string   `json:"netmask"`
	Gateway          string   `json:"gateway"`
	DNS              []string `json:"dns"`
}
