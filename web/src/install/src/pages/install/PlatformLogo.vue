<template>
  <v-stepper-content
    key="plot-form-content"
    :step="3"
  >
    <v-row justify="center">
      <v-col cols="12" sm="10" md="8" lg="6">
        <v-card
          class="mx-auto"
          max-width="500"
          :disabled="loading"
          :loading="loading"
        >
          <v-card-text>
            <v-file-input
              v-model="file"
              color="deep-purple accent-4"
              label="上传logo"
              placeholder="请选择文件"
              accept="image/png, image/jpeg"
              prepend-icon="mdi-camera"
              outlined
              show-size
              :rules="[rules.required]"
            >
              <template v-slot:selection="{ index, text }">
                <v-chip
                  v-if="index < 2"
                  color="deep-purple accent-4"
                  dark
                  label
                  small
                >
                  {{ text }}
                </v-chip>

                <span
                  v-else-if="index === 2"
                  class="overline grey--text text--darken-3 mx-2"
                >
                  +{{ files.length - 2 }} File(s)
                </span>
              </template>
            </v-file-input>
          </v-card-text>
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
  import { initLogo } from '../../api/install.service'

  export default {
    name: 'platform-logo',
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
        file: [],
        rules: {
          required: (value) => !!value || '请选择上传的logo文件.',
        }
      }
    },
    methods: {
      ...mapActions('app', ['showError', 'showSuccess']),
      onSubmit () {
        if (!this.file || this.file.length === 0) {
          this.showError({ error: { message: '请选择上传的logo文件' } })
          return
        }
        let formData = new FormData()
        formData.append('logo', this.file)
        this.loading = true
        initLogo(formData).then(res => {
          if (res && res.success) {
            this.nextStep('cors-step')
          }
        }).finally(() => {
          this.loading = false
        })
      }
    }
  }
</script>
