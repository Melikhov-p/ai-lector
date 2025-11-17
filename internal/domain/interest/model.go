package interest

import "strings"

// Interest - доменная модель интереса пользователя

type Interest struct {
	id    int64  // ID - идентификатор
	name  string // Name - название интереса
	emoji string // Emoji - смайлики интереса
}

// NewInterest создает новый интерес (без ID - он будет добавлен репозиторием)
func NewInterest(n string, emoji string) *Interest {
	return &Interest{
		name:  strings.TrimSpace(n),
		emoji: strings.TrimSpace(emoji),
	}
}

func NewInterestFromDB(id int64, name string, emoji string) *Interest {
	return &Interest{
		id:    id,
		name:  name,
		emoji: emoji,
	}
}

func (i *Interest) ID() int64 {
	return i.id
}

func (i *Interest) Name() string {
	return i.name
}

func (i *Interest) Emoji() string {
	return i.emoji
}

func (i *Interest) SetID(id int64) {
	i.id = id
}
