# Egg Carton CLI

## What This Is

A production-ready command-line tool for managing secrets in the Egg Carton managed vault service. Users download a single binary, authenticate once via OAuth, and immediately manage their secrets — no configuration, no environment variables, no infrastructure knowledge required. The Egg Carton team provides all backend infrastructure (AWS Cognito, API Gateway, Lambda, DynamoDB).

## Core Value

Zero-config secret management — download, `egg login`, and go. The CLI should feel as simple as any developer tool: one install, one auth step, and you're productive.

## Requirements

### Validated

<!-- Existing working functionality from codebase audit -->

- ✓ `egg login` — OAuth 2.0 PKCE authentication via AWS Cognito; tokens cached at `~/.eggcarton/credentials.json` — existing
- ✓ `egg lay <key> <value>` — Store a secret in the vault — existing
- ✓ `egg get [key]` — Retrieve one or all secrets — existing
- ✓ `egg break <key>` — Delete a secret — existing
- ✓ `egg hatch -- <cmd>` — Inject all secrets as environment variables and run a subprocess — existing
- ✓ Multi-platform binary builds (macOS arm64/amd64, Linux amd64, Windows amd64) — existing

### Active

<!-- Production hardening, missing features, and new capabilities -->

- [ ] Config (Cognito pool ID, client ID, domain, region, API endpoint) baked into binary at build time via ldflags — users never touch env vars or config files
- [ ] Test coverage for all critical paths: API client, config/token management, auth flow, all commands, and error paths
- [ ] User-facing documentation: command reference, quickstart, and in-binary `--help` text polished for a public audience
- [ ] `egg list` command — list all secret key names in the vault (ListEggs currently stubbed as "not implemented")
- [ ] `--role <role>` optional flag across all commands — lets AI agents and automation identify themselves from a preset list (e.g. `ci`, `developer`, `admin`, `ops`) for audit tracking on the backend
- [ ] HTTP client timeout (currently no timeout — requests can hang indefinitely)
- [ ] Deduplicate token load/validate/refresh logic across all commands into a shared helper
- [ ] Fix resource leaks: ensure `resp.Body.Close()` is deferred in all response paths
- [ ] Fix silent error suppression in login command (`existingTokens, _ :=`)
- [ ] Consistent, user-friendly error messages across all commands

### Out of Scope

- Multi-account / named profile support — single account per user is sufficient for v1
- Team/shared vault access — each user owns their secrets; collaboration is v2+
- Secret rotation scheduling — out of scope for CLI; backend concern
- Web UI — CLI only
- Cross-device token sync — users re-authenticate on each machine
- Permission-based roles — role flag is audit/identity only, not access control

## Context

The Egg Carton CLI is the public-facing interface for the Egg Carton secret vault service. The backend (Lambda functions, DynamoDB, KMS, API Gateway, Cognito user pool) lives in a separate private repository owned by the same team. This CLI is what end users and developers download and interact with.

Current state: Core secret CRUD and env injection are functional. The CLI loads configuration from a `.env` file, which is appropriate for development but must be replaced with baked-in config for public distribution. Test coverage is essentially zero (10 TODO stubs in `main_test.go`). The codebase has known bugs (no HTTP timeout, resource leaks) and structural debt (duplicated token logic across 4 command files).

Target users: developers, DevOps engineers, and AI agents/CI systems that need to securely manage secrets without embedding them in code or configuration files.

## Constraints

- **Tech stack**: Go + Cobra — must stay; existing infrastructure is built around it
- **Backend**: Separate repository; CLI only makes HTTP calls to existing Lambda endpoints — no backend changes in scope
- **Distribution**: Single binary, no installer, no runtime dependencies
- **Platforms**: macOS (arm64, amd64), Linux (amd64), Windows (amd64)
- **Auth**: AWS Cognito PKCE flow must remain; no client secrets in binary

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Bake config into binary via ldflags at build time | End users must not need to set env vars; PKCE means no secret is required in binary — only public identifiers (pool ID, client ID, domain, endpoint) | — Pending |
| Role flag is identity metadata only, not access control | Access control lives on the backend; CLI role flag is for audit trail (who/what made this call) | — Pending |
| Preset role list hardcoded in CLI | Keeps UX simple; consistent values for backend audit logs | — Pending |

---
*Last updated: 2026-02-27 after initialization*
