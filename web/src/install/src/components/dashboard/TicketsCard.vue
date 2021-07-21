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
        <div class="font-weight-bold">
          <copy-label :text="item.user.email" />
        </div>
      </template>

      <template v-slot:item.date="{ item }">
        <div>{{ item.date | formatDate }}</div>
      </template>

      <template v-slot:item.priority="{ item }">
        <v-chip
          label
          small
          class="font-weight-bold"
          :class="{
            'error': item.priority === 'High'
          }"
        >{{ item.priority }}</v-chip>
      </template>

      <template v-slot:item.status="{ item }">
        <div class="font-weight-bold d-flex align-center">
          <div v-if="item.status === 'CLOSED'" class="secondary--text">
            <v-icon small color="secondary">mdi-circle-medium</v-icon>
            <span>Closed</span>
          </div>
          <div v-if="item.status === 'OPEN'" class="success--text">
            <v-icon small color="success">mdi-circle-medium</v-icon>
            <span>Open</span>
          </div>
        </div>
      </template>

      <template v-slot:item.action="{ item }">
        <div class="actions">
          <v-btn small @click="open(item)">
            Open Ticket
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
        { text: 'Ticket Id', align: 'start', value: 'id' },
        {
          text: 'User',
          sortable: false,
          value: 'user'
        },
        { text: 'Priority', value: 'priority' },
        { text: 'Status', value: 'status' },
        { text: 'Create Date', value: 'date' },
        { text: '', sortable: false, align: 'right', value: 'action' }
      ],
      items: [
        {
          id: '423',
          user: {
            name: 'John Simon',
            email: 'johnsimon@blobhill.com',
            avatar: '/images/avatars/avatar1.svg'
          },
          date: '2020-05-10',
          priority: 'Low',
          status: 'OPEN'
        },
        {
          id: '424',
          user: {
            name: 'Greg Cool J',
            email: 'cool@caprimooner.com',
            avatar: '/images/avatars/avatar2.svg'
          },
          date: '2020-05-11',
          priority: 'High',
          status: 'CLOSED'
        },
        {
          id: '425',
          user: {
            name: 'Samantha Bush',
            email: 'bush@catloveisstilllove.com',
            avatar: '/images/avatars/avatar3.svg'
          },
          date: '2020-05-11',
          priority: 'Low',
          status: 'CLOSED'
        },
        {
          id: '426',
          user: {
            name: 'Ben Howard',
            email: 'ben@indiecoolers.com',
            avatar: '/images/avatars/avatar4.svg'
          },
          date: '2020-05-12',
          priority: 'Low',
          status: 'OPEN'
        },
        {
          id: '427',
          user: {
            name: 'Jack Candy',
            email: 'jack@candylooove.com',
            avatar: '/images/avatars/avatar5.svg'
          },
          date: '2020-05-13',
          priority: 'High',
          status: 'OPEN'
        }
      ]
    }
  },
  methods: {
    open(item) { }
  }
}
</script>
