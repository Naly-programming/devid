#!/usr/bin/env node

const { execSync } = require("child_process");
const fs = require("fs");
const path = require("path");
const https = require("https");
const os = require("os");

const REPO = "Naly-programming/devid";
const BIN_DIR = path.join(__dirname, "bin");

function getPlatform() {
  const platform = os.platform();
  switch (platform) {
    case "darwin": return "darwin";
    case "linux": return "linux";
    case "win32": return "windows";
    default: throw new Error(`Unsupported platform: ${platform}`);
  }
}

function getArch() {
  const arch = os.arch();
  switch (arch) {
    case "x64": return "amd64";
    case "arm64": return "arm64";
    default: throw new Error(`Unsupported architecture: ${arch}`);
  }
}

function getLatestVersion() {
  return new Promise((resolve, reject) => {
    https.get(`https://api.github.com/repos/${REPO}/releases/latest`, {
      headers: { "User-Agent": "devid-npm-installer" }
    }, (res) => {
      let data = "";
      res.on("data", (chunk) => data += chunk);
      res.on("end", () => {
        try {
          const json = JSON.parse(data);
          resolve(json.tag_name.replace(/^v/, ""));
        } catch (e) {
          reject(new Error("Failed to parse GitHub API response"));
        }
      });
    }).on("error", reject);
  });
}

function download(url, dest) {
  return new Promise((resolve, reject) => {
    const follow = (url) => {
      https.get(url, { headers: { "User-Agent": "devid-npm-installer" } }, (res) => {
        if (res.statusCode === 302 || res.statusCode === 301) {
          return follow(res.headers.location);
        }
        if (res.statusCode !== 200) {
          return reject(new Error(`Download failed: ${res.statusCode}`));
        }
        const file = fs.createWriteStream(dest);
        res.pipe(file);
        file.on("finish", () => { file.close(); resolve(); });
      }).on("error", reject);
    };
    follow(url);
  });
}

async function main() {
  const platform = getPlatform();
  const arch = getArch();
  const version = await getLatestVersion();
  const ext = platform === "windows" ? "zip" : "tar.gz";
  const url = `https://github.com/${REPO}/releases/download/v${version}/devid_${version}_${platform}_${arch}.${ext}`;

  console.log(`Installing devid v${version} (${platform}/${arch})...`);

  fs.mkdirSync(BIN_DIR, { recursive: true });

  const tmpFile = path.join(os.tmpdir(), `devid.${ext}`);
  await download(url, tmpFile);

  if (ext === "zip") {
    execSync(`powershell -Command "Expand-Archive -Force '${tmpFile}' '${BIN_DIR}'"`, { stdio: "inherit" });
  } else {
    execSync(`tar -xzf "${tmpFile}" -C "${BIN_DIR}"`, { stdio: "inherit" });
  }

  // Make executable
  const binary = path.join(BIN_DIR, platform === "windows" ? "devid.exe" : "devid");
  if (platform !== "windows") {
    fs.chmodSync(binary, 0o755);
  }

  fs.unlinkSync(tmpFile);
  console.log(`devid v${version} installed successfully.`);
}

main().catch((err) => {
  console.error(`Failed to install devid: ${err.message}`);
  process.exit(1);
});
