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
        <div className="glass-surface flex h-44 items-center justify-center rounded-lg border-dashed text-gray-500 dark:text-gray-400 sm:h-64">
          <span className="rounded-full border border-gray-200 bg-white/80 px-3 py-1.5 text-xs dark:border-gray-700 dark:bg-[#232426] dark:text-gray-300 md:px-4 md:text-sm">
            {hasError ? "网站数据暂时异常" : "暂无轮播图"}
          </span>
        </div>
      ) : (
        <Carousel
          autoplay
          arrows
          className="glass-surface overflow-hidden rounded-lg [&_.slick-dots-bottom]:bottom-3 [&_.slick-dots_li_button]:!h-1.5 [&_.slick-dots_li_button]:!rounded-full"
        >
          {items.map((item) => (
            <a
              key={item.id}
              href={`/posts/${item.id}`}
              target="_blank"
              rel="noopener noreferrer"
              className="group relative block h-44 focus:outline-none sm:h-64"
              tabIndex={0}
            >
              <div className="relative flex h-full w-full items-center justify-center bg-gray-100">
                {item.cover_img && (
                  <Image
                    src={item.cover_img}
                    alt={item.title}
                    fill
                    sizes="(max-width: 768px) 100vw, 66vw"
                    className="object-cover transition-transform duration-700 group-hover:scale-[1.03]"
                  />
                )}
                <div className="absolute inset-x-0 bottom-0 bg-gradient-to-t from-black/70 via-black/35 to-transparent px-4 pb-5 pt-12 text-white md:px-5 md:pb-6">
                  <div className="mb-2 inline-flex rounded-full border border-white/25 bg-white/15 px-2.5 py-1 text-[11px] font-medium text-white/90 backdrop-blur">
                    精选文章
                  </div>
                  <div
                    className="truncate text-base font-bold leading-snug md:text-xl"
                    title={item.title}
                  >
                    {item.title}
                  </div>
                  <div className="mt-1.5 line-clamp-2 text-xs text-white/85 md:text-sm">
                    {item.summary}
                  </div>
                </div>
              </div>
            </a>
          ))}
        </Carousel>
      )}
    </section>
  );
}
