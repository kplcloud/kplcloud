<template>
  <span>
    <span v-if="value === 0">
      {{ value }}%
    </span>
    <span v-else-if="value > 0" class="success--text">
      <v-icon small color="success">mdi-arrow-top-right</v-icon> {{ value }}%
    </span>
    <span v-else class="error--text">
      <v-icon small color="error">mdi-arrow-bottom-right</v-icon> {{ Math.abs(value) }}%
    </span>
  </span>
</template>

<script>
/*
|---------------------------------------------------------------------
| Trend Percent Component
|---------------------------------------------------------------------
|
| Shows and up or down arrow with green or red color presenting the
| trend of the given value
|
*/
export default {
  props: {
    value: {
      type: Number,
      default: 0
    }
  }
}
</script>
