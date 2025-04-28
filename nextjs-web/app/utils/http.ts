'use client';

import { useState, useEffect } from 'react';

// 基础接口定义
export interface IPageData<T> {
  pageNo: number;
  pageSize: number;
  totalPage: number;
  totalCount: number;
  list: T[];
}

export interface IListData<T> {
  list: T[];
}

export interface IBaseResponse {
  code: number;
  message: string;
}

export interface IResponse<T> {
  code: number;
  data?: T;
  message: string;
}

// 请求方法类型
type Methods = "GET" | "POST" | "DELETE" | "PUT";

// 请求选项接口
export interface RequestOptions extends RequestInit {
  timeout?: number;
  params?: Record<string, any>;
}

// 模拟简单的哈希函数，用于生成唯一的缓存键
function hash(items: any[]): string {
  return items.map(item => 
    typeof item === 'object' ? JSON.stringify(item) : String(item)
  ).join('_');
}

// HTTP请求类
class HttpRequest {
  private baseURL: string;
  
  constructor(baseURL: string = '') {
    this.baseURL = baseURL;
  }
  
  async request<T = any>(
    url: string,
    method: Methods,
    data?: any,
    options?: RequestOptions,
  ): Promise<IResponse<T>> {
    const requestKey = hash(["api-fetch", this.baseURL + url, method, data]);
    const requestOptions: RequestInit = {
      method,
      headers: {
        'Content-Type': 'application/json',
        ...(options?.headers || {})
      },
      ...options,
    };

    // 处理不同请求方法的数据
    if (method === "GET" || method === "DELETE") {
      // 处理查询参数
      const searchParams = new URLSearchParams();
      if (data) {
        Object.entries(data).forEach(([key, value]) => {
          if (value !== undefined && value !== null) {
            searchParams.append(key, String(value));
          }
        });
      }
      
      const queryString = searchParams.toString();
      url = queryString ? `${url}?${queryString}` : url;
    } else if (method === "POST" || method === "PUT") {
      // 处理请求体
      if (data) {
        requestOptions.body = JSON.stringify(data);
      }
    }

    // 完整的请求URL
    const fullURL = `${this.baseURL}${url}`;
    
    try {
      // 执行请求
      const response = await fetch(fullURL, requestOptions);
      
      // 处理非 2xx 响应
      if (!response.ok) {
        throw new Error(`HTTP error: ${response.status} - ${response.statusText}`);
      }
      
      // 解析响应JSON
      const result: IResponse<T> = await response.json();
      
      // 业务逻辑错误处理
      if (result.code !== 0 && result.code !== 200) {
        throw new Error(result.message || '请求失败');
      }
      
      return result;
    } catch (error) {
      console.error('Request error:', error);
      throw error;
    }
  }

  // 封装常用方法

  get<T = any>(url: string, params?: any, options?: RequestOptions) {
    return this.request<T>(url, "GET", params, options);
  }

  post<T = any>(url: string, data?: any, options?: RequestOptions) {
    return this.request<T>(url, "POST", data, options);
  }

  put<T = any>(url: string, data?: any, options?: RequestOptions) {
    return this.request<T>(url, "PUT", data, options);
  }

  delete<T = any>(url: string, params?: any, options?: RequestOptions) {
    return this.request<T>(url, "DELETE", params, options);
  }
}

// 创建一个默认实例
const httpRequest = new HttpRequest('/api');

export default httpRequest;

// React Hook用于数据获取
export function useHttpRequest<T>(
  url: string, 
  method: Methods = 'GET', 
  data?: any, 
  options?: RequestOptions
) {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);
  const [response, setResponse] = useState<IResponse<T> | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        const result = await httpRequest.request<T>(url, method, data, options);
        setResponse(result);
        setError(null);
      } catch (err) {
        setError(err as Error);
        setResponse(null);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [url, method, JSON.stringify(data), JSON.stringify(options)]);

  return { loading, error, response, data: response?.data };
}