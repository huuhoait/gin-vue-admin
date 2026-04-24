<template>
  <div>
    <div class="w-full h-screen bg-gray-50 flex items-center justify-center">
      <div class="flex flex-col items-center text-2xl gap-4">
        <img class="w-1/3" src="../../assets/404.png" />
        <p class="text-lg">This page was whisked away by mysterious forces</p>
        <p class="text-lg">
          Common cause: your role does not have access to this route. If you need it, ask an admin to assign it under Role Management.
        </p>
        <p>
          Repo: <a
            href="https://github.com/huuhoaitvn/gin-vue-admin"
            target="_blank"
            class="text-blue-600"
            >https://github.com/huuhoaitvn/gin-vue-admin</a
          >
        </p>
        <el-button @click="toDashboard">Back to home</el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
  import { useUserStore } from '@/pinia/modules/user'
  import { useRouter } from 'vue-router'
  import { emitter } from '@/utils/bus'

  defineOptions({
    name: 'Error'
  })

  const userStore = useUserStore()
  const router = useRouter()
  const toDashboard = () => {
    try {
      router.push({ name: userStore.userInfo.authority.defaultRouter })
    } catch (error) {
        emitter.emit('show-error', {
        code: '401',
        message: "Route permissions changed. Please sign in again.",
        fn: () => {
          userStore.ClearStorage()
          router.push({ name: 'Login', replace: true })
        }
      })
    }
  }
</script>
