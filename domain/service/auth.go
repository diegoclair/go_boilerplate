package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"

	"github.com/diegoclair/go-boilerplate/domain/entity"
	"github.com/diegoclair/go_utils-lib/v2/resterrors"
)

const (
	wrongLogin string = "Document or password are wrong"
)

type authService struct {
	svc *Service
}

func newAuthService(svc *Service) AuthService {
	return &authService{
		svc: svc,
	}
}

func (s *authService) Login(ctx context.Context, cpf, secret string) (account entity.Account, err error) {

	ctx, log := s.svc.log.NewSessionLogger(ctx)
	log.Info("Process Started")
	defer log.Info("Process Finished")

	account, err = s.svc.dm.Account().GetAccountByDocument(ctx, cpf)
	if err != nil {
		log.Error(err)
		return account, resterrors.NewUnauthorizedError(wrongLogin)
	}

	log.Info("account_id: ", account.ID)
	log.Info("account_uuid: ", account.UUID)
	log.Info("name: ", account.Name)

	hasher := md5.New()
	hasher.Write([]byte(secret))
	pass := hex.EncodeToString(hasher.Sum(nil))

	if pass != account.Secret {
		log.Error("wrong password")
		return account, resterrors.NewUnauthorizedError(wrongLogin)
	}

	return account, nil
}

func (s *authService) CreateSession(ctx context.Context, session entity.Session) (err error) {

	ctx, log := s.svc.log.NewSessionLogger(ctx)
	log.Info("Process Started")
	defer log.Info("Process Finished")

	err = s.svc.dm.Auth().CreateSession(ctx, session)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (s *authService) GetSessionByUUID(ctx context.Context, sessionUUID string) (session entity.Session, err error) {

	ctx, log := s.svc.log.NewSessionLogger(ctx)
	log.Info("Process Started")
	defer log.Info("Process Finished")

	session, err = s.svc.dm.Auth().GetSessionByUUID(ctx, sessionUUID)
	if err != nil {
		log.Error(err)
		return session, err
	}

	return session, nil
}
