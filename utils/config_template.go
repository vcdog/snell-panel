/*
 * @Author: Vincent Yang
 * @Date: 2025-02-02 18:06:00
 * @LastEditors: Vincent Yang
 * @LastEditTime: 2025-02-02 18:06:00
 * @FilePath: /snell-panel/utils/config_template.go
 * @Telegram: https://t.me/missuo
 * @GitHub: https://github.com/missuo
 *
 * Copyright © 2025 by Vincent, All Rights Reserved.
 */

package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LoadConfigTemplate loads a configuration template from the templates directory
func LoadConfigTemplate(templateName string) (string, error) {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %v", err)
	}

	// Construct the template file path
	templatePath := filepath.Join(cwd, "templates", templateName+".conf")

	// Read the template file
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template %s: %v", templateName, err)
	}

	return string(content), nil
}

// GenerateConfigFromTemplate generates a configuration file from a template
// by replacing placeholders with actual values
func GenerateConfigFromTemplate(template string, nodes []string, nodeNames []string) string {
	// Join nodes and node names
	nodesSection := strings.Join(nodes, "\n")
	nodeNamesStr := strings.Join(nodeNames, ", ")

	// Replace placeholders
	config := strings.ReplaceAll(template, "{{NODES}}", nodesSection)
	config = strings.ReplaceAll(config, "{{NODE_NAMES}}", nodeNamesStr)

	return config
}

// GetTemplateNameByRuleSet returns the appropriate template name based on ruleSet
func GetTemplateNameByRuleSet(ruleSet string, customRules []string) string {
	// If custom rules are provided and not empty, determine based on them
	if len(customRules) > 0 {
		// Count the number of unique rules
		ruleCount := len(customRules)
		
		// Determine template based on rule count
		if ruleCount <= 3 {
			return "minimal"
		} else if ruleCount <= 8 {
			return "balanced"
		} else {
			return "comprehensive"
		}
	}

	// Otherwise use the preset
	switch ruleSet {
	case "minimal":
		return "minimal"
	case "balanced":
		return "balanced"
	case "comprehensive":
		return "comprehensive"
	default:
		// Default to minimal if unknown
		return "minimal"
	}
}
