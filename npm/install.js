#!/usr/bin/env node

/**
 * Installation script for cntm npm wrapper
 * Downloads the appropriate binary for the current platform from GitHub releases
 */

const fs = require('fs');
const path = require('path');
const https = require('https');
const { execSync } = require('child_process');

// Package info
const PACKAGE_NAME = 'cntm';
const VERSION = require('./package.json').version;
const GITHUB_REPO = 'nghiadoan-work/claude-nia-tool-management-cli';

// Platform detection
const PLATFORM_MAP = {
  darwin: 'darwin',
  linux: 'linux',
  win32: 'windows'
};

const ARCH_MAP = {
  x64: 'amd64',
  arm64: 'arm64'
};

function getPlatform() {
  const platform = PLATFORM_MAP[process.platform];
  if (!platform) {
    throw new Error(`Unsupported platform: ${process.platform}`);
  }
  return platform;
}

function getArch() {
  const arch = ARCH_MAP[process.arch];
  if (!arch) {
    throw new Error(`Unsupported architecture: ${process.arch}`);
  }
  return arch;
}

function getBinaryName() {
  const platform = getPlatform();
  const arch = getArch();

  if (platform === 'windows') {
    return `cntm-${platform}-${arch}.exe`;
  }
  return `cntm-${platform}-${arch}`;
}

function getDownloadURL() {
  const platform = getPlatform();
  const arch = getArch();
  const ext = platform === 'windows' ? 'zip' : 'tar.gz';
  const archiveName = `${PACKAGE_NAME}-${VERSION}-${platform}-${arch}.${ext}`;

  return `https://github.com/${GITHUB_REPO}/releases/download/v${VERSION}/${archiveName}`;
}

function downloadFile(url, dest) {
  return new Promise((resolve, reject) => {
    console.log(`Downloading ${PACKAGE_NAME} v${VERSION} for ${getPlatform()}-${getArch()}...`);
    console.log(`URL: ${url}`);

    const file = fs.createWriteStream(dest);

    https.get(url, (response) => {
      // Follow redirects
      if (response.statusCode === 302 || response.statusCode === 301) {
        return https.get(response.headers.location, (redirectResponse) => {
          if (redirectResponse.statusCode !== 200) {
            reject(new Error(`Download failed with status ${redirectResponse.statusCode}`));
            return;
          }

          const totalBytes = parseInt(redirectResponse.headers['content-length'], 10);
          let downloadedBytes = 0;

          redirectResponse.on('data', (chunk) => {
            downloadedBytes += chunk.length;
            const percent = ((downloadedBytes / totalBytes) * 100).toFixed(1);
            process.stdout.write(`\rProgress: ${percent}%`);
          });

          redirectResponse.pipe(file);

          file.on('finish', () => {
            file.close();
            console.log('\nDownload complete!');
            resolve();
          });
        }).on('error', reject);
      }

      if (response.statusCode !== 200) {
        reject(new Error(`Download failed with status ${response.statusCode}`));
        return;
      }

      const totalBytes = parseInt(response.headers['content-length'], 10);
      let downloadedBytes = 0;

      response.on('data', (chunk) => {
        downloadedBytes += chunk.length;
        const percent = ((downloadedBytes / totalBytes) * 100).toFixed(1);
        process.stdout.write(`\rProgress: ${percent}%`);
      });

      response.pipe(file);

      file.on('finish', () => {
        file.close();
        console.log('\nDownload complete!');
        resolve();
      });
    }).on('error', (err) => {
      fs.unlink(dest, () => {});
      reject(err);
    });
  });
}

function extractArchive(archivePath, destDir) {
  console.log('Extracting archive...');

  const platform = getPlatform();

  try {
    if (platform === 'windows') {
      // Use tar on Windows (available in Windows 10+)
      execSync(`tar -xf "${archivePath}" -C "${destDir}"`, { stdio: 'inherit' });
    } else {
      // Use tar on Unix systems
      execSync(`tar -xzf "${archivePath}" -C "${destDir}"`, { stdio: 'inherit' });
    }
    console.log('Extraction complete!');
  } catch (error) {
    throw new Error(`Failed to extract archive: ${error.message}`);
  }
}

function makeExecutable(filePath) {
  if (getPlatform() !== 'windows') {
    fs.chmodSync(filePath, 0o755);
  }
}

async function install() {
  try {
    const binDir = path.join(__dirname, 'bin');
    const tmpDir = path.join(__dirname, 'tmp');

    // Create directories
    if (!fs.existsSync(binDir)) {
      fs.mkdirSync(binDir, { recursive: true });
    }
    if (!fs.existsSync(tmpDir)) {
      fs.mkdirSync(tmpDir, { recursive: true });
    }

    // Download archive
    const platform = getPlatform();
    const arch = getArch();
    const ext = platform === 'windows' ? 'zip' : 'tar.gz';
    const archiveName = `${PACKAGE_NAME}-${VERSION}-${platform}-${arch}.${ext}`;
    const archivePath = path.join(tmpDir, archiveName);
    const downloadURL = getDownloadURL();

    await downloadFile(downloadURL, archivePath);

    // Extract archive
    extractArchive(archivePath, tmpDir);

    // Move binary to bin directory
    const binaryName = getBinaryName();
    const extractedBinaryPath = path.join(tmpDir, binaryName);
    const finalBinaryPath = path.join(binDir, platform === 'windows' ? 'cntm.exe' : 'cntm');

    if (!fs.existsSync(extractedBinaryPath)) {
      throw new Error(`Binary not found in archive: ${extractedBinaryPath}`);
    }

    fs.copyFileSync(extractedBinaryPath, finalBinaryPath);
    makeExecutable(finalBinaryPath);

    // Cleanup
    fs.rmSync(tmpDir, { recursive: true, force: true });

    console.log(`\n✓ ${PACKAGE_NAME} v${VERSION} installed successfully!`);
    console.log(`\nRun 'npx ${PACKAGE_NAME} --help' to get started.`);

  } catch (error) {
    console.error('\n✗ Installation failed:', error.message);
    console.error('\nPlease try one of these alternatives:');
    console.error(`  1. Download directly from: https://github.com/${GITHUB_REPO}/releases`);
    console.error(`  2. Build from source: go install github.com/${GITHUB_REPO}@latest`);
    process.exit(1);
  }
}

// Run installation
if (require.main === module) {
  install();
}

module.exports = { install };
