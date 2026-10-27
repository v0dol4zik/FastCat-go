package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"time"
)

func hasFastfetch() bool {
	_, err := exec.LookPath("fastfetch")
	return err == nil
}

func homeConfig() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fatal("cannot determine home directory: " + err.Error())
	}
	return filepath.Join(home, ".config", "fastfetch")
}

func ensureHomeConfig() {
	if err := os.MkdirAll(homeConfig(), 0o755); err != nil {
		fatal(err.Error())
	}
}

// runFastfetch runs fastfetch interactively (wired to the real terminal).
func runFastfetch(args ...string) {
	cmd := exec.Command("fastfetch", args...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%sfastfetch exited with an error: %v%s\n", cRed, err, cReset)
	}
}

// findRepoRoot walks up from the cwd and the executable location,
// looking for the theme folders.
func findRepoRoot() string {
	var starts []string
	if wd, err := os.Getwd(); err == nil {
		starts = append(starts, wd)
	}
	if exe, err := os.Executable(); err == nil {
		starts = append(starts, filepath.Dir(exe))
	}
	for _, start := range starts {
		for d := start; ; d = filepath.Dir(d) {
			if isRepoRoot(d) {
				return d
			}
			if parent := filepath.Dir(d); parent == d {
				break
			}
		}
	}
	return ""
}

func isRepoRoot(d string) bool {
	for _, sub := range []string{"Small-Themes", "Large-Themes", "Visuals-Themes"} {
		if fi, err := os.Stat(filepath.Join(d, sub)); err == nil && fi.IsDir() {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------------------------
// File copying

func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if entry.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		if entry.Type() != 0 { // symlinks etc.
			if entry.Type()&os.ModeSymlink != 0 {
				link, err := os.Readlink(path)
				if err != nil {
					return err
				}
				os.Remove(target)
				return os.Symlink(link, target)
			}
			return nil
		}
		return copyFile(path, target)
	})
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	info, err := in.Stat()
	if err != nil {
		return err
	}
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode().Perm())
	if err != nil {
		return err
	}
	if _, err = io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}

// ---------------------------------------------------------------------------
// Backup (replaces backup.sh)

func backup() {
	src := homeConfig()
	fmt.Println(cDRed + "Backing Up..." + cCyan)
	if fi, err := os.Stat(src); err != nil || !fi.IsDir() {
		fmt.Println(cRed + "Nothing to back up: no fastfetch config found." + cReset)
		return
	}
	dst := "Backup-" + time.Now().Format("2006-01-02-15:04:05")
	if err := os.MkdirAll(dst, 0o755); err != nil {
		fatal(err.Error())
	}
	if err := copyDir(src, filepath.Join(dst, "fastfetch")); err != nil {
		fatal("backup failed: " + err.Error())
	}
	wd, _ := os.Getwd()
	fmt.Println(cRed + " Backed Up!" + cReset)
	fmt.Println("Saved to: " + filepath.Join(wd, dst))
}
