<template>
  <div class="bill-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>账单管理</span>
          <el-button type="primary" @click="generateBillDialog = true">生成账单</el-button>
        </div>
      </template>

      <!-- 账单列表 -->
      <el-table :data="billList" border v-loading="loading">
        <el-table-column prop="period" label="账期" width="120">
          <template #default="{ row }">
            <el-tag type="info">{{ row.period }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="bill_no" label="账单编号" width="180" />
        <el-table-column prop="order_count" label="订单数" width="100" align="center" />
        <el-table-column prop="total_amount" label="总金额(元)" width="150" align="right">
          <template #default="{ row }">
            <span class="amount">¥{{ (row.total_amount || 0).toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'warning'">
              {{ row.status === 1 ? '已结算' : '待结算' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" min-width="150" />
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="primary" @click="viewDetail(row)">详情</el-button>
            <el-button
              v-if="row.status === 0"
              size="small"
              type="success"
              @click="markSettled(row)"
            >标记结算</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        @size-change="loadBillList"
        @current-change="loadBillList"
        style="margin-top: 20px; text-align: right"
      />
    </el-card>

    <!-- 生成账单弹窗 -->
    <el-dialog v-model="generateBillDialog" title="生成账单" width="400px">
      <el-form :model="generateForm" label-width="80px">
        <el-form-item label="账期">
          <el-date-picker
            v-model="generateForm.month"
            type="month"
            placeholder="选择月份"
            format="YYYY-MM"
            value-format="YYYY-MM"
          />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="generateForm.remark" type="textarea" rows="3" placeholder="备注信息" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="generateBillDialog = false">取消</el-button>
        <el-button type="primary" @click="handleGenerateBill">确认生成</el-button>
      </template>
    </el-dialog>

    <!-- 账单详情弹窗 -->
    <el-dialog v-model="detailDialog" title="账单详情" width="900px">
      <div v-if="currentBill" class="bill-detail">
        <el-descriptions :column="3" border>
          <el-descriptions-item label="账期">{{ currentBill.period }}</el-descriptions-item>
          <el-descriptions-item label="账单编号">{{ currentBill.bill_no }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="currentBill.status === 1 ? 'success' : 'warning'">
              {{ currentBill.status === 1 ? '已结算' : '待结算' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="订单数量">{{ currentBill.order_count }}</el-descriptions-item>
          <el-descriptions-item label="总金额">
            <span class="amount">¥{{ (currentBill.total_amount || 0).toFixed(2) }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="备注" :span="3">{{ currentBill.remark || '-' }}</el-descriptions-item>
        </el-descriptions>

        <h4 style="margin: 20px 0 10px">订单明细</h4>
        <el-table :data="currentBill.orders || []" border size="small">
          <el-table-column prop="id" label="订单ID" width="80" />
          <el-table-column prop="order_sn" label="订单编号" width="180" />
          <el-table-column prop="total" label="金额" width="120" align="right">
            <template #default="{ row }">
              ¥{{ (row.total || 0).toFixed(2) }}
            </template>
          </el-table-column>
          <el-table-column prop="order_status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag size="small">
                {{ orderStatusText(row.order_status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="下单时间" />
        </el-table>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getBillList, generateBill, getBillDetail, updateBillStatus } from '@/api/bill'

const loading = ref(false)
const billList = ref([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 生成账单弹窗
const generateBillDialog = ref(false)
const generateForm = ref({
  month: '',
  remark: ''
})

// 详情弹窗
const detailDialog = ref(false)
const currentBill = ref(null)

// 订单状态文本
const orderStatusText = (status) => {
  const map = {
    0: '待确认',
    1: '已确认',
    2: '已完成',
    3: '已取消'
  }
  return map[status] || '未知'
}

// 加载账单列表
const loadBillList = async () => {
  loading.value = true
  try {
    const res = await getBillList({ page: page.value, pageSize: pageSize.value })
    if (res.code === 0) {
      billList.value = res.data.list || []
      total.value = res.data.total || 0
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

// 生成账单
const handleGenerateBill = async () => {
  if (!generateForm.value.month) {
    ElMessage.warning('请选择账期')
    return
  }

  const [year, month] = generateForm.value.month.split('-').map(Number)

  try {
    await generateBill({ year, month })
    ElMessage.success('账单生成成功')
    generateBillDialog.value = false
    loadBillList()
  } catch (e) {
    ElMessage.error('生成失败')
  }
}

// 查看详情
const viewDetail = async (row) => {
  try {
    const res = await getBillDetail(row.id)
    if (res.code === 0) {
      currentBill.value = res.data
      detailDialog.value = true
    }
  } catch (e) {
    ElMessage.error('加载详情失败')
  }
}

// 标记结算
const markSettled = async (row) => {
  try {
    await ElMessageBox.confirm('确定将该账单标记为已结算?', '提示', { type: 'warning' })
    await updateBillStatus(row.id, 1)
    ElMessage.success('标记成功')
    loadBillList()
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error('操作失败')
    }
  }
}

onMounted(() => {
  loadBillList()
})
</script>

<style lang="scss" scoped>
.bill-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.amount {
  color: #67c23a;
  font-weight: bold;
  font-size: 16px;
}

.bill-detail {
  .amount {
    color: #67c23a;
    font-weight: bold;
  }
}
</style>
