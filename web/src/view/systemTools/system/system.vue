<template>
  <div class="system">
    <el-form ref="form" :model="config" label-width="240px">
      <!--  System start  -->
      <el-tabs v-model="activeNames">
        <el-tab-pane :label="t('admin.systemtools.system.system')" name="1" class="mt-3.5">
          <el-form-item :label="t('admin.systemtools.system.port')">
            <el-input-number
              v-model="config.system.addr"
              :placeholder="t('admin.systemtools.system.enter_port')"
            />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.db_type')">
            <el-select v-model="config.system['db-type']" class="w-full">
              <el-option value="mysql" />
              <el-option value="pgsql" />
              <el-option value="mssql" />
              <el-option value="sqlite" />
              <el-option value="oracle" />
            </el-select>
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.oss_type')">
            <el-select v-model="config.system['oss-type']" class="w-full">
              <el-option value="local" :label="t('admin.systemtools.system.oss_local')" />
              <el-option value="qiniu" :label="t('admin.systemtools.system.oss_qiniu')" />
              <el-option value="tencent-cos" :label="t('admin.systemtools.system.oss_tencent')" />
              <el-option value="aliyun-oss" :label="t('admin.systemtools.system.oss_aliyun')" />
              <el-option value="huawei-obs" :label="t('admin.systemtools.system.oss_huawei')" />
              <el-option value="cloudflare-r2" :label="t('admin.systemtools.system.oss_cloudflare')" />
              <el-option value="minio">MinIO</el-option>
            </el-select>
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.multipoint_login')">
            <el-switch v-model="config.system['use-multipoint']" />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.enable_redis')">
            <el-switch v-model="config.system['use-redis']" />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.enable_mongo')">
            <el-switch v-model="config.system['use-mongo']" />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.strict_role_mode')">
            <el-switch v-model="config.system['use-strict-auth']" />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.ip_limit_count')">
            <el-input-number v-model.number="config.system['iplimit-count']" />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.ip_limit_duration')">
            <el-input-number v-model.number="config.system['iplimit-time']" />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.disable_auto_migrate')">
            <el-switch v-model="config.system['disable-auto-migrate']" />
          </el-form-item>
          <el-tooltip
            :content="t('admin.systemtools.system.router_prefix_tooltip')"
            placement="top-start"
          >
            <el-form-item :label="t('admin.systemtools.system.global_router_prefix')">
              <el-input
                v-model.trim="config.system['router-prefix']"
                :placeholder="t('admin.systemtools.system.enter_router_prefix')"
              />
            </el-form-item>
          </el-tooltip>
        </el-tab-pane>
        <el-tab-pane :label="t('admin.systemtools.system.jwt_signing_key')" name="2" class="mt-3.5">
          <el-form-item :label="t('admin.systemtools.system.jwt_signing_key')">
            <el-input
              v-model.trim="config.jwt['signing-key']"
              :placeholder="t('admin.systemtools.system.enter_signing_key')"
            >
              <template #append>
                <el-button @click="getUUID">{{ t('admin.systemtools.system.generate') }}</el-button>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.expires_in')">
            <el-input
              v-model.trim="config.jwt['expires-time']"
              :placeholder="t('admin.systemtools.system.enter_expires_in')"
            />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.buffer_time')">
            <el-input
              v-model.trim="config.jwt['buffer-time']"
              :placeholder="t('admin.systemtools.system.enter_buffer_time')"
            />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.issuer')">
            <el-input
              v-model.trim="config.jwt.issuer"
              :placeholder="t('admin.systemtools.system.enter_issuer')"
            />
          </el-form-item>
        </el-tab-pane>
        <el-tab-pane :label="t('admin.systemtools.system.zap_logger')" name="3" class="mt-3.5">
          <el-form-item :label="t('admin.systemtools.system.level')">
            <el-select v-model="config.zap.level">
              <el-option value="off" :label="t('admin.systemtools.system.level_off')" />
              <el-option value="fatal" :label="t('admin.systemtools.system.level_fatal')" />
              <el-option value="error" :label="t('admin.systemtools.system.level_error')" />
              <el-option value="warn" :label="t('admin.systemtools.system.level_warn')" />
              <el-option value="info" :label="t('admin.systemtools.system.level_info')" />
              <el-option value="debug" :label="t('admin.systemtools.system.level_debug')" />
              <el-option value="trace" :label="t('admin.systemtools.system.level_trace')" />
            </el-select>
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.format')">
            <el-select v-model="config.zap.format">
              <el-option value="console" label="console" />
              <el-option value="json" label="json" />
            </el-select>
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.prefix')">
            <el-input
              v-model.trim="config.zap.prefix"
              :placeholder="t('admin.systemtools.system.enter_prefix')"
            />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.directory')">
            <el-input
              v-model.trim="config.zap.director"
              :placeholder="t('admin.systemtools.system.enter_directory')"
            />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.encode_level')">
            <el-select v-model="config.zap['encode-level']" class="w-6/12">
              <el-option
                value="LowercaseLevelEncoder"
                label="LowercaseLevelEncoder"
              />
              <el-option
                value="LowercaseColorLevelEncoder"
                label="LowercaseColorLevelEncoder"
              />
              <el-option
                value="CapitalLevelEncoder"
                label="CapitalLevelEncoder"
              />
              <el-option
                value="CapitalColorLevelEncoder"
                label="CapitalColorLevelEncoder"
              />
            </el-select>
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.stacktrace_key')">
            <el-input
              v-model.trim="config.zap['stacktrace-key']"
              :placeholder="t('admin.systemtools.system.enter_stacktrace_key')"
            />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.retention_days')">
            <el-input-number v-model="config.zap['retention-day']" />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.show_line')">
            <el-switch v-model="config.zap['show-line']" />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.log_to_console')">
            <el-switch v-model="config.zap['log-in-console']" />
          </el-form-item>
        </el-tab-pane>
        <el-tab-pane
          :label="t('admin.systemtools.system.redis')"
          name="4"
          class="mt-3.5"
          v-if="config.system['use-redis']"
        >
          <el-form-item :label="t('admin.systemtools.system.db')">
            <el-input-number v-model="config.redis.db" min="0" max="16" />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.address')">
            <el-input
              v-model.trim="config.redis.addr"
              :placeholder="t('admin.systemtools.system.enter_address')"
            />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.password')">
            <el-input
              v-model.trim="config.redis.password"
              :placeholder="t('admin.systemtools.system.enter_password')"
            />
          </el-form-item>
        </el-tab-pane>
        <el-tab-pane :label="t('admin.systemtools.system.email')" name="5" class="mt-3.5">
          <el-form-item :label="t('admin.systemtools.system.to')">
            <el-input
              v-model="config.email.to"
              :placeholder="t('admin.systemtools.system.emails_comma')"
            />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.port')">
            <el-input-number v-model="config.email.port" />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.from')">
            <el-input
              v-model.trim="config.email.from"
              :placeholder="t('admin.systemtools.system.enter_sender_email')"
            />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.host')">
            <el-input
              v-model.trim="config.email.host"
              :placeholder="t('admin.systemtools.system.enter_host')"
            />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.use_ssl')">
            <el-switch v-model="config.email['is-ssl']" />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.use_login_auth')">
            <el-switch v-model="config.email['is-loginauth']" />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.secret')">
            <el-input
              v-model.trim="config.email.secret"
              :placeholder="t('admin.systemtools.system.enter_secret')"
            />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.test_email')">
            <el-button @click="email">{{ t('admin.systemtools.system.send_test_email') }}</el-button>
          </el-form-item>
        </el-tab-pane>
        <el-tab-pane
          :label="t('admin.systemtools.system.mongodb')"
          name="14"
          class="mt-3.5"
          v-if="config.system['use-mongo']"
        >
          <el-form-item :label="t('admin.systemtools.system.coll_name')">
            <el-input
              v-model.trim="config.mongo.coll"
              :placeholder="t('admin.systemtools.system.enter_coll_name')"
            />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.mongo_options')">
            <el-input
              v-model.trim="config.mongo.options"
              :placeholder="t('admin.systemtools.system.enter_mongo_options')"
            />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.database_name')">
            <el-input
              v-model.trim="config.mongo.database"
              :placeholder="t('admin.systemtools.system.enter_database_name')"
            />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.username')">
            <el-input
              v-model.trim="config.mongo.username"
              :placeholder="t('admin.systemtools.system.enter_username')"
            />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.password')">
            <el-input
              v-model.trim="config.mongo.password"
              :placeholder="t('admin.systemtools.system.enter_password')"
            />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.min_pool_size')">
            <el-input-number v-model="config.mongo['min-pool-size']" min="0" />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.max_pool_size')">
            <el-input-number
              v-model="config.mongo['max-pool-size']"
              min="100"
            />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.socket_timeout')">
            <el-input-number
              v-model="config.mongo['socket-timeout-ms']"
              min="0"
            />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.connect_timeout')">
            <el-input-number
              v-model="config.mongo['socket-timeout-ms']"
              min="0"
            />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.enable_zap_logs')">
            <el-switch v-model="config.mongo['is-zap']" />
          </el-form-item>
          <el-form-item
            v-for="(item, k) in config.mongo.hosts"
            :key="k"
            :label="t('admin.systemtools.system.node_label', { n: k + 1 })"
          >
            <div v-for="(_, k2) in item" :key="k2">
              <el-form-item :key="k + k2" :label="k2" label-width="60">
                <el-input
                  v-model.trim="item[k2]"
                :placeholder="k2 === 'host' ? t('admin.systemtools.system.enter_host') : t('admin.systemtools.system.enter_port')"
                />
              </el-form-item>
            </div>
            <el-form-item v-if="k > 0">
              <el-button
                type="danger"
                size="small"
                plain
                :icon="Minus"
                @click="removeNode(k)"
                class="ml-3"
              />
            </el-form-item>
          </el-form-item>
          <el-form-item>
            <el-button
              type="primary"
              size="small"
              plain
              :icon="Plus"
              @click="addNode"
            />
          </el-form-item>
        </el-tab-pane>
        <el-tab-pane :label="t('admin.systemtools.system.captcha')" name="7" class="mt-3.5">
          <el-form-item :label="t('admin.systemtools.system.key_length')">
            <el-input-number
              v-model="config.captcha['key-long']"
              :min="4"
              :max="6"
            />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.image_width')">
            <el-input-number v-model.number="config.captcha['img-width']" />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.image_height')">
            <el-input-number v-model.number="config.captcha['img-height']" />
          </el-form-item>
        </el-tab-pane>
        <el-tab-pane :label="t('admin.systemtools.system.db_config')" name="9" class="mt-3.5">
          <template v-if="config.system['db-type'] === 'mysql'">
            <el-form-item label="">
              <h3>MySQL</h3>
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.username')">
              <el-input
                v-model.trim="config.mysql.username"
                :placeholder="t('admin.systemtools.system.enter_username')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.password')">
              <el-input
                v-model.trim="config.mysql.password"
                :placeholder="t('admin.systemtools.system.enter_password')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.address')">
              <el-input
                v-model.trim="config.mysql.path"
                :placeholder="t('admin.systemtools.system.enter_address')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.database_name')">
              <el-input
                v-model.trim="config.mysql['db-name']"
                :placeholder="t('admin.systemtools.system.enter_database_name')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.prefix_label')">
              <el-input
                v-model.trim="config.mysql['prefix']"
                :placeholder="t('admin.systemtools.system.optional')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.singular_table')">
              <el-switch v-model="config.mysql['singular']" />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.engine')">
              <el-input
                v-model.trim="config.mysql['engine']"
                :placeholder="t('admin.systemtools.system.default_innodb')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.max_idle_conns')">
              <el-input-number
                v-model="config.mysql['max-idle-conns']"
                :min="1"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.max_open_conns')">
              <el-input-number
                v-model="config.mysql['max-open-conns']"
                :min="1"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.log_writes')">
              <el-switch v-model="config.mysql['log-zap']" />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.log_mode')">
              <el-select v-model="config.mysql['log-mode']">
                <el-option value="off" :label="t('admin.systemtools.system.level_off')" />
                <el-option value="fatal" :label="t('admin.systemtools.system.level_fatal')" />
                <el-option value="error" :label="t('admin.systemtools.system.level_error')" />
                <el-option value="warn" :label="t('admin.systemtools.system.level_warn')" />
                <el-option value="info" :label="t('admin.systemtools.system.level_info')" />
                <el-option value="debug" :label="t('admin.systemtools.system.level_debug')" />
                <el-option value="trace" :label="t('admin.systemtools.system.level_trace')" />
              </el-select>
            </el-form-item>
          </template>
          <template v-if="config.system['db-type'] === 'pgsql'">
            <el-form-item label="">
              <h3>PostgreSQL</h3>
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.username')">
              <el-input
                v-model="config.pgsql.username"
                :placeholder="t('admin.systemtools.system.enter_username')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.password')">
              <el-input
                v-model="config.pgsql.password"
                :placeholder="t('admin.systemtools.system.enter_password')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.address')">
              <el-input
                v-model.trim="config.pgsql.path"
                :placeholder="t('admin.systemtools.system.enter_address')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.database')">
              <el-input
                v-model.trim="config.pgsql['db-name']"
                :placeholder="t('admin.systemtools.system.enter_database')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.prefix_label')">
              <el-input
                v-model.trim="config.pgsql['prefix']"
                :placeholder="t('admin.systemtools.system.enter_prefix_label')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.singular_table')">
              <el-switch v-model="config.pgsql['singular']" />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.engine')">
              <el-input
                v-model.trim="config.pgsql['engine']"
                :placeholder="t('admin.systemtools.system.enter_engine')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.max_idle_conns')">
              <el-input-number v-model="config.pgsql['max-idle-conns']" />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.max_open_conns')">
              <el-input-number v-model="config.pgsql['max-open-conns']" />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.log_writes')">
              <el-switch v-model="config.pgsql['log-zap']" />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.log_mode')">
              <el-select v-model="config.pgsql['log-mode']">
                <el-option value="off" :label="t('admin.systemtools.system.level_off')" />
                <el-option value="fatal" :label="t('admin.systemtools.system.level_fatal')" />
                <el-option value="error" :label="t('admin.systemtools.system.level_error')" />
                <el-option value="warn" :label="t('admin.systemtools.system.level_warn')" />
                <el-option value="info" :label="t('admin.systemtools.system.level_info')" />
                <el-option value="debug" :label="t('admin.systemtools.system.level_debug')" />
                <el-option value="trace" :label="t('admin.systemtools.system.level_trace')" />
              </el-select>
            </el-form-item>
          </template>
          <template v-if="config.system['db-type'] === 'mssql'">
            <el-form-item label="">
              <h3>MsSQL</h3>
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.username')">
              <el-input
                v-model.trim="config.mssql.username"
                :placeholder="t('admin.systemtools.system.enter_username')"
              />
            </el-form-item>
            <el-form-item label.trim="Password">
              <el-input
                v-model.trim="config.mssql.password"
                :placeholder="t('admin.systemtools.system.enter_password')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.address')">
              <el-input
                v-model.trim="config.mssql.path"
                :placeholder="t('admin.systemtools.system.enter_address')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.port')">
              <el-input
                v-model.trim="config.mssql.port"
                :placeholder="t('admin.systemtools.system.enter_port')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.database')">
              <el-input
                v-model.trim="config.mssql['db-name']"
                :placeholder="t('admin.systemtools.system.enter_database')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.prefix_label')">
              <el-input
                v-model.trim="config.mssql['prefix']"
                :placeholder="t('admin.systemtools.system.enter_prefix_label')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.singular_table')">
              <el-switch v-model="config.mssql['singular']" />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.engine')">
              <el-input
                v-model.trim="config.mssql['engine']"
                :placeholder="t('admin.systemtools.system.enter_engine')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.max_idle_conns')">
              <el-input-number v-model="config.mssql['max-idle-conns']" />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.max_open_conns')">
              <el-input-number v-model="config.mssql['max-open-conns']" />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.log_writes')">
              <el-switch v-model="config.mssql['log-zap']" />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.log_mode')">
              <el-select v-model="config.mssql['log-mode']">
                <el-option value="off" :label="t('admin.systemtools.system.level_off')" />
                <el-option value="fatal" :label="t('admin.systemtools.system.level_fatal')" />
                <el-option value="error" :label="t('admin.systemtools.system.level_error')" />
                <el-option value="warn" :label="t('admin.systemtools.system.level_warn')" />
                <el-option value="info" :label="t('admin.systemtools.system.level_info')" />
                <el-option value="debug" :label="t('admin.systemtools.system.level_debug')" />
                <el-option value="trace" :label="t('admin.systemtools.system.level_trace')" />
              </el-select>
            </el-form-item>
          </template>
          <template v-if="config.system['db-type'] === 'sqlite'">
            <el-form-item label="">
              <h3>sqlite</h3>
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.username')">
              <el-input
                v-model.trim="config.sqlite.username"
                :placeholder="t('admin.systemtools.system.enter_username')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.password')">
              <el-input
                v-model.trim="config.sqlite.password"
                :placeholder="t('admin.systemtools.system.enter_password')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.address')">
              <el-input
                v-model.trim="config.sqlite.path"
                :placeholder="t('admin.systemtools.system.enter_address')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.port')">
              <el-input
                v-model.trim="config.sqlite.port"
                :placeholder="t('admin.systemtools.system.enter_port')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.database')">
              <el-input
                v-model.trim="config.sqlite['db-name']"
                :placeholder="t('admin.systemtools.system.enter_database')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.max_idle_conns')">
              <el-input-number v-model="config.sqlite['max-idle-conns']" />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.max_open_conns')">
              <el-input-number v-model="config.sqlite['max-open-conns']" />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.log_writes')">
              <el-switch v-model="config.sqlite['log-zap']" />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.log_mode')">
              <el-select v-model="config.sqlite['log-mode']">
                <el-option value="off" :label="t('admin.systemtools.system.level_off')" />
                <el-option value="fatal" :label="t('admin.systemtools.system.level_fatal')" />
                <el-option value="error" :label="t('admin.systemtools.system.level_error')" />
                <el-option value="warn" :label="t('admin.systemtools.system.level_warn')" />
                <el-option value="info" :label="t('admin.systemtools.system.level_info')" />
                <el-option value="debug" :label="t('admin.systemtools.system.level_debug')" />
                <el-option value="trace" :label="t('admin.systemtools.system.level_trace')" />
              </el-select>
            </el-form-item>
          </template>
          <template v-if="config.system['db-type'] === 'oracle'">
            <el-form-item label="">
              <h3>oracle</h3>
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.username')">
              <el-input
                v-model.trim="config.oracle.username"
                :placeholder="t('admin.systemtools.system.enter_username')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.password')">
              <el-input
                v-model.trim="config.oracle.password"
                :placeholder="t('admin.systemtools.system.enter_password')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.address')">
              <el-input
                v-model.trim="config.oracle.path"
                :placeholder="t('admin.systemtools.system.enter_address')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.database_name')">
              <el-input
                v-model.trim="config.oracle['db-name']"
                :placeholder="t('admin.systemtools.system.enter_database_name')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.prefix_label')">
              <el-input
                v-model.trim="config.oracle['prefix']"
                :placeholder="t('admin.systemtools.system.optional')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.singular_table')">
              <el-switch v-model="config.oracle['singular']" />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.engine')">
              <el-input
                v-model.trim="config.oracle['engine']"
                :placeholder="t('admin.systemtools.system.default_innodb')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.max_idle_conns')">
              <el-input-number
                v-model="config.oracle['max-idle-conns']"
                :min="1"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.max_open_conns')">
              <el-input-number
                v-model="config.oracle['max-open-conns']"
                :min="1"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.log_writes')">
              <el-switch v-model="config.oracle['log-zap']" />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.log_mode')">
              <el-select v-model="config.oracle['log-mode']">
                <el-option value="off" :label="t('admin.systemtools.system.level_off')" />
                <el-option value="fatal" :label="t('admin.systemtools.system.level_fatal')" />
                <el-option value="error" :label="t('admin.systemtools.system.level_error')" />
                <el-option value="warn" :label="t('admin.systemtools.system.level_warn')" />
                <el-option value="info" :label="t('admin.systemtools.system.level_info')" />
                <el-option value="debug" :label="t('admin.systemtools.system.level_debug')" />
                <el-option value="trace" :label="t('admin.systemtools.system.level_trace')" />
              </el-select>
            </el-form-item>
          </template>
        </el-tab-pane>
        <el-tab-pane :label="t('admin.systemtools.system.oss')" name="10" class="mt-3.5">
          <template v-if="config.system['oss-type'] === 'local'">
            <h2>{{ t('admin.systemtools.system.oss_local') }}</h2>
            <el-form-item :label="t('admin.systemtools.system.public_path')">
              <el-input
                v-model.trim="config.local.path"
                :placeholder="t('admin.systemtools.system.enter_public_path')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.storage_path')">
              <el-input
                v-model.trim="config.local['store-path']"
                :placeholder="t('admin.systemtools.system.enter_storage_path')"
              />
            </el-form-item>
          </template>
          <template v-if="config.system['oss-type'] === 'qiniu'">
            <h2>{{ t('admin.systemtools.system.oss_qiniu') }}</h2>
            <el-form-item :label="t('admin.systemtools.system.zone')">
              <el-input
                v-model.trim="config.qiniu.zone"
                :placeholder="t('admin.systemtools.system.enter_zone')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.bucket')">
              <el-input
                v-model.trim="config.qiniu.bucket"
                :placeholder="t('admin.systemtools.system.enter_bucket')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.cdn_domain')">
              <el-input
                v-model.trim="config.qiniu['img-path']"
                :placeholder="t('admin.systemtools.system.enter_cdn_domain')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.use_https')">
              <el-switch v-model="config.qiniu['use-https']">{{ t('admin.systemtools.system.on') }}</el-switch>
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.access_key')">
              <el-input
                v-model.trim="config.qiniu['access-key']"
                :placeholder="t('admin.systemtools.system.enter_access_key')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.secret_key')">
              <el-input
                v-model.trim="config.qiniu['secret-key']"
                :placeholder="t('admin.systemtools.system.enter_secret_key')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.use_cdn_upload')">
              <el-switch v-model="config.qiniu['use-cdn-domains']" />
            </el-form-item>
          </template>
          <template v-if="config.system['oss-type'] === 'tencent-cos'">
            <h2>{{ t('admin.systemtools.system.oss_tencent') }}</h2>
            <el-form-item :label="t('admin.systemtools.system.bucket')">
              <el-input
                v-model.trim="config['tencent-cos']['bucket']"
                :placeholder="t('admin.systemtools.system.enter_bucket')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.region')">
              <el-input
                v-model.trim="config['tencent-cos'].region"
                :placeholder="t('admin.systemtools.system.enter_region')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.secret_id')">
              <el-input
                v-model.trim="config['tencent-cos']['secret-id']"
                :placeholder="t('admin.systemtools.system.enter_secret_id')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.secret_key')">
              <el-input
                v-model.trim="config['tencent-cos']['secret-key']"
                :placeholder="t('admin.systemtools.system.enter_secret_key')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.path_prefix')">
              <el-input
                v-model.trim="config['tencent-cos']['path-prefix']"
                :placeholder="t('admin.systemtools.system.enter_path_prefix')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.base_url')">
              <el-input
                v-model.trim="config['tencent-cos']['base-url']"
                :placeholder="t('admin.systemtools.system.enter_base_url')"
              />
            </el-form-item>
          </template>
          <template v-if="config.system['oss-type'] === 'aliyun-oss'">
            <h2>{{ t('admin.systemtools.system.oss_aliyun') }}</h2>
            <el-form-item :label="t('admin.systemtools.system.endpoint')">
              <el-input
                v-model.trim="config['aliyun-oss'].endpoint"
                :placeholder="t('admin.systemtools.system.enter_endpoint')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.access_key_id')">
              <el-input
                v-model.trim="config['aliyun-oss']['access-key-id']"
                :placeholder="t('admin.systemtools.system.enter_access_key_id')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.access_key_secret')">
              <el-input
                v-model.trim="config['aliyun-oss']['access-key-secret']"
                :placeholder="t('admin.systemtools.system.enter_access_key_secret')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.bucket')">
              <el-input
                v-model.trim="config['aliyun-oss']['bucket-name']"
                :placeholder="t('admin.systemtools.system.enter_bucket')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.bucket_url')">
              <el-input
                v-model.trim="config['aliyun-oss']['bucket-url']"
                :placeholder="t('admin.systemtools.system.enter_bucket_url')"
              />
            </el-form-item>
          </template>
          <template v-if="config.system['oss-type'] === 'huawei-obs'">
            <h2>{{ t('admin.systemtools.system.oss_huawei') }}</h2>
            <el-form-item :label="t('admin.systemtools.system.path')">
              <el-input
                v-model.trim="config['hua-wei-obs'].path"
                :placeholder="t('admin.systemtools.system.enter_path')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.bucket')">
              <el-input
                v-model.trim="config['hua-wei-obs'].bucket"
                :placeholder="t('admin.systemtools.system.enter_bucket')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.endpoint')">
              <el-input
                v-model.trim="config['hua-wei-obs'].endpoint"
                :placeholder="t('admin.systemtools.system.enter_endpoint')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.access_key')">
              <el-input
                v-model.trim="config['hua-wei-obs']['access-key']"
                :placeholder="t('admin.systemtools.system.enter_access_key')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.secret_key')">
              <el-input
                v-model.trim="config['hua-wei-obs']['secret-key']"
                :placeholder="t('admin.systemtools.system.enter_secret_key')"
              />
            </el-form-item>
          </template>
          <template v-if="config.system['oss-type'] === 'cloudflare-r2'">
            <h2>{{ t('admin.systemtools.system.oss_cloudflare') }}</h2>
            <el-form-item :label="t('admin.systemtools.system.path')">
              <el-input
                v-model.trim="config['cloudflare-r2'].path"
                :placeholder="t('admin.systemtools.system.enter_path')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.bucket')">
              <el-input
                v-model.trim="config['cloudflare-r2'].bucket"
                :placeholder="t('admin.systemtools.system.enter_bucket')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.base_url')">
              <el-input
                v-model.trim="config['cloudflare-r2']['base-url']"
                :placeholder="t('admin.systemtools.system.enter_base_url')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.account_id')">
              <el-input
                v-model.trim="config['cloudflare-r2']['account-id']"
                :placeholder="t('admin.systemtools.system.enter_account_id')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.access_key_id2')">
              <el-input
                v-model.trim="config['cloudflare-r2']['access-key-id']"
                :placeholder="t('admin.systemtools.system.enter_access_key_id2')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.secret_access_key')">
              <el-input
                v-model.trim="config['cloudflare-r2']['secret-access-key']"
                :placeholder="t('admin.systemtools.system.enter_secret_access_key')"
              />
            </el-form-item>
          </template>
          <template v-if="config.system['oss-type'] === 'minio'">
            <h2>{{ t('admin.systemtools.system.minio') }}</h2>
            <el-form-item :label="t('admin.systemtools.system.endpoint')">
              <el-input
                v-model.trim="config.minio.endpoint"
                :placeholder="t('admin.systemtools.system.minio_endpoint_placeholder')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.access_key_id2')">
              <el-input
                v-model.trim="config.minio['access-key-id']"
                :placeholder="t('admin.systemtools.system.enter_access_key_id')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.access_key_secret2')">
              <el-input
                v-model.trim="config.minio['access-key-secret']"
                :placeholder="t('admin.systemtools.system.enter_access_key_secret2')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.bucket')">
              <el-input
                v-model.trim="config.minio['bucket-name']"
                :placeholder="t('admin.systemtools.system.enter_bucket')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.bucket_url')">
              <el-input
                v-model.trim="config.minio['bucket-url']"
                :placeholder="t('admin.systemtools.system.enter_bucket_url')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.base_path')">
              <el-input
                v-model.trim="config.minio['base-path']"
                :placeholder="t('admin.systemtools.system.enter_base_path')"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.system.enable_ssl')">
              <el-switch v-model="config.minio['use-ssl']" />
            </el-form-item>
          </template>
        </el-tab-pane>
        <el-tab-pane :label="t('admin.systemtools.system.excel_upload')" name="11" class="mt-3.5">
          <el-form-item :label="t('admin.systemtools.system.merge_output_dir')">
            <el-input
              v-model.trim="config.excel.dir"
              :placeholder="t('admin.systemtools.system.enter_output_dir')"
            />
          </el-form-item>
        </el-tab-pane>
        <el-tab-pane :label="t('admin.systemtools.system.autocode')" name="12" class="mt-3.5">
          <el-form-item :label="t('admin.systemtools.system.auto_restart_linux')">
            <el-switch v-model="config.autocode['transfer-restart']" />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.root_project_root')">
            <el-input v-model="config.autocode.root" disabled />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.server_backend_code_path')">
            <el-input
              v-model.trim="config.autocode['server']"
              :placeholder="t('admin.systemtools.system.enter_backend_code_path')"
            />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.sapi')">
            <el-input
              v-model.trim="config.autocode['server-api']"
              :placeholder="t('admin.systemtools.system.enter_backend_api_path')"
            />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.s_initialize')">
            <el-input
              v-model.trim="config.autocode['server-initialize']"
              :placeholder="t('admin.systemtools.system.enter_backend_init_path')"
            />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.s_model')">
            <el-input
              v-model.trim="config.autocode['server-model']"
              :placeholder="t('admin.systemtools.system.enter_backend_model_path')"
            />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.s_request')">
            <el-input
              v-model.trim="config.autocode['server-request']"
              :placeholder="t('admin.systemtools.system.enter_backend_request_path')"
            />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.s_router')">
            <el-input
              v-model.trim="config.autocode['server-router']"
              :placeholder="t('admin.systemtools.system.enter_backend_router_path')"
            />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.s_service')">
            <el-input
              v-model.trim="config.autocode['server-service']"
              :placeholder="t('admin.systemtools.system.enter_backend_service_path')"
            />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.web_frontend_path')">
            <el-input
              v-model.trim="config.autocode.web"
              :placeholder="t('admin.systemtools.system.enter_frontend_path')"
            />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.w_api')">
            <el-input
              v-model.trim="config.autocode['web-api']"
              :placeholder="t('admin.systemtools.system.enter_web_api_path')"
            />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.w_form')">
            <el-input
              v-model.trim="config.autocode['web-form']"
              :placeholder="t('admin.systemtools.system.enter_web_form_path')"
            />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.system.w_table')">
            <el-input
              v-model.trim="config.autocode['web-table']"
              :placeholder="t('admin.systemtools.system.enter_web_table_path')"
            />
          </el-form-item>
        </el-tab-pane>
      </el-tabs>
    </el-form>
    <div class="mt-4">
      <el-button type="primary" @click="update">{{ t('admin.systemtools.system.update') }}</el-button>
      <el-button type="primary" @click="reload">{{ t('admin.systemtools.system.reload_service') }}</el-button>
    </div>
  </div>
</template>

<script setup>
  import { getSystemConfig, reloadSystem, setSystemConfig } from '@/api/system'
  import { ref } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { Minus, Plus } from '@element-plus/icons-vue'
  import { emailTest } from '@/api/email'
  import { CreateUUID } from '@/utils/format'
  import { useI18n } from 'vue-i18n'

  const { t } = useI18n()

  defineOptions({
    name: 'Config'
  })

  const activeNames = ref('1')
  const config = ref({
    system: {
      'iplimit-count': 0,
      'iplimit-time': 0
    },
    jwt: {},
    mysql: {},
    mssql: {},
    sqlite: {},
    pgsql: {},
    oracle: {},
    excel: {},
    autocode: {},
    redis: {},
    mongo: {
      coll: '',
      options: '',
      database: '',
      username: '',
      password: '',
      'min-pool-size': '',
      'max-pool-size': '',
      'socket-timeout-ms': '',
      'connect-timeout-ms': '',
      'is-zap': false,
      hosts: [
        {
          host: '',
          port: ''
        }
      ]
    },
    qiniu: {},
    'tencent-cos': {},
    'aliyun-oss': {},
    'hua-wei-obs': {},
    'cloudflare-r2': {},
    minio: {},
    captcha: {},
    zap: {},
    local: {},
    email: {},
    timer: {
      detail: {}
    }
  })

  const initForm = async () => {
    const res = await getSystemConfig()
    if (res.code === 0) {
      config.value = res.data.config
    }
  }
  initForm()
  const reload = () => {
    ElMessageBox.confirm(t('admin.systemtools.system.reload_confirm'), t('admin.systemtools.system.warning'), {
      confirmButtonText: t('admin.common.confirm'),
      cancelButtonText: t('admin.common.cancel'),
      type: 'warning'
    })
      .then(async () => {
        const res = await reloadSystem()
        if (res.code === 0) {
          ElMessage({
            type: 'success',
            message: t('admin.systemtools.system.success')
          })
        }
      })
      .catch(() => {
        ElMessage({
          type: 'info',
          message: t('admin.systemtools.system.cancelled')
        })
      })
  }

  const update = async () => {
    const res = await setSystemConfig({ config: config.value })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: t('admin.systemtools.system.config_updated')
      })
      await initForm()
    }
  }

  const email = async () => {
    const res = await emailTest()
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: t('admin.systemtools.system.email_sent')
      })
      await initForm()
    } else {
      ElMessage({
        type: 'error',
        message: t('admin.systemtools.system.email_failed')
      })
    }
  }

  const getUUID = () => {
    config.value.jwt['signing-key'] = CreateUUID()
  }

  const addNode = () => {
    config.value.mongo.hosts.push({
      host: '',
      port: ''
    })
  }

  const removeNode = (index) => {
    config.value.mongo.hosts.splice(index, 1)
  }
</script>

<style lang="scss" scoped>
  .system {
    @apply bg-white p-9 rounded dark:bg-slate-900;
    h2 {
      @apply p-2.5 my-2.5 text-lg shadow;
    }
  }
</style>
