<template>
  <v-stepper-content
    key="db-content"
    :step="6"
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
              <v-switch
                v-model="allow"
                :label="`允许跨域 : ${allow.toString()}`"
              ></v-switch>
              <v-text-field
                ref="origin"
                v-if="allow===true"
                v-model="origin"
                :rules="[() => !!origin || 'origin为必填']"
                :error-messages="errorMessages"
                label="Origin"
                placeholder="*"
                required
                outlined
              ></v-text-field>
              <v-select
                ref="methods"
                v-model="methods"
                :items="methodItems"
                chips
                label="Methods"
                multiple
                outlined
                v-if="allow===true"
              ></v-select>
              <v-text-field
                ref="headers"
                v-if="allow===true"
                v-model="headers"
                :rules="[() => !!headers || 'headers为必填']"
                :error-messages="errorMessages"
                label="Headers"
                placeholder="Origin,Content-Type,Authorization,mode,cors,Token"
                required
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
  import { initCors } from '../../api/install.service'

  export default {
    name: 'build',
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
        allow: false,
        origin: '',
        methods: [],
        headers: '',

        methodItems: ['GET', 'POST', 'OPTIONS', 'PUT', 'DELETE'],
        errorMessages: '',
        formHasErrors: ''
      }
    },
    computed: {},
    watch: {},

    methods: {
      ...mapActions('app', ['showError', 'showSuccess']),
      onSubmit () {
        this.nextStep('repo-step')
        return
        if (this.allow === true) {
          this.formHasErrors = false
          if (!this.$refs.form.validate()) {
            this.formHasErrors = true
            return
          }
        }

        this.loading = true
        initCors({
          allow: this.allow,
          methods: this.methods,
          headers: this.headers,
          origin: this.origin,
        }).then(res => {
          if (res && res.success) {
            this.nextStep('repo-step')
          }
        }).finally(() => {
          this.loading = false
        })
      }
    }
  }
</script>
