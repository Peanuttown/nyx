package cfg

import "fmt"

type ClientCfg struct {
	SvrAddress    string `yaml:"svr_addr"` //ip:port
	DevIpWithMask string `yaml:"dev_ip_with_mask"`
	LogFile       string `yaml:"log_file"`
}

// 赋予默认值
func EndowDefault(cfg *ClientCfg) {
	if len(cfg.DevIpWithMask) == 0 {
		cfg.DevIpWithMask = "172.31.1.1/16"
	}
}

func (this *ClientCfg) Check() error {
	if len(this.DevIpWithMask) == 0 {
		return fmt.Errorf("tunnel 设备的 ip 地址为空")
	}

	if len(this.SvrAddress) == 0 {
		return fmt.Errorf("服务器地址为空")
	}
	return nil
}
