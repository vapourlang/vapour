package r

import (
	"encoding/json"
	"os/exec"

	"github.com/devOpifex/vapour/cache"
)

type Package struct {
	Name      string   `json:"name"`
	Functions []string `json:"functions"`
}

const BASEPACKAGES string = "BASEPACKAGES"

func Callr(cmd string) ([]byte, error) {
	out, err := exec.Command(
		"R",
		"-s",
		"-e",
		cmd,
	).Output()

	return out, err
}

func ListBaseFunctions() ([]Package, error) {
	c, ok := cache.Get(BASEPACKAGES)

	if ok {
		return c.([]Package), nil
	}

	var packages []Package

	output, err := Callr(
		`base_packages = getOption('defaultPackages')
		base_packages <- c(base_packages, "base")
		pkgs <- lapply(base_packages, function (pkg){
      fns <- ls(paste0('package:', pkg))
			fns <- fns[!grepl('<-', fns)]
			fns <- paste0(fns, collapse = '","')
			fns <- paste0('"functions":["', fns, '"]')
			pkg <- paste0('"name":"', pkg, '"')
			paste0('{', pkg, ',', fns, '}')
		})

		json <- paste0(pkgs, collapse = ",")
		cat(paste0("[", json, "]"))`,
	)

	if err != nil {
		return packages, err
	}

	err = json.Unmarshal(output, &packages)

	if err != nil {
		return packages, err
	}

	cache.Set(BASEPACKAGES, packages)

	return packages, err
}
