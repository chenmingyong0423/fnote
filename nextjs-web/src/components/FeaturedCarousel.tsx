"use client";
import { Carousel } from "antd";
import Image from "next/image";
import React from "react";
import type { CarouselItemVO } from "../api/carousel";

export default function FeaturedCarousel({ items }: { items: CarouselItemVO[] }) {
  return (
    <section>
      <Carousel autoplay arrows className="rounded-lg overflow-hidden shadow">
        {items.map((item) => (
          <a
            key={item.id}
            href={`/posts/${item.id}`}
            className="block relative h-56 sm:h-72 focus:outline-none"
            tabIndex={0}
          >
            <div className="w-full h-full flex items-center justify-center bg-gray-100 relative">
              {/* 图片展示 */}
              {item.cover_img && (
                <Image src={item.cover_img} alt={item.title} fill sizes="60" className="rounded-lg" priority />
              )}
              {/* 文字遮罩 */}
              <div className="absolute bottom-0 left-0 right-0 bg-black/50 text-white p-4 rounded-b-lg">
                <div className="text-lg font-bold truncate" title={item.title}>{item.title}</div>
                <div className="text-sm mt-1 line-clamp-2">{item.summary}</div>
              </div>
            </div>
          </a>
        ))}
      </Carousel>
    </section>
  );
}
