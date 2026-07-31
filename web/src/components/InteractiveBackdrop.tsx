"use client";

import React, { useEffect, useRef } from "react";

type PointerState = {
  active: boolean;
  x: number;
  y: number;
  easedX: number;
  easedY: number;
};

type Particle = {
  x: number;
  y: number;
  vx: number;
  vy: number;
  radius: number;
  phase: number;
  speed: number;
  drift: number;
  depth: number;
  alpha: number;
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
    const start = performance.now();
    let particles: Particle[] = [];
    const pointer: PointerState = {
      active: false,
      x: window.innerWidth * 0.5,
      y: window.innerHeight * 0.35,
      easedX: window.innerWidth * 0.5,
      easedY: window.innerHeight * 0.35,
    };

    const createParticles = () => {
      const area = width * height;
      const count = Math.max(56, Math.min(140, Math.floor(area / 10500)));
      particles = Array.from({ length: count }, () => ({
        x: Math.random() * width,
        y: Math.random() * height,
        vx: (Math.random() - 0.5) * 0.52,
        vy: -0.16 - Math.random() * 0.34,
        radius: 1.4 + Math.random() * 3.2,
        phase: Math.random() * Math.PI * 2,
        speed: 0.34 + Math.random() * 0.42,
        drift: 18 + Math.random() * 42,
        depth: 0.18 + Math.random() * 0.82,
        alpha: 0.52 + Math.random() * 0.38,
      }));
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
      createParticles();
    };

    const movePointer = (event: PointerEvent) => {
      pointer.active = true;
      pointer.x = event.clientX;
      pointer.y = event.clientY;
    };

    const leavePointer = () => {
      pointer.active = false;
    };

    const wrap = (value: number, max: number) => {
      if (value < -24) return value + max + 48;
      if (value > max + 24) return value - max - 48;
      return value;
    };

    const draw = (now: number) => {
      const motionScale = reduceMotion.matches ? 0.7 : 1;
      const time = ((now - start) / 1000) * motionScale;
      ctx.clearRect(0, 0, width, height);

      const isDark = document.documentElement.classList.contains("dark");
      const line = isDark ? "rgba(96, 165, 250, 0.08)" : "rgba(37, 99, 235, 0.06)";
      const accent = isDark ? "rgba(20, 184, 166, 0.08)" : "rgba(20, 184, 166, 0.06)";
      const glow = isDark ? "rgba(96, 165, 250, 0.13)" : "rgba(59, 130, 246, 0.1)";
      const particleColor = isDark ? "248, 250, 252" : "29, 78, 216";
      const nearLineColor = isDark ? "226, 232, 240" : "37, 99, 235";
      const particleAlphaBoost = isDark ? 1.35 : 1;
      const lineAlphaBoost = isDark ? 1.6 : 1;

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

      const targetX = pointer.active ? pointer.x : width * 0.52 + Math.sin(time * 0.12) * width * 0.18;
      const targetY = pointer.active ? pointer.y : height * 0.42 + Math.cos(time * 0.14) * height * 0.12;
      pointer.easedX += (targetX - pointer.easedX) * 0.055;
      pointer.easedY += (targetY - pointer.easedY) * 0.055;

      const parallaxX = (pointer.easedX / Math.max(width, 1) - 0.5) * 86;
      const parallaxY = (pointer.easedY / Math.max(height, 1) - 0.5) * 68;
      const particlePositions = particles.map((particle) => {
        particle.x = wrap(particle.x + particle.vx * particle.depth * motionScale, width);
        particle.y = wrap(particle.y + particle.vy * particle.depth * motionScale, height);

        const driftX =
          Math.sin(time * particle.speed + particle.phase) * particle.drift;
        const driftY =
          Math.cos(time * particle.speed * 0.86 + particle.phase * 1.7) *
          particle.drift *
          0.72;
        const naturalX = particle.x + driftX + parallaxX * particle.depth;
        const naturalY = particle.y + driftY + parallaxY * particle.depth;
        const localDx = naturalX - pointer.easedX;
        const localDy = naturalY - pointer.easedY;
        const distance = Math.hypot(localDx, localDy) || 1;
        const influence = pointer.active
          ? Math.max(0, 1 - distance / 380) * 150 * particle.depth * motionScale
          : 0;
        const x = wrap(
          naturalX + (localDx / distance) * influence,
          width
        );
        const y = wrap(
          naturalY + (localDy / distance) * influence,
          height
        );
        return { ...particle, x, y };
      });

      ctx.lineWidth = 1;
      for (let i = 0; i < particlePositions.length; i += 1) {
        const current = particlePositions[i];
        for (let j = i + 1; j < particlePositions.length; j += 1) {
          const next = particlePositions[j];
          const distance = Math.hypot(current.x - next.x, current.y - next.y);
          if (distance > 112) continue;
          const alpha = (1 - distance / 112) * 0.16 * lineAlphaBoost;
          ctx.strokeStyle = `rgba(${nearLineColor}, ${alpha})`;
          ctx.beginPath();
          ctx.moveTo(current.x, current.y);
          ctx.lineTo(next.x, next.y);
          ctx.stroke();
        }
      }

      particlePositions.forEach((particle) => {
        const shimmer = 0.75 + Math.sin(time * 0.8 + particle.phase) * 0.25;
        ctx.save();
        ctx.shadowBlur = (isDark ? 20 : 14) * particle.depth;
        ctx.shadowColor = `rgba(${particleColor}, ${0.32 * particle.alpha * particleAlphaBoost})`;
        ctx.fillStyle = `rgba(${particleColor}, ${Math.min(1, particle.alpha * shimmer * particleAlphaBoost)})`;
        ctx.beginPath();
        ctx.arc(particle.x, particle.y, particle.radius, 0, Math.PI * 2);
        ctx.fill();
        ctx.restore();
      });

      rafId = window.requestAnimationFrame(draw);
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
      className="pointer-events-none fixed inset-0 z-20 opacity-70"
    />
  );
}
