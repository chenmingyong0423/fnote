'use client';

import { useState, useRef, forwardRef, useImperativeHandle, useMemo } from 'react';
import { Button } from 'antd';
import { IReply, ICommentReplyRequest } from '../../api/comment';
import { useConfigContext } from '../../context/ConfigContext';
import dayjs from 'dayjs';
import CommentForm, { CommentFormRef } from './CommentForm';
import dynamic from 'next/dynamic';

// 动态导入 Markdown 组件，避免服务端渲染问题
const MDEditor = dynamic(() => import('@uiw/react-md-editor'), { 
  ssr: false,
  loading: () => <p>Loading editor...</p>
});

interface CommentReplyProps {
  replies: IReply[];
  author: string;
  commentId: string;
  onSubmitReply2Reply: (req: ICommentReplyRequest, commentId: string) => void;
}

export interface CommentReplyRef {
  clearReplyReq: () => void;
}

const CommentReply = forwardRef<CommentReplyRef, CommentReplyProps>((props, ref) => {
  const { replies: originalReplies = [], author = "", commentId, onSubmitReply2Reply } = props;
  const { isBlackMode } = useConfigContext();
  const [activeCommentIndex, setActiveCommentIndex] = useState<string>('');
  const replyFormRef = useRef<Record<string, CommentFormRef>>({});

  // 处理回复中的引用内容
  const replies = useMemo(() => {
    return originalReplies.map((rpy) => {
      const newRpy = { ...rpy }; // 创建一个新对象
      if (newRpy.reply_to_id !== "") {
        const replied = originalReplies.find((r) => r.id === newRpy.reply_to_id);
        if (replied) {
          newRpy.replied_content = replied.content;
        }
      }
      return newRpy;
    });
  }, [originalReplies]);

  // 清空回复表单
  const clearReplyReq = () => {
    if (activeCommentIndex && replyFormRef.current[activeCommentIndex]) {
      replyFormRef.current[activeCommentIndex].clearReq();
      setActiveCommentIndex('');
    }
  };

  // 提交回复
  const submitReply = (req: any, replyToId: string) => {
    const req4Reply: ICommentReplyRequest = {
      ...req,
      replyToId: replyToId,
    };
    onSubmitReply2Reply(req4Reply, commentId);
  };

  // 暴露方法给父组件
  useImperativeHandle(ref, () => ({
    clearReplyReq
  }));

  return (
    <>
      {replies.map((rpy) => (
        <div className="flex mb-5" key={rpy.id}>
          <div className="w-[8%] min-h-[100px] flex justify-center md:w-[8%] sm:w-[15%]">
            {rpy.picture ? (
              <img
                src={rpy.picture}
                alt=""
                className="w-12 h-12 rounded-full cursor-pointer hover:rotate-[360deg] ease-out duration-1000 lg:mr-0"
              />
            ) : (
              <div className="w-12 h-12 rounded-full lg:mr-0 text-gray-400 flex items-center justify-center">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" className="w-12 h-12">
                  <path fillRule="evenodd" d="M18.685 19.097A9.723 9.723 0 0021.75 12c0-5.385-4.365-9.75-9.75-9.75S2.25 6.615 2.25 12a9.723 9.723 0 003.065 7.097A9.716 9.716 0 0012 21.75a9.716 9.716 0 006.685-2.653zm-12.54-1.285A7.486 7.486 0 0112 15a7.486 7.486 0 015.855 2.812A8.224 8.224 0 0112 20.25a8.224 8.224 0 01-5.855-2.438zM15.75 9a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0z" clipRule="evenodd" />
                </svg>
              </div>
            )}
          </div>
          
          <div className="w-[91%] ml-[1%] flex flex-col md:w-[91%] sm:w-[84%]">
            <div className="text-gray-400 h-[55px] leading-[35px] flex gap-x-2 sm:leading-[25px]">
              <div className="flex gap-x-2 sm:flex-col">
                <div className="flex gap-x-2 sm:gap-x-1">
                  {rpy.website ? (
                    <a 
                      href={rpy.website}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-[#1E80FF] sm:truncate"
                    >
                      {rpy.name === author 
                        ? `${rpy.name}[作者]`
                        : rpy.name
                      }
                    </a>
                  ) : (
                    <span className="text-[#1E80FF] sm:truncate">
                      {rpy.name === author 
                        ? `${rpy.name}[作者]`
                        : rpy.name
                      }
                    </span>
                  )}
                  <span>回复</span>
                  <span className="text-[#1E80FF] sm:truncate">
                    {rpy.reply_to === author
                      ? `${rpy.reply_to}[作者]`
                      : rpy.reply_to
                    }
                  </span>
                </div>
                <span>
                  发表于 {dayjs(rpy.reply_time * 1000).format('YYYY-MM-DD HH:mm:ss')}
                </span>
              </div>
              
              <Button
                className="w-[60px] h-8 leading-8 hover:bg-gray-100 ml-auto"
                onClick={() => setActiveCommentIndex(activeCommentIndex === rpy.id ? '' : rpy.id)}
              >
                回复
              </Button>
            </div>
            
            <div>
              <div data-theme={isBlackMode ? 'dark' : 'light'}>
                <MDEditor.Markdown 
                  source={rpy.content}
                  className="markdown-body lg:p-12 p-0"
                />
              </div>
            </div>
            
            {rpy.replied_content && (
              <div className="bg-gray-100 h-10 leading-10 pl-4 truncate dark:text-gray-300 dark:bg-gray-800">
                {rpy.replied_content}
              </div>
            )}
            
            {activeCommentIndex === rpy.id && (
              <div>
                <CommentForm
                  commentId={rpy.id}
                  onSubmit={submitReply}
                  ref={(el) => {
                    if (el) {
                      replyFormRef.current[rpy.id] = el;
                    }
                  }}
                />
                <Button
                  className="w-[60px] h-8 leading-8 hover:bg-gray-100 mx-auto block"
                  onClick={() => setActiveCommentIndex('')}
                >
                  取消
                </Button>
              </div>
            )}
          </div>
        </div>
      ))}
      
      <style jsx global>{`
        .markdown-body {
          box-sizing: border-box;
          min-width: 200px;
          max-width: 980px;
          margin: 0 auto;
          padding: 45px;
        }
        
        @media (max-width: 767px) {
          .markdown-body {
            padding: 15px;
          }
        }
        
        .markdown-body a {
          color: black !important;
        }
        
        .dark .markdown-body a {
          color: rgba(255, 255, 255, 0.7) !important;
        }
      `}</style>
    </>
  );
});

CommentReply.displayName = 'CommentReply';

export default CommentReply;