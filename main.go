package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const (
	ERROR_SUCCESS      = 0
	ERROR_MORE_DATA    = 234
	RmForceShutdown    = 1
	CCH_RM_SESSION_KEY = 256
)

var (
	rstrtmgr                = syscall.NewLazyDLL("Rstrtmgr.dll")
	procRmStartSession      = rstrtmgr.NewProc("RmStartSession")
	procRmRegisterResources = rstrtmgr.NewProc("RmRegisterResources")
	procRmGetList           = rstrtmgr.NewProc("RmGetList")
	procRmShutdown          = rstrtmgr.NewProc("RmShutdown")
	procRmEndSession        = rstrtmgr.NewProc("RmEndSession")
)

func unlockFile(filePath string) error {
	var sessionHandle uint32
	sessionKey := make([]uint16, CCH_RM_SESSION_KEY)

	ret, _, _ := procRmStartSession.Call(
		uintptr(unsafe.Pointer(&sessionHandle)),
		0,
		uintptr(unsafe.Pointer(&sessionKey[0])),
	)

	if ret != ERROR_SUCCESS {
		return fmt.Errorf("RmStartSession failed with error: %d", ret)
	}

	defer procRmEndSession.Call(uintptr(sessionHandle))

	filePathPtr, err := syscall.UTF16PtrFromString(filePath)
	if err != nil {
		return err
	}

	ret, _, _ = procRmRegisterResources.Call(
		uintptr(sessionHandle),
		1,
		uintptr(unsafe.Pointer(&filePathPtr)),
		0,
		0,
		0,
		0,
	)

	if ret != ERROR_SUCCESS {
		return fmt.Errorf("RmRegisterResources failed with error: %d", ret)
	}

	var pnProcInfoNeeded, pnProcInfo, lpdwRebootReasons uint32

	ret, _, _ = procRmGetList.Call(
		uintptr(sessionHandle),
		uintptr(unsafe.Pointer(&pnProcInfoNeeded)),
		uintptr(unsafe.Pointer(&pnProcInfo)),
		0,
		uintptr(unsafe.Pointer(&lpdwRebootReasons)),
	)

	if ret != ERROR_SUCCESS && ret != ERROR_MORE_DATA {
		return fmt.Errorf("RmGetList failed with error: %d", ret)
	}

	if pnProcInfoNeeded > 0 {
		ret, _, _ = procRmShutdown.Call(
			uintptr(sessionHandle),
			RmForceShutdown,
			0,
		)

		if ret != ERROR_SUCCESS {
			return fmt.Errorf("RmShutdown failed with error: %d", ret)
		}
	} else {
		fmt.Println("File is not locked")
	}

	return nil
}

func fetchFile(filePath string) (int, error) {
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		err := unlockFile(filePath)
		if err != nil {
			fmt.Printf("Attempt %d failed to unlock: %v\n", i+1, err)
			continue
		}

		data, err := os.ReadFile(filePath)
		if err == nil {
			return len(data), nil
		}
		fmt.Printf("Attempt %d failed to read: %v\n", i+1, err)
	}
	return 0, fmt.Errorf("failed to fetch data after %d attempts", maxRetries)
}

func main() {
	filePath := os.Args[1]
	length, err := fetchFile(filePath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("File length: %d bytes\n", length)
}
