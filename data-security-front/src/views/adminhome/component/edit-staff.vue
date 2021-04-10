<template>
  <div class="edit-staff-container">
    <el-form ref="editStaffForm" :model="editStaffForm" :rules="editStaffRules" label-width="80px">
      <el-form-item label="姓名" prop="staff_name">
        <el-input placeholder="请输入姓名" v-model="editStaffForm.staff_name" auto-complete="off"/>
      </el-form-item>
      <el-form-item label="学历" prop="qualification">
        <el-select v-model="editStaffForm.qualification" placeholder="请选择">
          <el-option
            v-for="item in options"
            :key="item.value"
            :label="item.label"
            :value="item.value">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="体重" prop="weight">
        <el-input placeholder="请输入体重" v-model="editStaffForm.weight" auto-complete="off"/>
      </el-form-item>
      <el-form-item label="身高" prop="height">
        <el-input placeholder="请输入身高" v-model="editStaffForm.height" auto-complete="off"/>
      </el-form-item>
      <el-form-item label="身份证号" prop="id_card">
        <el-input placeholder="请输入身份证" v-model="editStaffForm.id_card" auto-complete="off"/>
      </el-form-item>
      <el-form-item label="薪资" prop="salary">
        <el-input placeholder="请输入薪资" v-model="editStaffForm.salary" auto-complete="off"/>
      </el-form-item>
      <el-form-item>
        <el-button type="primary"
                   @click="onEditStaff('editStaffForm')">提交
        </el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
import { updateStaff } from '@/api/staff'

export default {
  name: 'EditStaffIndex',
  components: {},
  props: ['deliver'],
  data () {
    return {
      editStaffForm: {
        staff_id: this.deliver.staff_id,
        staff_name: this.deliver.staff_name,
        height: this.deliver.height,
        weight: this.deliver.weight,
        qualification: this.deliver.qualification,
        id_card: this.deliver.id_card,
        salary: this.deliver.salary
      },
      options: [{
        label: '本科',
        value: '本科'
      }, {
        label: '硕士',
        value: '硕士'
      }, {
        label: '博士',
        value: '博士'
      }],
      editStaffRules: {
        staff_name: [
          {
            required: true,
            message: '姓名不可为空',
            trigger: 'blur'
          }
        ],
        qualification: [{
          required: true,
          message: '学历不可为空',
          trigger: 'blur'
        }],
        height: [{
          required: true,
          message: '身高不可为空',
          trigger: 'blur'
        }],
        weight: [{
          required: true,
          message: '体重不可为空',
          trigger: 'blur'
        }],
        id_card: [{
          required: true,
          message: '身份证号不可为空',
          trigger: 'blur'
        }],
        salary: [{
          required: true,
          message: '薪资不可为空',
          trigger: 'blur'
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
    onEditStaff (formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          updateStaff(this.editStaffForm, this.editStaffForm.staff_id)
            .then(res => {
              if (res.data.payload === '修改成功') {
                this.$message.success(res.data.payload)
                this.$emit('closeEditDialog')
              } else {
                this.$message.error(res.data.payload)
              }
            })
        } else {
          return false
        }
      })
    }
  }
}
</script>

<style scoped lang="less"></style>
