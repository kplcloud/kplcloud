<template>
  <v-stepper-content
    key="db-content"
    :step="1"
  >
    <v-row justify="center">
      <v-col cols="12" sm="10" md="8" lg="6">
        <v-card>
          <v-form
            ref="form"
            v-model="valid"
          >
            <v-card-text>
              <v-select
                :items="drives"
                v-model="drive"
                label="数据库驱动"
                :rules="[() => !!drive || '数据库驱动为必选']"
                required
              ></v-select>
              <v-text-field
                ref="host"
                v-model="host"
                :rules="[() => !!host || '数据库地址为必填']"
                :error-messages="errorMessages"
                label="数据库地址"
                placeholder="127.0.0.1"
                required
              ></v-text-field>
              <v-text-field
                ref="port"
                v-model="port"
                :rules="[() => !!port || '数据库端口为必填']"
                :error-messages="errorMessages"
                label="数据库端口"
                placeholder="3306"
                required
              ></v-text-field>

              <v-text-field
                ref="user"
                v-model="user"
                :rules="[() => !!user || '数据库用名为必填']"
                :error-messages="errorMessages"
                label="数据库用户"
                placeholder="root"
                required
              ></v-text-field>

              <v-text-field
                v-model="password"
                :append-icon="passwordShow ? 'mdi-eye' : 'mdi-eye-off'"
                :rules="[rules.required, rules.min]"
                :type="passwordShow ? 'text' : 'password'"
                name="input-10-1"
                label="数据库密码"
                hint="至少8个字符"
                counter
                @click:append="passwordShow = !passwordShow"
              ></v-text-field>

              <v-text-field
                ref="database"
                v-model="database"
                :rules="[() => !!database || '数据库为必填']"
                :error-messages="errorMessages"
                label="数据库名"
                placeholder=""
                required
              ></v-text-field>
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
              <v-btn color="primary" text>测试一下</v-btn>
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
  import { initDb } from '../../api/install.service'

  export default {
    name: 'db',
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
        drives: ['mysql'],
        drive: 'mysql',
        host: '127.0.0.1',
        port: 3306,
        user: '',
        password: '',
        passwordShow: false,
        rules: {
          required: (value) => !!value || 'Required.',
          min: (v) => v.length >= 4 || 'Min 4 characters',
        },
        database: '',
        errorMessages: '',
        formHasErrors: ''
      }
    },
    computed: {
    },
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
        this.loading = true;
        initDb({
          drive: this.drive,
          host: this.host,
          port: parseInt(this.port),
          user: this.user,
          password: this.password,
          database: this.database,
        }).then(res => {
          if(res && res.success) {
            // this.$router.push({})// install/platform
            this.nextStep(2)
          }
        }).finally(() => {
          this.loading = false;
        })
      }
    }
  }
</script>
