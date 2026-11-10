//go:build windows && (amd64 || arm64)

package output

// The device other processes play on: Windows keeps, for each app, the
// output device its sound goes to, as Settings › Sound › Volume mixer
// sets it, through its AudioPolicyConfig, which it doesn't document but
// which sound mixers use. The Video side's player is a WebView2 page,
// whose sound plays in a process of WebView2's own, which RouteProcesses
// sends to a device.

import (
	"errors"
	"runtime"
	"syscall"
	"unsafe"
)

var (
	combase                      = syscall.NewLazyDLL("combase.dll")
	procRoGetActivationFactory   = combase.NewProc("RoGetActivationFactory")
	procWindowsCreateString      = combase.NewProc("WindowsCreateString")
	procWindowsDeleteString      = combase.NewProc("WindowsDeleteString")
	procCreateToolhelp32Snapshot = kernel32.NewProc("CreateToolhelp32Snapshot")
	procProcess32FirstW          = kernel32.NewProc("Process32FirstW")
	procProcess32NextW           = kernel32.NewProc("Process32NextW")
	procCloseHandle              = kernel32.NewProc("CloseHandle")

	// The factory's interface since Windows 10 21H2, and before.
	iidPolicyConfig21H2      = guid{0xab3d4648, 0xe242, 0x459f, [8]byte{0xb0, 0x2f, 0x54, 0x1c, 0x70, 0x30, 0x63, 0x24}}
	iidPolicyConfigDownlevel = guid{0x2a59116d, 0x6c4f, 0x45e0, [8]byte{0xa7, 0x4f, 0x70, 0x7e, 0x3f, 0xef, 0x92, 0x58}}
)

const (
	policyConfigClass = "Windows.Media.Internal.AudioPolicyConfig"
	// setPersistedDefault is SetPersistedDefaultAudioEndpoint's slot in
	// the vtable: after IInspectable's six methods and nineteen others.
	setPersistedDefault = 6 + 19
	// devicePrefix and deviceSuffix make an endpoint's ID the device
	// interface the policy takes.
	devicePrefix = `\\?\SWD#MMDEVAPI#`
	deviceSuffix = `#{e6327cad-dcec-4949-ae8a-991e976a79d2}`

	th32csSnapProcess = 0x2
	eMultimedia       = 1
)

// processEntry is PROCESSENTRY32W.
type processEntry struct {
	size            uint32
	usage           uint32
	processID       uint32
	defaultHeapID   uintptr
	moduleID        uint32
	threads         uint32
	parentProcessID uint32
	priClassBase    int32
	flags           uint32
	exeFile         [260]uint16
}

// ProcessTree returns root and the processes it started, and theirs.
func ProcessTree(root uint32) []uint32 {
	snap, _, _ := syscall.SyscallN(procCreateToolhelp32Snapshot.Addr(), th32csSnapProcess, 0)
	if snap == 0 || snap == ^uintptr(0) {
		return []uint32{root}
	}
	defer syscall.SyscallN(procCloseHandle.Addr(), snap)
	children := map[uint32][]uint32{}
	var e processEntry
	e.size = uint32(unsafe.Sizeof(e))
	for r, _, _ := syscall.SyscallN(procProcess32FirstW.Addr(), snap, uintptr(unsafe.Pointer(&e))); r != 0; r, _, _ = syscall.SyscallN(procProcess32NextW.Addr(), snap, uintptr(unsafe.Pointer(&e))) {
		if e.processID != e.parentProcessID {
			children[e.parentProcessID] = append(children[e.parentProcessID], e.processID)
		}
	}
	out := []uint32{root}
	seen := map[uint32]bool{root: true}
	for i := 0; i < len(out); i++ {
		for _, c := range children[out[i]] {
			if !seen[c] {
				seen[c] = true
				out = append(out, c)
			}
		}
	}
	return out
}

func hstring(s string) (uintptr, error) {
	u, err := syscall.UTF16FromString(s)
	if err != nil {
		return 0, err
	}
	var h uintptr
	if r, _, _ := syscall.SyscallN(procWindowsCreateString.Addr(), uintptr(unsafe.Pointer(&u[0])), uintptr(len(u)-1), uintptr(unsafe.Pointer(&h))); int32(r) < 0 {
		return 0, check(r)
	}
	return h, nil
}

// RouteProcesses sends the sound of the processes pids to the output
// device id, an endpoint ID as Devices tells, or, for "", back to the
// system's default device.
func RouteProcesses(pids []uint32, id string) error {
	if err := procRoGetActivationFactory.Find(); err != nil {
		return errors.New("this version of Windows can't choose the device of an app")
	}
	done := make(chan error, 1)
	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		undo, err := initCOM()
		if err != nil {
			done <- err
			return
		}
		defer undo()
		done <- routeProcesses(pids, id)
	}()
	return <-done
}

func routeProcesses(pids []uint32, id string) error {
	class, err := hstring(policyConfigClass)
	if err != nil {
		return err
	}
	defer syscall.SyscallN(procWindowsDeleteString.Addr(), class)
	var factory unsafe.Pointer
	r, _, _ := syscall.SyscallN(procRoGetActivationFactory.Addr(), class, uintptr(unsafe.Pointer(&iidPolicyConfig21H2)), uintptr(unsafe.Pointer(&factory)))
	if int32(r) < 0 {
		r, _, _ = syscall.SyscallN(procRoGetActivationFactory.Addr(), class, uintptr(unsafe.Pointer(&iidPolicyConfigDownlevel)), uintptr(unsafe.Pointer(&factory)))
	}
	if int32(r) < 0 || factory == nil {
		return errors.New("this version of Windows can't choose the device of an app")
	}
	vtbl := *(*[setPersistedDefault + 3]uintptr)(*(*unsafe.Pointer)(factory))
	defer syscall.SyscallN(vtbl[2], uintptr(factory)) // Release
	var device uintptr
	if id != "" {
		if device, err = hstring(devicePrefix + id + deviceSuffix); err != nil {
			return err
		}
		defer syscall.SyscallN(procWindowsDeleteString.Addr(), device)
	}
	var first error
	for _, pid := range pids {
		for _, role := range []uintptr{eConsole, eMultimedia} {
			r, _, _ := syscall.SyscallN(vtbl[setPersistedDefault], uintptr(factory), uintptr(pid), eRender, role, device)
			if err := check(r); err != nil && first == nil {
				first = err
			}
		}
	}
	return first
}
