<template>
  <v-app>
    <v-app-bar
      app
      color="primary"
      dark
    >
      <div class="d-flex align-center">

        <h1>Skill-Share-Club</h1>
      </div>

      <v-spacer></v-spacer>

    </v-app-bar>

    <v-main class="main">
      <Calendar class="calendar-main"/>
    </v-main>
  </v-app>
</template>

<script lang="ts">
import Vue from 'vue';
import Calendar from './components/Calendar.vue';

export default Vue.extend({
  name: 'App',

  components: {
    Calendar,
  },

  data: () => ({
    //
  }),
});
</script>

<style scoped>
.calendar-main {
  height: 90vh;
}
</style>
