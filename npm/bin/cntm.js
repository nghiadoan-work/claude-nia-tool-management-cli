#!/usr/bin/env node

/**
 * Wrapper script for cntm binary
 * Executes the appropriate platform-specific binary
 */

const { spawn } = require('child_process');
const path = require('path');
const fs = require('fs');

function getBinaryPath() {
  const platform = process.platform;
  const binaryName = platform === 'win32' ? 'cntm.exe' : 'cntm';
  const binaryPath = path.join(__dirname, binaryName);

  if (!fs.existsSync(binaryPath)) {
    console.error('Error: cntm binary not found.');
    console.error('Please run: npm install cntm');
    process.exit(1);
  }

  return binaryPath;
}

function run() {
  const binaryPath = getBinaryPath();
  const args = process.argv.slice(2);

  // Spawn the binary with all arguments
  const child = spawn(binaryPath, args, {
    stdio: 'inherit',
    shell: false
  });

  // Forward exit code
  child.on('exit', (code) => {
    process.exit(code || 0);
  });

  // Handle errors
  child.on('error', (err) => {
    console.error('Failed to execute cntm:', err.message);
    process.exit(1);
  });
}

run();
