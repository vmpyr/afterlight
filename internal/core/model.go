package core

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type UserStatus string

const (
	StatusAlive   UserStatus = "ALIVE"
	StatusWarning UserStatus = "WARNING"
	StatusVerify  UserStatus = "VERIFICATION_REQUIRED"
	StatusDead    UserStatus = "CONFIRMED_DEAD"
)

type MessageType string

const (
	MsgText MessageType = "TEXT_MESSAGE"
	MsgFile MessageType = "FILE_UPLOAD"
)

type EncryptedBlob []byte

// Metadata field specific scanner
type Metadata map[string]string

func (m *Metadata) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		if value == nil {
			*m = make(map[string]string)
			return nil
		}
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, m)
}

func (m Metadata) Value() (driver.Value, error) {
	return json.Marshal(m)
}
