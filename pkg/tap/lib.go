package tap

import (
	"os"
	"unsafe"

	"github.com/Peanuttown/tzzGoUtil/sys/io"
	"golang.org/x/sys/unix"
)

type Tap struct {
	file *os.File
}

func NewTap(devName string) (TapI, error) {
	f, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	var ifr [unix.IFNAMSIZ + 64]byte
	copy(ifr[:], []byte(devName))
	*(*uint16)(unsafe.Pointer(&ifr[unix.IFNAMSIZ])) = unix.IFF_TAP | unix.IFF_NO_PI
	err = io.IOCtl(
		f.Fd(),
		unix.TUNSETIFF,
		uintptr(unsafe.Pointer(&ifr[0])),
	)
	if err != nil {
		return nil, err
	}
	return &Tap{
		file: os.NewFile(f.Fd(), devName),
	}, nil
}

func (this *Tap) Read(b []byte) (n int, err error) {
	return this.file.Read(b)
}
