<template>
  <div class="grid grid-cols-2 gap-2">
    <a
      v-for="item in wikis"
      :key="item.url"
      :href="item.url"
      class="text-sm text-black/70 dark:text-white/70 no-underline hover:text-[var(--el-color-primary)] dark:hover:text-white"
      target="_blank"
    >
      {{ item.title }}
    </a>
  </div>
</template>

<script setup>
  const wikis = [
    {
      title: 'Vue3',
      url: 'https://v3.cn.vuejs.org/guide/introduction.html'
    },
    {
      title: 'Gin docs',
      url: 'https://gin-gonic.com/'
    },
    {
      title: 'GVA docs',
      url: 'https://www.gin-vue-admin.com/'
    },
    {
      title: 'Plugin market',
      url: 'https://plugin.gin-vue-admin.com/'
    },
    {
      title: 'GitHub repo',
      url: 'https://github.com/huuhoaitvn/gin-vue-admin'
    }
  ]
</script>

<style scoped lang="scss"></style>
