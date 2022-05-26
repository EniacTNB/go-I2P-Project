<template>
  <div class="app-container">
    <el-table
      v-loading="listLoading"
      :data="list"
      element-loading-text="Loading"
      border
      fit
      highlight-current-row
    >
      <el-table-column align="center" label="Index" width="95">
        <template slot-scope="scope">
          {{ scope.$index + 1 }}
        </template>
      </el-table-column>
      <el-table-column label="ProbeName" width="150" align="center">
        <template slot-scope="scope">
          {{ scope.row.name }}
        </template>
      </el-table-column>
      <el-table-column label="IP" width="150" align="center">
        <template slot-scope="scope">
          <span>{{ scope.row.IP }}</span>
        </template>
      </el-table-column>
      
      <el-table-column label="Address" width="150" align="center">
        <template slot-scope="scope">
          {{ scope.row.Address }}
        </template>
      </el-table-column>
      <el-table-column label="Number" width="150" align="center">
        <template slot-scope="scope">
          {{ scope.row.routerinfo }}
        </template>
      </el-table-column>

      <el-table-column
        align="center"
        prop="created_at"
        label="Display_time"
        width="200"
      >
        <template slot-scope="scope">
          <i class="el-icon-time" />
          <span>{{ getDisplayTime(scope.row.time) }}</span>
        </template>
      </el-table-column>
      <el-table-column
        class-name="status-col"
        label="Status"
        width="110"
        align="center"
      >
        <template slot-scope="scope">
          <el-tag :type="'danger'" @click="checkRouterInfo(scope.row)">{{
            "查看详情"
          }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column
        class-name="status-col"
        label="Manage"
        width="110"
        align="center"
      >
        <template slot-scope="scope">
          <el-tag :type=indexList[scope.$index]  @click="changeProbeStatus(scope.$index)">{{
            valueList[scope.$index]
          }}</el-tag>
        </template>
      </el-table-column>
      
      
    </el-table>
  </div>
</template>

<script>
import { getList } from "@/api/table";
import Pagination from "@/layout/components/PageIndex";

export default {
  filters: {
    statusFilter(Status) {
      const statusMap = {
        ts: "success",
        draft: "gray",
        check: "danger",
      };
      return statusMap[Status];
    },
  },
  // components: { Pagination },
  computed: {
    // 最大页数
    pageMax() {
      return Math.ceil(this.total / this.limit);
    },
  },
  data() {
    return {
      // total: 200,
      // page: 1, // 当前页码
      size: 20, // 每页数量
      list: null,
      listLoading: true,
      dataList: [],
      requestForm: {
        page: this.page,
        size: this.size,
      },
      form: {
        region: "",
      },
      pageList: [], // 页码列表
    //   index:"success",
      indexList:["danger","danger","danger","success"],
      valueList:["stop","stop","stop","start"]
    };
  },
  created() {
    this.fetchData(this.size, this.page);
  },
  onLoad() {
    this.initData();
    this.getData(this.size, this.page);
  },

  methods: {
    checkRouterInfo(value) {
      // alert("${this.$route.fullPath}");
      localStorage.setItem("searchParam", JSON.stringify(value));
      this.$router.push(`/example/tree?id=`);
    },
    buttonType(index){
// 选中的button更改颜色，其余的都更改为灰色
        return index
      },
    changeProbeStatus(index){
        console.log("index is ")
        console.log(index)
        if(this.indexList[index]=='success'){
            console.log("this is success")
            // this.buttonType("danger")
            this.indexList[index] ="danger"
            this.valueList[index] = "stop"
        }else{
            // this.buttonType("success")
            this.indexList[index] ="success"
            this.valueList[index] = "start"
        }
    },
    initData() {
      this.pageList = []; // 清空页码
      var i = 1;
      do {
        this.pageList.push(i);
        i++;
      } while (i <= this.pageMax);
      this.pageList.length > 5 && // 最多显示5页
        (this.pageList = this.pageList.slice(0, 5));
    },
    // 子组件事件回调：分页
    pageChange(pageCurrent) {
      this.$emit("page-change", pageCurrent);
    },

    getData(size, page) {
    //   getList(this.requestForm).then((response) => {
    //     this.dataList = response.data;
    //   });
      // this.dataList = this.list;
      console.log("dddddataList");
      console.log(this.dataList);
    },
    pageChange(page) {
      this.page = page;
      this.getData();
    },
    fetchData(size, page) {
      this.listLoading = false;
    //   getList(this.requestForm).then((response) => {
    //     this.list = response.data;
    //     this.listLoading = false;
    //     console.log("list is...");
    //     console.log(this.list);
    //     // return response;
    //   });
    this.list = [{'index':1,'name':'SinaporeProbe','IP':'207.148.64.55','Address':'Singapore','routerinfo':15883,'time':1643381397000},
    {'index':1,'name':'AmericanProbe','IP':'45.63.21.151','Address':'New Jersey','routerinfo':20967,'time':1643381486000},
    {'index':1,'name':'UKProbe','IP':'165.225.94.253','Address':'London','routerinfo':22494,'time':1643381973000},
    {'index':1,'name':'LocalProbe','IP':'localhost','Address':'Beijing','routerinfo':305772,'time':1631485399200},]
    },
    getDisplayTime(ts) {
      var date = new Date(ts);
      var Year = date.getFullYear() + "-";
      var Month =
        (date.getMonth() + 1 < 10
          ? "0" + (date.getMonth() + 1)
          : date.getMonth() + 1) + "-";
      var Day = date.getDate() + " ";
      var hour = date.getHours() + ":";
      var minutes = date.getMinutes() + ":";
      var secodes = date.getSeconds();
      var res = Year + Month + Day + hour + minutes + secodes;
      return res;
    },
  },
  watch: {
    // 监听页码变化 -> 页码列表更新
    page(val) {
      if (val <= 3) {
        this.pageList = [];
        var i = 1;
        do {
          this.pageList.push(i);
          i++;
        } while (i <= this.pageMax);
        this.pageList.length > 5 && // 最多显示5页
          (this.pageList = this.pageList.slice(0, 5));
      } else if (val === this.pageMax) {
        this.pageList = [val - 2, val - 1, val];
      } else if (val === this.pageMax - 1) {
        this.pageList = [val - 2, val - 1, val, val + 1];
      } else {
        this.pageList = [val - 2, val - 1, val, val + 1, val + 2];
      }
    },
    // 监听页码变化 -> 总数更新
    total(val) {
      this.initData();
    },
  },
};
</script>
<style rel="stylesheet/scss" lang="scss" scoped>
.page {
  padding: 10px;
  background-color: #fff;

  &-block {
    display: inline-block;
    width: 30px;
    height: 28px;
    padding: 4px 8px;
    font-size: 0.8em;
    line-height: 18px;
    border: 1px solid #ddd;
    border-right: none;
    box-sizing: border-box;

    &:first-child {
      width: 34px;
      border-top-left-radius: 4px;
      border-bottom-left-radius: 4px;
    }

    &:last-child {
      width: 34px;
      border-right: 1px solid #ddd;
      border-top-right-radius: 4px;
      border-bottom-right-radius: 4px;
    }
  }
}
</style>