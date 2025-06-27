"use client";
import { useEffect } from "react";

export interface PostSeoProps {
  title: string;
  description?: string;
  keywords?: string;
  coverImg?: string;
}

export default function PostSeoClient({ title, description, keywords, coverImg }: PostSeoProps) {
  useEffect(() => {
    if (title) document.title = title;
    if (description) {
      let meta = document.querySelector('meta[name="description"]');
      if (!meta) {
        meta = document.createElement('meta');
        meta.setAttribute('name', 'description');
        document.head.appendChild(meta);
      }
      meta.setAttribute('content', description);
    }
    if (keywords) {
      let meta = document.querySelector('meta[name="keywords"]');
      if (!meta) {
        meta = document.createElement('meta');
        meta.setAttribute('name', 'keywords');
        document.head.appendChild(meta);
      }
      meta.setAttribute('content', keywords);
    }
    if (coverImg) {
      let meta = document.querySelector('meta[property="og:image"]');
      if (!meta) {
        meta = document.createElement('meta');
        meta.setAttribute('property', 'og:image');
        document.head.appendChild(meta);
      }
      meta.setAttribute('content', coverImg);
    }
  }, [title, description, keywords, coverImg]);
  return null;
}
