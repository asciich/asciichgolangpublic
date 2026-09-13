package kvmutilsgeneric

import "encoding/xml"

// Parsed representation of the relevant parts of 'virsh net-dumpxml <name>'.
type KvmNetworkXml struct {
	XMLName xml.Name `xml:"network"`
	Name    string   `xml:"name"`
	Forward struct {
		Mode string `xml:"mode,attr"`
	} `xml:"forward"`
	Ip struct {
		Address string `xml:"address,attr"`
		Netmask string `xml:"netmask,attr"`
		Dhcp    struct {
			Range struct {
				Start string `xml:"start,attr"`
				End   string `xml:"end,attr"`
			} `xml:"range"`
		} `xml:"dhcp"`
	} `xml:"ip"`
}

func (x *KvmNetworkXml) GetForwardMode() (string, error) {
	return x.Forward.Mode, nil
}

func (x *KvmNetworkXml) GetIpAddress() (string, error) {
	ip := x.Ip.Address
	return ip, nil
}

func (x *KvmNetworkXml) GetIpNetmask() (string, error) {
	netmask := x.Ip.Netmask
	return netmask, nil
}

func (x *KvmNetworkXml) GetIpDhcpRangeStart() (string, error) {
	start := x.Ip.Dhcp.Range.Start
	return start, nil
}

func (x *KvmNetworkXml) GetIpDhcpRangeEnd() (string, error) {
	end := x.Ip.Dhcp.Range.End
	return end, nil
}
