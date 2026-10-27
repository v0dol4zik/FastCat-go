package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Theme is a single theme directory (contains a "fastfetch" subfolder).
type Theme struct {
	Name string // directory name
	Dir  string // full path to the theme directory
}

// Category groups themes (Small-Themes, Large-Themes, Visuals-Themes).
type Category struct {
	Name   string // display name, e.g. "Small-Themes"
	Path   string // full path to the category directory
	Visual bool   // image-based themes (need a rendering protocol)
	Themes []Theme
}

func loadCategories(root string) ([]*Category, error) {
	specs := []struct {
		dir    string
		visual bool
	}{
		{"Small-Themes", false},
		{"Large-Themes", false},
		{"Visuals-Themes", true},
	}
	var cats []*Category
	for _, s := range specs {
		cat, err := loadCategory(filepath.Join(root, s.dir), s.visual)
		if err != nil {
			return nil, err
		}
		cats = append(cats, cat)
	}
	return cats, nil
}

// loadCategory scans a theme folder for subdirectories containing a
// "fastfetch" config dir — no hardcoded theme lists needed.
func loadCategory(path string, visual bool) (*Category, error) {
	fi, err := os.Stat(path)
	if err != nil || !fi.IsDir() {
		return nil, fmt.Errorf("theme folder %q not found", filepath.Base(path))
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	cat := &Category{Name: filepath.Base(path), Path: path, Visual: visual}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		dir := filepath.Join(path, e.Name())
		if st, err := os.Stat(filepath.Join(dir, "fastfetch")); err == nil && st.IsDir() {
			cat.Themes = append(cat.Themes, Theme{Name: e.Name(), Dir: dir})
		}
	}
	sort.Slice(cat.Themes, func(i, j int) bool { return cat.Themes[i].Name < cat.Themes[j].Name })
	return cat, nil
}

// ---------------------------------------------------------------------------
// Menu

func categoryMenu(cat *Category) {
	for {
		clear()
		categoryBanner(cat)
		input := prompt()
		switch strings.ToLower(input) {
		case "00":
			return // back to main menu
		case "p":
			if cat.Visual {
				visualPreviewLoop(cat)
			} else {
				textPreviewLoop(cat)
			}
		case "g":
			galleryLoop(cat)
		case "r":
			if cat.Visual {
				invalidSelection()
			} else {
				t := cat.Themes[rand.Intn(len(cat.Themes))]
				fmt.Printf("%sRandom pick: %s%s%s\n", cYellow, cGreen, t.Name, cReset)
				applyTextTheme(cat, t)
			}
		case "s":
			if cat.Visual {
				invalidSelection()
			} else {
				searchThemes(cat)
			}
		case "d":
			defaultTheme(cat)
		case "x":
			clear()
			fmt.Println(cRed + "GoodBye" + cReset)
			os.Exit(0)
		default:
			if n, ok := atoi(input); ok && n >= 1 && n <= len(cat.Themes) {
				t := cat.Themes[n-1]
				if cat.Visual {
					proto := promptLogoProtocol()
					applyVisualTheme(cat, t, proto)
				} else {
					applyTextTheme(cat, t)
				}
				os.Exit(0) // faithful to the original: exit after applying
			}
			fmt.Println(cDCyan + "Wrong transaction number!" + cReset)
		}
	}
}

func searchThemes(cat *Category) {
	fmt.Printf("%sSearch theme:%s\n", cYellow, cReset)
	fmt.Printf("%sm3tozz%s@%sfastcat\n%s↳%s ", cYellow, cDRed, cBlue, cDRed, cCyan)
	query := strings.ToLower(readLine())
	if query == "" {
		fmt.Println(cRed + "Empty search." + cReset)
		waitAnyKey()
		return
	}
	var matches []int
	for i, t := range cat.Themes {
		if strings.Contains(strings.ToLower(t.Name), query) {
			matches = append(matches, i)
		}
	}
	if len(matches) == 0 {
		fmt.Printf("%sNo themes found matching '%s'%s\n", cRed, query, cReset)
		waitAnyKey()
		return
	}
	fmt.Println(cBlue + "--- Search Results ---" + cReset)
	for _, i := range matches {
		fmt.Printf("  %s[%d]%s %s\n", cYellow, i+1, cReset, cat.Themes[i].Name)
	}
	fmt.Println(cBlue + "----------------------" + cReset)
	fmt.Printf("%sEnter number to apply, or press Enter to go back:%s ", cDCyan, cReset)
	choice := readLine()
	if n, ok := atoi(choice); ok && n >= 1 && n <= len(cat.Themes) {
		applyTextTheme(cat, cat.Themes[n-1])
		os.Exit(0)
	}
}

func defaultTheme(cat *Category) {
	if cat.Visual {
		fmt.Printf("\n%sWarning: There is no 'default' for visual themes.%s\n", cYellow, cReset)
		waitAnyKey()
		return
	}
	loader()
	os.RemoveAll(homeConfig())
	clear()
	runFastfetch()
	os.Exit(0)
}

// ---------------------------------------------------------------------------
// Applying themes

// applyTextTheme copies <theme>/fastfetch into ~/.config/fastfetch and runs fastfetch.
func applyTextTheme(cat *Category, t Theme) {
	src := filepath.Join(t.Dir, "fastfetch")
	if fi, err := os.Stat(src); err != nil || !fi.IsDir() {
		fmt.Printf("%sTheme directory '%s' not found!%s\n", cRed, t.Name, cReset)
		waitAnyKey()
		return
	}
	loader()
	if err := installTheme(src, homeConfig(), nil); err != nil {
		fatal(err.Error())
	}
	clear()
	runFastfetch()
}

// applyVisualTheme installs an image theme and rewrites the logo block
// to use the chosen rendering protocol.
func applyVisualTheme(cat *Category, t Theme, proto string) {
	src := filepath.Join(t.Dir, "fastfetch")
	loader()
	transform := func(content string) string {
		return transformVisualConfig(content, proto, homeConfig())
	}
	if err := installTheme(src, homeConfig(), transform); err != nil {
		fatal(err.Error())
	}
	clear()
	fmt.Println("Running fastfetch with the new theme. Errors (if any) will be shown below:")
	runFastfetch("--show-errors")
}

// installTheme wipes dst, copies src into it and optionally rewrites config.jsonc.
func installTheme(src, dst string, transform func(string) string) error {
	if err := os.RemoveAll(dst); err != nil {
		return err
	}
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}
	if err := copyDir(src, dst); err != nil {
		return err
	}
	if transform != nil {
		cfg := filepath.Join(dst, "config.jsonc")
		data, err := os.ReadFile(cfg)
		if err == nil {
			if err := os.WriteFile(cfg, []byte(transform(string(data))), 0o644); err != nil {
				return err
			}
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Visual (image) protocol handling

func promptLogoProtocol() string {
	fmt.Printf("\n%sChoose the image rendering protocol for this theme:%s\n", cYellow, cReset)
	fmt.Println("-------------------------------------------------------------------")
	fmt.Println(cGreen + "[1] auto" + cReset + "         - (Default) Tries to auto-detect. Works well on Kitty, Ghostty but can be unpredictable.")
	fmt.Println(cGreen + "[2] kitty" + cReset + "        - The standard Kitty protocol. Best for Kitty terminal.")
	fmt.Println(cGreen + "[3] kitty-direct" + cReset + " - Fastest protocol. Supported by: WezTerm, Warp, Kitty (PNG only), Ghostty (PNG only).")
	fmt.Println(cGreen + "[4] iterm" + cReset + "        - iTerm2 protocol. Supported by: iTerm2, WezTerm, Konsole.")
	fmt.Println(cGreen + "[5] sixel" + cReset + "        - For terminals with Sixel graphics support (e.g., foot, Contour).")
	fmt.Println("-------------------------------------------------------------------")
	protos := []string{"auto", "kitty", "kitty-direct", "iterm", "sixel"}
	for {
		fmt.Printf("%sm3tozz%s@%sfastcat\n%s↳%s ", cYellow, cDRed, cBlue, cDRed, cCyan)
		in := readLine()
		if n, ok := atoi(in); ok && n >= 1 && n <= len(protos) {
			return protos[n-1]
		}
		fmt.Println(cRed + "Invalid choice. Please try again." + cReset)
	}
}

var visualTypeRe = regexp.MustCompile(`"type":\s*"[^"]*"`)

// transformVisualConfig mirrors the original sed edits:
//   - inside the "logo": { ... } block replace "type": "..." with the protocol
//   - prefix every "source": " with the destination directory
func transformVisualConfig(content, protocol, dir string) string {
	lines := strings.Split(content, "\n")
	inLogo := false
	for i, line := range lines {
		if strings.Contains(line, `"logo"`) {
			inLogo = true
		}
		if inLogo {
			if strings.Contains(line, `"type"`) {
				lines[i] = visualTypeRe.ReplaceAllString(line, `"type": "`+protocol+`"`)
			}
			if strings.Contains(line, "}") {
				inLogo = false
			}
		}
	}
	out := strings.Join(lines, "\n")
	return strings.ReplaceAll(out, `"source": "`, `"source": "`+dir+string(filepath.Separator))
}

// ---------------------------------------------------------------------------
// Previews

func textPreviewLoop(cat *Category) {
	for {
		clear()
		fmt.Println(cYellow + "--- Preview Mode -------------------------------------------" + cReset)
		printGrid(cat.Themes, 4, 20)
		fmt.Println(cYellow + "------------------------------------------------------------" + cReset)
		fmt.Println(cDCyan + "Type a number to preview · [Q] return to menu" + cReset)
		in := prompt()
		switch strings.ToLower(in) {
		case "q":
			clear()
			return
		case "":
			invalidInput()
		default:
			if n, ok := atoi(in); ok && n >= 1 && n <= len(cat.Themes) {
				previewTextTheme(cat.Themes[n-1])
			} else {
				invalidInput()
			}
		}
	}
}

// showTextPreview renders a text theme from a temp dir (nothing is applied).
func showTextPreview(t Theme, title string) {
	tmp, err := os.MkdirTemp("", "fastcat-*")
	if err != nil {
		return
	}
	defer os.RemoveAll(tmp)
	if err := copyDir(filepath.Join(t.Dir, "fastfetch"), tmp); err != nil {
		fmt.Printf("%s%v%s\n", cRed, err, cReset)
		return
	}
	rewritePreviewConfig(tmp, false, "")

	clear()
	fmt.Println(cYellow + title + cReset)
	runFastfetch("--config", filepath.Join(tmp, "config.jsonc"))
}

func previewTextTheme(t Theme) {
	showTextPreview(t, "--- Preview Mode (not applied) ---")
	fmt.Println("\n" + cYellow + "--- End of Preview ---" + cReset)
	fmt.Println(cDCyan + "Press any key to return to menu..." + cReset)
	readKey()
}

func visualPreviewLoop(cat *Category) {
	proto := promptLogoProtocol()
	for {
		clear()
		fmt.Printf("%s--- Preview Mode (Protocol: %s) -----------------------%s\n", cYellow, proto, cReset)
		printGrid(cat.Themes, 3, 18)
		fmt.Println(cYellow + "-------------------------------------------------------------------" + cReset)
		fmt.Println(cDCyan + "Type a number to preview · [C] change protocol · [Q] return to menu" + cReset)
		in := prompt()
		switch strings.ToLower(in) {
		case "q":
			clear()
			return
		case "c":
			proto = promptLogoProtocol()
		case "":
			invalidInput()
		default:
			if n, ok := atoi(in); ok && n >= 1 && n <= len(cat.Themes) {
				previewVisualTheme(cat.Themes[n-1], proto)
			} else {
				invalidInput()
			}
		}
	}
}

// showVisualPreview renders a Visuals theme from a temp dir (nothing is applied).
func showVisualPreview(t Theme, proto, title string) {
	tmp, err := os.MkdirTemp("", "fastcat-*")
	if err != nil {
		return
	}
	defer os.RemoveAll(tmp)
	if err := copyDir(filepath.Join(t.Dir, "fastfetch"), tmp); err != nil {
		fmt.Printf("%s%v%s\n", cRed, err, cReset)
		return
	}
	rewritePreviewConfig(tmp, true, proto)

	clear()
	fmt.Println(cYellow + title + cReset)
	runFastfetch("--config", filepath.Join(tmp, "config.jsonc"), "--show-errors")
}

func previewVisualTheme(t Theme, proto string) {
	showVisualPreview(t, proto, "--- Preview: "+t.Name+" (not applied) ---")
	fmt.Println("\n" + cYellow + "--- End of Preview ---" + cReset)
	fmt.Println(cDCyan + "Press any key to continue..." + cReset)
	readKey()
}

// galleryLoop is the interactive theme browser: flip through previews with
// arrow keys, apply with Enter — no need to type theme numbers one by one.
func galleryLoop(cat *Category) {
	proto := ""
	if cat.Visual {
		proto = promptLogoProtocol()
	}
	i := 0
	for {
		t := cat.Themes[i]
		title := fmt.Sprintf("--- Gallery [%d/%d]: %s (not applied) ---", i+1, len(cat.Themes), t.Name)
		if cat.Visual {
			showVisualPreview(t, proto, title)
		} else {
			showTextPreview(t, title)
		}
		fmt.Printf("\n%s←/→ (или j/k) листать · Enter применить · R случайная · Q в меню%s\n", cDCyan, cReset)

		switch key := readKey(); key {
		case "right", "down", "j", "d", "n", " ":
			i = (i + 1) % len(cat.Themes)
		case "left", "up", "k", "a", "p":
			i = (i - 1 + len(cat.Themes)) % len(cat.Themes)
		case "r":
			i = rand.Intn(len(cat.Themes))
		case "enter":
			if cat.Visual {
				applyVisualTheme(cat, cat.Themes[i], proto)
			} else {
				applyTextTheme(cat, cat.Themes[i])
			}
			os.Exit(0) // faithful to the original: exit after applying
		case "q", "esc", "":
			clear()
			return
		}
	}
}

var imgExts = map[string]bool{".png": true, ".jpg": true, ".jpeg": true, ".webp": true, ".gif": true, ".bmp": true}

// rewritePreviewConfig fixes asset paths so a theme can run from a temp dir.
func rewritePreviewConfig(dir string, visual bool, proto string) {
	cfg := filepath.Join(dir, "config.jsonc")
	data, err := os.ReadFile(cfg)
	if err != nil {
		return
	}
	content := string(data)
	if visual {
		content = transformVisualConfig(content, proto, dir)
	} else {
		content = strings.ReplaceAll(content, "~/.config/fastfetch/", dir+string(filepath.Separator))
		content = strings.ReplaceAll(content, "$HOME/.config/fastfetch/", dir+string(filepath.Separator))
		// Bare image filenames -> absolute temp paths.
		if entries, err := os.ReadDir(dir); err == nil {
			for _, e := range entries {
				if !e.IsDir() && imgExts[strings.ToLower(filepath.Ext(e.Name()))] {
					content = strings.ReplaceAll(content,
						`"source": "`+e.Name()+`"`,
						`"source": "`+filepath.Join(dir, e.Name())+`"`)
				}
			}
		}
	}
	os.WriteFile(cfg, []byte(content), 0o644)
}

func invalidSelection() {
	fmt.Printf("\n%sInvalid selection!%s\n", cRed, cReset)
	waitAnyKey()
}

func invalidInput() {
	fmt.Println(cRed + "Invalid input!" + cReset)
	waitAnyKey()
}
