<template>
  <div class="visualization-container">
    <div class="choice">
      <el-form v-model="options">
        <el-select v-model="options.bill_value" placeholder="请选择bill值高于或低于"
                   v-show="(options.weight_value === '')&& (options.height_value === '')&&(options.weight === '')&&(options.height === '')">
          <el-option
            v-for="item in options.bill_options"
            :key="item.value"
            :label="item.label"
            :value="item.value"
            :disabled="item.disabled">
          </el-option>
        </el-select>
        <el-input v-model="options.bill" placeholder="请输入bill值"
                  v-show="(options.weight_value === '')&&(options.height_value === '')&&(options.weight === '')&&(options.height === '')"></el-input>
        <el-select v-model="options.weight_value" placeholder="请选择weight值高于或低于"
                   v-show="(options.bill_value === '')&&(options.bill === '')">
          <el-option
            v-for="item in options.weight_options"
            :key="item.value"
            :label="item.label"
            :value="item.value"
            :disabled="item.disabled">
          </el-option>
        </el-select>
        <el-input v-model="options.weight" placeholder="请输入weight值"
                  v-show="(options.bill_value === '')&&(options.bill === '')"></el-input>
        <el-select v-model="options.height_value" placeholder="请选择height值高于或低于"
                   v-show="(options.bill_value === '')&&(options.bill === '')">
          <el-option
            v-for="item in options.height_options"
            :key="item.value"
            :label="item.label"
            :value="item.value"
            :disabled="item.disabled">
          </el-option>
        </el-select>
        <el-input v-model="options.height" placeholder="请输入height值"
                  v-show="(options.bill_value === '')&&(options.bill === '')"></el-input>
      </el-form>
      <el-button type="primary" @click="submit_choice()">提交</el-button>
    </div>

    <div class="charts">
      <el-card class="box-card" shadow="hover">
        <div class="bar" id="Chart_bar" :style="{width:'500px',height:'500px'}"></div>
      </el-card>
      <el-card style="position: absolute; width: 40%; left: 58%" shadow="hover">
        <div class="pie" id="Chart_pie" :style="{width:'500px',height:'500px'}"></div>
      </el-card>

    </div>
  </div>
</template>

<script>
import '@/theme/westeros'
import { analysis } from '@/api/patient'

const echarts = require('echarts/lib/echarts')
require('echarts/lib/chart/bar')
require('echarts/lib/chart/pie')
export default {
  name: 'VisualizationIndex',
  components: {},
  props: {},
  data () {
    return {
      result_pie: {
        series: [
          {
            name: '数据统计',
            type: 'pie', // 设置图表类型为饼图
            radius: '55%', // 饼图的半径
            data: [],
            label: {
              normal: {
                show: true,
                position: 'inside',
                formatter: '{b}: {c}({d}%)',
                textStyle: {
                  fontSize: 16
                }
              }
            }
          }
        ]
      },
      result_bar: {
        xAxis: {
          data: []
        },
        yAxis: {},
        series: [{
          name: '查询数目',
          type: 'bar',
          data: [],
          barWidth: 100,
          itemStyle: { // 上方显示数值
            normal: {
              label: {
                show: true, // 开启显示
                position: 'top', // 在上方显示
                textStyle: { // 数值样式
                  color: 'black',
                  fontSize: 16
                }
              }
            }
          }
        }]
      },
      options: {
        bill_value: '',
        weight_value: '',
        height_value: '',
        bill: '',
        weight: '',
        height: '',
        bill_options: [{
          value: 'above',
          label: 'bill高于'
        }, {
          value: 'below',
          label: 'bill低于'
        }, {
          value: '',
          label: '未选择'
        }],
        weight_options: [{
          value: 'above',
          label: 'weight高于'
        }, {
          value: 'below',
          label: 'weight低于'
        }, {
          value: '',
          label: '未选择'
        }],
        height_options: [{
          value: 'above',
          label: 'height高于'
        }, {
          value: 'below',
          label: 'height低于'
        }, {
          value: '',
          label: '未选择'
        }]
      }
    }
  },
  computed: {},
  watch: {},
  created () {
  },
  mounted () {
  },
  methods: {
    submit_choice () {
      let payload = {}
      if ((this.options.bill_value === '') && (this.options.weight_value === '') && (this.options.height_value === '')) {
        this.$message.error('请选择一项选择框的内容')
      } else {
        payload = {
          weight: parseInt(this.options.weight),
          high: parseInt(this.options.high),
          bill: parseInt(this.options.bill)
        }
        analysis(payload)
          .then((response) => {
            const data = response.data.payload
            if (this.options.bill !== '') {
              // bill高于
              if (this.options.bill_value === 'above') {
                this.result_pie.series[0].data = [{
                  value: data.bill_above,
                  name: '查询结果'
                }, {
                  value: data.bill_below,
                  name: '其他'
                }]
                this.result_bar.xAxis.data = ['bill>' + this.options.bill, '其他', '总数据量']
                this.result_bar.series[0].data = [data.bill_above, data.bill_below, data.data_length]
                this.fetch()
              } else if (this.options.bill_value === 'below') {
                this.result_pie.series[0].data = [{
                  value: data.bill_below,
                  name: '查询结果'
                }, {
                  value: data.bill_above,
                  name: '其他'
                }]
                this.result_bar.xAxis.data = ['bill<' + this.options.bill, '其他', '总数据量']
                this.result_bar.series[0].data = [data.bill_below, data.bill_above, data.data_length]
                this.fetch()
              }
            } else if (this.options.weight_value === 'above') {
              //    height未选择
              if (this.options.height_value === '') {
                this.result_pie.series[0].data = [{
                  value: data.weight_above,
                  name: '查询结果'
                }, {
                  value: data.weight_below,
                  name: '其他'
                }]
                this.result_bar.xAxis.data = ['weight>' + this.options.weight, '其他', '总数据量']
                this.result_bar.series[0].data = [data.weight_above, data.weight_below, data.data_length]
                this.fetch()
              } else if (this.options.height_value === 'above') {
                const other = data.data_length - data.weight_above_high_above
                this.result_pie.series[0].data = [{
                  value: data.weight_above_high_above,
                  name: '查询结果'
                }, {
                  value: other,
                  name: '其他'
                }]
                this.result_bar.xAxis.data = ['weight>' + this.options.weight + ',height>' + this.options.height, '其他', '总数据量']
                this.result_bar.series[0].data = [data.weight_above_high_above, other, data.data_length]
                this.fetch()
              } else {
                const other = data.data_length - data.weight_above_high_below
                this.result_pie.series[0].data = [{
                  value: data.weight_above_high_below,
                  name: '查询结果'
                }, {
                  value: other,
                  name: '其他'
                }]
                this.result_bar.xAxis.data = ['weight>' + this.options.weight + ',height<' + this.options.height, '其他', '总数据量']
                this.result_bar.series[0].data = [data.weight_above_high_below, other, data.data_length]
                this.fetch()
              }
            } else if (this.options.weight_value === 'below') {
              // height未选择
              if (this.options.height_value === '') {
                this.result_pie.series[0].data = [{
                  value: data.weight_below,
                  name: '查询结果'
                }, {
                  value: data.weight_above,
                  name: '其他'
                }]
                this.result_bar.xAxis.data = ['weight<' + this.options.weight, '其他', '总数据量']
                this.result_bar.series[0].data = [data.weight_below, data.weight_above, data.data_length]
                this.fetch()
              } else if (this.options.height_value === 'above') {
                const other = data.data_length - data.weight_below_high_above
                this.result_pie.series[0].data = [{
                  value: data.weight_below_high_above,
                  name: '查询结果'
                }, {
                  value: other,
                  name: '其他'
                }]
                this.result_bar.xAxis.data = ['weight<' + this.options.weight + ',height>' + this.options.height, '其他', '总数据量']
                this.result_bar.series[0].data = [data.weight_below_high_above, other, data.data_length]
                this.fetch()
              } else {
                const other = data.data_length - data.weight_below_high_below
                this.result_pie.series[0].data = [{
                  value: data.weight_below_high_below,
                  name: '查询结果'
                }, {
                  value: other,
                  name: '其他'
                }]
                this.result_bar.xAxis.data = ['weight<' + this.options.weight + ',height<' + this.options.height, '其他', '总数据量']
                this.result_bar.series[0].data = [data.weight_below_high_below, other, data.data_length]
                this.fetch()
              }
            } else if (this.options.weight_value === '') {
              // height高于
              if (this.options.height_value === 'above') {
                this.result_pie.series[0].data = [{
                  value: data.high_above,
                  name: '查询结果'
                }, {
                  value: data.high_below,
                  name: '其他'
                }]
                this.result_bar.xAxis.data = ['height>' + this.options.height, '其他', '总数据量']
                this.result_bar.series[0].data = [data.high_above, data.high_below, data.data_length]
                this.fetch()
              } else {
                this.result_pie.series[0].data = [{
                  value: data.high_below,
                  name: '查询结果'
                }, {
                  value: data.high_above,
                  name: '其他'
                }]
                this.result_bar.xAxis.data = ['height<' + this.options.height, '其他', '总数据量']
                this.result_bar.series[0].data = [data.high_below, data.high_above, data.data_length]
                this.fetch()
              }
            }
          })
          .catch((error) => {
            console.log(error)
          })
      }
    },
    fetch () {
      // eslint-disable-next-line camelcase
      const chartPie = echarts.init(document.getElementById('Chart_pie'), 'westeros')
      // eslint-disable-next-line camelcase
      const chartBar = echarts.init(document.getElementById('Chart_bar'), 'westeros')
      chartPie.setOption(this.result_pie)
      chartBar.setOption(this.result_bar)
    }
  }
}
</script>

<style scoped lang="less">

.bar {
  float: left;
  margin-left: 10%;
  margin-top: 5%;
}

.pie {
  float: left;
  margin-left: 15%;
  margin-top: 5%;
}

.choice {
  display: flex;
}

.choice .el-input {
  width: 130px;
  margin: 15px;
}

.choice .el-button {
  height: 40px;
  margin: auto;
}

.text {
  font-size: 14px;
}

.item {
  padding: 18px 0;
}

.box-card {
  position: absolute;
  width: 40%;
}
</style>
