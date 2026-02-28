import { LinkedGroup, type GraphNode as GraphNodeT } from './LinkedGroup'

type RawItem = {
    name: string
    classify: string
    technology: string
    tier: string
    source: string
    target: string
}

export type GraphNode = {
    name: string
    value: number
    depth?: number
    raw?: RawItem
    [key: string]: any
}

export type GraphLink = {
    source: string
    target: string
    value: number
}

type LayoutOptions = {
    classifySpacing?: number
    technologySpacing?: number
    tierSpacing?: number
    nodeXSpacing?: number
    nodeYSpacing?: number
}

export class GraphGridLayout {
    private data: RawItem[]
    private opt: Required<LayoutOptions>

    constructor(data: RawItem[], options?: LayoutOptions) {
        this.data = data
        this.opt = {
            classifySpacing: 500, // 分类间距
            technologySpacing: 160, // 子分类间距
            tierSpacing: 220, // 层级间距
            nodeXSpacing: 200, // 节点间距
            nodeYSpacing: 70, // 节点间距
            ...options
        }
    }

    public build() {
        const nodes: GraphNode[] = []
        const links: GraphLink[] = []
        const colorList: string[] = ["#008000", "#d17c00", "#8900ce", "#0078d1", "#d50000", "#800000", "#000080"]

        // 数据预处理
        const nodeMap = new Map<string, RawItem>()
        const linksMap = new Map<string, boolean>()
        const lg = new LinkedGroup();
        this.data.forEach(d => {
            const linkKey = `${d.source}__${d.target}`
            nodeMap.set(d.name, d)
            if (!linksMap.has(linkKey)) {
                linksMap.set(linkKey, true) // 去重
                links.push({ source: d.source, target: d.target, value: 1 })
                lg.add({ source: d.source, target: d.target })
            }
        })
        console.log('关系链表组', lg.chains);
        // 广度优先遍历
        const bfs = (s: Set<GraphNodeT>, pT: number, p: number) => {
            s.forEach(node => {
                const n = nodeMap.get(node.id)
                if (!n) return;
                const t = Number(n.tier)
                nodes.push({
                    name: node.id,
                    value: 1,
                    depth: (t * 5) + p,
                    raw: n,
                    itemStyle: {
                        color: colorList[t]
                    },
                    tooltip: {
                        formatter: (params: any) => {
                            return `名称：${n.name}<br>分类：${n.classify}<br>子分类：${n.technology}<br>层级：${n.tier}<br>间级：${p}<br>来源：${n.source}<br>目标：${n.target}`
                        }
                    }
                })
                if (node.next.size !== 0) bfs(node.next, t, pT === t ? p + 1 : 1)
            })
        }
        lg.chains.forEach(chain => {
            let curNode = lg.getChainsHeads(chain)
            nodes.push({
                name: curNode.id,
                value: 1,
                depth: 0,
                raw: nodeMap.get(curNode.id)!,

            })
            bfs(curNode.next, Number(nodeMap.get(curNode.id)!.tier), 1)
        })
        console.log('节点列表', nodes);

        return {
            nodes,
            links: links.filter(l => nodeMap.has(l.source) && nodeMap.has(l.target)),
            graphic: []
        }
    }

    private buildBackground(
        structure: Map<string, Map<string, Map<number, RawItem[]>>>
    ) {
        const graphics: any[] = []
        let classifyOffsetX = 0
        let classifyIndex = 0

        structure.forEach((techMap, classify) => {
            const height =
                techMap.size * this.opt.technologySpacing

            graphics.push({
                type: 'rect',
                left: classifyOffsetX - 100,
                top: -100,
                shape: {
                    width: this.opt.classifySpacing,
                    height: height + 200
                },
                style: {
                    fill:
                        classifyIndex % 2 === 0
                            ? '#f7f9fc'
                            : '#eef3f9'
                },
                z: -2
            })

            classifyOffsetX += this.opt.classifySpacing
            classifyIndex++
        })

        return graphics
    }
}