<template>
  <v-container id="pricing" class="pt-4 pt-md-8 pb-10">
    <v-responsive max-width="1200" class="mx-auto text-center">
      <h2 class="text-h3 mb-2">Pricing</h2>
      <div class="text-h6 text-lg-h5">Lorem ipsum dolor sit amet consectetur adipisicing elit. Natus impedit error labore doloremque fugit.</div>
    </v-responsive>

    <v-row class="mt-10">
      <v-col
        v-for="(plan, i) in plans"
        :key="i"
        cols="12"
        lg="4"
      >
        <v-card
          class="text-body-1 pa-4"
          outlined
          elevation="2"
        >
          <div class="d-flex justify-space-between">
            <div class="mr-2">
              <div class="d-flex align-center">
                <div class="text-h4 font-weight-black">{{ plan.title }}</div>
                <div v-if="plan.featured" class="ml-2">
                  <v-chip small color="primary" class="font-weight-black">Popular</v-chip>
                </div>
              </div>
              <div class="mt-1">Lorem ipsum dolor sit amet consectetur adipisicing elit.</div>
            </div>
            <div class="text-right">
              <div class="d-flex align-center">
                <div class="text-h5 font-weight-light mr-1">$</div>
                <div class="text-h3 text-number font-weight-black">{{ plan.price }}</div>
              </div>
              <div class="justify-end overline">/month</div>
            </div>
          </div>

          <v-divider class="my-4"></v-divider>

          <div class="text-h6">
            <div v-for="(feat, j) in plan.features" :key="j" class="d-flex align-center justify-space-between my-1">
              <div>
                <span v-if="feat.value" class="font-weight-black mr-1">{{ feat.value }}</span>
                <span class="text-truncate font-weight-regular">{{ feat.label }}</span>
              </div>
              <v-icon color="primary">mdi-check</v-icon>
            </div>
          </div>
          <v-divider class="my-4"></v-divider>
          <v-btn
            x-large
            block
            :outlined="!plan.featured"
            :color="plan.featured ? 'primary' : ''"
            to="/auth/signup"
          >
            Subscribe
          </v-btn>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
export default {
  data() {
    return {
      plans: [{
        title: 'Basic',
        price: '29',
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
        title: 'Startup',
        price: '49',
        featured: true,
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
        }]
      }, {
        title: 'Enterprise',
        price: '149',
        features: [{
          value: '100',
          label: 'Accounts'
        }, {
          value: '10TB',
          label: 'Bandwith'
        }, {
          value: '10TB',
          label: 'Disk Space'
        }, {
          value: '',
          label: 'Unlimited Emails'
        }]
      }]
    }
  }
}
</script>
