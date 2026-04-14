export const getCode = (templateID) => {
  return `<template>
  <!-- Export component -->
  <ExportExcel templateId="${templateID}" :condition="condition" :limit="limit" :offset="offset" :order="order" />

  <!-- Import component: handleSuccess is called after import -->
  <ImportExcel templateId="${templateID}" @on-success="handleSuccess" />

  <!-- Export template -->
  <ExportTemplate templateId="${templateID}" />
</template>

<script setup>
import { ref } from 'vue';
// Export component
import ExportExcel from '@/components/exportExcel/exportExcel.vue';
// Import component
import ImportExcel from '@/components/exportExcel/importExcel.vue';
// Export template component
import ExportTemplate from '@/components/exportExcel/exportTemplate.vue';

const condition = ref({}); // query conditions
const limit = ref(10); // max rows
const offset = ref(0); // offset
const order = ref('id desc'); // sort order

const handleSuccess = (res) => {
  console.log(res);
  // import success callback
};
</script>`
}
