package protocol

import (
	"time"
)

type AssumeRoleRequest struct {
}

type RoleCredentials struct {
	Result          CredentialsResult `json:"result"`
	Expiration      time.Time         `json:"expiration"`
	AccessKeyId     string            `json:"access_key_id"`
	SecretAccessKey string            `json:"secret_access_key"`
	SessionToken    string            `json:"session_token"`
}

type CredentialsResult uint8

const (
	CredentialsResultSuccess    CredentialsResult = 0x00
	CredentialsResultProxyError CredentialsResult = 0xFE
	CredentialsResultForbidden  CredentialsResult = 0xFF
)
