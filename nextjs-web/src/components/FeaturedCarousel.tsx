"use client";
import { Carousel } from "antd";
import Image from "next/image";
import React from "react";

export interface FeaturedArticle {
  id: number;
  title: string;
  cover: string;
  description: string;
  link: string;
}

export default function FeaturedCarousel({ articles }: { articles: FeaturedArticle[] }) {
  return (
    <section>
      <Carousel autoplay className="rounded-lg overflow-hidden shadow">
        {articles.map((item) => (
          <div key={item.id} className="relative h-56 sm:h-72 flex items-center justify-center bg-gray-100">
            <Image src={item.cover} alt={item.title} fill className="object-contain" />
            <div className="absolute bottom-0 left-0 right-0 bg-black/50 text-white p-4">
              <a href={item.link} className="text-lg font-bold hover:underline">{item.title}</a>
              <div className="text-sm mt-1">{item.description}</div>
            </div>
          </div>
        ))}
      </Carousel>
    </section>
  );
}
