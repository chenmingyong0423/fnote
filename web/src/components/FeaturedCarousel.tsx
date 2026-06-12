"use client";
import { Carousel } from "antd";
import Image from "next/image";
import React from "react";
import type { CarouselItemVO } from "../api/carousel";

export default function FeaturedCarousel({
  items,
  hasError = false,
}: {
  items: CarouselItemVO[];
  hasError?: boolean;
}) {
  return (
    <section>
      {items.length === 0 ? (
        <div className="h-44 sm:h-72 flex items-center justify-center rounded-lg border border-dashed border-gray-200 bg-gray-100 text-gray-500 shadow dark:border-gray-700 dark:bg-[#141414] dark:text-gray-400">
          <span className="rounded-full border border-gray-200 bg-white/70 px-3 md:px-4 py-1.5 text-xs md:text-sm dark:border-gray-700 dark:bg-[#232426] dark:text-gray-300">
            {hasError ? "网站数据暂时异常" : "暂无轮播图"}
          </span>
        </div>
      ) : (
        <Carousel autoplay arrows className="rounded-lg overflow-hidden shadow">
          {items.map((item) => (
            <a
              key={item.id}
              href={`/posts/${item.id}`}
              target="_blank"
              rel="noopener noreferrer"
              className="block relative h-44 sm:h-72 focus:outline-none"
              tabIndex={0}
            >
              <div className="w-full h-full flex items-center justify-center bg-gray-100 relative">
                {/* 图片展示 */}
                {item.cover_img && (
                  <Image src={item.cover_img} alt={item.title} fill sizes="60" className="rounded-lg" priority />
                )}
                {/* 文字遮罩 */}
                <div className="absolute bottom-0 left-0 right-0 bg-black/55 text-white p-3 md:p-4 rounded-b-lg">
                  <div className="text-base md:text-lg font-bold truncate" title={item.title}>{item.title}</div>
                  <div className="text-xs md:text-sm mt-1 line-clamp-2">{item.summary}</div>
                </div>
              </div>
            </a>
          ))}
        </Carousel>
      )}
    </section>
  );
}
