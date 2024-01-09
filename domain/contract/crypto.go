package contract

type Crypto interface {
	HashPassword(password string) (string, error)
	CheckPassword(password, hashedPassword string) error
}
