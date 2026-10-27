package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// checkUpdate compares the local version with the upstream go/VERSION
// marker. In quiet mode it only prints a hint; in force mode it asks
// and syncs the repository via git.
func checkUpdate(root string, force bool) {
	remote := fetchRemoteVersion()
	if remote == "" {
		if force {
			fmt.Println(cYellow + "Could not check for updates (offline or upstream unreachable)." + cReset)
		}
		return
	}
	if remote == version {
		if force {
			fmt.Printf("%sFastCat is up to date (%s).%s\n", cGreen, version, cReset)
		}
		return
	}

	fmt.Printf("%snew version found: %s → %s%s\n", cYellow, version, remote, cReset)
	if !force {
		fmt.Println(cDCyan + "Sync later with: fastcat --update" + cReset)
		return
	}

	fmt.Printf("%sUpdate now? [y/N]:%s ", cYellow, cReset)
	if !strings.EqualFold(readLine(), "y") {
		fmt.Println(cDCyan + "Skipped. You can update later with: fastcat --update" + cReset)
		return
	}

	if fi, err := os.Stat(filepath.Join(root, ".git")); err != nil || !fi.IsDir() {
		fmt.Printf("%sNot a git repository — clone it again:%s\n%s\n", cRed, cReset,
			"git clone --depth 1 https://github.com/m3tozz/FastCat.git")
		return
	}
	fmt.Println("updating... syncing entire repository")
	branch := gitOut(root, "symbolic-ref", "--short", "HEAD")
	if branch == "" {
		branch = "main"
	}
	gitRun(root, "fetch", "origin")
	gitRun(root, "reset", "--hard", "origin/"+branch)
	fmt.Printf("%sUpdated! Rebuild the binary:%s go build -o fastcat .%s\n", cGreen, cCyan, cReset)
}

func fetchRemoteVersion() string {
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(versionURL)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ""
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<12))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(body))
}

func gitOut(dir string, args ...string) string {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func gitRun(dir string, args ...string) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%sgit %s failed: %v%s\n", cRed, strings.Join(args, " "), err, cReset)
	}
}
