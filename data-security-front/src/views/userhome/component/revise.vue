<template>
  <el-form
    ref="reviseForm"
    :model="reviseForm"
    :rules="rules"
    label-width="80px"
    class="revise-form">
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
</template>

<script>
import { getRsa, getHash } from '@/util/tools'
import { revise } from '@/api/user'

export default {
  name: 'ReviseIndex',
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
        } else {
          return false
        }
      })
    }
  }
}
</script>

<style scoped lang="less">
.revise-form {
  position: relative;
  padding: 30px;
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
