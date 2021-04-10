<template>
  <div class="visualization-container">
    <div class="choice">
      <el-form v-model="options">
        <el-select v-model="options.salary_value" placeholder="请选择salary值高于或低于"
                   v-show="(options.weight_value === '')&& (options.height_value === '')&&(options.weight === '')&&(options.height === '')">
          <el-option
            v-for="item in options.salary_options"
            :key="item.value"
            :label="item.label"
            :value="item.value"
            :disabled="item.disabled">
          </el-option>
        </el-select>
        <el-input v-model="options.salary" placeholder="请输入salary值"
                  v-show="(options.weight_value === '')&&(options.height_value === '')&&(options.weight === '')&&(options.height === '')"></el-input>
        <el-select v-model="options.weight_value" placeholder="请选择weight值高于或低于"
                   v-show="(options.salary_value === '')&&(options.salary === '')">
          <el-option
            v-for="item in options.weight_options"
            :key="item.value"
            :label="item.label"
            :value="item.value"
            :disabled="item.disabled">
          </el-option>
        </el-select>
        <el-input v-model="options.weight" placeholder="请输入weight值"
                  v-show="(options.salary_value === '')&&(options.salary === '')"></el-input>
        <el-select v-model="options.height_value" placeholder="请选择height值高于或低于"
                   v-show="(options.salary_value === '')&&(options.salary === '')">
          <el-option
            v-for="item in options.height_options"
            :key="item.value"
            :label="item.label"
            :value="item.value"
            :disabled="item.disabled">
          </el-option>
        </el-select>
        <el-input v-model="options.height" placeholder="请输入height值"
                  v-show="(options.salary_value === '')&&(options.salary === '')"></el-input>
      </el-form>
      <el-button type="primary" @click="submit_choice" style="position: relative; margin-left: 20px">提交</el-button>
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
import { getAnalysis } from '@/api/staff'

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
        salary_value: '',
        weight_value: '',
        height_value: '',
        salary: '',
        weight: '',
        height: '',
        salary_options: [{
          value: 'above',
          label: 'salary高于'
        }, {
          value: 'below',
          label: 'salary低于'
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
      if ((this.options.salary_value === '') && (this.options.weight_value === '') && (this.options.height_value === '')) {
        this.$message.error('请选择一项选择框的内容')
      } else {
        getAnalysis({
          weight: parseInt(this.options.weight),
          height: parseInt(this.options.height),
          salary: parseInt(this.options.salary)
        })
          .then(res => {
            const data = res.data.payload
            if (this.options.salary !== '') {
              // salary高于
              if (this.options.salary_value === 'above') {
                this.result_pie.series[0].data = [{
                  value: data.salary_above,
                  name: '查询结果'
                }, {
                  value: data.salary_below,
                  name: '其他'
                }]
                this.result_bar.xAxis.data = ['salary>' + this.options.salary, '其他', '总数据量']
                this.result_bar.series[0].data = [data.salary_above, data.salary_below, data.data_length]
                this.fetch()
              } else if (this.options.salary_value === 'below') {
                this.result_pie.series[0].data = [{
                  value: data.salary_below,
                  name: '查询结果'
                }, {
                  value: data.salary_above,
                  name: '其他'
                }]
                this.result_bar.xAxis.data = ['salary<' + this.options.salary, '其他', '总数据量']
                this.result_bar.series[0].data = [data.salary_below, data.salary_above, data.data_length]
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
                const other = data.data_length - data.weight_above_height_above
                this.result_pie.series[0].data = [{
                  value: data.weight_above_height_above,
                  name: '查询结果'
                }, {
                  value: other,
                  name: '其他'
                }]
                this.result_bar.xAxis.data = ['weight>' + this.options.weight + ',height>' + this.options.height, '其他', '总数据量']
                this.result_bar.series[0].data = [data.weight_above_height_above, other, data.data_length]
                this.fetch()
              } else {
                const other = data.data_length - data.weight_above_height_below
                this.result_pie.series[0].data = [{
                  value: data.weight_above_height_below,
                  name: '查询结果'
                }, {
                  value: other,
                  name: '其他'
                }]
                this.result_bar.xAxis.data = ['weight>' + this.options.weight + ',height<' + this.options.height, '其他', '总数据量']
                this.result_bar.series[0].data = [data.weight_above_height_below, other, data.data_length]
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
                const other = data.data_length - data.weight_below_height_above
                this.result_pie.series[0].data = [{
                  value: data.weight_below_height_above,
                  name: '查询结果'
                }, {
                  value: other,
                  name: '其他'
                }]
                this.result_bar.xAxis.data = ['weight<' + this.options.weight + ',height>' + this.options.height, '其他', '总数据量']
                this.result_bar.series[0].data = [data.weight_below_height_above, other, data.data_length]
                this.fetch()
              } else {
                const other = data.data_length - data.weight_below_height_below
                this.result_pie.series[0].data = [{
                  value: data.weight_below_height_below,
                  name: '查询结果'
                }, {
                  value: other,
                  name: '其他'
                }]
                this.result_bar.xAxis.data = ['weight<' + this.options.weight + ',height<' + this.options.height, '其他', '总数据量']
                this.result_bar.series[0].data = [data.weight_below_height_below, other, data.data_length]
                this.fetch()
              }
            } else if (this.options.weight_value === '') {
              // height高于
              if (this.options.height_value === 'above') {
                this.result_pie.series[0].data = [{
                  value: data.height_above,
                  name: '查询结果'
                }, {
                  value: data.height_below,
                  name: '其他'
                }]
                this.result_bar.xAxis.data = ['height>' + this.options.height, '其他', '总数据量']
                this.result_bar.series[0].data = [data.height_above, data.height_below, data.data_length]
                this.fetch()
              } else {
                this.result_pie.series[0].data = [{
                  value: data.height_below,
                  name: '查询结果'
                }, {
                  value: data.height_above,
                  name: '其他'
                }]
                this.result_bar.xAxis.data = ['height<' + this.options.height, '其他', '总数据量']
                this.result_bar.series[0].data = [data.height_below, data.height_above, data.data_length]
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
      const chartPie = echarts.init(document.getElementById('Chart_pie'), 'westeros')
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
