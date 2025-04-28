'use client';

import { useState, useRef, useImperativeHandle, forwardRef } from 'react';
import { Button } from 'antd';
import { ICommentRequest } from '../../api/comment';
import { useConfigContext } from '../../context/ConfigContext';
import { MD5 } from 'crypto-js';
import { showToast } from '../MyToast';
import dynamic from 'next/dynamic';

// 动态导入 Markdown 组件，避免服务端渲染问题
const MDEditor = dynamic(() => import('@uiw/react-md-editor'), { 
  ssr: false,
  loading: () => <p>Loading editor...</p>
});

// 邮箱验证函数
const isValidEmail = (email: string): boolean => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
};

interface CommentFormProps {
  commentId?: string;
  onSubmit: (comment: ICommentRequest, commentId: string) => void;
}

export interface CommentFormRef {
  clearReq: () => void;
}

const CommentForm = forwardRef<CommentFormRef, CommentFormProps>((props, ref) => {
  const { commentId = "", onSubmit } = props;
  const { isBlackMode } = useConfigContext();
  
  const [showUsernameTip, setShowUsernameTip] = useState(false);
  const [showEmailTip, setShowEmailTip] = useState(false);
  const [showWebsiteTip, setShowWebsiteTip] = useState(false);
  const [pic, setPic] = useState<string>("");
  const [isPreview, setIsPreview] = useState(false);
  
  const [commentReq, setCommentReq] = useState<ICommentRequest>({
    postId: "",
    username: "",
    email: "",
    website: "",
    content: ""
  });
  
  // 表单清空函数
  const clearReq = () => {
    setCommentReq({
      postId: "",
      username: "",
      email: "",
      website: "",
      content: ""
    });
    setIsPreview(false);
    setPic("");
  };
  
  // 暴露清空方法给父组件
  useImperativeHandle(ref, () => ({
    clearReq
  }));
  
  // 根据邮箱计算 Gravatar 头像
  const calculateMD54Email = () => {
    setShowEmailTip(false);
    if (commentReq.email !== "") {
      const hashedEmail = MD5(commentReq.email.trim().toLowerCase()).toString();
      setPic(`https://1.gravatar.com/avatar/${hashedEmail}`);
    } else {
      setPic("");
    }
  };
  
  // 提交评论
  const submit = () => {
    if (commentReq.content === "") {
      showToast("评论内容不能为空！", "error");
      return;
    } else if (commentReq.username === "") {
      showToast("昵称不能为空！", "error");
      return;
    } else if (commentReq.email === "") {
      showToast("邮箱不能为空！", "error");
      return;
    }
    
    if (commentReq.website !== "" && !commentReq.website?.startsWith("https://")) {
      showToast("个人站点格式不正确！", "error");
      return;
    }
    
    if (!isValidEmail(commentReq.email)) {
      showToast("邮箱格式不正确！", "error");
      return;
    }
    
    // 深拷贝请求对象
    const deepCopyReq: ICommentRequest = JSON.parse(JSON.stringify(commentReq));
    onSubmit(deepCopyReq, commentId);
  };
  
  // 更新字段的处理函数
  const handleInputChange = (field: keyof ICommentRequest, value: string) => {
    setCommentReq(prev => ({
      ...prev,
      [field]: value
    }));
  };
  
  return (
    <div className="py-2">
      <div>
        {!isPreview ? (
          <textarea
            rows={10}
            className="w-full border border-gray-300 bg-[#F9F9F9] outline-none focus:border-[#1E80FF] rounded-md p-2 box-border mb-3 dark:text-gray-300 dark:bg-gray-800"
            value={commentReq.content}
            onChange={(e) => handleInputChange('content', e.target.value)}
            maxLength={200}
          />
        ) : (
          <div>
            <div className="font-bold flex items-center gap-x-2">
              <span className="block text-black w-8 h-8">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" className="w-6 h-6">
                  <path d="M12 15a3 3 0 100-6 3 3 0 000 6z" />
                  <path fillRule="evenodd" d="M1.323 11.447C2.811 6.976 7.028 3.75 12.001 3.75c4.97 0 9.185 3.223 10.675 7.69.12.362.12.752 0 1.113-1.487 4.471-5.705 7.697-10.677 7.697-4.97 0-9.186-3.223-10.675-7.69a1.762 1.762 0 010-1.113zM17.25 12a5.25 5.25 0 11-10.5 0 5.25 5.25 0 0110.5 0z" clipRule="evenodd" />
                </svg>
              </span>
              <span>预览中...</span>
            </div>
            <br />
            <div data-theme={isBlackMode ? 'dark' : 'light'} className="markdown-preview">
              <div className="markdown-body lg:p-12 p-0 max-h-[170px] h-[170px] overflow-x-auto whitespace-nowrap border border-gray-300 p-2 rounded-md">
                <MDEditor.Markdown source={commentReq.content} />
              </div>
            </div>
          </div>
        )}
      </div>
      
      <div>
        <div className="h-10 text-center leading-10">个人信息</div>
        <div className="flex justify-between items-center w-full mb-3 md:flex-row flex-col md:gap-y-0 gap-y-2">
          <div className="w-[5%] md:block flex justify-center">
            {pic ? (
              <img
                src={pic}
                alt=""
                className="w-auto h-12 rounded-full cursor-pointer hover:rotate-[360deg] ease-out duration-1000"
              />
            ) : (
              <div className="w-12 h-12 rounded-full text-gray-400 flex items-center justify-center">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" className="w-12 h-12">
                  <path fillRule="evenodd" d="M18.685 19.097A9.723 9.723 0 0021.75 12c0-5.385-4.365-9.75-9.75-9.75S2.25 6.615 2.25 12a9.723 9.723 0 003.065 7.097A9.716 9.716 0 0012 21.75a9.716 9.716 0 006.685-2.653zm-12.54-1.285A7.486 7.486 0 0112 15a7.486 7.486 0 015.855 2.812A8.224 8.224 0 0112 20.25a8.224 8.224 0 01-5.855-2.438zM15.75 9a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0z" clipRule="evenodd" />
                </svg>
              </div>
            )}
          </div>
          
          <div className="md:w-[30%] w-full relative">
            <input
              type="text"
              placeholder="* 昵称"
              value={commentReq.username}
              onChange={(e) => handleInputChange('username', e.target.value)}
              className="w-full outline-none border border-gray-300 bg-[#F9F9F9] focus:border-[#1E80FF] rounded-md p-2 box-border dark:text-gray-300 dark:bg-gray-800"
              onFocus={() => setShowUsernameTip(true)}
              onBlur={() => setShowUsernameTip(false)}
            />
            {showUsernameTip && (
              <span className="popup-text absolute bg-[#555] left-[25%] bottom-[125%] text-white p-2 rounded-md h-[20px] animated-fadeIn">
                您的昵称？
              </span>
            )}
          </div>
          
          <div className="md:w-[30%] w-full relative">
            <input
              type="text"
              placeholder="* 邮箱"
              value={commentReq.email}
              onChange={(e) => handleInputChange('email', e.target.value)}
              className="w-full outline-none border border-gray-300 bg-[#F9F9F9] focus:border-[#1E80FF] rounded-md p-2 box-border dark:text-gray-300 dark:bg-gray-800"
              onFocus={() => setShowEmailTip(true)}
              onBlur={calculateMD54Email}
            />
            {showEmailTip && (
              <span className="popup-text absolute bg-[#555] left-[25%] bottom-[125%] text-white p-2 rounded-md h-[20px] animated-fadeIn">
                用于接收通知。
              </span>
            )}
          </div>
          
          <div className="md:w-[30%] w-full relative">
            <input
              type="text"
              placeholder="个人站点，以 https:// 开头"
              value={commentReq.website}
              onChange={(e) => handleInputChange('website', e.target.value)}
              className="w-full outline-none border border-gray-300 bg-[#F9F9F9] focus:border-[#1E80FF] rounded-md p-2 box-border dark:text-gray-300 dark:bg-gray-800"
              onFocus={() => setShowWebsiteTip(true)}
              onBlur={() => setShowWebsiteTip(false)}
            />
            {showWebsiteTip && (
              <span className="popup-text absolute bg-[#555] left-[25%] bottom-[125%] text-white p-2 rounded-md h-[20px] animated-fadeIn">
                不许打广告哟！
              </span>
            )}
          </div>
        </div>
        
        <div className="flex gap-x-2 justify-center">
          <Button
            type="primary"
            style={{ width: '60px', height: '32px', lineHeight: '32px', background: '#1E80FF' }}
            onClick={clearReq}
          >
            清空
          </Button>
          <Button
            type="primary"
            style={{ width: '60px', height: '32px', lineHeight: '32px', background: '#1E80FF' }}
            onClick={() => setIsPreview(!isPreview)}
          >
            {isPreview ? '编辑' : '预览'}
          </Button>
          <Button
            type="primary"
            style={{ width: '60px', height: '32px', lineHeight: '32px', background: '#1E80FF' }}
            onClick={submit}
          >
            提交
          </Button>
        </div>
      </div>
      
      <style jsx>{`
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

        .popup-text:before {
          content: "";
          position: absolute;
          left: 50%;
          bottom: -20px;
          border-width: 10px;
          border-style: solid;
          border-color: #555 transparent transparent transparent;
          transform: translateX(-50%);
        }

        @keyframes fadeIn {
          from {
            opacity: 0;
          }
          to {
            opacity: 1;
          }
        }

        .animated-fadeIn {
          opacity: 0;
          animation: fadeIn 0.5s ease-in-out forwards;
        }
      `}</style>
    </div>
  );
});

CommentForm.displayName = 'CommentForm';

export default CommentForm;