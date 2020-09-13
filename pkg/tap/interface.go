package tap

import (
	"fmt"
	"os/exec"
)

type TapI interface {
	Read([]byte) (int, error)
}

func SetIP(dev string, ip string) error {
	cmd := exec.Command("ip", "addr", "add", ip, "dev", dev)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s, %w", output, err)
	}
	return nil
}
