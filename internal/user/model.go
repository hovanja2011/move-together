package user

type User struct {
	ID           string `json:"id" bson:"_id,omitempty"`
	Username     string `json:"username" bson:"_username"`
	PasswordHash string `json:"-" bson:"password"`
	Email        string `json:"email" bson:"email"`
}

// паттерн ДТО
// пользак: его данные -> сервер -> оборачиваем их в модель -> сторадж?? -> отправяем на БД
type CreateUserDTO struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password"`
	Email        string `json:"email"`
}
