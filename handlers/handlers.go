/*
 * @Author: Vincent Yang
 * @Date: 2025-05-03 04:24:49
 * @LastEditors: Vincent Yang
 * @LastEditTime: 2025-07-05 20:40:15
 * @FilePath: /snell-panel/handlers/handlers.go
 * @Telegram: https://t.me/missuo
 * @GitHub: https://github.com/missuo
 *
 * Copyright © 2025 by Vincent, All Rights Reserved.
 */

package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"snell-panel/models"
	"snell-panel/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Handlers contains the HTTP request handlers
type Handlers struct {
	DB    *sql.DB
	Token string
}

// NewHandlers creates a new Handlers instance
func NewHandlers(db *sql.DB, token string) *Handlers {
	return &Handlers{
		DB:    db,
		Token: token,
	}
}

// CorsMiddleware returns a CORS middleware configured for the API
func CorsMiddleware() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}
	return cors.New(config)
}

// AuthMiddleware returns a middleware that checks for API token or JWT
func (h *Handlers) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Check Query Token (Legacy / Static configuration)
		providedToken := c.Query("token")
		if providedToken != "" && providedToken == h.Token {
			c.Next()
			return
		}

		// 2. Check Bearer Token (JWT)
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// Use h.Token as secret key if available, otherwise default
			secretKey := []byte(h.Token)
			if len(secretKey) == 0 {
				secretKey = []byte("snell-panel-secret-key")
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return secretKey, nil
			})

			if err == nil && token.Valid {
				// Store claims like UserID if needed
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					c.Set("userID", claims["id"])
					c.Set("username", claims["username"])
				}
				c.Next()
				return
			}
		}

		c.JSON(http.StatusUnauthorized, models.ApiResponse{
			Status:  "error",
			Message: "Unauthorized",
		})
		c.Abort()
		return
	}
}

// Login handles user authentication
func (h *Handlers) Login(c *gin.Context) {
	var loginReq models.LoginRequest
	if err := c.BindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Status:  "error",
			Message: "Invalid request format",
		})
		return
	}

	var user models.User
	err := h.DB.QueryRow("SELECT id, username, password_hash FROM users WHERE username = $1", loginReq.Username).Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		// User not found
		c.JSON(http.StatusUnauthorized, models.ApiResponse{
			Status:  "error",
			Message: "Invalid username or password",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginReq.Password))
	if err != nil {
		// Password mismatch
		c.JSON(http.StatusUnauthorized, models.ApiResponse{
			Status:  "error",
			Message: "Invalid username or password",
		})
		return
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
	})

	// Use h.Token as secret key
	secretKey := []byte(h.Token)
	if len(secretKey) == 0 {
		secretKey = []byte("snell-panel-secret-key")
	}

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status:  "error",
			Message: "Failed to generate token",
		})
		return
	}

	// Don't return password hash
	user.PasswordHash = ""

	c.JSON(http.StatusOK, models.ApiResponse{
		Status:  "success",
		Message: "Login successful",
		Data: models.LoginResponse{
			Token: tokenString,
			User:  user,
		},
	})
}

// Welcome handles the root route
func (h *Handlers) Welcome(c *gin.Context) {
	c.JSON(http.StatusOK, models.ApiResponse{
		Status:  "success",
		Message: "Welcome to Snell Panel. Please use the API to manage the entries.\n https://github.com/missuo/snell-panel",
	})
}

// NotFound handles not found routes
func (h *Handlers) NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, models.ApiResponse{
		Status:  "error",
		Message: "Path not found",
	})
}

// InsertEntry handles creating a new entry
func (h *Handlers) InsertEntry(c *gin.Context) {
	var entry models.Entry
	if err := c.BindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	// Resolve domain to IP if needed and get IP information
	// Keep original domain/IP in entry.IP, only use resolved IP for getting geo info
	originalIP := entry.IP
	_, ipInfo, err := utils.GetIPInfoFromDomainOrIP(entry.IP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status:  "error",
			Message: fmt.Sprintf("Failed to resolve domain/IP or get IP info: %v", err),
		})
		return
	}

	// Update entry with IP information (but keep original domain/IP)
	entry.IP = originalIP // Keep the original domain/IP address
	entry.CountryCode = ipInfo.CountryCode
	entry.ISP = ipInfo.ISP
	entry.ASN = ipInfo.ASN
	entry.NodeID = utils.GenerateUUID()

	// Set default version if not provided
	if entry.Version == "" {
		entry.Version = "4"
	}

	// Insert entry into database
	var id int
	err = h.DB.QueryRow(`
		 INSERT INTO entries (ip, port, psk, country_code, isp, asn, node_id, node_name, version) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
		 RETURNING id`,
		entry.IP, entry.Port, entry.PSK, entry.CountryCode, entry.ISP, entry.ASN, entry.NodeID, entry.NodeName, entry.Version).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	// Set ID in entry
	entry.ID = id

	c.JSON(http.StatusCreated, models.ApiResponse{
		Status:  "success",
		Message: "Entry created successfully",
		Data:    entry,
	})
}

// DeleteEntryByIP handles deleting an entry by IP
func (h *Handlers) DeleteEntryByIP(c *gin.Context) {
	ip := c.Param("ip")

	result, err := h.DB.Exec("DELETE FROM entries WHERE ip = $1", ip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, models.ApiResponse{
			Status:  "error",
			Message: "Entry not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Status:  "success",
		Message: "Entry deleted successfully",
	})
}

// DeleteEntryByNodeID handles deleting an entry by node ID
func (h *Handlers) DeleteEntryByNodeID(c *gin.Context) {
	nodeID := c.Param("node_id")

	result, err := h.DB.Exec("DELETE FROM entries WHERE node_id = $1", nodeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, models.ApiResponse{
			Status:  "error",
			Message: "Entry not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Status:  "success",
		Message: "Entry deleted successfully",
	})
}

// QueryAllEntries handles retrieving all entries
func (h *Handlers) QueryAllEntries(c *gin.Context) {
	rows, err := h.DB.Query(`
		 SELECT id, ip, port, psk, country_code, isp, asn, node_id, node_name, version 
		 FROM entries
		 ORDER BY id
	 `)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	defer rows.Close()

	var entries []models.Entry
	for rows.Next() {
		var entry models.Entry
		if err := rows.Scan(
			&entry.ID, &entry.IP, &entry.Port, &entry.PSK,
			&entry.CountryCode, &entry.ISP, &entry.ASN,
			&entry.NodeID, &entry.NodeName, &entry.Version,
		); err != nil {
			c.JSON(http.StatusInternalServerError, models.ApiResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}
		entries = append(entries, entry)
	}

	if len(entries) == 0 {
		c.JSON(http.StatusNotFound, models.ApiResponse{
			Status:  "warning",
			Message: "No entries found",
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Status:  "success",
		Message: "Entries retrieved successfully",
		Data:    entries,
	})
}

// GetSubscription handles generating a subscription string
func (h *Handlers) GetSubscription(c *gin.Context) {
	// Check format parameter
	format := c.Query("format")
	if format == "advanced" {
		h.GetAdvancedSubscription(c)
		return
	}

	// Get the via, filter, and flag parameters from query string
	via := c.Query("via")
	filter := c.Query("filter")
	flagParam := c.Query("flag")

	// Default flag to true, set to false only if explicitly set to "false"
	showFlag := true
	if flagParam == "false" {
		showFlag = false
	}

	var query string
	var args []interface{}

	if filter != "" {
		// Filter nodes by node name containing the keyword
		query = `
			SELECT ip, port, psk, country_code, isp, asn, node_id, node_name, version 
			FROM entries
			WHERE node_name LIKE $1
			ORDER BY id
		`
		args = []interface{}{"%" + filter + "%"}
	} else {
		query = `
			SELECT ip, port, psk, country_code, isp, asn, node_id, node_name, version 
			FROM entries
			ORDER BY id
		`
	}

	rows, err := h.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	defer rows.Close()

	var subscriptionLines []string
	for rows.Next() {
		var entry models.Entry
		if err := rows.Scan(
			&entry.IP, &entry.Port, &entry.PSK,
			&entry.CountryCode, &entry.ISP, &entry.ASN,
			&entry.NodeID, &entry.NodeName, &entry.Version,
		); err != nil {
			c.JSON(http.StatusInternalServerError, models.ApiResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		emojiFlag := utils.CountryCodeToFlagEmoji(entry.CountryCode)
		nodeName := entry.NodeName
		if nodeName == "" {
			if showFlag {
				nodeName = fmt.Sprintf("%s %s AS%d %s %s",
					emojiFlag, entry.CountryCode, entry.ASN, entry.ISP, entry.NodeID)
			} else {
				nodeName = fmt.Sprintf("%s AS%d %s %s",
					entry.CountryCode, entry.ASN, entry.ISP, entry.NodeID)
			}
		} else {
			if showFlag {
				nodeName = fmt.Sprintf("%s %s", emojiFlag, entry.NodeName)
			} else {
				nodeName = entry.NodeName
			}
		}

		// Add - xxx suffix to node name when via parameter is provided
		if via != "" {
			nodeName = fmt.Sprintf("%s - %s", nodeName, via)
		}

		var line string
		if via != "" {
			// Include underlying-proxy parameter when via is specified
			line = fmt.Sprintf("%s = snell, %s, %d, psk = %s, version = %s, underlying-proxy = %s",
				nodeName, entry.IP, entry.Port, entry.PSK, entry.Version, via)
		} else {
			line = fmt.Sprintf("%s = snell, %s, %d, psk = %s, version = %s",
				nodeName, entry.IP, entry.Port, entry.PSK, entry.Version)
		}
		subscriptionLines = append(subscriptionLines, line)
	}

	if len(subscriptionLines) == 0 {
		c.JSON(http.StatusNotFound, models.ApiResponse{
			Status:  "error",
			Message: "No entries found for subscription",
		})
		return
	}

	c.String(http.StatusOK, strings.Join(subscriptionLines, "\n"))
}

// GetAdvancedSubscription handles generating a full configuration file
func (h *Handlers) GetAdvancedSubscription(c *gin.Context) {
	ruleSet := c.Query("ruleSet")
	customRulesStr := c.Query("customRules")
	
	// Fetch all nodes
	rows, err := h.DB.Query(`
		SELECT ip, port, psk, country_code, isp, asn, node_id, node_name, version 
		FROM entries
		ORDER BY id
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	defer rows.Close()

	var nodes []string
	var nodeNames []string

	for rows.Next() {
		var entry models.Entry
		if err := rows.Scan(
			&entry.IP, &entry.Port, &entry.PSK,
			&entry.CountryCode, &entry.ISP, &entry.ASN,
			&entry.NodeID, &entry.NodeName, &entry.Version,
		); err != nil {
			continue
		}

		emojiFlag := utils.CountryCodeToFlagEmoji(entry.CountryCode)
		nodeName := entry.NodeName
		if nodeName == "" {
			nodeName = fmt.Sprintf("%s %s AS%d %s %s",
				emojiFlag, entry.CountryCode, entry.ASN, entry.ISP, entry.NodeID)
		} else {
			nodeName = fmt.Sprintf("%s %s", emojiFlag, entry.NodeName)
		}
		
		// Clean node name
		nodeName = strings.TrimSpace(nodeName)
		nodeNames = append(nodeNames, nodeName)

		line := fmt.Sprintf("%s = snell, %s, %d, psk = %s, version = %s",
			nodeName, entry.IP, entry.Port, entry.PSK, entry.Version)
		nodes = append(nodes, line)
	}

	if len(nodes) == 0 {
		c.String(http.StatusOK, "# No nodes found in database")
		return
	}

	// Build the configuration
	var sb strings.Builder

	sb.WriteString("[General]\n")
	sb.WriteString("loglevel = notify\n")
	sb.WriteString("dns-server = system, 223.5.5.5, 8.8.8.8\n\n")

	sb.WriteString("[Proxy]\n")
	for _, node := range nodes {
		sb.WriteString(node + "\n")
	}
	sb.WriteString("\n")

	sb.WriteString("[Proxy Group]\n")
	// Proxy Selection Group
	sb.WriteString("Proxy = select, Auto, " + strings.Join(nodeNames, ", ") + "\n")
	// Auto Test Group
	sb.WriteString("Auto = url-test, " + strings.Join(nodeNames, ", ") + ", url=http://www.gstatic.com/generate_204, interval=600\n")
	sb.WriteString("\n")

	sb.WriteString("[Rule]\n")
	
	// Add rules based on ruleSet
	switch ruleSet {
	case "minimal":
		sb.WriteString("GEOIP,CN,DIRECT\n")
		sb.WriteString("FINAL,Proxy\n")
	case "balanced":
		sb.WriteString("# Balanced Rules\n")
		sb.WriteString("DOMAIN-SUFFIX,google.com,Proxy\n")
		sb.WriteString("DOMAIN-SUFFIX,youtube.com,Proxy\n")
		sb.WriteString("DOMAIN-SUFFIX,github.com,Proxy\n")
		sb.WriteString("DOMAIN-SUFFIX,telegram.org,Proxy\n")
		sb.WriteString("DOMAIN-KEYWORD,google,Proxy\n")
		sb.WriteString("GEOIP,CN,DIRECT\n")
		sb.WriteString("FINAL,Proxy\n")
	case "comprehensive":
		sb.WriteString("# Comprehensive Rules\n")
		sb.WriteString("DOMAIN-SUFFIX,google.com,Proxy\n")
		sb.WriteString("DOMAIN-SUFFIX,youtube.com,Proxy\n")
		sb.WriteString("DOMAIN-SUFFIX,facebook.com,Proxy\n")
		sb.WriteString("DOMAIN-SUFFIX,twitter.com,Proxy\n")
		sb.WriteString("DOMAIN-SUFFIX,instagram.com,Proxy\n")
		sb.WriteString("DOMAIN-SUFFIX,github.com,Proxy\n")
		sb.WriteString("DOMAIN-SUFFIX,telegram.org,Proxy\n")
		sb.WriteString("DOMAIN-SUFFIX,netflix.com,Proxy\n")
		sb.WriteString("DOMAIN-KEYWORD,google,Proxy\n")
		sb.WriteString("GEOIP,CN,DIRECT\n")
		sb.WriteString("FINAL,Proxy\n")
	case "custom":
		if customRulesStr != "" {
			rules := strings.Split(customRulesStr, ",")
			for _, r := range rules {
				// Simple validation or just append
				// Expecting format like "DOMAIN-SUFFIX:example.com:Proxy" from frontend or just standard string?
				// Frontend sends raw strings if json, but comma joined?
				// Implementation assumes frontend sends valid rule strings or we format them
				sb.WriteString(strings.TrimSpace(r) + "\n")
			}
		}
		sb.WriteString("FINAL,Proxy\n")
	default:
		// Default to direct if unknown, or balanced
		sb.WriteString("GEOIP,CN,DIRECT\n")
		sb.WriteString("FINAL,Proxy\n")
	}

	c.String(http.StatusOK, sb.String())
}

// ModifyNodeByNodeID handles modifying a node by its NodeID
func (h *Handlers) ModifyNodeByNodeID(c *gin.Context) {
	nodeID := c.Param("id")

	var modifyReq models.ModifyRequest
	if err := c.BindJSON(&modifyReq); err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	// Build query based on which fields are provided
	query := "UPDATE entries SET"
	var args []interface{}
	var setStatements []string
	paramIndex := 1

	if modifyReq.NodeName != "" {
		setStatements = append(setStatements, fmt.Sprintf(" node_name = $%d", paramIndex))
		args = append(args, modifyReq.NodeName)
		paramIndex++
	}

	if modifyReq.IP != "" {
		// Resolve domain to IP if needed and get IP information
		// Keep original domain/IP in database, only use resolved IP for getting geo info
		_, ipInfo, err := utils.GetIPInfoFromDomainOrIP(modifyReq.IP)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ApiResponse{
				Status:  "error",
				Message: fmt.Sprintf("Failed to resolve domain/IP or get IP info: %v", err),
			})
			return
		}

		// Use original domain/IP address (not resolved IP)
		setStatements = append(setStatements, fmt.Sprintf(" ip = $%d", paramIndex))
		args = append(args, modifyReq.IP)
		paramIndex++

		// Update geolocation info
		setStatements = append(setStatements, fmt.Sprintf(" country_code = $%d", paramIndex))
		args = append(args, ipInfo.CountryCode)
		paramIndex++

		setStatements = append(setStatements, fmt.Sprintf(" isp = $%d", paramIndex))
		args = append(args, ipInfo.ISP)
		paramIndex++

		setStatements = append(setStatements, fmt.Sprintf(" asn = $%d", paramIndex))
		args = append(args, ipInfo.ASN)
		paramIndex++
	}

	// If no fields to update, return error
	if len(setStatements) == 0 {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Status:  "error",
			Message: "No fields to update",
		})
		return
	}

	// Combine all set statements
	query += strings.Join(setStatements, ",")
	query += fmt.Sprintf(" WHERE node_id = $%d", paramIndex)
	args = append(args, nodeID)

	// Execute update
	result, err := h.DB.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, models.ApiResponse{
			Status:  "error",
			Message: "Node ID not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Status:  "success",
		Message: "Node updated successfully",
	})
}
