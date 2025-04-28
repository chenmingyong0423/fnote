'use client';

import { useState, useRef, forwardRef, useImperativeHandle } from 'react';
import { List, Avatar, Button, Input, Form, Divider, Typography, Space, Tooltip } from 'antd';
import { UserOutlined, MessageOutlined, SendOutlined, LinkOutlined } from '@ant-design/icons';
import { useConfigContext } from '../../context/ConfigContext';
import { IComment, ICommentRequest, ICommentReplyRequest, IReply } from '../../api/comment';
import dayjs from 'dayjs';
import dynamic from 'next/dynamic';

const { TextArea } = Input;
const { Title, Text, Paragraph } = Typography;

// 动态导入 Markdown 编辑器
const MDEditor = dynamic(() => import('@uiw/react-md-editor'), {
  ssr: false,
  loading: () => <p>Loading editor...</p>
});

interface ListCommentPostProps {
  comments: IComment[];
  author: string;
  onSubmit: (req: ICommentRequest) => void;
  onSubmitReply: (req: ICommentReplyRequest, commentId: string) => void;
  onSubmitReply2Reply: (req: ICommentReplyRequest, commentId: string) => void;
}

export interface ListCommentPostRef {
  clearForm: () => void;
}

interface CommentFormProps {
  onSubmit: (content: string) => void;
  placeholder?: string;
  buttonText?: string;
  onCancel?: () => void;
}

// 评论输入表单
const CommentFormComponent: React.FC<CommentFormProps> = ({ 
  onSubmit, 
  placeholder = '请输入评论内容...',
  buttonText = '提交评论',
  onCancel
}) => {
  const [form] = Form.useForm();
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      setSubmitting(true);
      
      onSubmit(values.content);
      form.resetFields();
      
      setSubmitting(false);
    } catch (error) {
      console.error('Submit failed:', error);
    }
  };

  return (
    <Form form={form} layout="vertical">
      <Form.Item 
        name="content" 
        rules={[{ required: true, message: '请输入评论内容' }]}
      >
        <TextArea 
          rows={4} 
          placeholder={placeholder}
          style={{ resize: 'none' }}
        />
      </Form.Item>
      <Form.Item>
        <Space>
          <Button 
            type="primary" 
            onClick={handleSubmit} 
            loading={submitting}
            icon={<SendOutlined />}
          >
            {buttonText}
          </Button>
          {onCancel && (
            <Button onClick={onCancel}>
              取消
            </Button>
          )}
        </Space>
      </Form.Item>
    </Form>
  );
};

const ListCommentPost = forwardRef<ListCommentPostRef, ListCommentPostProps>((props, ref) => {
  const { comments = [], author = "", onSubmit, onSubmitReply, onSubmitReply2Reply } = props;
  const { isBlackMode } = useConfigContext();
  
  const [activeComment, setActiveComment] = useState<string | null>(null);
  const [activeReply, setActiveReply] = useState<string | null>(null);
  
  // 清空表单
  const clearForm = () => {
    setActiveComment(null);
    setActiveReply(null);
  };
  
  // 暴露方法给父组件
  useImperativeHandle(ref, () => ({
    clearForm
  }));
  
  // 提交主评论
  const handleSubmitComment = (content: string) => {
    if (!content.trim()) return;
    
    const commentReq: ICommentRequest = {
      content,
      username: '当前用户', // 通常应该从用户状态中获取
      email: 'user@example.com',
      website: '',
      postId: ''
    };
    
    onSubmit(commentReq);
  };
  
  // 提交回复
  const handleSubmitReply = (commentId: string, content: string) => {
    if (!content.trim()) return;
    
    const replyReq: ICommentReplyRequest = {
      content,
      username: '当前用户', // 通常应该从用户状态中获取
      email: 'user@example.com',
      website: '',
      postId: '',
      replyToId: commentId
    };
    
    onSubmitReply(replyReq, commentId);
    setActiveComment(null);
  };
  
  // 提交回复的回复
  const handleSubmitReplyToReply = (commentId: string, replyId: string, content: string) => {
    if (!content.trim()) return;
    
    const replyReq: ICommentReplyRequest = {
      content,
      username: '当前用户', // 通常应该从用户状态中获取
      email: 'user@example.com',
      website: '',
      postId: '',
      replyToId: replyId
    };
    
    onSubmitReply2Reply(replyReq, commentId);
    setActiveReply(null);
  };

  // 渲染回复列表
  const renderReplies = (replies: IReply[], commentId: string) => {
    if (!replies || replies.length === 0) return null;
    
    return (
      <List
        itemLayout="horizontal"
        dataSource={replies}
        className="mt-2 pl-6 border-l-2 border-gray-100 dark:border-gray-700"
        renderItem={(reply) => (
          <List.Item
            key={reply.id}
            actions={[
              <Button 
                key="reply-button" 
                type="text" 
                onClick={() => setActiveReply(activeReply === reply.id ? null : reply.id)}
              >
                回复
              </Button>
            ]}
          >
            <List.Item.Meta
              avatar={
                reply.picture ? 
                  <Avatar src={reply.picture} /> : 
                  <Avatar icon={<UserOutlined />} />
              }
              title={
                <Space>
                  {reply.website ? (
                    <a 
                      href={reply.website}
                      target="_blank" 
                      rel="noopener noreferrer"
                      className="text-blue-500 hover:underline"
                    >
                      {reply.name === author ? `${reply.name} [作者]` : reply.name}
                    </a>
                  ) : (
                    <Text strong className="text-blue-500">
                      {reply.name === author ? `${reply.name} [作者]` : reply.name}
                    </Text>
                  )}
                  <Text type="secondary">回复</Text>
                  <Text strong className="text-blue-500">
                    {reply.reply_to === author ? `${reply.reply_to} [作者]` : reply.reply_to}
                  </Text>
                </Space>
              }
              description={
                <Tooltip title={dayjs(reply.reply_time * 1000).format('YYYY-MM-DD HH:mm:ss')}>
                  <Text type="secondary">
                    {dayjs(reply.reply_time * 1000).fromNow()}
                  </Text>
                </Tooltip>
              }
            />
            
            <div className="w-full mt-2">
              {/* 回复内容 */}
              <div data-theme={isBlackMode ? 'dark' : 'light'}>
                <MDEditor.Markdown 
                  source={reply.content} 
                  className="prose dark:prose-invert max-w-none"
                />
              </div>
              
              {/* 引用的原评论内容 */}
              {reply.replied_content && (
                <div className="mt-2 p-3 bg-gray-50 dark:bg-gray-800 border-l-4 border-gray-200 dark:border-gray-700 rounded text-sm">
                  <Text type="secondary">引用内容：</Text>
                  <Paragraph ellipsis={{ rows: 2, expandable: true, symbol: '展开' }} className="mt-1">
                    {reply.replied_content}
                  </Paragraph>
                </div>
              )}
              
              {/* 回复框 */}
              {activeReply === reply.id && (
                <div className="mt-4 ml-8">
                  <CommentFormComponent 
                    onSubmit={(content) => handleSubmitReplyToReply(commentId, reply.id, content)}
                    placeholder="回复评论..."
                    buttonText="回复"
                    onCancel={() => setActiveReply(null)}
                  />
                </div>
              )}
            </div>
          </List.Item>
        )}
      />
    );
  };

  return (
    <div className="bg-white dark:bg-gray-800 rounded-md shadow-sm p-6">
      {/* 评论标题 */}
      <div className="flex items-center gap-2 mb-6 pb-2 border-b border-gray-200 dark:border-gray-700">
        <MessageOutlined style={{ fontSize: 24 }} />
        <Title level={4} style={{ margin: 0 }}>
          评论区
        </Title>
        <Text type="secondary" className="ml-2">
          共 {comments.length} 条评论
        </Text>
      </div>
      
      {/* 评论输入框 */}
      <div className="mb-8">
        <CommentFormComponent onSubmit={handleSubmitComment} />
      </div>
      
      {/* 评论列表 */}
      <List
        itemLayout="vertical"
        size="large"
        dataSource={comments}
        renderItem={(comment) => (
          <List.Item
            key={comment.id}
            actions={[
              <Button 
                key="reply-button" 
                type="text" 
                icon={<MessageOutlined />}
                onClick={() => setActiveComment(activeComment === comment.id ? null : comment.id)}
              >
                回复
              </Button>
            ]}
            className="border-b border-gray-100 dark:border-gray-700 pb-4 mb-4"
          >
            <List.Item.Meta
              avatar={
                comment.picture ? 
                  <Avatar size={40} src={comment.picture} /> : 
                  <Avatar size={40} icon={<UserOutlined />} />
              }
              title={
                <Space>
                  {comment.website ? (
                    <a 
                      href={comment.website} 
                      target="_blank" 
                      rel="noopener noreferrer"
                      className="text-blue-500 hover:underline flex items-center"
                    >
                      <span>{comment.username === author ? `${comment.username} [作者]` : comment.username}</span>
                      <LinkOutlined className="ml-1 text-xs" />
                    </a>
                  ) : (
                    <Text strong className="text-blue-500">
                      {comment.username === author ? `${comment.username} [作者]` : comment.username}
                    </Text>
                  )}
                </Space>
              }
              description={
                <Tooltip title={dayjs(comment.comment_time * 1000).format('YYYY-MM-DD HH:mm:ss')}>
                  <Text type="secondary">
                    {dayjs(comment.comment_time * 1000).fromNow()}
                  </Text>
                </Tooltip>
              }
            />
            
            {/* 评论内容 */}
            <div className="mt-4" data-theme={isBlackMode ? 'dark' : 'light'}>
              <MDEditor.Markdown 
                source={comment.content}
                className="prose dark:prose-invert max-w-none"
              />
            </div>
            
            {/* 回复输入框 */}
            {activeComment === comment.id && (
              <div className="mt-4 ml-12">
                <CommentFormComponent 
                  onSubmit={(content) => handleSubmitReply(comment.id, content)}
                  placeholder="回复评论..."
                  buttonText="回复"
                  onCancel={() => setActiveComment(null)}
                />
              </div>
            )}
            
            {/* 回复列表 */}
            {renderReplies(comment.replies, comment.id)}
          </List.Item>
        )}
      />
      
      {comments.length === 0 && (
        <div className="text-center py-10">
          <MessageOutlined style={{ fontSize: 48, opacity: 0.3 }} />
          <p className="mt-4 text-gray-500 dark:text-gray-400">暂无评论，成为第一个评论的人吧！</p>
        </div>
      )}
      
      <style jsx global>{`
        /* 自定义暗黑模式下的样式 */
        .dark .markdown-body {
          color: rgba(255, 255, 255, 0.85);
        }
        
        .dark .markdown-body a {
          color: #1890ff;
        }
        
        /* 响应式调整 */
        @media (max-width: 576px) {
          .ant-list-item-meta-avatar {
            margin-right: 8px;
          }
        }
      `}</style>
    </div>
  );
});

ListCommentPost.displayName = 'ListCommentPost';

export default ListCommentPost;