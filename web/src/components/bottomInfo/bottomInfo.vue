<!--
This file is protected by copyright. Do not modify without authorization.
If you have not obtained a license, contact the author to purchase one.
Unauthorized commercial use may lead to legal claims. Please respect intellectual property.
-->
<template>
  <div
    class="flex flex-col md:flex-row gap-2 items-center text-sm text-slate-700 dark:text-slate-500 justify-center py-2"
  >
    <div class="text-center">
      <span class="mr-1">Powered by</span>
      <span>
        <a
          class="font-bold text-active"
          href="https://github.com/huuhoaitvn/gin-vue-admin"
          >IIT</a
        >
      </span>
    </div>
    <slot />
    <div class="text-center">
      <span class="mr-1">Copyright</span>
      <span>
        <a
          class="font-bold text-active"
          href="https://github.com/huuhoaitvn"
          >IIT team</a
        >
      </span>
    </div>
  </div>
</template>

<script setup>
  defineOptions({
    name: 'BottomInfo'
  })

  console.log(
    `%c powered by %c huuhoaitvne %c`,
    'background:#0081ff; padding: 1px; border-radius: 3px 0 0 3px; color: #fff',
    'background:#354855; padding: 1px 5px; border-radius: 0 3px 3px 0; color: #fff; font-weight: bold;',
    'background:transparent'
  )
</script>
