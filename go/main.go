// FastCat — FastFetch Theme Pack (Go edition)
// Port of fastcat.sh / start.sh (Small, Large, Visuals themes).
package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	app     = "FastCat"
	version = "2.0.0-go"
	// Version marker for the Go edition (checked for updates).
	versionURL = "https://raw.githubusercontent.com/m3tozz/FastCat/main/go/VERSION"
	remoteURL  = "https://github.com/m3tozz/FastCat"
)

func main() {
	arg := ""
	if len(os.Args) > 1 {
		arg = strings.ToLower(strings.TrimSpace(os.Args[1]))
	}
	switch arg {
	case "--shell", "-s":
		runShell()
	case "--backup", "-b":
		backup()
	case "--icons", "-i":
		installIcons()
	case "--about", "-a":
		about()
	case "--version", "-v":
		fmt.Println(version)
	case "--update", "-u":
		root := findRepoRoot()
		if root == "" {
			fatal("theme folders not found — run FastCat from the repository (or its build dir)")
		}
		checkUpdate(root, true)
	case "--help", "-h", "help":
		help()
	default:
		help()
	}
}

func runShell() {
	if !hasFastfetch() {
		printFastfetchMissing()
		os.Exit(1)
	}
	root := findRepoRoot()
	if root == "" {
		fatal("theme folders not found — run FastCat from the repository (or its build dir)")
	}
	ensureHomeConfig()

	// Quick, non-blocking-ish update check (3s timeout, silent on failure).
	checkUpdate(root, false)

	cats, err := loadCategories(root)
	if err != nil {
		fatal(err.Error())
	}

	for {
		mainMenu(root, cats)
	}
}

func mainMenu(root string, cats []*Category) {
	clear()
	mainBanner()
	input := prompt()
	switch strings.ToLower(input) {
	case "1", "01":
		categoryMenu(cats[0])
	case "2", "02":
		categoryMenu(cats[1])
	case "3", "03":
		categoryMenu(cats[2])
	case "a":
		about()
		waitAnyKey()
	case "b":
		backup()
		waitAnyKey()
	case "i":
		installIcons()
		waitAnyKey()
	case "u":
		checkUpdate(root, true)
		waitAnyKey()
	case "x":
		clear()
		fmt.Println(cRed + "GoodBye" + cReset)
		os.Exit(0)
	default:
		fmt.Println(cDCyan + "Wrong transaction number!" + cReset)
	}
}

func help() {
	fmt.Printf(`Wrong usage, there are several arguments for fastcat

	--shell : run the FastCat interactive menu.
	--backup: back up your current fastfetch configuration.
	--icons : install CascadiaCode Nerd Font.
	--update: sync the FastCat repository via git.
	--version: show the version.
	--about: about the FastCat project.
	--help: show this page.
	----------------------------------
	ex usage: ./fastcat --shell
`)
}

func about() {
	fmt.Println(`
 _____         _    ____      _
|  ___|_ _ ___| |_ / ___|__ _| |_
| |_ / _` + "`" + ` / __| __| |   / _` + "`" + ` | __|FastFetch Theme Pack!
|  _| (_| \__ \ |_| |__| (_| | |_
|_|  \__,_|___/\__|\____/\__,_|\__|`)
	fmt.Println(cBlue + "###############################################################" + cReset)
	fmt.Println("    Github              : " + cRed + " github.com/m3tozz/FastCat" + cReset)
	fmt.Println("    Community Server    : " + cRed + " discord.com/invite/sQwYCZer95" + cReset)
	fmt.Println("    Version             : " + cRed + " " + version + " " + cReset)
	fmt.Println(cBlue + "###############################################################" + cReset)
}

func printFastfetchMissing() {
	clear()
	fmt.Println(cRed + "FastFetch Not Found!" + cYellow)
	fmt.Println("To use FastCat, you must first install FastFetch :)")
	fmt.Println("\n" + cRed + "FastFetch Installation:" + cReset)
	fmt.Println(cCyan + "Debian or Ubuntu using APT:" + cReset)
	fmt.Println("  " + cBlue + "→ sudo apt install fastfetch" + cReset)
	fmt.Println("\n" + cGreen + "Arch Linux using Pacman:" + cReset)
	fmt.Println("  " + cBlue + "→ sudo pacman -S fastfetch" + cReset)
	fmt.Println("\n" + cMagenta + "RHEL, Fedora, or CentOS using DNF:" + cReset)
	fmt.Println("  " + cBlue + "→ sudo dnf install fastfetch" + cReset)
}
