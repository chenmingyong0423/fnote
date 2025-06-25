import React from "react";
import FooterRecords from "./FooterRecords";

interface FooterProps {
  websiteRecords: string[];
}

const Footer: React.FC<FooterProps> = ({ websiteRecords }) => (
  <footer className="w-full bg-white border-t border-gray-200 mt-8 mb-4 mx-auto max-w-7xl flex flex-col items-center rounded-xl p-4">
    <div className="flex items-center justify-center w-full text-center text-gray-500 text-sm">
      © {new Date().getFullYear()} Copyright © 2024 - Designed by
      <a
        href="https://github.com/chenmingyong0423/fnote"
        className="ml-1 text-blue-600 hover:underline"
      >
        Fnote
      </a>
    </div>
    <FooterRecords websiteRecords={websiteRecords} />
  </footer>
);

export default Footer;
