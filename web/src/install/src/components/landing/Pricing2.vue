<template>
  <v-sheet>
    <v-container class="pt-4 pt-md-8 pb-10">
      <v-responsive max-width="700" class="mx-auto text-center">
        <div class="text-overline mb-3">Pricing Plans</div>
        <h2 class="text-h3 text-lg-h2">A plan for every price</h2>
        <div class="text-h6 text-lg-h5 mt-4">Lorem ipsum dolor sit amet consectetur adipisicing elit. Natus impedit error labore doloremque fugit.</div>
      </v-responsive>

      <v-row class="mt-10" dense>
        <v-col cols="12" lg="2"></v-col>
        <v-col
          v-for="(plan, i) in plans"
          :key="i"
          cols="12"
          md="6"
          lg="4"
        >
          <v-card
            class="text-body-1"
            outlined
            light
            elevation="2"
            :class="[plan.highlight && 'pricing-highlight']"
          >
            <v-sheet class="pa-4 text-center">
              <div class="text-h6 font-weight-bold text-uppercase mb-6">{{ plan.title }}</div>
              <div class="d-flex justify-center ml-n3">
                <div class="font-weight-regular text-h4 mr-1" style="margin-top: 10px">$</div>
                <div class="text-h1 font-weight-black text-number">{{ plan.price }}</div>
              </div>
              <div class="grey--text justify-end overline">/ month</div>
            </v-sheet>
            <v-divider></v-divider>
            <v-sheet class="pa-4" color="grey lighten-5">
              <div class="my-6 text-h6">
                <div v-for="(feat, j) in plan.features" :key="j" class="d-flex align-center my-1">
                  <v-icon color="success" class="mr-2">mdi-check</v-icon>
                  <span v-if="feat.value" class="font-weight-black mr-1">{{ feat.value }}</span>
                  <span class="font-weight-regular">{{ feat.label }}</span>
                </div>
              </div>
              <v-btn
                x-large
                block
                outlined
                :dark="plan.highlight"
                :class="[plan.highlight && 'primary']"
              >Subscribe</v-btn>
            </v-sheet>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </v-sheet>
</template>

<script>
export default {
  data() {
    return {
      plans: [{
        highlight: false,
        title: 'Basic',
        price: '19',
        features: [{
          value: '5',
          label: 'Accounts'
        }, {
          value: '100GB',
          label: 'Bandwith'
        }, {
          value: '10GB',
          label: 'Disk Space'
        }, {
          value: '',
          label: 'Unlimited Emails'
        }]
      }, {
        highlight: true,
        title: 'Startup',
        price: '29',
        features: [{
          value: '20',
          label: 'Accounts'
        }, {
          value: '1TB',
          label: 'Bandwith'
        }, {
          value: '100GB',
          label: 'Disk Space'
        }, {
          value: '',
          label: 'Unlimited Emails'
        }, {
          value: '',
          label: 'Support 24/7'
        }]
      }]
    }
  }
}
</script>

<style lang="scss" scoped>
@include media("md-and-up") {
  .pricing-highlight {
    z-index: 1;
    transform: scale(1.1);
  }
}
</style>
