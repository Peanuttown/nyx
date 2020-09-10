package tap

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

type Tap struct {
}

func NewTap(devName string) (TapI, error) {
	f, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	conn, err := f.SyscallConn()
	if err != nil {
		return nil, err
	}

	var ifr [unix.IFNAMSIZ + 64]byte
	copy(ifr[:], []byte(devName))
	*(*uint16)(unsafe.Pointer(&ifr[unix.IFNAMSIZ])) = unix.IFF_TUN
	var errno syscall.Errno
	err = conn.Control(
		func(fd uintptr) {
			_, _, errno = unix.Syscall(
				unix.SYS_IOCTL,
				fd,
				uintptr(unix.TUNSETIFF),
				uintptr(unsafe.Pointer(&ifr[0])),
			)
		},
	)
	if err != nil {
		return nil, err
	}
	if errno != 0 {
		err = fmt.Errorf("创建失败, %w", errno)
		return nil, err
	}
	return nil, nil
}
