<template>
  <div 
    class="fixed inset-0 bg-black/40 dark:bg-black/60 flex items-center justify-center z-[999]"
    @click.self="closeModal"
  >
    <div class="bg-white dark:bg-gray-800 rounded-xl shadow-dialog dark:shadow-lg w-full max-w-md mx-4 transform transition-all duration-300 ease-in-out border border-transparent dark:border-gray-700">
      <!-- Header -->
      <div class="p-5 border-b border-gray-100 dark:border-gray-700 flex justify-between items-center">
        <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-100">{{ displayData.title }}</h3>
        <div class="text-gray-400 dark:text-gray-300 hover:text-gray-600 dark:hover:text-gray-200 transition-colors cursor-pointer" @click="closeModal">
          <close class="h-6 w-6" />
        </div>
      </div>
      
      <!-- Content -->
      <div class="p-6 pt-0">
        <!-- Error type -->
        <div class="mb-4">
          <div class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase mb-2">Error type</div>
          <div class="flex items-center gap-2">
            <lock v-if="displayData.icon === 'lock'" :class="['w-5 h-5', displayData.color]" />
            <warn v-if="displayData.icon === 'warn'" :class="['w-5 h-5', displayData.color]" />
            <server v-if="displayData.icon === 'server'" :class="['w-5 h-5', displayData.color]" />
            <span class="font-medium text-gray-800 dark:text-gray-100">{{ displayData.type }}</span>
          </div>
        </div>
        
        <!-- Error details -->
        <div class="mb-6">
          <div class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase mb-2">Details</div>
          <div class="bg-gray-100 dark:bg-gray-900/40 rounded-lg p-3 text-sm text-gray-700 dark:text-gray-200 leading-relaxed">
            {{ displayData.message }}
          </div>
        </div>
        
        <!-- Tips -->
        <div v-if="displayData.tips">
          <div class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase mb-2">Tip</div>
          <div class="flex items-center gap-2">
            <idea class="text-blue-500 dark:text-blue-400 w-5 h-5" />
            <p class="text-sm text-gray-600 dark:text-gray-300">{{ displayData.tips }}</p>
          </div>
        </div>
      </div>
      
      <!-- Footer -->
      <div class="py-2 px-4 border-t border-gray-100 dark:border-gray-700 flex justify-end">
        <div class="px-4 py-2 bg-blue-600 dark:bg-blue-500 text-white dark:text-gray-100 rounded-lg hover:bg-blue-700 dark:hover:bg-blue-600 transition-colors font-medium text-sm shadow-sm cursor-pointer" @click="handleConfirm">
          OK
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  errorData: {
    type: Object,
    required: true
  }
});

const emits = defineEmits(['close', 'confirm']);

const presetErrors = {
  500: {
    title: 'API error detected',
    type: 'Internal server error',
    icon: 'server',
    color: 'text-red-500 dark:text-red-400',
    tips: 'This is often caused by backend panics. Check backend logs first. If it blocks usage, you can force logout and clear cache.'
  },
  404: {
    title: 'Resource not found',
    type: 'Not Found',
    icon: 'warn',
    color: 'text-orange-500 dark:text-orange-400',
    tips: 'This is often caused by an endpoint not being registered (or server not restarted), or mismatched request path/method. For generated code, also check for stray spaces.'
  },
  401: {
    title: 'Authentication failed',
    type: 'Invalid token',
    icon: 'lock',
    color: 'text-purple-500 dark:text-purple-400',
    tips: 'Your session has expired or is invalid. Please sign in again.'
  },
  'network': {
    title: 'Network error',
    type: 'Network Error',
    icon: 'fa-wifi-slash',
    color: 'text-gray-500 dark:text-gray-400',
    tips: 'Unable to reach the server. Please check your network connection.'
  }
};

const displayData = computed(() => {
  const preset = presetErrors[props.errorData.code];
  if (preset) {
    return {
      ...preset,
      message: props.errorData.message || 'No additional details provided.'
    };
  }

  return {
    title: 'Unknown error',
    type: 'Request error detected',
    icon: 'fa-question-circle',
    color: 'text-gray-400 dark:text-gray-300',
    message: props.errorData.message || 'An unknown error occurred.',
    tips: 'Check the console for more details.'
  };
});

const closeModal = () => {
   emits('close')
};

const handleConfirm = () => {
  emits('confirm', props.errorData.code);
  closeModal();
};
</script>
