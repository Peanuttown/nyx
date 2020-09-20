package client

import (
	"fmt"
	gonet "net"
	"os"
	"sync"
	"time"

	"github.com/Peanuttown/nyx/pkg/cfg"
	"github.com/Peanuttown/tzzGoUtil/log"
	"github.com/Peanuttown/tzzGoUtil/net"
	"golang.org/x/net/context"
)

type Client struct {
	cfg    *cfg.ClientCfg
	logger *log.Logger
	onConn chan gonet.Conn
}

func NewClient() *Client {
	return &Client{}
}

func (this *Client) initChannels() {
	this.onConn = make(chan gonet.Conn, 1)
}

func (this *Client) connectSvr() (conn gonet.Conn, err error) {
	conn, err = gonet.Dial("tcp", this.cfg.SvrAddress)
	if err != nil {
		return nil, err
	}
	this.onConn <- conn
	return conn, nil
}

func (this *Client) doServe(overChan chan bool, f func() error) error {
	defer func() {
		overChan <- true
	}()
	return f()

}

func (this *Client) Run(cfg *cfg.ClientCfg) error {
	this.cfg = cfg
	// init logger
	if len(this.cfg.LogFile) == 0 {
		this.logger = log.NewLogger()
	} else {
		file, err := os.OpenFile(this.cfg.LogFile, os.O_RDWR, 0777)
		if err != nil {
			return err
		}
		this.logger.SetOutput(file)
	}
	this.ServeLoop()
	return nil
}

func (this *Client) ServeLoop() {
	for {
		err := this.onceServe()
		if err != nil{
			this.logger.Error("once Serve over: ",err)
		}
		this.logger.Info("准备下一次服务")
		time.Sleep(time.Second*5)
	}
}

func (this *Client) onceServe() error {
	this.logger.Info("start once serve")
	var tun net.ITun
	var conn gonet.Conn
	defer func() {
		if tun != nil {
			tun.Close()
		}
		if conn != nil {
			conn.Close()
		}
	}()
	// connect server
	var err error
	conn, err = gonet.Dial("tcp", this.cfg.SvrAddress)
	if err != nil {
		return err
	}
	// create tun device
	var devName = "nyxTun"
	tun, err = net.NewTun(devName)
	if err != nil {
		return err
	}
	net.SetIp(devName, net.IpWithMask(this.cfg.DevIpWithMask))
	onceTask := NewOnceTask(func() {
		if tun != nil {
			tun.Close()
		}
		if conn != nil {
			conn.Close()
		}
	})

	onceTask.GoDo(
		func(ctx context.Context) error {
			this.logger.Info("start loop for read from tun")
			bytes := make([]byte, 1024)
			for {
				select {
				case <-ctx.Done():
					return nil
				default:
					_, err := tun.Read(bytes)
					if err != nil {
						this.logger.Error(err)
						return err
					}
					// TODO handle data from tun device
				}
			}
		},
	)

	onceTask.GoDo(
		func(ctx context.Context) error {
			var bytes = make([]byte, 1024)
			for {
				select {
				case <-ctx.Done():
				default:
					_, err := conn.Read(bytes)
					if err != nil {
						return err
					}
					// TODO handle data
				}
			}
		},
	)
	onceTask.Wait()
	err = onceTask.Error()
	if err != nil {
		this.logger.Error("once servce over , %v", err)
	} else {
		this.logger.Info("once servce")
	}
	return err
}

type OnceTask struct {
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
	l      sync.Mutex
	Err    error
}

func (this *OnceTask) SaveErr(err error) {
	this.l.Lock()
	defer this.l.Unlock()
	if err != nil {
		if this.Err == nil {
			this.Err = err
		} else {
			this.Err = fmt.Errorf("%v, %w", this.Err, err)
		}
	}
}

func (this *OnceTask) GoDo(f func(ctx context.Context) error) {
	this.wg.Add(1)
	go func() {
		defer func() {
			this.cancel()
			this.wg.Done()
		}()
		err := f(this.ctx)
		this.SaveErr(err)
	}()
}

func NewOnceTask(extraCancel func()) *OnceTask {
	task := &OnceTask{}
	task.wg = sync.WaitGroup{}
	var cancel context.CancelFunc
	task.ctx, cancel = context.WithCancel(context.Background())
	task.cancel = func() {
		cancel()
		if extraCancel != nil {
			extraCancel()
		}
	}
	return task

}

func (this *OnceTask) Wait() {
	this.wg.Wait()
}

func (this *OnceTask) Error() error {
	this.l.Lock()
	defer this.l.Unlock()
	return this.Err
}
