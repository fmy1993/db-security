<template>
  <div class="admin-ip-container">
    <el-row type="flex" class="row-bg">
      <el-col :span="12">
        <el-input v-model="addField"></el-input>
      </el-col>
      <el-col :span="2">
        <el-button type="primary" icon="el-icon-edit" @click="addIp">增加</el-button>
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
            @click="handleFree(scope.$index, scope.row)">解封
          </el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
import { adminAddIp, adminFreeIp, adminGetAllIps } from '@/api/admin'
import { refreshToken } from '@/util/tools'

export default {
  name: 'AdminIpIndex',
  components: {},
  props: {},
  data () {
    return {
      tableData: [],
      addField: '',
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
    this.onGetAllIps()
  },
  methods: {
    onGetAllIps () {
      adminGetAllIps()
        .then((response) => {
          refreshToken(response)
          if (response.data.payload === 'null') {
            this.$message.warning('已经没有更多了')
          } else {
            this.tableData = response.data.payload
          }
        })
        .catch((error) => {
          console.log(error)
        })
    },
    handleFree (index, row) {
      adminFreeIp(row)
        .then((response) => {
          if (response.data.payload === '解封成功') {
            this.$message.success(response.data.payload)
            this.tableData.splice(index, 1)
          } else {
            this.$message.error(response.data.payload)
          }
        })
        .catch((error) => {
          console.log(error)
        })
    },
    addIp () {
      if (this.addField !== '') {
        adminAddIp(this.addField)
          .then((response) => {
            refreshToken(response)
            if (response.data.payload === '封禁成功') {
              this.$message.success(response.data.payload)
              this.addField = ''
            } else {
              this.$message.error(response.data.payload)
            }
          })
          .catch((error) => {
            console.log(error)
          })
      } else {
        this.$message.warning('请填写')
      }
    }
  }
}
</script>

<style scoped lang="less">
.admin-ip-container {
  width: 100%;
  height: 100%;
  overflow: hidden;
}
</style>
