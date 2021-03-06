package entity

import (
	"database/sql"
	"sort"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/photoprism/photoprism/internal/form"
	"github.com/photoprism/photoprism/internal/service"
	"github.com/photoprism/photoprism/internal/service/webdav"
	"github.com/ulule/deepcopier"
)

// Account represents a remote service account for uploading, downloading or syncing media files.
type Account struct {
	ID           uint   `gorm:"primary_key"`
	AccName      string `gorm:"type:varchar(128);"`
	AccOwner     string `gorm:"type:varchar(128);"`
	AccURL       string `gorm:"type:varbinary(512);"`
	AccType      string `gorm:"type:varbinary(256);"`
	AccKey       string `gorm:"type:varbinary(256);"`
	AccUser      string `gorm:"type:varbinary(256);"`
	AccPass      string `gorm:"type:varbinary(256);"`
	AccError     string `gorm:"type:varbinary(512);"`
	AccShare     bool
	AccSync      bool
	RetryLimit   uint
	SharePath    string `gorm:"type:varbinary(256);"`
	ShareSize    string `gorm:"type:varbinary(16);"`
	ShareExpires uint
	ShareExif    bool
	ShareSidecar bool
	SyncPath     string `gorm:"type:varbinary(256);"`
	SyncInterval uint
	SyncUpload   bool
	SyncDownload bool
	SyncDelete   bool
	SyncRaw      bool
	SyncVideo    bool
	SyncSidecar  bool
	SyncStart    sql.NullTime
	SyncedAt     sql.NullTime `deepcopier:"skip"`
	CreatedAt    time.Time    `deepcopier:"skip"`
	UpdatedAt    time.Time    `deepcopier:"skip"`
	DeletedAt    *time.Time   `deepcopier:"skip" sql:"index"`
}

// CreateAccount creates a new account entity in the database.
func CreateAccount(form form.Account, db *gorm.DB) (model *Account, err error) {
	model = &Account{}

	err = model.Save(form, db)

	return model, err
}

// Save updates the entity using form data and stores it in the database.
func (m *Account) Save(form form.Account, db *gorm.DB) error {
	if err := deepcopier.Copy(m).From(form); err != nil {
		return err
	}

	if m.AccType != string(service.TypeWebDAV) {
		// TODO: Only WebDAV supported at the moment
		m.AccShare = false
		m.AccSync = false
	}

	if m.SharePath == "" {
		m.SharePath = "/"
	}

	if m.SyncPath == "" {
		m.SyncPath = "/"
	}

	return db.Save(m).Error
}

// Delete deletes the entity from the database.
func (m *Account) Delete(db *gorm.DB) error {
	return db.Delete(m).Error
}

// Ls returns a list of directories or albums in an account.
func (m *Account) Ls() (result []string, err error) {
	if m.AccType == string(service.TypeWebDAV) {
		c := webdav.Connect(m.AccURL, m.AccUser, m.AccPass)
		result, err = c.Directories("/", true)
	}

	sort.Strings(result)

	return result, err
}
