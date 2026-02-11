() => {
  const h2List = Array.from(document.querySelectorAll("h2"));
  if (h2List.length === 0) {
    return;
  }
  // 提取 h2 h3
  const extractH2H3 = (h) => {
    return h.querySelector(".mw-headline")?.innerText || "";
  };
  // 提取 p 纯文本
  const extractPureP = (p) => {
    return p.innerText.trim();
  };
  // 提取 p
  const extractP = (p) => {
    let output = "";
    function walk(node) {
      if (node.nodeType === Node.TEXT_NODE) {
        output += node.textContent;
      } else if (node.nodeType === Node.ELEMENT_NODE) {
        if (node.tagName === "IMG") {
          const src = node.getAttribute("src") || "";
          output += "[IMG:" + src + "]";
        }
        // 递归子节点
        node.childNodes.forEach((child) => walk(child));
      }
    }
    p.childNodes.forEach((child) => walk(child));
    return output;
  };
  // 提取表格值 - 适用不同列不同处理
  const extractTableValue = (cell, index) => {
    console.log(index, "extractTableValue", cell);
    switch (index) {
      case 0:
        const img = cell.querySelector("img");
        return `[IMG:${img?.getAttribute("src") || "err"}]`;
      case 1:
        const b = cell.querySelector("b");
        const div = cell.querySelector("div .hidem");
        return {
          name: b?.innerText.trim() || "err",
          description: div?.innerText.trim() || "err",
        };
      case 2:
      case 3:
        return cell.innerText.trim() || "0";
      case 4:
        const lis = cell.querySelectorAll("li");
        if (lis.length === 0) {
          // a - 仅有一个效果
          const img = cell.querySelector("img");
          if (!img) return [];
          return [`[IMG:${img?.getAttribute("src") || "err"}] ${cell.innerText.trim()}`];
        } else {
          // ul -多个效果
        }
        return "";
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
    if (thead) {
      const rows = Array.from(tbody.querySelectorAll("tr"));
      bodys = rows
        .map((row, index) => {
          const cells = Array.from(row.querySelectorAll("td"));
          if (index === 0) {
            // 如果第一行 第一个元素全空,可能就是第一行是表格头
            if (cells[0].innerHtml === "") {
              headers = cells.map((td) => {
                const text = td.innerText.trim();
                return text === "" ? "Icon" : text;
              });
              return null;
            }
          }
          let obj = {};
          cells.forEach((cell, index) => {
            let key = "";
            if (headers.length > 0) {
              key = headers[index] || "col_" + index;
            } else {
              key = "col_" + index;
            }
            obj[key] = extractTableValue(cell, index);
          });
          return obj;
        })
        .filter((item) => item !== null);
    }

    return bodys;
  };

  // 开始执行
  let result = [];
  for (const h2 of h2List) {
    // 主要科技内容
    // console.log("(标题)->h2", h2);
    if (h2.nextElementSibling?.tagName.toLowerCase() !== "p") continue;
    const p = h2.nextElementSibling;
    // console.log("(介绍内容)->h2 后面紧跟的是 p 标签", p);
    if (p.nextElementSibling?.tagName.toLowerCase() !== "table") continue;
    const table = p.nextElementSibling;
    // console.log("(科技列表)->p 后面紧跟的是 table 标签", table);
    result.push({
      title: extractH2H3(h2),
      description: extractPureP(p),
      descriptionRich: extractP(p),
      list: extractTable(table),
      additional: {},
    });
    break;
    // 可能的额外科技
    if (table.nextElementSibling?.tagName.toLowerCase() !== "h3") continue;
    const h3 = table.nextElementSibling;
    // console.log("(额外科技)->table 后面紧跟的是 h3 标签-还需继续检查是否存在 table", h3);
    if (h3.nextElementSibling?.tagName.toLowerCase() !== "p") continue;
    const extraP = h3.nextElementSibling;
    // console.log("(额外科技列表)->h3 后面紧跟的是 p 标签-继续处理数据", extraP);
    if (extraP.nextElementSibling?.tagName.toLowerCase() !== "table") continue;
    const extraTable = extraP.nextElementSibling;
    // console.log("(额外科技列表)->h3 后面紧跟的是 table 标签-继续处理数据", extraTable);
    result[result.length - 1].additional = {
      title: extractH2H3(h3),
      description: extractPureP(extraP),
      descriptionRich: extractP(extraP),
      list: extractTable(extraTable),
    };
  }
  console.log("result", result);
  return result;
};
