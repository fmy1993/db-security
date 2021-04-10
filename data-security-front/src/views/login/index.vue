<template>
  <div class="login-container">
    <el-form ref="loginForm" :model="loginForm" :rules="loginFormRules" label-width="80px" class="login-form">
      <span class="login-title">欢迎登录</span>
      <el-form-item label="账号" prop="phone">
        <el-input type="text" placeholder="请输入账号" v-model="loginForm.phone" auto-complete="off"/>
      </el-form-item>
      <el-form-item label="密码" prop="password">
        <el-input type="password" placeholder="请输入密码" v-model="loginForm.password" auto-complete="off"/>
      </el-form-item>
      <el-form-item label="验证码" prop="captcha" @keyup.enter.native="onSubmit('loginForm')">
        <el-row>
          <el-col :span="12">
            <el-input type="text" placeholder="请输入验证码" v-model="loginForm.captcha" auto-complete="off"/>
          </el-col>
          <el-col :span="12" style="height: 40px;">
            <img :src=base64Img alt="点击换一张" id="id_captcha" @click="onCaptcha">
          </el-col>
        </el-row>
      </el-form-item>
      <el-form-item class="login-form-button">
        <el-button type="primary" @click="onSubmit('loginForm')"
                   v-loading.fullscreen.lock="fullscreenLoading">提交
        </el-button>
        <el-button @click="resetForm('loginForm')">重置</el-button>
        <el-button type="primary" @click="$router.push('/register')">注册</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
import { getSalt, login, getCaptcha, getUserInfo } from '@/api/user'
import { getRsa, getHash } from '@/util/tools'

export default {
  name: 'LoginIndex',
  components: {},
  props: {},
  data () {
    return {
      salt: '',
      base64Img: '',
      captchaId: '',
      loginForm: {
        phone: '17356581383',
        password: '123',
        captcha: '1233'
      },
      fullscreenLoading: false,
      csrfToken: '',
      loginFormRules: {
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
        captcha: [
          {
            required: true,
            message: '验证码不可为空',
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
    this.onCaptcha()
  },
  methods: {
    async onSalt () {
      await getSalt()
        .then(response => {
          this.salt = response.data.payload
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
          this.openFullScreen()
          login({
            phone: this.loginForm.phone,
            password: getRsa(this.salt + getHash(this.loginForm.password), publicKey),
            captchaId: this.captchaId,
            captchaValue: this.loginForm.captcha
          })
            .then((response) => {
              if (response.data.msg !== 'success') {
                this.$message.error(response.data.payload)
                this.fullscreenLoading = false
                this.onCaptcha()
              } else {
                localStorage.setItem('Authorization', response.data.payload.token)
                this.fullscreenLoading = false
                this.$message({
                  message: '恭喜你，登录成功',
                  type: 'success'
                })
                this.onGetUserInfo()
              }
            })
        } else {
          return false
        }
      })
    },
    onCaptcha () {
      getCaptcha()
        .then(response => {
          this.base64Img = response.data.payload.data
          this.captchaId = response.data.payload.captchaId
        })
    },
    resetForm (formName) {
      this.$refs[formName].resetFields()
    },
    openFullScreen () {
      this.fullscreenLoading = true
      setTimeout(() => {
        this.fullscreenLoading = false
      }, 20000)
    },
    async onGetUserInfo () {
      await getUserInfo()
        .then(res => {
          localStorage.setItem('user', JSON.stringify(res.data.payload))
          if (res.data.payload.is_super_user === 1) {
            this.$router.push('/admin-home')
          } else {
            this.$router.push('/home')
          }
        })
    }
  }
}
</script>

<style scoped lang="less">
.login-container {
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

.login-form {
  position: fixed;
  width: 600px;
  height: 450px;
  padding: 30px;
  border: 1px solid #DCDFE6;
  border-radius: 5px;
  box-shadow: 0 0 2px #b1b4ba;

  .login-title {
    text-align: center;
    font-size: 40px;
    display: flex;
    justify-content: center;
    margin-top: 30px;
    margin-bottom: 50px;
  }

  .login-form-button {
    margin-bottom: 0;
    margin-top: 60px;
    display: flex;
    justify-content: center;
    align-items: center;
    margin-left: -80px;
  }
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
