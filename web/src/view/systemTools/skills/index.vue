<template>
  <div class="h-full">
    <warning-bar
        href="https://plugin.gin-vue-admin.com/license"
        title="Development-only feature for building the in-project skills library."
    />
    <el-row :gutter="12" class="h-full">
      <el-col :xs="24" :sm="8" :md="6" :lg="5" class="flex flex-col gap-4 h-full">
        <el-card shadow="never" class="!border-none shrink-0">
          <div class="font-bold mb-2">AI Tools</div>
          <div class="flex flex-wrap gap-2">
            <div
              v-for="tool in tools"
              :key="tool.key"
              class="px-3 py-1.5 rounded-md text-sm cursor-pointer transition-all border select-none"
              :class="activeTool === tool.key
                ? 'bg-[var(--el-color-primary)] text-white border-[var(--el-color-primary)] shadow-sm'
                : 'bg-white hover:bg-gray-50 text-gray-700 border-gray-200 dark:bg-gray-800 dark:text-gray-300 dark:border-gray-700 dark:hover:bg-gray-700'"
              @click="handleToolSelect(tool.key)"
            >
              {{ tool.label }}
            </div>
          </div>
        </el-card>

        <el-card shadow="never" class="!border-none shrink-0">
          <div class="flex justify-between items-center mb-2">
            <span class="font-bold">Global constraint</span>
            <el-button type="primary" link icon="Edit" @click="openGlobalConstraint">Edit</el-button>
          </div>
          <div class="text-xs text-gray-500">Path: {{ globalConstraintPath }}</div>
        </el-card>

        <el-card shadow="never" class="!border-none flex-1 mt-2 flex flex-col min-h-0">
          <div class="flex justify-between items-center mb-2">
            <span class="font-bold">Skills</span>
            <div class="flex gap-1">
              <el-button type="primary" link icon="Download" @click="openOnlineDrawer">Online</el-button>
              <el-button type="primary" link icon="Plus" @click="openCreateDialog">New</el-button>
            </div>
          </div>
          <el-input
            v-model="skillFilter"
            size="small"
            clearable
            placeholder="Search skills"
            class="mb-2"
            prefix-icon="Search"
          />
          <el-scrollbar class="h-[calc(100vh-380px)]">
            <el-menu :default-active="activeSkill" class="!border-none" @select="handleSkillSelect">
              <el-menu-item
                v-for="skill in filteredSkills"
                :key="skill"
                :index="skill"
                class="!h-10 !leading-10 !my-1 !mx-1 !rounded-[4px]"
              >
                <div class="w-full flex items-center justify-between min-w-0">
                  <div class="flex items-center min-w-0 gap-1">
                    <el-icon><Document /></el-icon>
                    <span class="truncate" :title="skill">{{ skill }}</span>
                  </div>
                  <el-button
                    type="danger"
                    link
                    icon="Delete"
                    @click.stop="handleDeleteSkill(skill)"
                  >
                    Delete
                  </el-button>
                </div>
              </el-menu-item>
            </el-menu>
          </el-scrollbar>
        </el-card>
      </el-col>

      <el-col :xs="24" :sm="16" :md="18" :lg="19" class="h-full">
        <el-card shadow="never" class="!border-none h-full flex flex-col">
          <template v-if="!activeSkill">
            <div class="h-full flex items-center justify-center">
              <el-empty description="Select or create a skill" />
            </div>
          </template>
          <template v-else>
            <div class="flex justify-between items-center mb-4 pb-4 border-b border-gray-100 dark:border-gray-800">
              <div class="text-lg font-bold flex items-center gap-2">
                <span>{{ activeSkill }}</span>
                <el-tag size="small" type="info">Skill</el-tag>
              </div>
              <div class="flex items-center gap-2">
                <el-button icon="Download" @click="packageCurrentSkill">Package</el-button>
                <el-button type="primary" icon="Check" @click="saveCurrentSkill">Save</el-button>
              </div>
            </div>

            <el-tabs v-model="activeTab" class="h-full">
              <el-tab-pane label="Skill config" name="config">
                <div
                  class="mt-4 mb-4 rounded-md border border-gray-100 dark:border-gray-800 bg-gray-50 dark:bg-gray-900/30 p-3 text-xs text-gray-600 dark:text-gray-300"
                >
                  <div class="font-medium text-gray-700 dark:text-gray-200 mb-2">Guidelines</div>
                  <ul class="list-disc pl-4 space-y-1">
                    <li>Name: use lowercase + kebab-case as the skill directory and /slash command.</li>
                    <li>Description: include triggers and keywords (used for auto-matching).</li>
                    <li>Content: organize by Instructions / Examples / Guidelines, and specify inputs/outputs/constraints.</li>
                    <li>For templates/specs/scripts, reference templates/... , references/... , scripts/... in the markdown.</li>
                  </ul>
                </div>
                <el-form :model="form" label-width="160px">
                  <el-form-item>
                    <template #label>
                      <div class="flex items-center">
                        Name
                        <el-tooltip content="Unique identifier. Recommended: lowercase kebab-case." placement="top">
                          <el-icon class="ml-1 cursor-pointer"><QuestionFilled /></el-icon>
                        </el-tooltip>
                      </div>
                    </template>
                    <el-input v-model="form.name" placeholder="e.g. code-comment-expert" />
                    <div class="text-xs text-gray-400 mt-1">Recommended: 2–4 words, no spaces, ASCII only.</div>
                  </el-form-item>
                  <el-form-item>
                    <template #label>
                      <div class="flex items-center">
                        Description
                        <el-tooltip content="Describe when to use it + triggers + keywords. This is the most important line." placement="top">
                          <el-icon class="ml-1 cursor-pointer"><QuestionFilled /></el-icon>
                        </el-tooltip>
                      </div>
                    </template>
                    <el-input
                      v-model="form.description"
                      placeholder="e.g. Add bilingual comments for code review/refactor/readability"
                    />
                    <div class="text-xs text-gray-400 mt-1">Include: task type, scenario, keywords.</div>
                  </el-form-item>
                  <el-form-item>
                    <template #label>
                      <div class="flex items-center">
                        Allowed Tools
                        <el-tooltip content="Limit tool scope, e.g. Bash(gh *), Read, Write" placement="top">
                          <el-icon class="ml-1 cursor-pointer"><QuestionFilled /></el-icon>
                        </el-tooltip>
                      </div>
                    </template>
                    <el-input v-model="form.allowedTools" placeholder="Optional, e.g. Bash(gh *), Read, Write" />
                    <div class="text-xs text-gray-400 mt-1">Optional. If empty, it will be removed on save.</div>
                  </el-form-item>
                  <el-form-item>
                    <template #label>
                      <div class="flex items-center">
                        Context
                        <el-tooltip content="fork = isolated context, suitable for complex tasks" placement="top">
                          <el-icon class="ml-1 cursor-pointer"><QuestionFilled /></el-icon>
                        </el-tooltip>
                      </div>
                    </template>
                    <el-input v-model="form.context" placeholder="Optional, e.g. fork" />
                    <div class="text-xs text-gray-400 mt-1">Optional. If empty, it will be removed on save.</div>
                  </el-form-item>
                  <el-form-item>
                    <template #label>
                      <div class="flex items-center">
                        Agent
                        <el-tooltip content="When context=fork, you can set a sub-agent, e.g. Explore" placement="top">
                          <el-icon class="ml-1 cursor-pointer"><QuestionFilled /></el-icon>
                        </el-tooltip>
                      </div>
                    </template>
                    <el-input v-model="form.agent" placeholder="Optional, e.g. Explore / Build" />
                    <div class="text-xs text-gray-400 mt-1">Optional. If empty, it will be removed on save.</div>
                  </el-form-item>
                  <el-form-item>
                    <template #label>
                      <div class="flex items-center">
                        Markdown content
                        <el-tooltip content="Keep it concise; put heavy details in templates/references/scripts." placement="top">
                          <el-icon class="ml-1 cursor-pointer"><QuestionFilled /></el-icon>
                        </el-tooltip>
                      </div>
                    </template>
                    <div class="mb-2 flex flex-wrap gap-2">
                      <el-button
                        v-for="block in quickBlocks"
                        :key="block.label"
                        size="small"
                        @click="appendMarkdown(block.content)"
                      >
                        {{ block.label }}
                      </el-button>
                      <el-button size="small" @click="insertFullTemplate">Insert full template</el-button>
                    </div>
                    <el-input
                      v-model="form.markdown"
                      type="textarea"
                      :rows="20"
                      :placeholder="markdownPlaceholder"
                    />
                    <div class="text-xs text-gray-400 mt-1">
                      Keep markdown short; put details in templates/references/scripts and reference them by relative paths.
                    </div>
                  </el-form-item>
                </el-form>
              </el-tab-pane>

              <el-tab-pane label="Scripts" name="scripts" class="mt-4">
                <div class="flex justify-between items-center mb-4">
                  <div class="text-sm text-gray-500 bg-gray-50 dark:bg-gray-800 px-3 py-1 rounded">Path: scripts/</div>
                  <el-button type="primary" icon="Plus" size="small" @click="openScriptDialog">Create script</el-button>
                </div>
                <div class="text-xs text-gray-500 mb-3">
                  Put executable logic / validation here. Reference <span class="font-mono">scripts/filename</span> in markdown (requires code execution enabled).
                </div>
                <el-table :data="scriptRows" style="width: 100%">
                  <el-table-column prop="name" label="File">
                    <template #default="scope">
                      <div class="flex items-center gap-2">
                        <el-icon><Document /></el-icon>
                        <span>{{ scope.row.name }}</span>
                      </div>
                    </template>
                  </el-table-column>
                  <el-table-column label="Actions" width="180">
                    <template #default="scope">
                      <el-button type="primary" link icon="Edit" @click="openScriptEditor(scope.row.name)">Edit</el-button>
                      <el-button type="primary" link @click="insertFileSnippet('script', scope.row.name)">Insert</el-button>
                    </template>
                  </el-table-column>
                </el-table>
                <el-empty v-if="scriptRows.length === 0" description="No scripts" />
              </el-tab-pane>

              <el-tab-pane label="Resources" name="resources">
                <div class="flex justify-between items-center mb-4 mt-4">
                  <div class="text-sm text-gray-500 bg-gray-50 dark:bg-gray-800 px-3 py-1 rounded">Path: resources/</div>
                  <el-button type="primary" icon="Plus" size="small" @click="openResourceDialog">Create resource</el-button>
                </div>
                <div class="text-xs text-gray-500 mb-3">
                  Put background info or glossary here. Reference <span class="font-mono">resources/filename</span> in markdown.
                </div>
                <el-table :data="resourceRows" style="width: 100%">
                  <el-table-column prop="name" label="File">
                    <template #default="scope">
                      <div class="flex items-center gap-2">
                        <el-icon><Document /></el-icon>
                        <span>{{ scope.row.name }}</span>
                      </div>
                    </template>
                  </el-table-column>
                  <el-table-column label="Actions" width="180">
                    <template #default="scope">
                      <el-button type="primary" link icon="Edit" @click="openResourceEditor(scope.row.name)">Edit</el-button>
                      <el-button type="primary" link @click="insertFileSnippet('resource', scope.row.name)">Insert</el-button>
                    </template>
                  </el-table-column>
                </el-table>
                <el-empty v-if="resourceRows.length === 0" description="No resources" />
              </el-tab-pane>

              <el-tab-pane label="References" name="references">
                <div class="flex justify-between items-center mb-4 mt-4">
                  <div class="text-sm text-gray-500 bg-gray-50 dark:bg-gray-800 px-3 py-1 rounded">Path: references/</div>
                  <el-button type="primary" icon="Plus" size="small" @click="openReferenceDialog">Create reference</el-button>
                </div>
                <div class="text-xs text-gray-500 mb-3">
                  Put standards/rules here. Reference <span class="font-mono">references/filename</span> in markdown.
                </div>
                <el-table :data="referenceRows" style="width: 100%">
                  <el-table-column prop="name" label="File">
                    <template #default="scope">
                      <div class="flex items-center gap-2">
                        <el-icon><Document /></el-icon>
                        <span>{{ scope.row.name }}</span>
                      </div>
                    </template>
                  </el-table-column>
                  <el-table-column label="Actions" width="180">
                    <template #default="scope">
                      <el-button type="primary" link icon="Edit" @click="openReferenceEditor(scope.row.name)">Edit</el-button>
                      <el-button type="primary" link @click="insertFileSnippet('reference', scope.row.name)">Insert</el-button>
                    </template>
                  </el-table-column>
                </el-table>
                <el-empty v-if="referenceRows.length === 0" description="No references" />
              </el-tab-pane>

              <el-tab-pane label="Templates" name="templates">
                <div class="flex justify-between items-center mb-4 mt-4">
                  <div class="text-sm text-gray-500 bg-gray-50 dark:bg-gray-800 px-3 py-1 rounded">Path: templates/</div>
                  <el-button type="primary" icon="Plus" size="small" @click="openTemplateDialog">Create template</el-button>
                </div>
                <div class="text-xs text-gray-500 mb-3">
                  Put output structure / code skeleton here. Reference <span class="font-mono">templates/filename</span> in markdown.
                </div>
                <el-table :data="templateRows" style="width: 100%">
                  <el-table-column prop="name" label="File">
                    <template #default="scope">
                      <div class="flex items-center gap-2">
                        <el-icon><Document /></el-icon>
                        <span>{{ scope.row.name }}</span>
                      </div>
                    </template>
                  </el-table-column>
                  <el-table-column label="Actions" width="180">
                    <template #default="scope">
                      <el-button type="primary" link icon="Edit" @click="openTemplateEditor(scope.row.name)">Edit</el-button>
                      <el-button type="primary" link @click="insertFileSnippet('template', scope.row.name)">Insert</el-button>
                    </template>
                  </el-table-column>
                </el-table>
                <el-empty v-if="templateRows.length === 0" description="No templates" />
              </el-tab-pane>
            </el-tabs>
          </template>
        </el-card>
      </el-col>
    </el-row>

    <el-dialog v-model="createDialogVisible" title="New skill" width="420px">
      <el-form :model="newSkill" label-width="100px">
        <el-form-item label="Skill name">
          <el-input v-model="newSkill.name" placeholder="e.g. code-comment-expert" />
          <div class="text-xs text-gray-400 mt-1">Lowercase letters/numbers/kebab-case only.</div>
        </el-form-item>
        <el-form-item label="Description">
          <el-input v-model="newSkill.description" placeholder="e.g. Improve code comments & readability for review/refactor" />
          <div class="text-xs text-gray-400 mt-1">Be specific about triggers and keywords.</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">Cancel</el-button>
        <el-button type="primary" @click="createSkill">Create</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="scriptDialogVisible" title="Create script" width="420px">
      <el-form :model="newScript" label-width="100px">
        <el-form-item label="Script type">
          <el-select v-model="newScript.type" placeholder="Select type">
            <el-option label="Python (.py)" value="py" />
            <el-option label="JavaScript (.js)" value="js" />
            <el-option label="Shell (.sh)" value="sh" />
          </el-select>
        </el-form-item>
        <el-form-item label="Filename">
          <el-input v-model="newScript.name" placeholder="e.g. lint" />
          <div class="text-xs text-gray-400 mt-1">No extension needed; it will be added automatically.</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="scriptDialogVisible = false">Cancel</el-button>
        <el-button type="primary" @click="createScript">Create</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="resourceDialogVisible" title="Create resource" width="420px">
      <el-form :model="newResource" label-width="100px">
        <el-form-item label="Filename">
          <el-input v-model="newResource.name" placeholder="e.g. glossary" />
          <div class="text-xs text-gray-400 mt-1">Will auto-append .md.</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="resourceDialogVisible = false">Cancel</el-button>
        <el-button type="primary" @click="createResource">Create</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="referenceDialogVisible" title="Create reference" width="420px">
      <el-form :model="newReference" label-width="100px">
        <el-form-item label="Filename">
          <el-input v-model="newReference.name" placeholder="e.g. style-guide" />
          <div class="text-xs text-gray-400 mt-1">Will auto-append .md.</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="referenceDialogVisible = false">Cancel</el-button>
        <el-button type="primary" @click="createReference">Create</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="templateDialogVisible" title="Create template" width="420px">
      <el-form :model="newTemplate" label-width="100px">
        <el-form-item label="Filename">
          <el-input v-model="newTemplate.name" placeholder="e.g. output-structure" />
          <div class="text-xs text-gray-400 mt-1">Will auto-append .md.</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="templateDialogVisible = false">Cancel</el-button>
        <el-button type="primary" @click="createTemplate">Create</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="downloadTargetDialogVisible" title="Choose download target" width="420px">
      <el-form label-width="90px">
        <el-form-item label="Download to">
          <el-select v-model="downloadTarget" placeholder="Select tool" class="w-full">
            <el-option
              v-for="item in downloadTargetOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <div class="text-xs text-gray-500">
          You can download to one AI tool, or choose "All tools".
        </div>
      </el-form>
      <template #footer>
        <el-button @click="closeDownloadTargetDialog">Cancel</el-button>
        <el-button type="primary" @click="confirmDownloadSkill">Download</el-button>
      </template>
    </el-dialog>

    <el-drawer v-model="editorVisible" size="70%" destroy-on-close :with-header="false">
      <div class="h-full flex flex-col p-4">
        <div class="flex justify-between items-center mb-4">
          <div class="text-lg font-bold flex items-center gap-2">
            <el-icon><Edit /></el-icon>
            {{ editorTitle }}
          </div>
          <div class="flex gap-2">
            <el-button @click="editorVisible = false">Cancel</el-button>
            <el-button type="primary" icon="Check" @click="saveEditor">Save</el-button>
          </div>
        </div>
        <div class="flex-1 overflow-hidden border border-gray-200 dark:border-gray-700 rounded-md shadow-inner">
          <v-ace-editor
            v-model:value="editorContent"
            :lang="editorLang"
            theme="github_dark"
            class="w-full h-full"
            :options="{ showPrintMargin: false, fontSize: 14 }"
          />
        </div>
      </div>
    </el-drawer>

    <!-- Online Skills drawer -->
    <el-drawer
      v-model="onlineDrawerVisible"
      size="90%"
      :show-close="false"
      destroy-on-close
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">Online Skills</span>
          <el-button @click="onlineDrawerVisible = false">Close</el-button>
        </div>
      </template>
      <div class="mb-4">
        <el-form :inline="true" :model="onlineSearchInfo">
          <el-form-item label="Name">
            <el-input v-model="onlineSearchInfo.name" placeholder="Search by name" clearable @keyup.enter="searchOnlineSkills" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" icon="Search" @click="searchOnlineSkills">Search</el-button>
            <el-button icon="Refresh" @click="resetOnlineSearch">Reset</el-button>
          </el-form-item>
        </el-form>
      </div>
      <el-table v-loading="onlineLoading" :data="onlineSkillList" stripe>
        <el-table-column label="Cover" width="80">
          <template #default="{ row }">
            <el-image
              v-if="row.picture"
              :src="row.picture"
              style="width: 50px; height: 50px"
              fit="cover"
              class="rounded"
            />
          </template>
        </el-table-column>
        <el-table-column label="Name" prop="name" min-width="160" show-overflow-tooltip>
          <template #default="{ row }">
            <a
              class="text-blue-500 hover:text-blue-700 cursor-pointer"
              :href="`https://plugin.gin-vue-admin.com/details/${row.ID}`"
              target="_blank"
            >{{ row.name }}</a>
          </template>
        </el-table-column>
        <el-table-column label="Summary" prop="resume" min-width="240" show-overflow-tooltip />
        <el-table-column label="Version" prop="actVersion" width="100" />
        <el-table-column label="Downloads" prop="downloadCount" width="90" />
        <el-table-column label="Actions" width="120" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="row.money === 0"
              type="primary"
              link
              icon="Download"
              :loading="downloadingIds.has(row.ID)"
              @click="handleDownloadSkill(row)"
            >Download</el-button>
            <a
              v-else
              class="text-blue-500 hover:text-blue-700 text-sm"
              :href="`https://plugin.gin-vue-admin.com/details/${row.ID}`"
              target="_blank"
            >Buy</a>
          </template>
        </el-table-column>
      </el-table>
      <div class="flex justify-center mt-4">
        <el-pagination
          :current-page="onlineSearchInfo.page"
          :page-size="onlineSearchInfo.pageSize"
          :page-sizes="[10, 20, 50]"
          :total="onlineTotal"
          layout="total, sizes, prev, pager, next"
          @current-change="handleOnlinePageChange"
          @size-change="handleOnlineSizeChange"
        />
      </div>
    </el-drawer>
  </div>
</template>

<script setup>
  import { computed, onMounted, reactive, ref } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { QuestionFilled, Document, Plus, Search, Check, Edit } from '@element-plus/icons-vue'
  import WarningBar from '@/components/warningBar/warningBar.vue'
  import {
    getSkillTools,
    getSkillList,
    getSkillDetail,
    saveSkill,
    deleteSkill,
    createSkillScript,
    getSkillScript,
    saveSkillScript,
    createSkillResource,
    getSkillResource,
    saveSkillResource,
    createSkillReference,
    getSkillReference,
    saveSkillReference,
    createSkillTemplate,
    getSkillTemplate,
    saveSkillTemplate,
    getGlobalConstraint,
    saveGlobalConstraint,
    packageSkill,
    downloadOnlineSkill
  } from '@/api/skills'
  import { getShopPluginList } from '@/api/plugin/api'
  import { VAceEditor } from 'vue3-ace-editor'
  import 'ace-builds/src-noconflict/mode-javascript'
  import 'ace-builds/src-noconflict/mode-python'
  import 'ace-builds/src-noconflict/mode-sh'
  import 'ace-builds/src-noconflict/mode-markdown'
  import 'ace-builds/src-noconflict/theme-github_dark'

  defineOptions({
    name: 'Skills'
  })

  const tools = ref([
    { key: 'copilot', label: 'Copilot' },
    { key: 'claude', label: 'Claude' },
    { key: 'cursor', label: 'Cursor' },
    { key: 'trae', label: 'Trae' },
    { key: 'codex', label: 'Codex' }
  ])
  const activeTool = ref('claude')
  const skills = ref([])
  const activeSkill = ref('')
  const skillFilter = ref('')
  const activeTab = ref('config')
  const globalConstraintExists = ref(false)

  const toolDirMap = {
    copilot: '.aone_copilot',
    claude: '.claude',
    cursor: '.cursor',
    trae: '.trae',
    codex: '.codex'
  }

  const globalConstraintPath = computed(() => {
    if (!activeTool.value) return 'skills/README.md'
    const toolDir = toolDirMap[activeTool.value] || `.${activeTool.value}`
    return `${toolDir}/skills/README.md`
  })

  const form = reactive({
    name: '',
    description: '',
    allowedTools: '',
    context: '',
    agent: '',
    markdown: ''
  })

  const markdownPlaceholder =
    'Suggested structure: # Skill Title -> ## Instructions -> ## Examples -> ## Guidelines.\n' +
    'Reference templates/... , references/... , scripts/... in markdown when needed.\n\n' +
    'Example:\n# Code Comment Expert\n## Instructions\n- Describe goal, inputs, outputs, steps.\n\n## Examples\n- Input: ...\n- Output: ...\n\n## Guidelines\n- Constraints, format, quality bar.\n'

  const quickBlocks = [
    { label: 'Title', content: '\n# Skill Title\n' },
    { label: 'Instructions', content: '\n## Instructions\n- Describe what to do and how.\n' },
    { label: 'Examples', content: '\n## Examples\n- Input: ...\n- Output: ...\n' },
    { label: 'Guidelines', content: '\n## Guidelines\n- Constraints, quality bar, and caveats.\n' },
    { label: 'Output format', content: '\n## Output Format\n1. ...\n2. ...\n' },
    { label: 'Reference template', content: '\nWhen you need structure, see templates/your-template.md.\n' },
    { label: 'Reference doc', content: '\nFor standards/terminology, see references/your-reference.md.\n' },
    { label: 'Run script', content: '\nFor automation, run scripts/your-script.py "{input}".\n' }
  ]

  const scripts = ref([])
  const resources = ref([])
  const references = ref([])
  const templates = ref([])

  const scriptRows = computed(() => skillsFilesToRows(scripts.value))
  const resourceRows = computed(() => skillsFilesToRows(resources.value))
  const referenceRows = computed(() => skillsFilesToRows(references.value))
  const templateRows = computed(() => skillsFilesToRows(templates.value))

  const createDialogVisible = ref(false)
  const scriptDialogVisible = ref(false)
  const resourceDialogVisible = ref(false)
  const referenceDialogVisible = ref(false)
  const templateDialogVisible = ref(false)

  const newSkill = reactive({
    name: '',
    description: ''
  })

  const newScript = reactive({
    name: '',
    type: 'py'
  })

  const newResource = reactive({
    name: ''
  })

  const newReference = reactive({
    name: ''
  })

  const newTemplate = reactive({
    name: ''
  })

  const editorVisible = ref(false)
  const editorContent = ref('')
  const editorFileName = ref('')
  const editorType = ref('script')
  const editorLang = ref('text')

  const editorTitle = computed(() => {
    if (!editorFileName.value) {
      return editorType.value === 'constraint' ? 'Global constraint' : 'File editor'
    }
    if (editorType.value === 'script') return `Script: ${editorFileName.value}`
    if (editorType.value === 'resource') return `Resource: ${editorFileName.value}`
    if (editorType.value === 'reference') return `Reference: ${editorFileName.value}`
    if (editorType.value === 'template') return `Template: ${editorFileName.value}`
    if (editorType.value === 'constraint') return `Global constraint: ${editorFileName.value}`
    return `File editor: ${editorFileName.value}`
  })

  const filteredSkills = computed(() => {
    if (!skillFilter.value) return skills.value
    return skills.value.filter((item) => item.toLowerCase().includes(skillFilter.value.toLowerCase()))
  })

  onMounted(async () => {
    await loadTools()
    await loadSkills()
  })

  async function loadTools() {
    try {
      const res = await getSkillTools()
      if (res.code === 0 && res.data?.tools?.length) {
        tools.value = res.data.tools
        if (!tools.value.find((item) => item.key === activeTool.value)) {
          activeTool.value = tools.value[0]?.key || 'claude'
        }
      }
    } catch (e) {
      ElMessage.warning('Failed to load tools; using defaults.')
    }
  }

  async function loadSkills() {
    if (!activeTool.value) return
    try {
      const res = await getSkillList({ tool: activeTool.value })
      if (res.code === 0) {
        skills.value = res.data?.skills || []
      }
    } catch (e) {
      ElMessage.error('Failed to load skills')
    }
  }

  async function loadSkillDetail(skillName) {
    if (!activeTool.value || !skillName) return
    try {
      const res = await getSkillDetail({ tool: activeTool.value, skill: skillName })
      if (res.code === 0) {
        const detail = res.data?.detail
        activeSkill.value = detail?.skill || skillName
        form.name = detail?.meta?.name || skillName
        form.description = detail?.meta?.description || ''
        form.allowedTools = detail?.meta?.allowedTools || ''
        form.context = detail?.meta?.context || ''
        form.agent = detail?.meta?.agent || ''
        form.markdown = detail?.markdown || ''
        scripts.value = detail?.scripts || []
        resources.value = detail?.resources || []
        references.value = detail?.references || []
        templates.value = detail?.templates || []
      }
    } catch (e) {
      ElMessage.error('Failed to load skill detail')
    }
  }

  async function openGlobalConstraint() {
    if (!activeTool.value) {
      ElMessage.warning('Please select a tool first')
      return
    }
    try {
      const res = await getGlobalConstraint({ tool: activeTool.value })
      if (res.code === 0) {
        globalConstraintExists.value = !!res.data?.exists
        if (!globalConstraintExists.value) {
          ElMessage.info('README.md not found. It will be created on save.')
        }
        openEditor('constraint', 'README.md', res.data?.content || '')
      }
    } catch (e) {
      ElMessage.error('Failed to read global constraint')
    }
  }

  function resetDetail() {
    activeSkill.value = ''
    form.name = ''
    form.description = ''
    form.allowedTools = ''
    form.context = ''
    form.agent = ''
    form.markdown = ''
    scripts.value = []
    resources.value = []
    references.value = []
    templates.value = []
    activeTab.value = 'config'
  }

  function handleToolSelect(key) {
    activeTool.value = key
    resetDetail()
    globalConstraintExists.value = false
    loadSkills()
  }

  function handleSkillSelect(skillName) {
    loadSkillDetail(skillName)
  }

  async function handleDeleteSkill(skillName) {
    if (!activeTool.value || !skillName) return
    try {
      await ElMessageBox.confirm(
        `Delete skill "${skillName}"? This will also delete its scripts/resources/references/templates files.`,
        'Confirm delete',
        {
          confirmButtonText: 'Delete',
          cancelButtonText: 'Cancel',
          type: 'warning'
        }
      )
    } catch (e) {
      return
    }

    try {
      const res = await deleteSkill({ tool: activeTool.value, skill: skillName })
      if (res.code !== 0) {
        return
      }
      if (activeSkill.value === skillName) {
        resetDetail()
      }
      await loadSkills()
      ElMessage.success('Deleted')
    } catch (e) {
      ElMessage.error('Delete failed')
    }
  }

  function openCreateDialog() {
    newSkill.name = ''
    newSkill.description = ''
    createDialogVisible.value = true
  }

  async function createSkill() {
    if (!newSkill.name.trim()) {
      ElMessage.warning('Please enter a skill name')
      return
    }
    const payload = {
      tool: activeTool.value,
      skill: newSkill.name.trim(),
      meta: {
        name: newSkill.name.trim(),
        description: newSkill.description.trim() || 'Please add a description',
        allowedTools: 'Bash(gh *)',
        context: 'fork',
        agent: 'Explore'
      },
      markdown: defaultSkillTemplate()
    }
    try {
      const res = await saveSkill(payload)
      if (res.code === 0) {
        ElMessage.success('Created')
        createDialogVisible.value = false
        await loadSkills()
        await loadSkillDetail(payload.skill)
      }
    } catch (e) {
      ElMessage.error('Create failed')
    }
  }

  async function saveCurrentSkill() {
    if (!activeSkill.value) return
    if (!form.name.trim()) {
      ElMessage.warning('Name is required')
      return
    }
    const payload = {
      tool: activeTool.value,
      skill: activeSkill.value,
      meta: {
        name: form.name.trim(),
        description: form.description.trim(),
        allowedTools: form.allowedTools.trim(),
        context: form.context.trim(),
        agent: form.agent.trim()
      },
      markdown: form.markdown
    }

    let syncTools = []
    try {
      await ElMessageBox.confirm('Sync to other AI client tools?', 'Sync', {
        confirmButtonText: 'Sync',
        cancelButtonText: 'Only current',
        type: 'warning'
      })
      syncTools = tools.value
        .map((item) => item.key)
        .filter((key) => key && key !== activeTool.value)
    } catch (e) {
      syncTools = []
    }

    if (syncTools.length) {
      payload.syncTools = syncTools
    }

    try {
      const res = await saveSkill(payload)
      if (res.code === 0) {
        ElMessage.success('Saved')
      }
    } catch (e) {
      ElMessage.error('Save failed')
    }
  }

  function extractFileNameFromDisposition(disposition) {
    if (!disposition) return ''
    const utf8Match = disposition.match(/filename\*=UTF-8''([^;]+)/i)
    if (utf8Match?.[1]) {
      try {
        return decodeURIComponent(utf8Match[1])
      } catch (e) {
        return utf8Match[1]
      }
    }
    const normalMatch = disposition.match(/filename="?([^";]+)"?/i)
    return normalMatch?.[1] || ''
  }

  async function packageCurrentSkill() {
    if (!activeTool.value || !activeSkill.value) {
      ElMessage.warning('Please select a skill first')
      return
    }
    try {
      const res = await packageSkill({ tool: activeTool.value, skill: activeSkill.value })
      const blob = res instanceof Blob ? res : (res?.data instanceof Blob ? res.data : null)
      if (!blob) {
        ElMessage.error('Package failed')
        return
      }
      const contentType = String(res?.headers?.['content-type'] || blob.type || '').toLowerCase()
      const disposition = String(res?.headers?.['content-disposition'] || '')
      const isZipResponse = contentType.includes('application/zip') || disposition.toLowerCase().includes('filename=')
      const isErrorBlob = contentType.includes('application/json') || contentType.includes('text/plain')
      if (!isZipResponse || isErrorBlob) {
        let msg = 'Package failed'
        try {
          const text = await blob.text()
          if (text) {
            try {
              const json = JSON.parse(text)
              msg = json?.msg || msg
            } catch (e) {
              msg = text
            }
          }
        } catch (e) {
          // ignore parse error
        }
        ElMessage.error(msg)
        return
      }
      const fileName = extractFileNameFromDisposition(disposition) || `${activeSkill.value}.zip`
      const url = window.URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = url
      link.download = fileName
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      window.URL.revokeObjectURL(url)
      ElMessage.success('Packaged')
    } catch (e) {
      ElMessage.error('Package failed')
    }
  }

  function appendMarkdown(content) {
    form.markdown = `${form.markdown || ''}${content}`
  }

  function insertFileSnippet(kind, fileName) {
    if (!fileName) return
    let snippet = ''
    switch (kind) {
      case 'script':
        snippet = `For automation, run scripts/${fileName} "{input}".`
        break
      case 'resource':
        snippet = `Background info: resources/${fileName}.`
        break
      case 'reference':
        snippet = `Follow: references/${fileName}.`
        break
      case 'template':
        snippet = `Output structure: templates/${fileName}.`
        break
      default:
        snippet = ''
    }
    if (!snippet) return
    appendMarkdown(`\n${snippet}\n`)
    ElMessage.success('Inserted into SKILL.md')
    activeTab.value = 'config'
  }

  function insertFullTemplate() {
    if (!form.markdown.trim()) {
      form.markdown = defaultSkillTemplate()
      return
    }
    form.markdown = `${form.markdown}\n${defaultSkillTemplate()}`
  }

  function openScriptDialog() {
    if (!activeSkill.value) {
      ElMessage.warning('Please select a skill first')
      return
    }
    newScript.name = ''
    newScript.type = 'py'
    scriptDialogVisible.value = true
  }

  async function createScript() {
    if (!newScript.name.trim()) {
      ElMessage.warning('Please enter a script filename')
      return
    }
    try {
      const res = await createSkillScript({
        tool: activeTool.value,
        skill: activeSkill.value,
        fileName: newScript.name.trim(),
        scriptType: newScript.type
      })
      if (res.code === 0) {
        scriptDialogVisible.value = false
        await loadSkillDetail(activeSkill.value)
        openEditor('script', res.data.fileName, res.data.content)
      }
    } catch (e) {
      ElMessage.error('Create script failed')
    }
  }

  async function openScriptEditor(fileName) {
    if (!fileName) return
    try {
      const res = await getSkillScript({
        tool: activeTool.value,
        skill: activeSkill.value,
        fileName
      })
      if (res.code === 0) {
        openEditor('script', fileName, res.data.content)
      }
    } catch (e) {
      ElMessage.error('Failed to read script')
    }
  }

  function openResourceDialog() {
    if (!activeSkill.value) {
      ElMessage.warning('Please select a skill first')
      return
    }
    newResource.name = ''
    resourceDialogVisible.value = true
  }

  async function createResource() {
    if (!newResource.name.trim()) {
      ElMessage.warning('Please enter a resource filename')
      return
    }
    try {
      const res = await createSkillResource({
        tool: activeTool.value,
        skill: activeSkill.value,
        fileName: newResource.name.trim()
      })
      if (res.code === 0) {
        resourceDialogVisible.value = false
        await loadSkillDetail(activeSkill.value)
        openEditor('resource', res.data.fileName, res.data.content)
      }
    } catch (e) {
      ElMessage.error('Create resource failed')
    }
  }

  async function openResourceEditor(fileName) {
    if (!fileName) return
    try {
      const res = await getSkillResource({
        tool: activeTool.value,
        skill: activeSkill.value,
        fileName
      })
      if (res.code === 0) {
        openEditor('resource', fileName, res.data.content)
      }
    } catch (e) {
      ElMessage.error('Failed to read resource')
    }
  }

  function openReferenceDialog() {
    if (!activeSkill.value) {
      ElMessage.warning('Please select a skill first')
      return
    }
    newReference.name = ''
    referenceDialogVisible.value = true
  }

  async function createReference() {
    if (!newReference.name.trim()) {
      ElMessage.warning('Please enter a reference filename')
      return
    }
    try {
      const res = await createSkillReference({
        tool: activeTool.value,
        skill: activeSkill.value,
        fileName: newReference.name.trim()
      })
      if (res.code === 0) {
        referenceDialogVisible.value = false
        await loadSkillDetail(activeSkill.value)
        openEditor('reference', res.data.fileName, res.data.content)
      }
    } catch (e) {
      ElMessage.error('Create reference failed')
    }
  }

  async function openReferenceEditor(fileName) {
    if (!fileName) return
    try {
      const res = await getSkillReference({
        tool: activeTool.value,
        skill: activeSkill.value,
        fileName
      })
      if (res.code === 0) {
        openEditor('reference', fileName, res.data.content)
      }
    } catch (e) {
      ElMessage.error('Failed to read reference')
    }
  }

  function openTemplateDialog() {
    if (!activeSkill.value) {
      ElMessage.warning('Please select a skill first')
      return
    }
    newTemplate.name = ''
    templateDialogVisible.value = true
  }

  async function createTemplate() {
    if (!newTemplate.name.trim()) {
      ElMessage.warning('Please enter a template filename')
      return
    }
    try {
      const res = await createSkillTemplate({
        tool: activeTool.value,
        skill: activeSkill.value,
        fileName: newTemplate.name.trim()
      })
      if (res.code === 0) {
        templateDialogVisible.value = false
        await loadSkillDetail(activeSkill.value)
        openEditor('template', res.data.fileName, res.data.content)
      }
    } catch (e) {
      ElMessage.error('Create template failed')
    }
  }

  async function openTemplateEditor(fileName) {
    if (!fileName) return
    try {
      const res = await getSkillTemplate({
        tool: activeTool.value,
        skill: activeSkill.value,
        fileName
      })
      if (res.code === 0) {
        openEditor('template', fileName, res.data.content)
      }
    } catch (e) {
      ElMessage.error('Failed to read template')
    }
  }

  function openEditor(type, fileName, content) {
    editorType.value = type
    editorFileName.value = fileName
    editorContent.value = content || ''
    editorLang.value = detectLang(fileName)
    editorVisible.value = true
  }

  async function saveEditor() {
    if (!editorFileName.value) return
    try {
      if (editorType.value === 'script') {
        const res = await saveSkillScript({
          tool: activeTool.value,
          skill: activeSkill.value,
          fileName: editorFileName.value,
          content: editorContent.value
        })
        if (res.code === 0) {
          ElMessage.success('Saved')
        }
      } else if (editorType.value === 'resource') {
        const res = await saveSkillResource({
          tool: activeTool.value,
          skill: activeSkill.value,
          fileName: editorFileName.value,
          content: editorContent.value
        })
        if (res.code === 0) {
          ElMessage.success('Saved')
        }
      } else if (editorType.value === 'reference') {
        const res = await saveSkillReference({
          tool: activeTool.value,
          skill: activeSkill.value,
          fileName: editorFileName.value,
          content: editorContent.value
        })
        if (res.code === 0) {
          ElMessage.success('Saved')
        }
      } else if (editorType.value === 'template') {
        const res = await saveSkillTemplate({
          tool: activeTool.value,
          skill: activeSkill.value,
          fileName: editorFileName.value,
          content: editorContent.value
        })
        if (res.code === 0) {
          ElMessage.success('Saved')
        }
      } else if (editorType.value === 'constraint') {
        let syncTools = []
        if (tools.value.length > 1) {
          try {
            await ElMessageBox.confirm('Sync to other AI client tools?', 'Sync', {
              confirmButtonText: 'Sync',
              cancelButtonText: 'Only current',
              type: 'warning'
            })
            syncTools = tools.value
              .map((item) => item.key)
              .filter((key) => key && key !== activeTool.value)
          } catch (e) {
            syncTools = []
          }
        }

        const res = await saveGlobalConstraint({
          tool: activeTool.value,
          content: editorContent.value,
          syncTools
        })
        if (res.code !== 0) {
          ElMessage.error('Save failed')
          return
        }
        globalConstraintExists.value = true
        ElMessage.success(syncTools.length ? 'Saved and synced' : 'Saved')
      }
    } catch (e) {
      ElMessage.error('Save failed')
    }
  }

  function detectLang(fileName) {
    if (!fileName) return 'text'
    const lower = fileName.toLowerCase()
    if (lower.endsWith('.py')) return 'python'
    if (lower.endsWith('.js')) return 'javascript'
    if (lower.endsWith('.sh')) return 'sh'
    if (lower.endsWith('.md')) return 'markdown'
    return 'text'
  }

  function defaultSkillTemplate() {
    return (
      '# Skill Title\n' +
      '## Instructions\n- Describe goal, inputs, outputs, steps.\n\n' +
      '## Examples\n- Input: ...\n- Output: ...\n\n' +
      '## Guidelines\n- Constraints, format, quality bar.\n\n' +
      '## Output Format\n1. ...\n2. ...\n'
    )
  }

  function skillsFilesToRows(list) {
    return (list || []).map((name) => ({ name }))
  }

  // ===== Online Skills =====
  const onlineDrawerVisible = ref(false)
  const onlineSkillList = ref([])
  const onlineTotal = ref(0)
  const onlineSearchInfo = reactive({ page: 1, pageSize: 10, name: '' })
  const onlineLoading = ref(false)
  const downloadingIds = reactive(new Set())
  const downloadTargetDialogVisible = ref(false)
  const downloadTarget = ref('')
  const downloadRow = ref(null)

  const ALL_TOOLS_DOWNLOAD_TARGET = '__all__'

  const downloadTargetOptions = computed(() => {
    const options = tools.value.map((item) => ({
      label: item.label || item.key,
      value: item.key
    }))
    options.push({
      label: 'All tools',
      value: ALL_TOOLS_DOWNLOAD_TARGET
    })
    return options
  })

  const pluginMarketLoginURL = 'https://plugin.gin-vue-admin.com'

  const isPluginMarketAuthError = (message) => {
    const msg = (message || '').toString()
    return msg.includes('plugin market login') || msg.includes('401')
  }

  const promptPluginMarketLogin = async () => {
    try {
      await ElMessageBox.confirm('Please log in to the plugin market before downloading. Go to login now?', 'Notice', {
        confirmButtonText: 'Go to plugin market',
        cancelButtonText: 'Cancel',
        type: 'warning'
      })
      window.open(pluginMarketLoginURL, '_blank')
    } catch (e) {
      // No extra notice on cancel
    }
  }

  const openOnlineDrawer = () => {
    onlineSearchInfo.page = 1
    onlineSearchInfo.pageSize = 10
    onlineSearchInfo.name = ''
    onlineDrawerVisible.value = true
    getOnlineSkills()
  }

  const getOnlineSkills = async () => {
    onlineLoading.value = true
    const res = await getShopPluginList({
      page: onlineSearchInfo.page,
      pageSize: onlineSearchInfo.pageSize,
      category: 6,
      name: onlineSearchInfo.name || undefined,
      updateTime: 1
    })
    if (res.code === 0) {
      onlineSkillList.value = res.data.list
      onlineTotal.value = res.data.total
    }
    onlineLoading.value = false
  }

  const searchOnlineSkills = () => {
    onlineSearchInfo.page = 1
    getOnlineSkills()
  }

  const resetOnlineSearch = () => {
    onlineSearchInfo.name = ''
    onlineSearchInfo.page = 1
    getOnlineSkills()
  }

  const handleOnlinePageChange = (page) => {
    onlineSearchInfo.page = page
    getOnlineSkills()
  }

  const handleOnlineSizeChange = (size) => {
    onlineSearchInfo.pageSize = size
    onlineSearchInfo.page = 1
    getOnlineSkills()
  }

  const getToolLabel = (key) => {
    return tools.value.find((item) => item.key === key)?.label || key
  }

  const closeDownloadTargetDialog = () => {
    downloadTargetDialogVisible.value = false
    downloadRow.value = null
  }

  const handleDownloadSkill = (row) => {
    downloadRow.value = row
    downloadTarget.value = activeTool.value || tools.value[0]?.key || ''
    downloadTargetDialogVisible.value = true
  }

  const confirmDownloadSkill = async () => {
    if (!downloadRow.value) {
      ElMessage.warning('No pending download found')
      return
    }
    const targetTools = downloadTarget.value === ALL_TOOLS_DOWNLOAD_TARGET
      ? tools.value.map((item) => item.key).filter(Boolean)
      : [downloadTarget.value].filter(Boolean)
    if (!targetTools.length) {
      ElMessage.warning('Please select a download target')
      return
    }

    const row = downloadRow.value
    closeDownloadTargetDialog()
    downloadingIds.add(row.ID)
    const successTools = []
    const failedTools = []
    try {
      for (const tool of targetTools) {
        try {
          const res = await downloadOnlineSkill({ tool, id: row.ID, version: row.actVersion })
          if (res.code === 0) {
            successTools.push(tool)
            continue
          }
          if (isPluginMarketAuthError(res.msg)) {
            await promptPluginMarketLogin()
            return
          }
          failedTools.push(`${getToolLabel(tool)}: ${res.msg || 'Download failed'}`)
        } catch (e) {
          const msg = e?.response?.data?.msg || e?.message || ''
          if (e?.response?.status === 401 || isPluginMarketAuthError(msg)) {
            await promptPluginMarketLogin()
            return
          }
          failedTools.push(`${getToolLabel(tool)}: Download failed`)
        }
      }

      if (successTools.includes(activeTool.value)) {
        await loadSkills()
      }

      if (failedTools.length === 0) {
        const successLabels = successTools.map((tool) => getToolLabel(tool)).join(', ')
        ElMessage({
          type: 'success',
          message: targetTools.length > 1 ? `${row.name} downloaded to: ${successLabels}` : `${row.name} downloaded`
        })
        return
      }

      if (successTools.length === 0) {
        ElMessage({ type: 'error', message: failedTools[0] || 'Download failed, please retry' })
        return
      }
      const successLabels = successTools.map((tool) => getToolLabel(tool)).join(', ')
      ElMessage({
        type: 'warning',
        message: `${row.name} partially downloaded. Success: ${successLabels}; Failed: ${failedTools.join('; ')}`
      })
    } finally {
      downloadingIds.delete(row.ID)
    }
  }
</script>

