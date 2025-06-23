import React from "react";

const Footer: React.FC = () => (
  <footer className="w-full bg-white text-center py-4 text-gray-500 text-sm border-t border-gray-200 mt-8 h-[60px]">
    © {new Date().getFullYear()} Copyright © 2024 - Designed by <a href="https://github.com/chenmingyong0423/fnote">Fnote</a>
  </footer>
);

export default Footer;
