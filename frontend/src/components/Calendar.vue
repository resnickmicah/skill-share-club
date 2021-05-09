<template>
  <v-row>
    <v-col>
      <!-- this height of 400 does not work for me. it either needs a border or 100% height to fill the parent. -->
      <!-- <v-sheet height="400"> -->
      <v-sheet class="main_sheet">
        <v-calendar
          ref="calendar"
          :now="today"
          :value="today"
          :events="events"
          color="primary"
          type="month"
          @click:date="heck"
        ></v-calendar>
      </v-sheet>
    </v-col>
  </v-row>
  
</template>

<script>
export default {
  methods: {
    heck: function (event) {
      console.log(event)
      let newEvent = prompt(`Do you want to add an event to ${event.date}?`)
      console.log(this.events)
      this.events.push({
        name: newEvent,
        start: event.date,
        end: event.date
      })
    }
  },
  data: () => ({
    today: "2021-05-04",
    events: [
      {
        name: "pair programming",
        start: "2021-05-02 09:00",
        end: "2021-05-02 11:00",
      },
      {
        name: "Weekly Meeting",
        start: "2021-05-03 09:00",
        end: "2021-05-03 10:00",
      },
      {
        name: `Thomas' Birthday`,
        start: "2021-05-07",
      },
      {
        name: "Mash Potatoes",
        start: "2021-05-08 12:30",
        end: "2021-05-08 15:30",
      },
    ],
  }),
  mounted() {
    this.$refs.calendar.scrollToTime("08:00");
  },
};
</script>

<style scoped>
.my-event {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  border-radius: 2px;
  background-color: #1867c0;
  color: #ffffff;
  border: 1px solid #1867c0;
  font-size: 12px;
  padding: 3px;
  cursor: pointer;
  margin-bottom: 1px;
  left: 4px;
  margin-right: 8px;
  position: relative;
}
.my-event.with-time {
  position: absolute;
  right: 4px;
  margin-right: 0px;
}
.main_sheet {
  height: 93vh;
  border: 2px solid black;
}
</style>
