<template>
  <v-card>
    <v-card-title>{{ label }}</v-card-title>
    <v-data-table
      :headers="headers"
      :items="items"
      hide-default-footer
    >
      <template v-slot:item.id="{ item }">
        <div class="font-weight-bold"># <copy-label :text="item.id" /></div>
      </template>

      <template v-slot:item.user="{ item }">
        <div class="d-flex align-center py-1">
          <v-avatar size="40" class="elevation-1 grey lighten-3">
            <v-img :src="item.user.avatar" />
          </v-avatar>
          <div class="ml-1">
            <div class="font-weight-bold">{{ item.user.name }}</div>
            <div class="caption">
              <copy-label :text="item.user.email" />
            </div>
          </div>
        </div>
      </template>

      <template v-slot:item.date="{ item }">
        <div>{{ item.date | formatDate }}</div>
      </template>

      <template v-slot:item.company="{ item }">
        <copy-label :text="item.company" />
      </template>

      <template v-slot:item.amount="{ item }">
        {{ item.amount | formatCurrency }}
      </template>

      <template v-slot:item.status="{ item }">
        <div class="font-weight-bold d-flex align-center text-no-wrap">
          <div v-if="item.status === 'PENDING'" class="warning--text">
            <v-icon small color="warning">mdi-circle-medium</v-icon>
            <span>Pending</span>
          </div>
          <div v-if="item.status === 'PAID'" class="success--text">
            <v-icon small color="success">mdi-circle-medium</v-icon>
            <span>Paid</span>
          </div>
        </div>
      </template>

      <template v-slot:item.action="{ item }">
        <div class="actions">
          <v-btn icon @click="open(item)">
            <v-icon>mdi-open-in-new</v-icon>
          </v-btn>
        </div>
      </template>
    </v-data-table>
  </v-card>
</template>

<script>
import CopyLabel from '../common/CopyLabel'

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
  components: {
    CopyLabel
  },
  props: {
    label: {
      type: String,
      default: ''
    }
  },
  data () {
    return {
      headers: [
        { text: 'Order Id', align: 'start', value: 'id' },
        {
          text: 'User',
          sortable: false,
          value: 'user'
        },
        { text: 'Date', value: 'date' },
        { text: 'Company', value: 'company' },
        { text: 'Amount', value: 'amount' },
        { text: 'Status', value: 'status' },
        { text: '', sortable: false, align: 'right', value: 'action' }
      ],
      items: [
        {
          id: '2837',
          user: {
            name: 'John Simon',
            email: 'johnsimon@blobhill.com',
            avatar: '/images/avatars/avatar1.svg'
          },
          date: '2020-05-10',
          company: 'BlobHill',
          amount: 52877,
          status: 'PAID'
        },
        {
          id: '2838',
          user: {
            name: 'Greg Cool J',
            email: 'cool@caprimooner.com',
            avatar: '/images/avatars/avatar2.svg'
          },
          date: '2020-05-11',
          company: 'Caprimooner',
          amount: 2123,
          status: 'PENDING'
        },
        {
          id: '2839',
          user: {
            name: 'Samantha Bush',
            email: 'bush@catloveisstilllove.com',
            avatar: '/images/avatars/avatar3.svg'
          },
          date: '2020-05-11',
          company: 'CatLovers',
          amount: 12313,
          status: 'PENDING'
        },
        {
          id: '2840',
          user: {
            name: 'Ben Howard',
            email: 'ben@indiecoolers.com',
            avatar: '/images/avatars/avatar4.svg'
          },
          date: '2020-05-12',
          company: 'IndieCoolers',
          amount: 9873,
          status: 'PAID'
        },
        {
          id: '2841',
          user: {
            name: 'Jack Candy',
            email: 'jack@candylooove.com',
            avatar: '/images/avatars/avatar5.svg'
          },
          date: '2020-05-13',
          company: 'CandyLooove',
          amount: 29573,
          status: 'PAID'
        }
      ]
    }
  },
  methods: {
    open(item) { }
  }
}
</script>
