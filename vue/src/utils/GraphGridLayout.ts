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
    raw: RawItem
    depth?: number
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
    private tierTree: Set<GraphNodeT>[]

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
        this.tierTree = []
    }


    public build() {
        const nodes: GraphNode[] = []
        const links: GraphLink[] = []

        // 数据预处理
        const nodeMap = new Map<string, RawItem>()
        const linksMap = new Map<string, boolean>()
        const lg = new LinkedGroup();
        this.data.forEach(d => {
            const linkKey = `${d.source}__${d.target}`
            nodeMap.set(d.name, d)
            if (!linksMap.has(linkKey)) {
                linksMap.set(linkKey, true) // 去重
                links.push({ source: d.source, target: d.target, value: 1, })
                lg.add({ source: d.source, target: d.target })
            }
        })
        nodeMap.forEach(node => {
            nodes.push({
                name: node.name,
                value: 1,
                raw: node
            })
        })
        console.log('关系链表组', lg.chains);
        console.log('节点映射', nodeMap);

        return {
            nodes,
            links,
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