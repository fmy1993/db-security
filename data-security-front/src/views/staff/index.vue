<template>
  <div class="staff-container">
    <div class="staff-container-search">
      <el-row :gutter="20">
        <el-col :span="12">
          <el-input
            @keyup.enter.native="onStaffInfo"
            placeholder="请输入内容"
            prefix-icon="el-icon-search"
            v-model="qualification">
          </el-input>
        </el-col>
        <el-col :span="2">
          <el-button type="primary" icon="el-icon-search" @click="onStaffInfo">搜索</el-button>
        </el-col>
        <el-col :span="2">
          <el-tooltip class="item" effect="dark" content="查询哪个身份的人最多" placement="top-end">
            <el-button type="primary" @click="onGetMostIdCard">查询</el-button>
          </el-tooltip>
        </el-col>
        <el-col :span="2"
                style="height: 38px; display: flex;align-items: center;justify-content: center">
          <span style="font-size: 20px">{{ mostProvince }}</span>
        </el-col>
      </el-row>
    </div>
    <el-table :data="tableData" class="staff-table" height="80%">
      <el-table-column type="index"></el-table-column>
      <el-table-column v-for="(item, key) in titles" :key="key" :label="item.name" :prop="item.value">
      </el-table-column>
    </el-table>
    <el-pagination
      class="staff-table-pagination"
      background
      layout="prev, pager, next, jumper"
      :current-page.sync="currentPage"
      :total="total"
      :page-size="50"
      @current-change="onStaffInfo"></el-pagination>
  </div>
</template>

<script>
import { getStaff, getMostIdcard } from '@/api/staff'

export default {
  name: 'StaffIndex',
  components: {},
  props: {},
  data () {
    return {
      tableData: [],
      qualification: '',
      currentPage: 1,
      total: 0,
      titles: [{
        name: '身高',
        value: 'height'
      }, {
        name: '体重',
        value: 'weight'
      }, {
        name: '学历',
        value: 'qualification'
      }, {
        name: '薪资',
        value: 'salary'
      }],
      mostProvince: ''
    }
  },
  computed: {},
  watch: {},
  created () {
  },
  mounted () {
    this.onStaffInfo()
  },
  methods: {
    onStaffInfo () {
      getStaff({
        page: this.currentPage,
        qualification: this.qualification
      })
        .then(res => {
          this.tableData = res.data.payload.staff
          this.total = res.data.payload.total
        })
    },
    onGetMostIdCard () {
      getMostIdcard()
        .then(res => {
          this.mostProvince = res.data.payload
        })
    }
  }
}
</script>

<style scoped lang="less">
.staff-container {
  color: black;
  height: 100%;
}

.staff-table-pagination {
  margin-top: 50px;
  display: flex;
  justify-content: center;
}
</style>
