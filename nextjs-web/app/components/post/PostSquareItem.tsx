import Image from 'next/image';
import Link from 'next/link';
import dayjs from 'dayjs';
import { IPost } from '../../api/post';

// 定义 props 接口
interface PostSquareItemProps {
  posts: IPost[];
}

export default function PostSquareItem({ posts = [] }: PostSquareItemProps) {
  // 添加默认值防止 posts 为 undefined
  const safeItems = posts || [];
  
  // 在 Next.js 中，我们使用环境变量而不是 useRuntimeConfig
  const baseUrl = process.env.NEXT_PUBLIC_DOMAIN || '';
  const apiHost = process.env.NEXT_PUBLIC_API_HOST || '';
  
  // 如果没有数据，显示空状态
  if (safeItems.length === 0) {
    return (
      <div className="flex justify-center items-center p-10 bg-white rounded-md dark:bg-gray-800">
        <p className="text-gray-500 dark:text-gray-400">暂无文章数据</p>
      </div>
    );
  }
  
  return (
    <div className="flex flex-wrap justify-between">
      {safeItems.map((item, index) => (
        <Link
          key={index}
          href={`${baseUrl}/posts/${item.sug}`}
          target="_blank"
          title={item.title}
          className="item group flex flex-col items-center box-border p-5 bg-white rounded-md w-[49%] md:w-[49%] sm:w-full h-[400px] cursor-pointer duration-100 custom_shadow dark:text-dtc dark:bg-gray-800 mb-5"
        >
          <div className="h-2/3 overflow-hidden relative w-full">
            <Image
              className="object-cover w-full h-full group-hover:scale-110 duration-500"
              src={item.cover_img?.startsWith('http') ? item.cover_img : `${apiHost}${item.cover_img}`}
              alt={item.title}
              width={400}
              height={300}
            />
            <div className="flex flex-wrap gap-y-2 gap-x-3 z-10 w-auto absolute top-3 ease-linear duration-200">
              {item.sticky_weight === 1 && (
                <span className="bg-[#22c55e] rounded-md text-white py-[0.2em] px-[0.8em]">
                  ↑置顶
                </span>
              )}
              {(item.categories || []).map((category, idx) => (
                <span
                  key={idx}
                  className="bg-[#2db7f5] rounded-md text-white py-[0.2em] px-[0.8em]"
                >
                  {category}
                </span>
              ))}
              {(item.tags || []).map((tag, idx) => (
                <span
                  key={idx}
                  className="bg-orange-500 rounded-md text-white py-[0.2em] px-[0.8em]"
                >
                  {tag}
                </span>
              ))}
            </div>
          </div>
          <div className="h-1/3 overflow-hidden relative w-full flex flex-col">
            <div className="mb-2 text-lg h-[60px]">
              {item.title}
            </div>
            <div className="flex-grow">
              <p className="leading-relaxed text-gray-500 truncate">
                {item.summary}
              </p>
            </div>
            <div className="flex gap-x-3 h-[40px] mt-auto">
              <div className="flex gap-x-1 items-center">
                <span>👁️</span><span>{item.visit_count || 0}</span>
              </div>
              <div className="flex gap-x-1 items-center">
                <span>👍</span><span>{item.like_count || 0}</span>
              </div>
              <div className="flex gap-x-1 items-center">
                <span>💬</span><span>{item.comment_count || 0}</span>
              </div>
              <div className="ml-auto flex gap-x-1 items-center">
                <span>{item.created_at ? dayjs(item.created_at * 1000).format("YYYY-MM-DD") : ''}</span>
              </div>
            </div>
          </div>
        </Link>
      ))}

      <style jsx>{`
        .custom_shadow {
          box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
          transition: all 0.3s cubic-bezier(.25,.8,.25,1);
        }
        
        .custom_shadow:hover {
          box-shadow: 0 10px 20px rgba(0, 0, 0, 0.1);
        }
      `}</style>
    </div>
  );
}