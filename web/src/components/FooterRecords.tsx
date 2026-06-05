"use client";
import React from "react";

interface FooterRecordsProps {
  websiteRecords: string[];
}

const FooterRecords: React.FC<FooterRecordsProps> = ({ websiteRecords }) => {
  if (!websiteRecords.length) return null;
  return (
    <div className="flex flex-wrap justify-center gap-x-4 gap-y-2 text-xs text-gray-500">
      {websiteRecords.map((item, idx) => (
        <span key={idx} dangerouslySetInnerHTML={{ __html: item }} />
      ))}
    </div>
  );
};

export default FooterRecords;
