package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	aws_config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/mister-webhooks/sts-assume-role-proxy/internal/config"
	"github.com/mister-webhooks/sts-assume-role-proxy/internal/unixsock"
	"github.com/mister-webhooks/sts-assume-role-proxy/protocol"
	"github.com/mister-webhooks/sts-assume-role-proxy/typedsocket"
	"github.com/urfave/cli/v2"
)

func server(config *config.ServerConfiguration) {
	if _, err := os.Stat(config.SocketPath); err == nil {
		if err := os.Remove(config.SocketPath); err != nil {
			log.Println("Error removing existing socket file:", err)
			return
		}
	}

	listener, err := typedsocket.NewTypedServer[*protocol.AssumeRoleRequest, protocol.RoleCredentials](func() (net.Listener, error) {
		return net.Listen("unix", config.SocketPath)
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

	log.Println("Server listening on", config.SocketPath)

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
		roleARN, ok := config.AccessTable.Lookup(pinfo.Namespace, pinfo.Uid)

		if ok {
			log.Printf("client at pid %d {%s:%d} is allowed access to role '%s'", pinfo.Pid, pinfo.Namespace, pinfo.Uid, roleARN)

			hostname, err := os.Hostname()

			if err != nil {
				return fmt.Errorf("error retrieving own hostname: %w", err)
			}

			cfg, err := aws_config.LoadDefaultConfig(ctx)

			if err != nil {
				return fmt.Errorf("could not load AWS credentials: %w", err)
			}

			stsService := sts.NewFromConfig(cfg)

			sessionIdentifier := fmt.Sprintf("%d@%s", pinfo.Pid, hostname)

			result, err := stsService.AssumeRole(ctx, &sts.AssumeRoleInput{RoleArn: &roleARN, RoleSessionName: &sessionIdentifier})

			if err != nil {
				return fmt.Errorf("could not request credentials on behalf of {%s:%d}: %w", pinfo.Namespace, pinfo.Uid, err)
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

func client(config *config.ServerConfiguration) {
	conn, err := typedsocket.Dial[protocol.AssumeRoleRequest, *protocol.RoleCredentials]("unix", config.SocketPath)

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

func loadConfig(filepath string) (*config.ServerConfiguration, error) {
	cfgData, err := os.ReadFile(filepath)

	if err != nil {
		return nil, fmt.Errorf("error reading configuration file: %s", err)
	}

	cfg, err := config.NewConfigurationFromYAML(cfgData)

	if err != nil {
		return nil, fmt.Errorf("error parsing configuration file: %s", err)
	}

	return cfg, nil
}

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "/usr/local/etc/sts_assume_role_proxy.conf",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "server",
				Usage: "run the sts-assume-role-proxy server",
				Action: func(ctx *cli.Context) error {
					cfg, err := loadConfig(ctx.String("config"))

					if err != nil {
						return err
					}

					server(cfg)
					return nil
				},
			}, {
				Name:  "client",
				Usage: "run the sts-assume-role-proxy client",
				Action: func(ctx *cli.Context) error {
					cfg, err := loadConfig(ctx.String("config"))

					if err != nil {
						log.Fatal(err)
					}

					client(cfg)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
