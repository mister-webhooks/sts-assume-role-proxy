package protocol

import (
	"encoding/binary"
	"fmt"
	"time"
)

type AssumeRoleRequest struct {
	Rolename string
}

type RoleCredentials struct {
	Result          uint8
	Expiration      time.Time
	AccessKeyId     string
	SecretAccessKey string
	SessionToken    string
}

func (arr AssumeRoleRequest) MarshalBinary() ([]byte, error) {
	return LVString(arr.Rolename).MarshalBinary()
}

func (arr *AssumeRoleRequest) UnmarshalBinary(data []byte) error {
	lvs := new(LVString)
	err := lvs.UnmarshalBinary(data)

	if err != nil {
		return err
	}

	arr.Rolename = string(*lvs)

	return nil
}

func (rc RoleCredentials) MarshalBinary() ([]byte, error) {
	encodedExpiration, err := rc.Expiration.MarshalBinary()

	if err != nil {
		return nil, err
	}

	encodedAccessKeyId, err := LVString(rc.AccessKeyId).MarshalBinary()

	if err != nil {
		return nil, err
	}

	encodedSecretAccessKey, err := LVString(rc.SecretAccessKey).MarshalBinary()

	if err != nil {
		return nil, err
	}

	encodedSessionToken, err := LVString(rc.SessionToken).MarshalBinary()

	if err != nil {
		return nil, err
	}

	encoded := make([]byte, 0, 1+4+len(encodedExpiration)+len(encodedAccessKeyId)+len(encodedSecretAccessKey)+len(encodedSessionToken))

	encoded = append(encoded, rc.Result)
	encoded = binary.BigEndian.AppendUint32(encoded, uint32(len(encodedExpiration)))
	encoded = append(encoded, encodedExpiration...)
	encoded = append(encoded, encodedAccessKeyId...)
	encoded = append(encoded, encodedSecretAccessKey...)
	encoded = append(encoded, encodedSessionToken...)

	return encoded, nil
}

func (rc *RoleCredentials) UnmarshalBinary(data []byte) error {
	lvs := new(LVString)

	rc.Result = data[0]
	data = data[1:]

	// Decode ExpiresAt
	rc.Expiration = *new(time.Time)
	expiresAtLength := binary.BigEndian.Uint32(data[0:4])
	data = data[4:]

	err := rc.Expiration.UnmarshalBinary(data[:expiresAtLength])

	if err != nil {
		return err
	}

	data = data[expiresAtLength:]

	// Decode AccessKeyId
	err = lvs.UnmarshalBinary(data)
	if err != nil {
		return err
	}

	rc.AccessKeyId = string(*lvs)
	data = data[lvs.EncodedLength():]

	// Decode SecretAccessKey
	err = lvs.UnmarshalBinary(data)
	if err != nil {
		return err
	}

	rc.SecretAccessKey = string(*lvs)
	data = data[lvs.EncodedLength():]

	// Decode SessionToken
	err = lvs.UnmarshalBinary(data)
	if err != nil {
		return err
	}

	rc.SessionToken = string(*lvs)
	data = data[lvs.EncodedLength():]

	if len(data) != 0 {
		return fmt.Errorf("decoding error: %d bytes of RoleCredentials unread", len(data))
	}

	return nil
}
