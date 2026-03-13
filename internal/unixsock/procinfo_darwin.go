package unixsock

/*
#include <sys/types.h>
#include <sys/socket.h>
#include <sys/un.h>
#include <unistd.h>
#include <sys/ucred.h>
#include <errno.h>
*/
import "C"
import (
	"fmt"
	"net"
	"unsafe"
)

func GetPeerInfo(unixConn *net.UnixConn) (*PeerInfo, error) {
	raw, err := unixConn.SyscallConn()
	if err != nil {
		return nil, fmt.Errorf("error opening raw connection: %s", err)
	}

	var callbackErr error

	pid_length := C.socklen_t(C.sizeof_pid_t)
	var pid C.pid_t

	xucred_len := C.socklen_t(C.sizeof_struct_xucred)
	var xucred C.struct_xucred

	controlErr := raw.Control(func(fd uintptr) {
		var rv C.int

		rv, callbackErr = C.getsockopt(C.int(fd), C.SOL_LOCAL, C.LOCAL_PEERPID, unsafe.Pointer(&pid), &pid_length)
		if rv == -1 {
			return
		}
		_, callbackErr = C.getsockopt(C.int(fd), C.SOL_LOCAL, C.LOCAL_PEERCRED, unsafe.Pointer(&xucred), &xucred_len)
	})

	if callbackErr != nil {
		return nil, fmt.Errorf("error getting socket info: %w", callbackErr)
	}

	if controlErr != nil {
		return nil, fmt.Errorf("error attaching callback: %w", controlErr)
	}

	return &PeerInfo{
		Namespace: "[root]",
		Pid:       uint(pid),
		Uid:       uint(xucred.cr_uid),
	}, nil
}
