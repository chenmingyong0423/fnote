'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { getLatestComments } from '../../api/comment';
import type { ILatestComment } from '../../api/comment';
import type { IResponse, IListData } from '../../utils/http';
import dayjs from 'dayjs';

const CommentHome = () => {
  const router = useRouter();
  const [comments, setComments] = useState<ILatestComment[]>([]);

  useEffect(() => {
    const fetchComments = async () => {
      try {
        const postRes = await getLatestComments();
        if (postRes && postRes.data) {
          const commentsData = postRes.data as IListData<ILatestComment>;
          setComments(commentsData.list || []);
        }
      } catch (error) {
        console.log(error);
      }
    };

    fetchComments();
  }, []);

  return (
    <div className="flex flex-col bg-white p-5 rounded-md dark:bg-gray-800 dark:text-gray-300">
      <div className="text-lg border-b border-b-gray-200 px-1 py-2">
        最新评论
      </div>
      {comments.map((item, index) => (
        <div
          key={index}
          className="border-b border-b-gray-200 pb-2 p-2"
        >
          <div className="flex h-15 my-2">
            <div>
              {item.picture ? (
                <img
                  src={item.picture}
                  alt=""
                  className="w-12 h-12 rounded-full cursor-pointer hover:rotate-[360deg] ease-out duration-1000 lg:mr-0"
                />
              ) : (
                <div className="w-12 h-12 rounded-full lg:mr-0 text-gray-400">
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" className="w-12 h-12">
                    <path fillRule="evenodd" d="M18.685 19.097A9.723 9.723 0 0021.75 12c0-5.385-4.365-9.75-9.75-9.75S2.25 6.615 2.25 12a9.723 9.723 0 003.065 7.097A9.716 9.716 0 0012 21.75a9.716 9.716 0 006.685-2.653zm-12.54-1.285A7.486 7.486 0 0112 15a7.486 7.486 0 015.855 2.812A8.224 8.224 0 0112 20.25a8.224 8.224 0 01-5.855-2.438zM15.75 9a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0z" clipRule="evenodd" />
                  </svg>
                </div>
              )}
            </div>
            <div className="flex flex-col items-start ml-3">
              <span className="text-base">{item.name}</span>
              <span className="text-gray-500 text-xs">
                {dayjs(item.created_at * 1000).format('YYYY-MM-DD HH:mm:ss')}
              </span>
            </div>
          </div>
          <div>
            <div className="py-1 truncate">
              {item.content}
            </div>
            <div
              className="flex gap-2 items-center text-gray-500 py-1 cursor-pointer hover:bg-green-100 dark:hover:bg-white/20 duration-100"
              onClick={() => router.push(`/posts/${item.post_id}`)}
            >
              <span>
                <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                </svg>
              </span>
              <span>{item.post_title}</span>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
};

export default CommentHome;