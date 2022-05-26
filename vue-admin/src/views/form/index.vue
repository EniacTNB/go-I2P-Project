<template>
  <div class="app-container">
    <div class="row">
      <div class="col-md-12">
        <!-- <h3>Lua在线运行工具</h3> -->
        <div class="alert alert-warning alert-dismissable">
          <h4>注意</h4>
          每次运行时，Lua虚拟机并不会刷新，前一次的变量都会被保存。如果需要彻底重新运行，请点击刷新按钮：
          <button
            class="btn btn-primary btn-outline-primary"
            onclick="newLuaState()"
          >
            刷新虚拟机</button
          ><br />
          <span id="saveText"
            >当前代码会自动保存在浏览器缓存，下次打开自动加载。</span
          ><br />
          Lua运行环境为Lua 5.3。集成moonscript，可直接require("moonscript")
        </div>
      </div>
    </div>
    <el-form ref="form" :model="form" label-width="120px">
      <el-form-item label="FileName">
        <el-input v-model="form.name" id="filenameinputing" />
      </el-form-item>
      <el-form-item label="I2P Point">
        <el-select
          v-model="form.region"
          placeholder="please select exper time"
          id="timeselector"
        >
          <el-option label="10min" value="10" />
          <el-option label="20min" value="20" />

          <el-option label="30min" value="30" />
          <el-option label="60min" value="60" />
          <el-option label="100min" value="100" />
          <el-option label="120min" value="120" />
        </el-select>
      </el-form-item>
      <el-form-item label="Run Exper At">
        <el-col :span="11">
          <el-date-picker
            v-model="form.date1"
            type="date"
            placeholder="Pick a date"
            style="width: 100%"
            id="date1pocker"
          />
        </el-col>
        <el-col :span="2" class="line">-</el-col>
        <el-col :span="11">
          <el-time-picker
            v-model="form.date2"
            type="fixed-time"
            placeholder="Pick a time"
            style="width: 100%"
            id="date2pocker"
          />
        </el-col>
      </el-form-item>
      <el-form-item label="public">
        <el-switch v-model="form.delivery" id="publicswitcher" />
      </el-form-item>
      <!-- <el-form-item label="experiment type">
        <el-checkbox-group v-model="form.type">
          <el-checkbox label="节点保留策略" name="type" />
          <el-checkbox label="节点选择策略" name="type" />
          <el-checkbox label="节点发布策略" name="type" />
          <el-checkbox label="节点获取策略" name="type" />
        </el-checkbox-group>
      </el-form-item> -->
      <el-form-item label="Exper Type">
        <el-radio-group v-model="form.resource" id="expertyper">
          <el-radio label="节点保留策略" />
          <el-radio label="节点选择策略" />
          <el-radio label="节点发布策略" />
          <el-radio label="节点获取策略" />
        </el-radio-group>
      </el-form-item>
      <el-form-item label="Description">
        <el-input v-model="form.desc" type="textarea" id="descinput" />
      </el-form-item>
      <el-form-item>
        <el-button  type="primary" @click="onSubmit" id="buttonsubmit"
          >Save</el-button
        >
        <el-button @click="onCancel" id="buttoncancel">Clear</el-button>
      </el-form-item>
    </el-form>

    <!-- <div class="row"> -->
    <el-col :span="17">
      <!-- <div class="col-md-6"> -->
      <h3>code:</h3>
      <codemirror
        ref="mycode"
        :value="curCode"
        :options="cmOptions"
        class="code"
        id="mycodes"
        @input="onCmCodeChange"
        style="position: relative; width: 100%; height: 370px"
      >
      </codemirror>
    </el-col>
    <el-col :span="7">
      <!-- <div class="col-md-6"> -->
      <h3>result:</h3>
      <codemirror
        ref="mycode1"
        :value="fileres"
        :options="cmOptions"
        class="code"
        id="mycodes1"
        readonly="true"
        style="position: relative; width: 100%; height: 370px"
      >
        虚拟机初始化ing...</codemirror
      >
      <!-- </div> -->
    </el-col>
    <button
      class="btn btn-primary btn-block btn-outline-primary"
      @click="editExper()"
      id="buttonedit"
    >
      Edit
    </button>
    <button
      class="btn btn-primary btn-block btn-outline-primary"
      @click="getResult()"
      id="buttonrun"
    >
      Run
    </button>
    <!-- </div> -->
  </div>
</template>

<script defer=true>
import { codemirror } from "vue-codemirror";
import "codemirror/theme/ambiance.css"; // 这里引入的是主题样式，根据设置的theme的主题引入，一定要引入！！
// import "codemirror/theme/default.css"; // 这里引入的是主题样式，根据设置的theme的主题引入，一定要引入！！
require("codemirror/mode/lua/lua"); // 这里引入的模式的js，根据设置的mode引入，一定要引入！！
import "codemirror/addon/edit/matchbrackets";
import "codemirror/addon/selection/active-line";
// 括号、引号编辑和删除时成对出现
import "codemirror/addon/edit/closebrackets";
// 折叠代码要用到一些玩意
import "codemirror/addon/fold/foldgutter.css";
import "codemirror/addon/fold/foldgutter";
import "codemirror/addon/fold/xml-fold";
import "codemirror/addon/fold/foldcode";
import "codemirror/addon/fold/brace-fold";
import "codemirror/addon/fold/indent-fold.js";
import "codemirror/addon/fold/markdown-fold.js";
import "codemirror/addon/fold/comment-fold.js";

import { createNewLab } from "@/api/table";
import { getFileContent } from "@/api/table";
import { getFileChanged } from "@/api/table";
export default {
  data() {
    return {
      form: {
        name: "",
        region: "",
        date1: "",
        date2: "",
        delivery: false,
        type: [],
        resource: "",
        desc: "",
      },
      dataJson: {},
      codedata: null,
      curCode: "dddlskadjf;lakjtiruat",
      fileres: "",
      cmOptions: {
        value: "",
        mode: "text/x-lua",
        // theme: "default",
        lineNumbers: true, // 是否显示行数
        line: true,
        readOnly: false,
      },
    };
  },
  created() {
    var param = localStorage.getItem("experParam");
    this.param = param;
    console.log("t------- param");
    console.log(param);
    this.setParam(param);
  },

  mounted() {
    var param = localStorage.getItem("experParam");
    this.param = param;
    console.log("this is param");
    console.log(param);
    this.makeElementAbled();
    if (param != null) {
      this.setParam(param);
    }
    localStorage.removeItem("experParam");

    this.timer = setInterval(() => {
      this.getFileChange();
    }, 60000);
  },
  methods: {
    getFileChange() {
      getFileChanged().then((response) => {
        // console.log("===================")
        // console.log(typeof(response.data))
        // // console.log(type(response.data))
        // console.log("===================")
        for (let resv in response.data) {
          this.fileres = this.fileres + response.data[resv] + "\n";
        }

        // console.log("res is");
        // console.log(this.fileres);
      });
    },
    setParam(pa) {
      this.dataJson = JSON.parse(pa);
      console.log("param in set function");
      console.log(pa);

      this.form["name"] = this.dataJson["file_name"];
      this.form["region"] = this.dataJson["exper_time"];
      // this.form["date1"] = this
      this.form["delivery"] = this.dataJson["isPublic"];
      var type = this.dataJson["exper_type"];

      if (type == 1) {
        this.form["resource"] = "节点保留策略";
      } else if (type == 2) {
        this.form["resource"] = "节点选择策略";
      } else if (type == 3) {
        this.form["resource"] = "节点发布策略";
      } else if (type == 4) {
        this.form["resource"] = "节点获取策略";
      }
      this.form["desc"] = this.dataJson["desc"];

      this.getFileContent(this.form["name"]);
      this.makeElementDisabled();
    },
    getFileContent(filename) {
      var param = {
        filename: filename,
      };
      getFileContent(param).then((response) => {
        this.curCode = response.data;
      });
    },
    makeElementDisabled() {
      let x = document.getElementById("filenameinputing");
      x.disabled = "true";
      console.log(x);
      let x1 = document.getElementById("timeselector");
      //
      x1.disabled = "true";
      let x2 = document.getElementById("date1pocker");
      //
      x2.disabled = "true";
      let x3 = document.getElementById("date2pocker");
      //
      x3.disabled = "true";
      let x4 = document.getElementById("publicswitcher");
      console.log(x4);
      //
      x4.disabled = "true";
      let x5 = document.getElementById("expertyper");
      //
      x5.disabled = "true";
      let x6 = document.getElementById("descinput");
      //
      x6.disabled = "true";
      let x7 = document.getElementById("buttonsubmit");
      //
      // x7.disabled = "true";
      x7.style.display = "none"
      let x8 = document.getElementById("buttoncancel");
      //
      // x8.disabled = true;
      x8.style.display = "none"

      let x9 = document.getElementById("buttonedit");
      // x9.style.display = "none"
      x9.innerText = "查看实验详情"

      let x10 = document.getElementById("buttonrun");
      x10.style.display = "none"
      // console.log(x9)
      console.log(x10)
    },
    makeElementAbled() {},
    onCmCodeChange(newCode) {
      this.curCode = newCode;
    },
    onSubmit() {
      this.$message("submit!");
    },
    onCancel() {
      this.$message({
        message: "cancel!",
        type: "warning",
      });
    },
    editExper() {
      var luaCodeText = this.curCode;

      // alert(luaCodeText);
      var luaFileName = this.form.name;
      var experPoint = this.form.resource;
      var experTime = this.form.region;
      var startTime = "0";
      var addTime = "0";
      var isPublic = this.form.delivery;
      var date1 = this.form.date1;
      var date2 = this.form.date2;
      console.log("Point:", experPoint);
      console.log("experTime:", experTime);
      console.log("startTime:", startTime);
      console.log("addTime:", addTime);
      console.log("isPublic:", isPublic);
      console.log("date1:", date1);
      console.log("date2:", date2);
      var descText = this.form.desc;

      // alert(luaFileName);
      var param = {
        filecontent: luaCodeText,
        filename: luaFileName,
        experPoint: experPoint,
        experTime: experTime,
        startTime: startTime,
        addTime: addTime,
        user: "admin",
        date1: date1,
        date2: date2,
        isPublic: isPublic,
        desc: descText,
      };
      createNewLab(param).then((response) => {
        this.list = response.data;
        this.listLoading = false;
        // return response;
      });
    },
    runExper() {},
  },
  components: {
    codemirror,
  },
};
</script>

<style scoped>
.line {
  text-align: center;
}
.col-md-12 {
  -webkit-box-flex: 0;
  -webkit-flex: 0 0 100%;
  -ms-flex: 0 0 100%;
  flex: 0 0 100%;
  max-width: 100%;
}
.row {
  display: -webkit-box;
  display: -webkit-flex;
  display: -ms-flexbox;
  display: flex;
  -webkit-flex-wrap: wrap;
  -ms-flex-wrap: wrap;
  flex-wrap: wrap;
  margin-right: -15px;
  margin-left: -15px;
}
.alert {
  padding: 0.75rem 1.25rem;
  margin-bottom: 1rem;
  border: 1px solid transparent;
  border-radius: 0.25rem;
}
.alert-warning {
  background-color: #fcf8e3;
  border-color: #faf2cc;
  color: #8a6d3b;
}

.alert-warning hr {
  border-top-color: #f7ecb5;
}

.alert-warning .alert-link {
  color: #66512c;
}

.btn {
  display: inline-block;
  font-weight: 400;
  line-height: 1.25;
  text-align: center;
  white-space: nowrap;
  vertical-align: middle;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  user-select: none;
  border: 1px solid transparent;
  padding: 0.5rem 1rem;
  font-size: 1rem;
  border-radius: 0.25rem;
  -webkit-transition: all 0.2s ease-in-out;
  -o-transition: all 0.2s ease-in-out;
  transition: all 0.2s ease-in-out;
}

.btn-primary:hover {
  color: #fff;
  background-color: #025aa5;
  border-color: #01549b;
}
.btn-outline-primary {
  color: #0275d8;
  background-image: none;
  background-color: transparent;
  border-color: #0275d8;
}

.col-md-6 {
  -webkit-box-flex: 0;
  -webkit-flex: 0 0 70%;
  -ms-flex: 0 0 70%;
  flex: 0 0 70%;
  max-width: 70%;
}
.col-md-60 {
  -webkit-box-flex: 0;
  -webkit-flex: 0 0 30%;
  -ms-flex: 0 0 30%;
  flex: 0 0 30%;
  max-width: 30%;
}

.form-control {
  display: block;
  width: 100%;
  padding: 0.5rem 0.75rem;
  font-size: 1rem;
  line-height: 1.25;
  color: #464a4c;
  background-color: #fff;
  background-image: none;
  -webkit-background-clip: padding-box;
  background-clip: padding-box;
  border: 1px solid rgba(0, 0, 0, 0.15);
  border-radius: 0.25rem;
  -webkit-transition: border-color ease-in-out 0.15s,
    -webkit-box-shadow ease-in-out 0.15s;
  transition: border-color ease-in-out 0.15s,
    -webkit-box-shadow ease-in-out 0.15s;
  -o-transition: border-color ease-in-out 0.15s, box-shadow ease-in-out 0.15s;
  transition: border-color ease-in-out 0.15s, box-shadow ease-in-out 0.15s;
  transition: border-color ease-in-out 0.15s, box-shadow ease-in-out 0.15s,
    -webkit-box-shadow ease-in-out 0.15s;
}
.btn-block {
  display: block;
  width: 100%;
}

.CodeMirror-line {
  font-family: Consolas, "Microsoft Yahei" !important;
}
</style>

