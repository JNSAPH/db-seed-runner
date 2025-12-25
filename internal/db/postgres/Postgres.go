package postgres

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

type Credentials struct {
	Host     string
	Port     int
	User     string
	Password string
}

type Postgres struct {
	credentaisls *Credentials
	db           *sql.DB
}

func NewPostgres(c *Credentials) *Postgres {
	return &Postgres{credentaisls: c}
}

func (p *Postgres) Connect() error {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		p.credentaisls.Host, p.credentaisls.Port, p.credentaisls.User, p.credentaisls.Password)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open db connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping db: %w", err)
	}

	p.db = db
	return nil
}

func (p *Postgres) Close() error {
	if p.db != nil {
		return p.db.Close()
	}
	return nil
}

func (p *Postgres) RunScript(filePath string) error {
	if p.db == nil {
		return fmt.Errorf("database not connected")
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	statements := splitSQLStatements(string(content))
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		_, err = p.db.Exec(stmt)
		if err != nil {
			return fmt.Errorf("failed to execute statement: %w", err)
		}
	}

	return nil
}

/**
 * Splits SQL script into individual statements.
 * Handles semicolons within single quotes and dollar-quoted strings.
 */
// splitSQLStatements splits a SQL script into individual statements.
// It implements a state machine to correctly handle semicolons within:
// - Single-quoted strings ('string')
// - Double-quoted identifiers ("identifier")
// - Dollar-quoted strings ($tag$string$tag$)
// - Line comments (-- comment)
// - Block comments (/* comment */)
func splitSQLStatements(sqlScript string) []string {
	var stmts []string
	var buf strings.Builder

	// State flags
	inString := false       // Inside '...'
	inQuotedID := false     // Inside "..."
	inDollar := false       // Inside $tag$...$tag$
	inLineComment := false  // Inside -- ...
	inBlockComment := false // Inside /* ... */

	dollarTag := "" // The specific tag for the current dollar-quoted string

	n := len(sqlScript)

	for i := 0; i < n; i++ {
		char := sqlScript[i]
		var nextChar byte
		if i+1 < n {
			nextChar = sqlScript[i+1]
		}

		// 1. Handle Line Comments (-- ...)
		if inLineComment {
			if char == '\n' {
				inLineComment = false
			}
			buf.WriteByte(char)
			continue
		}

		// 2. Handle Block Comments (/* ... */)
		if inBlockComment {
			if char == '*' && nextChar == '/' {
				inBlockComment = false
				buf.WriteByte(char)
				buf.WriteByte(nextChar)
				i++ // Skip next char
				continue
			}
			buf.WriteByte(char)
			continue
		}

		// 3. Handle Single-Quoted Strings ('...')
		if inString {
			if char == '\'' {
				// Check for escaped quote ('')
				if nextChar == '\'' {
					buf.WriteByte(char)
					buf.WriteByte(nextChar)
					i++ // Skip next char
					continue
				}
				inString = false
			}
			buf.WriteByte(char)
			continue
		}

		// 4. Handle Double-Quoted Identifiers ("...")
		if inQuotedID {
			if char == '"' {
				// Check for escaped quote ("")
				if nextChar == '"' {
					buf.WriteByte(char)
					buf.WriteByte(nextChar)
					i++ // Skip next char
					continue
				}
				inQuotedID = false
			}
			buf.WriteByte(char)
			continue
		}

		// 5. Handle Dollar-Quoted Strings ($tag$ ... $tag$)
		if inDollar {
			if char == '$' {
				// Check if this is the closing tag
				if strings.HasPrefix(sqlScript[i:], dollarTag) {
					inDollar = false
					buf.WriteString(dollarTag)
					i += len(dollarTag) - 1 // Skip the rest of the tag
					continue
				}
			}
			buf.WriteByte(char)
			continue
		}

		// --- State Transitions (Entering a state) ---

		// Start Line Comment
		if char == '-' && nextChar == '-' {
			inLineComment = true
			buf.WriteByte(char)
			buf.WriteByte(nextChar)
			i++
			continue
		}

		// Start Block Comment
		if char == '/' && nextChar == '*' {
			inBlockComment = true
			buf.WriteByte(char)
			buf.WriteByte(nextChar)
			i++
			continue
		}

		// Start String
		if char == '\'' {
			inString = true
			buf.WriteByte(char)
			continue
		}

		// Start Quoted Identifier
		if char == '"' {
			inQuotedID = true
			buf.WriteByte(char)
			continue
		}

		// Start Dollar Quote
		if char == '$' {
			// Look ahead to find the end of the tag
			end := strings.IndexByte(sqlScript[i+1:], '$')
			if end != -1 {
				tag := sqlScript[i : i+1+end+1]
				// Validate tag format (letters, numbers, underscore)
				isValidTag := true
				for _, c := range tag[1 : len(tag)-1] {
					if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_') {
						isValidTag = false
						break
					}
				}

				if isValidTag {
					inDollar = true
					dollarTag = tag
					buf.WriteString(tag)
					i += len(tag) - 1
					continue
				}
			}
		}

		// --- Statement Split ---

		// If we are not in any special state, a semicolon marks the end of a statement
		if char == ';' {
			stmts = append(stmts, buf.String())
			buf.Reset()
			continue
		}

		// Regular character
		buf.WriteByte(char)
	}

	// Append any remaining text as the last statement
	if buf.Len() > 0 {
		if strings.TrimSpace(buf.String()) != "" {
			stmts = append(stmts, buf.String())
		}
	}
	return stmts
}
