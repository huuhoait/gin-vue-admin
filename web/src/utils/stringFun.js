/* eslint-disable */
export const toUpperCase = (str) => {
  if (str[0]) {
    return str.replace(str[0], str[0].toUpperCase())
  } else {
    return ''
  }
}

export const toLowerCase = (str) => {
  if (str[0]) {
    return str.replace(str[0], str[0].toLowerCase())
  } else {
    return ''
  }
}

// Convert camelCase to snake_case
export const toSQLLine = (str) => {
  if (str === 'ID') return 'ID'
  return str.replace(/([A-Z])/g, '_$1').toLowerCase()
}

// Convert snake_case to camelCase
export const toHump = (name) => {
  return name.replace(/\_(\w)/g, function (all, letter) {
    return letter.toUpperCase()
  })
}
