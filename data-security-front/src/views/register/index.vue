<template>
  <div class="register-container">
    <el-form ref="registerForm" :model="registerForm" :rules="registerFormRules" label-width="80px"
             class="register-form">
      <span class="register-title">欢迎注册</span>
      <el-form-item label="账号" prop="phone">
        <el-input type="text" placeholder="请输入账号" v-model="registerForm.phone" auto-complete="off"/>
      </el-form-item>
      <el-form-item label="密码" prop="password">
        <el-input type="password" placeholder="请输入密码" v-model="registerForm.password" auto-complete="off"/>
      </el-form-item>
      <el-form-item label="确认密码" prop="checkPwd">
        <el-input type="password" placeholder="请再次输入密码" v-model="registerForm.checkPwd" autocomplete="off"></el-input>
      </el-form-item>
      <el-form-item label="校验码" prop="checkCode">
        <el-input type="password" placeholder="请输入校验码" v-model="registerForm.checkCode" autocomplete="off"></el-input>
      </el-form-item>
      <el-form-item class="register-form-button">
        <el-button type="primary" @click="onSubmit('registerForm')" :disabled="submitLoading">提交
        </el-button>
        <el-button @click="resetForm('registerForm')">重置</el-button>
        <el-button type="primary" @click="$router.push('/login')">登录</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
import { getSalt, register } from '@/api/user'
import { getRsa, getHash } from '@/util/tools'

export default {
  name: 'RegisterIndex',
  components: {},
  props: {},
  data () {
    const validatePass = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请再次输入密码'))
      } else if (value !== this.registerForm.password) {
        callback(new Error('两次输入密码不一致!'))
      } else {
        callback()
      }
    }
    return {
      salt: '',
      submitLoading: false,
      registerForm: {
        phone: '',
        password: '',
        checkPwd: '',
        checkCode: ''
      },
      registerFormRules: {
        phone: [
          {
            required: true,
            message: '账号不可为空',
            trigger: 'blur'
          },
          {
            min: 11,
            max: 11,
            message: '长度必须为11个字符',
            trigger: 'blur'
          }
        ],
        password: [
          {
            required: true,
            message: '密码不可为空',
            trigger: 'blur'
          }
        ],
        checkPwd: [
          {
            required: true,
            message: '密码不可为空',
            trigger: 'blur'
          },
          {
            validator: validatePass,
            trigger: 'blur'
          }
        ],
        checkCode: [
          {
            required: true,
            message: '请输入校验码',
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
    async onSalt () {
      await getSalt()
        .then(response => {
          this.salt = response.data.payload
        })
        .catch(err => {
          console.log(err)
        })
    },
    async onSubmit (formName) {
      this.$refs[formName].validate(async (valid) => {
        if (valid) {
          await this.onSalt()
          const publicKey = '-----BEGIN PUBLIC KEY-----\n' +
            'MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDOikBHakQyUwIVjiMsFwVVdvFX\n' +
            'PcZD3ph8Fcv4Ftl0Jhe5/McNB+w34vAui1k8d28J+42xuCTb8kKteGWOdXl9aPWX\n' +
            'ZE07+6fnZO3u/T6AC8WFU80haTBz6a0rDWnHPkYrkazRJkqUBjcNMlUSttWNYJfZ\n' +
            'SqyV0/rKHa6fK0dDTwIDAQAB\n' +
            '-----END PUBLIC KEY-----'
          register({
            phone: this.registerForm.phone,
            password: getRsa(this.salt + getHash(this.registerForm.password), publicKey),
            checkCode: getHash(this.registerForm.checkCode)
          })
            .then(response => {
              if (response.data.payload === '注册成功') {
                this.$message({
                  message: response.data.payload + ',快去登录吧',
                  type: 'success'
                })
                this.$router.push('/login')
              } else {
                this.$message.error(response.data.payload)
              }
            })
            .catch(error => {
              console.log(error)
            })
        }
      })
    },
    resetForm (formName) {
      this.$refs[formName].resetFields()
    }
  }
}
</script>

<style scoped lang="less">
.register-container {
  position: fixed;
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}

.register-form {
  position: fixed;
  width: 600px;
  height: 450px;
  padding: 30px;
  border: 1px solid #DCDFE6;
  border-radius: 5px;
  box-shadow: 0 0 2px #b1b4ba;

  .register-title {
    text-align: center;
    font-size: 40px;
    display: flex;
    justify-content: center;
    margin-top: 30px;
    margin-bottom: 40px;
  }

  .register-form-button {
    margin-bottom: 0px;
    margin-top: 29px;
    display: flex;
    justify-content: center;
    align-items: center;
    margin-left: -80px;
  }
}

.el-form-item {
  margin-bottom: 20px;

  .el-button {
    padding-left: 30px;
    padding-right: 30px;
    margin-left: 40px;
    margin-right: 40px;
  }
}
</style>
