package errcodes

import "github.com/diegoclair/apperr"

var (
	// Auth errors
	ErrInvalidCredentials  = apperr.Define(apperr.KindAuthentication, "AUTH_INVALID_CREDENTIALS", "document or password are wrong")
	ErrDeactivatedAccount  = apperr.Define(apperr.KindAuthentication, "AUTH_ACCOUNT_DEACTIVATED", "account is deactivated")
	ErrSessionNotFound     = apperr.Define(apperr.KindAuthentication, "AUTH_SESSION_NOT_FOUND", "session not found")
	ErrSessionBlocked      = apperr.Define(apperr.KindAuthentication, "AUTH_SESSION_BLOCKED", "session blocked")
	ErrSessionTokenMismatch = apperr.Define(apperr.KindAuthentication, "AUTH_SESSION_TOKEN_MISMATCH", "mismatched session token")
	ErrSessionExpired      = apperr.Define(apperr.KindAuthentication, "AUTH_SESSION_EXPIRED", "session has expired")

	// Account errors
	ErrCPFAlreadyInUse = apperr.Define(apperr.KindConflict, "ACCOUNT_CPF_EXISTS", "the CPF is already in use")

	// Transfer errors
	ErrInsufficientFunds       = apperr.Define(apperr.KindConflict, "TRANSFER_INSUFFICIENT_FUNDS", "your account doesn't have sufficient funds to do this operation")
	ErrSelfTransfer            = apperr.Define(apperr.KindValidation, "TRANSFER_SELF_TRANSFER", "you can't transfer to yourself")
	ErrInvalidDestinationAccount = apperr.Define(apperr.KindNotFound, "TRANSFER_DEST_ACCOUNT_NOT_FOUND", "invalid destination account")
)
