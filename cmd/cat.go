package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/jeeftor/license-manager/internal/config"
	"github.com/jeeftor/license-manager/internal/logger"
	"github.com/jeeftor/license-manager/internal/processor"
	"github.com/spf13/cobra"
)

var (
	cfgShowLineNumbers bool
	cfgShowSpecial     bool
)

func renderLegend() string {
	// Define colors and letters for different sections
	sections := []struct {
		letter string
		desc   string
		color  *color.Color
	}{
		{"P", "Preamble (build tags, package)", color.New(color.BgRed)},
		{"H", "License Header", color.New(color.BgGreen)},
		{"L", "License Body", color.New(color.BgYellow)},
		{"F", "License Footer", color.New(color.BgBlue)},
		{"C", "Code", color.New(color.BgMagenta)},
	}

	specialStyle := color.New(color.FgHiBlue)

	// Calculate the width needed for the main legend
	maxDescWidth := 0
	for _, s := range sections {
		if len(s.desc) > maxDescWidth {
			maxDescWidth = len(s.desc)
		}
	}

	var output strings.Builder
	output.WriteString("Legend:\n")

	// Determine if we have enough space for two columns (assuming 80 char terminal width)
	// Each section needs about 40 chars with padding
	termWidth := 80 // Could make this dynamic if needed
	useColumns := termWidth >= 80

	if useColumns {
		// Print sections in two columns
		halfLen := (len(sections) + 1) / 2
		for i := 0; i < halfLen; i++ {
			s1 := sections[i]
			prefix1 := s1.color.Add(color.FgBlack).Add(color.Bold).Sprintf(" %s ", s1.letter)

			line := fmt.Sprintf("%-4s %-30s", prefix1, s1.desc)

			// If there's a second column
			if i+halfLen < len(sections) {
				s2 := sections[i+halfLen]
				prefix2 := s2.color.Add(color.FgBlack).Add(color.Bold).Sprintf(" %s ", s2.letter)
				line += fmt.Sprintf("    %-4s %-30s", prefix2, s2.desc)
			}

			output.WriteString(line + "\n")
		}
	} else {
		// Single column format
		for _, s := range sections {
			prefix := s.color.Add(color.FgBlack).Add(color.Bold).Sprintf(" %s ", s.letter)
			output.WriteString(fmt.Sprintf("%-4s %s\n", prefix, s.desc))
		}
	}

	// Add special characters legend if enabled
	if cfgShowSpecial {
		output.WriteString("\nSpecial Characters:\n")
		if useColumns {
			output.WriteString(fmt.Sprintf("%-4s %-30s    %-4s %-30s\n",
				specialStyle.Sprint("·"), "Space",
				specialStyle.Sprint("↵"), "Line Feed (LF)"))
			output.WriteString(fmt.Sprintf("%-4s %-30s    %-4s %-30s\n",
				specialStyle.Sprint("⏎"), "Carriage Return (CR)",
				specialStyle.Sprint("⏎↵"), "CRLF"))
		} else {
			output.WriteString(fmt.Sprintf("%-4s %s\n", specialStyle.Sprint("·"), "Space"))
			output.WriteString(fmt.Sprintf("%-4s %s\n", specialStyle.Sprint("↵"), "Line Feed (LF)"))
			output.WriteString(fmt.Sprintf("%-4s %s\n", specialStyle.Sprint("⏎"), "Carriage Return (CR)"))
			output.WriteString(fmt.Sprintf("%-4s %s\n", specialStyle.Sprint("⏎↵"), "CRLF"))
		}
	}

	output.WriteString("\n")
	return output.String()
}

// renderSpecialChars replaces whitespace with visible characters
func renderSpecialChars(text string) string {
	if !cfgShowSpecial {
		return text
	}

	specialStyle := color.New(color.FgHiBlue)

	// Replace spaces with visible dot
	text = strings.ReplaceAll(text, " ", specialStyle.Sprint("·"))

	// Replace CRLF first (order matters)
	text = strings.ReplaceAll(text, "\r\n", specialStyle.Sprint("⏎↵"))

	// Then replace individual CR and LF
	text = strings.ReplaceAll(text, "\r", specialStyle.Sprint("⏎"))
	text = strings.ReplaceAll(text, "\n", specialStyle.Sprint("↵"))

	return text
}

// addLetterPrefix adds letters with colored backgrounds and optional line numbers
func addLetterPrefix(
	text string,
	letter string,
	col *color.Color,
	startLine int,
	showLineNumbers bool,
) string {

	lines := strings.Split(text, "\n")

	// Determine padding width based on total number of lines
	padWidth := 0
	if showLineNumbers {
		totalLines := startLine + len(lines)
		padWidth = 2 // minimum width
		if totalLines >= 100 {
			padWidth = 3
		}
		if totalLines >= 1000 {
			padWidth = 4
		}
	}

	// Create a background colored letter (black text on colored background)
	letterPrefix := col.Add(color.FgBlack).Add(color.Bold).Sprintf(" %s ", letter) + " "

	// Create line number style (white text on cyan background)
	lineNumStyle := color.New(color.FgBlack).Add(color.Bold).Add(color.BgCyan)

	// Format each line
	for i := range lines {
		lineNum := startLine + i + 1 // Start at 1 instead of 0

		// Process special characters if enabled
		lineContent := renderSpecialChars(lines[i])

		if lines[i] == "" {
			// For empty lines
			if showLineNumbers {
				// Format line number with cyan background
				lineNumStr := lineNumStyle.Sprintf("%*d", padWidth, lineNum)
				lines[i] = fmt.Sprintf("%s %s", lineNumStr, letterPrefix)
			} else {
				lines[i] = letterPrefix
			}
		} else {
			if showLineNumbers {
				// Format line number with cyan background
				lineNumStr := lineNumStyle.Sprintf("%*d", padWidth, lineNum)
				lines[i] = fmt.Sprintf("%s %s%s", lineNumStr, letterPrefix, lineContent)
			} else {
				lines[i] = letterPrefix + lineContent
			}
		}
	}

	// If showing special chars, add visible newline at end of each line except last
	if cfgShowSpecial && len(lines) > 0 {
		specialStyle := color.New(color.FgHiBlue)
		for i := 0; i < len(lines)-1; i++ {
			lines[i] = lines[i] + specialStyle.Sprint("↵")
		}
	}

	return strings.Join(lines, "\n")
}

// processCatPatterns handles pattern processing specifically for the cat command
func processCatPatterns(patterns []string) []string {
	if patterns == nil {
		return nil
	}

	var result []string
	for _, pattern := range patterns {
		// Split pattern on commas and process each part
		parts := strings.Split(pattern, ",")
		for _, part := range parts {
			if trimmed := strings.TrimSpace(part); trimmed != "" {
				result = append(result, trimmed)
			}
		}
	}
	return result
}

var catCmd = &cobra.Command{
	Use:   "cat [files...]",
	Short: "Display file with colored component blocks",
	Long: `Show file contents with colored block prefixes indicating different components (preamble, license, code).

Use -n or --line-numbers to display line numbers.
Use -s or --special to show special characters (spaces, CR, LF).

Files can be specified either:
- As arguments: license-manager cat -n -s file1.go file2.go
- Using --input flag: license-manager cat -n -s --input "*.go"
- Or both: license-manager cat -n -s --input "lib/*.go" cmd/*.go`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Process all input sources
		var allInputs []string

		// Process --input flag patterns
		if cfgInputs != nil {
			allInputs = append(allInputs, processCatPatterns(cfgInputs)...)
		}

		// Process command line arguments
		if len(args) > 0 {
			allInputs = append(allInputs, processCatPatterns(args)...)
		}

		// Ensure we have at least one input source
		if len(allInputs) == 0 {
			return fmt.Errorf("no input files specified - use arguments or --input flag")
		}

		appCfg := config.AppConfig{
			Inputs: strings.Join(allInputs, ","), // Join processed inputs for compatibility
			Skips:  ProcessPatterns(cfgSkips),    // Use original ProcessPatterns for skips

			// Style settings
			HeaderStyle:  cfgPresetStyle,
			CommentStyle: "go",

			LogLevel:    logger.ParseLogLevel(cfgLogLevel),
			IsPreCommit: false,
		}

		procCfg, err := appCfg.ToProcessorConfig()
		if err != nil {
			return err
		}

		// Define colors (using background colors) and letters for different sections
		sections := []struct {
			letter string
			desc   string
			color  *color.Color
		}{
			{"P", "Preamble (build tags, package)", color.New(color.BgRed)},
			{"H", "License Header", color.New(color.BgGreen)},
			{"B", "License Body", color.New(color.BgYellow)},
			{"F", "License Footer", color.New(color.BgBlue)},
			{"C", "Code", color.New(color.BgMagenta)},
		}

		// Print legend
		fmt.Print(renderLegend())

		// Create processor
		p := processor.NewFileProcessor(procCfg)
		files, err := p.PrepareOperation()
		if err != nil {
			return err
		}

		for _, file := range files {
			if len(files) > 1 {
				fmt.Printf("==> %s <==\n", file)
			}

			_, components, err := p.GetFileComponents(file)

			if err != nil {
				return fmt.Errorf("error processing file %s: %v", file, err)
			}

			// Track current line number
			currentLine := 0

			// Print each section with letter prefixes
			output := []string{}

			if components.Preamble != "" {
				output = append(
					output,
					addLetterPrefix(
						components.Preamble,
						"P",
						sections[0].color,
						currentLine,
						cfgShowLineNumbers,
					),
				)
				currentLine += len(strings.Split(components.Preamble, "\n"))
			}
			// process license block

			lic := components.FullLicenseBlock

			if lic != nil {
				licenseStrings := strings.Split(lic.String, "\n")

				for i, line := range licenseStrings {
					if i < lic.BodyStart {
						output = append(output, addLetterPrefix(
							line,
							"H",
							sections[1].color,
							currentLine+i,
							cfgShowLineNumbers,
						))
					} else if i < lic.FooterStart {
						output = append(output, addLetterPrefix(
							line,
							"L",
							sections[2].color,
							currentLine+i,
							cfgShowLineNumbers,
						))
					} else {
						output = append(output, addLetterPrefix(
							line,
							"F",
							sections[3].color,
							currentLine+i,
							cfgShowLineNumbers,
						))
					}
				}

				currentLine += len(strings.Split(components.FullLicenseBlock.String, "\n"))
			}
			if components.Rest != "" {
				output = append(
					output,
					addLetterPrefix(
						components.Rest,
						"C",
						sections[4].color,
						currentLine,
						cfgShowLineNumbers,
					),
				)
			}

			// Join all sections, preserving empty lines between them
			fmt.Println(strings.Join(output, "\n"))

			if len(files) > 1 {
				fmt.Println() // Add spacing between files
			}
		}

		cmd.SilenceUsage = true
		return nil
	},
}

func init() {
	rootCmd.AddCommand(catCmd)

	// Add line number flag
	catCmd.Flags().StringSliceVar(&cfgInputs, "input", nil, "input file patterns (optional)")
	catCmd.Flags().BoolVarP(&cfgShowLineNumbers, "line-numbers", "n", false, "show line numbers")
	catCmd.Flags().BoolVarP(&cfgShowSpecial, "special", "s", false, "show special characters")
}
