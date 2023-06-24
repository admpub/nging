package packer

import (
	"os/exec"
	"strings"
)

type Manager struct {
	Name       string
	InstallArg string
	UpdateArg  string
	RemoveArg  string
}

func Check(packageName string) bool {
	_, err := exec.LookPath(packageName)
	return err == nil
}

func Install(packageName string) error {
	mngr, err := Default()
	if err != nil {
		return err
	}

	c := mngr.Name + " " + mngr.InstallArg + " " + packageName
	err = Command(c)
	return err
}

func Remove(packageName string) error {
	mngr, err := Default()
	if err != nil {
		return err
	}

	c := mngr.Name + " " + mngr.RemoveArg + " " + packageName
	err = Command(c)
	return err
}

func Update() error {
	mngr, err := Default()
	if err != nil {
		return err
	}

	c := mngr.Name + " " + mngr.UpdateArg
	err = Command(c)
	return err
}

func Command(command string) error {
	args := strings.Fields(command)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = Stdout
	cmd.Stderr = Stderr
	err := cmd.Run()
	return err
}
