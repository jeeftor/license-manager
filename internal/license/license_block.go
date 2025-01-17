package license

import (
	"license-manager/internal/comment"
	"license-manager/internal/styles"
	"strings"
)

const (
	markerStart = "​" // Zero-width space
	markerEnd   = "‌" // Zero-width non-joiner
)

// LicenseBlock represents a complete license block with style information
type LicenseBlock struct {
	comment *comment.Comment
}

// NewLicenseBlock creates a new license block with the given style and content
func NewLicenseBlock(style styles.CommentLanguage, header, body, footer string) *LicenseBlock {
	c := comment.NewComment(style, styles.HeaderFooterStyle{
		Header: header,
		Footer: footer,
	}, body)

	return &LicenseBlock{
		comment: c,
	}
}

// String returns the complete license block as a string
func (lb *LicenseBlock) String() string {
	return lb.comment.String()
}

// Clone creates a deep copy of the license block
func (lb *LicenseBlock) Clone() *LicenseBlock {
	return &LicenseBlock{
		comment: lb.comment.Clone(),
	}
}

// GetStyle returns the current comment style
func (lb *LicenseBlock) GetStyle() styles.CommentLanguage {
	return lb.comment.GetStyle()
}

// SetStyle updates the comment style
func (lb *LicenseBlock) SetStyle(style styles.CommentLanguage) {
	lb.comment.SetStyle(style)
}

// GetBody returns the license body content
func (lb *LicenseBlock) GetBody() string {
	return lb.comment.GetBody()
}

// SetBody updates the license body content
func (lb *LicenseBlock) SetBody(body string) {
	lb.comment.SetBody(body)
}

// GetHeader returns the license header content
func (lb *LicenseBlock) GetHeader() string {
	return lb.comment.GetHeader()
}

// SetHeader updates the license header content
func (lb *LicenseBlock) SetHeader(header string) {
	lb.comment.SetHeaderFooterStyle(styles.HeaderFooterStyle{
		Header: header,
		Footer: lb.comment.GetFooter(),
	})
}

// GetFooter returns the license footer content
func (lb *LicenseBlock) GetFooter() string {
	return lb.comment.GetFooter()
}

// SetFooter updates the license footer content
func (lb *LicenseBlock) SetFooter(footer string) {
	lb.comment.SetHeaderFooterStyle(styles.HeaderFooterStyle{
		Header: lb.comment.GetHeader(),
		Footer: footer,
	})
}

// Helper functions for working with markers
func hasMarkers(text string) bool {
	return strings.Contains(text, markerStart) && strings.Contains(text, markerEnd)
}

func addMarkers(text string) string {
	if hasMarkers(text) {
		return text
	}
	return markerStart + text + markerEnd
}

func stripMarkers(text string) string {
	text = strings.ReplaceAll(text, markerStart, "")
	text = strings.ReplaceAll(text, markerEnd, "")
	return text
}
