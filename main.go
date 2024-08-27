package main

import (
	"fmt"
	"os"
)

const NIX_PATH string = "/etc/nixos/configuration.nix"

func install(packages []string, target string) {
	if packages[0] == "install" {
		packages = packages[1:]
	}

	pkgs := Pkgs{}.FromFile(target, "environment.systemPackages")
	for _, new_pkg := range packages {
		pkgs.pkgs = append(pkgs.pkgs, new_pkg)
	}

	pkgs.ToFile(target, "environment.systemPackages")
}

func remove(packages []string, target string) {
	if packages[0] == "remove" {
		packages = packages[1:]
	}

	pkgs := Pkgs{}.FromFile(target, "environment.systemPackages")

	new_pkgs := make([]string, 0)

	for _, pkg := range pkgs.pkgs {
		if !lstContains(packages, pkg) {
			new_pkgs = append(new_pkgs, pkg)
		}
	}

	pkgs.pkgs = new_pkgs
	pkgs.ToFile(target, "environment.systemPackages")
}

func list(target string) {
	pkgs := Pkgs{}.FromFile(target, "environment.systemPackages")

	fmt.Println("Installed Packages\n==================")
	for i, pkg := range pkgs.pkgs {
		fmt.Printf("%d)\t%s\n", i+1, pkg)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error, wrong arguments")
		os.Exit(1)
	}

	args := os.Args[1:]

	target := NIX_PATH
	if lstContains(args, "-t") {
		idx := lstIdxOf(args, "-t") + 1
		if idx >= len(args) {
			fmt.Println("Error, `-t` flag included but no target specified")
			os.Exit(1)
		}
		target = args[idx]
		args[idx] = ""
		args[idx-1] = ""
	}

	new_args := make([]string, 0)
	for _, i := range args {
		if i != "" {
			new_args = append(new_args, i)
		}
	}
	args = new_args

	if args[0] == "install" {
		fmt.Println("installing ", args[1:])
		install(args, target)
	} else if args[0] == "remove" {
		fmt.Println("removing ", args[1:])
		remove(args, target)
	} else if args[0] == "list" {
		list(target)
	}

	fmt.Println(args)

}
