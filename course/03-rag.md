# 第 3 课：RAG 检索增强

[课程目录](README.md) · [完整阶段讲义](../stage3_rag/README.md)

## 学习目标

- 理解 embedding、余弦相似度、Top-K 和切块。
- 实现内存向量索引与关键词检索。
- 使用 RRF 融合语义与关键词排名。

## 课程主线

RAG 的核心不是“接一个向量数据库”，而是把资料切成合适的检索单元，可靠地找出相关片段，再让模型基于这些片段回答。本课从数学基础一直做到可供 Agent 调用的搜索工具。

## 配套练习

| # | 练习 | 重点 |
|---|---|---|
| 1 | [余弦相似度](../stage3_rag/ex01_cosine/) | 向量归一化与零向量 |
| 2 | [Chunking](../stage3_rag/ex02_chunking/) | UTF-8、重叠与边界 |
| 3 | [向量索引](../stage3_rag/ex03_vector_index/) | Top-K、稳定排序 |
| 4 | [搜索工具](../stage3_rag/ex04_search_tool/) | 检索能力工具化 |
| 5 | [混合检索](../stage3_rag/ex05_hybrid_rrf/) | BM25、向量与 RRF |

## 完成标准

能够从零实现一个小规模、可测试的混合检索系统，并将它作为 Agent 工具暴露。
