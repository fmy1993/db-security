<template>
  <div class="admin-user-container">
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
            @click="handleBan(scope.$index, scope.row)">封禁
          </el-button>
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
import { getAllUsers, freezeUser, freeUser } from '@/api/user'

export default {
  name: 'AdminUserIndex',
  components: {},
  props: {},
  data () {
    return {
      tableData: [],
      titleData: [
        {
          name: '用户号',
          value: 'id'
        }, {
          name: '手机号',
          value: 'phone'
        }, {
          name: '密码',
          value: 'password'
        }, {
          name: '指纹',
          value: 'finger_print'
        }, {
          name: '管理员',
          value: 'is_super_user'
        }, {
          name: '加入日期',
          value: 'date_joined'
        }]
    }
  },
  computed: {},
  watch: {},
  created () {
  },
  mounted () {
    this.onGetAllUsers()
  },
  methods: {
    onGetAllUsers () {
      getAllUsers()
        .then((response) => {
          if (response.data.payload === 'null') {
            this.$message.warning('已经没有更多了')
          } else {
            this.tableData = response.data.payload
          }
        })
    },
    handleBan (index, row) {
      freezeUser(row.id)
        .then((res) => {
          if (res.data.msg === 'success') {
            this.$message.success(res.data.payload)
          } else {
            this.$message.error(res.data.payload)
          }
        })
    },
    handleFree (index, row) {
      freeUser(row.id)
        .then((res) => {
          if (res.data.msg === 'success') {
            this.$message.success(res.data.payload)
          } else {
            this.$message.error(res.data.payload)
          }
        })
    }
  }
}
</script>

<style scoped lang="less">
.admin-user-container {
  width: 100%;
  height: 100%;
  overflow: hidden;
}
</style>
