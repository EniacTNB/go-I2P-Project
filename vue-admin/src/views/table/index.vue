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
      <el-table-column label="Pubkey" width="400" align="center">
        <template slot-scope="scope">
          {{ scope.row.pubkey }}
        </template>
      </el-table-column>
      <el-table-column label="Signkey" width="400" align="center">
        <template slot-scope="scope">
          <span>{{ scope.row.signkey }}</span>
        </template>
      </el-table-column>
      <!-- <el-table-column label="published" width="400" align="center">
        <template slot-scope="scope">
          {{ scope.row.published }}
        </template>
      </el-table-column> -->
      <el-table-column label="signature" width="400" align="center">
        <template slot-scope="scope">
          {{ scope.row.signature }}
        </template>
      </el-table-column>

      <el-table-column
        align="center"
        prop="created_at"
        label="Display_time"
        width="300"
      >
        <template slot-scope="scope">
          <i class="el-icon-time" />
          <span>{{ getDisplayTime(scope.row.published) }}</span>
        </template>
      </el-table-column>
      <el-table-column
        class-name="status-col"
        label="Status"
        width="110"
        align="center"
      >
        <template slot-scope="scope">
          <el-tag
            :type="'danger'"
            @click="checkRouterInfo(scope.row)"
            >{{ "查看详情" }}</el-tag
          >
        </template>
      </el-table-column>
    </el-table>
    <div class="page">
      <span class="page-block center">{{ "«" }}</span>
      <span v-if="pageList[2] - 2 > 1" class="page-block center">...</span>
      <span
        v-for="pageNum in pageList"
        :key="pageNum"
        class="page-block center"
        :style="{ color: pageNum === page ? '#80bd01' : '#778087' }"
        >{{ pageNum }}</span
      >
      <span v-if="pageMax - pageList[2] > 2" class="page-block center"
        >...</span
      >
      <span class="page-block center">{{ "»" }}</span>
    </div>
    <!-- <pagination
      :total="dataList.length"
      :limit="size"
      :page="page"
      @page-change="pageChange"
    ></pagination> -->
  </div>
</template>

<script>
import { getList } from "@/api/table";
import Pagination from "@/layout/components/PageIndex";

export default {
  props: {
    // 内容总数
    total: {
      type: Number,
      default: 100,
    },
    // 每页数量
    limit: {
      type: Number,
      default: 10,
    },
    // 当前页码
    page: {
      type: Number,
      default: 1,
    },
  },
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
      pageList: [], // 页码列表
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
      getList(this.requestForm).then((response) => {
        this.dataList = response.data;
      });
      // this.dataList = this.list;
      console.log("dddddataList");
      console.log(this.dataList);
    },
    pageChange(page) {
      this.page = page;
      this.getData();
    },
    fetchData(size, page) {
      this.listLoading = true;
      getList(this.requestForm).then((response) => {
        this.list = response.data;
        this.listLoading = false;
        console.log("list is...");
        console.log(this.list);
        // return response;
      });
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