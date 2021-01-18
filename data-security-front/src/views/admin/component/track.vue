<template>
  <el-container class="admin-track-container">
    <el-header class="admin-track-title">盗版追踪</el-header>
    <el-main>
      <div class="admin-track-upload">
        <el-upload drag :action="uploadUrl"
                   class="admin-track-upload-box"
                   :on-success="getResult" :before-upload="beforeFileUpload">
          <i class="el-icon-upload"></i>
          <div class="el-upload__text">将文件拖到此处，或<em>点击上传</em></div>
        </el-upload>
      </div>
      <el-row class="admin-track-result">
        <el-col :span="8">
          <el-card :body-style="{ padding: '10px' }" shadow="hover" class="admin-track-result-image-box">
            <el-image :src="picSrc" class="admin-track-result-image">
              <div slot="error" class="image-slot">
                <i class="el-icon-picture-outline"></i>
              </div>
            </el-image>
            <span style="font-size: 20px;display: flex; margin-top: 30px">提取到的水印图片</span>
          </el-card>
        </el-col>
        <el-card class="admin-track-result-fp" shadow="hover">
          <el-row>
            <el-col :span="12">
              <span style="font-size: 22px;">指纹信息</span>
            </el-col>
            <el-col :span="12">
              <span v-text="fingerPrint"></span>
            </el-col>
          </el-row>
        </el-card>
        <el-card class="admin-track-result-user" shadow="hover">
          <span class="admin-track-result-title" v-text="resultUser"></span>
          <el-progress type="circle" :percentage="probability" class="admin-track-circle-result"></el-progress>
        </el-card>
      </el-row>
    </el-main>
  </el-container>
</template>

<script>
import { adminTrack } from '@/api/admin'

export default {
  name: 'app',
  components: {},
  props: {},
  data () {
    return {
      resultUser: '是kale的概率为',
      picSrc: '',
      probability: 100,
      fingerPrint: '212121212121',
      uploadUrl: '/util'
    }
  },
  computed: {},
  watch: {},
  created () {
  },
  mounted () {
  },
  methods: {
    open () {
      this.$message.error('文件有误')
    },
    getResult (response, file, filelist) {
      if (response.msg === '文件有误') {
        this.open()
      }
    },
    beforeFileUpload (file) {
      const fd = new FormData()
      fd.append('file', file)
      adminTrack(fd)
        .then((response) => {
          this.picSrc = response.data.payload.picSrc
          this.resultUser = '是' + response.data.payload.phone + '的概率为'
          this.fingerPrint = response.data.payload.fingerPrint
          this.$message.success('成功')
        })
        .catch((error) => {
          console.log(error)
        })
    }
  }
}
</script>

<style scoped lang="less">
.admin-track-container {
  width: 100%;
  height: 100%;

  .admin-track-title {
    position: absolute;
    width: 100%;
    left: 0;
    top: 0;
    background-color: #545c64;
    padding-bottom: 40px;
    padding-top: 40px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 50px;
    color: white;
  }
}

.admin-track-upload {
  position: fixed;
  top: 150px;
  width: 400px;
  height: 600px;
  border-radius: 5px;
  box-shadow: 0 0 2px #b1b4ba;
  display: flex;
  justify-content: center;

  .admin-track-upload-box {
    margin-top: 30px;
  }
}

.admin-track-result {
  position: fixed;
  top: 150px;
  left: 750px;
  width: 800px;
  height: 600px;
  border-radius: 5px;
  box-shadow: 0 0 2px #b1b4ba;
}

.admin-track-result-image-box {
  position: absolute;
  left: 60px;
  margin-top: 20px;
  height: 200px;
}

.admin-track-result-image {
  width: 100px;
  height: 100px;
  background-color: #F5F7FA;
  left: 50%;
  transform: translate(-50%, 0%);
}

.admin-track-result-fp {
  width: 400px;
  height: 200px;
  margin-top: 20px;
  position: absolute;
  right: 10px;
}
.admin-track-result-user {
  position: absolute;
  width: 90%;
  height: 300px;
  left: 5%;
  top: 250px;
  display: flex;
  justify-content: space-around;
  align-items: center;

  .admin-track-result-title {
    font-size: 30px;
    position: absolute;
    left: 40px;
    top: 50%;
    transform: translate(0%, -50%);
  }
  .admin-track-circle-result {
    position: absolute;
    right: 40px;
    top: 50%;
    transform: translate(0%, -50%);
  }
}
</style>
