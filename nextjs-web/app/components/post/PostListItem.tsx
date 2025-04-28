import Image from 'next/image';
import Link from 'next/link';
import dayjs from 'dayjs';
import { IPost } from '../../api/post';

// Define props interface
interface PostListItemProps {
  posts: IPost[];
}

export default function PostListItem({ posts = [] }: PostListItemProps) {
  // 添加默认值防止 posts 为 undefined
  const safeItems = posts || [];
  
  // 在 Next.js 中使用环境变量代替 useRuntimeConfig
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
    <div>
      {safeItems.map((item, index) => (
        <Link
          key={index}
          href={`${baseUrl}/posts/${item.sug}`}
          target="_blank"
          title={item.title}
          className="slide-up item group flex p-5 bg-white rounded-md h-[200px] cursor-pointer ease-linear duration-100 hover:drop-shadow-xl hover:-translate-y-[8px] dark:text-dtc dark:bg-gray-800 mb-5 relative"
        >
          <div className="w-1/3 overflow-hidden relative">
            <Image
              className="object-contain max-w-full h-full"
              src={item.cover_img?.startsWith('http') ? item.cover_img : `${apiHost}${item.cover_img}`}
              alt={item.title}
              width={200}
              height={150}
            />
            <div className="flex flex-wrap gap-x-3 w-full z-10 absolute top-3 -left-full group-hover:left-[1%] ease-linear duration-200">
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
          <div className="w-2/3 flex flex-col">
            <div className="mb-2 text-xl h-[50px]">
              {item.title}
            </div>
            <div className="flex-grow h-[100px] line-clamp-3">
              <p className="leading-relaxed text-gray-500">{item.summary}</p>
            </div>
            <div className="flex gap-x-3 mt-2 h-[20px]">
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
        .line-clamp-3 {
          display: -webkit-box;
          -webkit-box-orient: vertical;
          -webkit-line-clamp: 3;
          overflow: hidden;
        }

        .item::before {
          content: "";
          position: absolute;
          width: 99%;
          transform: scaleX(0);
          height: 2px;
          bottom: 0;
          left: 0.5%;
          background-color: #0087ca;
          transform-origin: bottom;
          transition: transform 0.3s ease-out;
        }

        .item:hover::before {
          transform: scaleX(1);
          transform-origin: right left;
        }

        @keyframes slideUp {
          0% {
            transform: translateY(100%);
          }
          100% {
            transform: translateY(0);
          }
        }

        .slide-up {
          animation: slideUp 0.5s ease;
          animation-iteration-count: 1;
        }
      `}</style>
    </div>
  );
}