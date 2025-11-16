#!/usr/bin/env node

const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');

// Detect platform and architecture
const platform = process.platform;
const arch = process.arch;

console.log(`Installing cntm for ${platform}-${arch}...`);

// Check if Go is installed
try {
  execSync('go version', { stdio: 'ignore' });
} catch (error) {
  console.error('Error: Go is not installed or not in PATH');
  console.error('Please install Go from https://golang.org/dl/');
  process.exit(1);
}

// Determine binary name based on platform
const binaryName = platform === 'win32' ? 'cntm.exe' : 'cntm';
const binaryPath = path.join(__dirname, '..', binaryName);

// Remove old binary if exists
if (fs.existsSync(binaryPath)) {
  fs.unlinkSync(binaryPath);
}

// Build the binary
try {
  console.log('Building cntm binary...');
  execSync(`go build -o ${binaryName}`, {
    cwd: path.join(__dirname, '..'),
    stdio: 'inherit',
  });

  // Make binary executable on Unix-like systems
  if (platform !== 'win32') {
    fs.chmodSync(binaryPath, 0o755);
  }

  console.log('✓ cntm installed successfully!');
  console.log(`\nRun 'cntm --help' to get started.`);
} catch (error) {
  console.error('Error: Failed to build cntm binary');
  console.error(error.message);
  process.exit(1);
}
