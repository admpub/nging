package winpty

import (
	"fmt"
	"syscall"
	"unsafe"
)

func createAgentCfg(flags uint32) (uintptr, error) {
	var errorPtr uintptr

	err := winpty_error_free.Find() // check if dll available
	if err != nil {
		return uintptr(0), err
	}

	defer winpty_error_free.Call(errorPtr)

	agentCfg, _, _ := winpty_config_new.Call(
		uintptr(flags),
		uintptr(0), // winpty expects a UINT64 so we need to pad the call on 386
		uintptr(unsafe.Pointer(errorPtr)),
	)

	if agentCfg == uintptr(0) {
		return 0, fmt.Errorf("Unable to create agent config, %s", GetErrorMessage(errorPtr))
	}

	return agentCfg, nil
}

func createSpawnCfg(flags uint32, appname, cmdline, cwd string, env []string) (uintptr, error) {
	var errorPtr uintptr
	defer winpty_error_free.Call(errorPtr)

	cmdLineStr, err := syscall.UTF16PtrFromString(cmdline)
	if err != nil {
		return 0, fmt.Errorf("Failed to convert cmd to pointer.")
	}

	appNameStr, err := syscall.UTF16PtrFromString(appname)
	if err != nil {
		return 0, fmt.Errorf("Failed to convert app name to pointer.")
	}

	cwdStr, err := syscall.UTF16PtrFromString(cwd)
	if err != nil {
		return 0, fmt.Errorf("Failed to convert working directory to pointer.")
	}

	envStr, err := UTF16PtrFromStringArray(env)

	if err != nil {
		return 0, fmt.Errorf("Failed to convert cmd to pointer.")
	}

	spawnCfg, _, _ := winpty_spawn_config_new.Call(
		uintptr(flags),
		uintptr(0), // winpty expects a UINT64 so we need to pad the call on 386
		uintptr(unsafe.Pointer(appNameStr)),
		uintptr(unsafe.Pointer(cmdLineStr)),
		uintptr(unsafe.Pointer(cwdStr)),
		uintptr(unsafe.Pointer(envStr)),
		uintptr(unsafe.Pointer(errorPtr)),
	)

	if spawnCfg == uintptr(0) {
		return 0, fmt.Errorf("Unable to create spawn config, %s", GetErrorMessage(errorPtr))
	}

	return spawnCfg, nil
}
