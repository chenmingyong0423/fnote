'use client';

import { useState, useRef, forwardRef, useImperativeHandle } from 'react';
import { Comment, List, Avatar, Button, Form, Input, Tooltip } from 'antd';
import { useConfigContext } from '../../context/ConfigContext';
import { IComment, ICommentRequest, ICommentReplyRequest } from '../../api/comment';
import dayjs from 'dayjs';
import dynamic from 'next/dynamic';
import { UserOutlined, MessageOutlined } from '@ant-design/icons';

// 引入 TextArea 组件
const { TextArea } = Input;

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
  clearForm: () => void;
}

// 评论表单组件
const CommentEditor = ({ onChange, onSubmit, submitting, value, onCancel = undefined }) => (
  <div>
    <Form.Item>
      <TextArea rows={4} onChange={onChange} value={value} />
    </Form.Item>
    <Form.Item>
      <Button htmlType="submit" loading={submitting} onClick={onSubmit} type="primary" style={{ marginRight: 8 }}>
        提交评论
      </Button>
      {onCancel && (
        <Button onClick={onCancel}>
          取消
        </Button>
      )}
    </Form.Item>
  </div>
);

const AntdCommentPost = forwardRef<CommentPostRef, CommentPostProps>((props, ref) => {
  const { comments = [], author = "", onSubmit, onSubmitReply, onSubmitReply2Reply } = props;
  const { isBlackMode } = useConfigContext();
  
  // 评论状态管理
  const [commentValue, setCommentValue] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [replyStates, setReplyStates] = useState<Record<string, { open: boolean, value: string }>>({});
  
  // 清空所有表单
  const clearForm = () => {
    setCommentValue('');
    setReplyStates({});
  };
  
  // 暴露方法给父组件
  useImperativeHandle(ref, () => ({
    clearForm
  }));
  
  // 主评论提交
  const handleCommentSubmit = () => {
    if (!commentValue.trim()) return;
    
    setSubmitting(true);
    
    const commentReq: ICommentRequest = {
      content: commentValue,
      username: '当前用户', // 通常应该从用户状态中获取
      email: 'user@example.com',
      website: '',
      postId: ''
    };
    
    onSubmit(commentReq);
    
    // 模拟提交后的状态
    setTimeout(() => {
      setSubmitting(false);
      setCommentValue('');
    }, 500);
  };
  
  // 回复评论
  const handleReplySubmit = (commentId: string) => {
    const replyState = replyStates[commentId];
    if (!replyState || !replyState.value.trim()) return;
    
    const replyReq: ICommentReplyRequest = {
      content: replyState.value,
      username: '当前用户',
      email: 'user@example.com',
      website: '',
      postId: '',
      replyToId: commentId
    };
    
    onSubmitReply(replyReq, commentId);
    
    // 更新状态
    setReplyStates((prev) => ({
      ...prev,
      [commentId]: { open: false, value: '' }
    }));
  };
  
  // 回复的回复
  const handleReplyToReply = (commentId: string, replyId: string) => {
    const replyState = replyStates[replyId];
    if (!replyState || !replyState.value.trim()) return;
    
    const replyReq: ICommentReplyRequest = {
      content: replyState.value,
      username: '当前用户',
      email: 'user@example.com',
      website: '',
      postId: '',
      replyToId: replyId
    };
    
    onSubmitReply2Reply(replyReq, commentId);
    
    // 更新状态
    setReplyStates((prev) => ({
      ...prev,
      [replyId]: { open: false, value: '' }
    }));
  };
  
  // 切换回复框状态
  const toggleReply = (id: string) => {
    setReplyStates((prev) => ({
      ...prev,
      [id]: { open: !(prev[id]?.open || false), value: prev[id]?.value || '' }
    }));
  };
  
  // 处理回复内容变化
  const handleReplyChange = (id: string, value: string) => {
    setReplyStates((prev) => ({
      ...prev,
      [id]: { open: prev[id]?.open || false, value }
    }));
  };
  
  // 渲染嵌套回复
  const renderReplies = (comment: IComment) => {
    if (!comment.replies || comment.replies.length === 0) return null;
    
    return comment.replies.map((reply) => (
      <Comment
        key={reply.id}
        author={
          <span>
            {reply.website ? (
              <a 
                href={reply.website}
                target="_blank"
                rel="noopener noreferrer"
                style={{ color: '#1E80FF' }}
              >
                {reply.name === author ? `${reply.name}[作者]` : reply.name}
              </a>
            ) : (
              <span style={{ color: '#1E80FF' }}>
                {reply.name === author ? `${reply.name}[作者]` : reply.name}
              </span>
            )}
            <span style={{ marginLeft: 8, marginRight: 8 }}>回复</span>
            <span style={{ color: '#1E80FF' }}>
              {reply.reply_to === author ? `${reply.reply_to}[作者]` : reply.reply_to}
            </span>
          </span>
        }
        avatar={
          reply.picture ? 
            <Avatar src={reply.picture} /> : 
            <Avatar icon={<UserOutlined />} />
        }
        content={
          <div className={isBlackMode ? 'dark' : ''}>
            <MDEditor.Markdown source={reply.content} />
            {reply.replied_content && (
              <div className="mt-2 p-2 bg-gray-50 dark:bg-gray-700 rounded text-sm">
                <div className="text-gray-500 dark:text-gray-400">原内容：</div>
                <div>{reply.replied_content}</div>
              </div>
            )}
          </div>
        }
        datetime={
          <Tooltip title={dayjs(reply.reply_time * 1000).format('YYYY-MM-DD HH:mm:ss')}>
            <span>{dayjs(reply.reply_time * 1000).fromNow()}</span>
          </Tooltip>
        }
        actions={[
          <span key="reply" onClick={() => toggleReply(reply.id)}>回复</span>
        ]}
      >
        {replyStates[reply.id]?.open && (
          <CommentEditor
            onChange={(e) => handleReplyChange(reply.id, e.target.value)}
            onSubmit={() => handleReplyToReply(comment.id, reply.id)}
            submitting={false}
            value={replyStates[reply.id]?.value || ''}
            onCancel={() => toggleReply(reply.id)}
          />
        )}
      </Comment>
    ));
  };

  return (
    <div className="bg-white p-4 dark:bg-gray-800 dark:text-gray-200 rounded shadow-sm">
      <div className="mb-6 flex items-center gap-2 border-b pb-2">
        <MessageOutlined style={{ fontSize: '24px' }} />
        <span className="text-2xl">评论</span>
      </div>
      
      {/* 评论输入表单 */}
      <div className="mb-6">
        <CommentEditor
          onChange={(e) => setCommentValue(e.target.value)}
          onSubmit={handleCommentSubmit}
          submitting={submitting}
          value={commentValue}
        />
      </div>
      
      {/* 评论列表 */}
      <List
        className="comment-list"
        header={`${comments.length} 条评论`}
        itemLayout="horizontal"
        dataSource={comments}
        renderItem={(comment) => (
          <Comment
            key={comment.id}
            author={
              <span>
                {comment.website ? (
                  <a 
                    href={comment.website}
                    target="_blank"
                    rel="noopener noreferrer"
                    style={{ color: '#1E80FF' }}
                  >
                    {comment.username === author ? `${comment.username}[作者]` : comment.username}
                  </a>
                ) : (
                  <span style={{ color: '#1E80FF' }}>
                    {comment.username === author ? `${comment.username}[作者]` : comment.username}
                  </span>
                )}
              </span>
            }
            avatar={
              comment.picture ? 
                <Avatar src={comment.picture} /> : 
                <Avatar icon={<UserOutlined />} />
            }
            content={
              <div className={isBlackMode ? 'dark' : ''}>
                <MDEditor.Markdown source={comment.content} />
              </div>
            }
            datetime={
              <Tooltip title={dayjs(comment.comment_time * 1000).format('YYYY-MM-DD HH:mm:ss')}>
                <span>{dayjs(comment.comment_time * 1000).fromNow()}</span>
              </Tooltip>
            }
            actions={[
              <span key="reply" onClick={() => toggleReply(comment.id)}>回复</span>
            ]}
          >
            {replyStates[comment.id]?.open && (
              <CommentEditor
                onChange={(e) => handleReplyChange(comment.id, e.target.value)}
                onSubmit={() => handleReplySubmit(comment.id)}
                submitting={false}
                value={replyStates[comment.id]?.value || ''}
                onCancel={() => toggleReply(comment.id)}
              />
            )}
            {renderReplies(comment)}
          </Comment>
        )}
      />
      
      <style jsx global>{`
        .markdown-body {
          box-sizing: border-box;
          min-width: 200px;
          max-width: 980px;
          margin: 0 auto;
        }
        
        .dark .markdown-body a {
          color: rgba(255, 255, 255, 0.7) !important;
        }
        
        .ant-comment-content-detail {
          padding-right: 16px;
        }
      `}</style>
    </div>
  );
});

AntdCommentPost.displayName = 'AntdCommentPost';

export default AntdCommentPost;