# Fastcat-go

<p align="center"><img src="https://github.com/user-attachments/assets/1f24e219-07fa-48e4-a5b5-2844e50c7d35"></p>

---

## What is FastCat?

FastCat is a collection of ready-made configurations for [fastfetch](https://github.com/fastfetch-cli/fastfetch). Launch the menu, pick the config you like and apply it with a single keypress.

Theme categories:

- **Small Themes** — minimalist text logos
- **Large Themes** — large ASCII art
- **Visuals Themes** — image-based themes (require a graphics-capable terminal: Kitty, WezTerm, iTerm2, Ghostty, etc.)

---

## <a href="https://www.instagram.com/textzuhree/"><img src="https://img.shields.io/badge/textzuhree-black?&logo=instagram&logoColor=purple" hspace="10"  align="right" /></a>Install (Automatic)

```bash
git clone --depth 1 https://github.com/m3tozz/FastCat.git && cd FastCat && bash ./fastcat.sh --help
```

### Install (Manual)

```bash
git clone https://github.com/m3tozz/FastCat.git
cd FastCat
bash ./fastcat.sh --help
```

### Go Edition (a single binary instead of bash scripts)

```bash
cd FastCat/go
go build -o fastcat .
./fastcat --shell
```

For details, see [go/README.md](go/README.md).

---

## Main Commands

| Command | Action |
|---|---|
| `--shell` | interactive menu (Small / Large / Visuals) |
| `--backup` | back up your current `~/.config/fastfetch` |
| `--icons` | install the CascadiaCode Nerd Font |
| `--version` / `--about` / `--help` | version / about / help |

For the bash edition: `bash ./fastcat.sh -s`, `bash ./fastcat.sh --backup`, `bash ./install-icons.sh`.

### Choosing a Theme

1. Launch the menu (`--shell`).
2. Pick a category: Small, Large or Visuals.
3. Browse the themes:
   - **`[G] Gallery`** — interactive gallery: `←`/`→` (or `j`/`k`) flip through previews, `Enter` applies the theme, `R` picks a random one;
   - **`[P] Preview`** — preview a theme by its number;
   - **`[S] Search`** — search themes by name;
   - **`[R] Random`** — a random theme.

A preview changes nothing on your system: the config runs from a temporary folder. Only applying a theme (Enter / entering its number) overwrites `~/.config/fastfetch`.

### Installing the Icon Font (CascadiaCode Nerd Font)

```bash
bash ./install-icons.sh        # bash edition
./go/fastcat --icons           # Go edition
```

<details>
<summary>Download the required icon pack manually (for Arch Linux)</summary>
Command: <code>git clone https://aur.archlinux.org/ttf-meslo-nerd-font-powerlevel10k.git && cd ttf-meslo-nerd-font-powerlevel10k && makepkg -si && cd ..</code>
</details>

### Backing Up Your Own fastfetch Config

```bash
bash ./fastcat.sh --backup
```

---

## Themes

> [!IMPORTANT]
> Since FastCat can now preview its themes, the [FastCat Themes site](https://m3tozz.github.io/FastCat-Themes/) is no longer needed.
> After selecting the **Large** or **Small** section, press **P (Preview)** or **G (Gallery)** and browse the themes right in your terminal.

<details>
<summary>How do I add my own fastfetch theme to the FastCat collection?</summary>

Fork this repository, add a folder with your theme (it must contain a <code>fastfetch/</code> subfolder with a <code>config.jsonc</code> inside) to <code>Small-Themes</code>, <code>Large-Themes</code> or <code>Visuals-Themes</code>, and open a pull request.
</details>

---

<p align="center"><img src="https://github.com/user-attachments/assets/aa96371a-8c47-4ce8-9c89-2e482e726bf0">
<a href="https://discord.com/invite/sQwYCZer95"><img src="https://img.shields.io/badge/Join the Community Discord Server-black?&logo=discord" hspace="10"  align="center" /></a>
<a href="https://matrix.to/#/#neocat-channels:matrix.org"><img src="https://img.shields.io/badge/Join the Community Matrix Server-black?&logo=matrix" hspace="10"  align="center" /></a></p>

---
