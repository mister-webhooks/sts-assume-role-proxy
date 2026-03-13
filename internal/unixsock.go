package unixsock

/*
#include <sys/types.h>
#include <sys/socket.h>
#include <sys/un.h>
#include <unistd.h>
#include <errno.h>


pid_t get_socket_pid(int socket) {
	pid_t peer_pid;
	socklen_t len = sizeof(peer_pid);
	if (getsockopt(socket, 0, LOCAL_PEERPID, &peer_pid, &len) == -1) {
    return -1;
	}
	return peer_pid;
}
*/
import "C"
import (
	"fmt"
	"net"
)

type Pid C.pid_t

func GetClientPID(conn net.Conn) (Pid, error) {
	unixConn, ok := conn.(*net.UnixConn)
	if !ok {
		return -1, fmt.Errorf("unexpected socket type, expected Unix connection")
	}

	raw, err := unixConn.SyscallConn()
	if err != nil {
		return -1, fmt.Errorf("error opening raw connection: %s", err)
	}

	var callbackErr error
	var pid C.pid_t

	controlErr := raw.Control(func(fd uintptr) {
		pid, callbackErr = C.get_socket_pid(C.int(fd))
	})

	if callbackErr != nil {
		return -1, fmt.Errorf("error getting socket info: %w", callbackErr)
	}

	if controlErr != nil {
		return -1, fmt.Errorf("error attaching callback: %w", controlErr)
	}

	return Pid(pid), nil
}
