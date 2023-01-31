package domain

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"gorm.io/gorm"
)

var Db *gorm.DB

type UserSessions struct {
	SessionData *webauthn.SessionData `json:"-"`
	SessionCred *webauthn.Credential  `json:"-"`
	DisplayName string
	Jwt         string
	Expiration  uint64 `json:"-"`
}

func (session *UserSessions) DeleteAfter(sessions map[string]*UserSessions) {
	for i := session.Expiration; i >= 0; i-- {
		time.Sleep(1)
	}

	log.Printf("user delete")
	delete(sessions, session.DisplayName)
}

type UserModel struct {
	Id            uint   `gorm:"primarykey;autoIncrement;not null"`
	Username      string `gorm:"type:varchar(255);not null"`
	Icon          []byte `gorm:"type:blob;"`
	Email         string `gorm:"type:varchar(255);"`
	Password      string `gorm:"type:varchar(255);"`
	Permissions   uint64 `gorm:"type:numeric;not null"`
	Credentials   string `gorm:"type:text"`
	webauthn.User `gorm:"-" json:"-"`
	Credentals    []webauthn.Credential `gorm:"-"`
}

func (user *UserModel) TableName() string {
	return "users"
}

func (user *UserModel) SaveCredentials() {
	// @todo asure that credentials are transform to string
	var publicKeys []string
	for _, v := range user.Credentals {
		b, _ := json.Marshal(v)

		publicKeys = append(publicKeys, string(b))
	}
	user.Credentials = strings.Join(publicKeys, ";")
	Db.Save(&user)
}

func (user *UserModel) ParseCredentials() {
	for _, v := range strings.Split(user.Credentials, ";") {
		cred := new(webauthn.Credential)
		err := json.Unmarshal([]byte(v), cred)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		user.Credentals = append(user.Credentals, *cred)
	}
}

func (user *UserModel) Create() {
	Db.Create(user)
}

/*
if user fin return true else false
*/
func (user *UserModel) Find() bool {
	tx := Db.Where("username = ?", user.Username).Find(user)
	return tx.RowsAffected != 0
}
func (user *UserModel) Get() *UserModel {

	tx := Db.Where("username = ?", user.Username).Find(user)
	if tx.RowsAffected == 0 {
		return nil
	}
	return user

}
func (user *UserModel) Delete() {
	Db.Delete(user)
}

func (user *UserModel) Update() {
	Db.Save(&user)
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
	u.Credentals = append(u.Credentals, cred)
}

// WebAuthnCredentials returns credentials owned by the user
func (u UserModel) WebAuthnCredentials() []webauthn.Credential {
	return u.Credentals
}

// CredentialExcludeList returns a CredentialDescriptor array filled
// with all the user's credentials
func (u UserModel) CredentialExcludeList() []protocol.CredentialDescriptor {

	credentialExcludeList := []protocol.CredentialDescriptor{}
	for _, cred := range u.Credentals {
		descriptor := protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: cred.ID,
		}
		credentialExcludeList = append(credentialExcludeList, descriptor)
	}

	return credentialExcludeList
}
