package models

// import "github.com/jinzhu/gorm"

type PushMessage struct {
	ID     int `gorm:"primary_key;autoIncrement"`
	From   string
	Target string
}

func (PushMessage) TableName() string {
	return "push_message"
}

func GetAllMessage() {}

func NewMessage(data map[string]interface{}) {
	pm := &PushMessage{
		From:   data["from"].(string),
		Target: data["target"].(string),
	}

	db.Create(pm)
}

func GetMessageByTarget(target string) *PushMessage {
	pm := &PushMessage{
		Target: target,
	}
	err := db.Where("target = ?", target).Find(&pm)
	if err != nil {
		return nil
	}
	return pm
}

func DeleteMessageByTarget() {}
