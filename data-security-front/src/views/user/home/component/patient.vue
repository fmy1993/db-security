<template>
  <div class="patient-container">
    <el-row type="flex" class="row-bg">
      <el-col :span="12">
        <el-input v-model="searchField"></el-input>
      </el-col>
      <el-col :span="2">
        <el-button type="primary" icon="el-icon-search" @click="onPatientInfo">搜索</el-button>
      </el-col>
    </el-row>
    <el-table
      :data="tableData"
      height="80%"
      class="patient-table">
      <el-table-column
        v-for="(item, key) in titleData"
        :key="key"
        :label="item.name"
        :prop="item.value"></el-table-column>
    </el-table>

    <el-pagination
      background
      class="patient-pagination"
      layout="prev, pager, next, jumper"
      :current-page="currentPage"
      :total="total"
      page-size="50"
      @current-change="handleCurrentChange"></el-pagination>
  </div>
</template>

<script>
import { patient, search } from '@/api/patient'
import { refreshToken } from '@/util/tools'

export default {
  name: 'UserHomePatientIndex',
  components: {},
  props: {},
  data () {
    return {
      tableData: [],
      searchField: '',
      currentPage: 1,
      total: 30000,
      titleData: [{
        name: '病号',
        value: 'id'
      }, {
        name: '手机',
        value: 'phone'
      }, {
        name: '体重',
        value: 'weight'
      }, {
        name: '身高',
        value: 'high'
      }, {
        name: '年龄',
        value: 'age'
      }, {
        name: '地址',
        value: 'address'
      }, {
        name: '账单',
        value: 'bill'
      }]
    }
  },
  computed: {},
  watch: {},
  created () {
  },
  mounted () {
    if (typeof localStorage.getItem('total') === 'undefined') {
      this.total = 30000
    } else {
      this.total = parseInt(localStorage.getItem('total'))
    }
    this.onPatientInfo()
  },
  methods: {
    onPatientInfo () {
      if (this.searchField === '') {
        patient(this.currentPage)
          .then(response => {
            refreshToken(response)
            if (response.data.payload === 'null') {
              this.$message.warning('已经没有更多了')
            } else {
              this.tableData = response.data.payload
            }
          })
          .catch(error => {
            console.log(error)
          })
      } else {
        search(this.currentPage, this.searchField)
          .then(response => {
            refreshToken(response)
            if (response.data.payload === 'null') {
              this.$message.warning('已经没有更多了')
            } else {
              this.tableData = response.data.payload
            }
          })
          .catch(error => {
            console.log(error)
          })
      }
    },
    handleCurrentChange (val) {
      this.currentPage = val
      this.onPatientInfo()
    }
  }
}
</script>

<style scoped lang="less">
.patient-container {
  width: 100%;
  height: 100%;
  overflow: hidden;

  .patient-table {
    margin-top: 50px;
  }

  .patient-pagination {
    margin-top: 20px;
    margin-left: 100px;
  }
}
</style>
