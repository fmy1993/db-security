<template>
  <div class="admin-patient-container">
    <div>
      <el-row type="flex" class="row-bg">
        <el-col :span="12">
          <el-input v-model="searchField"></el-input>
        </el-col>
        <el-col :span="2">
          <el-button type="primary" icon="el-icon-search" @click="onPatientInfo">搜索</el-button>
        </el-col>
        <el-col :span="2">
          <el-button type="primary" icon="el-icon-edit" @click="dialogFormVisible = true">增加</el-button>
        </el-col>
        <el-col :span="2">
          <el-button type="primary" @click="differentialPrivacy" v-loading.fullscreen.lock="fullscreenLoading">差分隐私
          </el-button>
        </el-col>
      </el-row>
    </div>
    <el-table
      :data="tableData"
      height="80%"
      class="admin-patient-table">
      <el-table-column
        v-for="(item, key) in titleData"
        :key="key"
        :label="item.name"
        :prop="item.value"></el-table-column>
      <el-table-column label="操作">
        <template slot-scope="scope">
          <el-button
            size="mini"
            style="padding-left: 10px; padding-right: 10px"
            @click="handleEdit(scope.$index, scope.row)">编辑
          </el-button>
          <el-button
            size="mini"
            type="danger"
            style="padding-left: 10px; padding-right: 10px"
            @click="handleDelete(scope.$index, scope.row)">删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      background
      class="admin-patient-pagination"
      layout="prev, pager, next, jumper"
      :current-page="currentPage"
      :total="total"
      page-size="50"
      @current-change="handleCurrentChange"></el-pagination>

    <el-dialog title="填写信息" :visible.sync="dialogFormVisible" append-to-body>
      <el-form :model="addForm">
        <el-form-item label="姓名">
          <el-input v-model="addForm.name" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="体重">
          <el-input v-model="addForm.weight" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="身高">
          <el-input v-model="addForm.high" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="年龄">
          <el-input v-model="addForm.age" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="身份证">
          <el-input v-model="addForm.id_card" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="手机">
          <el-input v-model="addForm.phone" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="地址">
          <el-input v-model="addForm.address" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="账单">
          <el-input v-model="addForm.bill" autocomplete="off"></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false">取 消</el-button>
        <el-button type="primary" @click="onAddPatient">确 定</el-button>
      </div>
    </el-dialog>

    <el-dialog title="修改信息" :visible.sync="updateDialogFormVisible" append-to-body>
      <el-form :model="updateForm">
        <el-form-item label="姓名">
          <el-input v-model="updateForm.name" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="体重">
          <el-input v-model="updateForm.weight" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="身高">
          <el-input v-model="updateForm.high" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="年龄">
          <el-input v-model="updateForm.age" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="身份证">
          <el-input v-model="updateForm.id_card" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="手机">
          <el-input v-model="updateForm.phone" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="地址">
          <el-input v-model="updateForm.address" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="账单">
          <el-input v-model="updateForm.bill" autocomplete="off"></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="updateDialogFormVisible = false">取 消</el-button>
        <el-button type="primary" @click="onUpdatePatient">确 定</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import {
  adminPatient,
  adminSearch,
  adminAddPatient,
  adminUpdatePatient,
  adminDeletePatient,
  adminDp
} from '@/api/admin'
import { refreshToken } from '@/util/tools'

export default {
  name: 'AdminPatientIndex',
  components: {},
  props: {},
  data () {
    return {
      tableData: [],
      searchField: '',
      currentPage: 1,
      total: 30000,
      dialogFormVisible: false,
      updateDialogFormVisible: false,
      fullscreenLoading: false,
      addForm: {
        name: '',
        weight: 0,
        high: 0,
        age: 0,
        id_card: '',
        phone: '',
        address: '',
        bill: ''
      },
      updateForm: {
        name: '',
        weight: 0,
        high: 0,
        age: 0,
        id_card: '',
        phone: '',
        address: '',
        bill: ''
      },
      titleData: [
        {
          name: '病号',
          value: 'id'
        }, {
          name: '姓名',
          value: 'name'
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
        }, {
          name: '身份证号',
          value: 'id_card'
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
    handleCurrentChange (val) {
      this.currentPage = val
      this.onPatientInfo()
    },
    onPatientInfo () {
      if (this.searchField === '') {
        adminPatient(this.currentPage)
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
      } else {
        adminSearch(this.currentPage, this.searchField)
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
      }
    },
    onAddPatient () {
      this.addForm.age = parseInt(this.addForm.age)
      this.addForm.bill = parseFloat(this.addForm.bill)
      this.addForm.weight = parseFloat(this.addForm.weight)
      this.addForm.high = parseInt(this.addForm.high)
      adminAddPatient(this.addForm)
        .then((response) => {
          refreshToken(response)
          if (response.data.payload === '添加成功') {
            this.$message.success(response.data.payload)
          } else {
            this.$message.error(response.data.payload)
          }
        })
        .catch((error) => {
          console.log(error)
        })
      this.dialogFormVisible = false
    },
    handleEdit (index, row) {
      this.updateDialogFormVisible = true
      this.updateForm = row
    },
    onUpdatePatient () {
      adminUpdatePatient(this.updateForm)
        .then((response) => {
          refreshToken(response)
          if (response.data.payload === '修改成功') {
            this.$message.success(response.data.payload)
          } else {
            this.$message.error(response.data.payload)
          }
        })
        .catch((error) => {
          console.log(error)
        })
      this.updateDialogFormVisible = false
    },
    handleDelete (index, row) {
      adminDeletePatient(row)
        .then((response) => {
          refreshToken(response)
          if (response.data.payload === '删除成功') {
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
    differentialPrivacy () {
      this.openFullScreen()
      adminDp()
        .then((response) => {
          if (response.data.payload === 'success') {
            this.$message.success(response.data.payload)
          } else {
            this.$message.error(response.data.payload)
          }
          this.fullscreenLoading = false
        })
        .catch((error) => {
          console.log(error)
        })
    },
    openFullScreen () {
      this.fullscreenLoading = true
      setTimeout(() => {
        this.fullscreenLoading = false
      }, 40000)
    }
  }
}
</script>

<style scoped lang="less">
.admin-patient-container {
  width: 100%;
  height: 100%;
  overflow: hidden;

  .admin-patient-table {
    margin-top: 50px;
  }

  .admin-patient-pagination {
    margin-top: 20px;
    margin-left: 100px;
  }
}
</style>
