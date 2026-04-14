import { useDictionaryStore } from '@/pinia/modules/dictionary'

/**
 * Generate dictionary cache key
 * @param {string} type - dictionary type
 * @param {number} depth - depth
 * @param {string|number|null} value - target node value
 * @returns {string} cache key
 */
const generateCacheKey = (type, depth, value) => {
  if (value !== null && value !== undefined) {
    return `${type}_value_${value}_depth_${depth}`
  }
  return depth === 0 ? `${type}_tree` : `${type}_depth_${depth}`
}

/**
 * Get dictionary data
 * @param {string} type - dictionary type (required)
 * @param {Object} options - optional params
 * @param {number} options.depth - depth, defaults to 0 (full tree)
 * @param {string|number|null} options.value - node value, returns its children, defaults to null
 * @returns {Promise<Array>} dictionary items
 * @example
 * // Get full dictionary tree
 * const dictTree = await getDict('user_status')
 *
 * // Get flattened list with given depth
 * const dictFlat = await getDict('user_status', {
 *  depth: 2
 * })
 *
 * // Get children for a given node
 * const children = await getDict('user_status', {
 *  value: 'active'
 * })
 */
export const getDict = async (
  type,
  options = {
    depth: 0,
    value: null
  }
) => {
  // Validate params
  if (!type || typeof type !== 'string') {
    console.warn('getDict: type must be a non-empty string')
    return []
  }

  if (typeof options.depth !== 'number' || options.depth < 0) {
    console.warn('getDict: depth must be a non-negative number')
    options.depth = 0
  }

  try {
    const dictionaryStore = useDictionaryStore()

    // Load data through the store
    await dictionaryStore.getDictionary(type, options.depth, options.value)

    // Build cache key
    const cacheKey = generateCacheKey(type, options.depth, options.value)

    // Read from cache
    const result = dictionaryStore.dictionaryMap[cacheKey]

    // Return array
    return Array.isArray(result) ? result : []
  } catch (error) {
    console.error('getDict: failed to load dictionary data', { type, options, error })
    return []
  }
}

// Dictionary label helper
export const showDictLabel = (
  dict,
  code,
  keyCode = 'value',
  valueCode = 'label'
) => {
  if (!dict) {
    return ''
  }
  const dictMap = {}
  dict.forEach((item) => {
    if (Reflect.has(item, keyCode) && Reflect.has(item, valueCode)) {
      dictMap[item[keyCode]] = item[valueCode]
    }
  })
  return Reflect.has(dictMap, code) ? dictMap[code] : ''
}
