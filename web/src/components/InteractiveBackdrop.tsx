"use client";

import React, { useEffect, useRef } from "react";

type PointerState = {
  active: boolean;
  x: number;
  y: number;
};

export default function InteractiveBackdrop() {
  const canvasRef = useRef<HTMLCanvasElement | null>(null);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    const reduceMotion = window.matchMedia("(prefers-reduced-motion: reduce)");

    const ctx = canvas.getContext("2d");
    if (!ctx) return;

    let width = 0;
    let height = 0;
    let rafId = 0;
    let start = performance.now();
    const pointer: PointerState = {
      active: false,
      x: window.innerWidth * 0.5,
      y: window.innerHeight * 0.35,
    };

    const resize = () => {
      const dpr = Math.min(window.devicePixelRatio || 1, 2);
      width = window.innerWidth;
      height = window.innerHeight;
      canvas.width = Math.floor(width * dpr);
      canvas.height = Math.floor(height * dpr);
      canvas.style.width = `${width}px`;
      canvas.style.height = `${height}px`;
      ctx.setTransform(dpr, 0, 0, dpr, 0, 0);
    };

    const movePointer = (event: PointerEvent) => {
      pointer.active = true;
      pointer.x = event.clientX;
      pointer.y = event.clientY;
    };

    const leavePointer = () => {
      pointer.active = false;
    };

    const draw = (now: number) => {
      const time = (now - start) / 1000;
      ctx.clearRect(0, 0, width, height);

      const isDark = document.documentElement.classList.contains("dark");
      const line = isDark ? "rgba(96, 165, 250, 0.08)" : "rgba(37, 99, 235, 0.06)";
      const accent = isDark ? "rgba(20, 184, 166, 0.08)" : "rgba(20, 184, 166, 0.06)";
      const glow = isDark ? "rgba(96, 165, 250, 0.1)" : "rgba(59, 130, 246, 0.08)";

      const grid = width < 640 ? 56 : 72;
      const offset = (time * 10) % grid;

      ctx.lineWidth = 1;
      ctx.strokeStyle = line;
      for (let x = -grid + offset; x < width + grid; x += grid) {
        ctx.beginPath();
        ctx.moveTo(x, 0);
        ctx.lineTo(x - 28, height);
        ctx.stroke();
      }

      ctx.strokeStyle = accent;
      for (let y = -grid + offset * 0.7; y < height + grid; y += grid) {
        ctx.beginPath();
        ctx.moveTo(0, y);
        ctx.lineTo(width, y - 18);
        ctx.stroke();
      }

      const focusX = pointer.active ? pointer.x : width * (0.45 + Math.sin(time * 0.18) * 0.16);
      const focusY = pointer.active ? pointer.y : height * (0.28 + Math.cos(time * 0.2) * 0.08);
      const gradient = ctx.createRadialGradient(focusX, focusY, 0, focusX, focusY, Math.max(width, height) * 0.42);
      gradient.addColorStop(0, glow);
      gradient.addColorStop(0.45, isDark ? "rgba(20, 184, 166, 0.04)" : "rgba(20, 184, 166, 0.035)");
      gradient.addColorStop(1, "rgba(255, 255, 255, 0)");
      ctx.fillStyle = gradient;
      ctx.fillRect(0, 0, width, height);

      if (!reduceMotion.matches) {
        rafId = window.requestAnimationFrame(draw);
      }
    };

    resize();
    window.addEventListener("resize", resize);
    window.addEventListener("pointermove", movePointer, { passive: true });
    window.addEventListener("pointerleave", leavePointer);
    rafId = window.requestAnimationFrame(draw);

    return () => {
      window.cancelAnimationFrame(rafId);
      window.removeEventListener("resize", resize);
      window.removeEventListener("pointermove", movePointer);
      window.removeEventListener("pointerleave", leavePointer);
    };
  }, []);

  return (
    <canvas
      ref={canvasRef}
      aria-hidden="true"
      className="pointer-events-none fixed inset-0 z-0 opacity-45"
    />
  );
}
