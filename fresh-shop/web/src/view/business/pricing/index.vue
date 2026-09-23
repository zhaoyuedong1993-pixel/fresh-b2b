<script>
export default {
  name: 'Pricing',
}
</script>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getGoodsList } from '@/api/goods'
import { getCategoryList } from '@/api/category'
import { batchUpdateCostPrice, applyPricing as applyPriceApi, getMarkupRate, batchParseText } from '@/api/pricing'

const loading = ref(false)
const goodsList = ref([])
const categoryList = ref([])
const categoryId = ref(null)
const keyword = ref('')
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const markupRate = ref(10)

const aiText = ref('')
const aiLoading = ref(false)
const aiResult = ref([])

const updatedCount = computed(() => aiResult.value.filter(i => i.updated).length)

const calculateSellPrice = (costPrice) => {
  if (!costPrice || costPrice <= 0) return '0.00'
  return (costPrice * (1 + markupRate.value / 100)).toFixed(2)
}

const calculateDiff = (costPrice) => {
  if (!costPrice || costPrice <= 0) return '0.00'
  return (costPrice * markupRate.value / 100).toFixed(2)
}

const loadCategory = async () => {
  try {
    const res = await getCategoryList()
    if (res.code === 0) {
      categoryList.value = res.data || []
    }
  } catch (e) {
    console.error(e)
  }
}

const loadGoods = async () => {
  loading.value = true
  try {
    const res = await getGoodsList({
      page: page.value,
      pageSize: pageSize.value,
      categoryId: categoryId.value,
      name: keyword.value
    })
    if (res.code === 0) {
      goodsList.value = res.data.list || []
      total.value = res.data.total || 0
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const loadMarkupRate = async () => {
  try {
    const res = await getMarkupRate()
    if (res.code === 0 && res.data) {
      markupRate.value = res.data.markupRate || 10
    }
  } catch (e) {
    console.error(e)
  }
}

const parseWithAI = async () => {
  if (!aiText.value.trim()) {
    ElMessage.warning('请输入采购价格信息')
    return
  }

  aiLoading.value = true
  aiResult.value = []

  try {
    const res = await batchParseText({ text: aiText.value })
    if (res.code === 0) {
      aiResult.value = res.data.items || []
      ElMessage.success('解析完成：共 ' + res.data.total + ' 个，成功匹配 ' + res.data.updated + ' 个')
      loadGoods()
    } else {
      ElMessage.error(res.msg || '解析失败')
    }
  } catch (e) {
    ElMessage.error('解析失败')
  } finally {
    aiLoading.value = false
  }
}

const batchUpdate = async () => {
  const updates = {}
  let count = 0

  goodsList.value.forEach(g => {
    if (g.costPrice && g.costPrice > 0) {
      updates[g.id] = g.costPrice
      count++
    }
  })

  if (count === 0) {
    ElMessage.warning('请先输入采购价')
    return
  }

  try {
    await ElMessageBox.confirm('已填写 ' + count + ' 个商品的采购价，确定批量更新吗？', '提示', { type: 'warning' })
    await batchUpdateCostPrice({ updates })
    ElMessage.success('采购价更新成功')
    loadGoods()
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error('更新失败')
    }
  }
}

const applyPricing = async () => {
  try {
    await ElMessageBox.confirm('确定应用当前采购价生成新的商户价吗？', '确认', { type: 'warning' })
    await applyPriceApi()
    ElMessage.success('价格已生效')
    loadGoods()
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error('应用失败')
    }
  }
}

onMounted(() => {
  loadCategory()
  loadGoods()
  loadMarkupRate()
})
</script>

<template>
  <div class="pricing-container">
    <el-card class="ai-card">
      <template #header>
        <div class="card-header">
          <span>AI 智能解析批量导入</span>
        </div>
      </template>

      <div class="ai-input-area">
        <el-input v-model="aiText" type="textarea" :rows="4" placeholder="粘贴采购价格信息，例如：今日莲藕1.5块一斤，豆芽1.5块一斤" class="ai-textarea" />
        <div class="ai-actions">
          <el-button type="primary" @click="parseWithAI" :loading="aiLoading">AI解析并导入</el-button>
          <el-button @click="aiText = ''">清空</el-button>
          <span class="ai-tip">AI会自动解析商品名称和价格</span>
        </div>
      </div>

      <div v-if="aiResult.length > 0" class="ai-result">
        <el-divider>解析结果</el-divider>
        <el-table :data="aiResult" size="small" border>
          <el-table-column prop="name" label="商品名称" width="150" />
          <el-table-column label="价格(元)" width="120">
            <template #default="{ row }"><span class="parsed-price">¥{{ row.price.toFixed(2) }}</span></template>
          </el-table-column>
          <el-table-column prop="unit" label="单位" width="80" />
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <el-tag v-if="row.updated" type="success" size="small">已更新</el-tag>
              <el-tag v-else type="info" size="small">未匹配</el-tag>
            </template>
          </el-table-column>
        </el-table>
        <div class="result-summary">共解析 {{ aiResult.length }} 个，成功匹配 {{ updatedCount }} 个</div>
      </div>
    </el-card>

    <el-card class="goods-card">
      <template #header>
        <div class="card-header">
          <span>商品采购价管理</span>
          <span class="markup-rate">加价比例: <el-tag type="success">{{ markupRate }}%</el-tag></span>
        </div>
      </template>

      <el-form :inline="true" class="filter-form">
        <el-form-item label="商品分类">
          <el-select v-model="categoryId" placeholder="全部" clearable @change="loadGoods">
            <el-option v-for="item in categoryList" :key="item.id" :label="item.title" :value="item.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="商品名称">
          <el-input v-model="keyword" placeholder="搜索商品" clearable @change="loadGoods" />
        </el-form-item>
        <el-form-item><el-button type="primary" @click="loadGoods">查询</el-button></el-form-item>
      </el-form>

      <div class="action-bar">
        <el-space>
          <el-button type="success" @click="batchUpdate">批量更新采购价</el-button>
          <el-button type="warning" @click="applyPricing">应用价格</el-button>
        </el-space>
      </div>

      <el-table :data="goodsList" border v-loading="loading" class="goods-table">
        <el-table-column prop="name" label="商品名称" min-width="180" />
        <el-table-column label="当前售价(元)" width="120">
          <template #default="{ row }"><span class="current-price">¥{{ (row.price || 0).toFixed(2) }}</span></template>
        </el-table-column>
        <el-table-column label="采购价(元)" width="180">
          <template #default="{ row }">
            <el-input-number v-model="row.costPrice" :min="0" :precision="2" :step="0.1" :controls="false" size="small" placeholder="输入采购价" style="width: 140px" />
          </template>
        </el-table-column>
        <el-table-column label="预计商户价(元)" width="140">
          <template #default="{ row }">
            <span v-if="row.costPrice > 0" class="sell-price">¥{{ calculateSellPrice(row.costPrice) }}</span>
            <span v-else class="no-price">-</span>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination v-model:current-page="page" v-model:page-size="pageSize" :total="total" :page-sizes="[10, 20, 50, 100]" layout="total, sizes, prev, pager, next, jumper" @size-change="loadGoods" @current-change="loadGoods" style="margin-top: 20px; text-align: right" />
    </el-card>
  </div>
</template>

<style lang="scss" scoped>
.pricing-container { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.markup-rate { font-size: 14px; color: #666; }
.filter-form { margin-bottom: 20px; }
.action-bar { margin-bottom: 20px; }
.ai-card { margin-bottom: 20px; }
.ai-input-area { margin-bottom: 15px; }
.ai-actions { display: flex; align-items: center; gap: 10px; }
.ai-tip { color: #999; font-size: 12px; margin-left: 10px; }
.ai-result { margin-top: 20px; }
.parsed-price { color: #67c23a; font-weight: bold; }
.result-summary { margin-top: 15px; padding: 10px; background: #f5f7fa; border-radius: 4px; color: #666; font-size: 14px; }
.goods-table .current-price { color: #909399; }
.goods-table .sell-price { color: #67c23a; font-weight: bold; }
.goods-table .no-price { color: #c0c4cc; }
</style>
