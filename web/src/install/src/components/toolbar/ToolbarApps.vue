<template>
  <v-menu offset-y left transition="slide-y-transition">
    <template v-slot:activator="{ on }">
      <v-btn icon v-on="on">
        <v-icon>mdi-view-grid-outline</v-icon>
      </v-btn>
    </template>

    <v-card class="d-flex flex-row flex-wrap" style="max-width: 280px">
      <div
        v-for="app in apps"
        :key="app.link"
        class="app-tile pa-3 text-center"
        style="flex: 0 50%"
        @click="navigateTo(app.link)"
      >
        <v-icon color="primary">{{ app.icon }}</v-icon>
        <div class="font-weight-bold mt-1">{{ app.key ? $t(app.key) : app.text }}</div>
        <div class="caption">{{ app.subtitleKey ? $t(app.subtitleKey) : app.subtitle }}</div>
      </div>
    </v-card>
  </v-menu>
</template>

<script>
import config from '../../configs'
/*
|---------------------------------------------------------------------
| Toolbar Apps Component
|---------------------------------------------------------------------
|
| Quickmenu for applications in the toolbar
|
*/
export default {
  data() {
    return {
      apps: config.toolbar.apps
    }
  },
  methods: {
    navigateTo(path) {
      if (this.$route.path !== path) this.$router.push(path)
    }
  }
}
</script>

<style lang="scss" scoped>
.app-tile {
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
  cursor: pointer;
  border-radius: 6px;
  background-color: var(--v-background-base);
  transition: transform 0.2s;

  &:hover {
    transform: scale(1.1);
  }
}
</style>
