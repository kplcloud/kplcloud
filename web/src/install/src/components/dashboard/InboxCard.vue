<template>
  <inbox-page class="email-card"></inbox-page>
</template>

<script>
import InboxPage from '../../apps/email/pages/InboxPage'

export default {
  components: {
    InboxPage
  }
}
</script>

<style lang="scss" scoped>
.email-card {
  max-height: 340px;
  overflow-y: scroll;
}
</style>
