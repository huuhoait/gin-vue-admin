import { findSysDictionary } from '@/api/sysDictionary'
import { getDictionaryTreeListByType } from '@/api/sysDictionaryDetail'

import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useDictionaryStore = defineStore('dictionary', () => {
  const dictionaryMap = ref({})

  const setDictionaryMap = (dictionaryRes) => {
    dictionaryMap.value = { ...dictionaryMap.value, ...dictionaryRes }
  }

  // Filter tree depth
  const filterTreeByDepth = (items, currentDepth, targetDepth) => {
    if (targetDepth === 0) {
      // depth=0 returns all data
      return items
    }

    if (currentDepth >= targetDepth) {
      // Reached target depth: remove children
      return items.map((item) => ({
        label: item.label,
        value: item.value,
        extend: item.extend
      }))
    }

    // Recurse through children
    return items.map((item) => ({
      label: item.label,
      value: item.value,
      extend: item.extend,
      children: item.children
        ? filterTreeByDepth(item.children, currentDepth + 1, targetDepth)
        : undefined
    }))
  }

  // Flatten tree to array (compatible with legacy flat format)
  const flattenTree = (items) => {
    const result = []

    const traverse = (nodes) => {
      nodes.forEach((item) => {
        result.push({
          label: item.label,
          value: item.value,
          extend: item.extend
        })

        if (item.children && item.children.length > 0) {
          traverse(item.children)
        }
      })
    }

    traverse(items)
    return result
  }

  // Normalize tree nodes to a standard shape
  const normalizeTreeData = (items) => {
    return items.map((item) => ({
      label: item.label,
      value: item.value,
      extend: item.extend,
      children:
        item.children && item.children.length > 0
          ? normalizeTreeData(item.children)
          : undefined
    }))
  }

  // Find node by value and return its children (depth-limited)
  const findNodeByValue = (
    items,
    targetValue,
    currentDepth = 1,
    maxDepth = 0
  ) => {
    for (const item of items) {
      // Found target node
      if (item.value === targetValue) {
        // maxDepth=0: return all children
        if (maxDepth === 0) {
          return item.children ? normalizeTreeData(item.children) : []
        }
        // Otherwise return children limited by depth
        if (item.children && item.children.length > 0) {
          return filterTreeByDepth(item.children, 1, maxDepth)
        }
        return []
      }

      // If current depth is less than max, continue searching in children
      if (
        item.children &&
        item.children.length > 0 &&
        (maxDepth === 0 || currentDepth < maxDepth)
      ) {
        const result = findNodeByValue(
          item.children,
          targetValue,
          currentDepth + 1,
          maxDepth
        )
        if (result !== null) {
          return result
        }
      }
    }
    return null
  }

  const getDictionary = async (type, depth = 0, value = null) => {
    // If value is provided, return children of the specified node
    if (value !== null) {
      // Build cache key including value and depth
      const cacheKey = `${type}_value_${value}_depth_${depth}`

      if (
        dictionaryMap.value[cacheKey] &&
        dictionaryMap.value[cacheKey].length
      ) {
        return dictionaryMap.value[cacheKey]
      }

      try {
        // Fetch full tree data
        const treeRes = await getDictionaryTreeListByType({ type })
        if (
          treeRes.code === 0 &&
          treeRes.data &&
          treeRes.data.list &&
          treeRes.data.list.length > 0
        ) {
          // Find node by value and return its children
          const targetNodeChildren = findNodeByValue(
            treeRes.data.list,
            value,
            1,
            depth
          )

          if (targetNodeChildren !== null) {
            let resultData
            if (depth === 0) {
              // depth=0: return full children tree
              resultData = targetNodeChildren
            } else {
              // Other depths: flatten children
              resultData = flattenTree(targetNodeChildren)
            }

            const dictionaryRes = {}
            dictionaryRes[cacheKey] = resultData
            setDictionaryMap(dictionaryRes)
            return dictionaryMap.value[cacheKey]
          } else {
            // If not found, return empty array
            return []
          }
        }
      } catch (error) {
        console.error('Failed to load dictionary data by value:', error)
        return []
      }
    }

    // Legacy behavior: when value is not provided
    // Build cache key including depth
    const cacheKey = depth === 0 ? `${type}_tree` : `${type}_depth_${depth}`

    if (dictionaryMap.value[cacheKey] && dictionaryMap.value[cacheKey].length) {
      return dictionaryMap.value[cacheKey]
    } else {
      try {
        // First try to get tree data
        const treeRes = await getDictionaryTreeListByType({ type })
        if (
          treeRes.code === 0 &&
          treeRes.data &&
          treeRes.data.list &&
          treeRes.data.list.length > 0
        ) {
          // Use tree data
          const treeData = treeRes.data.list

          let resultData
          if (depth === 0) {
            // depth=0: return full tree with normalized node shape
            resultData = normalizeTreeData(treeData)
          } else {
            // Other depths: filter by depth, then flatten
            const filteredData = filterTreeByDepth(treeData, 1, depth)
            resultData = flattenTree(filteredData)
          }

          const dictionaryRes = {}
          dictionaryRes[cacheKey] = resultData
          setDictionaryMap(dictionaryRes)
          return dictionaryMap.value[cacheKey]
        } else {
          // If no tree data, fall back to legacy flat list
          const res = await findSysDictionary({ type })
          if (res.code === 0) {
            const dictionaryRes = {}
            const dict = []
            res.data.resysDictionary.sysDictionaryDetails &&
              res.data.resysDictionary.sysDictionaryDetails.forEach((item) => {
                dict.push({
                  label: item.label,
                  value: item.value,
                  extend: item.extend
                })
              })
            dictionaryRes[cacheKey] = dict
            setDictionaryMap(dictionaryRes)
            return dictionaryMap.value[cacheKey]
          }
        }
      } catch (error) {
        console.error('Failed to load dictionary data:', error)
        // On error, fall back to legacy approach
        const res = await findSysDictionary({ type })
        if (res.code === 0) {
          const dictionaryRes = {}
          const dict = []
          res.data.resysDictionary.sysDictionaryDetails &&
            res.data.resysDictionary.sysDictionaryDetails.forEach((item) => {
              dict.push({
                label: item.label,
                value: item.value,
                extend: item.extend
              })
            })
          dictionaryRes[cacheKey] = dict
          setDictionaryMap(dictionaryRes)
          return dictionaryMap.value[cacheKey]
        }
      }
    }
  }

  return {
    dictionaryMap,
    setDictionaryMap,
    getDictionary
  }
})
