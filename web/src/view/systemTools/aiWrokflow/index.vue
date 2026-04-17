<template>
  <div class="gva-table-box ai-workflow-page space-y-4">
    <el-card shadow="never">
      <div class="flex flex-wrap items-start justify-between gap-4">
        <div class="max-w-3xl">
          <h1 class="mt-2 text-2xl font-semibold text-slate-900">
            AI Requirement Analysis & Prompt Workflow
          </h1>
          <p class="mt-3 text-sm leading-6 text-slate-500">
            Sessions are saved to the backend automatically. After refreshing the page, you can view history,
            replay results by node, and rollback to any Assistant node to continue the conversation.
          </p>
        </div>
        <div class="flex flex-wrap">
          <el-button :icon="MagicStick" @click="fillExample"
            >Fill example</el-button
          >
        </div>
      </div>
    </el-card>

    <div class="grid grid-cols-1 gap-4 xl:grid-cols-[380px_minmax(0,1fr)]">
      <el-card shadow="never" class="self-start">
        <template #header>
          <div class="flex items-center justify-between gap-3">
            <div>
              <p class="text-base font-semibold text-slate-800">Input & history</p>
              <p class="mt-1 text-xs text-slate-500">
                {{
                  currentSession.id
                    ? `Session #${currentSession.id}`
                    : 'Unsaved new session'
                }}
              </p>
            </div>
            <el-tag
              effect="plain"
              :type="activeTab === 'analysis' ? 'primary' : 'success'"
            >
              {{ activeTab === 'analysis' ? 'Current: Analysis' : 'Current: Workflow' }}
            </el-tag>
          </div>
        </template>

        <div class="space-y-4">
          <div
            class="space-y-3 rounded border border-slate-200 bg-white p-3 shadow-sm"
          >
            <div
              class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between"
            >
              <div>
                <p class="text-sm font-semibold text-slate-800">
                  Choose a mode first
                </p>
                <p class="mt-1 text-xs text-slate-500">
                  Many users skip this. You can switch by clicking the cards below.
                </p>
              </div>
              <span
                class="w-fit shrink-0 whitespace-nowrap rounded-full bg-slate-100 px-3 py-1 text-[11px] font-semibold text-slate-500"
                >Switch anytime</span
              >
            </div>
            <div class="grid grid-cols-1 gap-3 md:grid-cols-2 xl:grid-cols-1">
              <button
                type="button"
                class="rounded border px-4 py-4 text-left transition"
                :class="
                  activeTab === 'analysis'
                    ? 'border-sky-400 bg-sky-50 shadow-sm shadow-sky-100'
                    : 'border-slate-200 bg-slate-50 hover:border-slate-300 hover:bg-white'
                "
                @click="switchTab('analysis')"
              >
                <div class="flex items-center justify-between gap-3">
                  <div>
                    <p
                      class="text-base font-semibold"
                      :class="
                        activeTab === 'analysis'
                          ? 'text-sky-700'
                          : 'text-slate-800'
                      "
                    >
                      Analysis
                    </p>
                    <p
                      class="mt-2 text-xs leading-5"
                      :class="
                        activeTab === 'analysis'
                          ? 'text-sky-700/80'
                          : 'text-slate-500'
                      "
                    >
                      Break down requirements into modules, fields, and open questions before implementation.
                    </p>
                  </div>
                  <span
                    class="shrink-0 whitespace-nowrap rounded-full px-3 py-1 text-[11px] font-semibold leading-none"
                    :class="
                      activeTab === 'analysis'
                        ? 'bg-sky-600 text-white'
                        : 'bg-white text-slate-500 border border-slate-200'
                    "
                  >
                    {{ activeTab === 'analysis' ? 'Active' : 'Switch to this' }}
                  </span>
                </div>
              </button>

              <button
                type="button"
                class="rounded border px-4 py-4 text-left transition"
                :class="
                  activeTab === 'workflow'
                    ? 'border-emerald-400 bg-emerald-50 shadow-sm shadow-emerald-100'
                    : 'border-slate-200 bg-slate-50 hover:border-slate-300 hover:bg-white'
                "
                @click="switchTab('workflow')"
              >
                <div class="flex items-center justify-between gap-3">
                  <div>
                    <p
                      class="text-base font-semibold"
                      :class="
                        activeTab === 'workflow'
                          ? 'text-emerald-700'
                          : 'text-slate-800'
                      "
                    >
                      Workflow
                    </p>
                    <p
                      class="mt-2 text-xs leading-5"
                      :class="
                        activeTab === 'workflow'
                          ? 'text-emerald-700/80'
                          : 'text-slate-500'
                      "
                    >
                      Generate step-by-step prompts and an execution route from requirements or analysis results.
                    </p>
                  </div>
                  <span
                    class="shrink-0 whitespace-nowrap rounded-full px-3 py-1 text-[11px] font-semibold leading-none"
                    :class="
                      activeTab === 'workflow'
                        ? 'bg-emerald-600 text-white'
                        : 'bg-white text-slate-500 border border-slate-200'
                    "
                  >
                    {{ activeTab === 'workflow' ? 'Active' : 'Switch to this' }}
                  </span>
                </div>
              </button>
            </div>
          </div>

          <div
            class="space-y-3 rounded border border-slate-200 bg-slate-50 p-3"
          >
            <div class="flex items-center">
              <el-input
                v-model="historyKeywords[activeTab]"
                class="flex-1 mr-2"
                clearable
                placeholder="Search history"
                @keyup.enter="loadSessionList(activeTab)"
              />
              <el-button
                plain
                class="!px-3"
                :icon="RefreshRight"
                :loading="historyLoading"
                title="Refresh list"
                @click="loadSessionList(activeTab)"
              />
              <el-button
                plain
                class="!px-3"
                :icon="Plus"
                title="New session"
                @click="startNewConversation(activeTab)"
              />
            </div>
            <div
              v-if="currentSessionList.length"
              class="max-h-[280px] space-y-2 overflow-auto pr-1"
            >
              <div
                v-for="item in currentSessionList"
                :key="item.ID || item.id"
                class="cursor-pointer rounded border p-3"
                :class="
                  isSessionActive(item)
                    ? 'border-sky-300 bg-sky-50'
                    : 'border-slate-200 bg-white'
                "
                @click="openSession(item)"
              >
                <div class="flex items-start justify-between gap-3">
                  <div class="min-w-0 flex-1">
                    <div class="flex flex-wrap items-center gap-2">
                      <p class="truncate text-sm font-semibold text-slate-800">
                        {{ item.title || 'Untitled' }}
                      </p>
                      <el-tag size="small" effect="plain">{{
                        item.tab === 'analysis' ? 'Analysis' : 'Workflow'
                      }}</el-tag>
                    </div>
                    <p class="mt-2 text-xs leading-5 text-slate-500">
                      {{ item.summary || 'No summary yet' }}
                    </p>
                    <p class="mt-2 text-[11px] text-slate-400">
                      Updated: {{
                        formatTime(item.UpdatedAt || item.updatedAt)
                      }}
                    </p>
                  </div>
                  <el-button
                    link
                    type="danger"
                    :icon="Delete"
                    @click.stop="removeSession(item)"
                  />
                </div>
              </div>
            </div>
            <el-empty v-else description="No saved sessions for this tab yet." />
          </div>

          <el-form label-position="top">
            <template v-if="activeTab === 'analysis'">
              <el-form-item label="Requirement">
                <el-input
                  v-model="analysisForm.requirement"
                  type="textarea"
                  :rows="10"
                  maxlength="4000"
                  show-word-limit
                  placeholder="Example: Build a leave management system with employee requests, manager approvals, and HR reporting."
                />
              </el-form-item>
              <el-form-item label="Target type">
                <el-radio-group v-model="analysisForm.packageType">
                  <el-radio-button label="auto">Auto</el-radio-button>
                  <el-radio-button label="package">package</el-radio-button>
                  <el-radio-button label="plugin">plugin</el-radio-button>
                </el-radio-group>
              </el-form-item>
              <el-form-item label="Business scenario">
                <el-input
                  v-model="analysisForm.businessScene"
                  placeholder="Example: OA, ERP, internal admin"
                />
              </el-form-item>
              <el-form-item label="Extra constraints">
                <el-input
                  v-model="analysisForm.extraConstraints"
                  type="textarea"
                  :rows="4"
                  placeholder="Example: attachments, status dictionaries, export reports."
                />
              </el-form-item>
              <el-form-item label="Has client pages">
                <el-switch
                  v-model="analysisForm.hasClientPage"
                  active-text="Yes"
                  inactive-text="No"
                />
              </el-form-item>
              <template v-if="analysisForm.hasClientPage">
                <el-form-item label="Client page description">
                  <el-input
                    v-model="analysisForm.clientPageDescription"
                    type="textarea"
                    :rows="4"
                    placeholder="Example: an employee H5 page for submitting requests, tracking approval status, and uploading attachments."
                  />
                </el-form-item>
                <el-form-item label="Client page constraints">
                  <el-input
                    v-model="analysisForm.clientPageConstraints"
                    type="textarea"
                    :rows="4"
                    placeholder="Example: mobile-first, fewer fields, prominent approval status, attachments support images and PDF."
                  />
                </el-form-item>
              </template>
              <div class="flex flex-wrap">
                <el-button
                  type="primary"
                  :loading="analysisLoading"
                  :icon="MagicStick"
                  @click="runAnalysis"
                  >Run analysis</el-button
                >
                <el-button
                  :disabled="!hasAnalysisResult"
                  @click="pushAnalysisToWorkflow"
                  >Send to workflow</el-button
                >
                <el-button @click="startNewConversation('analysis')"
                  >Clear</el-button
                >
              </div>
            </template>

            <template v-else>
              <el-form-item label="Requirements or analysis result">
                <el-input
                  v-model="workflowForm.source"
                  type="textarea"
                  :rows="10"
                  maxlength="6000"
                  show-word-limit
                  placeholder="Paste requirements directly, or bring in the structured analysis result."
                />
              </el-form-item>
              <el-form-item label="Workflow goal">
                <el-select v-model="workflowForm.flowType" class="w-full">
                  <el-option label="GVA codegen" value="gva_codegen" />
                  <el-option label="GVA feature polish" value="gva_polish" />
                  <el-option label="MCP usage guide" value="mcp_assist" />
                </el-select>
              </el-form-item>
              <el-form-item label="Extra constraints">
                <el-input
                  v-model="workflowForm.extraConstraints"
                  type="textarea"
                  :rows="4"
                  placeholder="Example: provide copyable prompts for each step and explain expected outputs."
                />
              </el-form-item>
              <div class="flex flex-wrap">
                <el-button
                  type="primary"
                  :loading="workflowLoading"
                  :icon="MagicStick"
                  @click="runWorkflow"
                  >Generate</el-button
                >
                <el-button
                  :disabled="!workflowResult.steps.length"
                  :icon="DocumentCopy"
                  @click="copyAllWorkflowPrompts"
                  >Copy all</el-button
                >
                <el-button @click="startNewConversation('workflow')"
                  >Clear</el-button
                >
              </div>
            </template>

            <el-divider />

            <div class="space-y-3 rounded bg-slate-50 p-3">
              <div class="text-xs leading-6 text-slate-500">
                The current session is persisted to the backend. Click a history item to restore it. The latest session is auto-loaded after refresh.
              </div>
              <div class="flex flex-wrap">
                <el-button plain @click="startNewConversation(activeTab)"
                  >New session</el-button
                >
                <el-button
                  plain
                  :disabled="!currentSession.conversationId"
                  @click="
                    copyText(
                      currentSession.conversationId,
                      'Copied conversation_id'
                    )
                  "
                  >Copy session ID</el-button
                >
              </div>
              <div
                class="rounded-xl border border-dashed border-slate-300 bg-white px-3 py-2 text-xs text-slate-500"
              >
                conversation_id:
                <span class="break-all font-mono text-slate-700">{{
                  currentSession.conversationId ||
                  'This session has no conversation_id yet'
                }}</span>
              </div>
              <div
                class="rounded-xl border border-dashed border-slate-300 bg-white px-3 py-2 text-xs text-slate-500"
              >
                Current node:
                <span class="break-all font-mono text-slate-700">{{
                  currentSession.currentNodeId || 'Defaults to latest Assistant node'
                }}</span>
              </div>
            </div>

            <el-divider />

            <el-collapse>
              <el-collapse-item title="Extra passthrough payload" name="settings">
                <el-form-item label="Extra payload (JSON)">
                  <el-input
                    v-model="settings.extraPayload"
                    type="textarea"
                    :rows="4"
                    placeholder='Example: {"tenant":"default","provider":"dify"}'
                  />
                </el-form-item>
              </el-collapse-item>
            </el-collapse>
          </el-form>
        </div>
      </el-card>

      <el-card shadow="never" class="min-h-[760px]">
        <template #header>
          <div class="flex items-center justify-between gap-3">
            <div>
              <p class="text-base font-semibold text-slate-800">
                {{ activeTab === 'analysis' ? 'Analysis result' : 'Workflow result' }}
              </p>
              <p class="mt-1 text-xs text-slate-500">
                {{
                  activeTab === 'analysis'
                    ? 'Replay any version of the analysis result from session nodes.'
                    : 'Replay any version of the workflow result from session nodes.'
                }}
              </p>
            </div>
            <div class="flex gap-2">
              <el-button
                plain
                :loading="currentDumpLoading"
                :disabled="!canDumpCurrentMarkdown"
                @click="dumpCurrentMarkdown"
                >{{ activeTab === 'analysis' ? 'Dump analysis' : 'Dump prompts' }}</el-button
              >
              <el-button
                v-if="selectedAssistantMessage"
                plain
                @click="rollbackToMessage(selectedAssistantMessage.id)"
                >Rollback to current node</el-button
              >
              <el-button
                v-if="activeTab === 'analysis'"
                :disabled="!hasAnalysisResult"
                :icon="DocumentCopy"
                @click="copyAnalysisResult"
                >Copy structured result</el-button
              >
              <el-button
                v-else
                :disabled="!workflowResult.rawText"
                :icon="RefreshRight"
                @click="copyWorkflowRaw"
                >Copy raw result</el-button
              >
            </div>
          </div>
        </template>

        <div class="space-y-4">
          <div
            v-if="currentLoading"
            class="rounded border border-sky-200 bg-sky-50 px-4 py-3 text-sm text-sky-700"
          >
            <div class="mb-2 flex items-center gap-2">
              <el-icon class="animate-spin"><RefreshRight /></el-icon>
              <span class="font-semibold">AI thinking...</span>
            </div>
            <pre
              v-if="streamingPreviewText"
              ref="streamingPreviewRef"
              class="mt-2 max-h-[200px] overflow-auto whitespace-pre-wrap rounded bg-slate-950 p-3 text-xs leading-5 text-slate-100"
            >{{ streamingPreviewText }}</pre>
            <div v-else class="text-xs text-sky-500">Waiting for model response...</div>
          </div>
          <div class="rounded border border-slate-200 bg-slate-50 p-4">
            <div class="mb-3 flex items-center justify-between">
              <h3 class="text-sm font-semibold text-slate-800">Conversation</h3>
              <el-tag effect="plain"
                >{{ currentSession.messages.length }} messages</el-tag
              >
            </div>
            <div
              v-if="currentSession.messages.length"
              ref="conversationListRef"
              class="max-h-[320px] space-y-3 overflow-auto pr-1"
            >
              <div
                v-for="item in currentSessionDisplayMessages"
                :key="item.id"
                class="rounded px-4 py-3"
                :class="
                  item.role === 'user'
                    ? 'bg-sky-600 text-white'
                    : item.isSelected
                    ? 'border border-sky-300 bg-sky-50 text-slate-700'
                    : 'border border-slate-200 bg-white text-slate-700'
                "
              >
                <div class="mb-2 flex items-start justify-between gap-3">
                  <div>
                    <div
                      class="text-xs font-semibold uppercase tracking-[0.2em]"
                      :class="
                        item.role === 'user' ? 'text-sky-100' : 'text-slate-400'
                      "
                    >
                      {{ item.role === 'user' ? 'User' : 'Assistant' }}
                    </div>
                    <div
                      v-if="item.createdAt"
                      class="mt-1 text-[11px]"
                      :class="
                        item.role === 'user'
                          ? 'text-sky-100/80'
                          : 'text-slate-400'
                      "
                    >
                      {{ formatTime(item.createdAt) }}
                    </div>
                  </div>
                  <div v-if="item.role === 'assistant'" class="flex flex-wrap">
                    <el-button
                      link
                      type="primary"
                      @click.stop="selectMessageNode(item.id)"
                      >View node</el-button
                    >
                    <el-button
                      link
                      type="danger"
                      @click.stop="rollbackToMessage(item.id)"
                      >Rollback here</el-button
                    >
                  </div>
                </div>
                <div
                  v-if="item.role === 'user'"
                  class="whitespace-pre-wrap text-sm leading-6"
                >
                  {{ item.content }}
                </div>
                <div v-else class="space-y-3">
                  <template v-if="item.isStreaming">
                    <div
                      class="rounded-xl border border-sky-200 bg-sky-50 px-3 py-2 text-xs leading-5 text-sky-700"
                    >
                      Streaming response. Expanding content in real time...
                    </div>
                    <pre
                      class="overflow-auto whitespace-pre-wrap rounded bg-slate-950 p-4 text-xs leading-6 text-slate-100"
                      >{{ item.content || '...' }}</pre
                    >
                  </template>
                  <template v-else>
                  <div
                    class="rounded-xl border border-dashed border-slate-200 bg-white px-3 py-2 text-xs leading-5 text-slate-500"
                  >
                    {{
                      item.display.preview ||
                      'Assistant output is collapsed; expand to view when needed.'
                    }}
                  </div>
                  <div class="flex gap-2">
                    <el-tag
                      v-if="item.display.hasThink"
                      effect="plain"
                      type="warning"
                      class="cursor-pointer"
                      @click="openDrawer('Think', item.display.think)"
                    >Think</el-tag>
                    <el-tag
                      v-if="item.display.hasRawFile"
                      effect="plain"
                      class="cursor-pointer"
                      @click="openDrawer('Raw output', item.display.rawFile)"
                    >Raw output</el-tag>
                  </div>
                  </template>
                </div>
              </div>
            </div>
            <el-empty v-else description="No messages in this session yet." />
          </div>

          <div class="rounded border border-slate-200 bg-white p-4">
            <div class="mb-3 flex items-center justify-between">
              <h3 class="text-sm font-semibold text-slate-800">Follow up</h3>
              <span class="text-xs text-slate-400">{{
                currentSession.conversationId
                  ? 'Will reuse the current conversation_id'
                  : 'After rollback, a new session will start from the selected node'
              }}</span>
            </div>
            <el-input
              v-model="followUpInput"
              type="textarea"
              :rows="4"
              maxlength="3000"
              show-word-limit
              :placeholder="
                activeTab === 'analysis'
                  ? 'Example: add a return-from-leave flow and CC feature.'
                  : 'Example: break steps down further; add field confirmation and dictionary design.'
              "
            />
            <div class="mt-3 flex flex-wrap">
              <el-button
                type="primary"
                :loading="currentLoading"
                :disabled="!canSendFollowUp"
                @click="sendFollowUp"
                >Send follow-up</el-button
              >
              <el-button @click="followUpInput = ''">Clear input</el-button>
            </div>
          </div>

          <template v-if="activeTab === 'analysis'">
            <template v-if="hasAnalysisResult">
              <div class="rounded bg-slate-50 p-4">
                <p
                  class="text-xs font-semibold uppercase tracking-[0.2em] text-sky-600"
                >
                  Summary
                </p>
                <p
                  class="mt-2 whitespace-pre-wrap text-sm leading-6 text-slate-700"
                >
                  {{
                    analysisResult.summary || 'No summary parsed; keeping raw output.'
                  }}
                </p>
                <div class="mt-3 flex flex-wrap gap-2">
                  <el-tag
                    v-if="analysisResult.recommendedPackageType"
                    effect="plain"
                    type="primary"
                    >Recommended: {{ analysisResult.recommendedPackageType }}</el-tag
                  >
                  <el-tag effect="plain"
                    >Modules {{ analysisResult.modules.length }}</el-tag
                  >
                  <el-tag
                    v-if="analysisClientPages.length"
                    effect="plain"
                    type="success"
                    >Client pages {{ analysisClientPages.length }}</el-tag
                  >
                  <el-tag effect="plain"
                    >Open questions {{ analysisResult.missingInfo.length }}</el-tag
                  >
                </div>
              </div>
              <el-collapse class="mt-4">
                <el-collapse-item
                  v-if="analysisResult.missingInfo.length"
                  title="Missing info"
                  name="analysis-missing"
                  ><div class="flex flex-wrap gap-2">
                    <el-tag
                      v-for="item in analysisResult.missingInfo"
                      :key="item"
                      effect="plain"
                      type="warning"
                      >{{ item }}</el-tag
                    >
                  </div></el-collapse-item
                >
                <el-collapse-item
                  v-if="analysisResult.modules.length"
                  title="Module list"
                  name="analysis-modules"
                >
                  <div class="space-y-3">
                    <el-collapse>
                      <el-collapse-item
                        v-for="(module, index) in analysisResult.modules"
                        :key="`${module.name}-${index}`"
                        :name="`${module.name || index}`"
                        :title="`${index + 1}. ${module.label || module.name || `Module ${index + 1}`}`"
                      >
                        <div class="space-y-4">
                          <p class="text-sm leading-6 text-slate-600">
                            {{ module.description || 'No module description.' }}
                          </p>
                          <div
                            v-if="module.fields.length"
                            class="grid grid-cols-1 gap-3 md:grid-cols-2"
                          >
                            <div
                              v-for="field in module.fields"
                              :key="`${module.name}-${field.name}`"
                              class="rounded-2xl border border-slate-200 bg-slate-50 p-3"
                            >
                              <div
                                class="flex flex-wrap items-center justify-between gap-2"
                              >
                                <div>
                                  <p class="text-sm font-semibold text-slate-800">
                                    {{ field.label || field.name }}
                                  </p>
                                  <p class="mt-1 text-xs text-slate-500">
                                    {{ field.name }}
                                  </p>
                                </div>
                                <el-tag size="small" effect="plain">
                                  {{ field.type || 'string' }}
                                </el-tag>
                              </div>
                              <p class="mt-3 text-sm leading-6 text-slate-600">
                                {{ field.description || 'No field description.' }}
                              </p>
                              <div class="mt-3 flex flex-wrap gap-2">
                                <el-tag
                                  v-if="field.required"
                                  size="small"
                                  type="danger"
                                  effect="plain"
                                >
                                  Required
                                </el-tag>
                                <el-tag
                                  v-if="field.dictionary"
                                  size="small"
                                  type="success"
                                  effect="plain"
                                >
                                  Dictionary: {{ field.dictionary }}
                                </el-tag>
                                <el-tag
                                  v-if="field.relation"
                                  size="small"
                                  type="warning"
                                  effect="plain"
                                >
                                  Relation: {{ field.relation }}
                                </el-tag>
                              </div>
                            </div>
                          </div>
                        </div>
                      </el-collapse-item>
                    </el-collapse>
                  </div>
                </el-collapse-item>
                <el-collapse-item
                  v-if="analysisClientPages.length"
                  title="Client pages"
                  name="analysis-client-pages"
                >
                  <div class="space-y-3">
                    <el-collapse>
                      <el-collapse-item
                        v-for="(page, index) in analysisClientPages"
                        :key="`${page.name}-${index}`"
                        :name="`${page.name || index}`"
                        :title="`${index + 1}. ${page.label || page.name || `Page ${index + 1}`}`"
                      >
                        <div class="space-y-4">
                          <div class="flex flex-wrap gap-2">
                            <el-tag v-if="page.pageType" effect="plain" type="primary">
                              {{ page.pageType }}
                            </el-tag>
                            <el-tag
                              v-for="moduleName in page.targetModules"
                              :key="`${page.name}-${moduleName}`"
                              effect="plain"
                              type="success"
                            >
                              {{ moduleName }}
                            </el-tag>
                          </div>
                          <p class="text-sm leading-6 text-slate-600">
                            {{ page.description || 'No page description.' }}
                          </p>
                          <div
                            v-if="page.fields.length"
                            class="grid grid-cols-1 gap-3 md:grid-cols-2"
                          >
                            <div
                              v-for="field in page.fields"
                              :key="`${page.name}-${field.name}`"
                              class="rounded-2xl border border-slate-200 bg-slate-50 p-3"
                            >
                              <div class="flex flex-wrap items-center justify-between gap-2">
                                <div>
                                  <p class="text-sm font-semibold text-slate-800">
                                    {{ field.label || field.name }}
                                  </p>
                                  <p class="mt-1 text-xs text-slate-500">
                                    {{ field.name }}
                                  </p>
                                </div>
                                <el-tag size="small" effect="plain">
                                  {{ field.displayType || 'text' }}
                                </el-tag>
                              </div>
                              <p class="mt-3 text-sm leading-6 text-slate-600">
                                {{ field.description || 'No page field description.' }}
                              </p>
                              <div class="mt-3 flex flex-wrap gap-2">
                                <el-tag
                                  v-if="field.required"
                                  size="small"
                                  type="danger"
                                  effect="plain"
                                >
                                  Required
                                </el-tag>
                                <el-tag
                                  v-if="field.sourceModule || field.sourceField"
                                  size="small"
                                  type="info"
                                  effect="plain"
                                >
                                  {{ `${field.sourceModule || '-'}.${
                                    field.sourceField || '-'
                                  }` }}
                                </el-tag>
                              </div>
                            </div>
                          </div>
                          <div v-if="page.interactions.length" class="space-y-2">
                            <h4 class="text-sm font-semibold text-slate-800">Interactions</h4>
                            <div
                              v-for="item in page.interactions"
                              :key="`${page.name}-${item}`"
                              class="rounded border border-slate-200 bg-white px-4 py-3 text-sm leading-6 text-slate-600"
                            >
                              {{ item }}
                            </div>
                          </div>
                          <div v-if="page.relations.length" class="space-y-2">
                            <h4 class="text-sm font-semibold text-slate-800">Field mappings</h4>
                            <div
                              v-for="item in page.relations"
                              :key="`${page.name}-${item}`"
                              class="rounded border border-slate-200 bg-white px-4 py-3 text-sm leading-6 text-slate-600"
                            >
                              {{ item }}
                            </div>
                          </div>
                        </div>
                      </el-collapse-item>
                    </el-collapse>
                  </div>
                </el-collapse-item>
                <el-collapse-item
                  v-if="analysisResult.suggestions.length"
                  title="Implementation suggestions"
                  name="analysis-suggestions"
                  ><div class="space-y-2">
                    <div
                      v-for="item in analysisResult.suggestions"
                      :key="item"
                      class="rounded border border-slate-200 bg-white px-4 py-3 text-sm leading-6 text-slate-600"
                    >
                      {{ item }}
                    </div>
                  </div></el-collapse-item
                >
                </el-collapse>
            </template>
            <el-empty
              v-else
              description="Select a session node to view the corresponding version of the analysis result."
            />
          </template>

          <template v-else>
            <template
              v-if="workflowResult.steps.length || workflowResult.rawText"
            >
              <div class="rounded bg-slate-50 p-4">
                <p
                  class="text-xs font-semibold uppercase tracking-[0.2em] text-sky-600"
                >
                  Workflow Summary
                </p>
                <p
                  class="mt-2 whitespace-pre-wrap text-sm leading-6 text-slate-700"
                >
                  {{ workflowResult.summary || 'Prompt workflow generated.' }}
                </p>
                <div class="mt-3 flex flex-wrap gap-2">
                  <el-tag effect="plain"
                    >Steps {{ workflowResult.steps.length }}</el-tag
                  >
                </div>
              </div>
              <div v-if="workflowResult.steps.length" class="mt-4 space-y-3">
                <div
                  v-for="(step, index) in workflowResult.steps"
                  :key="`${step.title}-${index}`"
                  class="rounded border border-slate-200 bg-white p-4 shadow-sm shadow-slate-100"
                >
                  <div class="flex flex-wrap items-start justify-between gap-3">
                    <div class="flex items-start gap-3">
                      <div
                        class="flex h-9 w-9 items-center justify-center rounded-full bg-sky-600 text-sm font-semibold text-white"
                      >
                        {{ index + 1 }}
                      </div>
                      <div>
                        <h3 class="text-base font-semibold text-slate-800">
                          {{ step.title || `Step ${index + 1}` }}
                        </h3>
                        <p class="mt-1 text-sm leading-6 text-slate-500">
                          {{ step.goal || 'Run this prompt sequence in order.' }}
                        </p>
                      </div>
                    </div>
                    <div class="flex flex-wrap gap-2">
                      <el-tag
                        v-if="step.suggestedTool"
                        effect="plain"
                        type="success"
                        >{{ step.suggestedTool }}</el-tag
                      >
                      <el-tag
                        :type="step.autoExecutable ? 'danger' : 'info'"
                        effect="plain"
                        >{{
                          step.autoExecutable ? 'Auto-executable' : 'Manual review recommended'
                        }}</el-tag
                      >
                    </div>
                  </div>
                  <div class="mt-4 rounded bg-slate-950 p-4">
                    <div
                      class="mb-2 text-xs font-semibold uppercase tracking-[0.18em] text-slate-400"
                    >
                      Prompt
                    </div>
                    <pre
                      class="max-h-[260px] overflow-auto whitespace-pre-wrap text-xs leading-6 text-slate-100"
                      >{{ step.prompt || 'Model did not return a prompt for this step.' }}</pre
                    >
                  </div>
                  <div
                    v-if="step.expectedOutput"
                    class="mt-4 rounded border border-dashed border-slate-300 bg-slate-50 p-4"
                  >
                    <div
                      class="mb-2 text-xs font-semibold uppercase tracking-[0.18em] text-slate-400"
                    >
                      Expected Output
                    </div>
                    <p
                      class="whitespace-pre-wrap text-sm leading-6 text-slate-600"
                    >
                      {{ step.expectedOutput }}
                    </p>
                  </div>
                </div>
              </div>
            </template>
            <el-empty
              v-else
              description="Select a session node to view the corresponding version of the workflow result."
            />
          </template>
        </div>
      </el-card>
    </div>
    <el-drawer
      v-model="drawerVisible"
      :title="drawerTitle"
      size="50%"
      destroy-on-close
    >
      <pre
        class="h-full overflow-auto whitespace-pre-wrap rounded bg-slate-950 p-4 text-xs leading-6 text-slate-100"
      >{{ drawerContent }}</pre>
    </el-drawer>
  </div>
</template>

<script setup>
import { computed, nextTick, onMounted, reactive, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useI18n } from 'vue-i18n'
import {
  Delete,
  DocumentCopy,
  MagicStick,
  Plus,
  Promotion,
  RefreshRight
} from '@element-plus/icons-vue'

const { t } = useI18n()
import {
  analyzeRequirementByAISSEStream,
  dumpAIWorkflowMarkdown,
  deleteAIWorkflowSession,
  generatePromptFlowByAISSEStream,
  getAIWorkflowSessionDetail,
  getAIWorkflowSessionList,
  saveAIWorkflowSession
} from '@/api/autoCode'
import { useUserStore } from '@/pinia/modules/user'

defineOptions({ name: 'AIWorkflow' })

const router = useRouter()
const userStore = useUserStore()
const SETTINGS_KEY = 'gva_ai_workflow_settings'
const ACTIVE_SESSION_KEY = 'gva_ai_workflow_active_session_ids'
const TAB_MODE_MAP = {
  analysis: 'analysisChat',
  workflow: 'workflowPromptChat'
}
const FLOW_TYPE_LABEL_MAP = {
  gva_codegen: 'GVA codegen',
  gva_polish: 'GVA feature polish',
  mcp_assist: 'MCP usage guide'
}
const defaultSettings = { extraPayload: '' }
const newAnalysisForm = () => ({
  requirement: '',
  packageType: 'auto',
  businessScene: '',
  extraConstraints: '',
  hasClientPage: false,
  clientPageDescription: '',
  clientPageConstraints: ''
})
const newWorkflowForm = () => ({
  source: '',
  flowType: 'gva_codegen',
  extraConstraints: ''
})
const emptyAnalysis = () => ({
  summary: '',
  recommendedPackageType: '',
  missingInfo: [],
  suggestions: [],
  modules: [],
  clientPages: [],
  rawText: '',
  rawJson: ''
})
const emptyWorkflow = () => ({
  summary: '',
  steps: [],
  rawText: '',
  rawJson: ''
})
const newSession = (tab) => ({
  id: 0,
  tab,
  title: '',
  summary: '',
  conversationId: '',
  messageId: '',
  currentNodeId: '',
  formData: {},
  resultData: {},
  messages: []
})
const parseJSON = (value) => {
  try {
    return JSON.parse(value)
  } catch (error) {
    return null
  }
}
const clone = (value, fallback = {}) => {
  try {
    return JSON.parse(JSON.stringify(value ?? fallback))
  } catch (error) {
    return fallback
  }
}
const firstText = (...values) =>
  values.find((item) => typeof item === 'string' && item.trim()) || ''
const firstParagraph = (text) =>
  String(text || '')
    .split(/\n{2,}/)
    .map((item) => item.trim())
    .find(Boolean) || ''
const formatPayload = (payload) => {
  try {
    return JSON.stringify(payload, null, 2)
  } catch (error) {
    return String(payload || '')
  }
}
const normalizeText = (value) =>
  String(value || '')
    .replace(/\r\n/g, '\n')
    .trim()
const toArray = (value) => (Array.isArray(value) ? value : [])
const loadStorage = (key, fallback) => {
  try {
    return { ...fallback, ...JSON.parse(localStorage.getItem(key) || '{}') }
  } catch (error) {
    return { ...fallback }
  }
}
const activeTab = ref('analysis')
const followUpInput = ref('')
const historyLoading = ref(false)
const analysisLoading = ref(false)
const workflowLoading = ref(false)
const dumpLoading = reactive({ analysis: false, workflow: false })
const conversationListRef = ref(null)
const settings = reactive(loadStorage(SETTINGS_KEY, defaultSettings))
const activeSessionIds = reactive(
  loadStorage(ACTIVE_SESSION_KEY, { analysis: 0, workflow: 0 })
)
const historyKeywords = reactive({ analysis: '', workflow: '' })
const sessionLists = reactive({ analysis: [], workflow: [] })
const sessions = reactive({
  analysis: newSession('analysis'),
  workflow: newSession('workflow')
})
const analysisForm = reactive(newAnalysisForm())
const workflowForm = reactive(newWorkflowForm())
const analysisResult = ref(emptyAnalysis())
const workflowResult = ref(emptyWorkflow())
const currentSession = computed(() => sessions[activeTab.value])
const currentSessionList = computed(() => sessionLists[activeTab.value])
const currentLoading = computed(() =>
  activeTab.value === 'analysis' ? analysisLoading.value : workflowLoading.value
)
const streamingPreviewText = ref('')
const streamingPreviewRef = ref(null)
const drawerVisible = ref(false)
const drawerTitle = ref('')
const drawerContent = ref('')
const openDrawer = (title, content) => {
  drawerTitle.value = title
  drawerContent.value = content
  drawerVisible.value = true
}
const analysisClientPages = computed(() =>
  Array.isArray(analysisResult.value?.clientPages)
    ? analysisResult.value.clientPages
    : []
)
const hasAnalysisResult = computed(() =>
  Boolean(
    analysisResult.value.summary ||
      analysisResult.value.modules.length ||
      analysisClientPages.value.length ||
      analysisResult.value.missingInfo.length ||
      analysisResult.value.rawText
  )
)
const canSendFollowUp = computed(() =>
  Boolean(currentSession.value.messages.length && followUpInput.value.trim())
)
const selectedAssistantMessage = computed(
  () =>
    currentSession.value.messages.find(
      (item) =>
        item.role === 'assistant' &&
        item.id === currentSession.value.currentNodeId
    ) || null
)
const currentUser = computed(() =>
  String(
    userStore.userInfo.ID ||
      userStore.userInfo.uuid ||
      userStore.userInfo.id ||
      'gva-ai-workflow'
  )
)
const hasWorkflowResult = computed(() =>
  Boolean(
    workflowResult.value.summary ||
      workflowResult.value.steps.length ||
      workflowResult.value.rawText
  )
)
const canDumpCurrentMarkdown = computed(() =>
  activeTab.value === 'analysis' ? hasAnalysisResult.value : hasWorkflowResult.value
)
const currentDumpLoading = computed(() => dumpLoading[activeTab.value])
const isSessionActive = (item) =>
  Number(item.ID || item.id || 0) ===
  Number(activeSessionIds[activeTab.value] || 0)
const formatTime = (value) => {
  const date = new Date(value)
  return Number.isNaN(date.getTime())
    ? String(value || '-')
    : date.toLocaleString()
}

watch(
  () => ({ ...settings }),
  (value) => localStorage.setItem(SETTINGS_KEY, JSON.stringify(value)),
  { deep: true }
)
watch(
  () => ({ ...activeSessionIds }),
  (value) => localStorage.setItem(ACTIVE_SESSION_KEY, JSON.stringify(value)),
  { deep: true }
)
watch(
  () => ({ ...analysisForm }),
  (value) => {
    sessions.analysis.formData = clone(value)
  },
  { deep: true }
)
watch(
  () => ({ ...workflowForm }),
  (value) => {
    sessions.workflow.formData = clone(value)
  },
  { deep: true }
)
watch(
  analysisResult,
  (value) => {
    sessions.analysis.resultData = clone(value)
    sessions.analysis.summary = firstText(value.summary)
  },
  { deep: true }
)
watch(
  workflowResult,
  (value) => {
    sessions.workflow.resultData = clone(value)
    sessions.workflow.summary = firstText(value.summary)
  },
  { deep: true }
)
watch(
  () => [activeTab.value, currentSession.value.messages.length],
  async () => {
    await nextTick()
    if (conversationListRef.value)
      conversationListRef.value.scrollTop =
        conversationListRef.value.scrollHeight
  },
  { flush: 'post' }
)

const extractJson = (text) => {
  if (!text || typeof text !== 'string') return null
  const direct = parseJSON(text.trim())
  if (direct) return direct
  const code = text.match(/```(?:json)?\s*([\s\S]*?)```/i)?.[1]
  if (code) {
    const parsed = parseJSON(code.trim())
    if (parsed) return parsed
  }
  const starts = [text.indexOf('{'), text.indexOf('[')].filter((i) => i >= 0)
  const ends = [text.lastIndexOf('}'), text.lastIndexOf(']')].filter(
    (i) => i >= 0
  )
  return starts.length && ends.length
    ? parseJSON(text.slice(Math.min(...starts), Math.max(...ends) + 1))
    : null
}

const normalizeStrings = (value) => {
  if (!value) return []
  if (Array.isArray(value))
    return value
      .map((item) =>
        typeof item === 'string'
          ? item.trim()
          : firstText(
              item?.label,
              item?.name,
              item?.title,
              JSON.stringify(item)
            )
      )
      .filter(Boolean)
  if (typeof value === 'string')
    return value
      .split(/\n|,|;/)
      .map((item) => item.trim())
      .filter(Boolean)
  return []
}

const normalizeFields = (value) =>
  toArray(value).map((field, index) =>
    typeof field === 'string'
      ? {
          name: field,
          label: field,
          type: 'string',
          required: false,
          description: ''
        }
      : {
          name: firstText(
            field.name,
            field.fieldName,
            field.key,
            `field_${index + 1}`
          ),
          label: firstText(
            field.label,
            field.fieldLabel,
            field.title,
            field.name,
            `Field ${index + 1}`
          ),
          type: firstText(field.type, field.fieldType, 'string'),
          required: Boolean(field.required),
          description: firstText(field.description, field.comment, field.desc),
          dictionary: firstText(
            field.dictionary,
            field.dictionaryName,
            field.dict
          ),
          relation: firstText(
            field.relation,
            field.relationship,
            field.association
          )
        }
  )

const normalizeModules = (value) =>
  toArray(value).map((module, index) =>
    typeof module === 'string'
      ? {
          name: `module_${index + 1}`,
          label: module,
          description: '',
          fields: []
        }
      : {
          name: firstText(
            module.name,
            module.moduleName,
            `module_${index + 1}`
          ),
          label: firstText(
            module.label,
            module.title,
            module.cnName,
            module.name,
            `Module ${index + 1}`
          ),
          description: firstText(
            module.description,
            module.goal,
            module.summary
          ),
          fields: normalizeFields(
            module.fields || module.fieldList || module.columns || []
          )
        }
  )

const normalizeClientPages = (value) =>
  toArray(value).map((page, index) => ({
    name: firstText(page.name, page.pageName, `page_${index + 1}`),
    label: firstText(page.label, page.title, page.name, `Page ${index + 1}`),
    description: firstText(page.description, page.summary, page.goal),
    pageType: firstText(page.pageType, page.type, page.viewType),
    targetModules: normalizeStrings(
      page.targetModules || page.modules || page.relatedModules
    ),
    fields: toArray(page.fields).map((field, fieldIndex) => ({
      name: firstText(field.name, field.fieldName, `page_field_${fieldIndex + 1}`),
      label: firstText(
        field.label,
        field.title,
        field.name,
        `Page field ${fieldIndex + 1}`
      ),
      sourceModule: firstText(field.sourceModule, field.module, field.moduleName),
      sourceField: firstText(field.sourceField, field.field, field.backendField),
      displayType: firstText(field.displayType, field.component, field.type, 'text'),
      required: Boolean(field.required),
      description: firstText(field.description, field.comment, field.desc)
    })),
    interactions: normalizeStrings(
      page.interactions || page.actions || page.behaviors
    ),
    relations: normalizeStrings(
      page.relations || page.mappings || page.fieldRelations
    )
  }))

const normalizeSteps = (value, fallback) => {
  const list = toArray(value)
  if (!list.length) {
    return String(fallback || '')
      .split(/\n(?=(?:Step)\s*\d+|(?:\d+)\.\s)/i)
      .map((item) => item.trim())
      .filter(Boolean)
      .map((item, index) => ({
        title: `Step ${index + 1}`,
        goal: '',
        prompt: item,
        expectedOutput: '',
        suggestedTool: '',
        autoExecutable: false
      }))
  }
  return list.map((step, index) =>
    typeof step === 'string'
      ? {
          title: `Step ${index + 1}`,
          goal: '',
          prompt: step,
          expectedOutput: '',
          suggestedTool: '',
          autoExecutable: false
        }
      : {
          title: firstText(
            step.title,
            step.name,
            step.stepName,
            `Step ${index + 1}`
          ),
          goal: firstText(step.goal, step.description, step.objective),
          prompt: firstText(
            step.prompt,
            step.content,
            step.instruction,
            step.text
          ),
          expectedOutput: firstText(
            step.expectedOutput,
            step.expected,
            step.output
          ),
          suggestedTool: firstText(step.suggestedTool, step.tool, step.action),
          autoExecutable: Boolean(step.autoExecutable)
        }
  )
}

const normalizeAnalysis = (payload, rawText = '') => ({
  summary: firstText(
    payload?.summary,
    payload?.overview,
    payload?.analysisSummary,
    firstParagraph(rawText)
  ),
  recommendedPackageType: firstText(
    payload?.recommendedPackageType,
    payload?.packageType,
    payload?.targetType
  ),
  missingInfo: normalizeStrings(
    payload?.missingInfo ||
      payload?.missing_info ||
      payload?.questions ||
      payload?.clarifyQuestions
  ),
  suggestions: normalizeStrings(
    payload?.suggestions ||
      payload?.recommendations ||
      payload?.advice ||
      payload?.tips
  ),
  clientPages: normalizeClientPages(
    payload?.clientPages || payload?.pages || payload?.client_pages
  ),
  modules: normalizeModules(
    payload?.modules ||
      payload?.moduleList ||
      payload?.entities ||
      payload?.items
  ),
  rawText,
  rawJson: formatPayload(payload || rawText)
})

const normalizeWorkflow = (payload, rawText = '') => ({
  summary: firstText(
    payload?.summary,
    payload?.overview,
    payload?.workflowSummary,
    firstParagraph(rawText)
  ),
  steps: normalizeSteps(
    payload?.steps ||
      payload?.workflow ||
      payload?.prompts ||
      payload?.promptFlow ||
      payload?.promptList,
    rawText
  ),
  rawText,
  rawJson: formatPayload(payload || rawText)
})

const parseDisplay = (content) => {
  const text = normalizeText(content)
  const think = normalizeText(
    text.match(/<think(?:ing)?>\s*([\s\S]*?)<\/think(?:ing)?>/i)?.[1]
  )
  const rest = normalizeText(
    text.replace(/<think(?:ing)?>\s*[\s\S]*?<\/think(?:ing)?>/gi, '')
  )
  const raw = extractJson(rest || text)
    ? formatPayload(extractJson(rest || text))
    : rest || text
  return {
    think,
    rawFile: raw,
    preview: `${(raw || think).slice(0, 140)}${
      (raw || think).length > 140 ? '...' : ''
    }`,
    hasThink: Boolean(think),
    hasRawFile: Boolean(raw)
  }
}

const currentSessionDisplayMessages = computed(() =>
  currentSession.value.messages.map((item) => ({
    ...item,
    display: item.role === 'assistant' ? parseDisplay(item.content) : null,
    isStreaming:
      currentLoading.value &&
      item.role === 'assistant' &&
      item.id === currentSession.value.currentNodeId,
    isSelected:
      item.role === 'assistant' &&
      item.id === currentSession.value.currentNodeId
  }))
)

const STRUCTURED_RESULT_KEYS = [
  'summary',
  'overview',
  'analysisSummary',
  'workflowSummary',
  'recommendedPackageType',
  'packageType',
  'targetType',
  'modules',
  'moduleList',
  'entities',
  'items',
  'clientPages',
  'client_pages',
  'pages',
  'missingInfo',
  'missing_info',
  'questions',
  'clarifyQuestions',
  'suggestions',
  'recommendations',
  'advice',
  'tips',
  'steps',
  'workflow',
  'prompts',
  'promptFlow',
  'promptList'
]

const hasStructuredResultShape = (value) =>
  Boolean(
    value &&
      typeof value === 'object' &&
      STRUCTURED_RESULT_KEYS.some((key) =>
        Object.prototype.hasOwnProperty.call(value, key)
      )
  )

const unwrapStructuredResult = (value, depth = 0) => {
  if (depth > 4 || value == null) return null
  if (typeof value === 'string') {
    const parsed = extractJson(value)
    return parsed?.payload && typeof parsed.payload === 'object'
      ? parsed.payload
      : parsed
  }
  if (Array.isArray(value)) return value
  if (typeof value !== 'object') return null
  if (hasStructuredResultShape(value)) return value

  const candidates = [
    value?.structured,
    value?.payload,
    value?.outputs,
    value?.output,
    value?.result,
    value?.results,
    value?.data?.outputs,
    value?.data?.output,
    value?.data?.result,
    value?.data?.payload,
    value?.data,
    value?.answer,
    value?.text,
    value?.content
  ]

  for (const candidate of candidates) {
    const resolved = unwrapStructuredResult(candidate, depth + 1)
    if (resolved) return resolved
  }

  return null
}

const fallbackAnswerText = (value) => {
  if (typeof value === 'string') return value
  if (value && typeof value === 'object') return formatPayload(value)
  return ''
}

const unwrap = (response) => {
  if (!response) return {}
  if (typeof response.code !== 'undefined') return response.data || {}
  if (
    typeof response.status !== 'undefined' &&
    typeof response.data !== 'undefined'
  )
    return response.data
  return response
}

const resolveChat = (response) => {
  const raw = unwrap(response)
  const structuredCandidate = unwrapStructuredResult(raw)
  const answerText =
    firstText(raw.answer, raw.text, raw.output, raw.content) ||
    fallbackAnswerText(structuredCandidate)
  const parsed = extractJson(answerText)
  const structuredFromText =
    parsed?.payload && typeof parsed.payload === 'object'
      ? parsed.payload
      : parsed
  return {
    answerText,
    structured: structuredFromText || structuredCandidate,
    conversationId: firstText(raw.conversation_id, raw.conversationId),
    messageId: firstText(raw.message_id, raw.messageId)
  }
}

const messageSnapshot = (tab, message) => {
  const payload = Object.keys(message.snapshot || {}).length
    ? message.snapshot
    : extractJson(message.content)
  return tab === 'analysis'
    ? normalizeAnalysis(payload, firstText(payload?.rawText, message.content))
    : normalizeWorkflow(payload, firstText(payload?.rawText, message.content))
}

const applyMessage = (tab, message) => {
  if (!message) {
    if (tab === 'analysis') analysisResult.value = emptyAnalysis()
    else workflowResult.value = emptyWorkflow()
    return
  }
  const snapshot = messageSnapshot(tab, message)
  if (tab === 'analysis') analysisResult.value = snapshot
  else workflowResult.value = snapshot
  sessions[tab].resultData = clone(snapshot)
  sessions[tab].summary = firstText(snapshot.summary)
}

const applySession = (tab) => {
  const session = sessions[tab]
  if (tab === 'analysis')
    Object.assign(analysisForm, {
      ...newAnalysisForm(),
      ...clone(session.formData)
    })
  else
    Object.assign(workflowForm, {
      ...newWorkflowForm(),
      ...clone(session.formData)
    })
  const selected =
    session.messages.find(
      (item) => item.role === 'assistant' && item.id === session.currentNodeId
    ) ||
    [...session.messages].reverse().find((item) => item.role === 'assistant')
  if (selected) {
    session.currentNodeId = selected.id
    applyMessage(tab, selected)
  } else if (tab === 'analysis') {
    analysisResult.value = Object.keys(session.resultData || {}).length
      ? normalizeAnalysis(
          session.resultData,
          firstText(session.resultData?.rawText)
        )
      : emptyAnalysis()
  } else {
    workflowResult.value = Object.keys(session.resultData || {}).length
      ? normalizeWorkflow(
          session.resultData,
          firstText(session.resultData?.rawText)
        )
      : emptyWorkflow()
  }
}

const sessionTitle = (tab) => {
  const firstUser = sessions[tab].messages.find(
    (item) => item.role === 'user'
  )?.content
  return firstText(
    firstUser,
    tab === 'analysis' ? analysisForm.requirement : workflowForm.source,
    tab === 'analysis'
      ? analysisResult.value.summary
      : workflowResult.value.summary
  ).slice(0, 120)
}

const normalizeInlineText = (value) =>
  normalizeText(value).replace(/\s+/g, ' ').trim()

const formatModuleFieldForTransfer = (field) => {
  const meta = [field.type || 'string']
  if (field.required) meta.push('required')
  if (field.dictionary) meta.push(`dict:${field.dictionary}`)
  if (field.relation) meta.push(`rel:${field.relation}`)

  const lines = [`- ${field.label || field.name} (${field.name}) | ${meta.join(' | ')}`]

  if (field.description) {
    lines.push(`  Note: ${normalizeInlineText(field.description)}`)
  }

  return lines.join('\n')
}

const buildAnalysisTransferText = () => {
  const blocks = [
    '# Analysis result',
    'Generate a workflow directly from the analysis below. Do not output items like "unclear requirements" or "to be confirmed".'
  ]

  if (analysisForm.requirement.trim()) {
    blocks.push(`Requirement: ${normalizeInlineText(analysisForm.requirement)}`)
  }
  if (analysisForm.businessScene.trim()) {
    blocks.push(`Business scenario: ${normalizeInlineText(analysisForm.businessScene)}`)
  }
  if (analysisForm.packageType) {
    blocks.push(`Target type: ${analysisForm.packageType}`)
  }
  if (analysisForm.extraConstraints.trim()) {
    blocks.push(`Extra constraints: ${normalizeInlineText(analysisForm.extraConstraints)}`)
  }
  if (analysisForm.hasClientPage) {
    blocks.push('Has client pages: Yes')
  }
  if (analysisForm.clientPageDescription.trim()) {
    blocks.push(
      `Client page description: ${normalizeInlineText(analysisForm.clientPageDescription)}`
    )
  }
  if (analysisForm.clientPageConstraints.trim()) {
    blocks.push(
      `Client page constraints: ${normalizeInlineText(
        analysisForm.clientPageConstraints
      )}`
    )
  }
  if (analysisResult.value.summary) {
    blocks.push(`Summary: ${normalizeInlineText(analysisResult.value.summary)}`)
  }
  if (analysisResult.value.recommendedPackageType) {
    blocks.push(`Recommended type: ${analysisResult.value.recommendedPackageType}`)
  }
  if (analysisResult.value.suggestions.length) {
    blocks.push(`Suggestions: ${analysisResult.value.suggestions.join('; ')}`)
  }

  if (analysisResult.value.modules.length) {
    const moduleText = analysisResult.value.modules
      .map((module, index) => {
        const lines = [
          `${index + 1}. ${module.label || module.name || `Module ${index + 1}`}`
        ]

        if (module.description) {
          lines.push(`   Description: ${normalizeInlineText(module.description)}`)
        }

        if (module.fields.length) {
          lines.push('   Fields:')
          lines.push(
            ...module.fields.map((field) => `   ${formatModuleFieldForTransfer(field)}`)
          )
        }

        return lines.join('\n')
      })
      .join('\n\n')

    blocks.push(`Modules:\n${moduleText}`)
  }

  if (analysisClientPages.value.length) {
    const pageText = analysisClientPages.value
      .map((page, index) => {
        const lines = [
          `${index + 1}. ${page.label || page.name || `Page ${index + 1}`}`
        ]

        if (page.pageType) {
          lines.push(`   Page type: ${page.pageType}`)
        }
        if (page.description) {
          lines.push(`   Description: ${normalizeInlineText(page.description)}`)
        }
        if (page.targetModules.length) {
          lines.push(`   Related modules: ${page.targetModules.join(', ')}`)
        }
        if (page.fields.length) {
          lines.push('   Fields:')
          lines.push(
            ...page.fields.map((field) => {
              const source = [field.sourceModule, field.sourceField]
                .filter(Boolean)
                .join('.')
              const meta = [field.displayType || 'text']
              if (field.required) meta.push('required')
              if (source) meta.push(`map:${source}`)
              const text = `- ${field.label || field.name} (${field.name}) | ${meta.join(' | ')}`
              return `   ${text}${field.description ? `\n     Note: ${normalizeInlineText(field.description)}` : ''}`
            })
          )
        }
        if (page.interactions.length) {
          lines.push(`   Interactions: ${page.interactions.join('; ')}`)
        }
        if (page.relations.length) {
          lines.push(`   Field mappings: ${page.relations.join('; ')}`)
        }

        return lines.join('\n')
      })
      .join('\n\n')

    blocks.push(`Client pages:\n${pageText}`)
  }

  return blocks.join('\n\n')
}

const sessionPayload = (tab) => ({
  id: sessions[tab].id,
  tab,
  title: sessions[tab].title || sessionTitle(tab),
  summary:
    sessions[tab].summary || firstText(sessions[tab].resultData?.summary),
  conversationId: sessions[tab].conversationId,
  messageId: sessions[tab].messageId,
  currentNodeId: sessions[tab].currentNodeId,
  settings: { extraPayload: settings.extraPayload },
  formData: clone(sessions[tab].formData),
  resultData: clone(sessions[tab].resultData),
  messages: sessions[tab].messages.map((item) => ({
    id: item.id,
    role: item.role,
    content: item.content,
    snapshot: clone(item.snapshot),
    conversationId: item.conversationId,
    messageId: item.messageId,
    createdAt: item.createdAt
  }))
})

const dumpCurrentMarkdown = async () => {
  const tab = activeTab.value
  const hasResult =
    tab === 'analysis' ? hasAnalysisResult.value : hasWorkflowResult.value

  if (!hasResult) {
    ElMessage.warning(
      tab === 'analysis'
        ? t('admin.systemtools.aiworkflow.no_analysis_dump')
        : t('admin.systemtools.aiworkflow.no_workflow_dump')
    )
    return
  }

  try {
    await persistSession(tab, false)
  } catch (error) {
    // Dump uses in-memory data first; save failures should not block.
  }

  dumpLoading[tab] = true
  try {
    const data = unwrap(await dumpAIWorkflowMarkdown(sessionPayload(tab)))
    const result = data.result || data
    const path = firstText(result.filePath, result.relativePath)
    await ElMessageBox.alert(
      path || t('admin.systemtools.aiworkflow.markdown_dump_done'),
      tab === 'analysis' ? t('admin.systemtools.aiworkflow.analysis_dumped') : t('admin.systemtools.aiworkflow.prompts_dumped'),
      {
        confirmButtonText: t('admin.systemtools.aiworkflow.ok')
      }
    )
  } finally {
    dumpLoading[tab] = false
  }
}

const startNewConversation = (tab = activeTab.value) => {
  sessions[tab] = newSession(tab)
  activeSessionIds[tab] = 0
  if (tab === 'analysis') {
    Object.assign(analysisForm, newAnalysisForm())
    analysisResult.value = emptyAnalysis()
  } else {
    Object.assign(workflowForm, newWorkflowForm())
    workflowResult.value = emptyWorkflow()
  }
  if (activeTab.value === tab) followUpInput.value = ''
  ElMessage.success(
    tab === 'analysis' ? t('admin.systemtools.aiworkflow.new_analysis_session') : t('admin.systemtools.aiworkflow.new_workflow_session')
  )
}

const hydrateSession = (tab, raw) => {
  sessions[tab] = {
    ...newSession(tab),
    id: Number(raw?.ID || raw?.id || 0),
    tab: firstText(raw?.tab) || tab,
    title: firstText(raw?.title),
    summary: firstText(raw?.summary),
    conversationId: firstText(raw?.conversationId, raw?.conversation_id),
    messageId: firstText(raw?.messageId, raw?.message_id),
    currentNodeId: firstText(raw?.currentNodeId),
    formData: clone(raw?.formData),
    resultData: clone(raw?.resultData),
    messages: toArray(raw?.messages).map((item) => ({
      id:
        firstText(item?.id) ||
        `msg-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
      role: firstText(item?.role) || 'assistant',
      content: firstText(item?.content),
      snapshot: clone(item?.snapshot),
      conversationId: firstText(item?.conversationId, item?.conversation_id),
      messageId: firstText(item?.messageId, item?.message_id),
      createdAt: firstText(item?.createdAt)
    }))
  }
  activeSessionIds[tab] = sessions[tab].id
  if (activeTab.value === tab) applySession(tab)
}

const loadSessionList = async (tab = activeTab.value) => {
  historyLoading.value = true
  try {
    const data = unwrap(
      await getAIWorkflowSessionList({
        page: 1,
        pageSize: 50,
        tab,
        keyword: historyKeywords[tab]
      })
    )
    sessionLists[tab] = toArray(data.list)
  } finally {
    historyLoading.value = false
  }
}

const openSession = async (item) => {
  const data = unwrap(
    await getAIWorkflowSessionDetail({ id: Number(item.ID || item.id || 0) })
  )
  if (data.session) hydrateSession(activeTab.value, data.session)
  followUpInput.value = ''
}

const persistSession = async (tab, refreshList = true) => {
  const hasResult =
    tab === 'analysis' ? hasAnalysisResult.value : hasWorkflowResult.value
  if (!sessions[tab].messages.length && !hasResult) return null
  const data = unwrap(await saveAIWorkflowSession(sessionPayload(tab)))
  if (data.session) hydrateSession(tab, data.session)
  if (refreshList) await loadSessionList(tab)
  return data.session || null
}

const removeSession = async (item) => {
  await ElMessageBox.confirm(
    t('admin.systemtools.aiworkflow.delete_session_confirm'),
    t('admin.systemtools.aiworkflow.delete_session_title'),
    {
      type: 'warning',
      confirmButtonText: t('admin.common.delete'),
      cancelButtonText: t('admin.common.cancel')
    }
  )
  await deleteAIWorkflowSession({ id: Number(item.ID || item.id || 0) })
  if (isSessionActive(item)) startNewConversation(activeTab.value)
  await loadSessionList(activeTab.value)
  ElMessage.success(t('admin.systemtools.aiworkflow.session_deleted'))
}

const selectMessageNode = async (messageId) => {
  const message = currentSession.value.messages.find(
    (item) => item.role === 'assistant' && item.id === messageId
  )
  if (!message) return
  currentSession.value.currentNodeId = messageId
  applyMessage(activeTab.value, message)
  if (currentSession.value.id) await persistSession(activeTab.value, false)
}

const rollbackToMessage = async (messageId) => {
  const idx = currentSession.value.messages.findIndex(
    (item) => item.role === 'assistant' && item.id === messageId
  )
  if (idx < 0) return
  await ElMessageBox.confirm(
    t('admin.systemtools.aiworkflow.rollback_session_confirm'),
    t('admin.systemtools.aiworkflow.rollback_session_title'),
    {
      type: 'warning',
      confirmButtonText: t('admin.systemtools.aiworkflow.rollback_action'),
      cancelButtonText: t('admin.common.cancel')
    }
  )
  const message = currentSession.value.messages[idx]
  currentSession.value.messages = currentSession.value.messages.slice(
    0,
    idx + 1
  )
  currentSession.value.currentNodeId = messageId
  currentSession.value.conversationId = ''
  currentSession.value.messageId = ''
  applyMessage(activeTab.value, message)
  await persistSession(activeTab.value)
  ElMessage.success(t('admin.systemtools.aiworkflow.rolled_back_node'))
}

const parseExtraPayload = () => {
  if (!settings.extraPayload.trim()) return {}
  try {
    const parsed = JSON.parse(settings.extraPayload)
    return parsed && typeof parsed === 'object' ? parsed : {}
  } catch (error) {
    ElMessage.error(t('admin.systemtools.aiworkflow.invalid_extra_payload'))
    throw error
  }
}

const sendChat = async ({ tab, query, inputs, onProgress }) => {
  let extra = {}
  try {
    extra = parseExtraPayload()
  } catch (error) {
    return null
  }
  const safeExtra = { ...extra }
  delete safeExtra.mode
  delete safeExtra.query
  delete safeExtra.inputs
  delete safeExtra.user
  delete safeExtra.response_mode
  delete safeExtra.conversation_id
  const requestData = {
    ...safeExtra,
    mode: TAB_MODE_MAP[tab],
    query:
      query ||
      (tab === 'analysis'
        ? analysisForm.requirement
        : `Generate ${
            FLOW_TYPE_LABEL_MAP[workflowForm.flowType] || 'Prompt workflow'
          } based on the current input.`),
    inputs,
    user: currentUser.value,
    response_mode: 'streaming',
    scene: 'gva_ai_workflow'
  }
  if (!String(requestData.query || '').trim()) {
    ElMessage.error(t('admin.systemtools.aiworkflow.missing_query'))
    return null
  }
  if (sessions[tab].conversationId)
    requestData.conversation_id = sessions[tab].conversationId
  try {
    return resolveChat(
      await (tab === 'analysis'
        ? analyzeRequirementByAISSEStream(requestData, {
            onMessage: onProgress
          })
        : generatePromptFlowByAISSEStream(requestData, {
            onMessage: onProgress
          }))
    )
  } catch (error) {
    ElMessage.error(error?.message || t('admin.systemtools.aiworkflow.ai_request_failed'))
    return null
  }
}

const addMessage = (tab, role, content, extras = {}) => {
  const message = {
    id: `${role}-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
    role,
    content,
    snapshot: clone(extras.snapshot),
    conversationId: firstText(extras.conversationId),
    messageId: firstText(extras.messageId),
    createdAt: new Date().toISOString()
  }
  sessions[tab].messages.push(message)
  return message
}

const removeMessages = (tab, messageIds = []) => {
  const idSet = new Set(messageIds.filter(Boolean))
  sessions[tab].messages = sessions[tab].messages.filter(
    (item) => !idSet.has(item.id)
  )
}

const scrollConversationToBottom = async () => {
  await nextTick()
  if (conversationListRef.value) {
    conversationListRef.value.scrollTop = conversationListRef.value.scrollHeight
  }
}

const updateStreamMessage = (tab, message, payload = {}) => {
  if (!message) return
  if (payload.text) {
    message.content = payload.text
    streamingPreviewText.value = payload.text
    nextTick(() => {
      if (streamingPreviewRef.value) {
        streamingPreviewRef.value.scrollTop = streamingPreviewRef.value.scrollHeight
      }
    })
  }
  if (payload.conversationId) {
    message.conversationId = payload.conversationId
    sessions[tab].conversationId = payload.conversationId
  }
  if (payload.messageId) {
    message.messageId = payload.messageId
    sessions[tab].messageId = payload.messageId
  }
  scrollConversationToBottom()
}

const runAnalysis = async () => {
  if (!analysisForm.requirement.trim()) return ElMessage.warning(t('admin.systemtools.aiworkflow.enter_requirement_first'))
  analysisLoading.value = true
  streamingPreviewText.value = ''
  const streamUserMessage = addMessage(
    'analysis',
    'user',
    analysisForm.requirement
  )
  const streamAssistant = addMessage('analysis', 'assistant', '')
  sessions.analysis.currentNodeId = streamAssistant.id
  scrollConversationToBottom()
  try {
    const result = await sendChat({
      tab: 'analysis',
      query: analysisForm.requirement,
      inputs: {
        requirement: analysisForm.requirement,
        packageType: analysisForm.packageType,
        businessScene: analysisForm.businessScene,
        extraConstraints: analysisForm.extraConstraints,
        hasClientPage: analysisForm.hasClientPage,
        clientPageDescription: analysisForm.clientPageDescription,
        clientPageConstraints: analysisForm.clientPageConstraints
      },
      onProgress: (payload) =>{
        updateStreamMessage('analysis', streamAssistant, payload)
      }
    })
    if (!result) {
      if (!streamAssistant.content) {
        removeMessages('analysis', [streamUserMessage.id, streamAssistant.id])
      }
      return
    }
    const snapshot = normalizeAnalysis(result.structured, result.answerText)
    sessions.analysis.conversationId =
      result.conversationId || sessions.analysis.conversationId
    sessions.analysis.messageId =
      result.messageId || sessions.analysis.messageId
    streamAssistant.content =
      result.answerText || streamAssistant.content || t('admin.systemtools.aiworkflow.model_no_text')
    streamAssistant.snapshot = snapshot
    streamAssistant.conversationId =
      result.conversationId || streamAssistant.conversationId
    streamAssistant.messageId = result.messageId || streamAssistant.messageId
    sessions.analysis.currentNodeId = streamAssistant.id
    sessions.analysis.title = sessionTitle('analysis')
    analysisResult.value = snapshot
    await persistSession('analysis')
    ElMessage.success(t('admin.systemtools.aiworkflow.analysis_completed'))
  } finally {
    analysisLoading.value = false
    streamingPreviewText.value = ''
  }
}

const pushAnalysisToWorkflow = async () => {
  if (!hasAnalysisResult.value) return ElMessage.warning(t('admin.systemtools.aiworkflow.no_analysis_yet'))
  await switchTab('workflow')
  workflowForm.source = buildAnalysisTransferText()
  ElMessage.success(t('admin.systemtools.aiworkflow.sent_to_workflow'))
}

const runWorkflow = async () => {
  if (!workflowForm.source.trim())
    return ElMessage.warning(t('admin.systemtools.aiworkflow.enter_requirements_or_analysis'))
  workflowLoading.value = true
  streamingPreviewText.value = ''
  const streamUserMessage = addMessage('workflow', 'user', '')
  const streamAssistant = addMessage('workflow', 'assistant', '')
  sessions.workflow.currentNodeId = streamAssistant.id
  scrollConversationToBottom()
  try {
    const query = `Generate ${
      FLOW_TYPE_LABEL_MAP[workflowForm.flowType] || 'Prompt workflow'
    } based on the current input.`
    streamUserMessage.content = query
    const result = await sendChat({
      tab: 'workflow',
      query,
      inputs: {
        source: workflowForm.source,
        flowType: workflowForm.flowType,
        extraConstraints: workflowForm.extraConstraints
      },
      onProgress: (payload) =>{
        updateStreamMessage('workflow', streamAssistant, payload)
      }
    })
    if (!result) {
      if (!streamAssistant.content) {
        removeMessages('workflow', [streamUserMessage.id, streamAssistant.id])
      }
      return
    }
    const snapshot = normalizeWorkflow(result.structured, result.answerText)
    sessions.workflow.conversationId =
      result.conversationId || sessions.workflow.conversationId
    sessions.workflow.messageId =
      result.messageId || sessions.workflow.messageId
    streamAssistant.content =
      result.answerText || streamAssistant.content || t('admin.systemtools.aiworkflow.model_no_text')
    streamAssistant.snapshot = snapshot
    streamAssistant.conversationId =
      result.conversationId || streamAssistant.conversationId
    streamAssistant.messageId = result.messageId || streamAssistant.messageId
    sessions.workflow.currentNodeId = streamAssistant.id
    sessions.workflow.title = sessionTitle('workflow')
    workflowResult.value = snapshot
    await persistSession('workflow')
    ElMessage.success(t('admin.systemtools.aiworkflow.workflow_generated'))
  } finally {
    workflowLoading.value = false
    streamingPreviewText.value = ''
  }
}

const sendFollowUp = async () => {
  const query = followUpInput.value.trim()
  if (!query) return ElMessage.warning(t('admin.systemtools.aiworkflow.enter_followup'))
  if (!currentSession.value.messages.length)
    return ElMessage.warning(t('admin.systemtools.aiworkflow.no_session_content'))
  const tab = activeTab.value
  if (tab === 'analysis') analysisLoading.value = true
  else workflowLoading.value = true
  streamingPreviewText.value = ''
  const streamUserMessage = addMessage(tab, 'user', query)
  const streamAssistant = addMessage(tab, 'assistant', '')
  sessions[tab].currentNodeId = streamAssistant.id
  scrollConversationToBottom()
  try {
    const result = await sendChat({
      tab,
      query,
      inputs:
        tab === 'analysis'
          ? {
              requirement: analysisForm.requirement,
              packageType: analysisForm.packageType,
              businessScene: analysisForm.businessScene,
              extraConstraints: analysisForm.extraConstraints,
              hasClientPage: analysisForm.hasClientPage,
              clientPageDescription: analysisForm.clientPageDescription,
              clientPageConstraints: analysisForm.clientPageConstraints
            }
          : {
              source: workflowForm.source,
              flowType: workflowForm.flowType,
              extraConstraints: workflowForm.extraConstraints
            },
      onProgress: (payload) => updateStreamMessage(tab, streamAssistant, payload)
    })
    if (!result) {
      if (!streamAssistant.content) {
        removeMessages(tab, [streamUserMessage.id, streamAssistant.id])
      }
      return
    }
    const snapshot =
      tab === 'analysis'
        ? normalizeAnalysis(result.structured, result.answerText)
        : normalizeWorkflow(result.structured, result.answerText)
    sessions[tab].conversationId =
      result.conversationId || sessions[tab].conversationId
    sessions[tab].messageId = result.messageId || sessions[tab].messageId
    streamAssistant.content =
      result.answerText || streamAssistant.content || t('admin.systemtools.aiworkflow.model_no_text')
    streamAssistant.snapshot = snapshot
    streamAssistant.conversationId =
      result.conversationId || streamAssistant.conversationId
    streamAssistant.messageId = result.messageId || streamAssistant.messageId
    sessions[tab].currentNodeId = streamAssistant.id
    sessions[tab].title = sessionTitle(tab)
    if (tab === 'analysis') analysisResult.value = snapshot
    else workflowResult.value = snapshot
    followUpInput.value = ''
    await persistSession(tab)
    ElMessage.success(t('admin.systemtools.aiworkflow.followup_sent'))
  } finally {
    if (tab === 'analysis') analysisLoading.value = false
    else workflowLoading.value = false
    streamingPreviewText.value = ''
  }
}

const copyText = async (text, successMessage = t('admin.systemtools.aiworkflow.copied')) => {
  if (!text) return ElMessage.warning(t('admin.systemtools.aiworkflow.nothing_to_copy'))
  try {
    await navigator.clipboard.writeText(text)
  } catch (error) {
    const textarea = document.createElement('textarea')
    textarea.value = text
    textarea.style.position = 'fixed'
    textarea.style.left = '-9999px'
    document.body.appendChild(textarea)
    textarea.focus()
    textarea.select()
    document.execCommand('copy')
    document.body.removeChild(textarea)
  }
  ElMessage.success(successMessage)
}

const copyAnalysisResult = async () =>
  copyText(
    JSON.stringify(clone(analysisResult.value), null, 2),
    'Copied structured analysis result'
  )
const copyAllWorkflowPrompts = async () =>
  copyText(
    workflowResult.value.steps
      .map(
        (step, index) =>
          `# ${index + 1}. ${step.title || `Step ${index + 1}`}\n${
            step.prompt || ''
          }`
      )
      .join('\n\n'),
    'Copied all prompts'
  )
const copyWorkflowRaw = async () =>
  copyText(
    workflowResult.value.rawText || workflowResult.value.rawJson,
    'Copied raw result'
  )
const switchTab = async (tab) => {
  activeTab.value = tab
  followUpInput.value = ''
  applySession(tab)
  if (!sessionLists[tab].length) await loadSessionList(tab)
}
const fillExample = () => {
  switchTab('analysis')
  analysisForm.requirement =
    'Build a leave management system: employees submit leave requests, department managers approve, HR reviews. Support status filtering and export reports. Also include leave-type and approval-status dictionaries.'
  analysisForm.packageType = 'auto'
  analysisForm.businessScene = 'OA internal admin'
  analysisForm.extraConstraints =
    'Need attachment fields; may add return-from-leave later; prefer a structured plan that fits gin-vue-admin codegen.'
  analysisForm.hasClientPage = true
  analysisForm.clientPageDescription =
    'Need an employee H5 page to submit requests, track approval progress, upload attachments, and initiate return-from-leave.'
  analysisForm.clientPageConstraints =
    'Mobile-first; keep fields minimal; show approval status clearly; attachments support images and PDF.'
  workflowForm.source = ''
  workflowForm.flowType = 'gva_codegen'
  workflowForm.extraConstraints = 'Provide copyable prompts for each step and explain expected outputs.'
}
const goAutoCode = () => router.push({ name: 'autoCode' })

onMounted(async () => {
  await Promise.all([loadSessionList('analysis'), loadSessionList('workflow')])
  for (const tab of ['analysis', 'workflow']) {
    const preferred = Number(activeSessionIds[tab] || 0)
    const target = preferred ? { ID: preferred } : sessionLists[tab][0]
    if (target?.ID || target?.id) {
      const data = unwrap(
        await getAIWorkflowSessionDetail({ id: Number(target.ID || target.id) })
      )
      if (data.session) hydrateSession(tab, data.session)
    }
  }
  applySession(activeTab.value)
})
</script>

<style scoped>
.ai-workflow-page :deep(.el-collapse-item__header) {
  padding-left: 6px;
}
</style>
