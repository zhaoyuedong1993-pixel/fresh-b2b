<template>
  <div class="pricing-container">
    <!-- AI智能解析 -->
    <el-card class="ai-card">
      <template #header>
        <div class="card-header">
          <span>📝 AI 智能解析批量导入</span>
        </div>
      </template>

      <div class="ai-input-area">
        <el-input
          v-model="aiText"
          type="textarea"
          :rows="4"
          placeholder="粘贴采购价格信息，例如：
今日莲藕1.5块一斤，豆芽1.5块一斤，土豆1块钱一斤
生菜2元/斤，白菜1.5元/斤"
          class="ai-textarea"
        />
        <div class="ai-actions">
          <el-button type="primary" @click="parseWithAI" :loading="aiLoading">
            🤖 AI解析并导入
          </el-button>
          <el-button @click="aiText = ''">清空</el-button>
          <span class="ai-tip">AI会自动解析商品名称和价格，匹配后自动更新采购价</span>
        </div>
      </div>

      <!-- AI解析结果 -->
      <div v-if="aiResult.length > 0" class="ai-result">
        <el-divider>解析结果</el-divider>
        <el-table :data="aiResult" size="small" border>
          <el-table-column prop="name" label="商品名称" width="150" />
          <el-table-column prop="price" label="解析价格(元)" width="120">
            <template #default="{ row }">
              <span class="parsed-price">¥{{ row.price.toFixed(2) }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="unit" label="单位" width="80" />
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <el-tag v-if="row.updated" type="success" size="small">已更新</el-tag>
              <el-tag v-else type="info" size="small">未匹配</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作">
            <template #default="{ row }">
              <span v-if="row.updated" class="success-text">✓ 已更新采购价</span>
              <span v-else class="gray-text">商品未在系统中找到</span>
            </template>
          </el-table-column>
        </el-table>
        <div class="result-summary">
          共解析 {{ aiResult.length }} 个商品，成功匹配并更新 {{ updatedCount }} 个
        </div>
      </div>
    </el-card>

    <!-- 商品列表 -->
    <el-card class="goods-card">
      <template #header>
        <div class="card-header">
          <span>📦 商品采购价管理</span>
          <el-space>
            <span class="markup-rate">
              加价比例: <el-tag type="success">{{ markupRate }}%</el-tag>
            </span>
          </el-space>
        </div>
      </template>

      <!-- 筛选条件 -->
      <el-form :inline="true" class="filter-form">
        <el-form-item label="商品分类">
          <el-select v-model="categoryId" placeholder="全部" clearable @change="loadGoods">
            <el-option
              v-for="item in categoryList"
              :key="item.id"
              :label="item.title"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="商品名称">
          <el-input v-model="keyword" placeholder="搜索商品" clearable @change="loadGoods" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadGoods">查询</el-button>
        </el-form-item>
      </el-form>

      <!-- 操作按钮 -->
      <div class="action-bar">
        <el-space>
          <el-button type="success" @click="batchUpdate">💾 批量更新采购价</el-button>
          <el-button type="warning" @click="applyPricing">✨ 应用价格</el-button>
        </el-space>
        <span class="tip">先填写采购价，再点击"批量更新"，确认无误后点击"应用价格"生效</span>
      </div>

      <!-- 商品列表 -->
      <el-table :data="goodsList" border v-loading="loading" class="goods-table">
        <el-table-column prop="name" label="商品名称" min-width="180" />
        <el-table-column prop="categoryTitle" label="分类" width="100">
          <template #default="{ row }">
            <el-tag size="small">{{ row.category && row.category.title || '-' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="price" label="当前售价(元)" width="120">
          <template #default="{ row }">
            <span class="current-price">¥{{ (row.price || 0).toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="采购价(元)" width="180">
          <template #default="{ row }">
            <el-input-number
              v-model="row.costPrice"
              :min="0"
              :precision="2"
              :step="0.1"
              :controls="false"
              size="small"
              placeholder="输入采购价"
              style="width: 140px"
            />
          </template>
        </el-table-column>
        <el-table-column label="预计商户价(元)" width="140">
          <template #default="{ row }">
            <span v-if="row.costPrice > 0" class="sell-price">
              ¥{{ calculateSellPrice(row.costPrice) }}
            </span>
            <span v-else class="no-price">-</span>
          </template>
        </el-table-column>
        <el-table-column label="差额(元)" width="100">
          <template #default="{ row }">
            <span v-if="row.costPrice > 0" class="diff">
              +{{ calculateDiff(row.costPrice) }}
            </span>
            <span v-else>-</span>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="loadGoods"
        @current-change="loadGoods"
        style="margin-top: 20px; text-align: right"
      />
    </el-card>
  </div>
</template>

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

// AI 解析相关
const aiText = ref('')
const aiLoading = ref(false)
const aiResult = ref([])

const updatedCount = computed(() => aiResult.value.filter(i => i.updated).length)

// 计算商户价
const calculateSellPrice = (costPrice) => {
  if (!costPrice || costPrice <= 0) return '0.00'
  return (costPrice * (1 + markupRate.value / 100)).toFixed(2)
}

// 计算差额
const calculateDiff = (costPrice) => {
  if (!costPrice || costPrice <= 0) return '0.00'
  return (costPrice * markupRate.value / 100).toFixed(2)
}

// 加载分类
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

// 加载商品列表
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

// 加载加价比例
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

// AI解析
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
      ElMessage.success(`解析完成：共 ${res.data.total} 个，成功匹配 ${res.data.updated} 个`)
      // 刷新商品列表
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

// 批量更新采购价
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

  await ElMessageBox.confirm(
    `已填写 ${count} 个商品的采购价，确定批量更新吗？`,
    '提示',
    { type: 'warning' }
  )

  try {
    await batchUpdateCostPrice({ updates })
    ElMessage.success('采购价更新成功')
    loadGoods()
  } catch (e) {
    // 用户取消不报错
    if (e !== 'cancel') {
      ElMessage.error('更新失败')
    }
  }
}

// 应用价格
const applyPricing = async () => {
  await ElMessageBox.confirm(
    '确定应用当前采购价生成新的商户价吗？',
    '确认',
    { type: 'warning' }
  )

  try {
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

<style lang="scss" scoped>
.pricing-container {
  padding: 20px;
}

.ai-card {
  margin-bottom: 20px;
  :deep(.el-card__header) {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: #fff;
    padding: 12px 20px;
  }
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.ai-input-area {
  .ai-textarea {
    margin-bottom: 15px;
    :deep(.el-textarea__inner) {
      font-size: 15px;
      line-height: 1.8;
    }
  }
}

.ai-actions {
  display: flex;
  align-items: center;
  gap: 10px;

  .ai-tip {
    color: #999;
    font-size: 12px;
    margin-left: 10px;
  }
}

.ai-result {
  margin-top: 20px;

  .parsed-price {
    color: #67c23a;
    font-weight: bold;
  }

  .success-text {
    color: #67c23a;
  }

  .gray-text {
    color: #999;
  }

  .result-summary {
    margin-top: 15px;
    padding: 10px;
    background: #f5f7fa;
    border-radius: 4px;
    color: #666;
    font-size: 14px;
  }
}

.goods-card {
  :deep(.el-card__header) {
    background: #f5f7fa;
  }
}

.markup-rate {
  font-size: 14px;
  color: #666;
}

.filter-form {
  margin-bottom: 20px;
}

.action-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;

  .tip {
    font-size: 12px;
    color: #999;
  }
}

.goods-table {
  .current-price {
    color: #909399;
  }

  .sell-price {
    color: #67c23a;
    font-weight: bold;
  }

  .no-price {
    color: #c0c4cc;
  }

  .diff {
    color: #e6a23c;
  }
}
</style>
