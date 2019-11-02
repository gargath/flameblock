//nolint
package api

type NsolidHook struct {
	Time       string   `json:"time"`
	Event      string   `json:"event"`
	Agents     []Agent  `json:"agents"`
	Config     Config   `json:"config"`
	Assets     []string `json:"assets"`
	BlockedFor int64    `json:"blockedFor"`
	Stack      string   `json:"stack"`
	Threshold  int      `json:"threshold"`
}

type Agent struct {
	Id      string                 `json:"id"`
	Info    AgentInfo              `json:"info"`
	Metrics map[string]interface{} `json:"metrics"`
}

type AgentInfo struct {
	Id           string                 `json:"id"`
	App          string                 `json:"app"`
	AppVersion   string                 `json:"appVersion"`
	Tags         []string               `json:"tags"`
	Pid          int                    `json:"pid"`
	ProcessStart int64                  `json:"processStart"`
	NodeEnv      string                 `json:"nodeEnv"`
	ExecPath     string                 `json:"execPath"`
	Main         string                 `json:"main"`
	Arch         string                 `json:"arch"`
	Platform     string                 `json:"platform"`
	Hostname     string                 `json:"hostname"`
	TotalMem     int64                  `json:"totalMem"`
	Versions     map[string]interface{} `json:"versions"`
	CpuCores     int                    `json:"cpuCores"`
	CpuModel     string                 `json:"cpuModel"`
}

type Config struct {
	Event     string         `json:"event"`
	Actions   []ConfigAction `json:"actions"`
	Threshold int            `json:"threshold"`
}

type ConfigAction struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	Type string `json:"type"`
	Id   string `json:"id"`
}

// Flamegraph data

type Node struct {
	Name     string  `json:"name"`
	Value    int64   `json:"value"`
	Children []*Node `json:"children,omitempty"`
}
