# psmtp

A simple command-line tool for sending emails by piping message bodies through stdin.

**Summary:**  
psmtp requires Go 1.16 or later for module support and modern I/O APIs. It uses Gmail App Passwords (with 2-Step Verification enabled) to authenticate securely over SMTP. Configuration is stored in a JSON file—by default at `$HOME/.psmtp/config.json`—and you can build it with `go build` or install via `go install …@latest`. The project is released under the MIT License.

---

## Prerequisites

- **Go 1.16+** installed on your machine  
- A **Gmail App Password** (requires 2-Step Verification enabled)  
- **2-Step Verification** enabled on your Google account  

---

## Configuration

Create a JSON file, for example `~/.psmtp/config.json`:

```json
{
  "email":    "prrar83@gmail.com",
  "password": "your_app_password",
  "smtp_host":"smtp.gmail.com",
  "smtp_port":"587"
}

Or specify a custom path with `-c/--config`.

## Usage

```shell
# Send a one-line message
echo "Hello, world!" \
  | psmtp -s "Test Message" alice@example.com

# Send multi-recipient
cat body.txt \
  | psmtp bob@example.com,charlie@example.com
````

Options:
* `-c`, `--config` Path to JSON config file (default: `$HOME/.psmtp/config.json`)
* `-s`, `--subject` Email subject (optional)
* `recipient1,recipient2...` Positional comma-separated recipients

## Generating a Gmail App Password

1. **Enable 2-Step Verification**  
   Visit: https://support.google.com/accounts/answer/185839?hl=en

2. **Create an App Password**  
   1. Go to: https://myaccount.google.com/apppasswords  
   2. Sign in if prompted  
   3. Under **Select app**, choose **Other (Custom name)**  
   4. Enter `psmtp` and click **Generate**  
   5. Copy the 16-character password into your `config.json`
