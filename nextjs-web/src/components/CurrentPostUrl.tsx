"use client";
import React from "react";
import { usePathname } from "next/navigation";

const baseHost = process.env.NEXT_PUBLIC_BASE_HOST || process.env.BASE_HOST || "";

export const CurrentPostUrl: React.FC = () => {
  const pathname = usePathname();
  if (!baseHost) return null;
  return (
    <span className="break-all">{`${baseHost}${pathname}`}</span>
  );
};
