package service

import (
	"errors"
	"fmt"
	"time"

	"lazy-auth/app/errs"
	"lazy-auth/app/model"
	"lazy-auth/app/repository"
	"lazy-auth/app/zlog"
	"lazy-auth/common"
	"lazy-auth/config"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type authService struct {
	userRepository    repository.UserRepository
	roleRepository    repository.RoleRepository
	sessionRepository repository.SessionRepository
	configEnv         config.ConfigEnv
}

func NewAuthService(
	userRepository repository.UserRepository,
	roleRepository repository.RoleRepository,
	sessionRepository repository.SessionRepository,
	configEnv config.ConfigEnv,
) AuthService {
	return authService{
		userRepository:    userRepository,
		roleRepository:    roleRepository,
		sessionRepository: sessionRepository,
		configEnv:         configEnv,
	}
}

func (s authService) CreateRole(roleReq model.CreateRoleRequest) (*model.RoleResponse, error) {
	prepareCreateRole := repository.Role{
		Name:        roleReq.Name,
		Description: roleReq.Description,
	}

	role, err := s.roleRepository.Create(prepareCreateRole)
	if err != nil {
		zlog.Error(err)
		if err == gorm.ErrDuplicatedKey {
			return nil, errs.NewUnprocessableEntity("Role name duplicated")
		}
		zlog.Error(err)
		return nil, errs.NewUnexpectedError()
	}

	roleResponse := model.RoleResponse{
		Name:        role.Name,
		Description: role.Description,
	}

	return &roleResponse, nil
}

func (s authService) GetRoles(query model.QueryRole) (*model.RolePageResponse, error) {
	roles, total, err := s.roleRepository.GetAll(query)
	if err != nil {
		zlog.Error(err)
		return nil, errs.NewUnexpectedError()
	}

	rolesResponse := common.Map(roles, func(role repository.Role) model.RoleResponse {
		return model.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
		}
	})
	meta := common.BuildMetaPagination(&total, query.Limit, query.Offset)

	return &model.RolePageResponse{Meta: meta, Data: rolesResponse}, nil
}

func (s authService) Login(body model.LoginRequest) (*model.TokenResponse, error) {
	user, err := s.userRepository.GetByUsername(body.Username)
	if err != nil {
		return nil, errs.NewUnauthorizedError("username or password is incorrect")
	}

	ok := common.CheckPasswordHash(body.Password, user.PasswordHash)
	if !ok {
		return nil, errs.NewUnauthorizedError("username or password is incorrect")
	}

	user.LastAccessAt = time.Now()
	err = s.userRepository.Update(user)
	if err != nil {
		return nil, errs.NewUnexpectedError()
	}

	refreshToken := uuid.NewString()
	refreshTokenExpiresAt := common.AddTimeByDuration(s.configEnv.JwtRefreshTokenExpiresIn)
	session := repository.Session{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    refreshTokenExpiresAt,
	}
	err = s.sessionRepository.Create(&session)
	if err != nil {
		return nil, errs.NewUnauthorizedError("username or password is incorrect")
	}

	tokenExpiresAt := common.AddTimeByDuration(s.configEnv.JwtTokenExpiresIn)
	token := common.GenerateToken(
		user.ID,
		session.ID,
		s.configEnv.JwtTokenSecret,
		tokenExpiresAt,
	)

	refreshTokenAES, err := common.Encrypt(
		session.RefreshToken,
		s.configEnv.JwtRefreshTokenSecret,
	)
	if err != nil {
		zlog.Error(err)
		return nil, errs.NewUnexpectedError()
	}

	return &model.TokenResponse{
		AccessToken:           token,
		TokenType:             "Bearer",
		TokenExpiresAt:        tokenExpiresAt,
		RefreshToken:          refreshTokenAES,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
	}, nil
}

func (s authService) RefreshToken(
	refreshTokenC string,
) (*model.TokenResponse, error) {
	refreshToken, err := common.Decrypt(
		refreshTokenC,
		s.configEnv.JwtRefreshTokenSecret,
	)
	if err != nil {
		return nil, errs.NewUnauthorizedError("refresh token is invalid")
	}

	session, err := s.sessionRepository.GetByRefreshToken(refreshToken)
	if err != nil {
		return nil, errs.NewUnauthorizedError("refresh token is invalid")
	}

	session.RefreshToken = uuid.NewString()
	session.ExpiresAt = common.AddTimeByDuration(s.configEnv.JwtRefreshTokenExpiresIn)
	s.sessionRepository.Update(session)

	tokenExpiresAt := common.AddTimeByDuration(s.configEnv.JwtTokenExpiresIn)
	token := common.GenerateToken(
		session.UserID,
		session.ID,
		s.configEnv.JwtTokenSecret,
		tokenExpiresAt,
	)

	refreshTokenAES, err := common.Encrypt(
		session.RefreshToken,
		s.configEnv.JwtRefreshTokenSecret,
	)
	if err != nil {
		zlog.Error(err)
		return nil, errs.NewUnexpectedError()
	}

	return &model.TokenResponse{
		AccessToken:           token,
		TokenType:             "Bearer",
		TokenExpiresAt:        tokenExpiresAt,
		RefreshToken:          refreshTokenAES,
		RefreshTokenExpiresAt: session.ExpiresAt,
	}, nil
}

func (s authService) Logout(logoutReq model.LogoutRequest) error {
	refreshToken, err := common.Decrypt(
		logoutReq.RefreshToken,
		s.configEnv.JwtRefreshTokenSecret,
	)
	if err != nil {
		return errs.NewUnauthorizedError("refresh token is invalid")
	}

	session, err := s.sessionRepository.GetByRefreshToken(refreshToken)
	if err != nil {
		return errs.NewUnauthorizedError("refresh token is invalid")
	}

	if logoutReq.IsAll {
		err = s.sessionRepository.DeleteByUserId(session.UserID)
		if err != nil {
			zlog.Error(err)
			return errs.NewUnexpectedError()
		}
		return nil
	}

	err = s.sessionRepository.DeleteById(session.ID)
	if err != nil {
		zlog.Error(err)
		return errs.NewUnexpectedError()
	}
	return nil
}

func (s authService) ChangePassword(
	userId string,
	changePassReq model.ChangePasswordRequest,
) error {
	user, err := s.userRepository.GetById(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.NewNotFoundError("user not found")
		}
		zlog.Error(err)
		return errs.NewUnexpectedError()
	}

	ok := common.CheckPasswordHash(changePassReq.OldPassword, user.PasswordHash)
	if !ok {
		return errs.NewUnauthorizedError("old password is incorrect")
	}

	user.PasswordHash, _ = common.HashPassword(changePassReq.NewPassword)
	user.ChangePasswordAt = time.Now()
	err = s.userRepository.Update(user)
	if err != nil {
		zlog.Error(err)
		return errs.NewUnexpectedError()
	}

	return nil
}

func (s authService) ForgotPassword(forgotReq model.ForgotPasswordRequest) error {
	user, err := s.userRepository.GetByEmail(forgotReq.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		zlog.Error(err)
		return errs.NewUnexpectedError()
	}

	user.Ticket = fmt.Sprintf("reset_password:%s", uuid.NewString())
	user.TicketExpiresAt = common.AddTimeByDuration(s.configEnv.TicketExpiresIn)
	err = s.userRepository.Update(user)
	if err != nil {
		zlog.Error(err)
		return errs.NewUnexpectedError()
	}

	// TODO: Implement functionality to send tickets via email

	return nil
}

func (s authService) ResetPassword(resetReq model.ResetPasswordRequest) error {
	user, err := s.userRepository.GetByTicket(resetReq.Ticket)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.NewUnauthorizedError("ticket is invalid")
		}
		zlog.Error(err)
		return errs.NewUnexpectedError()
	}

	user.PasswordHash, _ = common.HashPassword(resetReq.Password)
	user.Ticket = ""
	user.TicketExpiresAt = time.Time{}
	user.ChangePasswordAt = time.Now()
	err = s.userRepository.Update(user)
	if err != nil {
		zlog.Error(err)
		return errs.NewUnexpectedError()
	}

	return nil
}
