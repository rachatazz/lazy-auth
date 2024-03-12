package service

import (
	"errors"
	"fmt"

	"lazy-auth/app/errs"
	"lazy-auth/app/model"
	"lazy-auth/app/repository"
	"lazy-auth/app/zlog"
	"lazy-auth/common"

	"gorm.io/gorm"
)

type userService struct {
	userRepository repository.UserRepository
	roleRepository repository.RoleRepository
}

func NewUserService(
	userRepository repository.UserRepository,
	roleRepository repository.RoleRepository,
) UserService {
	return userService{
		userRepository: userRepository,
		roleRepository: roleRepository,
	}
}

func (s userService) CreateUser(
	userReq model.CreateUserRequest,
	roleName string,
) (*model.UserResponse, error) {
	role, err := s.roleRepository.GetByName(roleName)
	if err != nil {
		return nil, errs.NewNotFoundError(fmt.Sprintf("Role %s not found", roleName))
	}

	passwordHash, _ := common.HashPassword(userReq.Password)
	user := repository.User{
		RoleID:       role.ID,
		Email:        userReq.Email,
		Username:     userReq.Username,
		PasswordHash: passwordHash,
		DisplayName:  userReq.DisplayName,
		FirstName:    userReq.FirstName,
		LastName:     userReq.LastName,
	}

	err = s.userRepository.Create(&user)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, errs.NewUnprocessableEntity("Username or email duplicated")
		}
		zlog.Error(err)
		return nil, errs.NewUnexpectedError()
	}

	userResponse := model.UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		VerifyFlag:  user.VerifyFlag,
	}

	return &userResponse, nil
}

func (s userService) GetUsers(query model.QueryUser) (*model.UserPageResponse, error) {
	users, total, err := s.userRepository.GetMany(query)
	if err != nil {
		zlog.Error(err)
		return nil, errs.NewUnexpectedError()
	}

	rolesResponse := common.Map(users, func(user repository.User) model.UserResponse {
		return model.UserResponse{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			DisplayName: user.DisplayName,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			VerifyFlag:  user.VerifyFlag,
		}
	})
	meta := common.BuildMetaPagination(&total, query.Limit, query.Offset)

	return &model.UserPageResponse{Meta: meta, Data: rolesResponse}, nil
}

func (s userService) GetUserById(id string) (*model.UserResponse, error) {
	user, err := s.userRepository.GetById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFoundError("user not found")
		}
		zlog.Error(err)
		return nil, errs.NewUnexpectedError()
	}

	userResponse := model.UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		VerifyFlag:  user.VerifyFlag,
	}
	return &userResponse, nil
}

func (s userService) UpdateUserById(
	userId string,
	userReq model.UpdateUserRequest,
) (*model.UserResponse, error) {
	user, err := s.userRepository.GetById(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFoundError("user not found")
		}
		zlog.Error(err)
		return nil, errs.NewUnexpectedError()
	}

	if userReq.DisplayName != nil {
		user.DisplayName = *userReq.DisplayName
	}

	if userReq.FirstName != nil {
		user.FirstName = *userReq.FirstName
	}

	if userReq.LastName != nil {
		user.LastName = *userReq.LastName
	}

	err = s.userRepository.Update(user)
	if err != nil {
		zlog.Error(err)
		return nil, errs.NewUnexpectedError()
	}

	userResponse := model.UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		VerifyFlag:  user.VerifyFlag,
	}
	return &userResponse, nil
}

func (s userService) ChangePassword(
	userId string,
	changePassReq model.ChangePasswordRequest,
) (*model.UserResponse, error) {
	user, err := s.userRepository.GetById(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFoundError("user not found")
		}
		zlog.Error(err)
		return nil, errs.NewUnexpectedError()
	}

	ok := common.CheckPasswordHash(changePassReq.OldPassword, user.PasswordHash)
	if !ok {
		return nil, errs.NewUnauthorizedError("old password is incorrect")
	}

	user.PasswordHash, _ = common.HashPassword(changePassReq.NewPassword)
	err = s.userRepository.Update(user)
	if err != nil {
		zlog.Error(err)
		return nil, errs.NewUnexpectedError()
	}

	userResponse := model.UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		VerifyFlag:  user.VerifyFlag,
	}

	return &userResponse, nil
}
