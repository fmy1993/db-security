<template>
  <div class="user-revise-container">
    <div class="revise-head">
      <el-button type="danger" @click="onLogout">退出</el-button>
    </div>
    <el-form
      ref="reviseForm"
      :model="reviseForm"
      :rules="rules"
      label-width="80px"
      class="revise-form">
      <span class="revise-title">修改密码</span>
      <el-form-item
        label="旧密码"
        prop="oldPwd">
        <el-input
          type="password"
          placeholder="请输入原始密码"
          v-model="reviseForm.oldPwd"
          auto-complete="off"/>
      </el-form-item>
      <el-form-item
        label="新密码"
        prop="newPwd">
        <el-input
          type="password"
          placeholder="请输入新密码"
          v-model="reviseForm.newPwd"
          auto-complete="off"/>
      </el-form-item>
      <el-form-item
        label="重复密码"
        prop="checkNewPwd">
        <el-input
          type="password"
          placeholder="请重复新密码"
          v-model="reviseForm.checkNewPwd"
          auto-complete="off"/>
      </el-form-item>
      <el-form-item class="revise-form-button">
        <el-button
          type="primary"
          @click="reviseSub('reviseForm')">提交
        </el-button>
        <el-button
          @click="resetForm('reviseForm')">重置
        </el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
import { getRsa, getHash } from '@/util/tools'
import { logout, revise } from '@/api/user'

export default {
  name: 'UserReviseIndex',
  components: {},
  props: {},
  data () {
    const validatePwd = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请再次输入密码'))
      } else if (value !== this.reviseForm.newPwd) {
        callback(new Error('两次输入密码不一致!'))
      } else {
        callback()
      }
    }
    return {
      reviseForm: {
        oldPwd: '',
        newPwd: '',
        checkNewPwd: ''
      },
      rules: {
        oldPwd: [
          {
            required: true,
            message: '请输入原密码',
            trigger: 'blur'
          }
        ],
        newPwd: [
          {
            required: true,
            message: '请输入新密码',
            trigger: 'blur'
          }
        ],
        checkNewPwd: [
          {
            required: true,
            message: '请重复新密码',
            trigger: 'blur'
          },
          {
            validator: validatePwd,
            trigger: 'blur'
          }
        ]
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
    resetForm (formName) {
      this.$refs[formName].resetFields()
    },
    reviseSub (formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          const publicKey = '-----BEGIN PUBLIC KEY-----\n' +
            'MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDOikBHakQyUwIVjiMsFwVVdvFX\n' +
            'PcZD3ph8Fcv4Ftl0Jhe5/McNB+w34vAui1k8d28J+42xuCTb8kKteGWOdXl9aPWX\n' +
            'ZE07+6fnZO3u/T6AC8WFU80haTBz6a0rDWnHPkYrkazRJkqUBjcNMlUSttWNYJfZ\n' +
            'SqyV0/rKHa6fK0dDTwIDAQAB\n' +
            '-----END PUBLIC KEY-----'
          revise({
            old_password: getRsa(getHash(this.reviseForm.oldPwd), publicKey),
            new_password: getRsa(getHash(this.reviseForm.newPwd), publicKey)
          })
            .then((response) => {
              if (response.data.payload === '修改成功') {
                this.$message.success(response.data.payload)
                localStorage.clear()
                this.$router.push('/login')
              } else {
                this.$message.error(response.data.payload)
              }
            })
            .catch((error) => {
              console.log(error)
            })
        } else {
          return false
        }
      })
    },
    onLogout () {
      this.$confirm('退出后再登录需要重新输入密码, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
        .then(() => {
          logout()
            .then(response => {
              if (response.data.payload === '退出成功') {
                this.$message({
                  type: 'success',
                  message: '退出成功!'
                })
                localStorage.clear()
              } else {
                this.$message({
                  type: 'error',
                  message: '你还没登录'
                })
                localStorage.clear()
              }
            })
            .catch(error => {
              console.log(error)
            })
          this.$router.push('/login')
        })
        .catch(() => {
          this.$message({
            type: 'info',
            message: '已取消退出'
          })
          // this.$router.push('/user-home/revise')
        })
    }
  }
}
</script>

<style scoped lang="less">
.user-revise-container {
  height: 100%;
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;

  .revise-title {
    text-align: center;
    font-size: 40px;
    display: flex;
    justify-content: center;
    margin-top: 30px;
    margin-bottom: 30px;
  }
}

.revise-head {
  position: fixed;
  top: 30px;
  right: 30px;

  .el-button {
    padding-left: 30px;
    padding-right: 30px;
  }
}

.revise-form {
  position: fixed;
  width: 800px;
  height: 400px;
  padding: 30px;
  border: 1px solid #DCDFE6;
  border-radius: 5px;
  box-shadow: 0 0 2px #b1b4ba;
}

.revise-form-button {
  margin-bottom: 0;
  margin-top: 40px;
  display: flex;
  justify-content: center;
  align-items: center;
  margin-left: -80px;
}

.el-form-item {
  margin-bottom: 30px;

  .el-button {
    padding-left: 30px;
    padding-right: 30px;
    margin-left: 40px;
    margin-right: 40px;
  }
}
</style>
