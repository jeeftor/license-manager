package processor

import "strings"

const (
	// Invisible markers for license blocks
	markerStart = "\u200B" // Zero-Width Space
	markerEnd   = "\u200C" // Zero-Width Non-Joiner
)

type HeaderFooterStyle struct {
	Name        string
	Description string
	Header      string
	Footer      string
}

func addMarkers(text string) string {
	return markerStart + text + markerEnd
}

var PresetStyles = map[string]HeaderFooterStyle{
	"simple": {
		Name:        "Simple",
		Description: "Simple lines before and after",
		Header:      addMarkers("----------------------------------------"),
		Footer:      addMarkers("----------------------------------------"),
	},
	"hash": {
		Name:        "Hash",
		Description: "Hash symbol borders",
		Header:      addMarkers("######################################"),
		Footer:      addMarkers("######################################"),
	},
	"box": {
		Name:        "Box",
		Description: "Box style with corners",
		Header:      addMarkers("+------------------------------------+"),
		Footer:      addMarkers("+------------------------------------+"),
	},
	"brackets": {
		Name:        "Brackets",
		Description: "Square bracket style",
		Header:      addMarkers("[ License Start ]----------------------"),
		Footer:      addMarkers("[ License End ]------------------------"),
	},
	"stars": {
		Name:        "Stars",
		Description: "Star border style",
		Header:      addMarkers("****************************************"),
		Footer:      addMarkers("****************************************"),
	},
	"modern": {
		Name:        "Modern",
		Description: "Modern style with unicode blocks",
		Header:      addMarkers("■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■"),
		Footer:      addMarkers("■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■"),
	},
	"arrows": {
		Name:        "Arrows",
		Description: "Arrow style borders",
		Header:      addMarkers(">>>>>>> LICENSE HEADER >>>>>>>>>>>>>>>>>>"),
		Footer:      addMarkers("<<<<<<<<<<<<<<<<<<<< LICENSE END <<<<<<<"),
	},
	"elegant": {
		Name:        "Elegant",
		Description: "Elegant style with dashes and brackets",
		Header:      addMarkers("---[ Begin License ]-------------------"),
		Footer:      addMarkers("---[ End License ]---------------------"),
	},
	"wave": {
		Name:        "Wave",
		Description: "Wavy style border",
		Header:      addMarkers("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~"),
		Footer:      addMarkers("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~"),
	},
	"dots": {
		Name:        "Dots",
		Description: "Dotted border style",
		Header:      addMarkers("......................................"),
		Footer:      addMarkers("......................................"),
	},
	"equals": {
		Name:        "Equals",
		Description: "Equal signs border",
		Header:      addMarkers("======================================"),
		Footer:      addMarkers("======================================"),
	},
	"decorative": {
		Name:        "Decorative",
		Description: "Decorative style with unicode symbols",
		Header:      addMarkers("♦═══════════[ LICENSE ]══════════════♦"),
		Footer:      addMarkers("♦══════════[ END LICENSE ]═══════════♦"),
	},
	"minimal": {
		Name:        "Minimal",
		Description: "Minimal style with thin lines",
		Header:      addMarkers("─────────────────────────────────────"),
		Footer:      addMarkers("─────────────────────────────────────"),
	},
	"classic": {
		Name:        "Classic",
		Description: "Classic style with asterisks and title",
		Header:      addMarkers("**** BEGIN LICENSE BLOCK ****"),
		Footer:      addMarkers("**** END LICENSE BLOCK ****"),
	},
	"angular": {
		Name:        "Angular",
		Description: "Angular style with forward slashes",
		Header:      addMarkers("//////////// LICENSE START ////////////"),
		Footer:      addMarkers("//////////// LICENSE END //////////////"),
	},
	"banner": {
		Name:        "Banner",
		Description: "Banner style with pipe borders",
		Header:      addMarkers("|====================================|"),
		Footer:      addMarkers("|====================================|"),
	},
	"retro": {
		Name:        "Retro",
		Description: "Retro style with plus signs",
		Header:      addMarkers("++++++++++++++++++++++++++++++++++++"),
		Footer:      addMarkers("++++++++++++++++++++++++++++++++++++"),
	},
	"clean": {
		Name:        "Clean",
		Description: "Clean style with triple dashes",
		Header:      addMarkers("-----------------------------------"),
		Footer:      addMarkers("-----------------------------------"),
	},
	"double": {
		Name:        "Double",
		Description: "Double line border",
		Header:      addMarkers("══════════════════════════════════"),
		Footer:      addMarkers("══════════════════════════════════"),
	},
	"branded": {
		Name:        "Branded",
		Description: "Corporate style with copyright symbol",
		Header:      addMarkers("© ─────────[ LICENSE ]──────────── ©"),
		Footer:      addMarkers("© ────────[ END LICENSE ]────────── ©"),
	},
}

// Helper functions to detect markers
func hasMarkers(text string) bool {
	return strings.Contains(text, markerStart) && strings.Contains(text, markerEnd)
}

func stripMarkers(text string) string {
	text = strings.ReplaceAll(text, markerStart, "")
	text = strings.ReplaceAll(text, markerEnd, "")
	return text
}

func findLicenseBlock(content string) (start, end int) {
	start = strings.Index(content, markerStart)
	if start == -1 {
		return -1, -1
	}

	// Look for the end marker after the second marker start
	secondStart := strings.Index(content[start+1:], markerStart)
	if secondStart == -1 {
		return -1, -1
	}
	secondStart += start + 1

	// Find the end marker after the second start marker
	end = strings.Index(content[secondStart:], markerEnd)
	if end == -1 {
		return -1, -1
	}
	end += secondStart + len(markerEnd)

	return start, end
}

// GetPresetStyle returns a HeaderFooterStyle by name, or a default if not found
func GetPresetStyle(name string) HeaderFooterStyle {
	if style, ok := PresetStyles[name]; ok {
		return style
	}
	// Return simple style as default
	return PresetStyles["simple"]
}

// ListPresetStyles returns a slice of all available style names
func ListPresetStyles() []string {
	var styles []string
	for name := range PresetStyles {
		styles = append(styles, name)
	}
	return styles
}
