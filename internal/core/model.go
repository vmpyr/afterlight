package core

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
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
	MsgS3   MessageType = "S3_OBJECT_LINK"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	Email         string     `json:"email"`
	CurrentStatus UserStatus `json:"current_status"`
	CreatedAt     time.Time  `json:"created_at"`
}

type CreateVaultRequest struct {
	VaultName string `json:"vault_name"`
	Hint      string `json:"hint,omitempty"`
	KdfSalt   string `json:"kdf_salt"`
}

type EncryptedBlob []byte
type CreateArtifactRequest struct {
	MessageType   MessageType   `json:"message_type"`
	EncryptedBlob EncryptedBlob `json:"encrypted_blob"`
	IV            string        `json:"iv"`
}

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
