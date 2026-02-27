# 🥚 EggCarton CLI

> **v0.1.0 — alpha**

A simple, secure command-line tool for managing application secrets — without the complexity or cost of AWS Secrets Manager.

EggCarton keeps your secrets encrypted end-to-end and accessible through a clean five-command interface. Authenticate once, then store, retrieve, and inject secrets directly into any process you run.

**Homepage:** [egg-carton.vercel.app](https://egg-carton.vercel.app/)

---

## Why EggCarton?

Managing secrets is painful. AWS Secrets Manager is powerful but expensive and over-engineered for most teams. `.env` files get committed, shared over Slack, and duplicated across machines.

EggCarton gives you:

- **One place** for all your secrets, accessible from any machine
- **End-to-end encryption** — secrets are encrypted before storage and decrypted only for you
- **Zero-config injection** — run any command with your secrets automatically in scope
- **Simple auth** — log in with OAuth in your browser, no API keys to manage

---

## Security

- Authentication uses **OAuth 2.0 with PKCE** — no client secrets are ever stored or transmitted
- Your access tokens are cached locally at `~/.eggcarton/credentials.json` with `0600` permissions (readable only by you)
- Tokens expire and are automatically refreshed; you won't be prompted to log in repeatedly
- All secrets are **encrypted at rest and in transit** on the backend — the server never holds plaintext values
- Secrets are scoped to your user account and inaccessible to others

---

## Installation

EggCarton CLI is currently distributed as source. You'll need [Go 1.21+](https://go.dev/dl/) installed.

```bash
git clone https://github.com/owenHochwald/egg-carton-cli.git
cd egg-carton-cli
go build -o egg .
```

Move the binary somewhere on your `PATH`:

```bash
mv egg /usr/local/bin/egg
```

Verify it works:

```bash
egg --help
```

---

## Quick Start

```bash
# 1. Authenticate — opens your browser
egg login

# 2. Store a secret
egg lay DB_HOST localhost
egg lay DB_USER admin
egg lay DB_PASS s3cr3t

# 3. Retrieve a secret
egg get DB_HOST
# → Value: localhost

# 4. Run a command with all your secrets injected
egg hatch -- go run main.go

# 5. Delete a secret you no longer need
egg break DB_PASS
```

---

## Commands

| Command | Alias | Description |
|---|---|---|
| `egg login` | — | Authenticate via OAuth (opens browser) |
| `egg lay <key> <value>` | `add` | Encrypt and store a secret |
| `egg get [key]` | — | Retrieve one secret, or list all |
| `egg hatch -- <cmd>` | `run` | Inject secrets as env vars and run a command |
| `egg break <key>` | — | Permanently delete a secret |

### `egg login`

Opens your default browser to complete authentication. Tokens are saved locally — you only need to log in once per session (or until your refresh token expires).

```bash
egg login
```

### `egg lay` / `egg add`

Store a new secret. The key can be any string; conventionally use `UPPER_SNAKE_CASE` to match environment variable conventions.

```bash
egg lay API_KEY abc123
egg add STRIPE_SECRET sk_live_...
```

### `egg get`

Retrieve a single secret by key, or omit the key to list everything in your vault.

```bash
egg get API_KEY          # prints the value for API_KEY
egg get                  # lists all secrets with keys and timestamps
```

### `egg hatch` / `egg run`

Fetches all your secrets, uppercases the keys, and injects them as environment variables into the subprocess. The process inherits your current shell environment plus your secrets — nothing leaks into the parent shell after the command finishes.

```bash
egg hatch -- node server.js
egg hatch -- ./deploy.sh
egg run -- env | grep API    # inspect what gets injected
```

### `egg break`

Permanently deletes a secret from your vault. This action is irreversible.

```bash
egg break OLD_API_KEY
```

---

## Example Workflow

Replace a `.env` file in a Node project:

```bash
# Migrate your existing .env
while IFS='=' read -r key value; do
  egg lay "$key" "$value"
done < .env

# Run your app — no .env file needed
egg hatch -- npm start

# Share nothing — teammates authenticate and pull their own secrets
```

## License

MIT — see [LICENSE](./LICENSE)

---

[egg-carton.vercel.app](https://egg-carton.vercel.app/)
