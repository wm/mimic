package application

//App represents an application that is used locally and remotely that uses a database
type App interface {
	String() string
	DumpDb() Database
	RestoreDb() Database
	Path() string
	RestoreFromDump(string) error
	RunAfterRakeTasks() error
}

//Database represents the database used by an Applicaiton
type Database interface {
	Name() string
	Host() string
	Port() string
	User() string
	Dump(string, string) error
	Drop() error
	Create() error
	Restore(string) error
}
