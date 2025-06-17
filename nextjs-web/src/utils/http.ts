// src/utils/http.ts

/**
 * 通用 HTTP 请求方法，服务端和客户端均可用。
 * @param input 请求地址或 Request 对象
 * @param init fetch 配置项
 */
export async function request<T = any>(
  input: RequestInfo,
  init?: RequestInit
): Promise<T> {
  const res = await fetch(input, init);
  if (!res.ok) {
    // 可根据需要自定义错误处理
    throw new Error(res.statusText);
  }
  return res.json();
}
