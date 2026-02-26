<template>
    <div ref="chartRef" class="graph-chart"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import * as echarts from 'echarts'
import type { ECharts, EChartsOption } from 'echarts'
import type { GraphNode, GraphLink } from '@/utils/GraphGridLayout'

/**
 * 组件 Props
 */
const props = defineProps<{
    nodes: GraphNode[]
    links: GraphLink[]
    graphic?: any[]
}>()

const chartRef = ref<HTMLDivElement | null>(null)
let chartInstance: ECharts | null = null

/**
 * 构建图表配置
 */
const getOption = (): EChartsOption => {
    return {
        graphic: [...(props.graphic ?? [])],
        tooltip: {
            trigger: 'item'
        },
        series: [
            {
                type: 'sankey',
                roam: true, // 支持缩放拖拽
                data: props.nodes,
                links: props.links,
                emphasis: {
                    focus: 'adjacency'
                },
                lineStyle: {
                    color: 'gradient',
                    curveness: 0.5
                }
            }
        ]
    }
}

/**
 * 初始化
 */
const initChart = () => {
    if (!chartRef.value) return
    chartInstance = echarts.init(chartRef.value)
    chartInstance.setOption(getOption())
}

/**
 * 响应式更新
 */
watch(
    () => [props.nodes, props.links],
    () => {
        if (chartInstance) {
            chartInstance.setOption(getOption())
        }
    },
    { deep: true }
)

/**
 * resize 监听
 */
const handleResize = () => {
    chartInstance?.resize()
}

onMounted(() => {
    initChart()
    window.addEventListener('resize', handleResize)
})

onBeforeUnmount(() => {
    window.removeEventListener('resize', handleResize)
    chartInstance?.dispose()
})
</script>

<style scoped>
.graph-chart {
    width: 100%;
    height: 100dvh;
}
</style>