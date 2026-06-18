"use client";

import { Carousel } from "antd";
import type { CarouselRef } from "antd/es/carousel";
import Image from "next/image";
import React from "react";
import type { CarouselItemVO } from "../api/carousel";

const DEFAULT_TEXT_COLOR = "#ffffff";

function resolveTextAppearance(value: string) {
  const color = /^#[0-9a-f]{6}$/i.test(value) ? value : DEFAULT_TEXT_COLOR;
  const red = Number.parseInt(color.slice(1, 3), 16);
  const green = Number.parseInt(color.slice(3, 5), 16);
  const blue = Number.parseInt(color.slice(5, 7), 16);
  const luminance = (red * 0.299 + green * 0.587 + blue * 0.114) / 255;
  const isDarkText = luminance < 0.48;

  return {
    color,
    isDarkText,
    textShadow: isDarkText
      ? "0 1px 3px rgba(255, 255, 255, 0.9)"
      : "0 1px 3px rgba(0, 0, 0, 0.8)",
  };
}

export default function FeaturedCarousel({
  items,
  hasError = false,
}: {
  items: CarouselItemVO[];
  hasError?: boolean;
}) {
  const carouselRef = React.useRef<CarouselRef | null>(null);
  const carouselWrapperRef = React.useRef<HTMLDivElement | null>(null);
  const lastWheelAtRef = React.useRef(0);

  React.useEffect(() => {
    const wrapper = carouselWrapperRef.current;
    if (!wrapper) return;

    const handleWheel = (event: WheelEvent) => {
      if (window.innerWidth < 768 || items.length <= 1) return;

      const delta =
        Math.abs(event.deltaY) >= Math.abs(event.deltaX)
          ? event.deltaY
          : event.deltaX;
      if (Math.abs(delta) < 12) return;

      event.preventDefault();
      event.stopPropagation();

      const now = window.performance.now();
      if (now - lastWheelAtRef.current < 550) return;

      lastWheelAtRef.current = now;
      if (delta > 0) {
        carouselRef.current?.next();
      } else {
        carouselRef.current?.prev();
      }
    };

    wrapper.addEventListener("wheel", handleWheel, { passive: false });

    return () => {
      wrapper.removeEventListener("wheel", handleWheel);
    };
  }, [items.length]);

  return (
    <section>
      {items.length === 0 ? (
        <div className="glass-surface flex h-44 items-center justify-center rounded-lg border-dashed text-gray-500 dark:text-gray-400 sm:h-64">
          <span className="rounded-full border border-gray-200 bg-white/80 px-3 py-1.5 text-xs dark:border-gray-700 dark:bg-[#232426] dark:text-gray-300 md:px-4 md:text-sm">
            {hasError ? "网站数据暂时异常" : "暂无轮播图"}
          </span>
        </div>
      ) : (
        <div ref={carouselWrapperRef}>
          <Carousel
            ref={carouselRef}
            autoplay
            arrows
            className="glass-surface overflow-hidden rounded-lg [&_.slick-dots-bottom]:bottom-3 [&_.slick-dots_li_button]:!h-1.5 [&_.slick-dots_li_button]:!rounded-full"
          >
            {items.map((item) => {
              const textAppearance = resolveTextAppearance(item.color);
              const overlayClass = textAppearance.isDarkText
                ? "from-white/85 via-white/50"
                : "from-black/75 via-black/40";

              return (
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
                    <div
                      className={`absolute inset-x-0 bottom-0 bg-gradient-to-t ${overlayClass} to-transparent px-4 pb-5 pt-12 md:px-5 md:pb-6`}
                    >
                      <div className="mb-2 inline-flex rounded-full border border-white/20 bg-black/45 px-2.5 py-1 text-[11px] font-medium text-white backdrop-blur">
                        精选文章
                      </div>
                      <div
                        className="truncate text-base font-bold leading-snug md:text-xl"
                        style={{
                          color: textAppearance.color,
                          textShadow: textAppearance.textShadow,
                        }}
                        title={item.title}
                      >
                        {item.title}
                      </div>
                      <div
                        className="mt-1.5 line-clamp-2 text-xs opacity-90 md:text-sm"
                        style={{
                          color: textAppearance.color,
                          textShadow: textAppearance.textShadow,
                        }}
                      >
                        {item.summary}
                      </div>
                    </div>
                  </div>
                </a>
              );
            })}
          </Carousel>
        </div>
      )}
    </section>
  );
}
