/*
	psmtp.go: simple sendmail program
	          2025-05-01 - prrar

	Example of config.json:

	{
    	"email": "username@gmail.com",
    	"password": "app-generated-password",
    	"smtp_host": "smtp.gmail.com",
    	"smtp_port": "587"
	}

	Usage: psmtp [options] recipient1,recipient2...
  		-c string
        	Shorthand for --config
  		-config string
        	Path to JSON config file. If empty, defaults to $HOME/.psmtp/config.json
  		-s string
	        Shorthand for --subject
  		-subject string
        	Email subject (optional)
  		-to string
        	Comma-separated recipient emails (positional override)

	Example: echo "Hello, world!" | psmtp -s "Test Message" destination@gmail.com
*/

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
)

// Config holds credentials and SMTP server details for sending email.
type Config struct {
	Email    string `json:"email"`     // Sender's email address
	Password string `json:"password"`  // SMTP auth password or app-specific password
	SMTPHost string `json:"smtp_host"` // SMTP server hostname (e.g., smtp.gmail.com)
	SMTPPort string `json:"smtp_port"` // SMTP server port (e.g., "587")
}

// parseFlags reads and validates command-line arguments.
// - First non-flag argument is treated as recipient list (comma-separated).
// - --config or -c sets config file path.
// - --subject or -s sets the email subject.
func parseFlags() (cfgPath, subject string, toList []string) {
	var toCSV string

	// Config file flags
	flag.StringVar(&cfgPath, "config", "",
		"Path to JSON config file. If empty, defaults to $HOME/.psmtp/config.json")
	flag.StringVar(&cfgPath, "c", "",
		"Shorthand for --config")

	// Subject flags
	flag.StringVar(&subject, "subject", "",
		"Email subject (optional)")
	flag.StringVar(&subject, "s", "",
		"Shorthand for --subject")

	// Recipient flag (optional; can use positional)
	flag.StringVar(&toCSV, "to", "",
		"Comma-separated recipient emails (positional override)")

	// Custom usage message
	flag.Usage = func() {
		prog := filepath.Base(os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(),
			"Usage: %s [options] recipient1,recipient2...\n", prog)
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "Example: echo \"Hello, world!\" | %s -s \"Test Message\" destination@gmail.com\n", prog)
	}

	flag.Parse()

	// If no --to, take first positional argument as recipient list
	if toCSV == "" && len(flag.Args()) > 0 {
		toCSV = flag.Args()[0]
	}

	// toCSV is now required
	if toCSV == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Default config path if none provided
	if cfgPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("could not determine home directory: %v", err)
		}
		cfgPath = filepath.Join(home, ".psmtp", "config.json")
	}

	// Build toList slice
	for _, addr := range strings.Split(toCSV, ",") {
		if trimmed := strings.TrimSpace(addr); trimmed != "" {
			toList = append(toList, trimmed)
		}
	}

	return
}

// loadConfig reads JSON configuration from `path` and validates required fields.
func loadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("reading config file: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parsing config file: %w", err)
	}
	if cfg.Email == "" || cfg.Password == "" ||
		cfg.SMTPHost == "" || cfg.SMTPPort == "" {
		return Config{}, fmt.Errorf("incomplete config: missing required fields")
	}
	return cfg, nil
}

// readBody reads the entire email body from stdin.
func readBody() ([]byte, error) {
	return io.ReadAll(os.Stdin)
}

// buildMessage constructs the raw email content with headers and body using CRLF.
func buildMessage(
	from, subject string,
	toList []string,
	body []byte,
) []byte {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "From: %s\r\n", from)
	fmt.Fprintf(&buf, "To: %s\r\n", strings.Join(toList, ", "))
	fmt.Fprintf(&buf, "Subject: %s\r\n", subject)
	buf.WriteString("\r\n")
	buf.Write(body)
	return buf.Bytes()
}

// sendEmail authenticates with SMTP and sends the message.
func sendEmail(cfg Config, toList []string, msg []byte) error {
	auth := smtp.PlainAuth("", cfg.Email, cfg.Password, cfg.SMTPHost)
	return smtp.SendMail(
		cfg.SMTPHost+":"+cfg.SMTPPort,
		auth,
		cfg.Email,
		toList,
		msg,
	)
}

func main() {
	cfgPath, subject, toList := parseFlags()

	cfg, err := loadConfig(cfgPath)
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	body, err := readBody()
	if err != nil {
		log.Fatalf("Reading body error: %v", err)
	}

	msg := buildMessage(cfg.Email, subject, toList, body)

	if err := sendEmail(cfg, toList, msg); err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}

	log.Println("Email sent successfully.")
}
