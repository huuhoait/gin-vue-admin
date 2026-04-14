import { useParamsStore } from '@/pinia/modules/params'
/*
 * Helper to load a param value.
 * Usage: await getParams('key')
 */
export const getParams = async (key) => {
  const paramsStore = useParamsStore()
  await paramsStore.getParams(key)
  return paramsStore.paramsMap[key]
}
