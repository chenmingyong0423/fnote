'use client';

import { useState, useRef, useEffect } from 'react';
import dayjs from 'dayjs';
import { useConfigContext } from '../context/ConfigContext';
import { Button } from 'antd';

export default function Notice() {
  const config = useConfigContext();
  const [visible, setVisible] = useState(false);
  const [showTheSecondMarquee, setShowTheSecondMarquee] = useState(false);
  
  const marqueeContainerRef = useRef<HTMLDivElement>(null);
  const marqueeContentRef = useRef<HTMLSpanElement>(null);
  const marqueeContent2Ref = useRef<HTMLSpanElement>(null);

  const closeModal = () => {
    setVisible(false);
  };

  const checkMarquee = () => {
    if (marqueeContentRef.current && marqueeContainerRef.current) {
      if (marqueeContentRef.current.offsetWidth > marqueeContainerRef.current.offsetWidth) {
        marqueeContentRef.current.classList.add("marquee-animation");
        setShowTheSecondMarquee(true);
        if (marqueeContent2Ref.current) {
          marqueeContent2Ref.current.classList.add("marquee-animation");
        }
      }
    }
  };

  const stopMarquee = () => {
    if (marqueeContentRef.current) {
      marqueeContentRef.current.style.animationPlayState = "paused";
      if (marqueeContent2Ref.current) {
        marqueeContent2Ref.current.style.animationPlayState = "paused";
      }
    }
  };

  const startMarquee = () => {
    if (marqueeContentRef.current) {
      marqueeContentRef.current.style.animationPlayState = "running";
      if (marqueeContent2Ref.current) {
        marqueeContent2Ref.current.style.animationPlayState = "running";
      }
    }
  };

  useEffect(() => {
    checkMarquee();
  }, []);

  // 从ConfigContext中获取notice数据
  const notice = config.notice_info || { 
    show: true, 
    content: '博客升级改版中，感兴趣的小伙伴可以一起参与开发设计！',
    publish_time: Date.now() / 1000 
  };

  // 如果没有通知或通知设置为不显示，则不渲染组件
  if (!notice || !notice.show) {
    return null;
  }

  return (
    <div>
      <div
        onClick={() => setVisible(true)}
        className="flex items-center dark:text-dtc dark:bg-gray-800 h-[60px] bg-white mb-5 rounded-md cursor-pointer px-5 ease-linear duration-100 md:shadow-md md:hover:-translate-y-2"
      >
        <div className="text-orange-500 text-2xl">🔊</div>
        <div
          ref={marqueeContainerRef}
          className="ml-5 font-bold overflow-hidden whitespace-nowrap w-full"
          onMouseEnter={stopMarquee}
          onMouseLeave={startMarquee}
        >
          <span ref={marqueeContentRef} className="inline-block">
            [{dayjs(notice.publish_time * 1000).format("YYYY-MM-DD")}] {notice.title || '通知'}
          </span>
          {showTheSecondMarquee && (
            <span
              ref={marqueeContent2Ref}
              className="inline-block ml-5"
            >
              {notice.title || '通知'}
            </span>
          )}
        </div>
      </div>

      {/* 模态框 */}
      {visible && (
        <div
          className="fixed z-50 inset-0 bg-black bg-opacity-40 flex items-center justify-center p-4 shadow-lg"
          onClick={closeModal}
        >
          <div
            className="bg-white p-6 rounded-md shadow-lg md:min-w-[400px] max-w-[80%] dark:text-dtc dark:bg-gray-900"
            onClick={(e) => e.stopPropagation()}
          >
            <div className="text-right text-sm text-gray-500 mb-4">
              发布时间: {dayjs(notice.publish_time * 1000).format("YYYY-MM-DD HH:mm:ss")}
            </div>

            <h2 className="text-xl font-bold mb-4">
              {notice.title || '通知'}
            </h2>
            <p className="indent-8 leading-loose">
              {notice.content}
            </p>
            <div className="text-center mt-5">
              <Button
                className="w-[10%] p-2 m-auto mt-5 text-white bg-[#1E80FF]"
                onClick={closeModal}
              >
                关闭
              </Button>
            </div>
          </div>
        </div>
      )}

      <style jsx>{`
        @keyframes marquee {
          0% {
            transform: translateX(0);
          }
          100% {
            transform: translateX(-100%);
          }
        }

        .marquee-animation {
          animation: marquee 15s linear infinite;
          animation-play-state: running;
        }
      `}</style>
    </div>
  );
}