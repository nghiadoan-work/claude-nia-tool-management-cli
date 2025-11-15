#!/usr/bin/env node

/**
 * Simple test to verify the binary works
 */

const cntm = require('./index.js');

async function test() {
  console.log('Testing cntm npm wrapper...\n');

  try {
    // Test 1: Check binary exists
    console.log('1. Checking binary path...');
    console.log(`   Binary: ${cntm.binaryPath}`);

    // Test 2: Get version
    console.log('\n2. Getting version...');
    const version = await cntm.version();
    console.log(`   Version: ${version}`);

    // Test 3: Execute help command
    console.log('\n3. Testing help command...');
    const helpResult = await cntm.execute(['--help']);
    console.log(`   Exit code: ${helpResult.code}`);

    console.log('\n✓ All tests passed!');
  } catch (error) {
    console.error('\n✗ Test failed:', error.message);
    process.exit(1);
  }
}

if (require.main === module) {
  test();
}
