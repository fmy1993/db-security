<template>
  <el-container class="user-home-container">
    <el-header class="user-home-container-header">
        <span class="user-home-container-header-title">
        <i :class="{
            'el-icon-s-fold': isCollapse,
            'el-icon-s-unfold': !isCollapse
          }" @click="isCollapse = !isCollapse">数据发布平台</i>
      </span>
      <el-dropdown style="display: flex; align-items: center">
        <el-avatar size="medium" :src="circleUrl" style="margin-right: 10px"></el-avatar>
        <div class="info-title">
          <span style="color: white">{{ phone }}</span>
          <i class="el-icon-arrow-down"></i>
          <el-dropdown-menu slot="dropdown">
            <el-dropdown-item @click.native="reviseDialogVisible = true">修改密码</el-dropdown-item>
            <el-dropdown-item @click.native="onLogout">退出</el-dropdown-item>
          </el-dropdown-menu>
        </div>
      </el-dropdown>
    </el-header>
    <el-container>
      <el-dialog
        title="修改密码"
        :visible.sync="reviseDialogVisible"
        width="60%"
        append-to-body
        :before-close="handleClose">
        <revise></revise>
      </el-dialog>
      <el-aside width="auto">
        <el-menu
          router
          class="user-home-container-aside"
          :collapse="isCollapse"
          default-active="2"
          :default-openeds="['1', '2']"
          text-color="#fff"
          background-color="#3f3f3f"
          active-text-color="white">
          <el-submenu index="1">
            <template slot="title">
              <i class="el-icon-document"></i>
              <span>数据表格</span>
            </template>
            <el-menu-item-group>
              <el-menu-item index="/user-home/staff">员工信息
              </el-menu-item>
            </el-menu-item-group>
          </el-submenu>
          <el-submenu index="2">
            <template slot="title">
              <i class="el-icon-download"></i>
              <span>下载数据</span>
            </template>
            <el-menu-item-group>
              <el-menu-item
                index="/user-home/staff"
                @click="fileDownload('staff')">病人表
              </el-menu-item>
            </el-menu-item-group>
          </el-submenu>
          <el-menu-item index="/user-home/visualization">
            <i class="el-icon-s-marketing"></i>
            <span slot="title">数据可视</span>
          </el-menu-item>
        </el-menu>
      </el-aside>
      <el-main>
        <router-view></router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script>
import { download } from '@/api/staff'
import { logout } from '@/api/user'
import Revise from '@/views/userhome/component/revise'

export default {
  name: 'UserHomeIndex',
  components: { Revise },
  props: {},
  data () {
    return {
      circleUrl: 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png',
      isCollapse: false,
      reviseDialogVisible: false,
      phone: JSON.parse(localStorage.getItem('user')).phone
    }
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
    },
    onLogout () {
      this.$confirm('是否退出?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        logout()
          .then(res => {
            if (res.data.msg === 'success') {
              this.$message.success('退出成功')
              localStorage.removeItem('Authorization')
              this.$router.push('/login')
            } else {
              this.$message.error('退出失败')
            }
          })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消退出'
        })
      })
    },
    handleClose () {
      this.reviseDialogVisible = false
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

  .user-home-container-header {
    display: flex;
    align-items: center;
    background-color: rgba(43, 43, 43, 0.9);
    color: white;
    border-bottom: 1px solid black;
    justify-content: space-between;
  }

  .user-home-container-header-title {
    font-size: 20px;
  }

  .user-home-container-aside:not(.el-menu--collapse) {
    width: 250px;
    height: 100%;
    background-color: #3f3f3f;
  }

  .user-home-container-aside {
    width: 70px;
    height: 100%;
    background-color: #3f3f3f;
  }
}
</style>
