<template>
  <v-card class="d-flex flex-column flex-grow-1">

    <!-- loading spinner -->
    <div v-if="loading" class="d-flex flex-grow-1 align-center justify-center">
      <v-progress-circular indeterminate color="primary"></v-progress-circular>
    </div>

    <!-- information -->
    <div v-else class="d-flex flex-column flex-grow-1">
      <v-card-title>
        <div>{{ label }}</div>
        <v-spacer></v-spacer>
        <div>
          <v-select
            v-model="selectedInterval"
            solo
            dense
            hide-details
            hide-selected
            :items="intervals"
          ></v-select>
        </div>
      </v-card-title>

      <div class="chart-wrap">
        <apexchart
          type="donut"
          width="85%"
          :options="chartOptions"
          :series="series"
        ></apexchart>
      </div>
    </div>
  </v-card>
</template>

<script>
import VueApexCharts from 'vue-apexcharts'

/*
|---------------------------------------------------------------------
| DEMO Dashboard Card Component
|---------------------------------------------------------------------
|
| Demo card component to be used to gather some ideas on how to build
| your own dashboard component
|
*/
export default {
  props: {
    series: {
      type: Array,
      default: () => ([])
    },
    label: {
      type: String,
      default: ''
    },
    color: {
      type: String,
      default: '#333333'
    },
    value: {
      type: Number,
      default: 0
    },
    percentage: {
      type: Number,
      default: 0
    },
    percentageLabel: {
      type: String,
      default: 'vs. last week'
    },
    options: {
      type: Object,
      default: () => ({})
    },
    loading: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      selectedInterval: 'Last 7 days',
      intervals: [
        'Last 7 days',
        'Last 28 days',
        'Last month',
        'Last year'
      ]
    }
  },
  computed: {
    chartOptions() {
      return {
        chart: {
          type: 'donut',
          animations: {
            speed: 400
          },
          background: 'transparent'
        },
        stroke: {
          show: true,
          colors: [this.$vuetify.theme.isDark ? '#333' : '#fff'],
          width: 1,
          dashArray: 0
        },
        plotOptions: {
          pie: {
            expandOnClick: false,
            donut: {
              size: '74%'
            }
          }
        },
        theme: {
          mode: this.$vuetify.theme.isDark ? 'dark' : 'light'
        },
        labels: ['Referrals', 'Organic Search', 'Social Media', 'Others'],
        dataLabels: {
          enabled: false
        },
        colors: ['#2364aa', '#3da5d9', '#73bfb8', '#fec601', '#ea7317'],
        fill: {
          colors: ['#2364aa', '#3da5d9', '#73bfb8', '#fec601', '#ea7317']
        },
        legend: {
          offsetY: 40,
          fontSize: '13px',
          fontFamily: 'Quicksand',
          fontWeight: 700
        },
        responsive: [{
          breakpoint: 480,
          options: {
            chart: {
              width: 200
            },
            legend: {
              offsetY: 0,
              position: 'bottom'
            }
          }
        }],
        ...this.options
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.chart-wrap {
  max-width: 560px;
  max-height: 280px;
}
</style>
