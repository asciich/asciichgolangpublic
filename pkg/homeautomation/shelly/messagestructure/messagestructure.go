package messagestructure

import "encoding/json"

// Top-level envelope for all Shelly RPC messages
type ShellyMessage struct {
	Src    string          `json:"src"`
	Dst    string          `json:"dst"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

// NotifyStatus params — sent periodically with sensor readings
type NotifyStatusParams struct {
	Ts          float64          `json:"ts"`
	Temperature *TemperatureComp `json:"temperature:0,omitempty"`
	Humidity    *HumidityComp    `json:"humidity:0,omitempty"`
	DevicePower *DevicePower     `json:"devicepower:0,omitempty"`
}

type TemperatureComp struct {
	ID float64 `json:"id"`
	TC float64 `json:"tC"` // Celsius
	TF float64 `json:"tF"` // Fahrenheit
}

type HumidityComp struct {
	ID float64 `json:"id"`
	RH float64 `json:"rh"` // Relative humidity %
}

type DevicePower struct {
	ID      float64 `json:"id"`
	Battery struct {
		V       float64 `json:"V"`
		Percent float64 `json:"percent"`
	} `json:"battery"`
	External struct {
		Present bool `json:"present"`
	} `json:"external"`
}

// NotifyEvent params — sent on specific events
type NotifyEventParams struct {
	Ts     float64 `json:"ts"`
	Events []struct {
		Component string  `json:"component"`
		Event     string  `json:"event"`
		Ts        float64 `json:"ts"`
	} `json:"events"`
}

// NotifyFullStatus params — sent on wake-up with full device state
type NotifyFullStatusParams struct {
	Ts          float64          `json:"ts"`
	Temperature *TemperatureComp `json:"temperature:0,omitempty"`
	Humidity    *HumidityComp    `json:"humidity:0,omitempty"`
	DevicePower *DevicePower     `json:"devicepower:0,omitempty"`
	BLE         *BLEStatus       `json:"ble,omitempty"`
	Cloud       *CloudStatus     `json:"cloud,omitempty"`
	HTUI        *HTUIStatus      `json:"ht_ui,omitempty"`
	MQTT        *MQTTStatus      `json:"mqtt,omitempty"`
	Sys         *SysStatus       `json:"sys,omitempty"`
	WiFi        *WiFiStatus      `json:"wifi,omitempty"`
	WS          *WSStatus        `json:"ws,omitempty"`
}

type BLEStatus struct{}

type CloudStatus struct {
	Connected bool `json:"connected"`
}

type HTUIStatus struct{}

type MQTTStatus struct {
	Connected bool `json:"connected"`
}

type SysStatus struct {
	MAC              string          `json:"mac"`
	RestartRequired  bool            `json:"restart_required"`
	Time             *string         `json:"time"`
	UnixTime         *float64        `json:"unixtime"`
	Uptime           float64         `json:"uptime"`
	RAMSize          float64         `json:"ram_size"`
	RAMFree          float64         `json:"ram_free"`
	FSSize           float64         `json:"fs_size"`
	FSFree           float64         `json:"fs_free"`
	CfgRev           float64         `json:"cfg_rev"`
	KVSRev           float64         `json:"kvs_rev"`
	WebhookRev       float64         `json:"webhook_rev"`
	AvailableUpdates json.RawMessage `json:"available_updates"`
	WakeupReason     *WakeupReason   `json:"wakeup_reason,omitempty"`
	WakeupPeriod     float64         `json:"wakeup_period"`
	ResetReason      float64         `json:"reset_reason"`
}

type WakeupReason struct {
	Boot  string `json:"boot"`
	Cause string `json:"cause"`
}

type WiFiStatus struct {
	StaIP  string `json:"sta_ip"`
	Status string `json:"status"`
	SSID   string `json:"ssid"`
	RSSI   int    `json:"rssi"`
}

type WSStatus struct {
	Connected bool `json:"connected"`
}
