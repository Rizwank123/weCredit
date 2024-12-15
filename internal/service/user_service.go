package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/weCredit/internal/domain"
	"github.com/weCredit/internal/pkg/config"
	"github.com/weCredit/internal/pkg/security"
	"github.com/weCredit/internal/pkg/util"
)

type UserService struct {
	au  util.AppUtil
	cfg config.WeCreditConfig
	lcr domain.LoginCodeRepository
	scm security.Manager
	tr  domain.Transactioner
	usr domain.UserRepository
}

func NewUserService(au util.AppUtil, cfg config.WeCreditConfig, lcr domain.LoginCodeRepository, scm security.Manager, tr domain.Transactioner, usr domain.UserRepository) domain.UserService {
	return &UserService{
		au:  au,
		cfg: cfg,
		lcr: lcr,
		scm: scm,
		tr:  tr,
		usr: usr,
	}
}

// FindByUserName implements domain.UserService.
func (s *UserService) FindByUserName(username string) (result domain.User, err error) {
	return s.usr.FindByUserName(context.Background(), username)
}

// Login implements domain.UserService.
func (s *UserService) Login(in domain.LoginInput) (result domain.LoginOutput, err error) {
	usr, err := s.usr.FindByUserName(context.Background(), in.UserName)
	ctx := context.Background()
	ctx, err = s.tr.Begin(ctx)
	if err != nil {
		return result, err
	}
	defer func() {
		s.tr.Rollback(ctx, err)
	}()

	if err != nil {
		return result, errors.New("use not found please register first ")
	}
	loginCode, err := s.lcr.FindByUsername(context.TODO(), usr.UserName)
	if err != nil {
		return result, errors.New("login code not found init login first ")
	}
	if loginCode.Code != in.Otp {
		return result, errors.New("invalid otp")
	}
	// Check if the OTP has expired
	if time.Now().After(loginCode.ExpiryTime) {
		return result, errors.New("otp has expired")
	}
	// generate token
	ti := security.TokenMetadata{
		UserID: usr.ID.String(),
		Role:   usr.Role,
	}
	token, err := s.scm.GenerateAuthToken(ti)
	if err != nil {
		return result, err

	}
	err = s.lcr.DeleteByUsername(ctx, in.UserName)
	if err != nil {
		log.Println("Failed to delete login code:", err)
	}
	err = s.tr.Commit(ctx)
	if err != nil {
		return result, err
	}

	return domain.LoginOutput{
		Token:     token,
		ExpiresIn: int64(s.cfg.AuthExpiryPeriod),
	}, nil

}

// RegisterUser implements domain.UserService.
func (s *UserService) RegisterUser(in domain.RegisterUserInput) (result domain.User, err error) {
	result = domain.User{
		UserName: in.UserName,
		FullName: in.FullName,
		Role:     string(in.Role),
	}
	err = s.usr.CreateUser(context.Background(), &result)
	if err != nil {
		return result, err
	}
	return result, err
}

func (s *UserService) InitLogin(in domain.InitLoginInput) (err error) {
	usr, err := s.usr.FindByUserName(context.Background(), in.UserName)
	if err != nil {
		return err
	}
	if usr.ID.IsNil() {
		return errors.New("user not exists")
	}
	loginCode, err := s.lcr.FindByUsername(context.TODO(), usr.UserName)
	if err != nil && !errors.Is(err, domain.DataNotFoundError{}) {
		return err

	}
	var otp string
	if loginCode.Code != "" {
		otp = s.au.GenerateOTP(6)
		id := loginCode.ID
		loginCode := domain.LoginCode{
			Code:       otp,
			Status:     domain.LoginCodeStatusPENDING,
			Username:   usr.UserName,
			ExpiryTime: time.Now().Add(5 * time.Minute),
		}
		err = s.lcr.Update(context.TODO(), id, &loginCode)
		if err != nil {
			return err

		}
	} else {
		otp = s.au.GenerateOTP(6)
		loginCode := domain.LoginCode{
			Code:       otp,
			Status:     domain.LoginCodeStatusPENDING,
			Username:   usr.UserName,
			ExpiryTime: time.Now().Add(5 * time.Minute),
		}
		err = s.lcr.Create(context.TODO(), &loginCode)

	}
	if err != nil {
		return err
	}
	sms := domain.OtpMessage{
		To:  in.UserName,
		Otp: otp,
	}
	err = s.au.SendOtp(s.cfg, sms)
	return err

}

func (s *UserService) FindByID(id uuid.UUID) (result domain.User, err error) {
	return s.usr.FindByID(context.Background(), id)
}
