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
	log.Info("Login: Process Started")
	defer log.Info("Login: Process Finished")

	encryptedDocumentNumber, err := s.svc.cipher.Encrypt(cpf)
	if err != nil {
		log.Error("Login: ", err)
		return account, err
	}

	account, err = s.svc.dm.Account().GetAccountByDocument(ctx, encryptedDocumentNumber)
	if err != nil {
		log.Error("Login: ", err)
		return account, resterrors.NewUnauthorizedError(wrongLogin)
	}

	log.Info("Login: account_id: ", account.ID)
	log.Info("Login: account_uuid: ", account.UUID)
	log.Info("Login: name: ", account.Name)

	hasher := md5.New()
	hasher.Write([]byte(secret))
	pass := hex.EncodeToString(hasher.Sum(nil))

	if pass != account.Secret {
		log.Error("Login: wrong password")
		return account, resterrors.NewUnauthorizedError(wrongLogin)
	}

	return account, nil
}
