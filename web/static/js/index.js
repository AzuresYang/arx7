// new Vue({
//     el: '#app',
//     data: {
//       message: 'Hello Vue.js!'
//     }
//   })
//基于准备好的DOM，初始化echarts实例

//指定图表的配置项和数据
function getRootPath() {
  return "http://127.0.0.1:8888/"
}

var isRealTimeMonitor = false
var tRealTime
var myChart = echarts.init(document.getElementById('chart-monitor'));
var monitor_option = {
  title: {
    text: 'ARX监控数据',
    textStyle: {
      color: '#0080ff',
      fontWeight: 'bold',
    }
  },
  //提示框组件
  tooltip: {
    //坐标轴触发，主要用于柱状图，折线图等
    trigger: 'axis',
  },
  //图例
  legend: {
    data: ['曲线1']
  },
  //横轴
  xAxis: {
    data: []
  },
  //纵轴
  yAxis: {},
  //系列列表。每个系列通过type决定自己的图表类型
  series: [{
    name: '曲线1',
    //折线图
    type: 'line',
    smooth: true, //数据光滑过度
    symbol: 'none', //下一个数据点
    stack: 'a',
    // areaStyle: { // 这个是填充区域
    //   normal: {
    //     color: 'blue'
    //   }
    // },
    data: []
  }]
};

// 日期操作
Date.prototype.DateLess = function (strInterval, Number) {
  var dtTmp = this;
  switch (strInterval) {
    case 's': return new Date(Date.parse(dtTmp) - (1000 * Number));
    case 'n': return new Date(Date.parse(dtTmp) - (60000 * Number));
    case 'h': return new Date(Date.parse(dtTmp) - (3600000 * Number));
    case 'd': return new Date(Date.parse(dtTmp) - (86400000 * Number));
    case 'w': return new Date(Date.parse(dtTmp) - ((86400000 * 7) * Number));
    case 'q': return new Date(dtTmp.getFullYear(), (dtTmp.getMonth()) - Number * 3, dtTmp.getDate(), dtTmp.getHours(), dtTmp.getMinutes(), dtTmp.getSeconds());
    case 'm': return new Date(dtTmp.getFullYear(), (dtTmp.getMonth()) - Number, dtTmp.getDate(), dtTmp.getHours(), dtTmp.getMinutes(), dtTmp.getSeconds());
    case 'y': return new Date((dtTmp.getFullYear() - Number), dtTmp.getMonth(), dtTmp.getDate(), dtTmp.getHours(), dtTmp.getMinutes(), dtTmp.getSeconds());
  }
}

Date.prototype.format = function (format) {
  var o = {
    "M+": this.getMonth() + 1, //month
    "d+": this.getDate(),    //day
    "h+": this.getHours(),   //hour
    "m+": this.getMinutes(), //minute
    "s+": this.getSeconds(), //second
    "q+": Math.floor((this.getMonth() + 3) / 3),  //quarter
    "S": this.getMilliseconds() //millisecond
  }
  if (/(y+)/.test(format)) format = format.replace(RegExp.$1,
    (this.getFullYear() + "").substr(4 - RegExp.$1.length));
  for (var k in o) if (new RegExp("(" + k + ")").test(format))
    format = format.replace(RegExp.$1,
      RegExp.$1.length == 1 ? o[k] :
        ("00" + o[k]).substr(("" + o[k]).length));
  return format;
}

$(document).ready(function () {
  myChart.setOption(monitor_option)
});

// $("ifRealTime").is(':checked');
$("#bt_querymonitor").click(function () {

  // 定时刷新， 实时监控数据
  if ($("#ifRealTime").is(':checked')) {
    isRealTimeMonitor = true
    console.log("real time monitor ")
    tRealTime = window.setInterval(function () {

      var end_time = new Date();
      var start_time = end_time.DateLess('n', 3)
      $("#datetimepicker1").datetimepicker("setDate", start_time);

      $("#datetimepicker2").datetimepicker("setDate", end_time);
      console.log("end_time", end_time, "input:", $("#end_time").val())
      var data = $("#monitor_form").serialize();
      console.log("post data:" + data)
      refreshMonitorInfo(data)
    }, 1000);
  } else {
    isRealTimeMonitor = false
    var end_time = new Date();
    var start_time = end_time.DateLess('m', 5)

    var data = $("#monitor_form").serialize();
    console.log("post data:" + data)
    window.clearInterval(tRealTime);
    refreshMonitorInfo(data);
  }
})

// 停止刷新
$("#ifRealTime").click(function () {
  if ($("#ifRealTime").is(':checked') == false) {
    window.clearInterval(tRealTime);
  }
})
function refreshMonitorInfo(data) {
  var url = getRootPath() + "get/monitor"
  $.post(url,
    data,
    function (data, status) {

      //  data = "'"+data+"'";
      console.log("get resp status:", status)
      // console.log("resp:", data)
      var resp = JSON.parse(data);

      // console.log("resp data:", resp.Data)
      //var moni_data = JSON.parse(resp.Data)
      if (resp.Status == 0) {
        myChart.hideLoading();
        series = resp.Data.series
        xAxisData = resp.Data.xAxsis
        // console.log("last xAsis:" + xAxisData)
        // console.log("last series:" + series)
        myChart.setOption({
          xAxis: {
            data: xAxisData
          },
          //纵轴
          yAxis: {},
          //系列列表。每个系列通过type决定自己的图表类型
          series: [{
            data: series
          }]
        })
      } else {
        alert("查询错误:" + resp.Msg)
      }

    });
}


$(function () {

  //1.初始化Table
  var oTable = new SpiderTableInit();
  oTable.Init();

  var operateTable = new OperateTableInit();
  operateTable.Init();

});


var SpiderTableInit = function () {
  var oTableInit = new Object();
  //初始化Table
  oTableInit.Init = function () {
    var url = getRootPath() + "get/pods"
    console.log("url", url)
    $('#tb_spider').bootstrapTable({
      url: url,         //请求后台的URL（*）
      method: 'post',                      //请求方式（*）
      toolbar: '#cluster_toolbar',                //工具按钮用哪个容器
      striped: true,                      //是否显示行间隔色
      cache: false,                       //是否使用缓存，默认为true，所以一般情况下需要设置一下这个属性（*）
      pagination: true,                   //是否显示分页（*）
      sortable: false,                     //是否启用排序
      sortOrder: "asc",                   //排序方式
      queryParams: oTableInit.queryParams,//传递参数（*）
      sidePagination: "server",           //分页方式：client客户端分页，server服务端分页（*）
      pageNumber: 1,                       //初始化加载第一页，默认第一页
      pageSize: 5,                       //每页的记录行数（*）
      pageList: [5, 10, 25, 50, 100],        //可供选择的每页的行数（*）
      search: true,                       //是否显示表格搜索，此搜索是客户端搜索，不会进服务端，所以，个人感觉意义不大
      strictSearch: true,
      showColumns: true,                  //是否显示所有的列
      showRefresh: true,                  //是否显示刷新按钮
      minimumCountColumns: 2,             //最少允许的列数
      clickToSelect: true,                //是否启用点击选中行
      height: 500,                        //行高，如果没有设置height属性，表格自动根据记录条数觉得表格高度
      // uniqueId: "ID",                     //每一行的唯一标识，一般为主键列
      showToggle: true,                    //是否显示详细视图和列表视图的切换按钮
      cardView: false,                    //是否显示详细视图
      detailView: false,                   //是否显示父子表
      columns: [{
        field: 'SpiderName',
        title: '爬虫名'
      }, {
        field: 'NodeStatus',
        title: '节点状态'
      }, {
        field: 'RunStatus',
        title: '运行状态'
      }, {
        field: 'Age',
        title: '运行时间'
      }, {
        field: 'NodeAddr',
        title: '节点IP'
      },
      {
        field: 'Desc',
        title: '描述'
      },
      ]
    });
  };

  //得到查询的参数
  oTableInit.queryParams = function (params) {
    var temp = {   //这里的键的名字和控制器的变量名必须一直，这边改动，控制器也需要改成一样的

    };
    return temp;
  };
  return oTableInit;
};


var OperateTableInit = function () {
  var oTableInit = new Object();
  //初始化Table
  oTableInit.Init = function () {
    $('#tb_operate').bootstrapTable({
      toolbar: '#cluster_operate_toolbar',  //工具按钮用哪个容器
      striped: true,                      //是否显示行间隔色
      cache: false,                       //是否使用缓存，默认为true，所以一般情况下需要设置一下这个属性（*）
      pagination: true,                   //是否显示分页（*）
      sortable: false,                     //是否启用排序
      queryParams: oTableInit.queryParams,//传递参数（*）
      sidePagination: "client",           //分页方式：client客户端分页，server服务端分页（*）
      pageNumber: 1,                       //初始化加载第一页，默认第一页
      pageSize: 5,                       //每页的记录行数（*）
      pageList: [5, 10, 25, 50, 100],        //可供选择的每页的行数（*）
      search: true,                       //是否显示表格搜索，此搜索是客户端搜索，不会进服务端，所以，个人感觉意义不大
      strictSearch: true,
      showColumns: true,                  //是否显示所有的列
      showRefresh: false,                  //是否显示刷新按钮
      minimumCountColumns: 2,             //最少允许的列数
      clickToSelect: true,                //是否启用点击选中行
      height: 500,                        //行高，如果没有设置height属性，表格自动根据记录条数觉得表格高度
      // uniqueId: "ID",                     //每一行的唯一标识，一般为主键列
      showToggle: true,                    //是否显示详细视图和列表视图的切换按钮
      cardView: false,                    //是否显示详细视图
      detailView: false,                   //是否显示父子表
      columns: [{
        field: 'Time',
        width: "15%",
        title: '时间'
      }, {
        field: 'Operate',
        width: "15%",
        title: '操作'
      }, {
        field: 'Msg',
        title: '操作结果'
      },
      ]
    });
  };

  //得到查询的参数
  oTableInit.queryParams = function (params) {
    var temp = {   //这里的键的名字和控制器的变量名必须一直，这边改动，控制器也需要改成一样的

    };
    return temp;
  };
  return oTableInit;
};


function updateOperateTable(operate, msg) {
  var time = new Date()
  record = { Time: time.format("yyyy-MM-dd hh:mm:ss"), Operate: operate, Msg: msg }
  $('#tb_operate').bootstrapTable('prepend', record);
}


$("#bt_queryPods").click(function () {
  var data = "post no data"
  var url = getRootPath() + "get/pods"
  console.log("url", url)
  $.post(url,
    data,
    function (data, status) {
      $('#tb_spider').bootstrapTable('load', data);
      updateOperateTable("查询集群信息", "succ")
    });
});

// 获取spider状态
$("#bt_cluster_status").click(function () {
  var data = $("#form_cluster_status").serialize();
  var url = getRootPath() + "cluster/spiderstatus"
  console.log("url", url)
  $.post(url,
    data,
    function (data, status) {

      var resp = JSON.parse(data);
      updateOperateTable("获取爬虫状态", resp.Msg)
      if (resp.Status == 0) {
        console.log("get resp Data:", resp)
        // $('#tb_spider').bootstrapTable('refresh'); 
      } else {
        alert("查询错误:" + resp.Msg)
      }
    });
});

// 部署爬虫
$("#bt_cluster_deployment").click(function () {
  var data = $("#form_cluster_deployment").serialize();
  var url = getRootPath() + "cluster/deployment"
  console.log("url", url)
  $.post(url,
    data,
    function (data, status) {
      var resp = JSON.parse(data);
      updateOperateTable("部署爬虫", resp.Msg)
      if (resp.Status == 0) {
        console.log("get resp Data:", resp)
        // $('#tb_spider').bootstrapTable('refresh'); 
      } else {
        alert("部署爬虫错误:" + resp.Msg)
      }
    });
});

// 启动爬虫
$("#bt_cluster_start").click(function () {
  var data = $("#form_cluster_start").serialize();
  // var data = new FormData(document.querySelector("form_cluster_start"));//获取form值
  console.log("start data:", data)
  var url = getRootPath() + "cluster/start"
  console.log("url", url)
  // $.post(url,
  //   data,
  //   function (data, status) {

  //     var resp = JSON.parse(data);
  //     updateOperateTable("启动爬虫", resp.Msg)
  //     if (resp.Status == 0) {
  //       console.log("get resp Data:", resp)
  //      //  $('#tb_spider').bootstrapTable('refresh'); 
  //     } else {
  //       alert("启动爬虫错误:" + resp.Msg)
  //     }
  //   });
  
  // var fd = new FormData(document.querySelector("form_cluster_start"));
  var fd = new FormData();
  fd.append('config', $('#config')[0].files[0]);
  fd.append('spidername',$("#stspidername").val())
  $.ajax({
    url: url,
    type: "POST",
    data: fd,
    processData: false,  // 不处理数据
    contentType: false,   // 不设置内容类型
    success: function (data) {
      var resp = JSON.parse(data);
      updateOperateTable("启动爬虫", resp.Msg)
      if (resp.Status == 0) {
        console.log("get resp Data:", resp)
        //  $('#tb_spider').bootstrapTable('refresh'); 
      } else {
        alert("启动爬虫错误:" + resp.Msg)
      }
    }
  });
});


// 扩缩容
$("#bt_cluster_scale").click(function () {
  var data = $("#form_cluster_scale").serialize();
  var url = getRootPath() + "cluster/scale"
  console.log("url", url)
  $.post(url,
    data,
    function (data, status) {

      var resp = JSON.parse(data);
      updateOperateTable("爬虫扩缩容", resp.Msg)
      if (resp.Status == 0) {
        console.log("get resp Data:", resp)
        // $('#tb_spider').bootstrapTable('refresh'); 
      } else {
        alert("爬虫扩缩容操作错误:" + resp.Msg)
      }
    });
});

// 删除
$("#bt_cluster_delete").click(function () {
  var data = $("#form_cluster_delete").serialize();
  var url = getRootPath() + "cluster/delete"
  console.log("url", url)
  $.post(url,
    data,
    function (data, status) {

      var resp = JSON.parse(data);
      updateOperateTable("删除爬虫", resp.Msg)
      if (resp.Status == 0) {
        console.log("get resp Data:", resp)
        $('#tb_spider').bootstrapTable('refresh');
      } else {
        alert("删除爬虫操作错误:" + resp.Msg)
      }
    });
});