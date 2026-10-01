<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="searchForm" :inline="true" :model="searchInfo">
        <el-form-item label="公司名称">
          <el-input v-model="searchInfo.name" placeholder="公司名称" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="getTableData">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openAddDialog">新增公司</el-button>
      </div>
      <el-table :data="tableData" border>
        <el-table-column align="left" label="ID" min-width="80" prop="id" />
        <el-table-column align="left" label="公司名称" min-width="150" prop="name" />
        <el-table-column align="left" label="客户类型" min-width="100" prop="companyType">
          <template #default="{ row }">
            <el-tag :type="row.companyType === 'monthly' ? 'primary' : 'success'">
              {{ row.companyType === 'monthly' ? '月度' : '零售' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="联系人" min-width="100" prop="contact" />
        <el-table-column align="left" label="联系电话" min-width="120" prop="phone" />
        <el-table-column align="left" label="地址" min-width="200" prop="address" show-overflow-tooltip />
        <el-table-column align="left" label="邀请码" min-width="120" prop="invitationCode" />
        <el-table-column align="left" label="加价比例" min-width="100" prop="markupRate">
          <template #default="{ row }">
            {{ row.markupRate }}%
          </template>
        </el-table-column>
        <el-table-column align="left" label="状态" min-width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="创建时间" min-width="160" prop="createdAt">
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" min-width="150" fixed="right">
          <template #default="scope">
            <el-button type="primary" link icon="edit" @click="openEditDialog(scope.row)">编辑</el-button>
            <el-button type="danger" link icon="delete" @click="deleteCompanyFunc(scope.row)">删除</el-button>
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

    <!-- 新增/编辑弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogFlag === 'add' ? '新增公司' : '编辑公司'"
      width="500px"
    >
      <el-form ref="formRef" :model="form" label-width="100px" :rules="rules">
        <el-form-item label="公司名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入公司名称" />
        </el-form-item>
        <el-form-item label="客户类型" prop="companyType">
          <el-radio-group v-model="form.companyType">
            <el-radio label="monthly">月度</el-radio>
            <el-radio label="retail">零售</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="联系人" prop="contact">
          <el-input v-model="form.contact" placeholder="请输入联系人" />
        </el-form-item>
        <el-form-item label="联系电话" prop="phone">
          <el-input v-model="form.phone" placeholder="请输入联系电话" />
        </el-form-item>
        <el-form-item label="地址" prop="address">
          <el-input v-model="form.address" placeholder="请输入地址" />
        </el-form-item>
        <el-form-item v-if="dialogFlag === 'add'" label="管理员账号" prop="username">
          <el-input v-model="form.username" placeholder="登录用户名" />
        </el-form-item>
        <el-form-item v-if="dialogFlag === 'add'" label="管理员密码" prop="password">
          <el-input v-model="form.password" placeholder="登录密码" show-password />
        </el-form-item>
        <el-form-item label="加价比例" prop="markupRate">
          <el-input-number v-model="form.markupRate" :min="0" :max="100" /> %
        </el-form-item>
        <el-form-item label="截单时间" prop="cutoffTime">
          <el-time-picker v-model="form.cutoffTime" format="HH:mm" value-format="HH:mm" placeholder="选择时间" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getCompanyList, createCompany, updateCompany, deleteCompany } from '@/api/company'
import { formatDate } from '@/utils/format'

const tableData = ref([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const searchInfo = ref({})

const dialogVisible = ref(false)
const dialogFlag = ref('add')
const formRef = ref(null)
const form = ref({
  id: '',
  name: '',
  contact: '',
  phone: '',
  address: '',
  username: '',
  password: '',
  companyType: 'monthly',
  markupRate: 10,
  cutoffTime: '22:00',
  status: 1
})

const rules = {
  name: [{ required: true, message: '请输入公司名称', trigger: 'blur' }],
  phone: [{ required: true, message: '请输入联系电话', trigger: 'blur' }],
  username: [{ required: true, message: '请输入管理员账号', trigger: 'blur' }],
  password: [{ required: true, message: '请输入管理员密码', trigger: 'blur' }]
}

const getTableData = async () => {
  const res = await getCompanyList({
    page: page.value,
    pageSize: pageSize.value,
    ...searchInfo.value
  })
  if (res.code === 0) {
    tableData.value = res.data.list || []
    total.value = res.data.total || 0
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

const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

const openAddDialog = () => {
  dialogFlag.value = 'add'
  form.value = {
    name: '',
    contact: '',
    phone: '',
    address: '',
    username: '',
    password: '',
    companyType: 'monthly',
    markupRate: 10,
    cutoffTime: '22:00',
    status: 1
  }
  dialogVisible.value = true
}

const openEditDialog = (row) => {
  dialogFlag.value = 'edit'
  form.value = {
    id: row.id,
    name: row.name,
    contact: row.contact || '',
    phone: row.phone || '',
    address: row.address || '',
    companyType: row.companyType || 'monthly',
    markupRate: row.markupRate || 10,
    cutoffTime: row.cutoffTime || '22:00',
    status: row.status
  }
  dialogVisible.value = true
}

const submitForm = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return

    try {
      if (dialogFlag.value === 'add') {
        await createCompany(form.value)
        ElMessage.success('创建成功')
      } else {
        await updateCompany(form.value)
        ElMessage.success('更新成功')
      }
      dialogVisible.value = false
      getTableData()
    } catch (e) {
      ElMessage.error(e.message || '操作失败')
    }
  })
}

const deleteCompanyFunc = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定删除公司「${row.name}」吗？删除后不可恢复。`,
      '删除确认',
      { confirmButtonText: '确定删除', cancelButtonText: '取消', type: 'warning' }
    )
  } catch {
    return
  }

  try {
    await deleteCompany(row.id)
    ElMessage.success('删除成功')
    getTableData()
  } catch (e) {
    ElMessage.error(e.message || '删除失败')
  }
}

onMounted(() => {
  getTableData()
})
</script>
