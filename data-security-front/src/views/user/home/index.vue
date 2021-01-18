<template>
  <div class="user-home-container">
    <el-container class="user-home-container-container">
      <el-aside class="user-home-container-aside" width="250px">
        <el-menu
          router
          default-active="2-1"
          :default-openeds="['1', '2','3']"
          background-color="#545c64"
          text-color="#fff"
          active-text-color="#ffd04b">
          <span class="user-home-container-aside-title">用户数据界面</span>
          <el-submenu index="1">
            <template slot="title">
              <i class="el-icon-user"></i>
              <span>个人主页</span>
            </template>
            <el-menu-item-group>
              <el-menu-item index="/user-home/revise">信息修改
              </el-menu-item>
            </el-menu-item-group>
          </el-submenu>
          <el-submenu index="2">
            <template slot="title">
              <i class="el-icon-document"></i>
              <span>数据表格</span>
            </template>
            <el-menu-item-group>
              <el-menu-item index="/user-home/patient">病人信息
              </el-menu-item>
            </el-menu-item-group>
          </el-submenu>
          <el-submenu index="3">
            <template slot="title">
              <i class="el-icon-download"></i>
              <span>下载数据</span>
            </template>
            <el-menu-item-group>
              <el-menu-item
                index="/user-home/patient"
                @click="fileDownload('病人表')">病人表
              </el-menu-item>
            </el-menu-item-group>
          </el-submenu>
          <el-menu-item index="/user-home/visualization">
            <i class="el-icon-s-marketing"></i>
            <span slot="title">数据可视</span>
          </el-menu-item>
        </el-menu>
      </el-aside>
      <el-main id="body">
        <router-view></router-view>
      </el-main>
    </el-container>
  </div>
</template>

<script>
import { download } from '@/api/patient'

export default {
  name: 'UserHomeIndex',
  components: {},
  props: {},
  data () {
    return {}
  },
  computed: {},
  watch: {},
  created () {
  },
  mounted () {
  },
  methods: {
    fileDownload (fileName) {
      download()
        .then((response) => {
          const { data } = response
          const url = window.URL.createObjectURL(new Blob([data]))
          const link = document.createElement('a')
          link.style.display = 'none'
          link.href = url
          link.setAttribute('download', fileName + '.csv')
          document.body.appendChild(link)
          link.click()
        })
        .catch((error) => {
          console.log(error)
        })
    }
  }
}
</script>

<style scoped lang="less">
.user-home-container {
  position: fixed;
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;

  .user-home-container-container {
    height: 100%;

    .user-home-container-aside {
      width: 250px;
      height: 100%;
      background-color: #545c64;
    }
  }

}

.user-home-container-aside-title {
  color: white;
  font-size: 30px;
  display: flex;
  justify-content: center;
  margin-top: 30px;
  padding-bottom: 70px;
}
</style>
