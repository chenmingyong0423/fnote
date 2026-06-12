import React from "react";
import FooterRecords from "./FooterRecords";

interface FooterProps {
  websiteRecords: string[];
}

const Footer: React.FC<FooterProps> = ({ websiteRecords }) => (
  <footer className="w-[calc(100%-2rem)] md:w-full bg-white dark:bg-[#141414] border-t border-gray-200 dark:border-[#303030] mt-8 mb-4 mx-auto max-w-7xl flex flex-col items-center rounded-xl p-4">
    <div className="flex flex-wrap items-center justify-center gap-x-1 gap-y-1 w-full text-center text-gray-500 text-xs md:text-sm">
      © {new Date().getFullYear()} Copyright © 2024 - Designed by
      <a
        href="https://github.com/chenmingyong0423/fnote"
        className="text-blue-600 hover:underline"
      >
        Fnote
      </a>
    </div>
    <FooterRecords websiteRecords={websiteRecords} />
  </footer>
);

export default Footer;
