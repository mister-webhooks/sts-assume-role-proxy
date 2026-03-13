package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/mister-webhooks/sts-assume-role-proxy/internal/unixsock"
	"github.com/mister-webhooks/sts-assume-role-proxy/protocol"
	"github.com/mister-webhooks/sts-assume-role-proxy/typedsocket"
	"mellium.im/sysexit"
)

const SOCKET_PATH = "/tmp/sts-assume-role-proxy.sock"

var AccessTable = map[string]map[uint]string{
	"[root]": {
		501: "arn:aws:iam::350784047695:role/webhooksd-task-role",
	},
}

func server() {
	if _, err := os.Stat(SOCKET_PATH); err == nil {
		if err := os.Remove(SOCKET_PATH); err != nil {
			log.Println("Error removing existing socket file:", err)
			return
		}
	}

	listener, err := typedsocket.NewTypedServer[*protocol.AssumeRoleRequest, protocol.RoleCredentials](func() (net.Listener, error) {
		return net.Listen("unix", SOCKET_PATH)
	})

	if err != nil {
		log.Fatalf("error creating listener: %v", err)
	}
	defer listener.Close()

	// Handle signals for graceful shutdown to ensure the socket file is removed
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	go func(c chan os.Signal) {
		<-c
		log.Println("Caught signal: shutting down and removing socket file.")
		listener.Close() // This will unlink the socket file
		os.Exit(0)
	}(sigc)

	log.Println("Server listening on", SOCKET_PATH)

	listener.Serve(context.Background(), func(ctx context.Context, tc *typedsocket.TypedConnection[protocol.RoleCredentials, *protocol.AssumeRoleRequest]) error {
		/*
		 * Obtain peer info
		 */
		unixConn, ok := (tc.Conn()).(*net.UnixConn)
		if !ok {
			return fmt.Errorf("unexpected socket type, expected Unix connection")
		}

		pinfo, err := unixsock.GetPeerInfo(unixConn)

		if err != nil {
			return fmt.Errorf("could not obtain peer pid: %w", err)
		}

		/*
		 * Accept request
		 */
		req := new(protocol.AssumeRoleRequest)

		err = tc.Recv(req)

		if err != nil {
			return fmt.Errorf("recv error: %w", err)
		}

		log.Printf("client %+v sent %+v", pinfo, req)

		/*
		 * Determine if the request is serviceable
		 */
		roleARN, ok := AccessTable[pinfo.Namespace][pinfo.Uid]

		if ok {
			log.Printf("client at pid %d {%s:%d} is allowed access to role '%s'", pinfo.Pid, pinfo.Namespace, pinfo.Uid, roleARN)

			hostname, err := os.Hostname()

			if err != nil {
				return fmt.Errorf("error retrieving own hostname: %w", err)
			}

			cfg, err := config.LoadDefaultConfig(ctx)

			if err != nil {
				return fmt.Errorf("could not load AWS credentials: %w", err)
			}

			stsService := sts.NewFromConfig(cfg)

			sessionIdentifier := fmt.Sprintf("%d@%s", pinfo.Pid, hostname)

			result, err := stsService.AssumeRole(ctx, &sts.AssumeRoleInput{RoleArn: &roleARN, RoleSessionName: &sessionIdentifier})

			if err != nil {
				return err
			}

			return tc.Send(protocol.RoleCredentials{
				Result:          0x0,
				Expiration:      *result.Credentials.Expiration,
				AccessKeyId:     *result.Credentials.AccessKeyId,
				SecretAccessKey: *result.Credentials.SecretAccessKey,
				SessionToken:    *result.Credentials.SessionToken,
			})
		} else {
			log.Printf("client at pid %d {%s:%d} is not allowed access", pinfo.Pid, pinfo.Namespace, pinfo.Uid)

			return tc.Send(protocol.RoleCredentials{
				Result:          0xFF,
				AccessKeyId:     "",
				SecretAccessKey: "",
				SessionToken:    "",
			})
		}
	})
}

func client() {
	log.Printf("current process id: %d", os.Getpid())
	conn, err := typedsocket.Dial[protocol.AssumeRoleRequest, *protocol.RoleCredentials]("unix", SOCKET_PATH)

	if err != nil {
		log.Fatal(err)
	}

	err = conn.Send(protocol.AssumeRoleRequest{})

	if err != nil {
		log.Fatal(err)
	}

	reply := new(protocol.RoleCredentials)

	err = conn.Recv(reply)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("received: %+v", reply)
}

func usage() {
	fmt.Println("usage: sts-assume-role-proxy <server | client>")
	os.Exit(int(sysexit.ErrUsage))
}

func main() {
	if len(os.Args) == 1 {
		usage()
	}

	switch os.Args[1] {
	case "server":
		server()
	case "client":
		client()
	default:
		usage()
	}
}
