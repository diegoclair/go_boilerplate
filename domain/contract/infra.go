package contract

type Config interface {
	// GetMysqlDSN returns the mysql data source name
	GetMysqlDSN() string
}
