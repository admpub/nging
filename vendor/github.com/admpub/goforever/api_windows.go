//go:build windows

package goforever

import (
	"fmt"
	"os"
	"runtime"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	advapi32                    = syscall.NewLazyDLL("advapi32.dll")
	procCreateProcessWithLogonW = advapi32.NewProc("CreateProcessWithLogonW")
	logonProc                   = advapi32.NewProc("LogonUserW")
	impersonateProc             = advapi32.NewProc("ImpersonateLoggedOnUser")
)

const (
	Logon32LogonNetwork     = uintptr(3)
	Logon32LogonInteractive = uintptr(2)
	logon32ProviderDefault  = uintptr(0)

	logonWithProfile        uint32 = 0x00000001
	logonNetCredentialsOnly uint32 = 0x00000002
	createDefaultErrorMode  uint32 = 0x04000000
	createNewProcessGroup   uint32 = 0x00000200
)

type tokenType int

const (
	TokenUnknown tokenType = iota
	TokenPrimary
	TokenImpersonation
	TokenLinked
)

// CreateProcessWithLogonW ...
func CreateProcessWithLogonW(
	username *uint16,
	domain *uint16,
	password *uint16,
	logonFlags uint32,
	applicationName *uint16,
	commandLine *uint16,
	creationFlags uint32,
	environment *uint16,
	currentDirectory *uint16,
	startupInfo *syscall.StartupInfo,
	processInformation *syscall.ProcessInformation) error {
	r1, _, e1 := procCreateProcessWithLogonW.Call(
		uintptr(unsafe.Pointer(username)),
		uintptr(unsafe.Pointer(domain)),
		uintptr(unsafe.Pointer(password)),
		uintptr(logonFlags),
		uintptr(unsafe.Pointer(applicationName)),
		uintptr(unsafe.Pointer(commandLine)),
		uintptr(creationFlags),
		uintptr(unsafe.Pointer(environment)), // env
		uintptr(unsafe.Pointer(currentDirectory)),
		uintptr(unsafe.Pointer(startupInfo)),
		uintptr(unsafe.Pointer(processInformation)))
	runtime.KeepAlive(username)
	runtime.KeepAlive(domain)
	runtime.KeepAlive(password)
	runtime.KeepAlive(applicationName)
	runtime.KeepAlive(commandLine)
	runtime.KeepAlive(environment)
	runtime.KeepAlive(currentDirectory)
	runtime.KeepAlive(startupInfo)
	runtime.KeepAlive(processInformation)
	if int(r1) == 0 {
		return os.NewSyscallError("CreateProcessWithLogonW", e1)
	}
	return nil
}

// ListToEnvironmentBlock ...
func ListToEnvironmentBlock(list *[]string) (*uint16, error) {
	if list == nil {
		return nil, nil
	}
	size := 1
	ulines := make([][]uint16, len(*list))
	for i, v := range *list {
		uline, err := syscall.UTF16FromString(v)
		if err != nil {
			return nil, err
		}
		size += len(uline)
		ulines[i] = uline
	}
	result := make([]uint16, size)
	tail := 0
	for i := range *list {
		uline := ulines[i]
		copy(result[tail:], uline)
		tail += len(uline)
	}
	result[tail] = 0
	return &result[0], nil
}

// CreateProcessWithLogon creates a process giving user credentials
// Ref: https://github.com/hosom/honeycred/blob/master/honeycred.go
func CreateProcessWithLogon(username string, password string, domain string, path string, cmdLine string, workDir string) (*syscall.ProcessInformation, error) {
	user, err := syscall.UTF16PtrFromString(username)
	if err != nil {
		return nil, err
	}
	dom, err := syscall.UTF16PtrFromString(domain)
	if err != nil {
		return nil, err
	}
	pass, err := syscall.UTF16PtrFromString(password)
	if err != nil {
		return nil, err
	}
	logonFlags := logonWithProfile // changed
	applicationName, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return nil, err
	}
	commandLine, err := syscall.UTF16PtrFromString(cmdLine)
	if err != nil {
		return nil, err
	}
	creationFlags := createDefaultErrorMode
	environment, err := ListToEnvironmentBlock(nil)
	if err != nil {
		return nil, err
	}
	if len(workDir) == 0 {
		workDir = `c:\programdata`
	}
	currentDirectory, err := syscall.UTF16PtrFromString(workDir)
	if err != nil {
		return nil, err
	}
	startupInfo := &syscall.StartupInfo{}
	processInfo := &syscall.ProcessInformation{}

	err = CreateProcessWithLogonW(
		user,
		dom,
		pass,
		logonFlags,
		applicationName,
		commandLine,
		creationFlags,
		environment,
		currentDirectory,
		startupInfo,
		processInfo)
	return processInfo, err
}

// LogonUser attempts to log a user on to the local computer to generate a token.
func LogonUser(user, pass string, logonType ...uintptr) (token syscall.Token, err error) {
	// ".\0" meaning "this computer:
	domain := [2]uint16{uint16('.'), 0}
	var pu, pp []uint16
	if pu, err = syscall.UTF16FromString(user); err != nil {
		return
	}
	if pp, err = syscall.UTF16FromString(pass); err != nil {
		return
	}

	var logonTyp uintptr
	if len(logonType) > 0 {
		logonTyp = logonType[0]
	}
	if logonTyp == 0 {
		logonTyp = Logon32LogonNetwork
		//logonTyp = Logon32LogonInteractive
	}

	if rc, _, ec := syscall.SyscallN(logonProc.Addr(),
		uintptr(unsafe.Pointer(&pu[0])),
		uintptr(unsafe.Pointer(&domain[0])),
		uintptr(unsafe.Pointer(&pp[0])),
		logonTyp,
		logon32ProviderDefault,
		uintptr(unsafe.Pointer(&token))); rc == 0 {
		err = error(ec)
	}
	return
}

// Impersonate attempts to impersonate the user.
func Impersonate(user string, pass string) error {
	token, err := LogonUser(user, pass)
	if err != nil {
		return err
	}
	defer token.Close()
	return TokenImpersonate(token)
}

func TokenImpersonate(token syscall.Token) error {
	if rc, _, ec := syscall.SyscallN(impersonateProc.Addr(), uintptr(token), 0, 0); rc == 0 {
		return error(ec)
	}
	return nil
}

func DuplicateTokenEx(token syscall.Token, tokenType tokenType) (syscall.Token, error) {
	defer token.Close()
	t := windows.Token(token)
	var duplicatedToken windows.Token
	switch tokenType {
	case TokenPrimary:
		if err := windows.DuplicateTokenEx(t, windows.MAXIMUM_ALLOWED, nil, windows.SecurityDelegation, windows.TokenPrimary, &duplicatedToken); err != nil {
			return 0, fmt.Errorf("error while DuplicateTokenEx: %w", err)
		}
	case TokenImpersonation:
		if err := windows.DuplicateTokenEx(t, windows.MAXIMUM_ALLOWED, nil, windows.SecurityImpersonation, windows.TokenImpersonation, &duplicatedToken); err != nil {
			return 0, fmt.Errorf("error while DuplicateTokenEx: %w", err)
		}
	case TokenLinked:
		if err := windows.DuplicateTokenEx(t, windows.MAXIMUM_ALLOWED, nil, windows.SecurityDelegation, windows.TokenPrimary, &duplicatedToken); err != nil {
			return 0, fmt.Errorf("error while DuplicateTokenEx: %w", err)
		}
		dt, err := duplicatedToken.GetLinkedToken()
		windows.CloseHandle(windows.Handle(duplicatedToken))
		if err != nil {
			return 0, fmt.Errorf("error while getting LinkedToken: %w", err)
		}
		duplicatedToken = dt
	}
	if windows.Handle(duplicatedToken) == windows.InvalidHandle {
		return 0, fmt.Errorf("invalid duplicated token")
	}
	return syscall.Token(duplicatedToken), nil
}
