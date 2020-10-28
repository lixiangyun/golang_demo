package main

import (
	"fmt"
	ole "github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"os"
	"os/user"
	"path/filepath"
)

func CreateShortcut(dst, src string) error {
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
	return fmt.Sprintf("%s\\%s", dir, os.Args[0]), nil
}

func main()  {
	path, err := DeskTopPath()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	src, err := CurPath()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	dest := fmt.Sprintf("%s\\Demo.lnk", path)

	fmt.Println(dest)
	fmt.Println(src)

	err = CreateShortcut(dest, src)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
