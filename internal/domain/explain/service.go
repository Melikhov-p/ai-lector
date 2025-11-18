package explain

import (
	"fmt"
	"strings"

	"github.com/Melikhov-p/ai-lector/internal/domain/user"
)

// GetExplainRequestString собирает строку запроса пользователя на объяснение (предмет: string, тема: string, интересы: []string...)
func GetExplainRequestString(u *user.User, theme, subject string) string {
	var (
		req       string
		intersStr string
	)

	intersStr = strings.Join(u.InterestsNames(), ", ")

	req = fmt.Sprintf("предмет: %s, тема: %s, интересы: %s", subject, theme, intersStr)

	return req
}
