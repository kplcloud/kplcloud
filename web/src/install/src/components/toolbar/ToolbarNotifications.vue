<template>
  <v-menu offset-y left transition="slide-y-transition">
    <template v-slot:activator="{ on }">
      <v-badge
        bordered
        content="6"
        offset-x="22"
        offset-y="22"
      >
        <v-btn icon v-on="on">
          <v-icon>mdi-bell-outline</v-icon>
        </v-btn>
      </v-badge>
    </template>

    <!-- dropdown card -->
    <v-card>
      <v-list three-line dense max-width="400">
        <v-subheader class="pa-2 font-weight-bold">Notifications</v-subheader>
        <div v-for="(item, index) in items" :key="index">
          <v-divider v-if="index > 0 && index < items.length" inset></v-divider>

          <v-list-item @click="">
            <v-list-item-avatar size="32" :color="item.color">
              <v-icon dark small>{{ item.icon }}</v-icon>
            </v-list-item-avatar>

            <v-list-item-content>
              <v-list-item-title v-text="item.title"></v-list-item-title>
              <v-list-item-subtitle class="caption" v-text="item.subtitle"></v-list-item-subtitle>
            </v-list-item-content>
            <v-list-item-action class="align-self-center">
              <v-list-item-action-text v-text="item.time"></v-list-item-action-text>
            </v-list-item-action>
          </v-list-item>
        </div>
      </v-list>

      <div class="text-center py-2">
        <v-btn small>See all</v-btn>
      </div>
    </v-card>
  </v-menu>
</template>

<script>
/*
|---------------------------------------------------------------------
| Toolbar Notifications Component
|---------------------------------------------------------------------
|
| Quickmenu to check out notifications
|
*/
export default {
  data() {
    return {
      items: [
        {
          title: 'Brunch this weekend?',
          color: 'primary',
          icon: 'mdi-account-circle',
          subtitle: 'Lorem ipsum dolor sit amet consectetur, adipisicing elit. Sint, repudiandae?',
          time: '3 min'
        },
        {
          title: 'Summer BBQ',
          color: 'success',
          icon: 'mdi-email-outline',
          subtitle: 'Lorem ipsum dolor sit amet consectetur, adipisicing elit. Sint, repudiandae?',
          time: '3 min'
        },
        {
          title: 'Oui oui',
          color: 'teal lighten-1',
          icon: 'mdi-airplane-landing',
          subtitle: 'Lorem ipsum dolor sit amet consectetur, adipisicing elit. Sint, repudiandae?',
          time: '4 min'
        },
        {
          title: 'Disk capacity is at maximum',
          color: 'teal accent-3',
          icon: 'mdi-server',
          subtitle: 'Lorem ipsum dolor sit amet consectetur, adipisicing elit. Sint, repudiandae?',
          time: '3 hr'
        },
        {
          title: 'Recipe to try',
          color: 'blue-grey lighten-2',
          icon: 'mdi-noodles',
          subtitle: 'Lorem ipsum dolor sit amet consectetur, adipisicing elit. Sint, repudiandae?',
          time: '8 hr'
        }
      ]
    }
  }
}
</script>
