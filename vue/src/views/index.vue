<template>
    <template v-if="links.length > 0">
        <SankeyChart :nodes="nodes" :links="links" :graphic="graphic" />
    </template>
</template>
<script setup lang='ts'>
import { onMounted, ref } from 'vue'
import SankeyChart from './component/SankeyChart.vue'
import { GraphGridLayout } from '@/utils/GraphGridLayout'
import type { GraphNode, GraphLink } from '@/utils/GraphGridLayout'

const links = ref<GraphLink[]>([])
const nodes = ref<GraphNode[]>([])
const graphic = ref<any[]>([])


onMounted(async () => {
    const res = await fetch('/test_links.json')
    const raw = await res.json()

    const layout = new GraphGridLayout(raw)
    const result = layout.build()

    nodes.value = result.nodes
    links.value = result.links
    graphic.value = result.graphic
})

</script>
<style lang='scss' scoped></style>