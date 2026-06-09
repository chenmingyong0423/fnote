// src/utils/http.ts

type ApiLikeResponse<T = unknown> = {
  code?: number;
  message?: string;
  data?: T;
};

export class BackendUnavailableError extends Error {
  constructor(message = "后端服务暂不可用") {
    super(message);
    this.name = "BackendUnavailableError";
  }
}

export function isBackendUnavailableError(error: unknown): error is BackendUnavailableError {
  return error instanceof BackendUnavailableError || (
    error instanceof Error && error.name === "BackendUnavailableError"
  );
}

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
  let res: Response;
  try {
    res = await fetch(url, init);
  } catch {
    throw new BackendUnavailableError("无法连接后端服务，请确认后端已经启动。");
  }

  let body: ApiLikeResponse | null = null;
  try {
    body = (await res.json()) as ApiLikeResponse;
  } catch {
    body = null;
  }

  if (!res.ok) {
    const backendMessage =
        body?.message?.trim() || res.statusText || "请求失败";
    throw new Error(backendMessage);
  }

  return (body as T) ?? ({} as T);
}
