const {spawnSync} = require('node:child_process');
const result = spawnSync('yarn', ['audit', '--json'], {encoding:'utf8'});
let summary;
for (const line of result.stdout.split('\n')) {
  try { const row=JSON.parse(line); if(row.type==='auditSummary') summary=row.data.vulnerabilities; } catch {}
}
if (!summary) { console.error('Dependency scan unavailable'); process.exit(1); }
console.log(JSON.stringify(summary));
process.exit(Object.values(summary).some(count => count > 0) ? 1 : 0);
