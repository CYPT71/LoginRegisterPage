package main

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"log"
	"strings"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
)

type UserModel struct {
	Id            uint   `gorm:"primarykey;autoIncrement;not null"`
	Username      string `gorm:"type:varchar(255);not null"`
	Email         string `gorm:"type:varchar(255);"`
	Password      string `gorm:"type:varchar(255);"`
	Roles         uint64 `gorm:"type:numeric;not null"`
	Credentials   string `gorm:"type:text`
	webauthn.User `gorm:"-" json:"-"`
	credentals    []webauthn.Credential `gorm:"-" json:"-"`
}

func (user *UserModel) TableName() string {
	return "users"
}

func (user *UserModel) saveCredentials() {
	// @todo asure that credentials are transform to string
	var publicKeys []string
	for _, v := range user.credentals {
		b, _ := json.Marshal(v)

		publicKeys = append(publicKeys, string(b))
	}
	user.Credentials = strings.Join(publicKeys, ";")
	db.Save(&user)
}

func (user *UserModel) parseCredentials() {
	for _, v := range strings.Split(user.Credentials, ";") {
		cred := new(webauthn.Credential)
		err := json.Unmarshal([]byte(v), cred)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		user.credentals = append(user.credentals, *cred)
	}
}

func (user *UserModel) Create() {
	db.Create(user)
}

func (user *UserModel) Find() bool {
	tx := db.Where("username = ?", user.Username).Find(user)
	return tx.RowsAffected != 0
}
func (user *UserModel) Get() *UserModel {

	tx := db.Where("username = ?", user.Username).Find(user)
	if tx.RowsAffected == 0 {
		return nil
	}
	return user

}
func (user *UserModel) Delete() {
	db.Delete(user)
}

func (user *UserModel) Update() {
	db.Save(&user)
}

func randomUint64() uint64 {
	buf := make([]byte, 8)
	rand.Read(buf)
	return binary.LittleEndian.Uint64(buf)
}

// WebAuthnID returns the user's ID
func (u UserModel) WebAuthnID() []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, uint64(u.Id))
	return buf
}

// WebAuthnName returns the user's username
func (u UserModel) WebAuthnName() string {
	return u.Username
}

// WebAuthnDisplayName returns the user's display name
func (u UserModel) WebAuthnDisplayName() string {
	return u.Username
}

// WebAuthnIcon is not (yet) implemented
func (u UserModel) WebAuthnIcon() string {
	return ""
}

// AddCredential associates the credential to the user
func (u *UserModel) AddCredential(cred webauthn.Credential) {
	u.credentals = append(u.credentals, cred)
}

// WebAuthnCredentials returns credentials owned by the user
func (u UserModel) WebAuthnCredentials() []webauthn.Credential {
	return u.credentals
}

// CredentialExcludeList returns a CredentialDescriptor array filled
// with all the user's credentials
func (u UserModel) CredentialExcludeList() []protocol.CredentialDescriptor {

	credentialExcludeList := []protocol.CredentialDescriptor{}
	for _, cred := range u.credentals {
		descriptor := protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: cred.ID,
		}
		credentialExcludeList = append(credentialExcludeList, descriptor)
	}

	return credentialExcludeList
}
