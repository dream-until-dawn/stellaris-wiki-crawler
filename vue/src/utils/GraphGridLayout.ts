type RawItem = {
    name: string
    classify: string
    technology: string
    tier: string
    source: string
    target: string
}

type GraphNode = {
    id: string
    name: string
    x: number
    y: number
    raw: RawItem
}

type GraphLink = {
    source: string
    target: string
}

type LayoutOptions = {
    classifySpacing?: number
    tierSpacing?: number
    techSpacing?: number
    nodeSpacing?: number
}

export class GraphGridLayout {
    private data: RawItem[]
    private opt: Required<LayoutOptions>

    constructor(data: RawItem[], options?: LayoutOptions) {
        this.data = data
        this.opt = {
            classifySpacing: 500,
            tierSpacing: 220,
            techSpacing: 160,
            nodeSpacing: 70,
            ...options
        }
    }

    public build() {
        const nodes: GraphNode[] = []
        const links: GraphLink[] = []

        // 1️⃣ 建立节点Map
        const nodeMap = new Map<string, RawItem>()
        this.data.forEach(d => {
            nodeMap.set(d.name, d)
            links.push({ source: d.source, target: d.target })
        })

        // 2️⃣ 分层结构：classify -> technology -> tier
        const structure = new Map<
            string,
            Map<string, Map<number, RawItem[]>>
        >()

        nodeMap.forEach(node => {
            const tier = Number(node.tier)

            if (!structure.has(node.classify))
                structure.set(node.classify, new Map())

            const techMap = structure.get(node.classify)!

            if (!techMap.has(node.technology))
                techMap.set(node.technology, new Map())

            const tierMap = techMap.get(node.technology)!

            if (!tierMap.has(tier))
                tierMap.set(tier, [])

            tierMap.get(tier)!.push(node)
        })

        // 3️⃣ 计算坐标
        let classifyOffsetX = 0

        structure.forEach((techMap, classify) => {
            let techIndex = 0

            techMap.forEach((tierMap, technology) => {

                const baseY = techIndex * this.opt.techSpacing

                tierMap.forEach((nodeList, tier) => {

                    const baseX =
                        classifyOffsetX +
                        tier * this.opt.tierSpacing

                    // 居中对齐
                    const totalHeight =
                        (nodeList.length - 1) * this.opt.nodeSpacing

                    const startY =
                        baseY - totalHeight / 2

                    nodeList.forEach((node, index) => {
                        nodes.push({
                            id: node.name,
                            name: node.name,
                            x: baseX,
                            y: startY + index * this.opt.nodeSpacing,
                            raw: node
                        })
                    })
                })

                techIndex++
            })

            classifyOffsetX += this.opt.classifySpacing
        })

        return {
            nodes,
            links,
            graphic: this.buildBackground(structure)
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
                techMap.size * this.opt.techSpacing

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