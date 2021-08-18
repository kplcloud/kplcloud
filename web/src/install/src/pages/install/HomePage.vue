<template>
  <div class="w-full">
    <v-card class="pa-6" flat>
      <v-stepper v-model="e1">
        <v-stepper-header>
          <template v-for="st in steps">
            <v-stepper-step :key="st.key" :complete="e1 > st.step" :step="st.step">{{st.title}}</v-stepper-step>
            <v-divider :key="st.step"></v-divider>
          </template>

        </v-stepper-header>

        <v-stepper-items>
          <db :nextStep="nextStep" :v-if="e1===1"/>
          <platform :nextStep="nextStep" :v-if="e1===2"/>
          <platform-logo :nextStep="nextStep" :v-if="e1===3"/>
          <cors :nextStep="nextStep" :v-if="e1===4"/>
          <redis :nextStep="nextStep" :v-if="e1===5"/>
          <build :nextStep="nextStep" :v-if="e1===6"/>
          <repo :nextStep="nextStep" :v-if="e1===7"/>

        </v-stepper-items>
      </v-stepper>
      <v-divider class="my-3"></v-divider>
      <div>
        <div class="text-h6">{{ $t('ecommerce.description') }}</div>
        <div class="text-body-1 my-2">如有疑问请联系管理员.</div>
      </div>
    </v-card>
  </div>
</template>

<script>
  import { mapActions } from 'vuex'

  import Db from './Db'
  import Platform from './Platform'
  import PlatformLogo from './PlatformLogo'
  import Cors from './Cors'
  import Redis from './Redis'
  import Build from './Build'
  import Repo from './Repo'

  export default {
    components: {
      Db,
      Platform,
      PlatformLogo,
      Cors,
      Redis,
      Build,
      Repo,
    },
    data () {
      return {
        e1: 1,
        steps: [
          { key: 'db-step', step: 1, title: '数据库配置' },
          { key: 'platform-step', step: 2, title: '平台初始化' },
          { key: 'logo-step', step: 3, title: '设置Logo' },
          { key: 'cors-step', step: 4, title: '跨域设置' },
          { key: 'redis-step', step: 5, title: 'Redis设置' },
          { key: 'build-step', step: 6, title: '构建设置' },
          { key: 'repo-step', step: 7, title: '镜像仓库' },
          { key: 'reboot-step', step: 8, title: '开始使用' },
        ]
      }
    },
    computed: {
      // ...mapState('app', ['toast']),
    },
    watch: {
      steps (val) {
        if (this.e1 > val) {
          this.e1 = val
        }
      }
    },
    created () {
      let step = 1
      this.steps.map((item, key) => {
        if (this.$route.params.step === item.key) {
          step = item.step
          return
        }
      })
      this.e1 = step
    },
    mounted () {
      console.log(this.$route.params)
    },
    methods: {
      ...mapActions('app', ['showError', 'showSuccess']),
      nextStep (step) {
        let n = this.e1
        this.steps.map((item, key) => {
          if (step === item.key) {
            n = item.step
          }
        })
        if (n === this.steps.length) {
          this.e1 = 1
        } else {
          this.e1 = n
        }
        this.$router.push({ params: { step: step } })
      }
    }
  }
</script>
