// Extend Date to format it to a string
// Month(M), day(d), hour(h), minute(m), second(s), quarter(q) can use 1-2 placeholders
// Year(y) can use 1-4 placeholders; millisecond(S) can only use 1 placeholder (1-3 digits)
// (new Date()).Format("yyyy-MM-dd hh:mm:ss.S") ==> 2006-07-02 08:09:04.423
// (new Date()).Format("yyyy-M-d h:m:s.S")      ==> 2006-7-2 8:9:4.18
// eslint-disable-next-line no-extend-native
Date.prototype.Format = function(fmt) {
  const o = {
    'M+': this.getMonth() + 1, // month
    'd+': this.getDate(), // day
    'h+': this.getHours(), // hour
    'm+': this.getMinutes(), // minute
    's+': this.getSeconds(), // second
    'q+': Math.floor((this.getMonth() + 3) / 3), // quarter
    'S': this.getMilliseconds() // millisecond
  }
  const reg = /(y+)/
  if (reg.test(fmt)) {
    const t = reg.exec(fmt)[1]
    fmt = fmt.replace(
      t,
      (this.getFullYear() + '').substring(4 - t.length)
    )
  }
  for (let k in o) {
    const regx = new RegExp('(' + k + ')')
    if (regx.test(fmt)) {
      const t = regx.exec(fmt)[1]
      fmt = fmt.replace(
        t,
        t.length === 1 ? o[k] : ('00' + o[k]).substring(('' + o[k]).length)
      )
    }
  }
  return fmt
}

export function formatTimeToStr(times, pattern) {
  let d = new Date(times).Format('yyyy-MM-dd hh:mm:ss')
  if (pattern) {
    d = new Date(times).Format(pattern)
  }
  return d.toLocaleString()
}
