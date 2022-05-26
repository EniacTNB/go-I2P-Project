<template>
  <!-- <div class="dashboard-container"> -->
  <!-- <div class="dashboard-text"> -->

  <!-- <div style="width: 540px; height: 422px; overflow: hidden; border: 0px"> -->
  <div style="margin: -25px 0px 0px 0px">
    <iframe
      src="http://207.246.81.104:81/ipmap"
      height="1000"
      width="100%"
    ></iframe>
  </div>
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
