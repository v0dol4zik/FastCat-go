package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const fontURL = "https://github.com/ryanoasis/nerd-fonts/releases/download/v3.2.1/CascadiaCode.zip"

// installIcons downloads and installs the CascadiaCode Nerd Font
// (replaces install-icons.sh; unzip is no longer needed — extraction is pure Go).
func installIcons() {
	const iconsBanner = `
    ______           __  ______      __
   / ____/___ ______/ /_/ ____/___ _/ /_
  / /_  / __ X / ___/ __/ /   / __ X / __/
 / __/ / /_/ (__  ) /_/ /___/ /_/ / /_
/_/    Y__,_/____/Y__/Y____/Y__,_/Y__/
 FastFetch Theme Pack`
	banner := strings.ReplaceAll(iconsBanner, "X", "`")
	banner = strings.ReplaceAll(banner, "Y", "\\")
	fmt.Println(cCyan + banner + "\n" + cReset)

	fmt.Printf("Do you want to install CascadiaCode Nerd Font? (Y/N): ")
	choice := readLine()
	if !strings.EqualFold(choice, "y") {
		fmt.Println(cRed + "Installation canceled by user." + cReset)
		return
	}

	var fontDir string
	if runtime.GOOS == "darwin" {
		home, _ := os.UserHomeDir()
		fontDir = filepath.Join(home, "Library", "Fonts")
	} else {
		home, _ := os.UserHomeDir()
		fontDir = filepath.Join(home, ".fonts")
	}
	fmt.Printf("%sEnsuring font directory exists: %s...%s\n", cGreen, fontDir, cReset)
	if err := os.MkdirAll(fontDir, 0o755); err != nil {
		fatal(err.Error())
	}

	fmt.Printf("%sDownloading CascadiaCode Nerd Font...%s\n", cBlue, cReset)
	client := &http.Client{Timeout: 10 * time.Minute}
	resp, err := client.Get(fontURL)
	if err != nil {
		fatal("download failed: " + err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fatal(fmt.Sprintf("download failed: HTTP %d", resp.StatusCode))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fatal("download failed: " + err.Error())
	}

	fmt.Printf("%sExtracting font files...%s\n", cBlue, cReset)
	zr, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		fatal("bad zip archive: " + err.Error())
	}
	for _, f := range zr.File {
		if f.FileInfo().IsDir() {
			continue
		}
		if err := extractZipFile(f, filepath.Join(fontDir, f.Name)); err != nil {
			fatal("extract failed: " + err.Error())
		}
	}

	fmt.Printf("%sRefreshing font cache...%s\n", cBlue, cReset)
	if path, err := exec.LookPath("fc-cache"); err == nil {
		cmd := exec.Command(path, "-fv")
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		cmd.Run()
	} else {
		fmt.Println(cYellow + "'fc-cache' not found — skipped (you may need to refresh the font cache manually)." + cReset)
	}

	fmt.Printf("\nInstallation complete! Please set 'CaskaydiaCove Nerd Font' in your terminal.\n")
}

func extractZipFile(f *zip.File, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, rc); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}
