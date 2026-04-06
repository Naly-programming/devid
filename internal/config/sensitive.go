package config

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
)

// Patterns that look like secrets.
var sensitivePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)^sk[-_]`),                        // API keys (sk-ant-, sk_live_, etc)
	regexp.MustCompile(`(?i)^(ghp|gho|ghu|ghs|ghr)_`),       // GitHub tokens
	regexp.MustCompile(`(?i)^xox[bpsa]-`),                    // Slack tokens
	regexp.MustCompile(`(?i)^eyJ[a-zA-Z0-9]`),               // JWT tokens
	regexp.MustCompile(`(?i)(password|passwd|secret|token)=`), // key=value pairs
	regexp.MustCompile(`(?i)^AKIA[0-9A-Z]`),                 // AWS access keys
	regexp.MustCompile(`(?i)^[a-f0-9]{40}$`),                // 40-char hex (git SHAs, but also tokens)
	regexp.MustCompile(`-----BEGIN (RSA |EC )?PRIVATE KEY`),  // Private keys
}

// SensitiveWarning describes a potential secret found in identity data.
type SensitiveWarning struct {
	Section string
	Field   string
	Value   string
	Pattern string
}

// CheckSensitive scans the non-private sections of an identity for values
// that look like secrets. Returns warnings for any matches found.
func CheckSensitive(id *Identity) []SensitiveWarning {
	// Encode to TOML, then strip the [private] section and check remaining values
	clean := id.WithoutPrivate()

	var buf strings.Builder
	toml.NewEncoder(&buf).Encode(clean)

	// Also do a field-level check on string values
	var warnings []SensitiveWarning

	warnings = append(warnings, checkSection("identity", map[string]string{
		"name": id.Identity.Name, "tone": id.Identity.Tone,
		"comments": id.Identity.Comments, "responses": id.Identity.Responses,
		"pace": id.Identity.Pace,
	})...)

	warnings = append(warnings, checkSection("conventions", map[string]string{
		"formatting": id.Conventions.Formatting, "pr_style": id.Conventions.PRStyle,
		"commit_style": id.Conventions.CommitStyle, "error_handling": id.Conventions.ErrorHandling,
		"naming": id.Conventions.Naming,
	})...)

	warnings = append(warnings, checkSection("ai", map[string]string{
		"verbosity": id.AI.Verbosity, "confirmation": id.AI.Confirmation,
		"suggestions": id.AI.Suggestions, "code_comments": id.AI.CodeComments,
		"tests": id.AI.Tests,
	})...)

	// Check string slices
	for _, items := range []struct {
		section string
		field   string
		values  []string
	}{
		{"stack", "primary", id.Stack.Primary},
		{"stack", "secondary", id.Stack.Secondary},
		{"stack", "data", id.Stack.Data},
		{"stack", "infra", id.Stack.Infra},
		{"stack.avoid", "items", id.Stack.Avoid.Items},
	} {
		for _, v := range items.values {
			if w := checkValue(items.section, items.field, v); w != nil {
				warnings = append(warnings, *w)
			}
		}
	}

	// Check learned entries
	for _, entry := range id.Learned.Entries {
		if w := checkValue("learned", "entries", entry); w != nil {
			warnings = append(warnings, *w)
		}
	}

	return warnings
}

func checkSection(section string, fields map[string]string) []SensitiveWarning {
	var warnings []SensitiveWarning
	for field, value := range fields {
		if w := checkValue(section, field, value); w != nil {
			warnings = append(warnings, *w)
		}
	}
	return warnings
}

func checkValue(section, field, value string) *SensitiveWarning {
	if value == "" {
		return nil
	}
	for _, pat := range sensitivePatterns {
		if pat.MatchString(value) {
			return &SensitiveWarning{
				Section: section,
				Field:   field,
				Value:   truncate(value, 20),
				Pattern: pat.String(),
			}
		}
	}
	return nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// FormatWarnings returns a human-readable string for sensitive data warnings.
func FormatWarnings(warnings []SensitiveWarning) string {
	if len(warnings) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("\nWARNING: potential secrets found in non-private sections:\n")
	for _, w := range warnings {
		b.WriteString(fmt.Sprintf("  [%s].%s = %q\n", w.Section, w.Field, w.Value))
	}
	b.WriteString("Move sensitive values to [private] to exclude them from distributed files.\n")
	return b.String()
}
