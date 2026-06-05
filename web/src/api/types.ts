// 通用 API 响应类型定义
export interface Response<T> {
  code: number;
  message: string;
  data: T;
}
