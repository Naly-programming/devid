#!/usr/bin/env node

const { execFileSync } = require("child_process");
const path = require("path");
const os = require("os");

const ext = os.platform() === "win32" ? ".exe" : "";
const binary = path.join(__dirname, `devid${ext}`);

try {
  execFileSync(binary, process.argv.slice(2), { stdio: "inherit" });
} catch (err) {
  if (err.status !== undefined) {
    process.exit(err.status);
  }
  console.error(`Failed to run devid: ${err.message}`);
  console.error("Try reinstalling: npm install -g devid-cli");
  process.exit(1);
}
