const { instantiate } = require('@wapc/host');
const { encode, decode } = require('@msgpack/msgpack');
const { promises: fs } = require('fs');
const path = require('path');

async function main() {
  // Read wasm off the local disk as Uint8Array
  const buffer = await fs.readFile(path.join(__dirname,'..', 'build', 'wapc-test.wasm'));

  // Instantiate a WapcHost with the bytes read off disk
  const host = await instantiate(buffer);

  // Encode the payload with MessagePack
  const payload = encode({ name: 'marco' });

  // Invoke the operation in the wasm guest
  const result = await host.invoke('sayHello', payload);

  // Decode the results using MessagePack
  const decoded = decode(result);

  // log to the console
  console.log(`Result: ${decoded}`);
}

main().catch(err => console.error(err));