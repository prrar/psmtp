# psmtp

A simple command-line tool for sending emails by piping message bodies through stdin.

---

## Usage

```shell
# Send a one-line message
echo "Hello, world!" | psmtp -s "Test Message" alice@example.com

# Send multi-recipient
cat body.txt | psmtp bob@example.com,charlie@example.com
````

|Option|Description|
|----|----|
|`-c`, `--config`|Path to JSON config file (default: `$HOME/.psmtp/config.json`)|
|`-s`, `--subject`|Email subject (optional)|
|`recipient1,recipient2...`|Comma-separated recipients|

---

## Configuration

Create a JSON file, for example `~/.psmtp/config.json`:

```json
{
  "email":    "john@gmail.com",
  "password": "your_app_password",
  "smtp_host":"smtp.gmail.com",
  "smtp_port":"587"
}
```
Or specify a custom path with `-c / --config`.

---

## Generating a Gmail App Password

1. **Enable 2-Step Verification**  
   Visit: https://support.google.com/accounts/answer/185839?hl=en

2. **Create an App Password**  
   1. Go to: https://myaccount.google.com/apppasswords  
   2. Sign in if prompted  
   3. Under **Select app**, choose **Other (Custom name)**  
   4. Enter `psmtp` and click **Generate**  
   5. Copy the 16-character password into your `config.json`

---
