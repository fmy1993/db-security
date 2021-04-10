<template>
  <el-container class="track-container">
    <el-main>
      <div class="track-upload">
        <el-upload drag
                   :action="uploadUrl"
                   :headers="headers"
                   class="admin-track-upload-box"
                   :on-success="getResult">
          <i class="el-icon-upload"></i>
          <div class="el-upload__text">将文件拖到此处，或<em>点击上传</em></div>
          <div class="el-upload__tip" slot="tip">只能上传csv文件</div>
        </el-upload>
      </div>
      <el-row class="track-result">
        <el-col :span="8">
          <el-card :body-style="{ padding: '10px' }" shadow="hover" class="track-result-image-box">
            <el-image :src="pic" class="track-result-image">
              <div slot="error" class="image-slot">
                <i class="el-icon-picture-outline"></i>
              </div>
            </el-image>
            <span style="font-size: 20px;display: flex; margin-top: 30px">提取到的水印图片</span>
          </el-card>
        </el-col>
        <el-card class="track-result-fp" shadow="hover">
          <el-row>
            <el-col :span="12">
              <span style="font-size: 22px;">指纹信息</span>
            </el-col>
            <el-col :span="12">
              <span v-text="fingerPrint"></span>
            </el-col>
          </el-row>
        </el-card>
        <el-card class="track-result-user" shadow="hover">
          <span class="track-result-title" v-text="resultUser"></span>
          <el-progress type="circle" :percentage="probability" class="track-circle-result"></el-progress>
        </el-card>
      </el-row>
    </el-main>
  </el-container>
</template>

<script>

import request from '@/util/request'

export default {
  name: 'app',
  components: {},
  props: {},
  data () {
    return {
      resultUser: '',
      pic: '',
      probability: 100,
      fingerPrint: '212121212121',
      uploadUrl: request.defaults.baseURL + '/track',
      headers: {
        Authorization: localStorage.getItem('Authorization')
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
    getResult (res, file, filelist) {
      if (res.payload === '文件有误') {
        this.$message.error(res.payload)
      } else {
        this.pic = res.payload.pic
        this.resultUser = '是' + res.payload.user + '的概率为'
        this.fingerPrint = res.payload.fingerprint
        this.$message.success('成功')
      }
    }
  }
}
</script>

<style scoped lang="less">
.track-container {
  width: 100%;
  height: 100%;
}

.track-upload {
  position: fixed;
  top: 150px;
  width: 400px;
  height: 600px;
  border-radius: 5px;
  box-shadow: 0 0 2px #b1b4ba;
  display: flex;
  justify-content: center;

  .track-upload-box {
    margin-top: 30px;
  }
}

.track-result {
  position: fixed;
  top: 150px;
  left: 750px;
  width: 800px;
  height: 600px;
  border-radius: 5px;
  box-shadow: 0 0 2px #b1b4ba;
}

.track-result-image-box {
  position: absolute;
  left: 60px;
  margin-top: 20px;
  height: 200px;
}

.track-result-image {
  width: 100px;
  height: 100px;
  background-color: #F5F7FA;
  left: 50%;
  transform: translate(-50%, 0%);
}

.track-result-fp {
  width: 400px;
  height: 200px;
  margin-top: 20px;
  position: absolute;
  right: 10px;
}

.track-result-user {
  position: absolute;
  width: 90%;
  height: 300px;
  left: 5%;
  top: 250px;
  display: flex;
  justify-content: space-around;
  align-items: center;

  .track-result-title {
    font-size: 30px;
    position: absolute;
    left: 40px;
    top: 50%;
    transform: translate(0%, -50%);
  }

  .track-circle-result {
    position: absolute;
    right: 40px;
    top: 50%;
    transform: translate(0%, -50%);
  }
}
</style>
