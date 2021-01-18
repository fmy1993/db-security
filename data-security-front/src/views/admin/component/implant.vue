<template>
  <div class="admin-implant-container">
    <el-table
      :data="tableData"
      height="80%">
      <el-table-column
        v-for="(item, key) in titleData"
        :key="key"
        :label="item.name"
        :prop="item.value"></el-table-column>
    </el-table>
  </div>
</template>

<script>
import { adminGetAllIis } from '@/api/admin'
import { refreshToken } from '@/util/tools'

export default {
  name: 'AdminImplantIndex',
  components: {},
  props: {},
  data () {
    return {
      tableData: [],
      titleData: [
        {
          name: 'id',
          value: 'id'
        }, {
          name: '索引表名',
          value: 'index_name'
        }, {
          name: '手机',
          value: 'phone'
        }, {
          name: '时间',
          value: 'datetime'
        }]
    }
  },
  computed: {},
  watch: {},
  created () {
  },
  mounted () {
    this.onGetAllIis()
  },
  methods: {
    onGetAllIis () {
      adminGetAllIis()
        .then((response) => {
          refreshToken(response)
          if (response.data.payload === 'null') {
            this.$message.warning('已经没有更多了')
          } else {
            this.tableData = response.data.payload
          }
        })
        .catch((error) => {
          console.log(error)
        })
    }
  }
}
</script>

<style scoped lang="less"></style>
