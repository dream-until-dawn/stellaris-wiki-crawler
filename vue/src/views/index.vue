<template>
    <GraphChart :nodes="nodes" :links="links" />
</template>
<script setup lang='ts'>
import { ref } from 'vue';
import GraphChart, { type GraphNode, type GraphLink } from './component/GraphChart.vue'

const nodes: GraphNode[] = [
    { id: '1', name: '节点A' },
    { id: '2', name: '节点B' },
    { id: '3', name: '节点C' }
]

const links: GraphLink[] = [
    { source: '1', target: '2' },
    { source: '2', target: '3' }
]
</script>
<style lang='scss' scoped></style>