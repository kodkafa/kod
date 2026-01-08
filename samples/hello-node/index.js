const args = process.argv.slice(2);
let name = "NodeUser";

for (let i = 0; i < args.length; i++) {
    if (args[i] === "--name" && args[i + 1]) {
        name = args[i + 1];
    }
}

console.log("╔════════════════════════════════════╗");
console.log(`║  Name: ${name.padEnd(27)} ║`);
console.log("║  Hello from KODKAFA Node.js!       ║");
console.log("╚════════════════════════════════════╝");
