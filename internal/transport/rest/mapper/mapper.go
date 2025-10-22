package mapper

import (
	"github.com/Melikhov-p/ai-lector/internal/domain/user"
	"github.com/Melikhov-p/ai-lector/internal/transport/rest/dto"
)

// FromUserToDTO - из доменной модели в транспортную пользователя
func FromUserToDTO(usr *user.User) dto.UserDTO {
	return dto.UserDTO{
		ID:        usr.ID(),
		Email:     usr.Email(),
		FirstName: usr.FirstName(),
		LastName:  usr.LastName(),
		Phone:     usr.Phone(),
	}
}

// FromSearchUserDTOtoUserFilter из траспортной модели поиска в доменный фильтр
func FromSearchUserDTOtoUserFilter(d dto.SearchUserDTO) user.UserFilter {
	return user.UserFilter{
		Email:     d.Email,
		Phone:     d.Phone,
		FirstName: d.FirstName,
		LastName:  d.LastName,
	}
}
