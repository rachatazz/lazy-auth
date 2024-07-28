package service

import (
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
	refreshTokenExpiresAt := common.AddTimeByDuration(s.configEnv.JwtRefreshTokenExpired)
	session := repository.Session{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    refreshTokenExpiresAt,
	}
	err = s.sessionRepository.Create(&session)
	if err != nil {
		return nil, errs.NewUnauthorizedError("username or password is incorrect")
	}

	tokenExpiresAt := common.AddTimeByDuration(s.configEnv.JwtTokenExpired)
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
	session.ExpiresAt = common.AddTimeByDuration(s.configEnv.JwtRefreshTokenExpired)
	s.sessionRepository.Update(session)

	tokenExpiresAt := common.AddTimeByDuration(s.configEnv.JwtTokenExpired)
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

func (s authService) ForgotPassword(model.ForgotPasswordRequest) error {
	return nil
}

func (s authService) ResetPassword(model.ResetPasswordRequest) (*model.UserResponse, error) {
	return nil, nil
}
