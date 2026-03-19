package unixsock

/*
#cgo LDFLAGS: -lutil
#include <sys/types.h>
#include <sys/socket.h>
#include <sys/ucred.h>
#include <sys/un.h>
#include <stdio.h>
#include <unistd.h>
#include <sys/user.h>
#include <libutil.h>
#include <sys/uio.h>
#include <stdlib.h>
#include <sys/param.h>
#include <sys/jail.h>
#include <sys/cdefs.h>

pid_t get_xucred_pid(const struct xucred *creds) {
	return creds->cr_pid;
}

int get_jail_name(int jid, char *name, int namelen) {
	struct iovec iov[4];

	iov[0].iov_base = __DECONST(char *, "jid");
	iov[0].iov_len = sizeof("jid");
	iov[1].iov_base = &jid;
	iov[1].iov_len = sizeof(jid);

	iov[2].iov_base = __DECONST(char *, "name");
	iov[2].iov_len = sizeof("name");
	iov[3].iov_base = name;
	iov[3].iov_len = namelen;

	return jail_get(iov, 4, 0);
}

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

	xucred_len := C.socklen_t(C.sizeof_struct_xucred)
	var xucred C.struct_xucred
	var procinfo *C.struct_kinfo_proc
	var jailname string

	controlErr := raw.Control(func(fd uintptr) {
		var rv C.int

		rv, callbackErr = C.getsockopt(C.int(fd), C.SOL_LOCAL, C.LOCAL_PEERCRED, unsafe.Pointer(&xucred), &xucred_len)
		if rv == -1 {
			return
		}

		procinfo, callbackErr = C.kinfo_getproc(C.get_xucred_pid(&xucred))

		c_jailname := make([]C.char, C.MAXHOSTNAMELEN)

		if procinfo.ki_jid == 0 {
			jailname = "[root]"
			return
		}

		rv, callbackErr = C.get_jail_name(procinfo.ki_jid, &c_jailname[0], C.MAXHOSTNAMELEN)
		if rv == -1 {
			return
		}

		jailname = C.GoString(&c_jailname[0])
	})

	if callbackErr != nil {
		return nil, fmt.Errorf("error getting socket info: %w", callbackErr)
	}

	if controlErr != nil {
		return nil, fmt.Errorf("error attaching callback: %w", controlErr)
	}

	return &PeerInfo{
		Namespace: jailname,
		Pid:       uint(C.get_xucred_pid(&xucred)),
		Uid:       uint(xucred.cr_uid),
	}, nil
}
