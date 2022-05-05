package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"

	"time"

	"github.com/diegoclair/go-boilerplate/domain/entity"
	"github.com/diegoclair/go-boilerplate/infra/auth"
	"github.com/diegoclair/go_utils-lib/v2/resterrors"
)

const (
	tokenExpirationTime        = 30 * time.Minute
	wrongLogin          string = "Document or password are wrong"
)

type authService struct {
	svc *Service
}

func newAuthService(svc *Service) AuthService {
	return &authService{
		svc: svc,
	}
}

func (s *authService) Login(ctx context.Context, cpf, secret string) (retVal entity.Authentication, err error) {

	ctx, log := s.svc.log.NewSessionLogger(ctx)
	log.Info("Login: Process Started")
	defer log.Info("Login: Process Finished")

	encryptedDocumentNumber, err := s.svc.cipher.Encrypt(cpf)
	if err != nil {
		log.Error("Login: ", err)
		return retVal, err
	}

	account, err := s.svc.dm.Account().GetAccountByDocument(ctx, encryptedDocumentNumber)
	if err != nil {
		log.Error("Login: ", err)
		return retVal, resterrors.NewUnauthorizedError(wrongLogin)
	}

	log.Info("Login: account_id: ", account.ID)
	log.Info("Login: account_uuid: ", account.UUID)
	log.Info("Login: name: ", account.Name)

	hasher := md5.New()
	hasher.Write([]byte(secret))
	pass := hex.EncodeToString(hasher.Sum(nil))

	if pass != account.Secret {
		log.Error("Login: wrong password")
		return retVal, resterrors.NewUnauthorizedError(wrongLogin)
	}

	issuedAt := time.Now()
	expiresAt := issuedAt.Add(tokenExpirationTime)

	claims := &entity.TokenData{}
	claims.AccountUUID = account.UUID
	claims.LoggedIn = true
	claims.IssuedAt = issuedAt.Unix()
	claims.ExpiresAt = expiresAt.Unix()
	claims.Issuer = "ST-go-boilerplate"

	newToken, err := auth.GenerateToken(s.svc.cfg.App.Auth, claims)
	if err != nil {
		log.Error("Login: error to generate token: ", err)
		return retVal, err
	}

	retVal.Token = newToken
	retVal.ServerTime = issuedAt.Unix()
	retVal.ValidTime = expiresAt.Unix()

	return retVal, nil
}
