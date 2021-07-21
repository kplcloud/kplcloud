import Vue from 'vue'
import VueECharts from 'vue-echarts'

/**
 * Vue ECharts
 * https://github.com/ecomfe/vue-echarts
 *
 */
import 'echarts/lib/component/tooltip'
import 'echarts/lib/component/legend'
import 'echarts/lib/component/polar'
import 'echarts/lib/chart/bar'
import 'echarts/lib/chart/line'
import 'echarts/lib/chart/pie'
import 'echarts/lib/chart/radar'

Vue.component('e-charts', VueECharts)
