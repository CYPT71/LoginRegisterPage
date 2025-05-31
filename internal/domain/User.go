package domain

import (
	"encoding/binary"
	"encoding/json"
	"log"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

type UserSessions struct {
	SessionData *webauthn.SessionData `json:"-"`
	SessionCred *webauthn.Credential  `json:"-"`
	UserIp      []string              `json:"-"`
	DisplayName string
	Jwt         string
	Expiration  time.Duration `json:"-"`
}

func (session *UserSessions) DeleteAfter(sessions map[string]*UserSessions) {

	timer := time.NewTimer(session.Expiration)

	go func() {
		<-timer.C // Wait for the timer to expire
		log.Printf("Session expired for user: %s", session.DisplayName)

		// Check user conditions for deletion
		user := UserModel{
			Username: session.DisplayName,
		}
		userModel := user.Get()

		if userModel.Password == "" && userModel.Incredentials == "" {
			// userModel.Delete()
			log.Printf("User deleted: %s", session.DisplayName)
		}

		// Delete the session from the sessions map
		delete(sessions, session.DisplayName)

	}()

}

type UserModel struct {
	Id            uint   `gorm:"primarykey;autoIncrement;not null"`
	Icon          string `gorm:"type:varchar(255);"`
	Username      string `gorm:"type:varchar(255);not null"`
	Email         string `gorm:"type:varchar(255);"`
	Password      string `gorm:"type:varchar(255);"`
	Permission    uint64 `gorm:"type:bigint"`
	Incredentials string `gorm:"column:credentials type:text"`
	webauthn.User `gorm:"-" json:"-"`
	Credentials   []webauthn.Credential `gorm:"-"`
}

func (user *UserModel) TableName() string {
	return "users"
}

func GetAllUsers() ([]UserModel, error) {
	var users []UserModel

	tx := Db.Find(&users)
	if tx.Error != nil {
		log.Println("Error:", tx.Error)
		return nil, tx.Error
	}

	return users, nil
}

func (user *UserModel) SaveCredentials() error {
	// @todo asure that credentials are transform to string
	var publicKeys []string
	for _, v := range user.Credentials {
		b, _ := json.Marshal(v)

		publicKeys = append(publicKeys, string(b))
	}
	user.Incredentials = strings.Join(publicKeys, ";")
	tx := Db.Save(&user)

	return tx.Error
}

func (user *UserModel) ParseCredentials() {
	for _, v := range strings.Split(user.Incredentials, ";") {
		cred := new(webauthn.Credential)
		err := json.Unmarshal([]byte(v), cred)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		user.Credentials = append(user.Credentials, *cred)
	}
}

func (user *UserModel) SetPassword(password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)
	return nil
}

func (user *UserModel) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil
}

func (user *UserModel) Create() error {
	user.Permission = Permissions["owner"]
	tx := Db.Create(user)

	return tx.Error
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

/* func randomUint64() uint64 {
	buf := make([]byte, 8)
	rand.Read(buf)
	return binary.LittleEndian.Uint64(buf)
} */

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
	u.Credentials = append(u.Credentials, cred)
}

// WebAuthnCredentials returns credentials owned by the user
func (u UserModel) WebAuthnCredentials() []webauthn.Credential {
	return u.Credentials
}

// CredentialExcludeList returns a CredentialDescriptor array filled
// with all the user's credentials
func (u UserModel) CredentialExcludeList() []protocol.CredentialDescriptor {

	credentialExcludeList := []protocol.CredentialDescriptor{}
	for _, cred := range u.Credentials {
		descriptor := protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: cred.ID,
		}
		credentialExcludeList = append(credentialExcludeList, descriptor)
	}

	return credentialExcludeList
}
