'use client';

import { useState, useRef, forwardRef, useImperativeHandle } from 'react';
import { Button } from 'antd';
import { useConfigContext } from '../../context/ConfigContext';
import CommentForm, { CommentFormRef } from './CommentForm';
import CommentReply, { CommentReplyRef } from './CommentReply';
import { IComment, ICommentRequest, ICommentReplyRequest } from '../../api/comment';
import dayjs from 'dayjs';
import dynamic from 'next/dynamic';

// 动态导入 Markdown 编辑器
const MDEditor = dynamic(() => import('@uiw/react-md-editor'), {
  ssr: false,
  loading: () => <p>Loading editor...</p>
});

interface CommentPostProps {
  comments: IComment[];
  author: string;
  onSubmit: (req: ICommentRequest) => void;
  onSubmitReply: (req: ICommentReplyRequest, commentId: string) => void;
  onSubmitReply2Reply: (req: ICommentReplyRequest, replyToId: string) => void;
}

export interface CommentPostRef {
  clearReq: () => void;
  clearReplyReq: () => void;
  clearReply2ReplyReq: () => void;
}

const CommentPost = forwardRef<CommentPostRef, CommentPostProps>((props, ref) => {
  const { comments = [], author = "", onSubmit, onSubmitReply, onSubmitReply2Reply } = props;
  const { isBlackMode } = useConfigContext();
  const [activeCommentIndex, setActiveCommentIndex] = useState<string>('');
  
  // 引用表单组件
  const commentFormRef = useRef<CommentFormRef>(null);
  const commentReplyFormRef = useRef<Record<string, CommentFormRef>>({});
  const replyListRefs = useRef<Record<string, CommentReplyRef>>({});

  // 清空主评论表单
  const clearReq = () => {
    if (commentFormRef.current) {
      commentFormRef.current.clearReq();
    }
  };

  // 清空回复评论的表单
  const clearReplyReq = () => {
    if (activeCommentIndex && commentReplyFormRef.current[activeCommentIndex]) {
      commentReplyFormRef.current[activeCommentIndex].clearReq();
      setActiveCommentIndex('');
    }
  };

  // 清空回复的回复表单
  const clearReply2ReplyReq = () => {
    Object.values(replyListRefs.current).forEach(ref => {
      if (ref) {
        ref.clearReplyReq();
      }
    });
  };

  // 暴露方法给父组件
  useImperativeHandle(ref, () => ({
    clearReq,
    clearReplyReq,
    clearReply2ReplyReq
  }));

  // 处理评论提交
  const handleSubmit = (req: ICommentRequest) => {
    onSubmit(req);
  };

  // 处理回复提交
  const handleSubmitReply = (req: ICommentRequest, commentId: string) => {
    const req4Reply: ICommentReplyRequest = {
      ...req
    };
    onSubmitReply(req4Reply, commentId);
  };

  // 处理回复的回复提交
  const handleSubmitReply2Reply = (req: ICommentReplyRequest, replyToId: string) => {
    onSubmitReply2Reply(req, replyToId);
  };

  return (
    <div className="bg-white p-2 dark:bg-gray-800 dark:text-gray-200">
      <div className="h-[50px] leading-[50px] pl-3 border-b border-gray-200 mb-2 flex items-center gap-x-1">
        <span className="text-2xl block text-black dark:text-gray-200">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" className="w-6 h-6">
            <path fillRule="evenodd" d="M4.804 21.644A6.707 6.707 0 006 21.75a6.721 6.721 0 003.583-1.029c.774.182 1.584.279 2.417.279 5.322 0 9.75-3.97 9.75-9 0-5.03-4.428-9-9.75-9s-9.75 3.97-9.75 9c0 2.409 1.025 4.587 2.674 6.192.232.226.277.428.254.543a3.73 3.73 0 01-.814 1.686.75.75 0 00.44 1.223zM8.25 10.875a1.125 1.125 0 100 2.25 1.125 1.125 0 000-2.25zM10.875 12a1.125 1.125 0 112.25 0 1.125 1.125 0 01-2.25 0zm4.875-1.125a1.125 1.125 0 100 2.25 1.125 1.125 0 000-2.25z" clipRule="evenodd" />
          </svg>
        </span>
        <span className="text-2xl">评论</span>
      </div>

      <div className="mb-2 border-b border-gray-200 p-1">
        <CommentForm
          ref={commentFormRef}
          onSubmit={handleSubmit}
        />
      </div>

      {comments.map((comment) => (
        <div className="flex mb-5" key={comment.id}>
          <div className="w-[8%] min-h-[180px] flex justify-center md:w-[8%] sm:w-[15%]">
            {comment.picture ? (
              <img
                src={comment.picture}
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
            <div className="text-gray-500 h-[55px] leading-[35px] flex gap-x-2 sm:leading-[25px]">
              <div className="flex gap-x-2 sm:flex-col">
                {comment.website ? (
                  <a 
                    href={comment.website}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-blue-500 hover:underline"
                  >
                    {comment.username === author 
                      ? `${comment.username}[作者]`
                      : comment.username
                    }
                  </a>
                ) : (
                  <span className="text-blue-500">
                    {comment.username === author 
                      ? `${comment.username}[作者]`
                      : comment.username
                    }
                  </span>
                )}
                <span>
                  发表于 {dayjs(comment.comment_time * 1000).format('YYYY-MM-DD HH:mm:ss')}
                </span>
              </div>
              
              <Button
                className="ml-auto hover:bg-gray-100"
                onClick={() => setActiveCommentIndex(activeCommentIndex === comment.id ? '' : comment.id)}
              >
                回复
              </Button>
            </div>
            
            <div>
              <div data-theme={isBlackMode ? 'dark' : 'light'}>
                <MDEditor.Markdown 
                  source={comment.content}
                  className="markdown-body lg:p-12 p-0"
                />
              </div>
            </div>
            
            {comment.replies && comment.replies.length > 0 && (
              <div>
                <CommentReply
                  replies={comment.replies}
                  author={author}
                  commentId={comment.id}
                  onSubmitReply2Reply={handleSubmitReply2Reply}
                  ref={(el) => {
                    if (el) {
                      replyListRefs.current[comment.id] = el;
                    }
                  }}
                />
              </div>
            )}
            
            {activeCommentIndex === comment.id && (
              <div>
                <CommentForm
                  commentId={comment.id}
                  onSubmit={handleSubmitReply}
                  ref={(el) => {
                    if (el) {
                      commentReplyFormRef.current[comment.id] = el;
                    }
                  }}
                />
                <Button
                  className="mt-2 hover:bg-gray-100"
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
    </div>
  );
});

CommentPost.displayName = 'CommentPost';

export default CommentPost;