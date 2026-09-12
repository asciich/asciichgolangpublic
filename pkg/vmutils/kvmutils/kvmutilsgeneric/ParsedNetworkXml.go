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
