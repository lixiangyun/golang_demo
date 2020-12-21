package main

import (
	"fmt"
	ole "github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

/*
	Windows Scripting Host 为我们提供了一个 WshShortcut 对象，可以使用此对象来创建快捷方式或获取快捷方式的信息
	WshShortcut 有以下属性：
	TargetPath:获取或设置快捷方式指向的目标文件的路径
	FullName:获取或设置快捷方式的路径
	Description:获取或设置快捷方式的说明
	IconLocation:获取或设置快捷方式图标的位置
	WindowStyle:获取或设置启动目标程序时所使用的窗口样式
	HotKey:获取或设置用于启动目标程序的热键，只热键只能激活 WINDOWS 桌面和开始菜单中的快捷方式
	WorkingDirectory:获取或设置目标程序的工作目录
	Arguments:获取一个 WshArgument 对象的集合
*/

func CreateShortcut(dst, src, icon string, pwd string) error {
	err := ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_SPEED_OVER_MEMORY)
	if err != nil {
		return err
	}

	oleShellObject, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		return err
	}
	defer oleShellObject.Release()

	wshell, err := oleShellObject.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return err
	}
	defer wshell.Release()

	cs, err := oleutil.CallMethod(wshell, "CreateShortcut", dst)
	if err != nil {
		return err
	}

	idispatch := cs.ToIDispatch()

	defer idispatch.Release()

	oleutil.PutProperty(idispatch, "TargetPath", src)
	oleutil.PutProperty(idispatch, "IconLocation", icon)

	oleutil.PutProperty(idispatch, "WindowStyle", 7)
	oleutil.PutProperty(idispatch, "Minimized", 7)
	oleutil.PutProperty(idispatch, "Maximized", 0)
	oleutil.PutProperty(idispatch, "Normal", 4)

	oleutil.PutProperty(idispatch, "WorkingDirectory", pwd)
	oleutil.PutProperty(idispatch, "Arguments", "-c config")

	oleutil.CallMethod(idispatch, "Save")
	return nil
}

func DeskTopPath() (string, error) {
	ur, err := user.Current()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s\\Desktop", ur.HomeDir), nil
}

func AppName() string {
	path := os.Args[0]
	idx := strings.LastIndex(path, "\\")
	if idx != -1 {
		return path[idx+1:]
	}
	return path
}

func StartUpPath() (string, error) {
	dir := os.Getenv("APPDATA")
	path := fmt.Sprintf("%s\\Microsoft\\Windows\\Start Menu\\Programs\\Startup", dir)
	file, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if file.IsDir() {
		return path, nil
	}
	return "", fmt.Errorf("path is not dir")
}

func CurPath() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}
	return dir, nil
}


func main()  {
	args := os.Args[1:]
	if len(args) > 0 {
		fmt.Printf("cmd: %v\n", args)
	}

	//path, err := DeskTopPath()
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}

	srcDir, err := CurPath()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	startup, err := StartUpPath()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	src := fmt.Sprintf("%s\\%s", srcDir, AppName())
	icon := fmt.Sprintf("%s\\icons.ico", srcDir)
	//dest := fmt.Sprintf("%s\\Demo.lnk", path)
	dest := fmt.Sprintf("%s\\Demo.lnk", startup)

	fmt.Println(dest)
	fmt.Println(src)
	fmt.Println(startup)

	err = CreateShortcut(dest, src, icon, srcDir)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for  {
		time.Sleep(time.Hour)
	}
}
