<template>
  <div>
    <div class="gva-table-box">
      <el-table :data="tableData" border v-loading="loading">
        <el-table-column align="left" label="ID" min-width="80" prop="id" />
        <el-table-column align="left" label="手机号" min-width="120" prop="phone" />
        <el-table-column align="left" label="用户昵称" min-width="100" prop="nickName" />
        <el-table-column align="left" label="申请时间" min-width="160" prop="applyTime">
          <template #default="{ row }">
            {{ formatDate(row.applyTime) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" min-width="150" fixed="right">
          <template #default="scope">
            <el-button type="success" size="small" @click="handleApprove(scope.row)">通过</el-button>
            <el-button type="danger" size="small" @click="handleReject(scope.row)">拒绝</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getUserList, setUserInfo } from '@/api/user'
import { formatDate } from '@/utils/format'

const loading = ref(false)
const tableData = ref([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const getTableData = async () => {
  loading.value = true
  try {
    // 获取待审核用户（auditStatus=0）
    const res = await getUserList({
      page: page.value,
      pageSize: pageSize.value,
      auditStatus: 0
    })
    if (res.code === 0) {
      tableData.value = res.data.list || []
      total.value = res.data.total || 0
    }
  } finally {
    loading.value = false
  }
}

const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

const handleApprove = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定通过「${row.phone}」的入会申请吗？`,
      '审批确认',
      { confirmButtonText: '确定通过', cancelButtonText: '取消', type: 'success' }
    )
  } catch {
    return
  }

  try {
    await setUserInfo({ id: row.ID, auditStatus: 1 })
    ElMessage.success('审批通过')
    getTableData()
  } catch (e) {
    ElMessage.error(e.message || '操作失败')
  }
}

const handleReject = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定拒绝「${row.phone}」的入会申请吗？`,
      '拒绝确认',
      { confirmButtonText: '确定拒绝', cancelButtonText: '取消', type: 'warning' }
    )
  } catch {
    return
  }

  try {
    await setUserInfo({ id: row.ID, auditStatus: 4 }) // 4=已拒绝
    ElMessage.success('已拒绝')
    getTableData()
  } catch (e) {
    ElMessage.error(e.message || '操作失败')
  }
}

onMounted(() => {
  getTableData()
})
</script>
