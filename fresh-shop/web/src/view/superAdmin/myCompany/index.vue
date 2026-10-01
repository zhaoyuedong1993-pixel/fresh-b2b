<template>
  <div>
    <div class="gva-table-box">
      <el-descriptions title="公司信息" :column="2" border v-if="companyInfo">
        <el-descriptions-item label="公司名称">{{ companyInfo.name }}</el-descriptions-item>
        <el-descriptions-item label="客户类型">
          <el-tag :type="companyInfo.companyType === 'monthly' ? 'primary' : 'success'">
            {{ companyInfo.companyType === 'monthly' ? '月度' : '零售' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="联系人">{{ companyInfo.contact || '无' }}</el-descriptions-item>
        <el-descriptions-item label="联系电话">{{ companyInfo.phone || '无' }}</el-descriptions-item>
        <el-descriptions-item label="地址" :span="2">{{ companyInfo.address || '无' }}</el-descriptions-item>
        <el-descriptions-item label="邀请码">{{ companyInfo.invitationCode || '无' }}</el-descriptions-item>
        <el-descriptions-item label="加价比例">{{ companyInfo.markupRate }}%</el-descriptions-item>
        <el-descriptions-item label="截单时间">{{ companyInfo.cutoffTime }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="companyInfo.status === 1 ? 'success' : 'danger'">
            {{ companyInfo.status === 1 ? '启用' : '禁用' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatDate(companyInfo.createdAt) }}</el-descriptions-item>
      </el-descriptions>

      <div style="margin-top: 20px" v-if="companyInfo">
        <el-button type="primary" @click="openEditDialog">编辑公司信息</el-button>
      </div>
    </div>

    <!-- 编辑弹窗 -->
    <el-dialog v-model="dialogVisible" title="编辑公司信息" width="500px">
      <el-form ref="formRef" :model="form" label-width="100px" :rules="rules">
        <el-form-item label="公司名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入公司名称" />
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
        <el-form-item label="加价比例" prop="markupRate">
          <el-input-number v-model="form.markupRate" :min="0" :max="100" /> %
        </el-form-item>
        <el-form-item label="截单时间" prop="cutoffTime">
          <el-time-picker v-model="form.cutoffTime" format="HH:mm" value-format="HH:mm" placeholder="选择时间" />
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
import { ElMessage } from 'element-plus'
import { getCompanyDetail, updateCompany } from '@/api/company'
import { formatDate } from '@/utils/format'
import { useUserStore } from '@/pinia/modules/user'

const userStore = useUserStore()
const companyInfo = ref(null)
const dialogVisible = ref(false)
const formRef = ref(null)

const form = ref({
  id: '',
  name: '',
  contact: '',
  phone: '',
  address: '',
  markupRate: 10,
  cutoffTime: '22:00'
})

const rules = {
  name: [{ required: true, message: '请输入公司名称', trigger: 'blur' }],
  phone: [{ required: true, message: '请输入联系电话', trigger: 'blur' }]
}

const getCompanyInfo = async () => {
  // 管理员只能查看自己的公司
  const companyId = userStore.userInfo.companyId
  if (!companyId) {
    ElMessage.warning('您还没有归属的公司')
    return
  }

  const res = await getCompanyDetail(companyId)
  if (res.code === 0) {
    companyInfo.value = res.data
  }
}

const openEditDialog = () => {
  form.value = {
    id: companyInfo.value.id,
    name: companyInfo.value.name,
    contact: companyInfo.value.contact || '',
    phone: companyInfo.value.phone || '',
    address: companyInfo.value.address || '',
    markupRate: companyInfo.value.markupRate || 10,
    cutoffTime: companyInfo.value.cutoffTime || '22:00'
  }
  dialogVisible.value = true
}

const submitForm = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return

    try {
      await updateCompany(form.value)
      ElMessage.success('更新成功')
      dialogVisible.value = false
      getCompanyInfo()
    } catch (e) {
      ElMessage.error(e.message || '更新失败')
    }
  })
}

onMounted(() => {
  getCompanyInfo()
})
</script>
