<template>
  <!-- <div class="dashboard-container"> -->
  <!-- <div class="dashboard-text"> -->
  <div style="margin: 0px 0px 0px -183px">
    <iframe
      src="http://208.167.248.208:5601/app/kibana#/visualize?_g=(refreshInterval:(pause:!t,value:0),time:(from:'2022-01-05T11:48:30.232Z',mode:absolute,to:'2022-01-12T21:51:08.799Z'))"
      height="800"
      width="100%"
    ></iframe>
  </div>
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
