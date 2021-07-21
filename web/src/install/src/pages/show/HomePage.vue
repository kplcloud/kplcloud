<template>
  <div class="w-full">
    <v-card class="pa-3" flat>
      <v-row>
        <v-col cols="12" md="5">
          <div class="d-flex">
            <div class="d-flex align-center flex-grow-1 justify-center">
              <img
                v-if="imageUrl"
                :src="imageUrl"
                class="rounded"
                style="max-width: 100%; max-height: 460px"
              />
            </div>
          </div>
        </v-col>
        <v-col cols="12" md="4">
          <div class="d-flex"></div>
          <div class="font-weight-bold text-h5">{{data.name}}</div>
          <v-divider class="my-2"></v-divider>
          <div class="d-flex align-center text-h6">
            <div>查阅次数: {{data.count}}</div>
          </div>
          <div class="mt-3 text-body-1">
            分享人: {{data.sharer}}
          </div>
          <v-divider class="my-2"></v-divider>
          <div class="font-weight-bold mb-1">关于这个文件</div>
          <ul>
            <li>名称: {{data.name}}</li>
            <li>类型: {{data.suffix}}</li>
            <li>大小:</li>
            <li>是否公开: {{data.public}}</li>
            <li>分享时间: {{data.createdAt | formatDate('YYYY/MM/DD HH:m:ss')}}</li>
            <li>过期时间: {{data.expireTime | formatDate('YYYY/MM/DD HH:m:ss')}}</li>
          </ul>
        </v-col>
        <v-col cols="12" md="3">
          <v-card outlined class="pa-2">
            <div class="d-flex flex-column">
              <!--              <div class="text-body-1 font-weight-bold">{{ $t('ecommerce.shipping') }}</div>-->
              <!--              <div class="text-body-1">免费下载</div>-->

              <div class="text-h6 success--text my-3">限时预览</div>
              <v-select
                :items="[1]"
                :value="1"
                :label="$t('ecommerce.quantity')"
                outlined
                dense
              ></v-select>

              <v-btn color="primary" block large class="mb-2" @click="() => openFile()">预览</v-btn>
            </div>
          </v-card>
        </v-col>
      </v-row>
      <v-divider class="my-3"></v-divider>
      <div>
        <div class="text-h6">{{ $t('ecommerce.description') }}</div>
        <div class="text-body-1 my-2">注意是否是敏感文件.</div>
      </div>
    </v-card>
  </div>
</template>

<script>
  import { mapState, mapActions } from 'vuex'

  import { showInfo } from '@/api/show.service'

  const imageSuffix = [
    'jpg',
    'png',
    'jpeg',
    'gif',
  ]
  export default {
    components: {},
    data () {
      return {
        isLoading: false,
        data: {},
        imageUrl: ''
      }
    },
    computed: {
      // ...mapState('app', ['toast']),
    },
    mounted () {
      const path = (window.location.pathname).split('/')
      let code = ''
      for (let i in path) {
        if (path[i] === '' || path[i] === 's') {
          continue
        }
        code = path[i]
        break
      }
      this.getInfo(code)
    },
    methods: {
      ...mapActions('app', ['showError', 'showSuccess']),
      openFile () {
        window.open(this.data.urls[0].url)
      },
      getInfo (code) {
        this.isLoading = true
        showInfo({ code: code }).then(res => {
          this.data = res.data
          if (imageSuffix.indexOf(res.data.suffix) !== -1) {
            this.imageUrl = res.data.urls[0].url
          } else {
            this.imageUrl = require(`@/assets/images/icons/${res.data.suffix}.png`)
          }
          this.showSuccess('success')
        }).finally(() => {
          this.isLoading = false
        }).catch((err) => {
          this.$router.push("404")
          this.showError({ message: '错误', error: { message: err } })
        })
      }
    }
  }
</script>
