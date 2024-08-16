package r

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/devOpifex/vapour/cache"
)

type Package struct {
	Name      string   `json:"name"`
	Functions []string `json:"functions"`
}

const (
	BASEPACKAGES = "BASEPACKAGES"
	FUNCTION     = "FUNCTION"
)

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

func PackageHasFunction(pkg, fn string) (bool, error) {
	key := pkg + fn
	c, ok := cache.Get(key)

	if ok {
		return c.(bool), nil
	}

	output, err := Callr(
		fmt.Sprintf("res <- tryCatch(%v::%s);cat(inherits(res, 'error'))", pkg, fn),
	)

	if err != nil {
		return false, err
	}

	ok = string(output) == "FALSE"

	cache.Set(key, ok)

	return ok, err
}

func PackageIsInstalled(pkg string) (bool, error) {
	key := "package::" + pkg
	c, ok := cache.Get(key)

	if ok {
		return c.(bool), nil
	}

	output, err := Callr(
		fmt.Sprintf("res <- requireNamespace('%v');cat(res)", pkg),
	)

	if err != nil {
		return false, err
	}

	ok = string(output) == "TRUE"

	cache.Set(key, ok)

	return ok, err
}
