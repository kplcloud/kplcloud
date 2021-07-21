<template>
  <v-sheet color="transparent">
    <v-container class="py-4 py-lg-10">
      <div class="text-center">
        <div class="text-overline mb-3">Support</div>
        <h2 class="text-h3 text-lg-h2">Frequently Asked Questions</h2>
        <v-responsive max-width="1200" class="mx-auto">
          <div class="text-h6 text-lg-h5 mt-4">Lorem ipsum dolor sit amet consectetur adipisicing elit. Natus impedit error labore doloremque fugit! Dolor fugit molestiae vero quos quisquam nobis, eos debitis magni omnis ea incidunt amet voluptate dignissimos!</div>
        </v-responsive>
      </div>
      <v-row class="mt-6">
        <v-col
          v-for="(item, i) in faq"
          :key="i"
          cols="12"
          lg="6"
        >
          <v-card class="pa-3">
            <div class="text-h5">
              {{ item.question }}
            </div>
            <div class="mt-2">
              {{ item.answer }}
            </div>
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
      faq: [{
        question: 'What happened to the developer who couldn\'t access git?',
        answer: 'He now has commitement issues.'
      }, {
        question: 'What do computers and air conditioners have in common?',
        answer: 'They both become useless when you open window.'
      }, {
        question: 'How do you tell an introvert developer from and extrovert developer?',
        answer: 'An extrovert developer looks at your shoes when he talks to you.'
      }, {
        question: 'What is the programmers\'s favourite hangout place?',
        answer: 'Foo Bar'
      }, {
        question: 'Why do Python developers need glasses?',
        answer: 'Because they don\'t C#.'
      }, {
        question: 'What\'s the object-oriented way to become wealthy?',
        answer: 'Inheritance'
      }]
    }
  }
}
</script>
