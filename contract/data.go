package contract

// DataManager holds the methods that manipulates the main data.
type DataManager interface {
	MySQL() MySQLRepo
}
