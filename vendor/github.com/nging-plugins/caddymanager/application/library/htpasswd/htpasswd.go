package htpasswd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/webx-top/com"
)

const (
	// PasswordSeparator separates passwords from hashes
	PasswordSeparator = ":"
	// LineSeparator separates password records
	LineSeparator = "\n"
)

// MaxHtpasswdFilesize if your htpassd file is larger than 8MB, then your are doing it wrong
const MaxHtpasswdFilesize = 8 * 1024 * 1024

// ParseHtpasswdFile parse htpasswd file
func ParseHtpasswdFile(file string) (users Accounts, err error) {
	var fi os.FileInfo
	fi, err = os.Stat(file)
	if err != nil {
		return
	}
	if fi.Size() > MaxHtpasswdFilesize {
		err = errors.New("this file is too large, use a database instead")
		return
	}

	users = Accounts{}
	var lineNumber int
	err = com.SeekFileLines(file, func(line string) error {
		lineNumber++
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			// skipping empty lines
			return nil
		}
		parts := strings.Split(line, PasswordSeparator)
		if len(parts) != 2 {
			err := errors.New(fmt.Sprintln("invalid line", lineNumber, "unexpected number of parts split by", PasswordSeparator, len(parts), "instead of 2 in\"", line, "\""))
			return err
		}
		for i, part := range parts {
			parts[i] = strings.TrimSpace(part)
		}
		_, alreadyExists := users[parts[0]]
		if alreadyExists {
			err := errors.New("invalid htpasswords file - user " + parts[0] + " was already defined")
			return err
		}
		users[parts[0]] = parts[1]
		return nil
	})
	return
}

// RemoveUser remove an existing user from a file, returns an error, if the user does not \
// exist in the file
func RemoveUser(file, user string) error {
	passwords, err := ParseHtpasswdFile(file)
	if err != nil {
		return err
	}
	_, ok := passwords[user]
	if !ok {
		return os.ErrNotExist
	}
	delete(passwords, user)
	return passwords.WriteToFile(file)
}

// SetPasswordHash directly set a hash for a user in a file
func SetPasswordHash(file, user, hash string) error {
	if len(hash) == 0 {
		return errors.New("you might want to rethink your hashing algorithm, it left you with an empty hash")
	}
	passwords, err := ParseHtpasswdFile(file)
	if err != nil {
		return err
	}
	passwords[user] = hash
	return passwords.WriteToFile(file)
}

// SetPassword set password for a user with a given hashing algorithm
func SetPassword(file, name, password string, algo Algorithm) error {
	passwords, err := ParseHtpasswdFile(file)
	if err != nil {
		return err
	}
	err = passwords.SetPassword(name, password, algo)
	if err != nil {
		return err
	}
	return passwords.WriteToFile(file)
}
