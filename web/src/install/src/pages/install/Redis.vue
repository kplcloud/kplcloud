<template>
  <v-stepper-content
    key="db-content"
    :step="5"
  >
    <v-row justify="center">
      <v-col cols="12" sm="10" md="8" lg="6">
        <v-card
          :disabled="loading"
          :loading="loading"
        >
          <v-form
            ref="form"
            v-model="valid"
          >
            <v-card-text>
              <v-text-field
                ref="hosts"
                v-model="hosts"
                :rules="[() => !!hosts || 'hosts为必填']"
                :error-messages="errorMessages"
                label="Hosts"
                placeholder="127.0.0.1:6379 ps:集群用','隔开"
                required
                outlined
              ></v-text-field>
              <v-text-field
                v-model="password"
                :append-icon="passwordShow ? 'mdi-eye' : 'mdi-eye-off'"
                :type="passwordShow ? 'text' : 'password'"
                name="input-10-1"
                label="密码"
                hint="可不填"
                counter
                outlined
                @click:append="passwordShow = !passwordShow"
              ></v-text-field>
              <v-text-field
                ref="database"
                v-model="database"
                label="DB"
                placeholder="0"
                outlined
              ></v-text-field>
              <v-text-field
                ref="prefix"
                v-model="prefix"
                label="前缀"
                placeholder="kplcloud"
                outlined
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
  import { initRedis } from '../../api/install.service'

  export default {
    name: 'redis',
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
        hosts: '',
        password: '',
        passwordShow: false,
        database: 0,
        prefix: 'kplcloud',

        errorMessages: '',
        formHasErrors: ''
      }
    },
    computed: {},
    watch: {},

    methods: {
      ...mapActions('app', ['showError', 'showSuccess']),
      onSubmit () {
        if (this.allow === true) {
          this.formHasErrors = false
          if (!this.$refs.form.validate()) {
            this.formHasErrors = true
            return
          }
        }

        this.loading = true
        initRedis({
          hosts: this.hosts,
          password: this.password,
          database: parseInt(this.database),
          prefix: this.prefix,
        }).then(res => {
          if (res && res.success) {
            this.nextStep('build-step')
          }
        }).finally(() => {
          this.loading = false
        })
      }
    }
  }
</script>
