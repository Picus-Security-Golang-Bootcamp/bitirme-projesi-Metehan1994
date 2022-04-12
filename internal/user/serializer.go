package user

import (
	"github.com/Metehan1994/final-project/internal/api"
	"github.com/Metehan1994/final-project/internal/models"
	"github.com/google/uuid"
)

func userToResponse(u *models.User) api.User {

	return api.User{
		ID:       u.ID.String(),
		Email:    &u.Email,
		Username: &u.Username,
		IsAdmin:  u.IsAdmin,
		Password: &u.Password,
	}
}

func signedUpUserToDBUser(s *api.SignUp) models.User {
	return models.User{
		ID:        uuid.New(),
		Email:     *s.Email,
		Username:  *s.Username,
		IsAdmin:   false,
		Password:  *s.Password,
		FirstName: *s.FirstName,
		LastName:  *s.LastName,
	}
}

// func responseToUser(u *api.User) *models.User {
// 	return &models.User{
// 		Username: *u.Username,
// 		Email:    *u.Email,
// 		IsAdmin:  u.IsAdmin,
// 		ID:       uuid.New(),
// 	}
// }

func RoleConvertToStringSlice(isAdmin bool) []string {
	var roles []string
	roles = append(roles, "customer")
	if isAdmin {
		roles = append(roles, "admin")
	}
	return roles
}
