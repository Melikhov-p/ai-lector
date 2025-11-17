package mapper

import (
	"github.com/Melikhov-p/ai-lector/internal/domain/interest"
	"github.com/Melikhov-p/ai-lector/internal/domain/user"
	"github.com/Melikhov-p/ai-lector/internal/transport/rest/dto"
)

// FromUserToDTO - из доменной модели в транспортную пользователя
func FromUserToDTO(usr *user.User) dto.UserDTO {
	usrDTO := dto.UserDTO{
		ID:            usr.ID(),
		Email:         usr.Email(),
		FirstName:     usr.FirstName(),
		LastName:      usr.LastName(),
		TrialRequests: usr.TrialRequests(),
		Phone:         usr.Phone(),
	}

	for _, i := range usr.Interests() {
		usrDTO.Interests = append(usrDTO.Interests, FromInterestToInterestDTO(i))
	}

	return usrDTO
}

func FromInterestToInterestDTO(inter *interest.Interest) *dto.InterestDTO {
	interDTO := dto.InterestDTO{
		Id:    inter.ID(),
		Name:  inter.Name(),
		Emoji: inter.Emoji(),
	}

	return &interDTO
}

// // FromSearchUserDTOtoUserFilter из транспортной модели поиска в доменный фильтр
// func FromSearchUserDTOtoUserFilter(d dto.SearchUserDTO) user.UserFilter {
// 	return user.UserFilter{
// 		Email:     d.Email,
// 		Phone:     d.Phone,
// 		FirstName: d.FirstName,
// 		LastName:  d.LastName,
// 	}
// }
