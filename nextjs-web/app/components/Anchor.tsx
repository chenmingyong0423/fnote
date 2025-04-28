'use client'; // 添加客户端渲染标记

import { useRef, useState, useEffect } from 'react';

// 定义 Toc 类型
export type Toc = {
  id: string;
  depth: number;
  text: string;
  children: Toc[];
};

// 定义组件 Props
interface AnchorProps {
  toc: Toc[];
  lineIndex: string;
}

export default function Anchor({ toc, lineIndex }: AnchorProps) {
  const [anchors, setAnchors] = useState<Toc[]>([]);
  const [lineIdx, setLineIdx] = useState<string>('');
  const containerRef = useRef<HTMLDivElement>(null);

  // 将 toc 展开并合并到 anchors 中
  const expandToc = (tocItems: Toc[], depth: number = 1) => {
    const result: Toc[] = [];
    tocItems.forEach((item) => {
      result.push(item);
      if (item.children && item.children.length > 0) {
        result.push(...expandToc(item.children, depth + 1));
      }
    });
    return result;
  };

  // 监听 toc 变化，更新 anchors
  useEffect(() => {
    if (toc && toc.length > 0) {
      setAnchors(expandToc(toc, 1));
    } else {
      setAnchors([]);
    }
  }, [toc]);

  // 处理点击事件
  const handleClick = (toc: Toc) => {
    setLineIdx(toc.id);
  };

  // 监听 lineIndex 变化
  useEffect(() => {
    if (lineIndex !== lineIdx) {
      setLineIdx(lineIndex);

      // 在状态更新后处理滚动操作
      setTimeout(() => {
        const container = containerRef.current;
        if (!container) return;

        const activeElement = container.querySelector(".active");
        if (activeElement) {
          const elementPosition = activeElement.getBoundingClientRect().top - container.getBoundingClientRect().top;
          
          // 判断元素是否超过当前 container 高度的一半
          if (elementPosition > container.clientHeight / 2) {
            container.scrollTop = elementPosition - container.clientHeight / 2 + container.scrollTop;
          } else {
            container.scrollTop = 0;
          }
        }
      }, 0);
    }
  }, [lineIndex, lineIdx]);

  return (
    <div
      ref={containerRef}
      className="flex flex-col bg-white rounded-md max-h-[700px] overflow-y-auto"
    >
      <div
        className="leading-6 text-lg border-b border-gray-200 pb-1.5 dark:text-white p-2 pl-4 pt-5"
      >
        目录
      </div>
      {anchors.map((anchor, index) => (
        <div
          key={index}
          style={{
            padding: `10px 0 10px ${anchor.depth * anchor.depth * anchor.depth}px`,
          }}
        >
          <a
            style={{ cursor: 'pointer' }}
            className={`pl-6 ${
              anchor.id === lineIdx
                ? 'anchor_border text-base text-[#1e80ff] font-bold active'
                : ''
            }`}
            href={`#${anchor.id}`}
            onClick={() => handleClick(anchor)}
          >
            {anchor.text}
          </a>
        </div>
      ))}

      <style jsx>{`
        .anchor_border {
          border-left: 2px solid #1e80ff;
        }
      `}</style>
    </div>
  );
}