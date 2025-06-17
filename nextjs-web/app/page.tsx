'use client';

import { useEffect, useState, useRef } from 'react';
import Notice from './components/Notice';
import Profile from './components/Profile';
import ExternalLink from './components/ExternalLink';
import PostListItem from './components/post/PostListItem';
import PostSquareItem from './components/post/PostSquareItem';
import PostCarousel from './components/post/PostCarousel';
import { useConfigContext } from './context/ConfigContext';
import { IPost } from './api/post';
import { getPosts, getLatestPosts } from './api/post';
import { GetCarousel, CarouselVO } from './api/config';
import { getLatestComments, ILatestComment } from './api/comment';
import { log } from 'console';

// 用于跟踪API调用状态的全局对象
const API_REQUEST_STATUS = {
  latestPosts: false,
  latestComments: false,
  carousel: false
};


export default function Home() {
  const [isMobile, setIsMobile] = useState<boolean>(false);
  const [latestPosts, setLatestPosts] = useState<IPost[]>([]);
  const [carouselData, setCarouselData] = useState<CarouselVO[]>([]);
  const [latestComments, setLatestComments] = useState<ILatestComment[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  
  // 使用 ref 来防止重复调用API
  const dataFetchedRef = useRef(false);
  
  const config = useConfigContext();
  const apiHost = process.env.NEXT_PUBLIC_API_HOST || '';
  
  // 设置响应式布局检测 - 使用媒体查询更可靠
  useEffect(() => {
    // 使用媒体查询检测移动设备
    const mobileQuery = window.matchMedia('(max-width: 767px)');
    
    // 初始检查
    setIsMobile(mobileQuery.matches);
    
    // 添加媒体查询事件监听
    const handleMobileChange = (event: MediaQueryListEvent) => {
      setIsMobile(event.matches);
    };
    
    // 现代浏览器使用 addEventListener
    if (mobileQuery.addEventListener) {
      mobileQuery.addEventListener('change', handleMobileChange);
      return () => mobileQuery.removeEventListener('change', handleMobileChange);
    } 
    // 旧版浏览器兼容性支持
    else {
      mobileQuery.addListener(handleMobileChange);
      return () => mobileQuery.removeListener(handleMobileChange);
    }
  }, []);

  // 获取数据
  useEffect(() => {
    // 如果已经获取过数据，不再重复获取
    if (dataFetchedRef.current) return;
    
    const fetchData = async () => {
      setLoading(true);
      let hasError = false;
      // 获取轮播图数据
      if (!API_REQUEST_STATUS.carousel) {
        try {
          API_REQUEST_STATUS.carousel = true;
          const carouselResponse = await GetCarousel();
          if (carouselResponse?.data) {
            setCarouselData(carouselResponse.data || []);
          }
        } catch (error) {
          hasError = true;
          setCarouselData([]);
          console.error('获取轮播图数据失败:', error);
        }
      }
      // 获取最新文章
      if (!API_REQUEST_STATUS.latestPosts) {
        try {
          API_REQUEST_STATUS.latestPosts = true;
          console.log('Fetching latest posts...');
          const latestResponse = await getLatestPosts();
          if (latestResponse?.data) {
            // 处理 cover_img 路径
            const serverHost = process.env.SERVER_HOST || '';
            const posts = (latestResponse.data.list || []).map((post: IPost) => ({
              ...post,
              cover_img: post.cover_img && !post.cover_img.startsWith('http')
                ? serverHost + post.cover_img
                : post.cover_img
            }));
            setLatestPosts(posts);
          }
        } catch (error) {
          hasError = true;
          setLatestPosts(mockPosts);
          console.error('获取最新文章失败:', error);
        }
      }
      // 获取最新评论
      if (!API_REQUEST_STATUS.latestComments) {
        try {
          API_REQUEST_STATUS.latestComments = true;
          console.log('Fetching latest comments...');
          const commentsResponse = await getLatestComments();
          if (commentsResponse?.data) {
            setLatestComments(commentsResponse.data || []);
            console.log('Latest comments fetched successfully');
          }
        } catch (error) {
          hasError = true;
          setLatestComments([]);
          console.error('获取最新评论失败:', error);
        }
      }
      // 标记数据已获取
      dataFetchedRef.current = true;
      setLoading(false);
      if (hasError) {
        // 可以在这里做全局错误提示
      }
    };
    
    fetchData();
    
    // 组件卸载时重置状态
    return () => {
      API_REQUEST_STATUS.latestPosts = false;
      API_REQUEST_STATUS.latestComments = false;
      API_REQUEST_STATUS.carousel = false;
    };
  }, []);
  
  // 模拟数据
  const mockPosts: IPost[] = [];

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-[#03080c]">
      {/* 主内容 */}
      <div className="container mx-auto px-4 py-8">
        {/* 移动端上使用单列布局 */}
        <div className="flex flex-col md:flex-row gap-8">
          {/* 左侧垂直布局 - 轮播图、公告、文章列表 */}
          <main className="w-full md:w-2/3">
            {/* 轮播图区域 */}
            <div className="mb-8">
              <PostCarousel carousel={carouselData} />
            </div>
            
            {/* 公告区域 */}
            <div className="mb-8">
              <Notice />
            </div>
            
            {/* 最新文章标题 */}
            <div className="mb-4 flex justify-between items-center">
              <h2 className="text-2xl font-bold dark:text-white">最新文章</h2>
            </div>
            
            {/* 最新文章列表 */}
            {loading ? (
              <div className="flex justify-center items-center h-64">
                <p className="text-gray-500 dark:text-gray-400">加载中...</p>
              </div>
            ) : (
              <>
                {!isMobile ? (
                  <PostListItem posts={latestPosts.length > 0 ? latestPosts : mockPosts} />
                ) : (
                  <PostSquareItem posts={latestPosts.length > 0 ? latestPosts : mockPosts} />
                )}
              </>
            )}
          </main>

          {/* 右侧布局 - 仅在非移动端显示 */}
          {!isMobile && (
            <aside className="w-full md:w-1/3">
              {/* 个人资料卡片 */}
              <div className="mb-6">
                <Profile />
              </div>
              
              {/* 社交链接 */}
              <div className="mb-6">
                <ExternalLink />
              </div>
              
              {/* 最新评论 */}
              <div className="bg-white rounded-lg shadow-md p-6 mb-6 dark:bg-gray-800">
                <h3 className="text-lg font-bold border-b pb-3 mb-4 dark:text-white dark:border-gray-700">最新评论</h3>
                <div className="space-y-4">
                  {Array.isArray(latestComments) && latestComments.length > 0 ? (
                    latestComments.map((comment, index) => (
                      <div key={index} className="flex gap-3">
                        {comment.picture ? (
                          <img
                            src={comment.picture.startsWith('http') ? comment.picture : `${apiHost}${comment.picture}`} 
                            alt={comment.name}
                            width={40}
                            height={40}
                            className="rounded-full"
                          />
                        ) : (
                          <div className="w-10 h-10 rounded-full bg-gray-200 flex items-center justify-center">
                            {comment.name.charAt(0)}
                          </div>
                        )}
                        <div>
                          <div className="flex items-center gap-2">
                            <span className="font-medium dark:text-white">{comment.name}</span>
                            <span className="text-xs text-gray-500 dark:text-gray-400">
                              {new Date(comment.created_at * 1000).toLocaleDateString()}
                            </span>
                          </div>
                          <a 
                            href={`/posts/${comment.post_id}`} 
                            className="text-sm text-gray-700 dark:text-gray-300 hover:text-blue-500"
                          >
                            <p className="line-clamp-2">{comment.content}</p>
                            <p className="text-xs text-gray-500 mt-1">文章：{comment.post_title}</p>
                          </a>
                        </div>
                      </div>
                    ))
                  ) : (
                    <div className="text-center py-4">
                      <p className="text-gray-500 dark:text-gray-400">暂无评论</p>
                    </div>
                  )}
                </div>
              </div>
            </aside>
          )}
        </div>
      </div>
    </div>
  );
}
