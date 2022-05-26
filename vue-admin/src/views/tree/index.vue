<template>
  <div class="app-container">
    <el-input
      v-model="filterText"
      placeholder="Filter keyword"
      style="margin-bottom: 30px"
    />

    <el-tree
      ref="tree2"
      :data="data3"
      :props="defaultProps"
      :filter-node-method="filterNode"
      class="filter-tree"
      default-expand-all
    />
  </div>
</template>

<script>
import { getRouterInfo } from "@/api/table";
export default {
  data() {
    return {
      filterText: "",
      dataJson: null,
      data3: [],
      data2: [
        {
          id: 1,
          label: "Level one 1",
          children: [
            {
              id: 4,
              label: "Level two 1-1",
              children: [
                {
                  id: 9,
                  label: "Level three 1-1-1",
                },
                {
                  id: 10,
                  label: "Level three 1-1-2",
                },
              ],
            },
          ],
        },
        {
          id: 2,
          label: "Level one 2",
          children: [
            {
              id: 5,
              label: "Level two 2-1",
            },
            {
              id: 6,
              label: "Level two 2-2",
            },
          ],
        },
        {
          id: 3,
          label: "Level one 3",
          children: [
            {
              id: 7,
              label: "Level two 3-1",
            },
            {
              id: 8,
              label: "Level two 3-2",
            },
          ],
        },
      ],
      defaultProps: {
        children: "children",
        label: "label",
      },
    };
  },
  watch: {
    filterText(val) {
      this.$refs.tree2.filter(val);
    },
  },
  created() {
    var param = localStorage.getItem("searchParam");
    this.param = param;
    console.log(param);
    this.setTreeData(param);
  },

  methods: {
    setTreeData(param) {
      this.dataJson = JSON.parse(param);
      //开始构造数据树
      var ids = 0;
      console.log(this.dataJson["pubkey"]);
      var pubkey = this.dataJson["pubkey"];
      var pubkeyValueJson = {
        // id: ids,
        label: pubkey,
      };
      ids = ids + 1;
      var pubkeyJson = {
        // id: ids,
        label: "pubkey",
        children: [pubkeyValueJson],
      };
      ids = ids + 1;
      this.data3.push(pubkeyJson);

      console.log(this.dataJson["signkey"]);
      var signkey = this.dataJson["signkey"];
      var signkeyValueJson = {
        // id: ids,
        lable: signkey,
      };
      ids = ids + 1;
      var signkeyJson = {
        // id: ids,
        label: "signkey",
        children: [signkeyValueJson],
      };
      ids = ids + 1;
      this.data3.push(signkeyJson);
      console.log(signkeyValueJson);
      console.log(signkeyJson);
      console.log(this.dataJson["published"]);
      var published = this.dataJson["published"];
      var displayTime = this.getDisplayTime(published);
      var publishValueJson = {
        label: displayTime,
      };
      var publishJson = {
        label: "published",
        children: [publishValueJson],
      };
      this.data3.push(publishJson);
      console.log(this.dataJson["signature"]);
      var signature = this.dataJson["signature"];
      var signatureValueJson = {
        label: signature,
      };
      var signatureJson = {
        label: "signature",
        children: [signatureValueJson],
      };
      this.data3.push(signatureJson);

      //-------options------------
      var capsValue = this.dataJson["options"]["caps"];
      var capsValueJson = {
        label: capsValue,
      };
      var capsJson = {
        label: "caps",
        children: [capsValueJson],
      };

      var netIdValue = this.dataJson["options"]["netId"];
      var netIdValueJson = {
        label: netIdValue,
      };
      var netIdJson = {
        label: "netId",
        children: [netIdValueJson],
      };

      var netdbKnownLeaseSetsValue =
        this.dataJson["options"]["netdb.knownLeaseSets"];
      var netDbKValueJson = {
        label: netdbKnownLeaseSetsValue,
      };
      var netDbKJson = {
        label: "netdb.knownLeaseSets",
        children: [netDbKValueJson],
      };

      var netDbKRValue = this.dataJson["options"]["netdb.knownRouters"];
      var netDbKRValueJson = {
        label: netDbKRValue,
      };
      var netDbKRJson = {
        label: "netdb.konwnRouters",
        children: [netDbKRValueJson],
      };

      var rversionValue = this.dataJson["options"]["router.version"];
      var rversionValueJson = {
        label: rversionValue,
      };
      var rversionJson = {
        label: "router.version",
        children: [rversionValueJson],
      };

      var optionJson = {
        label: "options",
        children: [rversionJson, netDbKRJson, netDbKJson, netIdJson, capsJson],
      };
      this.data3.push(optionJson);

      //------cert
      var signatureTypeValue = this.dataJson["cert"]["signature_type"];
      var signatureTypeValueJson = {
        label: signatureTypeValue,
      };
      var signatyreTypeJson = {
        label: "signature_type",
        children: [signatureTypeValueJson],
      };

      var cryptoType = this.dataJson["cert"]["crypto_type"];
      var cryptoTypeValueJson = {
        label: cryptoType,
      };
      var cryptoTypeJson = {
        label: "crypto_type",
        children: [cryptoTypeValueJson],
      };

      var certJson = {
        label: "cert",
        children: [signatyreTypeJson, cryptoTypeJson],
      };

      this.data3.push(certJson);

      //--------addrs
      var addrs = this.dataJson["addrs"];
      for (var i = 0, l = addrs.length; i < l; i++) {
        var costValue = addrs[i]["cost"];
        var costValueJson = {
          label: costValue,
        };
        var costJson = {
          label: "cost",
          children: [costValueJson],
        };

        var transportValue = addrs[i]["transport"];
        var transportValueJson = {
          label: transportValue,
        };
        var transportJson = {
          label: "transport",
          children: [transportValueJson],
        };

        var hostValue = addrs[i]["options"]["host"];
        var hostValueJson = {
          label: hostValue,
        };
        var hostJson = {
          label: "host",
          children: [hostValueJson],
        };

        var portValue = addrs[i]["options"]["port"];
        var portValueJson = {
          label: portValue,
        };
        var portJson = {
          label: "port",
          children: [portValueJson],
        };

        var capsValue = addrs[i]["options"]["caps"];
        var capsValueJson = {
          label: capsValue,
        };
        var capsJson = {
          label: "caps",
          children: [capsValueJson],
        };

        var keyValue = addrs[i]["options"]["key"];
        var keyValueJson = {
          label: keyValue,
        };
        var keyJson = {
          label: "key",
          children: [keyValueJson],
        };

        var iValue = addrs[i]["options"]["i"];
        var iValueJson = {
          label: iValue,
        };
        var iJson = {
          label: "i",
          children: [iValueJson],
        };

        var sValue = addrs[i]["options"]["s"];
        var sValueJson = {
          label: sValue,
        };
        var sJson = {
          label: "s",
          children: [sValueJson],
        };

        var vValue = addrs[i]["options"]["v"];
        var vValueJson = {
          label: vValue,
        };
        var vJson = {
          label: "v",
          children: [vValueJson],
        };

        var optionsJson = {
          label: "opitons",
          children: [
            keyJson,
            capsJson,
            portJson,
            hostJson,
            iJson,
            sJson,
            vJson,
          ],
        };

        var addrsJson = {
          label: "addrs" + (i + 1),
          children: [transportJson, costJson, optionsJson],
        };

        this.data3.push(addrsJson);
      }
    },
    filterNode(value, data) {
      if (!value) return true;
      return data.label.indexOf(value) !== -1;
    },
    getMetaData(pubkey) {
      var param = {
        pubkey: pubkey,
      };
      getRouterInfo(param).then((response) => {
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
};
</script>

