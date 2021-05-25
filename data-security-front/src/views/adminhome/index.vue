<template>
  <el-container class="admin-container">
    <el-header class="admin-container-header">
        <span class="admin-container-header-title">
        <i :class="{
            'el-icon-s-fold': isCollapse,
            'el-icon-s-unfold': !isCollapse
          }" @click="isCollapse = !isCollapse">管理平台</i>
      </span>
      <el-dropdown style="display: flex; align-items: center">
        <el-avatar size="medium" :src="circleUrl" style="margin-right: 10px"></el-avatar>
        <div class="info-title">
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
          class="admin-container-aside"
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
              <el-menu-item index="/admin-home/staff">员工信息
              </el-menu-item>
            </el-menu-item-group>
            <el-menu-item-group>
              <el-menu-item index="/admin-home/user">用户信息
              </el-menu-item>
            </el-menu-item-group>
            <el-menu-item-group>
              <el-menu-item index="/admin-home/ip">封禁ip
              </el-menu-item>
            </el-menu-item-group>
            <el-menu-item-group>
              <el-menu-item index="/admin-home/record">下载记录
              </el-menu-item>
            </el-menu-item-group>
          </el-submenu>
          <el-submenu index="2">
            <template slot="title">
              <i class="el-icon-document"></i>
              <span>盗版追踪</span>
            </template>
            <el-menu-item-group>
              <el-menu-item index="/admin-home/track">盗版追踪
              </el-menu-item>
            </el-menu-item-group>
          </el-submenu>
        </el-menu>
      </el-aside>
      <el-main>
        <router-view></router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script>
import { logout } from '@/api/user'
import Revise from '@/views/userhome/component/revise'

export default {
  name: 'AdminHomeIndex',
  components: { Revise },
  props: {},
  data () {
    return {
      circleUrl: 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png',
      isCollapse: false,
      reviseDialogVisible: false
    }
  },
  computed: {},
  watch: {},
  created () {
  },
  mounted () {
  },
  methods: {
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
.admin-container {
  position: fixed;
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;

  .admin-container-header {
    display: flex;
    align-items: center;
    background-color: rgba(43, 43, 43, 0.9);
    color: white;
    border-bottom: 1px solid black;
    justify-content: space-between;
  }

  .admin-container-header-title {
    font-size: 20px;
  }

  .admin-container-aside:not(.el-menu--collapse) {
    width: 250px;
    height: 100%;
    background-color: #3f3f3f;
  }

  .admin-container-aside {
    width: 70px;
    height: 100%;
    background-color: #3f3f3f;
  }
}
</style>
