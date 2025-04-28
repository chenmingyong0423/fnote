'use client';

import { useConfigContext } from '../context/ConfigContext';
import { showToast } from './MyToast';

export default function ExternalLink() {
  const config = useConfigContext();
  
  const copyExternalLink = async (content: string) => {
    await navigator.clipboard.writeText(content);
    showToast(content, 'success');
  };
  
  return (
    <div className="flex items-center justify-center gap-x-3 text-gray-500 text-xl pt-5 w-full">
      {config.social_info_list.map((icon, index) => (
        <div
          key={index}
          className="group w-12 h-12 bg-white rounded-full flex items-center overflow-hidden md:hover:w-[200px] cursor-pointer duration-500 dark:text-dtc dark:bg-gray-800"
        >
          <div className="bg-white rounded-full dark:text-dtc dark:bg-gray-800 p-3 w-10 flex items-center justify-center shadow-2xl shadow-black group-hover:bg-[#1E80FF] group-hover:text-white duration-300 dark:group-hover:bg-[#1E80FF]/50">
            {icon.link ? (
              <a
                className={icon.icon}
                href={icon.link}
                target="_blank"
                rel="noopener noreferrer"
              >
                {/* 使用图标 */}
                {!icon.icon.includes('i-') && <span>{icon.name.charAt(0)}</span>}
              </a>
            ) : (
              <span
                className={icon.icon}
                onClick={() => copyExternalLink(`${icon.name}: ${icon.link || ''}`)}
              >
                {/* 使用图标 */}
                {!icon.icon.includes('i-') && <span>{icon.name.charAt(0)}</span>}
              </span>
            )}
          </div>
          <div className="text-base w-full text-center truncate select-none dark:text-dtc">
            {icon.link ? (
              <a href={icon.link} target="_blank" rel="noopener noreferrer">
                {icon.name}
              </a>
            ) : (
              <span onClick={() => copyExternalLink(`${icon.name}: ${icon.link || ''}`)}>
                {icon.name}
              </span>
            )}
          </div>
        </div>
      ))}
    </div>
  );
}