package processor

type HeaderFooterStyle struct {
	Name        string
	Description string
	Header      string
	Footer      string
}

var PresetStyles = map[string]HeaderFooterStyle{
	"simple": {
		Name:        "Simple",
		Description: "Simple lines before and after",
		Header:      "----------------------------------------",
		Footer:      "----------------------------------------",
	},
	"hash": {
		Name:        "Hash",
		Description: "Hash symbol borders",
		Header:      "######################################",
		Footer:      "######################################",
	},
	"box": {
		Name:        "Box",
		Description: "Box style with corners",
		Header:      "+------------------------------------+",
		Footer:      "+------------------------------------+",
	},
	"brackets": {
		Name:        "Brackets",
		Description: "Square bracket style",
		Header:      "[ License Start ]----------------------",
		Footer:      "[ License End ]------------------------",
	},
	"stars": {
		Name:        "Stars",
		Description: "Star border style",
		Header:      "****************************************",
		Footer:      "****************************************",
	},
	"modern": {
		Name:        "Modern",
		Description: "Modern style with unicode blocks",
		Header:      "■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■",
		Footer:      "■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■",
	},
	"arrows": {
		Name:        "Arrows",
		Description: "Arrow style borders",
		Header:      ">>>>>>> LICENSE HEADER >>>>>>>>>>>>>>>>>>",
		Footer:      "<<<<<<<<<<<<<<<<<<<< LICENSE END <<<<<<<",
	},
	"elegant": {
		Name:        "Elegant",
		Description: "Elegant style with dashes and brackets",
		Header:      "---[ Begin License ]-------------------",
		Footer:      "---[ End License ]---------------------",
	},
	"wave": {
		Name:        "Wave",
		Description: "Wavy style border",
		Header:      "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~",
		Footer:      "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~",
	},
	"dots": {
		Name:        "Dots",
		Description: "Dotted border style",
		Header:      "......................................",
		Footer:      "......................................",
	},
	"equals": {
		Name:        "Equals",
		Description: "Equal signs border",
		Header:      "======================================",
		Footer:      "======================================",
	},
	"decorative": {
		Name:        "Decorative",
		Description: "Decorative style with unicode symbols",
		Header:      "♦═══════════[ LICENSE ]══════════════♦",
		Footer:      "♦══════════[ END LICENSE ]═══════════♦",
	},
	"minimal": {
		Name:        "Minimal",
		Description: "Minimal style with thin lines",
		Header:      "─────────────────────────────────────",
		Footer:      "─────────────────────────────────────",
	},
	"classic": {
		Name:        "Classic",
		Description: "Classic style with asterisks and title",
		Header:      "**** BEGIN LICENSE BLOCK ****",
		Footer:      "**** END LICENSE BLOCK ****",
	},
	"angular": {
		Name:        "Angular",
		Description: "Angular style with forward slashes",
		Header:      "//////////// LICENSE START ////////////",
		Footer:      "//////////// LICENSE END //////////////",
	},
	"banner": {
		Name:        "Banner",
		Description: "Banner style with pipe borders",
		Header:      "|====================================|",
		Footer:      "|====================================|",
	},
	"retro": {
		Name:        "Retro",
		Description: "Retro style with plus signs",
		Header:      "++++++++++++++++++++++++++++++++++++",
		Footer:      "++++++++++++++++++++++++++++++++++++",
	},
	"clean": {
		Name:        "Clean",
		Description: "Clean style with triple dashes",
		Header:      "-----------------------------------",
		Footer:      "-----------------------------------",
	},
	"double": {
		Name:        "Double",
		Description: "Double line border",
		Header:      "══════════════════════════════════",
		Footer:      "══════════════════════════════════",
	},
	"branded": {
		Name:        "Branded",
		Description: "Corporate style with copyright symbol",
		Header:      "© ─────────[ LICENSE ]──────────── ©",
		Footer:      "© ────────[ END LICENSE ]────────── ©",
	},
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
