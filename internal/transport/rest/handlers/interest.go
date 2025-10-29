package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Melikhov-p/ai-lector/internal/app"
	"github.com/Melikhov-p/ai-lector/internal/consts"
	"github.com/Melikhov-p/ai-lector/internal/domain/interest"
	"github.com/Melikhov-p/ai-lector/internal/transport/rest/dto"
	"github.com/Melikhov-p/ai-lector/internal/transport/rest/mapper"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type interestHandlers struct {
	log         *zap.Logger
	interestApp *app.InterestApp
}

func newInterestHandlers(log *zap.Logger, interestApp *app.InterestApp) *interestHandlers {
	return &interestHandlers{
		log:         log,
		interestApp: interestApp,
	}
}

// GetInterestsList получить все интересы
func (i *interestHandlers) GetInterestsList(w http.ResponseWriter, r *http.Request) {
	var (
		interests []*interest.Interest
		outDTO    dto.InterestListDTO
		err       error
	)

	interests, err = i.interestApp.GetInterests(r.Context())
	if err != nil {
		i.log.Error("Error getting interests", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, inter := range interests {
		outDTO.Interests = append(outDTO.Interests, mapper.FromInterestToInterestDTO(inter))
	}

	if err = json.NewEncoder(w).Encode(outDTO); err != nil {
		i.log.Error("Error encoding interests", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// GetInterest получаем конкретный интерес
func (i *interestHandlers) GetInterest(w http.ResponseWriter, r *http.Request) {
	var (
		inter         *interest.Interest
		interID       int
		interIDstring string
		err           error
	)

	interIDstring = chi.URLParam(r, consts.InterestIDURLParam.String())
	if interIDstring == "" {
		i.log.Error("interestID is required")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	interID, err = strconv.Atoi(interIDstring)
	if err != nil {
		i.log.Error("Error converting interestID to int", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	inter, err = i.interestApp.GetInterestByID(r.Context(), int64(interID))
	if err != nil {
		i.log.Error("Error getting interest", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(mapper.FromInterestToInterestDTO(inter)); err != nil {
		i.log.Error("Error encoding interest", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
