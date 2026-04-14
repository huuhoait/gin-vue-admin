#!/usr/bin/env node
/**
 * translate-go-cjk.mjs — Bulk CJK → English for gin-vue-admin Go sources
 *
 * Covers:
 *   - // line comments (ignores // inside strings on the same line)
 *   - gorm:"...;comment:...;..." values after comment:
 *   - interpreted strings "..." that contain CJK (errors, messages)
 *   - raw string literals `...` including multiline (MCP descriptions, etc.)
 *
 * Translation (real English text):
 *   - engine=libre (default): set LIBRETRANSLATE_URL to a LibreTranslate base URL, POST /translate
 *   - engine=http: TRANSLATE_HTTP_URL expects POST { texts: string[] } → { translations: string[] }
 *
 * Placeholder only (NOT translation — no English meaning):
 *   - engine=copy: replaces each CJK run with the literal token EN_ for offline / lint-only use.
 *     Do not use copy if you want readable English; use libre or http instead.
 *
 * Examples (real translation; hosted API usually needs a key):
 *   export LIBRETRANSLATE_URL=https://libretranslate.com
 *   export LIBRETRANSLATE_API_KEY=...   # from https://portal.libretranslate.com
 *   node services/admin/scripts/translate-go-cjk.mjs --dry-run
 *   node services/admin/scripts/translate-go-cjk.mjs --write
 *
 * Options:
 *   --dry-run       print files that would change (default)
 *   --write         overwrite files
 *   --include-tests include *_test.go
 *   --engine=libre|http|copy   (default: libre; copy = each contiguous CJK run → one EN_, no network)
 *   --root=PATH     scan tree (repeatable). Default: services/admin/server
 *   --max-files=N   limit for testing
 *   --delay-ms=N    pause between HTTP translations (default 80)
 */
import { readdir, readFile, writeFile } from 'node:fs/promises'
import { join, relative, extname } from 'node:path'
import { fileURLToPath } from 'node:url'

const ADMIN_ROOT = fileURLToPath(new URL('..', import.meta.url))

const CJK_RE = /[\u4e00-\u9fff\u3040-\u30ff\uac00-\ud7af\u3000-\u303f\uff00-\uffef]/
/** One or more CJK codepoints → single placeholder (avoids EN_EN_EN_… per character). */
const CJK_RUN_RE = /[\u4e00-\u9fff\u3040-\u30ff\uac00-\ud7af\u3000-\u303f\uff00-\uffef]+/g
/** Collapse repeated copy-mode placeholders from older runs (EN_EN_EN_ → EN_). */
const EN_PLACEHOLDER_RUN_RE = /(?:EN_)+/g

const cache = new Map()

function parseArgs() {
  const args = process.argv.slice(2)
  const opts = {
    dryRun: true,
    write: false,
    engine: process.env.TRANSLATE_ENGINE || 'libre',
    libreUrl: (process.env.LIBRETRANSLATE_URL || '').replace(/\/$/, ''),
    httpUrl: process.env.TRANSLATE_HTTP_URL || '',
    apiKey: process.env.LIBRETRANSLATE_API_KEY || process.env.TRANSLATE_API_KEY || '',
    skipTests: true,
    roots: [],
    maxFiles: 0,
    delayMs: 80
  }
  for (let i = 0; i < args.length; i++) {
    const a = args[i]
    if (a === '--write') {
      opts.write = true
      opts.dryRun = false
    } else if (a === '--dry-run') opts.dryRun = true
    else if (a === '--include-tests') opts.skipTests = false
    else if (a === '--engine') {
      opts.engine = args[++i] || opts.engine
    } else if (a.startsWith('--engine=')) opts.engine = a.slice('--engine='.length)
    else if (a.startsWith('--libre-url=')) opts.libreUrl = a.slice('--libre-url='.length).replace(/\/$/, '')
    else if (a.startsWith('--http-url=')) opts.httpUrl = a.slice('--http-url='.length)
    else if (a.startsWith('--root=')) opts.roots.push(a.slice('--root='.length))
    else if (a.startsWith('--max-files='))
      opts.maxFiles = Number(a.slice('--max-files='.length)) || 0
    else if (a.startsWith('--delay-ms='))
      opts.delayMs = Number(a.slice('--delay-ms='.length)) || 0
  }
  if (!opts.roots.length) opts.roots = [join(ADMIN_ROOT, 'server')]
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
    if (ent.isDirectory()) {
      if (ent.name === 'vendor' || ent.name === 'node_modules') continue
      await walk(abs, acc)
    } else acc.push(abs)
  }
  return acc
}

function sleep(ms) {
  return new Promise((r) => setTimeout(r, ms))
}

/** Returns { code, comment } if // starts a line comment; null if none */
function splitLineComment(line) {
  let i = 0
  let inStr = false
  let strQuote = ''
  let escape = false
  let rawTick = false

  while (i < line.length) {
    const c = line[i]

    if (rawTick) {
      if (c === '`') rawTick = false
      i++
      continue
    }

    if (!inStr) {
      if (c === '`') {
        rawTick = true
        i++
        continue
      }
      if (c === '"' || c === "'") {
        inStr = true
        strQuote = c
        i++
        continue
      }
      if (c === '/' && line[i + 1] === '/') {
        return { code: line.slice(0, i), comment: line.slice(i + 2) }
      }
      i++
      continue
    }

    if (escape) {
      escape = false
      i++
      continue
    }
    if (c === '\\' && strQuote !== '`') {
      escape = true
      i++
      continue
    }
    if (c === strQuote) {
      inStr = false
      strQuote = ''
    }
    i++
  }
  return null
}

async function translateLibre(text, opts) {
  const url = `${opts.libreUrl}/translate`
  const body = { q: text, source: 'auto', target: 'en', format: 'text' }
  if (opts.apiKey) body.api_key = opts.apiKey
  const res = await fetch(url, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body)
  })
  if (!res.ok) {
    const txt = await res.text()
    let hint = ''
    if (res.status === 400 && /api key|Visit https:\/\/portal/i.test(txt)) {
      hint =
        '\nHint: hosted LibreTranslate requires LIBRETRANSLATE_API_KEY (https://portal.libretranslate.com). ' +
        'Or point LIBRETRANSLATE_URL at a self-hosted instance with no key.'
    }
    throw new Error(`LibreTranslate ${res.status}: ${txt.slice(0, 300)}${hint}`)
  }
  const data = await res.json()
  if (typeof data.translatedText !== 'string')
    throw new Error('LibreTranslate: invalid response (no translatedText)')
  return data.translatedText
}

async function translateHttpBatch(texts, opts) {
  const res = await fetch(opts.httpUrl, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ texts })
  })
  if (!res.ok) throw new Error(`HTTP ${res.status}: ${(await res.text()).slice(0, 300)}`)
  const data = await res.json()
  if (!Array.isArray(data.translations) || data.translations.length !== texts.length) {
    throw new Error('Expected { translations: string[] } with same length as input')
  }
  return data.translations
}

async function translateText(text, opts) {
  if (!CJK_RE.test(text)) return text
  if (cache.has(text)) return cache.get(text)

  let out = text
  if (opts.engine === 'copy') {
    out = text.replace(CJK_RUN_RE, 'EN_').replace(EN_PLACEHOLDER_RUN_RE, 'EN_')
  } else if (opts.engine === 'libre') {
    if (!opts.libreUrl) {
      throw new Error('Set LIBRETRANSLATE_URL (base URL of LibreTranslate)')
    }
    out = await translateLibre(text, opts)
  } else if (opts.engine === 'http') {
    if (!opts.httpUrl) throw new Error('Set TRANSLATE_HTTP_URL')
    const [t] = await translateHttpBatch([text], opts)
    out = t
  } else {
    throw new Error(`Unknown --engine=${opts.engine}`)
  }

  if (opts.delayMs) await sleep(opts.delayMs)
  cache.set(text, out)
  return out
}

/** Extract double-quoted string spans (one line, Go interpreted strings) */
function extractQuotedStrings(line) {
  const out = []
  let i = 0
  while (i < line.length) {
    if (line[i] !== '"') {
      i++
      continue
    }
    const start = i
    i++
    let inner = ''
    let escape = false
    while (i < line.length) {
      const c = line[i]
      if (escape) {
        inner += '\\' + c
        escape = false
        i++
        continue
      }
      if (c === '\\') {
        inner += c
        escape = true
        i++
        continue
      }
      if (c === '"') {
        i++
        if (CJK_RE.test(inner)) out.push({ start, end: i, inner })
        break
      }
      inner += c
      i++
    }
  }
  return out
}

function rebuildLineWithQuotedReplacements(line, replacements) {
  if (replacements.length === 0) return line
  const sorted = [...replacements].sort((a, b) => b.start - a.start)
  let out = line
  for (const r of sorted) {
    out = out.slice(0, r.start) + `"${r.newInner}"` + out.slice(r.end)
  }
  return out
}

async function processLine(line, opts) {
  let work = line

  const sc = splitLineComment(work)
  if (sc && CJK_RE.test(sc.comment)) {
    const body = sc.comment.replace(/^\s*/, '')
    const tr = await translateText(body, opts)
    work = sc.code + '// ' + tr
  }

  if (CJK_RE.test(work) && work.includes('gorm:"')) {
    const innerMatch = [...work.matchAll(/gorm:"([^"]*)"/g)]
    for (const m of innerMatch) {
      const inner = m[1]
      if (!CJK_RE.test(inner)) continue
      const parts = inner.split(';')
      const rebuilt = []
      for (const p of parts) {
        const cm = p.match(/^comment:(.*)$/)
        if (cm && CJK_RE.test(cm[1])) {
          const tr = await translateText(cm[1], opts)
          rebuilt.push(`comment:${tr}`)
        } else {
          rebuilt.push(p)
        }
      }
      work = work.replace(`gorm:"${inner}"`, `gorm:"${rebuilt.join(';')}"`)
    }
  }

  if (CJK_RE.test(work) && work.includes('"')) {
    const qs = extractQuotedStrings(work)
    const repl = []
    for (const q of qs) {
      const tr = await translateText(q.inner, opts)
      repl.push({ ...q, newInner: tr })
    }
    if (repl.length) work = rebuildLineWithQuotedReplacements(work, repl)
  }

  return work
}

/** Skip `"` ... `"` with Go escape rules (no newlines inside string). */
function skipGoDoubleQuotedString(s, start) {
  let i = start + 1
  let escape = false
  while (i < s.length) {
    const ch = s[i]
    if (escape) {
      if (ch === 'u' && i + 4 < s.length) i += 5
      else if (ch === 'U' && i + 8 < s.length) i += 9
      else if (ch === 'x' && i + 2 < s.length) i += 3
      else i++
      escape = false
      continue
    }
    if (ch === '\\') {
      escape = true
      i++
      continue
    }
    if (ch === '"') return i + 1
    i++
  }
  return i
}

/** Skip `'...'` rune literal. */
function skipGoRuneLiteral(s, start) {
  let i = start + 1
  if (i >= s.length) return start + 1
  if (s[i] === '\\') {
    i++
    if (i >= s.length) return i
    const ch = s[i]
    if (ch === 'u' && i + 4 < s.length) i += 5
    else if (ch === 'U' && i + 8 < s.length) i += 9
    else if (ch === 'x' && i + 2 < s.length) i += 3
    else i++
  } else {
    const cp = s.codePointAt(i)
    i += cp > 0xffff ? 2 : 1
  }
  if (i < s.length && s[i] === "'") i++
  return i
}

/**
 * Translate CJK inside raw backtick strings only (handles multiline `...`).
 * Skips `"`, `'`, //, /* so backticks inside those are not treated as raw delimiters.
 */
async function translateRawStringsInCode(content, opts) {
  let out = ''
  let i = 0
  const n = content.length
  while (i < n) {
    const c = content[i]
    if (c === '"') {
      const end = skipGoDoubleQuotedString(content, i)
      out += content.slice(i, end)
      i = end
      continue
    }
    if (c === '`') {
      i++
      const bodyStart = i
      while (i < n && content[i] !== '`') i++
      const inner = content.slice(bodyStart, i)
      i++
      let newInner = inner
      if (CJK_RE.test(inner)) newInner = await translateText(inner, opts)
      out += '`' + newInner + '`'
      continue
    }
    if (c === "'") {
      const end = skipGoRuneLiteral(content, i)
      out += content.slice(i, end)
      i = end
      continue
    }
    if (c === '/' && i + 1 < n && content[i + 1] === '/') {
      const lineStart = i
      i += 2
      while (i < n && content[i] !== '\n' && content[i] !== '\r') i++
      out += content.slice(lineStart, i)
      continue
    }
    if (c === '/' && i + 1 < n && content[i + 1] === '*') {
      const blockStart = i
      i += 2
      while (i < n - 1 && !(content[i] === '*' && content[i + 1] === '/')) i++
      i = Math.min(i + 2, n)
      out += content.slice(blockStart, i)
      continue
    }
    out += c
    i++
  }
  return out
}

async function processFile(abs, opts) {
  const first = await readFile(abs, 'utf8')
  let work = await translateRawStringsInCode(first, opts)
  const lines = work.split(/\r?\n/)
  const out = []
  for (const line of lines) {
    if (!CJK_RE.test(line)) {
      out.push(line)
      continue
    }
    out.push(await processLine(line, opts))
  }
  const nl = first.includes('\r\n') ? '\r\n' : '\n'
  const text = out.join(nl)
  return first.endsWith('\n') ? text + (text.endsWith(nl) ? '' : nl) : text
}

async function main() {
  const opts = parseArgs()
  if (opts.engine === 'copy') {
    console.warn(
      'translate-go-cjk: --engine copy only inserts EN_ placeholders (not English). ' +
        'For real translation, omit --engine copy and set LIBRETRANSLATE_URL, then run --write.\n'
    )
  }
  if (opts.engine === 'libre' && !opts.libreUrl) {
    console.error(
      'translate-go-cjk: set LIBRETRANSLATE_URL to the LibreTranslate server base URL.\n' +
        'Example: export LIBRETRANSLATE_URL=https://libretranslate.com\n' +
        'Hosted instances often require: export LIBRETRANSLATE_API_KEY=... (see https://portal.libretranslate.com)\n' +
        'Offline placeholder mode (not translation): --engine copy'
    )
    process.exit(2)
  }

  let files = []
  for (const root of opts.roots) {
    files.push(...(await walk(root)))
  }
  files = files.filter((f) => extname(f) === '.go')
  if (opts.skipTests) files = files.filter((f) => !/_test\.go$/i.test(f))
  if (opts.maxFiles > 0) files = files.slice(0, opts.maxFiles)

  let changed = 0
  for (const abs of files) {
    const before = await readFile(abs, 'utf8')
    if (!CJK_RE.test(before)) continue
    const after = await processFile(abs, opts)
    if (after === before) continue

    changed++
    const rel = relative(process.cwd(), abs)
    if (opts.dryRun) {
      console.log(`[would change] ${rel}`)
    } else {
      await writeFile(abs, after, 'utf8')
      console.log(`[written] ${rel}`)
    }
  }

  console.log(
    `\nDone. files ${opts.dryRun ? 'to change' : 'changed'}: ${changed}, cached translations: ${cache.size}`
  )
  if (opts.engine === 'copy' && changed > 0) {
    console.warn(
      'translate-go-cjk: copy mode does not translate. To replace EN_ with English, run with LIBRETRANSLATE_URL and default engine (libre), or fix strings manually.'
    )
  }
}

main().catch((e) => {
  console.error(e)
  process.exit(2)
})
