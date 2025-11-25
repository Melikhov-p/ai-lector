package subscription

import "time"

// Subscription модель подписки
type Subscription struct {
	id       int64         // Внутренний ид
	name     string        // Название
	price    float64       // Цена (без скидок)
	duration time.Duration // Длительность
}

// NewSubscription модель подписки
func NewSubscription(name string, price float64, duration time.Duration) *Subscription {
	return &Subscription{
		name:     name,
		price:    price,
		duration: duration,
	}
}

// NewSubscriptionFromDB подписка из БД
func NewSubscriptionFromDB(id int64, name string, price float64, duration time.Duration) *Subscription {
	return &Subscription{
		id:       id,
		name:     name,
		price:    price,
		duration: duration,
	}
}

func (s *Subscription) ID() int64 {
	return s.id
}

func (s *Subscription) Name() string {
	return s.name
}

func (s *Subscription) Price() float64 {
	return s.price
}

func (s *Subscription) Duration() time.Duration {
	return s.duration
}

// UserSubscription связка подписки с пользователем
type UserSubscription struct {
	id        int64     // ИД
	userID    int64     // Пользователь
	subID     int64     // Подписка
	createdAt time.Time // Когда оформил
	endAt     time.Time // Когда закончится
	payed     bool      // Оплачена ли
}

// NewUserSubscription новая подписка пользователя
func NewUserSubscription(userID int64, sub *Subscription) *UserSubscription {
	return &UserSubscription{
		userID:    userID,
		subID:     sub.ID(),
		createdAt: time.Now(),
		endAt:     time.Now().Add(sub.duration),
		payed:     false,
	}
}

func (u *UserSubscription) ID() int64 {
	return u.id
}

func (u *UserSubscription) SubID() int64 {
	return u.subID
}

func (u *UserSubscription) CreatedAt() time.Time {
	return u.createdAt
}

func (u *UserSubscription) EndAt() time.Time {
	return u.endAt
}

func (u *UserSubscription) Payed() bool {
	return u.payed
}

func (u *UserSubscription) UserID() int64 {
	return u.userID
}
