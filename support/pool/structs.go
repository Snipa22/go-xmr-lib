package pool

type ConfigItem struct {
	Module    string      `json:"module"`
	Item      string      `json:"item"`
	ItemValue interface{} `json:"item_value"`
}

type ConfigPush []ConfigItem

type PortConfig struct {
	PoolPort   int    `json:"poolPort"`
	Difficulty int    `json:"difficulty"`
	PortDesc   string `json:"portDesc"`
	PortType   string `json:"portType"`
	Hidden     bool   `json:"hidden"`
	Ssl        bool   `json:"ssl"`
}

type PortPush []PortConfig

type BannedConfig struct {
	Address string `json:"address"`
	Proxy   bool   `json:"proxy"`
	Banned  bool   `json:"banned"`
}

type BannedPush []BannedConfig

type PoolInit struct {
	Ports []struct {
		Port       int    `json:"port"`
		Difficulty int    `json:"difficulty"`
		Desc       string `json:"desc"`
		PortType   string `json:"portType"`
		Hidden     bool   `json:"hidden"`
		Ssl        bool   `json:"ssl"`
	} `json:"ports"`
	PoolId   int    `json:"pool_id"`
	PoolIp   string `json:"pool_ip"`
	Hostname string `json:"hostname"`
}

type PoolCheckin struct {
	CurrentTime int64          `json:"current_time"`
	PoolId      int            `json:"pool_id"`
	BlockId     int            `json:"block_id"`
	Ports       map[string]int `json:"ports"`
}
