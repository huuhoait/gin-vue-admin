## GVA Email Sending Plugin
#### Developer: GIN-VUE-ADMIN Official

### Usage

#### 1. Register the plugin at the end of the Routers method in initialize/router.go using the security mode you need
    Example:
    The plugin can read values from the GVA config file or accept hardcoded values. Adding a dedicated config structure and passing its values in is recommended.
	PluginInit(PrivateGroup, email.CreateEmailPlug(
		global.GVA_CONFIG.Email.To,
		global.GVA_CONFIG.Email.From,
		global.GVA_CONFIG.Email.Host,
		global.GVA_CONFIG.Email.Secret,
		global.GVA_CONFIG.Email.Nickname,
		global.GVA_CONFIG.Email.Port,
		global.GVA_CONFIG.Email.IsSSL,
		global.GVA_CONFIG.Email.IsLoginAuth,
		))

    You can also pass literal values directly:

    PluginInit(PrivateGroup, email.CreateEmailPlug(
    "a@qq.com",
    "b@qq.com",
    "smtp.qq.com",
    "global.GVA_CONFIG.Email.Secret",
    "login-secret",
    465,
    true,
    true,
    ))

### 2. Configuration

#### 2-1 Global config struct
    // Form and Secret are usually the username and password

    type Email struct {
	    To          string  // Recipient: multiple addresses separated by commas, e.g. a@qq.com,b@qq.com. In production pass this as a parameter. This value is mainly used for error-monitoring emails
	    From        string  // Sender: the mailbox you send from
	    Host        string  // Server address, e.g. smtp.qq.com. Check your mailbox's SMTP docs
	    Secret      string  // Secret used to log in. Prefer an SMTP application key over the mailbox password
	    Nickname    string  // Sender nickname, optional
	    Port        int     // Port. Check your mailbox's SMTP docs; typically 465
	    IsSSL       bool    // Whether to enable SSL
	    IsLoginAuth bool    // Whether to use LOGIN auth (required by IBM, Microsoft mail servers, etc.)
    }
#### 2-2 Request struct
    // Form and Secret are usually the username and password

    type Email struct {
        To      string `json:"to"`      // Recipient
        Subject string `json:"subject"` // Email subject
        Body    string `json:"body"`    // Email body
    }


### 3. API helpers

    utils.EmailTest(subject, body) - Send a test email
    Example: utils.EmailTest("Test subject", "Test body")
    utils.ErrorToEmail(subject, body) - Error monitoring
    Example: utils.ErrorToEmail("Test subject", "Test body")
    utils.Email(recipients, subject, body) - Send to one or more recipients (comma-separated)
    Example: utils.Email("a.qq.com,b.qq.com", "Test subject", "Test body")

### 4. Directly callable endpoints

    Test endpoint: /email/emailTest [post] (swagger-enabled)

    Send endpoint: /email/emailSend [post] (swagger-enabled)
    Request:
    type Email struct {
        To      string `json:"to"`      // Recipient
        Subject string `json:"subject"` // Email subject
        Body    string `json:"body"`    // Email body
    }
   
