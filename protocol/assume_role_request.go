package protocol

import (
	"time"
)

type AssumeRoleRequest struct {
}

type RoleCredentials struct {
	Result          uint8     `json:"result"`
	Expiration      time.Time `json:"expiration"`
	AccessKeyId     string    `json:"access_key_id"`
	SecretAccessKey string    `json:"secret_access_key"`
	SessionToken    string    `json:"session_token"`
}
