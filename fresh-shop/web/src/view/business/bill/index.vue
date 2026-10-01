<template>
  <div class="bill-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>月度账单</span>
        </div>
      </template>

      <!-- 账单列表 -->
      <el-table :data="billList" border v-loading="loading" @row-click="viewDetail">
        <el-table-column prop="period" label="账期" width="120">
          <template #default="{ row }">
            <el-tag type="info">{{ row.period }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="order_count" label="订单数" width="120" align="center" />
        <el-table-column prop="total_amount" label="总金额(元)" width="150" align="right">
          <template #default="{ row }">
            <span class="amount">¥{{ (Number(row.total_amount) || 0).toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="settlement_type" label="结算类型" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="row.settlement_type === 2 ? 'success' : 'warning'">
              {{ row.settlement_type === 2 ? '已结清' : '未结清' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120">
          <template #default="{ row }">
            <el-button size="small" type="primary" @click.stop="viewDetail(row)">查看详情</el-button>
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

    <!-- 账单详情弹窗 -->
    <el-dialog v-model="detailDialog" title="账单详情" width="900px">
      <div v-if="currentBill">
        <el-descriptions :column="3" border>
          <el-descriptions-item label="账期">{{ currentBill.period }}</el-descriptions-item>
          <el-descriptions-item label="订单数量">{{ currentBill.order_count }}</el-descriptions-item>
          <el-descriptions-item label="总金额">
            <span class="amount">¥{{ (Number(currentBill.total_amount) || 0).toFixed(2) }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="结算状态">
            <el-tag :type="currentBill.status === 1 ? 'success' : 'warning'">
              {{ currentBill.status === 1 ? '已结算' : '待结算' }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>

        <h4 style="margin: 20px 0 10px">订单明细</h4>
        <el-table :data="currentBill.orders || []" border size="small">
          <el-table-column prop="order_sn" label="订单编号" width="200" />
          <el-table-column prop="shipment_name" label="收货人" width="120" />
          <el-table-column prop="shipment_mobile" label="电话" width="140" />
          <el-table-column prop="total" label="金额" width="120" align="right">
            <template #default="{ row }">
              ¥{{ (Number(row.total) || 0).toFixed(2) }}
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="下单时间" />
        </el-table>

        <div style="margin-top: 20px; text-align: right;">
          <el-button type="success" @click="markSettled" :disabled="currentBill.status === 1">
            标记已结算
          </el-button>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getBillList, getBillDetail, updateBillStatus } from '@/api/bill'
import { useUserStore } from '@/pinia/modules/user'

const userStore = useUserStore()

const loading = ref(false)
const billList = ref([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const detailDialog = ref(false)
const currentBill = ref(null)

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

// 查看详情
const viewDetail = async (row) => {
  try {
    const res = await getBillDetail({ period: row.period })
    if (res.code === 0) {
      currentBill.value = res.data
      detailDialog.value = true
    }
  } catch (e) {
    ElMessage.error('加载详情失败')
  }
}

// 标记结算
const markSettled = async () => {
  if (!currentBill.value) return
  try {
    await ElMessageBox.confirm('确定将该账期标记为已结算?', '提示', { type: 'warning' })
    await updateBillStatus({ period: currentBill.value.period, status: 1 })
    ElMessage.success('标记成功')
    detailDialog.value = false
    await loadBillList()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('操作失败')
  }
}

onMounted(async () => {
  await loadBillList()
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
</style>
