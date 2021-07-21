<template>
  <tasks-page class="todo-card"></tasks-page>
</template>

<script>
import TasksPage from '../../apps/todo/pages/TasksPage'

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
    TasksPage
  }
}
</script>

<style lang="scss" scoped>
.todo-card {
  height: 100%;
  max-height: 360px;
  overflow-y: scroll;
}
</style>
