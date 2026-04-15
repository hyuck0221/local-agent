#!/usr/bin/env node
// Shim that runs the Go binary downloaded by postinstall.
const { spawnSync } = require('node:child_process');
const { existsSync } = require('node:fs');
const path = require('node:path');

const binName = process.platform === 'win32' ? 'local-agent.exe' : 'local-agent';
const binPath = path.join(__dirname, '..', 'vendor', binName);

if (!existsSync(binPath)) {
  console.error('local-agent binary missing. Re-run `npm install local-agent`.');
  process.exit(1);
}

const result = spawnSync(binPath, process.argv.slice(2), { stdio: 'inherit' });
process.exit(result.status ?? 1);
