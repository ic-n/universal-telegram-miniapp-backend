package records

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ic-n/universal-telegram-miniapp-backend/server/pkg/miniapp"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Chat struct {
	PK        uint      `json:"-" gorm:"column:pk;primaryKey"`
	ChatID    int       `json:"ChatID" gorm:"column:chat_id;unique"`
	Type      string    `json:"Type"`
	Label     string    `json:"Label"`
	Username  string    `json:"Username,omitempty"`
	PhotoURL  string    `json:"PhotoURL,omitempty"`
	Admin     bool      `json:"-"`
	CreatedAt time.Time `json:"CreatedAt"`
}

type Record struct {
	PK         uint           `json:"PK" gorm:"column:pk;primaryKey"`
	ChatPK     uint           `json:"-" gorm:"column:chat_pk"`
	Chat       Chat           `json:"Chat" gorm:"foreignKey:ChatPK;references:PK"`
	Collection string         `json:"Collection"`
	Value      string         `json:"Value"`
	Metadata   datatypes.JSON `json:"Metadata"`
	CreatedAt  time.Time      `json:"CreatedAt"`
}

type Database struct {
	*gorm.DB
}

func New(db *gorm.DB) (*Database, error) {
	d := Database{db}

	err := db.AutoMigrate(&Chat{}, &Record{})
	if err != nil {
		return nil, err
	}

	return &d, nil
}

type Request struct {
	Secret string `json:"Secret"`
	Limit  int    `json:"Limit"`
	Record Record `json:"Record"`
}

func (db *Database) ForgetChat(chatID int64) error {
	var chat Chat
	err := db.Model(&chat).
		Where("chat_id = ?", chatID).
		First(&chat).
		Error
	if err != nil {
		return err
	}

	err = db.
		Model(&Record{}).
		Where("chat_pk = ?", chat.PK).
		Unscoped().
		Delete(&Record{}).
		Error
	if err != nil {
		return err
	}

	err = db.
		Unscoped().
		Delete(chat).
		Error
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) HandlerAddRecord(ctx *gin.Context) {
	var req Request
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chat, err := db.auth(&req)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	rec, err := db.AddRecord(chat, req.Record)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, rec)
}

func (db *Database) AddRecord(chat Chat, rec Record) (Record, error) {
	rec.Chat = chat
	rec.ChatPK = chat.PK

	if err := db.Create(&rec).Error; err != nil {
		return Record{}, nil
	}

	return rec, nil
}

func (db *Database) HandlerRecords(ctx *gin.Context) {
	var req Request
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chat, err := db.auth(&req)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	records, err := db.Records(chat.PK, req.Record.Collection, req.Limit)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, records)
}

func (db *Database) Records(chatPK uint, collection string, limit int) ([]Record, error) {
	var records []Record
	err := db.
		Preload("Chat").
		Where("chat_pk = ?", chatPK).
		Where("collection = ?", collection).
		Order("created_at DESC").
		Limit(limit).
		Find(&records).
		Error
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (db *Database) auth(req *Request) (Chat, error) {
	if err := miniapp.Verify(req.Secret); err != nil {
		return Chat{}, err
	}

	id, err := miniapp.Parse(req.Secret)
	if err != nil {
		return Chat{}, err
	}

	photo := req.Record.Chat.PhotoURL
	if id.User.ID != 0 {
		label := id.User.FirstName
		if id.User.LastName != "" {
			label += " " + id.User.LastName
		}
		req.Record.Chat = Chat{
			ChatID:   int(id.User.ID),
			Type:     "user",
			Label:    label,
			Username: id.User.Username,
			PhotoURL: id.User.PhotoURL,
		}
	} else {
		req.Record.Chat = Chat{
			ChatID:   int(id.Chat.ID),
			Type:     string(id.Chat.Type),
			Label:    id.Chat.Title,
			Username: id.Chat.Username,
			PhotoURL: id.Chat.PhotoURL,
		}
	}
	if photo != "" {
		req.Record.Chat.PhotoURL = photo
	}

	var chat Chat
	err = db.Model(&chat).
		Where("chat_id = ?", req.Record.Chat.ChatID).
		First(&chat).
		Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return Chat{}, err
		}

		if err := db.CreateChat(req.Record.Chat); err != nil {
			return Chat{}, err
		}
	}

	return chat, nil
}

func (db *Database) CreateChat(chat Chat) error {
	err := db.Model(&chat).
		Create(&chat).
		Error
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) Chat(chatID int64) (Chat, error) {
	var chat Chat
	err := db.Model(&chat).
		Where("chat_id = ?", chatID).
		First(&chat).
		Error
	if err != nil {
		return Chat{}, err
	}
	return chat, nil
}
