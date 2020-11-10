package main

import (
	"bufio"
	"fmt"
	"github.com/astaxie/beego/logs"
	flag "github.com/spf13/pflag"
	"log"
	"net"
	"net/smtp"
	"os"
	"strings"
	"time"
)

var (
	IP      string
	Timeout int
	Help    bool
)

func init()  {
	flag.BoolVar(&Help, "help", false, "usage help")
	flag.IntVar(&Timeout, "timeout", 60, "timeout")
	flag.StringVar(&IP, "ip", "127.0.0.1:25", "smtp server address")
}

type DebugConn struct {
	conn net.Conn
}

func (d *DebugConn)Read(b []byte) (n int, err error) {
	cnt, err := d.conn.Read(b)
	if cnt > 0 {
		logs.Info("DEBUG Read: %s", b[:cnt])
	}
	return cnt, err
}

func (d *DebugConn)Write(b []byte) (n int, err error) {
	logs.Info("DEBUG Write: %s", b[:])
	return d.conn.Write(b)
}

func (d *DebugConn)Close() error {
	return d.conn.Close()
}

func (d *DebugConn)LocalAddr() net.Addr {
	return d.conn.LocalAddr()
}

func (d *DebugConn)RemoteAddr() net.Addr {
	return d.conn.RemoteAddr()
}

func (d *DebugConn)SetDeadline(t time.Time) error {
	return d.conn.SetDeadline(t)
}

func (d *DebugConn)SetReadDeadline(t time.Time) error {
	return d.conn.SetReadDeadline(t)
}

func (d *DebugConn)SetWriteDeadline(t time.Time) error {
	return d.conn.SetWriteDeadline(t)
}

func Debug(conn net.Conn) net.Conn {
	return &DebugConn{conn: conn}
}

func SmtpDial(addr string, timeout time.Duration) (*smtp.Client, error, net.Conn) {
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return nil, err, nil
	}
	host, _, _ := net.SplitHostPort(addr)
	cli, err  := smtp.NewClient(Debug(conn), host)
	return cli, err, conn
}

func main()  {
	flag.Parse()
	if Help {
		flag.Usage()
		return
	}

	_, err, conn := SmtpDial(IP, time.Duration(Timeout) * time.Second)
	if err != nil {
		log.Println(err.Error())
		return
	}

	go func() {
		var buff [1024]byte
		for  {
			cnt, err := conn.Read(buff[:])
			if err != nil {
				return
			}
			fmt.Println(string(buff[:cnt]))
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text: ")

	for  {
		text, _ := reader.ReadString('\n')
		text = strings.ReplaceAll(text, "\r", "")
		text = strings.ReplaceAll(text, "\n", "")
		//log.Printf("%s\r\n", text)

		_, err = conn.Write([]byte(fmt.Sprintf("%s\r\n", text)))
		if err != nil {
			log.Println(err.Error())
		}
	}
}