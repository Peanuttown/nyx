package tap

import (
    "strings"
	"fmt"
	"os/exec"
    "github.com/Peanuttown/tzzGoUtil/process"
)

const(
    ETH_PREAMBLE_LEN=7
    ETH_SFD_LEN=1
    ETH_MAC_LEN=6
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

func DeviceUp(devName string)(error){
    _,errOut,err := process.ExeOutput("ip","link","set",devName,"up")
    if err != nil{
        return fmt.Errorf("%v, %w",errOut,err)
    }
    return nil
}

type EthernetFrame struct{
    preamble [7]byte
    sfd byte
    dstMAC string
    srcMAC string
}

func (this *EthernetFrame) String()(string){
    // preamble
    s := strings.Builder{}
    s.WriteString("preamble:\n")
    for _,b := range this.preamble{
        s.WriteString(fmt.Sprintf("%#b    ",b))
    }
    s.WriteString(fmt.Sprintf("sfd: %#b\n",this.sfd))
    s.WriteString(fmt.Sprintf("dstMAC : %s\n",this.dstMAC))
    s.WriteString(fmt.Sprintf("srcMAC: %s\n",this.srcMAC))
    return s.String()
}

func macByteToString(macBytes [ETH_MAC_LEN]byte)string{
    var list = make([]string,0,len(macBytes))
    for _,v := range macBytes{
        list = append(list,fmt.Sprintf("%x",v))
    }
    return strings.Join(list,":")
}

func DecodeEthernetFrame(b []byte)(*EthernetFrame,error){
    var index int
    var preamble [ETH_PREAMBLE_LEN]byte
    copy(preamble[:],b[index:])//0
    index+=ETH_PREAMBLE_LEN
    sfd:=b[index]//7
    index+=1
    var dstMAC [ETH_MAC_LEN]byte
    copy(dstMAC[:],b[index:]) //8
    index+=ETH_MAC_LEN
    var srcMAC [ETH_MAC_LEN]byte
    copy(srcMAC[:],b[index:])
    index+=ETH_MAC_LEN
    // parse mac
    return &EthernetFrame{
        preamble:preamble,
        sfd:sfd,
        dstMAC :macByteToString(dstMAC),
        srcMAC:macByteToString(srcMAC),
    },nil
}
