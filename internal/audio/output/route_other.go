//go:build !(windows && (amd64 || arm64))

package output

import "errors"

// ProcessTree returns root and the processes it started: root alone here.
func ProcessTree(root uint32) []uint32 { return []uint32{root} }

// RouteProcesses sends the sound of processes to a device, on Windows.
func RouteProcesses(pids []uint32, id string) error {
	return errors.New("only Windows can choose the device of an app")
}
