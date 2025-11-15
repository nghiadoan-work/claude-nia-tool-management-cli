/**
 * cntm - Claude Nia Tool Management CLI
 *
 * This is an npm wrapper for the cntm Go binary.
 * For programmatic usage, spawn the binary directly.
 */

const { spawn } = require('child_process');
const path = require('path');
const fs = require('fs');

function getBinaryPath() {
  const platform = process.platform;
  const binaryName = platform === 'win32' ? 'cntm.exe' : 'cntm';
  return path.join(__dirname, 'bin', binaryName);
}

/**
 * Execute cntm with given arguments
 * @param {string[]} args - Command line arguments
 * @param {object} options - Spawn options
 * @returns {Promise<{stdout: string, stderr: string, code: number}>}
 */
function execute(args = [], options = {}) {
  return new Promise((resolve, reject) => {
    const binaryPath = getBinaryPath();

    if (!fs.existsSync(binaryPath)) {
      reject(new Error('cntm binary not found. Please reinstall the package.'));
      return;
    }

    const defaultOptions = {
      stdio: 'pipe',
      shell: false,
      ...options
    };

    const child = spawn(binaryPath, args, defaultOptions);

    let stdout = '';
    let stderr = '';

    if (child.stdout) {
      child.stdout.on('data', (data) => {
        stdout += data.toString();
      });
    }

    if (child.stderr) {
      child.stderr.on('data', (data) => {
        stderr += data.toString();
      });
    }

    child.on('error', (err) => {
      reject(err);
    });

    child.on('exit', (code) => {
      resolve({
        stdout: stdout.trim(),
        stderr: stderr.trim(),
        code: code || 0
      });
    });
  });
}

/**
 * Get cntm version
 * @returns {Promise<string>}
 */
async function version() {
  const result = await execute(['version', '--output', 'json']);
  const data = JSON.parse(result.stdout);
  return data.version;
}

/**
 * Search for tools
 * @param {string} query - Search query
 * @returns {Promise<object[]>}
 */
async function search(query) {
  const result = await execute(['search', query, '--json']);
  return JSON.parse(result.stdout);
}

/**
 * Install a tool
 * @param {string} toolName - Tool name (with optional @version)
 * @returns {Promise<void>}
 */
async function install(toolName) {
  await execute(['install', toolName]);
}

/**
 * List installed tools
 * @returns {Promise<object>}
 */
async function list() {
  const result = await execute(['list', '--json']);
  return JSON.parse(result.stdout);
}

module.exports = {
  execute,
  version,
  search,
  install,
  list,
  binaryPath: getBinaryPath()
};
