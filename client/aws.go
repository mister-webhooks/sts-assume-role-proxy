package client

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/mister-webhooks/sts-assume-role-proxy/protocol"
	"github.com/mister-webhooks/sts-assume-role-proxy/typedsocket"
)

// A proxyCredentialsProvider is an AWS CredentialsProvider that uses the sts-assume-role-proxy
// to obtain credentials
type proxyCredentialsProvider struct {
	socketPath string
}

// Create a new CredentialsProvider that uses the sts-assume-role-proxy at socketPath
func NewProxyCredentialsProvider(socketPath string) *proxyCredentialsProvider {
	return &proxyCredentialsProvider{socketPath: socketPath}
}

// CredentialProvder interface implementation
func (pcp proxyCredentialsProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	conn, err := typedsocket.Dial[protocol.AssumeRoleRequest, *protocol.RoleCredentials]("unix", pcp.socketPath)

	if err != nil {
		return aws.Credentials{}, err
	}

	if err = conn.Send(protocol.AssumeRoleRequest{}); err != nil {
		return aws.Credentials{}, err
	}

	reply := new(protocol.RoleCredentials)

	if err = conn.Recv(reply); err != nil {
		return aws.Credentials{}, nil
	}

	return aws.Credentials{
		AccessKeyID:     reply.AccessKeyId,
		SecretAccessKey: reply.SecretAccessKey,
		SessionToken:    reply.SessionToken,
		CanExpire:       true,
		Expires:         reply.Expiration,
	}, nil
}
