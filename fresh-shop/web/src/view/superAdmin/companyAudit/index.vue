<template>
  <div>
    <div class="gva-table-box">
      <el-table :data="tableData" border v-loading="loading">
        <el-table-column align="left" label="ID" min-width="80" prop="id" />
        <el-table-column align="left" label="公司名称" min-width="150" prop="name" />
        <el-table-column align="left" label="联系人" min-width="100" prop="contact" />
        <el-table-column align="left" label="联系电话" min-width="120" prop="phone" />
        <el-table-column align="left" label="地址" min-width="200" prop="address" show-overflow-tooltip />
        <el-table-column align="left" label="申请时间" min-width="160" prop="applyTime">
          <template #default="{ row }">
            {{ formatDate(row.applyTime) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" min-width="200" fixed="right">
          <template #default="scope">
            <el-button type="success" @click="handleAudit(scope.row, true)">通过</el-button>
            <el-button type="danger" @click="handleRefuse(scope.row)">拒绝</el-button>
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

    <!-- 审核弹窗 -->
    <el-dialog v-model="auditDialogVisible" title="审核入驻申请" width="400px">
      <el-form :model="auditForm" label-width="100px">
        <el-form-item label="管理员密码">
          <el-input v-model="auditForm.password" type="password" show-password placeholder="设置管理员登录密码" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="auditForm.remark" type="textarea" rows="3" placeholder="审核备注（可选）" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="auditDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmAudit">确认通过</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getPendingCompanyList, auditCompany } from '@/api/company'
import { formatDate } from '@/utils/format'

const loading = ref(false)
const tableData = ref([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const auditDialogVisible = ref(false)
const auditForm = ref({
  id: '',
  password: '',
  remark: ''
})
const currentAuditRow = ref(null)

const getTableData = async () => {
  loading.value = true
  try {
    const res = await getPendingCompanyList({
      page: page.value,
      pageSize: pageSize.value
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

const handleAudit = (row) => {
  currentAuditRow.value = row
  auditForm.value = {
    id: row.id,
    password: '',
    remark: ''
  }
  auditDialogVisible.value = true
}

const handleRefuse = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定拒绝「${row.name}」的入驻申请吗？`,
      '拒绝确认',
      { confirmButtonText: '确定拒绝', cancelButtonText: '取消', type: 'warning' }
    )
  } catch {
    return
  }

  try {
    await auditCompany({ id: row.id, pass: false, remark: '' })
    ElMessage.success('已拒绝')
    getTableData()
  } catch (e) {
    ElMessage.error(e.message || '操作失败')
  }
}

const confirmAudit = async () => {
  if (!auditForm.value.password) {
    ElMessage.warning('请设置管理员密码')
    return
  }

  try {
    await auditCompany({
      id: auditForm.value.id,
      pass: true,
      password: auditForm.value.password,
      remark: auditForm.value.remark
    })
    ElMessage.success('审核通过')
    auditDialogVisible.value = false
    getTableData()
  } catch (e) {
    ElMessage.error(e.message || '操作失败')
  }
}

onMounted(() => {
  getTableData()
})
</script>
