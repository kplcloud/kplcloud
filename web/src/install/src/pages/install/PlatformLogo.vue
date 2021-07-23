<template>
  <v-stepper-content
    key="plot-form-content"
    :step="3"
  >
    <v-card
      class="mb-12"
      color="grey lighten-1"
      height="200px"
    >
      平台配置
    </v-card>

    <v-btn
      color="primary"
      @click="nextStep(2)"
    >
      下一步
    </v-btn>

    <v-btn text>取消</v-btn>
  </v-stepper-content>
</template>

<script>
  import { mapActions } from 'vuex'

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
    methods: {
      ...mapActions('app', ['showError', 'showSuccess']),
    }
  }
</script>
