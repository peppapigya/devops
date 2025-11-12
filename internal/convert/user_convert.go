package convert

import (
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/dal/model"
)

func UserDtoToUserDO(userDto *dto.UserSaveRequest) *model.SystemUser {
	return &model.SystemUser{
		ID:        userDto.ID,
		Username:  userDto.Username,
		Nickname:  userDto.Nickname,
		Remark:    userDto.Remark,
		DeptID:    userDto.DeptID,
		PostIds:   userDto.PostIds,
		Email:     userDto.Email,
		Mobile:    userDto.Mobile,
		Sex:       userDto.Sex,
		Avatar:    userDto.Avatar,
		Status:    userDto.Status,
		LoginIP:   userDto.LoginIP,
		LoginDate: userDto.LoginDate,
	}
}
