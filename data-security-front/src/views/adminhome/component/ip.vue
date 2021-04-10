<template>
  <div class="ip-container">
    <el-row type="flex" class="row-bg" :gutter="20">
      <el-col :span="12">
        <el-input v-model="freezeIpField"></el-input>
      </el-col>
      <el-col :span="2">
        <el-button type="primary" icon="el-icon-edit" @click="onFreezeIp">增加</el-button>
      </el-col>
    </el-row>
    <el-table
      :data="tableData"
      style="width: 100%;height: 100%"
      height="80%">
      <el-table-column
        v-for="(item, key) in titleData"
        :key="key"
        :label="item.name"
        :prop="item.value"></el-table-column>
      <el-table-column label="操作">
        <template slot-scope="scope">
          <el-button
            size="mini"
            type="danger"
            @click="onFreeIp(scope.$index, scope.row)">解封
          </el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
import { getAllIp, freeIp, freezeIp } from '@/api/ip'

export default {
  name: 'AdminIpIndex',
  components: {},
  props: {},
  data () {
    return {
      tableData: [],
      freezeIpField: '',
      titleData: [
        {
          name: 'id',
          value: 'id'
        }, {
          name: 'ip',
          value: 'ip'
        }]
    }
  },
  computed: {},
  watch: {},
  created () {
  },
  mounted () {
    this.onGetAllIp()
  },
  methods: {
    onGetAllIp () {
      getAllIp()
        .then(res => {
          if (res.data.payload === 'null') {
            this.$message.warning('已经没有更多了')
          } else {
            this.tableData = res.data.payload
          }
        })
    },
    onFreeIp (index, row) {
      freeIp(row.id)
        .then(res => {
          if (res.data.msg === 'success') {
            this.$message.success(res.data.payload)
            this.tableData.splice(index, 1)
          } else {
            this.$message.error(res.data.payload)
          }
        })
    },
    onFreezeIp () {
      if (this.freezeIpField !== '') {
        freezeIp({
          ip: this.freezeIpField
        })
          .then(res => {
            if (res.data.msg === 'success') {
              this.$message.success(res.data.payload)
              this.freezeIpField = ''
              this.onGetAllIp()
            } else {
              this.$message.error(res.data.payload)
            }
          })
      } else {
        this.$message.warning('不可为空，请填写')
      }
    }
  }
}
</script>

<style scoped lang="less">
.ip-container {
  width: 100%;
  height: 100%;
  overflow: hidden;
}
</style>
