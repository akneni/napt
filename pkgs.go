package main

import (
	"fmt"
	"os"
	"strings"
)

type Pkgs struct {
	pkgs         []string
	spacing_char rune
	num_spacing  int
	with_pkgs    bool
}

func (self Pkgs) toString() string {
	lst := make([]string, 0)
	for _, pkg := range self.pkgs {
		prefix := strings.Repeat(string(self.spacing_char), int(self.num_spacing))
		if !self.with_pkgs {
			pkg = "pkgs." + pkg
		}
		payload := prefix + pkg

		lst = append(lst, payload)
	}
	return strings.Join(lst, "\n")
}

func (self Pkgs) ToFile(filename, target_name string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading from nix config file: ", err)
		os.Exit(1)
	}

	text := string(data)
	lines := strings.Split(text, "\n")

	start_idx := -1

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, target_name) && strings.Contains(line, "=") && strings.HasSuffix(line, "[") {
			start_idx = i + 1

			for strings.Contains(line, "  ") {
				line = strings.ReplaceAll(line, "  ", " ")
			}
			for strings.Contains(line, " ;") {
				line = strings.ReplaceAll(line, " ;", ";")
			}
			break
		}
	}

	if start_idx == -1 {
		fmt.Println("Error parsing target file, target not found")
		os.Exit(1)
	}

	new_lines := make([]string, len(lines[:start_idx]))
	copy(new_lines, lines[:start_idx])

	for _, pkg := range strings.Split(self.toString(), "\n") {
		new_lines = append(new_lines, pkg)
	}

	end_reached := false
	for _, line := range lines[start_idx:] {
		if !end_reached {
			og_line := line
			line = strings.TrimSpace(line)
			for strings.Contains(line, " ;") {
				line = strings.ReplaceAll(line, " ;", ";")
			}
			if strings.TrimSpace(line) == "];" {
				end_reached = true
				new_lines = append(new_lines, og_line)
			}
		} else {
			new_lines = append(new_lines, line)
		}
	}

	final_text := strings.Join(new_lines, "\n")
	err = os.WriteFile(filename, []byte(final_text), 'w')
	if err != nil {
		fmt.Println("Error writing to file.")
		os.Exit(1)
	}
}

func (Pkgs) FromFile(filename, target_name string) Pkgs {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading from nix config file: ", err)
		os.Exit(1)
	}

	text := string(data)
	lines := strings.Split(text, "\n")

	start_idx := -1
	with_pkgs := false

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, target_name) && strings.Contains(line, "=") && strings.HasSuffix(line, "[") {
			start_idx = i + 1

			for strings.Contains(line, "  ") {
				line = strings.ReplaceAll(line, "  ", " ")
			}
			for strings.Contains(line, " ;") {
				line = strings.ReplaceAll(line, " ;", ";")
			}

			if strings.Contains(line, "with pkgs;") {
				with_pkgs = true
			}

			break
		}
	}

	if start_idx == -1 {
		fmt.Println("Error parsing target file, target not found")
		os.Exit(1)
	}

	lst := make([]string, 0)

	for i := start_idx; i < len(lines); i++ {
		line := lines[i]
		for strings.Contains(line, " ;") {
			line = strings.ReplaceAll(line, " ;", ";")
		}
		if strings.TrimSpace(line) == "];" {
			break
		}

		lst = append(lst, lines[i])
	}

	pkgs := Pkgs{}.fromList(lst)
	pkgs.with_pkgs = with_pkgs

	return pkgs

}

func (Pkgs) fromList(lst []string) Pkgs {
	if len(lst) == 0 {
		return Pkgs{
			pkgs:         lst,
			spacing_char: 0,
			num_spacing:  '\t',
			with_pkgs:    false,
		}
	}

	num_spacing := 0
	var spacing_char rune = 0

	for _, ch := range lst[0] {
		if num_spacing == 0 && spacing_char == 0 {
			if ch == ' ' {
				spacing_char = ' '
			} else if ch == '\t' {
				spacing_char = '\t'
			} else {
				break
			}
			num_spacing += 1
			continue
		}

		if ch == spacing_char {
			num_spacing += 1
		} else {
			break
		}
	}

	for i, pkg := range lst {
		lst[i] = strings.TrimSpace(pkg)
	}

	with_pkgs := true

	if lst[0][:4] == "pkgs." {
		with_pkgs = false
		for i, pkg := range lst {
			lst[i] = pkg[:4]
		}
	}

	return Pkgs{
		pkgs:         lst,
		spacing_char: spacing_char,
		num_spacing:  num_spacing,
		with_pkgs:    with_pkgs,
	}
}
