package server

import "io"

// For each client that connects to the server, a new FTPDriver is required.
// Create an implementation if this interface and provide it to FTPServer.
type DriverFactory interface {
	NewDriver() (Driver, error)
}

// You will create an implementation of this interface that speaks to your
// chosen persistence layer. graval will create a new instance of your
// driver for each client that connects and delegate to it as required.
type Driver interface {
	// Init init
	Init(*Conn)

	// params  - a file path
	// returns - a time indicating when the requested path was last modified
	//         - an error if the file doesn't exist or the user lacks
	//           permissions
	Stat(string) (FileInfo, error)

	// params  - path
	// returns - true if the current user is permitted to change to the
	//           requested path
	ChangeDir(string) error

	// params  - path, function on file or subdir found
	// returns - error
	//           path
	ListDir(string, func(FileInfo) error) error

	// params  - path
	// returns - true if the directory was deleted
	DeleteDir(string) error

	// params  - path
	// returns - true if the file was deleted
	DeleteFile(string) error

	// params  - from_path, to_path
	// returns - true if the file was renamed
	Rename(string, string) error

	// params  - path
	// returns - true if the new directory was created
	MakeDir(string) error

	// params  - path
	// returns - a string containing the file data to send to the client
	GetFile(string, int64) (int64, io.ReadCloser, error)

	// params  - desination path, an io.Reader containing the file data
	// returns - true if the data was successfully persisted
	PutFile(string, io.Reader, bool) (int64, error)
}
