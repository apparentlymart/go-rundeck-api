package rundeck

import (
	"encoding/xml"
	"time"
)

type SystemInfo struct {
	XMLName       xml.Name        `xml:"system"`
	ServerTime    SystemTimestamp `xml:"timestamp"`
	Rundeck       About           `xml:"rundeck"`
	OS            SystemOS        `xml:"os"`
	JVM           SystemJVM       `xml:"jvm"`
	Stats         SystemStats     `xml:"stats"`
}

type About struct {
	XMLName    xml.Name `xml:"rundeck"`
	Version    string   `xml:"version"`
	ApiVersion int64    `xml:"apiversion"`
	Build      string   `xml:"build"`
	Node       string   `xml:"node"`
	Base       string   `xml:"base"`
	ServerUUID string   `xml:"serverUUID,omitempty"`
}

type SystemTimestamp struct {
	Epoch       string `xml:"epoch,attr"`
	EpochUnit   string `xml:"unit,attr"`
	DateTimeStr string `xml:"datetime"`
}

type SystemOS struct {
	Architecture string `xml:"arch"`
	Name         string `xml:"name"`
	Version      string `xml:"version"`
}

type SystemJVM struct {
	Name                  string `xml:"name"`
	Vendor                string `xml:"vendor"`
	Version               string `xml:"version"`
	ImplementationVersion string `xml:"implementationVersion"`
}

type SystemStats struct {
	XMLName   xml.Name             `xml:"stats"`
	Uptime    SystemUptime         `xml:"uptime"`
	CPU       SystemCPUStats       `xml:"cpu"`
	Memory    SystemMemoryUsage    `xml:"memory"`
	Scheduler SystemSchedulerStats `xml:"scheduler"`
	Threads   SystemThreadStats    `xml:"threads"`
}

type SystemUptime struct {
	XMLName       xml.Name        `xml:"uptime"`
	Duration      string          `xml:"duration,attr"`
	DurationUnit  string          `xml:"unit,attr"`
	BootTimestamp SystemTimestamp `xml:"since"`
}

type SystemCPUStats struct {
	XMLName     xml.Name `xml:"cpu"`
	LoadAverage struct {
		Unit  string  `xml:"unit,attr"`
		Value float64 `xml:",chardata"`
	} `xml:"loadAverage"`
	ProcessorCount int64 `xml:"processors"`
}

type SystemMemoryUsage struct {
	XMLName xml.Name `xml:"memory"`
	Unit    string   `xml:"unit,attr"`
	Max     int64    `xml:"max"`
	Free    int64    `xml:"free"`
	Total   int64    `xml:"total"`
}

type SystemSchedulerStats struct {
	RunningJobCount int64 `xml:"running"`
}

type SystemThreadStats struct {
	ActiveThreadCount int64 `xml:"active"`
}

func (c *Client) GetSystemInfo() (*SystemInfo, error) {
	sysInfo := &SystemInfo{}
	err := c.get([]string{"system", "info"}, nil, sysInfo)
	return sysInfo, err
}

func (ts *SystemTimestamp) DateTime() time.Time {
	// Assume the server will always give us a valid timestamp,
	// so we don't need to handle the error case.
	// (Famous last words?)
	t, _ := time.Parse(time.RFC3339, ts.DateTimeStr)
	return t
}
