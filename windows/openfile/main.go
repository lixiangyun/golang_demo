package main

import (
	"github.com/astaxie/beego/logs"
	"os/exec"
)

func OpenBrowserWeb(url string)  {
	cmd := exec.Command("rundll32","url.dll,FileProtocolHandler", url)
	err := cmd.Run()
	if err != nil {
		logs.Error("run cmd fail, %s", err.Error())
	}
}

func main()  {
	OpenBrowserWeb("https://www.baidu.com/")
}
