<script setup>
import { ref, watchEffect } from 'vue';
import { useAppStore } from '@/pinia/modules/app.js';
import { storeToRefs } from 'pinia';

const props = defineProps({
  // Logo size (rem)
  size: {
    type: Number,
    default: 2
  }
})

const darkLogoPath = "/logo.png";  // No dark-mode logo by default; adjust path if needed.
const lightLogoPath = "/logo.png";

const appStore = useAppStore();
const { isDark } = storeToRefs(appStore);

const logoSrc = ref('');
const showTextPlaceholder = ref(false);

// Check whether image exists
function checkImageExists(url) {
  return new Promise((resolve) => {
    const tryToLoad = new Image();
    tryToLoad.onload = () => resolve(true);
    tryToLoad.onerror = () => resolve(false);
    tryToLoad.src = url;
  });
}

watchEffect(async () => {
  showTextPlaceholder.value = false; // reset placeholder state

  // In dark mode, try dark logo first
  if (isDark.value && await checkImageExists(darkLogoPath)) {
    logoSrc.value = darkLogoPath;
    return;
  }

  if (await checkImageExists(lightLogoPath)) {
    logoSrc.value = lightLogoPath;
    return
  }

  // No usable logo found
  showTextPlaceholder.value = true;
  console.error(
    'Error: logo.png (or logo-dark.png) was not found in public directory.'
  );
  console.warn(
    'Fix: put logo.png and/or logo-dark.png under /public, or ensure the path is correct.'
  );
});

// Use 16px as baseline size
const SPACING = 16
function getSize() {
  return {
    width: `${props.size * SPACING}px`,
    height: `${props.size * SPACING}px`,
  }
}
</script>

<template>
  <img v-if="!showTextPlaceholder && logoSrc" :src="logoSrc" :alt="$GIN_VUE_ADMIN.appName" class="object-contain"
    :style="{
      ...getSize()
    }" :class="{
      'filter invert-[90%] hue-rotate-180 brightness-110':
        isDark && logoSrc === '/logo.png',
    }" />
  <div v-else-if="showTextPlaceholder"
    class="rounded-full bg-gray-300 dark:bg-gray-600 flex items-center justify-center text-gray-700 dark:text-gray-200 font-bold text-xs"
    :style="{
      ...getSize()
    }">
    GVA
  </div>
</template>
