<template>
  <div class="admin-container">
    <el-container class="admin-home-container-container">
      <el-aside width="250px" class="admin-home-container-aside">
        <el-menu
          router
          default-active="2-1"
          :default-openeds="['1', '2','3']"
          background-color="#545c64"
          text-color="#fff"
          active-text-color="#ffd04b">
          <span class="admin-home-container-aside-title">后台界面</span>
          <el-submenu index="1">
            <template slot="title">
              <i class="el-icon-user"></i>
              <span>个人主页</span>
            </template>
            <el-menu-item-group>
              <el-menu-item index="/admin-home/revise">修改信息
              </el-menu-item>
            </el-menu-item-group>
          </el-submenu>
          <el-submenu index="2">
            <template slot="title">
              <i class="el-icon-document"></i>
              <span>数据表格</span>
            </template>
            <el-menu-item-group>
              <el-menu-item index="/admin-home/patient">病人信息
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
              <el-menu-item index="/admin-home/implant">嵌入记录
              </el-menu-item>
            </el-menu-item-group>
          </el-submenu>
          <el-submenu index="3">
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
      <el-main class="admin-body">
        <router-view></router-view>
      </el-main>
    </el-container>
  </div>
</template>

<script>
export default {
  name: 'AdminIndex',
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
  methods: {}
}
</script>

<style scoped lang="less">
.admin-container {
  position: fixed;
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;

  .admin-home-container-container {
    height: 100%;

    .admin-home-container-aside {
      width: 250px;
      height: 100%;
      background-color: #545c64;
    }
  }

}

.admin-home-container-aside-title {
  color: white;
  font-size: 30px;
  display: flex;
  justify-content: center;
  margin-top: 30px;
  padding-bottom: 70px;
}

.admin-body {
  position: relative;
  top: 0;
}
</style>
