package api39

import (
	"os"
	"os/exec"
)

func RebuildWithHugo(hugo, path string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	if err = os.Chdir(path); err != nil {
		return err
	}
	defer os.Chdir(cwd)

	cmd := exec.Command(hugo)

	return cmd.Run()
}

func RebuildWithMake(path string, targets ...string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	if err = os.Chdir(path); err != nil {
		return err
	}
	defer os.Chdir(cwd)

	mk, err := exec.LookPath("make")
	if err != nil {
		return err
	}

	cmd := exec.Command(mk, targets...)
	return cmd.Run()
}
