package r

import "os/exec"

func Callr(cmd string) ([]byte, error) {
	out, err := exec.Command(
		"R",
		"-s",
		"-e",
		cmd,
	).Output()

	return out, err
}
