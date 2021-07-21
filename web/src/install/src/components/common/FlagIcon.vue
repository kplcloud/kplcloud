<template>
  <span class="flag-icon" :class="[`flag-icon-${flag}`, { 'flag-round': round }]"></span>
</template>

<script>
import 'flag-icon-css/sass/flag-icon.scss'

/*
|---------------------------------------------------------------------
| Flag Icon Component
|---------------------------------------------------------------------
|
| Show by country code the respective flag for that country
| https://github.com/lipis/flag-icon-css
|
*/
export default {
  props: {
    // Country ISO
    flag: {
      type: String,
      default: 'us'
    },
    // Circle flag or square
    round: {
      type: Boolean,
      default: false
    }
  }
}
</script>

<style lang="scss" scoped>
.flag-icon {
  height: 22px;
  width: 22px;

  &.flag-round {
    background-size: cover;
    border-radius: 100%;
    height: 26px;
    width: 26px;
  }
}
</style>
