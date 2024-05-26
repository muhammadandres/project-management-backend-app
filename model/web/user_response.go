package web

import (
	"manajemen_tugas_master/model/domain"
)

func CreateResponseUser(userModel *domain.User) WebResponse {
	return WebResponse{
		Code:    200,
		Message: "Success",
		Data: domain.User{
			ID:    userModel.ID,
			Email: userModel.Email,
		},
	}
}
