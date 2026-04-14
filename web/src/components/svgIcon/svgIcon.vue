<template>
  <template v-if="localIcon">
    <svg aria-hidden="true" width="1em" height="1em" v-bind="bindAttrs">
      <use :xlink:href="'#' + localIcon" rel="external nofollow" />
    </svg>
  </template>
  <template v-else-if="icon">
    <Icon :icon="icon" v-bind="bindAttrs" />
  </template>
</template>

<script setup>
import { computed, useAttrs } from 'vue';
import { Icon } from '@iconify/vue'

/**
 * Usage:
 * Local icons (all available local icons are printed in console):
 * <SvgIcon localIcon="lock" class="text-red-500 text-3xl" />
 * 
 * Online icons (see: https://icones.js.org/ or https://icon-sets.iconify.design/):
 * <SvgIcon icon="mingcute:love-fill" class="text-red-500 text-3xl" />
 */
defineProps({
  // Use local registered svg icon by symbol id
  localIcon: {
    type: String,
    required: false,
    default: ''
  },
  // Iconify icon name, e.g. 'mdi:home'
  icon: {
    type: String,
    required: false,
    default: ''
  }
})
const attrs = useAttrs();

const bindAttrs = computed(() => ({
  class: (attrs.class) || '',
  style: (attrs.style) || ''
}))
</script>
