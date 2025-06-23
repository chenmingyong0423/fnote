"use client";
import { Carousel } from "antd";
import Image from "next/image";
import React from "react";
import type { CarouselItemVO } from "../api/carousel";

export default function FeaturedCarousel({ items }: { items: CarouselItemVO[] }) {
  return (
    <section>
      <Carousel autoplay className="rounded-lg overflow-hidden shadow">
        {items.map((item) => (
          <div key={item.id} className="relative h-56 sm:h-72 flex items-center justify-center bg-gray-100">
            {/* 图片展示 */}
            {item.cover_img && (
              <Image src={item.cover_img} alt={item.title} fill sizes="(max-width: 768px) 100vw, 100vw" className="object-fill rounded-lg" />
            )}
            {/* 文字遮罩 */}
            <div className="absolute bottom-0 left-0 right-0 bg-black/50 text-white p-4 rounded-b-lg">
              <div className="text-lg font-bold truncate" title={item.title}>{item.title}</div>
              <div className="text-sm mt-1 line-clamp-2">{item.summary}</div>
            </div>
          </div>
        ))}
      </Carousel>
    </section>
  );
}
