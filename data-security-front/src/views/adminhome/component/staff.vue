<template>
  <div class="staff-container">
    <div class="staff-container-search">
      <el-row :gutter="20" v-loading="dpLoading">
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
          <el-button type="primary" @click="onAddStaff">增加</el-button>
        </el-col>
        <el-col :span="2">
          <el-button type="primary" @click="onDifferentialPrivacy">添加噪音</el-button>
        </el-col>
        <el-dialog title="添加员工"
                   :visible.sync="addDialogVisible"
                   append-to-body
                   :before-close="handleClose">
          <add-staff @closeAddDialog="addDialogVisible = false"></add-staff>
        </el-dialog>
        <el-dialog title="编辑员工"
                   :visible.sync="editDialogVisible"
                   append-to-body
                   :before-close="handleClose">
          <edit-staff @closeEditDialog="editDialogVisible = false" :deliver="editDeliverData"></edit-staff>
        </el-dialog>
      </el-row>
    </div>
    <el-table :data="tableData" class="staff-table" height="80%">
      <el-table-column type="index"></el-table-column>
      <el-table-column v-for="(item, key) in titles" :key="key" :label="item.name" :prop="item.value">
      </el-table-column>
      <el-table-column label="操作">
        <template slot-scope="scope">
          <el-button
            size="mini"
            @click="onEdit(scope.$index, scope.row)">编辑
          </el-button>
          <el-button
            size="mini"
            type="danger"
            @click="onDelete(scope.$index, scope.row)">删除
          </el-button>
        </template>
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
import { getOriStaff, deleteStaff, differentialPrivacy } from '@/api/staff'
import AddStaff from '@/views/adminhome/component/add-staff'
import EditStaff from '@/views/adminhome/component/edit-staff'

export default {
  name: 'StaffIndex',
  components: {
    AddStaff,
    EditStaff
  },
  props: {},
  data () {
    return {
      tableData: [],
      qualification: '',
      currentPage: 1,
      total: 0,
      titles: [{
        name: 'id',
        value: 'staff_id'
      }, {
        name: '姓名',
        value: 'staff_name'
      }, {
        name: '身份证号',
        value: 'id_card'
      }, {
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
      addDialogVisible: false,
      editDialogVisible: false,
      editDeliverData: {},
      dpLoading: false
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
      getOriStaff({
        page: this.currentPage,
        qualification: this.qualification
      })
        .then(res => {
          this.tableData = res.data.payload.staffData
          this.total = res.data.payload.total
        })
    },
    onAddStaff () {
      this.addDialogVisible = true
    },
    handleClose () {
      this.addDialogVisible = false
      this.editDialogVisible = false
    },
    onEdit (index, row) {
      this.editDeliverData = row
      this.editDialogVisible = true
    },
    onDelete (index, row) {
      this.$confirm('此操作将永久删除该条记录, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deleteStaff(row.staff_id)
          .then(res => {
            if (res.data.payload === '删除成功') {
              this.$message.success('删除成功')
              this.tableData.splice(index, 1)
            } else {
              this.$message.error('删除失败')
            }
          })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消删除'
        })
      })
    },
    onDifferentialPrivacy () {
      this.dpLoading = true
      differentialPrivacy()
        .then(res => {
          if (res.data.msg === 'success') {
            this.$message.success('添加噪音成功')
            this.dpLoading = false
          } else {
            this.$message.error(res.data.payload)
            this.dpLoading = false
          }
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
