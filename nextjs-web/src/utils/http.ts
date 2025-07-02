// src/utils/http.ts

/**
 * 通用 HTTP 请求方法，服务端和客户端均可用。
 * @param input 请求地址或 Request 对象
 * @param init fetch 配置项
 */
export async function request<T = never>(
  input: RequestInfo,
  init?: RequestInit
): Promise<T> {
  let url = input;
  
  // 如果在服务端，且 input 是字符串且以 /api 或 /static 开头，则拼接后端地址
  if (typeof window === 'undefined' && typeof input === 'string' && (input.startsWith('/api') || input.startsWith('/static'))) {
    const serverHost = process.env.SERVER_HOST || 'http://localhost:8080';
    // 去掉 /api 前缀再拼接
    let path = input;
    if (path.startsWith('/api')) {
      path = path.replace(/^\/api/, '');
    }
    url = serverHost.replace(/\/$/, '') + path;
  }
  const res = await fetch(url, init);
  if (!res.ok) {
    // 可根据需要自定义错误处理
    throw new Error(res.statusText);
  }
  return res.json();
}
