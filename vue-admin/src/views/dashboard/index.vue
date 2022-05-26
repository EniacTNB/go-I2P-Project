<template>
  <!-- <div class="dashboard-container"> -->
  <!-- <div class="dashboard-text"> -->
  <iframe
    src="http://208.167.248.208:5601/app/kibana#/dashboard/5bdbbc10-64b8-11ec-81fe-7b7afb4e36b4?embed=true&_g=(refreshInterval%3A(pause%3A!t%2Cvalue%3A0)%2Ctime%3A(from%3A'2022-01-05T11%3A48%3A30.232Z'%2Cmode%3Aabsolute%2Cto%3A'2022-01-12T21%3A51%3A08.799Z'))"
    height="1200"
    width="100%"
  ></iframe>
  <!-- </div> -->
  <!-- <div class="dashboard-text">name: {{ name }}</div> -->
  <!-- </div> -->
</template>

<script>
import { mapGetters } from "vuex";

export default {
  name: "Dashboard",
  computed: {
    ...mapGetters(["name"]),
  },
};
</script>

// <style lang="scss" scoped>
// // .dashboard {
// //   &-container {
// //     height: 100%;
// //   }
// //   &-text {
// //     font-size: 30px;
// //     line-height: 46px;
// //   }
// // }
// //
//
//
//
</style>
