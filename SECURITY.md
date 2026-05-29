# Security Policy

## Supported Versions

| Version | Supported |
|---------|-----------|
| 0.x.x   | ✅ Active  |

## Reporting a Vulnerability

**Please do NOT report security vulnerabilities through public GitHub Issues.**

Instead, please send an email to **security@autodev.dev** with:

1. Description of the vulnerability
2. Steps to reproduce
3. Potential impact
4. Suggested fix (if any)

You will receive a response within **48 hours**, and we aim to release a patch within **7 days** for critical issues.

## Security Best Practices for Contributors

- Never hardcode credentials, API keys, or secrets
- Use environment variables for sensitive configuration
- Validate and sanitize all user input
- Prefer HTTPS for all network requests
- Pin dependency versions in production builds
- Follow the principle of least privilege for file system operations

## Disclosure Policy

We follow responsible disclosure. After a fix is released, we will:

1. Credit the reporter (unless they prefer anonymity)
2. Publish a security advisory on GitHub
3. Update the CHANGELOG

Thank you for helping keep AutoDev secure! 🔒
