#!/usr/bin/env node
/**
 * extract-cjk.mjs (Story 8.3)
 *
 * Scan Admin source trees for CJK characters. Two scopes:
 *   --scope=fe  (default)  -> web/src (vue/js/ts/jsx/tsx)
 *   --scope=be             -> server/{api/v1,service,middleware,model/common/response} (.go)
 *   --scope=all            -> run both scopes; exit non-zero if any hit (when --lint)
 *
 * Flags:
 *   --json   machine-readable report
 *   --lint   exit 1 if any hit found
 *
 * Intentionally read-only. Use it to build the iterative conversion TODO list
 * — never have it auto-rewrite sources.
 */
import { readdir, readFile } from 'node:fs/promises'
import { join, relative, extname } from 'node:path'
import { fileURLToPath } from 'node:url'

const ADMIN_ROOT = fileURLToPath(new URL('..', import.meta.url))

const SCOPES = {
  fe: {
    root: join(ADMIN_ROOT, 'web/src'),
    exts: new Set(['.vue', '.js', '.ts', '.jsx', '.tsx']),
    allowlist: [/\/i18n\/locales\/zh-CN\./, /\/node_modules\//]
  },
  be: {
    root: join(ADMIN_ROOT, 'server'),
    // Restrict to packages that emit user-facing strings. Plugin/source trees
    // are deliberately excluded (vendor code / DB seed data).
    includePrefixes: [
      join(ADMIN_ROOT, 'server/api/v1'),
      join(ADMIN_ROOT, 'server/service'),
      join(ADMIN_ROOT, 'server/middleware'),
      join(ADMIN_ROOT, 'server/model/common/response')
    ],
    exts: new Set(['.go']),
    allowlist: [/_test\.go$/]
  }
}

const CJK_RE = /[\u4e00-\u9fff\u3000-\u303f\uff00-\uffef]/

function parseArgs() {
  const args = process.argv.slice(2)
  const opts = { scope: 'fe', json: false, lint: false }
  for (const a of args) {
    if (a.startsWith('--scope=')) opts.scope = a.slice('--scope='.length)
    else if (a === '--json') opts.json = true
    else if (a === '--lint') opts.lint = true
  }
  return opts
}

async function walk(dir, acc = []) {
  let entries
  try {
    entries = await readdir(dir, { withFileTypes: true })
  } catch {
    return acc
  }
  for (const ent of entries) {
    const abs = join(dir, ent.name)
    if (ent.isDirectory()) await walk(abs, acc)
    else acc.push(abs)
  }
  return acc
}

function isAllowlisted(path, allowlist) {
  return allowlist.some((rx) => rx.test(path))
}

function passesInclude(path, includes) {
  if (!includes || includes.length === 0) return true
  return includes.some((p) => path.startsWith(p + '/') || path === p)
}

async function scanFile(abs) {
  const content = await readFile(abs, 'utf8')
  const hits = []
  content.split(/\r?\n/).forEach((line, idx) => {
    if (CJK_RE.test(line)) hits.push({ line: idx + 1, text: line.trim() })
  })
  return hits
}

async function runScope(name, cfg) {
  const files = await walk(cfg.root)
  const report = []
  for (const abs of files) {
    if (cfg.exts && !cfg.exts.has(extname(abs))) continue
    if (isAllowlisted(abs, cfg.allowlist || [])) continue
    if (!passesInclude(abs, cfg.includePrefixes)) continue
    const hits = await scanFile(abs)
    if (hits.length > 0) {
      report.push({
        scope: name,
        file: relative(process.cwd(), abs),
        hits: hits.length,
        samples: hits.slice(0, 3)
      })
    }
  }
  return { scanned: files.length, report }
}

function printHuman(scope, { scanned, report }) {
  const totalHits = report.reduce((n, r) => n + r.hits, 0)
  console.log(`[${scope}] scanned ${scanned} files`)
  console.log(`[${scope}] files with CJK: ${report.length}, total CJK lines: ${totalHits}`)
  for (const entry of report.slice(0, 15)) {
    console.log(`  ${entry.file} (${entry.hits})`)
    for (const s of entry.samples) {
      console.log(`     L${s.line}: ${s.text.slice(0, 120)}`)
    }
  }
  if (report.length > 15) console.log(`  ... and ${report.length - 15} more files`)
  console.log('')
}

async function main() {
  const opts = parseArgs()
  const scopes =
    opts.scope === 'all' ? ['fe', 'be'] : opts.scope === 'be' ? ['be'] : ['fe']

  let anyHits = false
  const jsonReport = {}
  for (const s of scopes) {
    const cfg = SCOPES[s]
    if (!cfg) {
      console.error(`unknown scope: ${s}`)
      process.exit(2)
    }
    const result = await runScope(s, cfg)
    jsonReport[s] = result
    if (result.report.length > 0) anyHits = true
    if (!opts.json) printHuman(s, result)
  }

  if (opts.json) process.stdout.write(JSON.stringify(jsonReport, null, 2) + '\n')

  if (opts.lint && anyHits) {
    console.error(`admin-i18n-lint: CJK characters found (scope=${opts.scope}). See above.`)
    process.exit(1)
  }
}

main().catch((err) => {
  console.error(err)
  process.exit(2)
})
