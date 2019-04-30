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

var xAxisData = ["衬衫", "羊毛衫", "雪纺衫", "裤子", "高跟鞋", "袜子"]
var series = [1, 5, 20, 5, 1, 8]
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
Date.prototype.DateLess = function(strInterval, Number) {   
  var dtTmp = this;  
  switch (strInterval) {   
      case 's' :return new Date(Date.parse(dtTmp) - (1000 * Number));  
      case 'n' :return new Date(Date.parse(dtTmp) - (60000 * Number));  
      case 'h' :return new Date(Date.parse(dtTmp) - (3600000 * Number));  
      case 'd' :return new Date(Date.parse(dtTmp) - (86400000 * Number));  
      case 'w' :return new Date(Date.parse(dtTmp) - ((86400000 * 7) * Number));  
      case 'q' :return new Date(dtTmp.getFullYear(), (dtTmp.getMonth()) - Number*3, dtTmp.getDate(), dtTmp.getHours(), dtTmp.getMinutes(), dtTmp.getSeconds());  
      case 'm' :return new Date(dtTmp.getFullYear(), (dtTmp.getMonth()) - Number, dtTmp.getDate(), dtTmp.getHours(), dtTmp.getMinutes(), dtTmp.getSeconds());  
      case 'y' :return new Date((dtTmp.getFullYear() - Number), dtTmp.getMonth(), dtTmp.getDate(), dtTmp.getHours(), dtTmp.getMinutes(), dtTmp.getSeconds());  
  }  
}  

$(document).ready(function () {
  myChart.setOption(monitor_option)
});


$("#bt_querymonitor").click(function () {
  var url = getRootPath() + "get/monitor"
  var data = $("#monitor_form").serialize();
  console.log("post data:" + data)
  console.log("query monitor:" + url)
  $.post(url,
    data,

    function (data, status) {

      //  data = "'"+data+"'";
      console.log("resp:", data)
      var resp = JSON.parse(data);

      console.log("resp data:", resp.Data)
      //var moni_data = JSON.parse(resp.Data)
      if (resp.Status == 0) {

        // monitor_option.series.data = resp.Data.series;
        // monitor_option.xAxis.data = resp.Data.xAxis;
        // monitor_option.series.data = ["衬衫", "羊毛衫", "雪纺衫", "裤子"];//resp.Data.series;
        // monitor_option.xAxis.data = [5, 20, 36, 10];
        //使用刚指定的配置项和数据显示图表
        // var myChart = echarts.init(document.getElementById('chart-monitor'));
        myChart.hideLoading();
        series = resp.Data.series
        xAxisData = resp.Data.xAxsis
        console.log("last xAsis:" + xAxisData)
        console.log("last series:" + series)
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
})