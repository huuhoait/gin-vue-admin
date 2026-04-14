<template>
  <div>
    <warning-bar
        href="https://plugin.gin-vue-admin.com/license"
        title="This feature is available for licensed users only. Click to purchase a license."
    />
    <div class="gva-search-box">
      <div class="text-xl mb-2 text-gray-600">
        AI Frontend Engineer <a
          class="text-blue-600 text-sm ml-4"
          href="https://plugin.gin-vue-admin.com/#/layout/userInfo/center"
          target="_blank"
      >Get AI-Path</a
      >
      </div>
      
      <!-- Options mode -->
      <div class="mb-4">
        <div class="mb-3">
          <div class="text-base font-medium mb-2">Page purpose</div>
          <el-radio-group v-model="pageType" class="mb-2" @change="handlePageTypeChange">
            <el-radio label="Corporate website">Corporate website</el-radio>
            <el-radio label="E-commerce page">E-commerce page</el-radio>
            <el-radio label="Personal blog">Personal blog</el-radio>
            <el-radio label="Product introduction">Product introduction</el-radio>
            <el-radio label="Campaign landing page">Campaign landing page</el-radio>
            <el-radio label="Other">Other</el-radio>
          </el-radio-group>
          <el-input v-if="pageType === 'Other'" v-model="pageTypeCustom" placeholder="Enter page purpose" class="w-full" />
        </div>
        
        <div class="mb-3">
          <div class="text-base font-medium mb-2">Main content sections</div>
          <el-checkbox-group v-model="contentBlocks" class="flex flex-wrap gap-2 mb-2">
            <el-checkbox label="Banner carousel">Banner carousel</el-checkbox>
            <el-checkbox label="Product/Service overview">Product/Service overview</el-checkbox>
            <el-checkbox label="Features">Features</el-checkbox>
            <el-checkbox label="Customer cases">Customer cases</el-checkbox>
            <el-checkbox label="Team">Team</el-checkbox>
            <el-checkbox label="Contact form">Contact form</el-checkbox>
            <el-checkbox label="News/Blog list">News/Blog list</el-checkbox>
            <el-checkbox label="Pricing table">Pricing table</el-checkbox>
            <el-checkbox label="FAQ">FAQ</el-checkbox>
            <el-checkbox label="Testimonials">Testimonials</el-checkbox>
            <el-checkbox label="Statistics">Statistics</el-checkbox>
            <el-checkbox label="Product list">Product list</el-checkbox>
            <el-checkbox label="Product cards">Product cards</el-checkbox>
            <el-checkbox label="Cart">Cart</el-checkbox>
            <el-checkbox label="Checkout">Checkout</el-checkbox>
            <el-checkbox label="Order tracking">Order tracking</el-checkbox>
            <el-checkbox label="Categories">Categories</el-checkbox>
            <el-checkbox label="Hot picks">Hot picks</el-checkbox>
            <el-checkbox label="Limited-time deals">Limited-time deals</el-checkbox>
            <el-checkbox label="Other">Other</el-checkbox>
          </el-checkbox-group>
          <el-input v-if="contentBlocks.includes('Other')" v-model="contentBlocksCustom" placeholder="Enter other sections" class="w-full" />
        </div>
        
        <div class="mb-3">
          <div class="text-base font-medium mb-2">Style preference</div>
          <el-radio-group v-model="stylePreference" class="mb-2">
            <el-radio label="Minimal">Minimal</el-radio>
            <el-radio label="Tech">Tech</el-radio>
            <el-radio label="Warm">Warm</el-radio>
            <el-radio label="Professional">Professional</el-radio>
            <el-radio label="Creative">Creative</el-radio>
            <el-radio label="Retro">Retro</el-radio>
            <el-radio label="Luxury">Luxury</el-radio>
            <el-radio label="Other">Other</el-radio>
          </el-radio-group>
          <el-input v-if="stylePreference === 'Other'" v-model="stylePreferenceCustom" placeholder="Enter style preference" class="w-full" />
        </div>
        
        <div class="mb-3">
          <div class="text-base font-medium mb-2">Layout</div>
          <el-radio-group v-model="layoutDesign" class="mb-2">
            <el-radio label="Single column">Single column</el-radio>
            <el-radio label="Two column">Two column</el-radio>
            <el-radio label="Three column">Three column</el-radio>
            <el-radio label="Grid">Grid</el-radio>
            <el-radio label="Gallery">Gallery</el-radio>
            <el-radio label="Masonry">Masonry</el-radio>
            <el-radio label="Card-based">Card-based</el-radio>
            <el-radio label="Sidebar + content">Sidebar + content</el-radio>
            <el-radio label="Split screen">Split screen</el-radio>
            <el-radio label="Full-page scrolling">Full-page scrolling</el-radio>
            <el-radio label="Mixed">Mixed</el-radio>
            <el-radio label="Responsive">Responsive</el-radio>
            <el-radio label="Other">Other</el-radio>
          </el-radio-group>
          <el-input v-if="layoutDesign === 'Other'" v-model="layoutDesignCustom" placeholder="Enter layout preference" class="w-full" />
        </div>
        
        <div class="mb-3">
          <div class="text-base font-medium mb-2">Color scheme</div>
          <el-radio-group v-model="colorScheme" class="mb-2">
            <el-radio label="Blues">Blues</el-radio>
            <el-radio label="Greens">Greens</el-radio>
            <el-radio label="Reds">Reds</el-radio>
            <el-radio label="Grayscale">Grayscale</el-radio>
            <el-radio label="Pure black/white">Pure black/white</el-radio>
            <el-radio label="Warm tones">Warm tones</el-radio>
            <el-radio label="Cool tones">Cool tones</el-radio>
            <el-radio label="Other">Other</el-radio>
          </el-radio-group>
          <el-input v-if="colorScheme === 'Other'" v-model="colorSchemeCustom" placeholder="Enter color scheme" class="w-full" />
        </div>
      </div>
      
      <!-- Detailed description input -->
      <div class="relative">
        <div class="text-base font-medium mb-2">Detailed description (optional)</div>
        <el-input
            v-model="prompt"
            :maxlength="2000"
            :placeholder="placeholder"
            :rows="5"
            resize="none"
            type="textarea"
            @blur="handleBlur"
            @focus="handleFocus"
        />
        <div class="flex absolute right-2 bottom-2">
          <el-tooltip effect="light">
            <template #content>
              <div>
                This feature is available for licensed users only. Go to <a
                  class="text-blue-600"
                  href="https://plugin.gin-vue-admin.com/license"
                  target="_blank"
              >Purchase a license</a
              >
              </div>
            </template>
            <el-button
                type="primary"
                @click="llmAutoFunc()"
            >
              <el-icon size="18">
                <ai-gva/>
              </el-icon>
              Generate
            </el-button>
          </el-tooltip>
        </div>
      </div>
    </div>
    <div>
      <div v-if="!outPut">
        <el-empty :image-size="200"/>
      </div>
      <div v-if="outPut && htmlFromLLM">
        <el-tabs type="border-card">
          <el-tab-pane label="Preview">
            <div class="h-[500px] overflow-auto bg-gray-50 p-4 rounded">
              <div v-if="!loadedComponents" class="text-gray-500 text-center py-4">
                Loading component...
              </div>
              <component
                v-else
                :is="loadedComponents" 
                class="vue-component-container w-full"
              />
            </div>
          </el-tab-pane>
          <el-tab-pane label="Source code">
            <div class="relative h-[500px] overflow-auto bg-gray-50 p-4 rounded">
              <el-button 
                type="primary" 
                :icon="DocumentCopy" 
                class="absolute top-2 right-2 px-2 py-1" 
                @click="copySnippet(htmlFromLLM)" 
                plain
              >
                Copy
              </el-button>
              <pre class="mt-10 whitespace-pre-wrap">{{ htmlFromLLM }}</pre>
            </div>
          </el-tab-pane>
        </el-tabs>
      </div>
    </div>
  </div>
</template>

<script setup>
import { llmAuto } from '@/api/autoCode'
import { ref, reactive, markRaw } from 'vue'
import * as Vue from "vue";
import WarningBar from '@/components/warningBar/warningBar.vue'
import { ElMessage } from 'element-plus'
import { defineAsyncComponent } from 'vue'
import { DocumentCopy } from '@element-plus/icons-vue'
import { loadModule } from "vue3-sfc-loader";

defineOptions({
  name: 'Picture'
})

const handleFocus = () => {
  document.addEventListener('keydown', handleKeydown);
}

const handleBlur = () => {
  document.removeEventListener('keydown', handleKeydown);
}

const handleKeydown = (event) => {
  if ((event.ctrlKey || event.metaKey) && event.key === 'Enter') {
    llmAutoFunc()
  }
}

// Copy helper: write string to clipboard
const copySnippet = (vueString) => {
  navigator.clipboard.writeText(vueString)
      .then(() => {
        ElMessage({
          message: 'Copied',
          type: 'success',
        })
      })
      .catch(err => {
        ElMessage({
          message: 'Copy failed',
          type: 'warning',
        })
      })
}

// Options mode state
const pageType = ref('Corporate website')
const pageTypeCustom = ref('')
const contentBlocks = ref(['Banner carousel', 'Product/Service overview'])
const contentBlocksCustom = ref('')
const stylePreference = ref('Minimal')
const stylePreferenceCustom = ref('')
const layoutDesign = ref('Responsive')
const layoutDesignCustom = ref('')
const colorScheme = ref('Blues')
const colorSchemeCustom = ref('')

// Recommended mapping from page purpose to content sections
const pageTypeContentMap = {
  'Corporate website': ['Banner carousel', 'Product/Service overview', 'Features', 'Customer cases', 'Contact form'],
  'E-commerce page': ['Banner carousel', 'Product list', 'Product cards', 'Cart', 'Categories', 'Hot picks', 'Limited-time deals', 'Checkout', 'Testimonials'],
  'Personal blog': ['Banner carousel', 'News/Blog list', 'Testimonials', 'Contact form'],
  'Product introduction': ['Banner carousel', 'Product/Service overview', 'Features', 'Pricing table', 'FAQ'],
  'Campaign landing page': ['Banner carousel', 'Features', 'Contact form', 'Statistics']
}

const prompt = ref('')

// Whether output exists
const outPut = ref(false)
// Holds generated Vue component code
const htmlFromLLM = ref("")

// Loaded component
const loadedComponents = ref(null)

const loadVueComponent = async (vueCode) => {
  try {
    // Use an in-memory virtual path
    const fakePath = `virtual:component-0.vue`
    
    const component = defineAsyncComponent({
      loader: async () => {
        try {
          const options = {
            moduleCache: {
              vue: Vue,
            },
            getFile(url) {
              // Handle all possible URL formats (relative, absolute, etc)
              // Extract file name and ignore query params
              const fileName = url.split('/').pop().split('?')[0]
              const componentFileName = fakePath.split('/').pop()
              
              // If the name matches our component, or url fully matches fakePath
              if (fileName === componentFileName || url === fakePath || 
                  url === `./component/0.vue`) {
                return Promise.resolve({
                  type: '.vue',
                  getContentData: () => vueCode
                })
              }
              
              console.warn('Unknown file requested:', url)
              return Promise.reject(new Error(`File not found: ${url}`))
            },
            addStyle(textContent) {
              // Don't append styles to document.head; return style content.
              // We'll attach it to Shadow DOM later.
              return textContent
            },
            handleModule(type, source, path, options) {
              // Default handler
              return undefined
            },
            log(type, ...args) {
              console.log(`[vue3-sfc-loader] [${type}]`, ...args)
            }
          }
          
          // Load the component
          const comp = await loadModule(fakePath, options)
          return comp.default || comp
        } catch (error) {
          console.error('Component load details:', error)
          throw error
        }
      },
      loadingComponent: {
        template: '<div>Loading...</div>'
      },
      errorComponent: {
        props: ['error'],
        template: '<div>Component failed to load: {{ error && error.message }}</div>',
        setup(props) {
          console.error('Error component received:', props.error)
          return {}
        }
      },
      // Timeout and retry options
      timeout: 30000,
      delay: 200,
      suspensible: false,
      onError(error, retry, fail) {
        console.error('Load error details:', error)
        fail()
      }
    })

    // Wrapper component to isolate styles via Shadow DOM
    const ShadowWrapper = {
      name: 'ShadowWrapper',
      setup() {
        return {}
      },
      render() {
        return Vue.h('div', { class: 'shadow-wrapper' })
      },
      mounted() {
        // Create Shadow DOM
        const shadowRoot = this.$el.attachShadow({ mode: 'open' })
        
        // Container element
        const container = document.createElement('div')
        container.className = 'shadow-container'
        shadowRoot.appendChild(container)
        
        // Extract styles from SFC
        const styleContent = vueCode.match(/<style[^>]*>([\s\S]*?)<\/style>/i)?.[1] || ''
        
        // Add styles to Shadow DOM
        if (styleContent) {
          const style = document.createElement('style')
          style.textContent = styleContent
          shadowRoot.appendChild(style)
        }
        
        // Mount Vue app into Shadow DOM container
        const app = Vue.createApp({
          render: () => Vue.h(component)
        })
        app.mount(container)
      }
    }

    loadedComponents.value = markRaw(ShadowWrapper)
    return ShadowWrapper
  } catch (error) {
    console.error('Component create error:', error)
    return null
  }
}

// When page purpose changes, update section selection
const handlePageTypeChange = (value) => {
  if (value !== 'Other' && pageTypeContentMap[value]) {
    contentBlocks.value = [...pageTypeContentMap[value]]
  }
}

const llmAutoFunc = async () => {
  // Build full prompt including options mode selection
  let fullPrompt = ''
  
  // Page purpose
  fullPrompt += `Page purpose: ${pageType.value === 'Other' ? pageTypeCustom.value : pageType.value}\n`
  
  // Content sections
  fullPrompt += 'Main content sections: '
  const blocks = contentBlocks.value.filter(block => block !== 'Other')
  if (contentBlocksCustom.value) {
    blocks.push(contentBlocksCustom.value)
  }
  fullPrompt += blocks.join(', ') + '\n'
  
  // Style preference
  fullPrompt += `Style preference: ${stylePreference.value === 'Other' ? stylePreferenceCustom.value : stylePreference.value}\n`
  
  // Layout
  fullPrompt += `Layout: ${layoutDesign.value === 'Other' ? layoutDesignCustom.value : layoutDesign.value}\n`
  
  // Color scheme
  fullPrompt += `Color scheme: ${colorScheme.value === 'Other' ? colorSchemeCustom.value : colorScheme.value}\n`
  
  // User's detailed description
  if (prompt.value) {
    fullPrompt += `\nDetailed description: ${prompt.value}`
  }
  
  const res = await llmAuto({web: fullPrompt, mode: 'createWeb'})
  if (res.code === 0) {
    outPut.value = true
    // Store generated component code
    htmlFromLLM.value = res.data.text
    // Load generated component
    await loadVueComponent(res.data.text)
  }
}

const placeholder = ref(`Add any extra requirements, e.g. emphasized elements, reference websites, or interaction effects.`)
</script>
