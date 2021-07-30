<template>
  <v-stepper-content
    key="db-content"
    :step="2"
  >
    <v-row justify="center">
      <v-col cols="12" sm="10" md="8" lg="6">
        <v-card>
          <v-form
            ref="form"
            v-model="valid"
          >
            <v-card-text>
              <v-text-field
                ref="appName"
                v-model="appName"
                :rules="[() => !!appName || '名称为必填']"
                :error-messages="errorMessages"
                label="名称"
                placeholder="kplcloud"
                required
              ></v-text-field>
              <!--         uploadPath, debug-->
              <v-text-field
                ref="adminName"
                v-model="adminName"
                :rules="[() => !!adminName || '管理员为必填']"
                :error-messages="errorMessages"
                label="管理员账号"
                placeholder="admin@kplcloud.com"
                required
              ></v-text-field>
              <v-text-field
                v-model="adminPassword"
                :append-icon="passwordShow ? 'mdi-eye' : 'mdi-eye-off'"
                :rules="[rules.required, rules.min]"
                :type="passwordShow ? 'text' : 'password'"
                name="input-10-1"
                label="管理员密码"
                hint="至少4个字符"
                counter
                @click:append="passwordShow = !passwordShow"
              ></v-text-field>
              <v-text-field
                ref="domain"
                v-model="domain"
                :rules="[() => !!domain || '域名为必填']"
                :error-messages="errorMessages"
                label="域名"
                placeholder="https://kplcloud.nsini.com"
                required
              ></v-text-field>
              <v-text-field
                ref="domainSuffix"
                v-model="domainSuffix"
                :rules="[() => !!domainSuffix || '域名后缀为必填']"
                :error-messages="errorMessages"
                label="域名"
                placeholder="https://%s.%s.nsini.com"
                required
              ></v-text-field>
              <v-text-field
                ref="logPath"
                v-model="logPath"
                :error-messages="errorMessages"
                label="日志路径"
                placeholder="/var/log/kplcloud.log"
              ></v-text-field>
              <v-select
                :items="levels"
                v-model="logLevel"
                label="日志级别"
                required
              ></v-select>
              <v-text-field
                ref="uploadPath"
                v-model="uploadPath"
                :rules="[() => !!uploadPath || '文件上传路径为必填']"
                :error-messages="errorMessages"
                label="上传路径"
                placeholder="/data/upload"
                required
              ></v-text-field>
              <v-switch
                v-model="debug"
                :label="`Debug : ${debug.toString()}`"
              ></v-switch>
            </v-card-text>
            <v-card-actions>
              <v-spacer></v-spacer>
              <v-slide-x-reverse-transition>
                <v-tooltip
                  v-if="formHasErrors"
                  left
                >
                  <template v-slot:activator="{ on, attrs }">
                    <v-btn
                      icon
                      class="my-0"
                      v-bind="attrs"
                      @click="resetForm"
                      v-on="on"
                    >
                      <v-icon>mdi-refresh</v-icon>
                    </v-btn>
                  </template>
                  <span>取消</span>
                </v-tooltip>
              </v-slide-x-reverse-transition>
            </v-card-actions>
          </v-form>
        </v-card>
      </v-col>
    </v-row>


    <v-btn
      color="primary"
      @click="onSubmit"
    >
      下一步
    </v-btn>

    <v-btn text>取消</v-btn>
  </v-stepper-content>
</template>

<script>
  import { mapActions } from 'vuex'
  import { initPlatform } from '../../api/install.service'

  export default {
    name: 'platform',
    props: {
      title: {
        type: String,
        default: ''
      },
      nextStep: {
        type: Function,
      }
    },
    data () {
      return {
        loading: false,
        valid: true,
        appName: '',
        adminName: '',
        adminPassword: '',
        domain: '',
        domainSuffix: '',
        logPath: '',
        levels: ['all', 'debug', 'info', 'warning', 'error'],
        logLevel: '',
        uploadPath: '',
        passwordShow: false,
        debug: true,
        rules: {
          required: (value) => !!value || 'Required.',
          min: (v) => v.length >= 4 || 'Min 4 characters',
        },
        errorMessages: '',
        formHasErrors: ''
      }
    },
    computed: {},
    watch: {},

    methods: {
      ...mapActions('app', ['showError', 'showSuccess']),
      resetForm () {
        this.errorMessages = []
        this.formHasErrors = false

        Object.keys(this.form).forEach((f) => {
          this.$refs[f].reset()
        })
      },
      onSubmit () {
        this.formHasErrors = false
        if (!this.$refs.form.validate()) {
          this.formHasErrors = true
          return
        }
        this.loading = true
        initPlatform({
          appName: this.appName,
          adminName: this.adminName,
          adminPassword: this.adminPassword,
          domain: this.domain,
          domainSuffix: this.domainSuffix,
          logPath: this.logPath,
          logLevel: this.logLevel,
          uploadPath: this.uploadPath,
          debug: this.debug === true,
        }).then(res => {
          if (res && res.success) {
            // this.$router.push({})// install/platform
            this.nextStep(3)
          }
        }).finally(() => {
          this.loading = false
        })
      }
    }
  }
</script>
