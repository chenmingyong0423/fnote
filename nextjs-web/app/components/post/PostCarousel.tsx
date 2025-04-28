'use client';

import { useRef } from 'react';
import Link from 'next/link';
import { Carousel } from 'antd';
import type { CarouselProps } from 'antd';
import { LeftOutlined, RightOutlined } from '@ant-design/icons';
import { CarouselVO } from '../../api/config';

interface PostCarouselProps {
  carousel?: CarouselVO[] | null;
}

const PostCarousel: React.FC<PostCarouselProps> = ({ carousel }) => {
  // 确保 carousel 是数组类型
  const safeCarousel: CarouselVO[] = Array.isArray(carousel) ? carousel : [];
  
  // 使用 any 类型作为临时解决方案
  const carouselRef = useRef<any>(null);
  
  // 获取环境配置
  const apiHost = process.env.NEXT_PUBLIC_API_HOST || '';
  
  // Carousel 设置
  const settings: CarouselProps = {
    dots: true,
    infinite: true,
    speed: 500,
    slidesToShow: 1,
    slidesToScroll: 1,
    autoplay: true,
    autoplaySpeed: 4000,
    pauseOnHover: true,
    adaptiveHeight: true,
  };

  // 切换到下一张幻灯片
  const nextSlide = () => {
    if (carouselRef.current) {
      carouselRef.current.next();
    }
  };

  // 切换到上一张幻灯片
  const prevSlide = () => {
    if (carouselRef.current) {
      carouselRef.current.prev();
    }
  };

  // 自定义箭头按钮样式
  const arrowStyle = {
    position: 'absolute',
    zIndex: 1,
    top: '50%',
    transform: 'translateY(-50%)',
    width: 40,
    height: 40,
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    background: 'rgba(0,0,0,0.2)',
    borderRadius: '50%',
    cursor: 'pointer',
    color: 'white',
    fontSize: 20,
    transition: 'all 0.3s',
  } as const;

  // 如果没有数据，显示空状态
  if (safeCarousel.length === 0) {
    return (
      <div className="w-full h-[340px] md:h-[340px] relative rounded-md flex items-center justify-center border border-gray-200 bg-white dark:text-gray-300 dark:bg-gray-800 dark:border-gray-700">
        <p className="text-gray-500 dark:text-gray-400">暂无轮播图数据</p>
      </div>
    );
  }

  return (
    <div className="w-full h-[340px] md:h-[340px] relative rounded-md overflow-hidden border border-gray-200 bg-white dark:text-gray-300 dark:bg-gray-800 dark:border-gray-700">
      {/* 左右箭头导航 */}
      <div 
        className="hidden md:flex hover:bg-opacity-60"
        style={{ ...arrowStyle, left: 10 }}
        onClick={prevSlide}
      >
        <LeftOutlined />
      </div>
      
      <div 
        className="hidden md:flex hover:bg-opacity-60"
        style={{ ...arrowStyle, right: 10 }}
        onClick={nextSlide}
      >
        <RightOutlined />
      </div>
      
      {/* 标签 */}
      <span className="z-10 absolute top-0 left-4 bg-[#2db7f5] rounded-b-md text-white text-sm md:text-xs py-[0.2em] px-[0.8em]">
        推荐
      </span>
      
      {/* Ant Design 轮播组件 */}
      <Carousel
        ref={carouselRef}
        {...settings}
        className="h-full carousel-custom"
      >
        {safeCarousel.map((item, index) => (
          <div key={index} className="h-full">
            <Link
              className="relative block w-full h-full slide-up group cursor-pointer ease-linear duration-100 text-center"
              href={`/posts/${item.id}`}
              target="_blank"
              title={item.title}
            >
              <img
                className="w-full h-full object-cover"
                src={item.cover_img?.startsWith('http') ? item.cover_img : `${apiHost}${item.cover_img}`}
                alt={item.title}
              />
              <div
                className="w-[90%] flex flex-col items-center justify-center absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2"
              >
                <div
                  className="text-2xl font-bold md:text-xl"
                  style={{ color: item.color || '#000' }}
                >
                  {item.title}
                </div>
                <div
                  className="text-lg md:text-base"
                  style={{ color: item.color || '#000' }}
                >
                  {item.summary}
                </div>
              </div>
            </Link>
          </div>
        ))}
      </Carousel>
      
      <style jsx global>{`
        .carousel-custom .slick-dots {
          bottom: 20px;
        }
        .carousel-custom .slick-dots li button {
          width: 30px;
          height: 3px;
          background: white;
          opacity: 0.5;
        }
        .carousel-custom .slick-dots li.slick-active button {
          opacity: 1;
        }
        .carousel-custom .slick-dots li button:before {
          display: none;
        }
      `}</style>
    </div>
  );
};

export default PostCarousel;