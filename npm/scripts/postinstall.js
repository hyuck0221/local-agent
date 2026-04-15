#!/usr/bin/env node
// Downloads the platform-appropriate local-agent binary from GitHub Releases
// into npm/vendor/ at install time. Skips cleanly on unsupported platforms so
// `npm install` never hard-fails.
const { mkdirSync, createWriteStream, chmodSync, existsSync } = require('node:fs');
const { pipeline } = require('node:stream/promises');
const { Readable } = require('node:stream');
const path = require('node:path');
const zlib = require('node:zlib');
const tar = require('node:child_process');

const pkg = require('../package.json');
const REPO = process.env.LOCAL_AGENT_REPO || 'hyuck0221/local-agent';
const VERSION = process.env.LOCAL_AGENT_VERSION || `v${pkg.version}`;
const DOWNLOAD_BASE = process.env.LOCAL_AGENT_DOWNLOAD_BASE || '';

const platMap = { darwin: 'darwin', linux: 'linux', win32: 'windows' };
const archMap = { x64: 'amd64', arm64: 'arm64' };

const os = platMap[process.platform];
const arch = archMap[process.arch];
if (!os || !arch) {
  console.warn(`[local-agent] unsupported platform ${process.platform}/${process.arch}; skipping binary install`);
  process.exit(0);
}

const stripped = VERSION.replace(/^v/, '');
const ext = os === 'windows' ? 'zip' : 'tar.gz';
const asset = `local-agent_${stripped}_${os}_${arch}.${ext}`;
const url = DOWNLOAD_BASE
  ? `${DOWNLOAD_BASE}/${asset}`
  : `https://github.com/${REPO}/releases/download/${VERSION}/${asset}`;

const vendorDir = path.join(__dirname, '..', 'vendor');
mkdirSync(vendorDir, { recursive: true });
const binName = os === 'windows' ? 'local-agent.exe' : 'local-agent';
const binPath = path.join(vendorDir, binName);

if (existsSync(binPath)) {
  process.exit(0);
}

(async () => {
  console.log(`[local-agent] downloading ${asset}`);
  const res = await fetch(url, { redirect: 'follow' });
  if (!res.ok) throw new Error(`download failed: ${res.status} ${res.statusText}`);

  const archivePath = path.join(vendorDir, asset);
  await pipeline(Readable.fromWeb(res.body), createWriteStream(archivePath));

  if (ext === 'tar.gz') {
    tar.execFileSync('tar', ['-xzf', archivePath, '-C', vendorDir], { stdio: 'inherit' });
  } else {
    // Node 22+ has built-in zip via `unzip` fallback; rely on tar which on
    // modern Windows ships with the OS and handles zip natively.
    tar.execFileSync('tar', ['-xf', archivePath, '-C', vendorDir], { stdio: 'inherit' });
  }

  if (os !== 'windows') chmodSync(binPath, 0o755);
  console.log(`[local-agent] installed ${binPath}`);
})().catch((err) => {
  console.error(`[local-agent] postinstall failed: ${err.message}`);
  process.exit(1);
});
