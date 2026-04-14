<template>
  <vue-office-pdf
    :src="pdf"
    @rendered="renderedHandler"
    @error="errorHandler"
  />
</template>
<script>
  export default {
    name: 'Pdf'
  }
</script>
<script setup>
  import { ref, watch } from 'vue'

  // Import VueOfficePdf component
  import VueOfficePdf from '@vue-office/pdf'
  // Import styles
  import '@vue-office/docx/lib/index.css'
  console.log('pdf===>')
  const props = defineProps({
    modelValue: {
      type: String,
      default: () => ''
    }
  })
  const pdf = ref(null)
  watch(
    () => props.modelValue,
    (val) => (pdf.value = val),
    { immediate: true }
  )
  const renderedHandler = () => {
    console.log('pdf loaded')
  }
  const errorHandler = () => {
    console.log('pdf error')
  }
</script>
