<template>
  <div class="w-full">
    <v-card class="pa-12" flat>
      <v-stepper v-model="e1">
        <v-stepper-header>
          <template v-for="st in steps">
            <v-stepper-step :key="st.key" :complete="e1 > st.step" :step="st.step">{{st.title}}</v-stepper-step>
            <v-divider :key="st.step"></v-divider>
          </template>

        </v-stepper-header>

        <v-stepper-items>
          <db :nextStep="nextStep(2)" :v-if="e1===1"/>
          <plat-form :click="nextStep(3)" :v-if="e1===2"/>
          <plat-form :click="nextStep(3)" :v-if="e1===3"/>
          <plat-form :click="nextStep(3)" :v-if="e1===4"/>
          <plat-form :click="nextStep(3)" :v-if="e1===5"/>
          <plat-form :click="nextStep(3)" :v-if="e1===6"/>
          <plat-form :click="nextStep(3)" :v-if="e1===7"/>
        </v-stepper-items>
      </v-stepper>
      <v-divider class="my-3"></v-divider>
      <div>
        <div class="text-h6">{{ $t('ecommerce.description') }}</div>
        <div class="text-body-1 my-2">注意是否是敏感文件.</div>
      </div>
    </v-card>
  </div>
</template>

<script>
  import { mapActions } from 'vuex'

  import Db from './Db'
  import PlatForm from './PlatForm'

  export default {
    components: {
      Db,
      PlatForm,
    },
    data () {
      return {
        e1: 1,
        steps: [
          { key: 'db-step', step: 1, title: '数据库配置' },
          { key: 'plot-form-step', step: 2, title: '平台初始化' },
          { key: 'logo-step', step: 3, title: '设置Logo' },
          { key: 'cors-step', step: 4, title: '跨域设置' },
          { key: 'redis-step', step: 5, title: 'Redis设置' },
          { key: 'build-step', step: 6, title: '构建设置' },
          { key: 'repo-step', step: 7, title: '镜像仓库' },
          { key: 'reboot-step', step: 8, title: '重启' },
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
    mounted () {
      // const path = (window.location.pathname).split('/')
      // let code = ''
      // for (let i in path) {
      //   if (path[i] === '' || path[i] === 's') {
      //     continue
      //   }
      //   code = path[i]
      //   break
      // }
    },
    methods: {
      ...mapActions('app', ['showError', 'showSuccess']),
      nextStep (n) {
        if (n === this.steps.length) {
          this.e1 = 1
        } else {
          this.e1 = n + 1
        }
      }
    }
  }
</script>
