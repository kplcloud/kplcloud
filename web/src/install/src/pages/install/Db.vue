<template>
  <v-stepper-content
    key="db-content"
    :step="1"
  >
    <v-card
      class="mb-12"
      color="grey lighten-1"
      height="200px"
    >
      数据库设置
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
    methods: {
      ...mapActions('app', ['showError', 'showSuccess']),
      // nextStep (n) {
      //   console.log(n)
      // }
    }
  }
</script>
