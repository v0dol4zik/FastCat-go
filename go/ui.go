package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// ANSI colors.
const (
	cReset   = "\033[0m"
	cRed     = "\033[1;31m"
	cGreen   = "\033[1;32m"
	cYellow  = "\033[1;33m"
	cBlue    = "\033[1;34m"
	cMagenta = "\033[1;35m"
	cCyan    = "\033[1;36m"
	cWhite   = "\033[1;37m"
	cDRed    = "\033[0;31m"
	cDGreen  = "\033[0;32m"
	cDCyan   = "\033[0;36m"
)

var stdin = bufio.NewReader(os.Stdin)

func clear() { fmt.Print("\033[2J\033[H") }

func readLine() string {
	line, err := stdin.ReadString('\n')
	if err != nil && strings.TrimSpace(line) == "" {
		// EOF (Ctrl+D / end of piped input) — exit instead of looping forever.
		fmt.Fprintln(os.Stderr, "\nEOF — exiting.")
		os.Exit(0)
	}
	return strings.TrimSpace(line)
}

func prompt() string {
	fmt.Printf("%sm3tozz%s@%sfastcat\n%s↳%s ", cYellow, cDRed, cBlue, cDRed, cCyan)
	return readLine()
}

func waitAnyKey() {
	fmt.Printf("%sPress any key to continue...%s\n", cDCyan, cReset)
	readKey()
}

// loader reproduces the original progress-bar animation.
func loader() {
	steps := []struct {
		bar  string
		pct  string
		wait time.Duration
	}{
		{"[#####               ]", "25%", 300 * time.Millisecond},
		{"[###############     ]", "75%", 200 * time.Millisecond},
		{"[####################]", "100%", 200 * time.Millisecond},
	}
	for _, s := range steps {
		fmt.Printf("\rFastCat - FastFetch Theme Pack! %s %s completed.", s.bar, s.pct)
		time.Sleep(s.wait)
	}
	fmt.Print("\r\033[K")
}

// atoi is a strict positive integer parser.
func atoi(s string) (int, bool) {
	if s == "" {
		return 0, false
	}
	n := 0
	for _, r := range s {
		if r < '0' || r > '9' {
			return 0, false
		}
		n = n*10 + int(r-'0')
		if n > 1_000_000 {
			return 0, false
		}
	}
	return n, true
}

func fatal(msg string) {
	fmt.Fprintln(os.Stderr, cRed+msg+cReset)
	os.Exit(1)
}

// ---------------------------------------------------------------------------
// Banners

func mainBanner() {
	fmt.Print(cCyan + `
⠀⠀⠀⢠⣾⣷⣦⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⣰⣿⣿⣿⣿⣷⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⢰⣿⣿⣿⣿⣿⣿⣷⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⢀⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⣦⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⣤⣀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣶⣤⣄⣀⣀⣤⣤⣶⣾⣿⣿⣿⡷
⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠁
⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠁⠀
⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣣⣿⣿⣿⣣⠏⠀⠀⠀
⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣠⠏⠀⠀⠀⠀    ` + cDRed + ` ______        _    _____      _   ` + cCyan + `
⣿⣿⣿⡇⠀⡾⠻⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠁⠀⠀⠀⠀⠀    ` + cDRed + ` |  ____|      | |  / ____|    | |  ` + cCyan + `
⣿⣿⣿⣧⡀⠁⣀⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡇⠀⠀⠀⠀⠀⠀    ` + cDRed + ` | |__ __ _ ___| |_| |     __ _| |_ ` + cCyan + `
⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡟⠉⢹⠉⠙⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀    ` + cDRed + ` |  __/ _` + "`" + ` / __| __| |    / _` + "`" + ` | __|` + cCyan + `
⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⣀⠀⣀⣼⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀   ` + cDRed + `  | | | (_| \__ \ |_| |___| (_| | |_ ` + cCyan + `
⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠋⠀⠀⠀⠀⠀⠀⠀⠀    ` + cDRed + ` |_|  \__,_|___/\__|\_____\__,_|\__|` + cCyan + `
⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠛⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀    ` + cDCyan + `     FastFetch Theme Pack.` + cCyan + `
⣿⣿⣿⣿⠿⣿⣿⣿⣿⣿⣿⣿⠿⠋⢃⠈⠢⡁⠒⠄⡀⠈⠁⠀⠀⠀⠀⠀⠀⠀
⣿⣿⠟⠁⠀⠀⠈⠉⠉⠁⠀⠀⠀⠀⠈⠆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠋⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠘⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀

`)
	fmt.Println(cBlue + "[01]" + cDGreen + "Small Themes " + cMagenta + "[02]" + cDGreen + "Large Themes " + cYellow + "[03]" + cDGreen + "Visuals Themes")
	fmt.Println(cRed + "[A]About [B]Backup [I]Icons [U]Update [x]Exit" + cReset)
}

func categoryBanner(cat *Category) {
	name := cat.Name
	fmt.Println(cDCyan)
	fmt.Println(cDRed + "______        _   _____       _   ")
	fmt.Println("\033[0;33m|  ___|      | | /  __ \\     | |  ")
	fmt.Println("\033[0;34m| |_ __ _ ___| |_| /  \\/ __ _| |_ ")
	fmt.Println("\033[0;35m|  _/ _` / __| __| |    / _` | __|" + cDRed + name)
	fmt.Println(cDCyan + "| || (_| \\__ \\ |_| \\__/\\ (_| | |_ ")
	fmt.Println(cDRed + "\\_| \\__,_|___/\\__|\\____/\\__,_|\\__|")
	fmt.Println(cReset)
	printGrid(cat.Themes, 4, 20)
	extra := ""
	if cat.Visual {
		extra = "\033[3m\033[92mThese themes require an image-supporting terminal emulator." + cReset + "\n"
	}
	fmt.Println(cYellow + "[P]Preview [G]Gallery [R]Random [S]Search [x]Exit [00]Menu [D]Default-Theme")
	fmt.Println(extra + cReset)
}

// printGrid renders the numbered theme grid, cycling index colors like the original.
var gridColors = []string{cBlue, cMagenta, cCyan, cDRed, cYellow, cDGreen}

func printGrid(themes []Theme, perRow, width int) {
	for i, t := range themes {
		col := gridColors[i%len(gridColors)]
		fmt.Printf("  %s[%02d]%s %-*s", col, i+1, cReset, width, t.Name)
		if (i+1)%perRow == 0 {
			fmt.Println()
		}
	}
	if len(themes)%perRow != 0 {
		fmt.Println()
	}
	fmt.Println()
}
