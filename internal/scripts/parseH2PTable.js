() => {
  // 提取 h2 h3
  const extractH2H3 = (h) => {
    return h.querySelector(".mw-headline")?.innerText || "";
  };
  // 提取目标节点的内容->输出富文本
  const extractRichText = (node) => {
    const walk = (current) => {
      if (!current) return "";
      // 1️⃣ 文本节点
      if (current.nodeType === Node.TEXT_NODE) {
        return current.nodeValue || "";
      }
      // 2️⃣ 元素节点
      if (current.nodeType === Node.ELEMENT_NODE) {
        // br
        if (current.tagName === "BR") {
          return "\n";
        }
        // 块级元素
        if (["P", "DIV"].includes(current.tagName)) {
          let result = "";
          current.childNodes.forEach((child) => {
            result += walk(child);
          });
          return result + "\n";
        }
        // 遇到 IMG
        if (current.tagName === "IMG") {
          const src = current.getAttribute("src") || "";
          return `[IMG:${src}]`;
        }
        // 其他元素
        let result = "";
        // 深度优先遍历子节点（保证顺序）
        current.childNodes.forEach((child) => {
          result += walk(child);
        });
        return result;
      }
      return "";
    };
    const output = walk(node);
    const isVaild = !/^[\n\t]+$/g.test(output);
    return isVaild ? output : "";
  };
  // 提取表格值 - 适用不同列不同处理
  const extractTableValue = (cell, colName) => {
    // console.log(colName, cell);
    switch (colName) {
      case "Icon": // 图标
        const img = cell.querySelector("img");
        return `[IMG:${img?.getAttribute("src") || "err"}]`;
      case "Technology": // 科技名称
        const b = cell.querySelector("b");
        const div = cell.querySelector("div .hidem");
        return {
          name: b?.innerText.trim() || "err",
          description: div?.innerText.trim() || "err",
        };
      case "Cost": // 成本
      case "Tier": // 等级
        return cell.innerText.trim() || "0";
      case "Effects / unlocks": // 效果
      case "Prerequisites": // 前置需求
      case "Empire": // 帝国
      case "DLC": // dlc
        const lis = cell.querySelectorAll("li");
        if (lis.length === 0) {
          // 没有使用ul - 仅有一个
          const text = extractRichText(cell).trim();
          return text ? [text] : [];
        } else {
          // 有多个
          return Array.from(lis)
            .map((li) => extractRichText(li).trim())
            .filter(Boolean);
        }
      case "Draw weight": // 权重
        // 提取出li
        const weightLis = cell.querySelectorAll("li");
        // 拿到无展开情况的复制元素
        const copyEl = cell.cloneNode(true);
        copyEl.querySelectorAll(".mw-collapsible").forEach((el) => el.remove());
        const text = copyEl.innerText.trim();

        if (weightLis.length === 0) {
          // 没有使用ul - 仅有一个
          return text ? [text] : [];
        } else {
          // 有多个效果
          return (text ? [text] : []).concat(
            Array.from(weightLis)
              .map((li) => extractRichText(li).trim())
              .filter(Boolean)
          );
        }
      default:
        return cell.innerText.trim();
    }
  };
  // 提取表格
  const extractTable = (table) => {
    // 表格头
    const thead = table.querySelector("thead");
    let headers = [];
    if (thead) {
      headers = Array.from(thead.querySelectorAll("th")).map((th) => {
        const text = th.innerText.trim();
        return text === "" ? "Icon" : text;
      });
    }

    // 表格体
    const tbody = table.querySelector("tbody");
    let bodys = [];
    if (tbody) {
      const rows = Array.from(tbody.querySelectorAll("tr"));
      bodys = rows
        .map((row) => {
          const ths = Array.from(row.querySelectorAll("th"));
          if (ths.length !== 0) {
            headers = Array.from(ths).map((th) => {
              const text = th.innerText.trim();
              return text === "" ? "Icon" : text;
            });
            return null; // 跳过表头
          }

          const cells = Array.from(row.querySelectorAll("td"));
          let obj = {};
          cells.forEach((cell, index) => {
            let key = "";
            if (headers.length > 0) {
              key = headers[index] || "col_" + index;
            } else {
              key = "col_" + index;
            }
            obj[key] = extractTableValue(cell, headers[index] || "unknown");
          });
          return obj;
        })
        .filter((item) => item !== null);
    }

    return bodys;
  };
  // -----------------------------------------------------------------------------------------------
  // 开始执行
  const tableList = Array.from(document.querySelectorAll("table"));
  if (tableList.length === 0) {
    return;
  }
  let dataResult = [];
  let result = [];
  for (const table of tableList) {
    const p = table.previousElementSibling;
    if (p?.tagName?.toLowerCase() !== "p") continue;
    let data = {
      title: null,
      description: [p],
      table: table,
    };
    let previous = p?.previousElementSibling;
    let previousTagName = previous?.tagName?.toLowerCase() || "null";
    if (previousTagName === "null") continue; // 没有上一个元素直接跳过本次

    while (previousTagName !== "null" && previousTagName !== "h2" && previousTagName !== "h3") {
      data.description.unshift(previous);
      previous = previous?.previousElementSibling;
      previousTagName = previous?.tagName?.toLowerCase() || "null";
    }
    data.title = previous;
    dataResult.push(data);
    // 开始处理data
    result.push({
      title: extractH2H3(previous),
      description: data.description.map((node) => extractRichText(node)),
      table: extractTable(table),
    });
  }
  console.log("dataResult", dataResult);
  console.log("result", result);
  let count = 0;
  for (const item of result) {
    count += item?.table?.length || 0;
  }
  console.log("result count", count);
  return result;
};
