export type GraphNodeT = {
    id: string;
    next: Set<GraphNodeT>;
    prev: Set<GraphNodeT>;
}

export class GraphNode {
    id: string;
    next: Set<GraphNodeT>;
    prev: Set<GraphNodeT>;
    constructor(id: string) {
        this.id = id;
        this.next = new Set();   // 指向谁
        this.prev = new Set();   // 被谁指向 
    }
}

export class LinkedGroup {
    nodes: Map<string, GraphNode>;
    chains: Set<GraphNode>[];
    constructor() {
        this.nodes = new Map();   // id -> GraphNode
        this.chains = [];         // 多条链（数组形式保存）
    }

    /**
     * 添加一条有向边 { source, target }
     */
    add({ source, target }: { source: string, target: string }) {
        const s = this._getOrCreateNode(source);
        const t = this._getOrCreateNode(target);

        // 1️⃣ 防止成环
        if (this._hasPath(t, s)) {
            throw new Error("检测到成环，拒绝添加");
        }

        // 2️⃣ 建立连接
        if (!s.next.has(t)) {
            s.next.add(t);
            t.prev.add(s);
        }

        // 3️⃣ 重建链分组
        this._rebuildChains();
    }

    /**
     * 获取或创建节点
     */
    _getOrCreateNode(id: string): GraphNode {
        if (!this.nodes.has(id)) {
            this.nodes.set(id, new GraphNode(id));
        }
        return this.nodes.get(id)!;
    }

    /**
     * DFS 检查是否存在路径（用于防环）
     */
    _hasPath(start: GraphNode, target: GraphNode, visited = new Set()) {
        if (start === target) return true;
        visited.add(start);

        for (let next of start.next) {
            if (!visited.has(next)) {
                if (this._hasPath(next, target, visited)) {
                    return true;
                }
            }
        }
        return false;
    }

    /**
     * 重建所有链
     */
    _rebuildChains() {
        this.chains = [];
        const visited: Set<GraphNode> = new Set();

        for (let node of this.nodes.values()) {
            if (!visited.has(node)) {
                const chain: Set<GraphNode> = new Set();
                this._dfsCollect(node, chain, visited);
                this.chains.push(chain);
            }
        }

        this._mergeChainsIfPossible();
    }

    /**
     * DFS 收集连通块
     */
    _dfsCollect(node: GraphNode, chain: Set<GraphNode>, visited: Set<GraphNode>) {
        visited.add(node);
        chain.add(node);

        for (let n of node.next) {
            if (!visited.has(n)) {
                this._dfsCollect(n, chain, visited);
            }
        }
        for (let p of node.prev) {
            if (!visited.has(p)) {
                this._dfsCollect(p, chain, visited);
            }
        }
    }

    /**
     * 自动合并首尾可连接链
     */
    _mergeChainsIfPossible() {
        let merged = true;

        while (merged) {
            merged = false;

            outer:
            for (let i = 0; i < this.chains.length; i++) {
                for (let j = i + 1; j < this.chains.length; j++) {
                    if (this._canMerge(this.chains[i]!, this.chains[j]!)) {
                        const mergedChain = new Set([
                            ...this.chains[i]!,
                            ...this.chains[j]!
                        ]);

                        this.chains.splice(j, 1);
                        this.chains.splice(i, 1);
                        this.chains.push(mergedChain);

                        merged = true;
                        break outer;
                    }
                }
            }
        }
    }

    /**
     * 判断两个链是否可以首尾连接
     */
    _canMerge(chainA: Set<GraphNode>, chainB: Set<GraphNode>) {
        for (let nodeA of chainA) {
            for (let nodeB of chainB) {
                if (nodeA.next.has(nodeB) || nodeB.next.has(nodeA)) {
                    return true;
                }
            }
        }
        return false;
    }

    /**
     * 获取当前链信息
     */
    getChains() {
        return this.chains.map(chain =>
            [...chain].map(n => ({
                id: n.id,
                next: [...n.next].map(n => n.id),
                prev: [...n.prev].map(n => n.id)
            }))
        );
    }

    // 给出链的头部
    getChainsHeads(chain: Set<GraphNode>): GraphNode {
        for (let node of chain) {
            if (node.prev.size === 0) {
                return node;
            }
        }
        return null!;
    }
}