#!/usr/bin/env node

const { spawn, execSync, execFileSync } = require("child_process");
const path = require("path");
const fs = require("fs");
const os = require("os");
const https = require("https");
const http = require("http");

// Determine OS & Arch mapping
const platformMap = {
  darwin: "darwin",
  linux: "linux",
  win32: "windows",
};

const archMap = {
  x64: "amd64",
  arm64: "arm64",
};

const platform = platformMap[process.platform];
const arch = archMap[process.arch];

if (!platform || !arch) {
  console.error(
    `[autodev] Unsupported platform/architecture: ${process.platform}/${process.arch}`,
  );
  process.exit(1);
}

const ext = platform === "windows" ? "zip" : "tar.gz";
const binaryName = platform === "windows" ? "autodev.exe" : "autodev";

// Version: prefer the latest GitHub release tag; fall back to package.json
const pkgJson = require("../package.json");
const fallbackVersion = `v${pkgJson.version}`;

function getLatestReleaseTag() {
  return new Promise((resolve) => {
    const options = {
      hostname: "api.github.com",
      path: "/repos/HEETMEHTA18/autodev/releases/latest",
      headers: {
        "User-Agent": "autodev-npm-cli",
      },
      timeout: 5000,
    };

    https
      .get(options, (res) => {
        let body = "";
        res.on("data", (chunk) => (body += chunk));
        res.on("end", () => {
          try {
            const json = JSON.parse(body);
            const versionRegex = /^v?\d+\.\d+\.\d+(-[a-zA-Z0-9.]+)?$/;
            if (json.tag_name && versionRegex.test(json.tag_name)) {
              resolve(json.tag_name);
              return;
            }
          } catch (_) {}
          resolve(fallbackVersion);
        });
      })
      .on("error", () => {
        resolve(fallbackVersion);
      });
  });
}

// Resolve target paths
const binDir = __dirname;
const binaryPath = path.join(binDir, binaryName);

// Development fallback paths
const devPaths = [
  path.join(__dirname, "..", "..", "cli", "bin", binaryName),
  path.join(__dirname, "..", "..", "..", "bin", binaryName),
  path.join(__dirname, "..", "..", "..", "packages", "cli", "bin", binaryName),
];

let activeBinaryPath = binaryPath;

// Check if we are running in local dev mode and have a compiled binary
for (const devPath of devPaths) {
  if (fs.existsSync(devPath)) {
    activeBinaryPath = devPath;
    break;
  }
}

/**
 * Download a file using Node.js built-in https module (follows redirects).
 * This avoids issues with curl/wget not being available or behaving
 * differently in sandboxed npx environments.
 */
function download(url, destPath, maxRedirects = 5) {
  return new Promise((resolve, reject) => {
    if (maxRedirects <= 0) return reject(new Error("Too many redirects"));

    const client = url.startsWith("https") ? https : http;
    client
      .get(url, (res) => {
        // Follow redirects (GitHub releases return 302)
        if (
          res.statusCode >= 300 &&
          res.statusCode < 400 &&
          res.headers.location
        ) {
          return download(res.headers.location, destPath, maxRedirects - 1)
            .then(resolve)
            .catch(reject);
        }

        if (res.statusCode !== 200) {
          return reject(new Error(`HTTP ${res.statusCode} from ${url}`));
        }

        const fileStream = fs.createWriteStream(destPath);
        res.pipe(fileStream);
        fileStream.on("finish", () => {
          fileStream.close();
          resolve();
        });
        fileStream.on("error", reject);
      })
      .on("error", reject);
  });
}

async function downloadBinary() {
  let version = await getLatestReleaseTag();
  console.log(
    `\n[autodev] Native binary not found. Downloading AutoDev ${version} for ${platform}/${arch}...`,
  );

  if (!fs.existsSync(binDir)) {
    fs.mkdirSync(binDir, { recursive: true });
  }

  // Construct download URL
  const archiveName = `autodev_${platform}_${arch}`;
  const archiveFile = `${archiveName}.${ext}`;
  let url = `https://github.com/HEETMEHTA18/autodev/releases/download/${version}/${archiveFile}`;

  const tempFile = path.join(
    os.tmpdir(),
    `autodev_download_${Date.now()}.${ext}`,
  );

  // Download using Node.js built-in HTTPS (handles redirects properly)
  try {
    console.log(`[autodev] Downloading from: ${url}`);
    await download(url, tempFile);
  } catch (err) {
    const stableFallback = "v0.3.2";
    if (version !== stableFallback) {
      console.warn(
        `\n[autodev] Failed to download version ${version}: ${err.message}`,
      );
      console.warn(
        `[autodev] Falling back to last known stable release: ${stableFallback}...`,
      );
      version = stableFallback;
      url = `https://github.com/HEETMEHTA18/autodev/releases/download/${version}/${archiveFile}`;
      try {
        console.log(`[autodev] Downloading from: ${url}`);
        await download(url, tempFile);
      } catch (retryErr) {
        console.error(
          `\n[autodev] Error downloading stable release asset: ${retryErr.message}`,
        );
        console.error(`[autodev] Please verify your network connection.`);
        process.exit(1);
      }
    } else {
      console.error(
        `\n[autodev] Error downloading release asset: ${err.message}`,
      );
      console.error(`[autodev] URL: ${url}`);
      process.exit(1);
    }
  }

  try {
    // Verify the file was actually downloaded
    if (!fs.existsSync(tempFile)) {
      throw new Error("Download completed but file not found on disk.");
    }
    const stat = fs.statSync(tempFile);
    if (stat.size < 1000) {
      throw new Error(
        `Downloaded file is too small (${stat.size} bytes), likely an error page.`,
      );
    }
    console.log(
      `[autodev] Downloaded ${(stat.size / 1024 / 1024).toFixed(1)} MB`,
    );
  } catch (err) {
    console.error(`\n[autodev] Error verifying download: ${err.message}`);
    process.exit(1);
  }

  // Extract
  console.log(`[autodev] Extracting binary...`);
  try {
    if (ext === "zip") {
      if (process.platform === "win32") {
        const escapedTempFile = tempFile.replace(/'/g, "''");
        const escapedBinDir = binDir.replace(/'/g, "''");
        execSync(
          `powershell -Command "Expand-Archive -Path '${escapedTempFile}' -DestinationPath '${escapedBinDir}' -Force"`,
          { stdio: "inherit" },
        );
      } else {
        execFileSync("unzip", ["-o", tempFile, "-d", binDir], {
          stdio: "inherit",
        });
      }
    } else {
      execFileSync("tar", ["-xzf", tempFile, "-C", binDir], {
        stdio: "inherit",
      });
    }

    // Clean up temp archive
    if (fs.existsSync(tempFile)) {
      fs.unlinkSync(tempFile);
    }

    // Set execution permissions on Linux/macOS
    if (process.platform !== "win32" && fs.existsSync(binaryPath)) {
      fs.chmodSync(binaryPath, 0o755);
    }
    console.log(`[autodev] Installation successful.\n`);
  } catch (err) {
    // Clean up temp file on error too
    if (fs.existsSync(tempFile)) {
      fs.unlinkSync(tempFile);
    }
    console.error(`\n[autodev] Error extracting archive: ${err.message}`);
    process.exit(1);
  }
}

async function main() {
  // If no local dev binary is found and the packaged binary is missing, download it
  if (activeBinaryPath === binaryPath && !fs.existsSync(binaryPath)) {
    await downloadBinary();
  }

  // Forward execution asynchronously to handle long-running / interactive processes properly
  const args = process.argv.slice(2);
  const child = spawn(activeBinaryPath, args, { stdio: "inherit" });

  // Forward termination signals to the child process (critical for long-running servers / MCP)
  const signals = ["SIGINT", "SIGTERM", "SIGHUP", "SIGQUIT"];
  signals.forEach((signal) => {
    process.on(signal, () => {
      if (!child.killed) {
        child.kill(signal);
      }
    });
  });

  child.on("close", (code, signal) => {
    if (code !== null) {
      process.exit(code);
    } else if (signal) {
      // Exit with standard 128 + signal number
      const signalCodes = { SIGINT: 2, SIGTERM: 15, SIGHUP: 1, SIGQUIT: 3 };
      process.exit(128 + (signalCodes[signal] || 0));
    } else {
      process.exit(0);
    }
  });

  child.on("error", (err) => {
    console.error(`[autodev] Failed to run binary: ${err.message}`);
    process.exit(1);
  });
}

main().catch((err) => {
  console.error(`[autodev] Fatal error: ${err.message}`);
  process.exit(1);
});
