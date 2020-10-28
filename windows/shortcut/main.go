package main

import (
	"fmt"
	ole "github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"os"
	"os/user"
	"path/filepath"
)

func CreateShortcut(dst, src, icon string) error {
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

	oleutil.PutProperty(idispatch, "WorkingDirectory", 0)
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

	path, err := DeskTopPath()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	srcDir, err := CurPath()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	src := fmt.Sprintf("%s\\%s", srcDir, os.Args[0])
	icon := fmt.Sprintf("%s\\icons.ico", srcDir)

	dest := fmt.Sprintf("%s\\Demo.lnk", path)

	fmt.Println(dest)
	fmt.Println(src)

	err = CreateShortcut(dest, src, icon)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
