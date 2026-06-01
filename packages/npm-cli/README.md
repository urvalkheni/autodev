# AutoDev NPM CLI Wrapper

This package provides a lightweight Node.js wrapper around the native Go compiled binaries for AutoDev. It allows developers to run AutoDev seamlessly using standard Node package execution engines like `npx` or `pnpm dlx`, or by installing it globally via npm.

## Usage

### Run via npx (No Installation Needed)

```bash
npx @heetmehta18/autodev --help
```

### Install Globally

```bash
npm install -g @heetmehta18/autodev
autodev --help
```

### Run Model Context Protocol (MCP) Server (Continuous Server)

The wrapper fully supports long-running stream-based commands like the AutoDev MCP server. To start the server and connect your AI coding tools (e.g. Claude Desktop, Cursor) to your local environment:

```bash
# Run directly via npx
npx @heetmehta18/autodev mcp start

# Run globally (if installed)
autodev mcp start
```

## How It Works

1. **Platform Detection:** The JavaScript wrapper reads `process.platform` and `process.arch` to map the user's platform to the target release asset names (e.g. `linux/amd64`, `windows/arm64`, `darwin/arm64`).
2. **Dynamic Download:** If the native binary is not yet cached locally in this package's `bin/` directory, the wrapper automatically downloads the correct compressed release (`.tar.gz` or `.zip`) directly from the corresponding GitHub Release tag (matching the `package.json` version).
3. **Execution Delegation:** The wrapper spawns the native binary as a subprocess, forwarding all arguments, stdio streams, and exit codes. Future runs bypass the download step entirely for instant execution.
4. **Development DX Mode:** During local development, if a compiled binary is found under `packages/cli/bin/autodev`, the wrapper forwards execution directly to the local dev build, bypassing remote GitHub requests.

## Publishing

To publish updates to the npm registry:

1. Ensure the package version matches the GitHub release tag:
   ```bash
   pnpm --filter=@heetmehta18/autodev version <new-version>
   ```
2. Publish to npm:
   ```bash
   pnpm --filter=@heetmehta18/autodev publish
   ```
