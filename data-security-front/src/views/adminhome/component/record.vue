<template>
  <div class="record-container">
    <el-table
      :data="tableData"
      style="width: 100%;height: 100%"
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
import { getAllRecords } from '@/api/record'

export default {
  name: 'AdminRecordIndex',
  components: {},
  props: {},
  data () {
    return {
      tableData: [],
      titleData: [
        {
          name: 'record_id',
          value: 'record_id'
        }, {
          name: 'user_id',
          value: 'user_id'
        }, {
          name: '用户名',
          value: 'phone'
        }, {
          name: '下载时间',
          value: 'download_time'
        }]
    }
  },
  computed: {},
  watch: {},
  created () {
  },
  mounted () {
    this.onGetAllRecords()
  },
  methods: {
    onGetAllRecords () {
      getAllRecords()
        .then(res => {
          if (res.data.payload === 'null') {
            this.$message.warning('已经没有更多了')
          } else {
            this.tableData = res.data.payload
          }
        })
    }
  }
}
</script>

<style scoped lang="less">
.record-container {
  width: 100%;
  height: 100%;
  overflow: hidden;
}
</style>
