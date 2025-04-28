'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { List, Avatar, Typography, Space, Tooltip, Card } from 'antd';
import { UserOutlined, CommentOutlined, EditOutlined } from '@ant-design/icons';
import { getLatestComments } from '../../api/comment';
import type { ILatestComment } from '../../api/comment';
import type { IResponse, IListData } from '../../utils/http';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';

// 配置 dayjs 相对时间插件
dayjs.extend(relativeTime);

const { Title, Text, Paragraph } = Typography;

const ListCommentHome = () => {
  const router = useRouter();
  const [comments, setComments] = useState<ILatestComment[]>([]);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    const fetchComments = async () => {
      try {
        setLoading(true);
        const postRes = await getLatestComments();
        if (postRes && postRes.data) {
          const commentsData = postRes.data as IListData<ILatestComment>;
          setComments(commentsData.list || []);
        }
      } catch (error) {
        console.error('获取最新评论失败:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchComments();
  }, []);

  const handlePostClick = (postId: string) => {
    router.push(`/posts/${postId}`);
  };

  return (
    <Card 
      title={
        <Space>
          <CommentOutlined />
          <span>最新评论</span>
        </Space>
      }
      bordered={false}
      className="shadow-sm"
    >
      <List
        itemLayout="horizontal"
        dataSource={comments}
        loading={loading}
        locale={{ emptyText: '暂无评论' }}
        renderItem={(comment) => (
          <List.Item>
            <List.Item.Meta
              avatar={
                comment.picture ? 
                  <Avatar src={comment.picture} className="transform transition-all duration-1000 hover:rotate-[360deg]" /> : 
                  <Avatar icon={<UserOutlined />} className="transform transition-all duration-1000 hover:rotate-[360deg]" />
              }
              title={<Text strong>{comment.name}</Text>}
              description={
                <Space direction="vertical" size={2} style={{ width: '100%' }}>
                  <Tooltip title={dayjs(comment.created_at * 1000).format('YYYY-MM-DD HH:mm:ss')}>
                    <Text type="secondary" className="text-xs">
                      {dayjs(comment.created_at * 1000).fromNow()}
                    </Text>
                  </Tooltip>
                  <Paragraph ellipsis={{ rows: 2 }} className="mb-1">
                    {comment.content}
                  </Paragraph>
                  <div 
                    onClick={() => handlePostClick(comment.post_id)}
                    className="flex items-center gap-1 text-blue-500 hover:text-blue-600 cursor-pointer text-sm p-1 hover:bg-gray-50 dark:hover:bg-white/10 rounded transition-colors"
                  >
                    <EditOutlined />
                    <span>{comment.post_title}</span>
                  </div>
                </Space>
              }
            />
          </List.Item>
        )}
      />
    </Card>
  );
};

export default ListCommentHome;