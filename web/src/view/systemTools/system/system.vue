<template>
  <div class="system">
    <el-form ref="form" :model="config" label-width="240px">
      <!--  System start  -->
      <el-tabs v-model="activeNames">
        <el-tab-pane label="System" name="1" class="mt-3.5">
          <el-form-item label="Port">
            <el-input-number
              v-model="config.system.addr"
              placeholder="Enter port"
            />
          </el-form-item>
          <el-form-item label="Database type">
            <el-select v-model="config.system['db-type']" class="w-full">
              <el-option value="mysql" />
              <el-option value="pgsql" />
              <el-option value="mssql" />
              <el-option value="sqlite" />
              <el-option value="oracle" />
            </el-select>
          </el-form-item>
          <el-form-item label="OSS type">
            <el-select v-model="config.system['oss-type']" class="w-full">
              <el-option value="local" label="Local" />
              <el-option value="qiniu" label="Qiniu" />
              <el-option value="tencent-cos" label="Tencent COS" />
              <el-option value="aliyun-oss" label="Alibaba OSS" />
              <el-option value="huawei-obs" label="Huawei OBS" />
              <el-option value="cloudflare-r2" label="cloudflare R2" />
              <el-option value="minio">MinIO</el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="Multipoint login">
            <el-switch v-model="config.system['use-multipoint']" />
          </el-form-item>
          <el-form-item label="Enable Redis">
            <el-switch v-model="config.system['use-redis']" />
          </el-form-item>
          <el-form-item label="Enable MongoDB">
            <el-switch v-model="config.system['use-mongo']" />
          </el-form-item>
          <el-form-item label="Strict role mode">
            <el-switch v-model="config.system['use-strict-auth']" />
          </el-form-item>
          <el-form-item label="IP limit count">
            <el-input-number v-model.number="config.system['iplimit-count']" />
          </el-form-item>
          <el-form-item label="IP limit duration">
            <el-input-number v-model.number="config.system['iplimit-time']" />
          </el-form-item>
          <el-form-item label="Disable auto-migrate">
            <el-switch v-model="config.system['disable-auto-migrate']" />
          </el-form-item>
          <el-tooltip
            content="After changing this, also update VITE_BASE_PATH in the frontend env file."
            placement="top-start"
          >
            <el-form-item label="Global router prefix">
              <el-input
                v-model.trim="config.system['router-prefix']"
                placeholder="Enter router prefix"
              />
            </el-form-item>
          </el-tooltip>
        </el-tab-pane>
        <el-tab-pane label="JWT signing key" name="2" class="mt-3.5">
          <el-form-item label="JWT signing key">
            <el-input
              v-model.trim="config.jwt['signing-key']"
              placeholder="Enter signing key"
            >
              <template #append>
                <el-button @click="getUUID">Generate</el-button>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item label="Expires in">
            <el-input
              v-model.trim="config.jwt['expires-time']"
              placeholder="Enter expires in"
            />
          </el-form-item>
          <el-form-item label="Buffer time">
            <el-input
              v-model.trim="config.jwt['buffer-time']"
              placeholder="Enter buffer time"
            />
          </el-form-item>
          <el-form-item label="Issuer">
            <el-input
              v-model.trim="config.jwt.issuer"
              placeholder="Enter issuer"
            />
          </el-form-item>
        </el-tab-pane>
        <el-tab-pane label="Zap logger" name="3" class="mt-3.5">
          <el-form-item label="Level">
            <el-select v-model="config.zap.level">
              <el-option value="off" label="Off" />
              <el-option value="fatal" label="Fatal" />
              <el-option value="error" label="Error" />
              <el-option value="warn" label="Warn" />
              <el-option value="info" label="Info" />
              <el-option value="debug" label="Debug" />
              <el-option value="trace" label="Trace" />
            </el-select>
          </el-form-item>
          <el-form-item label="Format">
            <el-select v-model="config.zap.format">
              <el-option value="console" label="console" />
              <el-option value="json" label="json" />
            </el-select>
          </el-form-item>
          <el-form-item label="Prefix">
            <el-input
              v-model.trim="config.zap.prefix"
              placeholder="Enter prefix"
            />
          </el-form-item>
          <el-form-item label="Directory">
            <el-input
              v-model.trim="config.zap.director"
              placeholder="Enter directory"
            />
          </el-form-item>
          <el-form-item label="Encode level">
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
          <el-form-item label="Stacktrace key">
            <el-input
              v-model.trim="config.zap['stacktrace-key']"
              placeholder="Enter stacktrace key"
            />
          </el-form-item>
          <el-form-item label="Retention days">
            <el-input-number v-model="config.zap['retention-day']" />
          </el-form-item>
          <el-form-item label="Show line">
            <el-switch v-model="config.zap['show-line']" />
          </el-form-item>
          <el-form-item label="Log to console">
            <el-switch v-model="config.zap['log-in-console']" />
          </el-form-item>
        </el-tab-pane>
        <el-tab-pane
          label="Redis"
          name="4"
          class="mt-3.5"
          v-if="config.system['use-redis']"
        >
          <el-form-item label="DB">
            <el-input-number v-model="config.redis.db" min="0" max="16" />
          </el-form-item>
          <el-form-item label="Address">
            <el-input
              v-model.trim="config.redis.addr"
              placeholder="Enter address"
            />
          </el-form-item>
          <el-form-item label="Password">
            <el-input
              v-model.trim="config.redis.password"
              placeholder="Enter password"
            />
          </el-form-item>
        </el-tab-pane>
        <el-tab-pane label="Email" name="5" class="mt-3.5">
          <el-form-item label="To">
            <el-input
              v-model="config.email.to"
              placeholder="Multiple emails separated by commas"
            />
          </el-form-item>
          <el-form-item label="Port">
            <el-input-number v-model="config.email.port" />
          </el-form-item>
          <el-form-item label="From">
            <el-input
              v-model.trim="config.email.from"
              placeholder="Enter sender email"
            />
          </el-form-item>
          <el-form-item label="host">
            <el-input
              v-model.trim="config.email.host"
              placeholder="Enter host"
            />
          </el-form-item>
          <el-form-item label="Use SSL">
            <el-switch v-model="config.email['is-ssl']" />
          </el-form-item>
          <el-form-item label="Use LoginAuth">
            <el-switch v-model="config.email['is-loginauth']" />
          </el-form-item>
          <el-form-item label="secret">
            <el-input
              v-model.trim="config.email.secret"
              placeholder="Enter secret"
            />
          </el-form-item>
          <el-form-item label="Test email">
            <el-button @click="email">Send test email</el-button>
          </el-form-item>
        </el-tab-pane>
        <el-tab-pane
          label="MongoDB"
          name="14"
          class="mt-3.5"
          v-if="config.system['use-mongo']"
        >
          <el-form-item label="Collection name (optional)">
            <el-input
              v-model.trim="config.mongo.coll"
              placeholder="Enter collection name"
            />
          </el-form-item>
          <el-form-item label="MongoDB options">
            <el-input
              v-model.trim="config.mongo.options"
              placeholder="Enter MongoDB options"
            />
          </el-form-item>
          <el-form-item label="Database name">
            <el-input
              v-model.trim="config.mongo.database"
              placeholder="Enter database name"
            />
          </el-form-item>
          <el-form-item label="Username">
            <el-input
              v-model.trim="config.mongo.username"
              placeholder="Enter username"
            />
          </el-form-item>
          <el-form-item label="Password">
            <el-input
              v-model.trim="config.mongo.password"
              placeholder="Enter password"
            />
          </el-form-item>
          <el-form-item label="Min pool size">
            <el-input-number v-model="config.mongo['min-pool-size']" min="0" />
          </el-form-item>
          <el-form-item label="Max pool size">
            <el-input-number
              v-model="config.mongo['max-pool-size']"
              min="100"
            />
          </el-form-item>
          <el-form-item label="Socket timeout (ms)">
            <el-input-number
              v-model="config.mongo['socket-timeout-ms']"
              min="0"
            />
          </el-form-item>
          <el-form-item label="Connect timeout (ms)">
            <el-input-number
              v-model="config.mongo['socket-timeout-ms']"
              min="0"
            />
          </el-form-item>
          <el-form-item label="Enable Zap logs">
            <el-switch v-model="config.mongo['is-zap']" />
          </el-form-item>
          <el-form-item
            v-for="(item, k) in config.mongo.hosts"
            :key="k"
            :label="`Node ${k + 1}`"
          >
            <div v-for="(_, k2) in item" :key="k2">
              <el-form-item :key="k + k2" :label="k2" label-width="60">
                <el-input
                  v-model.trim="item[k2]"
                :placeholder="k2 === 'host' ? 'Enter host' : 'Enter port'"
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
        <el-tab-pane label="Captcha" name="7" class="mt-3.5">
          <el-form-item label="Key length">
            <el-input-number
              v-model="config.captcha['key-long']"
              :min="4"
              :max="6"
            />
          </el-form-item>
          <el-form-item label="Image width">
            <el-input-number v-model.number="config.captcha['img-width']" />
          </el-form-item>
          <el-form-item label="Image height">
            <el-input-number v-model.number="config.captcha['img-height']" />
          </el-form-item>
        </el-tab-pane>
        <el-tab-pane label="Database config" name="9" class="mt-3.5">
          <template v-if="config.system['db-type'] === 'mysql'">
            <el-form-item label="">
              <h3>MySQL</h3>
            </el-form-item>
            <el-form-item label="Username">
              <el-input
                v-model.trim="config.mysql.username"
                placeholder="Enter username"
              />
            </el-form-item>
            <el-form-item label="Password">
              <el-input
                v-model.trim="config.mysql.password"
                placeholder="Enter password"
              />
            </el-form-item>
            <el-form-item label="Address">
              <el-input
                v-model.trim="config.mysql.path"
                placeholder="Enter address"
              />
            </el-form-item>
            <el-form-item label="Database name">
              <el-input
                v-model.trim="config.mysql['db-name']"
                placeholder="Enter database name"
              />
            </el-form-item>
            <el-form-item label="Prefix">
              <el-input
                v-model.trim="config.mysql['prefix']"
                placeholder="Optional"
              />
            </el-form-item>
            <el-form-item label="Singular table">
              <el-switch v-model="config.mysql['singular']" />
            </el-form-item>
            <el-form-item label="Engine">
              <el-input
                v-model.trim="config.mysql['engine']"
                placeholder="Default: InnoDB"
              />
            </el-form-item>
            <el-form-item label="maxIdleConns">
              <el-input-number
                v-model="config.mysql['max-idle-conns']"
                :min="1"
              />
            </el-form-item>
            <el-form-item label="maxOpenConns">
              <el-input-number
                v-model="config.mysql['max-open-conns']"
                :min="1"
              />
            </el-form-item>
            <el-form-item label="Log writes">
              <el-switch v-model="config.mysql['log-zap']" />
            </el-form-item>
            <el-form-item label="Log mode">
              <el-select v-model="config.mysql['log-mode']">
                <el-option value="off" label="Off" />
                <el-option value="fatal" label="Fatal" />
                <el-option value="error" label="Error" />
                <el-option value="warn" label="Warn" />
                <el-option value="info" label="Info" />
                <el-option value="debug" label="Debug" />
                <el-option value="trace" label="Trace" />
              </el-select>
            </el-form-item>
          </template>
          <template v-if="config.system['db-type'] === 'pgsql'">
            <el-form-item label="">
              <h3>PostgreSQL</h3>
            </el-form-item>
            <el-form-item label="Username">
              <el-input
                v-model="config.pgsql.username"
                placeholder="Enter username"
              />
            </el-form-item>
            <el-form-item label="Password">
              <el-input
                v-model="config.pgsql.password"
                placeholder="Enter password"
              />
            </el-form-item>
            <el-form-item label="Address">
              <el-input
                v-model.trim="config.pgsql.path"
                placeholder="Enter address"
              />
            </el-form-item>
            <el-form-item label="Database">
              <el-input
                v-model.trim="config.pgsql['db-name']"
                placeholder="Enter database"
              />
            </el-form-item>
            <el-form-item label="Prefix">
              <el-input
                v-model.trim="config.pgsql['prefix']"
                placeholder="Enter prefix"
              />
            </el-form-item>
            <el-form-item label="Singular table">
              <el-switch v-model="config.pgsql['singular']" />
            </el-form-item>
            <el-form-item label="Engine">
              <el-input
                v-model.trim="config.pgsql['engine']"
                placeholder="Enter engine"
              />
            </el-form-item>
            <el-form-item label="maxIdleConns">
              <el-input-number v-model="config.pgsql['max-idle-conns']" />
            </el-form-item>
            <el-form-item label="maxOpenConns">
              <el-input-number v-model="config.pgsql['max-open-conns']" />
            </el-form-item>
            <el-form-item label="Log writes">
              <el-switch v-model="config.pgsql['log-zap']" />
            </el-form-item>
            <el-form-item label="Log mode">
              <el-select v-model="config.pgsql['log-mode']">
                <el-option value="off" label="Off" />
                <el-option value="fatal" label="Fatal" />
                <el-option value="error" label="Error" />
                <el-option value="warn" label="Warn" />
                <el-option value="info" label="Info" />
                <el-option value="debug" label="Debug" />
                <el-option value="trace" label="Trace" />
              </el-select>
            </el-form-item>
          </template>
          <template v-if="config.system['db-type'] === 'mssql'">
            <el-form-item label="">
              <h3>MsSQL</h3>
            </el-form-item>
            <el-form-item label="Username">
              <el-input
                v-model.trim="config.mssql.username"
                placeholder="Enter username"
              />
            </el-form-item>
            <el-form-item label.trim="Password">
              <el-input
                v-model.trim="config.mssql.password"
                placeholder="Enter password"
              />
            </el-form-item>
            <el-form-item label="Address">
              <el-input
                v-model.trim="config.mssql.path"
                placeholder="Enter address"
              />
            </el-form-item>
            <el-form-item label="Port">
              <el-input
                v-model.trim="config.mssql.port"
                placeholder="Enter port"
              />
            </el-form-item>
            <el-form-item label="Database">
              <el-input
                v-model.trim="config.mssql['db-name']"
                placeholder="Enter database"
              />
            </el-form-item>
            <el-form-item label="Prefix">
              <el-input
                v-model.trim="config.mssql['prefix']"
                placeholder="Enter prefix"
              />
            </el-form-item>
            <el-form-item label="Singular table">
              <el-switch v-model="config.mssql['singular']" />
            </el-form-item>
            <el-form-item label="Engine">
              <el-input
                v-model.trim="config.mssql['engine']"
                placeholder="Enter engine"
              />
            </el-form-item>
            <el-form-item label="maxIdleConns">
              <el-input-number v-model="config.mssql['max-idle-conns']" />
            </el-form-item>
            <el-form-item label="maxOpenConns">
              <el-input-number v-model="config.mssql['max-open-conns']" />
            </el-form-item>
            <el-form-item label="Log writes">
              <el-switch v-model="config.mssql['log-zap']" />
            </el-form-item>
            <el-form-item label="Log mode">
              <el-select v-model="config.mssql['log-mode']">
                <el-option value="off" label="Off" />
                <el-option value="fatal" label="Fatal" />
                <el-option value="error" label="Error" />
                <el-option value="warn" label="Warn" />
                <el-option value="info" label="Info" />
                <el-option value="debug" label="Debug" />
                <el-option value="trace" label="Trace" />
              </el-select>
            </el-form-item>
          </template>
          <template v-if="config.system['db-type'] === 'sqlite'">
            <el-form-item label="">
              <h3>sqlite</h3>
            </el-form-item>
            <el-form-item label="Username">
              <el-input
                v-model.trim="config.sqlite.username"
                placeholder="Enter username"
              />
            </el-form-item>
            <el-form-item label="Password">
              <el-input
                v-model.trim="config.sqlite.password"
                placeholder="Enter password"
              />
            </el-form-item>
            <el-form-item label="Address">
              <el-input
                v-model.trim="config.sqlite.path"
                placeholder="Enter address"
              />
            </el-form-item>
            <el-form-item label="Port">
              <el-input
                v-model.trim="config.sqlite.port"
                placeholder="Enter port"
              />
            </el-form-item>
            <el-form-item label="Database">
              <el-input
                v-model.trim="config.sqlite['db-name']"
                placeholder="Enter database"
              />
            </el-form-item>
            <el-form-item label="maxIdleConns">
              <el-input-number v-model="config.sqlite['max-idle-conns']" />
            </el-form-item>
            <el-form-item label="maxOpenConns">
              <el-input-number v-model="config.sqlite['max-open-conns']" />
            </el-form-item>
            <el-form-item label="Log writes">
              <el-switch v-model="config.sqlite['log-zap']" />
            </el-form-item>
            <el-form-item label="Log mode">
              <el-select v-model="config.sqlite['log-mode']">
                <el-option value="off" label="Off" />
                <el-option value="fatal" label="Fatal" />
                <el-option value="error" label="Error" />
                <el-option value="warn" label="Warn" />
                <el-option value="info" label="Info" />
                <el-option value="debug" label="Debug" />
                <el-option value="trace" label="Trace" />
              </el-select>
            </el-form-item>
          </template>
          <template v-if="config.system['db-type'] === 'oracle'">
            <el-form-item label="">
              <h3>oracle</h3>
            </el-form-item>
            <el-form-item label="Username">
              <el-input
                v-model.trim="config.oracle.username"
                placeholder="Enter username"
              />
            </el-form-item>
            <el-form-item label="Password">
              <el-input
                v-model.trim="config.oracle.password"
                placeholder="Enter password"
              />
            </el-form-item>
            <el-form-item label="Address">
              <el-input
                v-model.trim="config.oracle.path"
                placeholder="Enter address"
              />
            </el-form-item>
            <el-form-item label="Database name">
              <el-input
                v-model.trim="config.oracle['db-name']"
                placeholder="Enter database name"
              />
            </el-form-item>
            <el-form-item label="Prefix">
              <el-input
                v-model.trim="config.oracle['prefix']"
                placeholder="Optional"
              />
            </el-form-item>
            <el-form-item label="Singular table">
              <el-switch v-model="config.oracle['singular']" />
            </el-form-item>
            <el-form-item label="Engine">
              <el-input
                v-model.trim="config.oracle['engine']"
                placeholder="Default: InnoDB"
              />
            </el-form-item>
            <el-form-item label="maxIdleConns">
              <el-input-number
                v-model="config.oracle['max-idle-conns']"
                :min="1"
              />
            </el-form-item>
            <el-form-item label="maxOpenConns">
              <el-input-number
                v-model="config.oracle['max-open-conns']"
                :min="1"
              />
            </el-form-item>
            <el-form-item label="Log writes">
              <el-switch v-model="config.oracle['log-zap']" />
            </el-form-item>
            <el-form-item label="Log mode">
              <el-select v-model="config.oracle['log-mode']">
                <el-option value="off" label="Off" />
                <el-option value="fatal" label="Fatal" />
                <el-option value="error" label="Error" />
                <el-option value="warn" label="Warn" />
                <el-option value="info" label="Info" />
                <el-option value="debug" label="Debug" />
                <el-option value="trace" label="Trace" />
              </el-select>
            </el-form-item>
          </template>
        </el-tab-pane>
        <el-tab-pane label="OSS" name="10" class="mt-3.5">
          <template v-if="config.system['oss-type'] === 'local'">
            <h2>Local</h2>
            <el-form-item label="Public path">
              <el-input
                v-model.trim="config.local.path"
                placeholder="Enter public path"
              />
            </el-form-item>
            <el-form-item label="Storage path">
              <el-input
                v-model.trim="config.local['store-path']"
                placeholder="Enter storage path"
              />
            </el-form-item>
          </template>
          <template v-if="config.system['oss-type'] === 'qiniu'">
            <h2>Qiniu</h2>
            <el-form-item label="Zone">
              <el-input
                v-model.trim="config.qiniu.zone"
                placeholder="Enter zone"
              />
            </el-form-item>
            <el-form-item label="Bucket">
              <el-input
                v-model.trim="config.qiniu.bucket"
                placeholder="Enter bucket"
              />
            </el-form-item>
            <el-form-item label="CDN domain">
              <el-input
                v-model.trim="config.qiniu['img-path']"
                placeholder="Enter CDN domain"
              />
            </el-form-item>
            <el-form-item label="Use HTTPS">
              <el-switch v-model="config.qiniu['use-https']">On</el-switch>
            </el-form-item>
            <el-form-item label="accessKey">
              <el-input
                v-model.trim="config.qiniu['access-key']"
                placeholder="Enter accessKey"
              />
            </el-form-item>
            <el-form-item label="secretKey">
              <el-input
                v-model.trim="config.qiniu['secret-key']"
                placeholder="Enter secretKey"
              />
            </el-form-item>
            <el-form-item label="Use CDN domains for upload">
              <el-switch v-model="config.qiniu['use-cdn-domains']" />
            </el-form-item>
          </template>
          <template v-if="config.system['oss-type'] === 'tencent-cos'">
            <h2>Tencent COS</h2>
            <el-form-item label="Bucket">
              <el-input
                v-model.trim="config['tencent-cos']['bucket']"
                placeholder="Enter bucket"
              />
            </el-form-item>
            <el-form-item label="Region">
              <el-input
                v-model.trim="config['tencent-cos'].region"
                placeholder="Enter region"
              />
            </el-form-item>
            <el-form-item label="secretID">
              <el-input
                v-model.trim="config['tencent-cos']['secret-id']"
                placeholder="Enter secretID"
              />
            </el-form-item>
            <el-form-item label="secretKey">
              <el-input
                v-model.trim="config['tencent-cos']['secret-key']"
                placeholder="Enter secretKey"
              />
            </el-form-item>
            <el-form-item label="Path prefix">
              <el-input
                v-model.trim="config['tencent-cos']['path-prefix']"
                placeholder="Enter path prefix"
              />
            </el-form-item>
            <el-form-item label="Base URL">
              <el-input
                v-model.trim="config['tencent-cos']['base-url']"
                placeholder="Enter base URL"
              />
            </el-form-item>
          </template>
          <template v-if="config.system['oss-type'] === 'aliyun-oss'">
            <h2>Alibaba OSS</h2>
            <el-form-item label="Endpoint">
              <el-input
                v-model.trim="config['aliyun-oss'].endpoint"
                placeholder="Enter endpoint"
              />
            </el-form-item>
            <el-form-item label="accessKeyId">
              <el-input
                v-model.trim="config['aliyun-oss']['access-key-id']"
                placeholder="Enter accessKeyId"
              />
            </el-form-item>
            <el-form-item label="accessKeySecret">
              <el-input
                v-model.trim="config['aliyun-oss']['access-key-secret']"
                placeholder="Enter accessKeySecret"
              />
            </el-form-item>
            <el-form-item label="Bucket">
              <el-input
                v-model.trim="config['aliyun-oss']['bucket-name']"
                placeholder="Enter bucket"
              />
            </el-form-item>
            <el-form-item label="Bucket URL">
              <el-input
                v-model.trim="config['aliyun-oss']['bucket-url']"
                placeholder="Enter bucket URL"
              />
            </el-form-item>
          </template>
          <template v-if="config.system['oss-type'] === 'huawei-obs'">
            <h2>Huawei OBS</h2>
            <el-form-item label="Path">
              <el-input
                v-model.trim="config['hua-wei-obs'].path"
                placeholder="Enter path"
              />
            </el-form-item>
            <el-form-item label="Bucket">
              <el-input
                v-model.trim="config['hua-wei-obs'].bucket"
                placeholder="Enter bucket"
              />
            </el-form-item>
            <el-form-item label="Endpoint">
              <el-input
                v-model.trim="config['hua-wei-obs'].endpoint"
                placeholder="Enter endpoint"
              />
            </el-form-item>
            <el-form-item label="accessKey">
              <el-input
                v-model.trim="config['hua-wei-obs']['access-key']"
                placeholder="Enter accessKey"
              />
            </el-form-item>
            <el-form-item label="secretKey">
              <el-input
                v-model.trim="config['hua-wei-obs']['secret-key']"
                placeholder="Enter secretKey"
              />
            </el-form-item>
          </template>
          <template v-if="config.system['oss-type'] === 'cloudflare-r2'">
            <h2>Cloudflare R2</h2>
            <el-form-item label="Path">
              <el-input
                v-model.trim="config['cloudflare-r2'].path"
                placeholder="Enter path"
              />
            </el-form-item>
            <el-form-item label="Bucket">
              <el-input
                v-model.trim="config['cloudflare-r2'].bucket"
                placeholder="Enter bucket"
              />
            </el-form-item>
            <el-form-item label="Base URL">
              <el-input
                v-model.trim="config['cloudflare-r2']['base-url']"
                placeholder="Enter Base URL"
              />
            </el-form-item>
            <el-form-item label="Account ID">
              <el-input
                v-model.trim="config['cloudflare-r2']['account-id']"
                placeholder="Enter account ID"
              />
            </el-form-item>
            <el-form-item label="Access Key ID">
              <el-input
                v-model.trim="config['cloudflare-r2']['access-key-id']"
                placeholder="Enter access key ID"
              />
            </el-form-item>
            <el-form-item label="Secret Access Key">
              <el-input
                v-model.trim="config['cloudflare-r2']['secret-access-key']"
                placeholder="Enter secret access key"
              />
            </el-form-item>
          </template>
          <template v-if="config.system['oss-type'] === 'minio'">
            <h2>MinIO</h2>
            <el-form-item label="Endpoint">
              <el-input
                v-model.trim="config.minio.endpoint"
                placeholder="Enter endpoint, e.g. 127.0.0.1:9000"
              />
            </el-form-item>
            <el-form-item label="Access Key ID">
              <el-input
                v-model.trim="config.minio['access-key-id']"
                placeholder="Enter Access Key ID"
              />
            </el-form-item>
            <el-form-item label="Access Key Secret">
              <el-input
                v-model.trim="config.minio['access-key-secret']"
                placeholder="Enter Access Key Secret"
              />
            </el-form-item>
            <el-form-item label="Bucket">
              <el-input
                v-model.trim="config.minio['bucket-name']"
                placeholder="Enter bucket"
              />
            </el-form-item>
            <el-form-item label="Bucket URL">
              <el-input
                v-model.trim="config.minio['bucket-url']"
                placeholder="Enter bucket URL"
              />
            </el-form-item>
            <el-form-item label="Base Path">
              <el-input
                v-model.trim="config.minio['base-path']"
                placeholder="Enter Base Path"
              />
            </el-form-item>
            <el-form-item label="Enable SSL">
              <el-switch v-model="config.minio['use-ssl']" />
            </el-form-item>
          </template>
        </el-tab-pane>
        <el-tab-pane label="Excel upload" name="11" class="mt-3.5">
          <el-form-item label="Merge output dir">
            <el-input
              v-model.trim="config.excel.dir"
              placeholder="Enter output dir"
            />
          </el-form-item>
        </el-tab-pane>
        <el-tab-pane label="Autocode" name="12" class="mt-3.5">
          <el-form-item label="Auto restart (linux)">
            <el-switch v-model="config.autocode['transfer-restart']" />
          </el-form-item>
          <el-form-item label="root (project root)">
            <el-input v-model="config.autocode.root" disabled />
          </el-form-item>
          <el-form-item label="Server (backend code path)">
            <el-input
              v-model.trim="config.autocode['server']"
              placeholder="Enter backend code path"
            />
          </el-form-item>
          <el-form-item label="SApi (backend api path)">
            <el-input
              v-model.trim="config.autocode['server-api']"
              placeholder="Enter backend api path"
            />
          </el-form-item>
          <el-form-item label="SInitialize (backend initialize path)">
            <el-input
              v-model.trim="config.autocode['server-initialize']"
              placeholder="Enter backend initialize path"
            />
          </el-form-item>
          <el-form-item label="SModel (backend model path)">
            <el-input
              v-model.trim="config.autocode['server-model']"
              placeholder="Enter backend model path"
            />
          </el-form-item>
          <el-form-item label="SRequest (backend request path)">
            <el-input
              v-model.trim="config.autocode['server-request']"
              placeholder="Enter backend request path"
            />
          </el-form-item>
          <el-form-item label="SRouter (backend router path)">
            <el-input
              v-model.trim="config.autocode['server-router']"
              placeholder="Enter backend router path"
            />
          </el-form-item>
          <el-form-item label="SService (backend service path)">
            <el-input
              v-model.trim="config.autocode['server-service']"
              placeholder="Enter backend service path"
            />
          </el-form-item>
          <el-form-item label="Web (frontend path)">
            <el-input
              v-model.trim="config.autocode.web"
              placeholder="Enter frontend path"
            />
          </el-form-item>
          <el-form-item label="WApi (web api path)">
            <el-input
              v-model.trim="config.autocode['web-api']"
              placeholder="Enter web api path"
            />
          </el-form-item>
          <el-form-item label="WForm (web form path)">
            <el-input
              v-model.trim="config.autocode['web-form']"
              placeholder="Enter web form path"
            />
          </el-form-item>
          <el-form-item label="WTable (web table path)">
            <el-input
              v-model.trim="config.autocode['web-table']"
              placeholder="Enter web table path"
            />
          </el-form-item>
        </el-tab-pane>
      </el-tabs>
    </el-form>
    <div class="mt-4">
      <el-button type="primary" @click="update">Update</el-button>
      <el-button type="primary" @click="reload">Reload service</el-button>
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
    ElMessageBox.confirm('Reload service?', 'Warning', {
      confirmButtonText: 'Confirm',
      cancelButtonText: 'Cancel',
      type: 'warning'
    })
      .then(async () => {
        const res = await reloadSystem()
        if (res.code === 0) {
          ElMessage({
            type: 'success',
            message: 'Success'
          })
        }
      })
      .catch(() => {
        ElMessage({
          type: 'info',
          message: 'Cancelled'
        })
      })
  }

  const update = async () => {
    const res = await setSystemConfig({ config: config.value })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: 'Config updated'
      })
      await initForm()
    }
  }

  const email = async () => {
    const res = await emailTest()
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: 'Email sent'
      })
      await initForm()
    } else {
      ElMessage({
        type: 'error',
        message: 'Email failed'
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
